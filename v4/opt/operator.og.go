// Code generated by optgen; DO NOT EDIT.

package opt

const (
	UnknownOp Operator = iota

	SubqueryOp
	VariableOp
	ConstOp
	PlaceholderOp
	ListOp
	OrderedListOp
	FilterListOp
	ProjectionsOp
	ExistsOp
	AndOp
	OrOp
	NotOp
	EqOp
	LtOp
	GtOp
	LeOp
	GeOp
	NeOp
	InOp
	NotInOp
	LikeOp
	NotLikeOp
	ILikeOp
	NotILikeOp
	SimilarToOp
	NotSimilarToOp
	RegMatchOp
	NotRegMatchOp
	RegIMatchOp
	NotRegIMatchOp
	IsDistinctFromOp
	IsNotDistinctFromOp
	IsOp
	IsNotOp
	AnyOp
	SomeOp
	AllOp
	BitandOp
	BitorOp
	BitxorOp
	PlusOp
	MinusOp
	MultOp
	DivOp
	FloorDivOp
	ModOp
	PowOp
	ConcatOp
	LShiftOp
	RShiftOp
	UnaryPlusOp
	UnaryMinusOp
	UnaryComplementOp
	FunctionOp
	TrueOp
	FalseOp
	ScanOp
	ValuesOp
	SelectOp
	ProjectOp
	InnerJoinOp
	LeftJoinOp
	RightJoinOp
	FullJoinOp
	SemiJoinOp
	AntiJoinOp
	InnerJoinApplyOp
	LeftJoinApplyOp
	RightJoinApplyOp
	FullJoinApplyOp
	SemiJoinApplyOp
	AntiJoinApplyOp
	GroupByOp
	UnionOp
	IntersectOp
	ExceptOp
	SortOp
	ArrangeOp
)

const opNames = "unknownsubqueryvariableconstplaceholderlistordered-listfilter-listprojectionsexistsandornoteqltgtlegeneinnot-inlikenot-likei-likenot-i-likesimilar-tonot-similar-toreg-matchnot-reg-matchreg-i-matchnot-reg-i-matchis-distinct-fromis-not-distinct-fromisis-notanysomeallbitandbitorbitxorplusminusmultdivfloor-divmodpowconcatl-shiftr-shiftunary-plusunary-minusunary-complementfunctiontruefalsescanvaluesselectprojectinner-joinleft-joinright-joinfull-joinsemi-joinanti-joininner-join-applyleft-join-applyright-join-applyfull-join-applysemi-join-applyanti-join-applygroup-byunionintersectexceptsortarrange"

var opIndexes = [...]uint32{0, 7, 15, 23, 28, 39, 43, 55, 66, 77, 83, 86, 88, 91, 93, 95, 97, 99, 101, 103, 105, 111, 115, 123, 129, 139, 149, 163, 172, 185, 196, 211, 227, 247, 249, 255, 258, 262, 265, 271, 276, 282, 286, 291, 295, 298, 307, 310, 313, 319, 326, 333, 343, 354, 370, 378, 382, 387, 391, 397, 403, 410, 420, 429, 439, 448, 457, 466, 482, 497, 513, 528, 543, 558, 566, 571, 580, 586, 590, 597}
