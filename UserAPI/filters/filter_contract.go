package filters

import "go.mongodb.org/mongo-driver/bson"

type IFilter interface {
	Build() bson.M
}

type LogicalOperator string
type ConditionalOperator string

const (
	AND LogicalOperator = "$and"
	OR  LogicalOperator = "$or"
)

const (
	NOT_EQUAL ConditionalOperator = "$ne"
	EQUAL     ConditionalOperator = "$eq"
)
