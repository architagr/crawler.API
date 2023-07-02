package filters

import "go.mongodb.org/mongo-driver/bson"

type PhoneFilter struct {
	filter bson.M
}

func (filter *PhoneFilter) Build() bson.M {
	return filter.filter
}

func InitPhoneFilter(existingFilter IFilter, logicalOperator LogicalOperator, comparisionOperator ConditionalOperator, val any) IFilter {
	conditionalOperation := bson.M{"phone": ConditionalOperatorFactory(comparisionOperator, val).Build()}

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
