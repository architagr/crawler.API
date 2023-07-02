package filters

import "go.mongodb.org/mongo-driver/bson"

type UsernameFilter struct {
	filter bson.M
}

func (filter *UsernameFilter) Build() bson.M {
	return filter.filter
}

func InitUsernameFilter(existingFilter IFilter, logicalOperator LogicalOperator, comparisionOperator ConditionalOperator, val any) IFilter {
	conditionalOperation := bson.M{"username": ConditionalOperatorFactory(comparisionOperator, val).Build()}

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
