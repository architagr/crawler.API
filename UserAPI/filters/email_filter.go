package filters

import "go.mongodb.org/mongo-driver/bson"

type EmailFilter struct {
	filter bson.M
}

func (filter *EmailFilter) Build() bson.M {
	return filter.filter
}

func InitEmailFilter(existingFilter IFilter, logicalOperator LogicalOperator, comparisionOperator ConditionalOperator, val any) IFilter {
	conditionalOperation := bson.M{"email": ConditionalOperatorFactory(comparisionOperator, val).Build()}

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
