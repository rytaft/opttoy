package v3

import (
	"bytes"
)

func init() {
	registerOperator(groupByOp, "groupBy", groupBy{})
}

type groupBy struct{}

func (groupBy) kind() operatorKind {
	return relationalKind
}

func (groupBy) format(e *expr, buf *bytes.Buffer, level int) {
	formatRelational(e, buf, level)
	formatExprs(buf, "groupings", e.groupings(), level)
	formatExprs(buf, "aggregations", e.aggregations(), level)
	formatExprs(buf, "filters", e.filters(), level)
	formatExprs(buf, "inputs", e.inputs(), level)
}

func (groupBy) initKeys(e *expr, state *queryState) {
}

func (g groupBy) updateProps(e *expr) {
	e.inputVars = g.requiredInputVars(e)
	e.inputVars &^= (e.props.outputVars() | e.providedInputVars())
	for _, input := range e.inputs() {
		e.inputVars |= input.inputVars
	}

	e.props.applyFilters(e.filters())

	// TODO(peter): update keys
}

func (groupBy) requiredInputVars(e *expr) bitmap {
	var v bitmap
	for _, filter := range e.filters() {
		v |= filter.inputVars
	}
	for _, aggregate := range e.aggregations() {
		v |= aggregate.inputVars
	}
	for _, grouping := range e.groupings() {
		v |= grouping.inputVars
	}
	return v
}

func (groupBy) equal(a, b *expr) bool {
	aAggregations, bAggregations := a.aggregations(), b.aggregations()
	if len(aAggregations) != len(bAggregations) {
		return false
	}
	for i := range aAggregations {
		if !aAggregations[i].equal(bAggregations[i]) {
			return false
		}
	}

	aGroupings, bGroupings := a.groupings(), b.groupings()
	if len(aGroupings) != len(bGroupings) {
		return false
	}
	for i := range aGroupings {
		if !aGroupings[i].equal(bAggregations[i]) {
			return false
		}
	}
	return true
}
