package filters

import "go.mongodb.org/mongo-driver/bson"

type IConditionalOperatorFilterStatergy interface {
	Build() bson.M
}

type EqualConditionalOperatorFilterStatergy struct {
	filter bson.M
}

func (statergy *EqualConditionalOperatorFilterStatergy) Build() bson.M {
	return statergy.filter
}

func InitEqualConditionalOperatorFilterStatergy(obj any) IConditionalOperatorFilterStatergy {
	return &EqualConditionalOperatorFilterStatergy{
		filter: bson.M{string(EQUAL): obj},
	}
}

type NotEqualConditionalOperatorFilterStatergy struct {
	filter bson.M
}

func (statergy *NotEqualConditionalOperatorFilterStatergy) Build() bson.M {
	return statergy.filter
}

func InitNotEqualConditionalOperatorFilterStatergy(obj any) IConditionalOperatorFilterStatergy {
	return &EqualConditionalOperatorFilterStatergy{
		filter: bson.M{string(NOT_EQUAL): obj},
	}
}

func ConditionalOperatorFactory(operator ConditionalOperator, obj any) ILogicalOperatorFilterStatergy {
	switch operator {
	case EQUAL:
		return InitEqualConditionalOperatorFilterStatergy(obj)
	case NOT_EQUAL:
		return InitNotEqualConditionalOperatorFilterStatergy(obj)
	}
	return InitEqualConditionalOperatorFilterStatergy(obj)
}
