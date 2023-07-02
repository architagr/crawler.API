package filters

import "go.mongodb.org/mongo-driver/bson"

type IdFilter struct {
	filter bson.M
}

func (filter *IdFilter) Build() bson.M {
	return filter.filter
}

func InitIdFilter(existingFilter IFilter, logicalOperator LogicalOperator, comparisionOperator ConditionalOperator, val any) IFilter {
	conditionalOperation := bson.M{"_id": ConditionalOperatorFactory(comparisionOperator, val).Build()}

	if existingFilter == nil {
		return &EmailFilter{
			filter: conditionalOperation,
		}
	}

	logicalOperation := LogicalOperatorFactory(logicalOperator, existingFilter.Build(), conditionalOperation)

	return &EmailFilter{
		filter: logicalOperation.Build(),
	}
}
