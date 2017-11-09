package v3

import (
	"bytes"
	"fmt"
	"math/bits"

	"github.com/cockroachdb/cockroach/pkg/sql/parser"
)

// queryState holds per-query state such as the tables referenced by the query
// and the mapping from table name to the column index for those tables columns
// within the query.
type queryState struct {
	catalog map[string]*table
	tables  map[string]bitmapIndex
	nextVar bitmapIndex
}

type columnProps struct {
	name   string
	tables []string
	index  bitmapIndex
	// TODO(peter): value constraints.
}

func (c columnProps) hasColumn(tableName, colName string) bool {
	if colName != c.name {
		return false
	}
	if tableName == "" {
		return true
	}
	return c.hasTable(tableName)
}

func (c columnProps) hasTable(tableName string) bool {
	for _, t := range c.tables {
		if t == tableName {
			return true
		}
	}
	return false
}

func (c columnProps) newVariableExpr(tableName string, props *logicalProps) *expr {
	if tableName == "" {
		if len(c.tables) > 0 {
			tableName = c.tables[0]
		}
	}
	col := &parser.ColumnItem{
		TableName: parser.TableName{
			TableName:               parser.Name(tableName),
			DBNameOriginallyOmitted: true,
		},
		ColumnName: parser.Name(c.name),
	}

	e := &expr{
		op:      variableOp,
		private: col.String(),
	}
	e.inputVars.set(c.index)
	e.updateProps()
	return e
}

type foreignKeyProps struct {
	src  bitmap
	dest bitmap
}

type logicalProps struct {
	columns []columnProps

	// Bitmap indicating which output columns cannot be NULL. The NULL-ability of
	// columns flows from the inputs and can also be derived from filters that
	// are NULL-intolerant.
	notNullCols bitmap

	// A column set is a key if no two rows are equal after projection onto that
	// set. A requirement for a column set to be a key is for no columns in the
	// set to be NULL-able. This requirement stems from the property of NULL
	// where NULL != NULL. The simplest example of a key is the primary key for a
	// table (recall that all of the columns of the primary key are defined to be
	// NOT NULL).
	//
	// A weak key is a set of columns where no two rows containing non-NULL
	// values are equal after projection onto that set. A UNIQUE index on a table
	// is a weak key and possibly a key if all of the columns are NOT NULL. A
	// weak key is a key if "(weakKeys[i] & notNullColumns) == weakKeys[i]".
	weakKeys []bitmap

	// A foreign key is a set of columns in the source table that uniquely
	// identify a single row in the destination table. A foreign key thus refers
	// to a primary key or unique key in the destination table. If the source
	// columns are NOT NULL a foreign key can prove the existence of a row in the
	// destination table and can also be used to infer the cardinality of joins
	// when joining on the foreign key. Consider the schema:
	//
	//   CREATE TABLE departments (
	//     dept_id INT PRIMARY KEY,
	//     name STRING
	//   );
	//
	//   CREATE TABLE employees (
	//     emp_id INT PRIMARY KEY,
	//     dept_id INT REFERENCES d (dept_id),
	//     name STRING,
	//     salary INT
	//   );
	//
	// And the query:
	//
	//   SELECT e.name, e.salary
	//   FROM employees e, departments d
	//   WHERE e.dept_id = d.dept_id
	//
	// The foreign key constraint specifies that employees.dept_id must match a
	// value in departments.dept_id or be NULL. Because departments.dept_id is NOT
	// NULL (due to being part of the primary key), we know the only rows from
	// employees that will not be in the join are those with a NULL dept_id. So we
	// can transform the query into:
	//
	//   SELECT e.name, e.salary
	//   FROM employees e
	//   WHERE e.dept_id IS NOT NULL
	foreignKeys []foreignKeyProps
}

func (p *logicalProps) String() string {
	var buf bytes.Buffer
	p.format(&buf, 0)
	return buf.String()
}

func (p *logicalProps) format(buf *bytes.Buffer, level int) {
	indent := spaces[:2*level]
	fmt.Fprintf(buf, "%scolumns:", indent)
	for _, col := range p.columns {
		buf.WriteString(" ")
		if tables := col.tables; len(tables) > 1 {
			buf.WriteString("{")
			for j, table := range tables {
				if j > 0 {
					buf.WriteString(",")
				}
				buf.WriteString(table)
			}
			buf.WriteString("}")
		} else if len(tables) == 1 {
			buf.WriteString(tables[0])
		}
		buf.WriteString(".")
		buf.WriteString(col.name)
		buf.WriteString(":")
		fmt.Fprintf(buf, "%d", col.index)
		if p.notNullCols.get(col.index) {
			buf.WriteString("*")
		}
	}
	buf.WriteString("\n")
	for _, key := range p.weakKeys {
		var prefix string
		if (key & p.notNullCols) != key {
			prefix = "weak "
		}
		fmt.Fprintf(buf, "%s%skey: %s\n", indent, prefix, key)
	}
	for _, fkey := range p.foreignKeys {
		fmt.Fprintf(buf, "%sforeign key: %s -> %s\n", indent, fkey.src, fkey.dest)
	}
}

// fingerprint returns a string which uniquely identifies the logical
// properties within the context of a query.
func (p *logicalProps) fingerprint() string {
	// TODO(peter): The fingerprint is unique, but human readable. A binary
	// format encoding columns and keys using varints might be faster and more
	// compact.
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, col := range p.columns {
		if i > 0 {
			buf.WriteString(" ")
		}
		fmt.Fprintf(&buf, "%d", col.index)
		if p.notNullCols.get(col.index) {
			buf.WriteString("*")
		}
	}
	buf.WriteString("]")
	if len(p.weakKeys) > 0 {
		buf.WriteString(" [")
		for i, key := range p.weakKeys {
			if i > 0 {
				buf.WriteString(" ")
			}
			fmt.Fprintf(&buf, "%s", key)
		}
		buf.WriteString("]")
	}
	if len(p.foreignKeys) > 0 {
		buf.WriteString(" [")
		for i, fkey := range p.foreignKeys {
			if i > 0 {
				buf.WriteString(" ")
			}
			fmt.Fprintf(&buf, "%s->%s", fkey.src, fkey.dest)
		}
		buf.WriteString("]")
	}
	return buf.String()
}

func (p *logicalProps) newColumnExpr(name string) *expr {
	for _, col := range p.columns {
		if col.name == name {
			return col.newVariableExpr(col.tables[0], p)
		}
	}
	return nil
}

func (p *logicalProps) newColumnExprByIndex(index bitmapIndex) *expr {
	for _, col := range p.columns {
		if col.index == index {
			return col.newVariableExpr(col.tables[0], p)
		}
	}
	return nil
}

// Add additional not-NULL columns based on the filtering expressions.
func (p *logicalProps) applyFilters(filters []*expr) {
	for _, filter := range filters {
		// TODO(peter): !isNullTolerant(filter)
		for v := filter.inputVars; v != 0; {
			i := bitmapIndex(bits.TrailingZeros64(uint64(v)))
			v.clear(i)
			p.notNullCols.set(i)
		}
	}
}

// A filter is compatible with the logical properties for an expression if all
// of the input variables used by the filter are provided by the columns.
func (p *logicalProps) isFilterCompatible(filter *expr) bool {
	v := filter.inputVars
	for i := 0; v != 0 && i < len(p.columns); i++ {
		v.clear(p.columns[i].index)
	}
	return v == 0
}

func (p *logicalProps) outputVars() bitmap {
	var b bitmap
	for _, col := range p.columns {
		b.set(col.index)
	}
	return b
}

func updateProps(e *expr) {
	for _, input := range e.inputs() {
		updateProps(input)
	}
	e.updateProps()
}

func initKeys(e *expr, state *queryState) {
	for _, input := range e.inputs() {
		initKeys(input, state)
	}
	e.initKeys(state)
}
