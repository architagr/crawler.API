package filters

import "go.mongodb.org/mongo-driver/bson"

type ILogicalOperatorFilterStatergy interface {
	Build() bson.M
}

type AndLogicalOperatorFilterStatergy struct {
	filter bson.M
}

func (statergy *AndLogicalOperatorFilterStatergy) Build() bson.M {
	return statergy.filter
}

func InitAndLogicalOperatorFilterStatergy(filters ...bson.M) ILogicalOperatorFilterStatergy {
	return &AndLogicalOperatorFilterStatergy{
		filter: bson.M{
			string(AND): filters,
		},
	}
}

type OrLogicalOperatorFilterStatergy struct {
	filter bson.M
}

func (statergy *OrLogicalOperatorFilterStatergy) Build() bson.M {
	return statergy.filter
}

func InitOrLogicalOperatorFilterStatergy(filters ...bson.M) ILogicalOperatorFilterStatergy {
	return &AndLogicalOperatorFilterStatergy{
		filter: bson.M{
			string(OR): filters,
		},
	}
}

func LogicalOperatorFactory(operator LogicalOperator, filter ...bson.M) ILogicalOperatorFilterStatergy {
	switch operator {
	case AND:
		return InitAndLogicalOperatorFilterStatergy(filter...)
	case OR:
		return InitOrLogicalOperatorFilterStatergy(filter...)
	}
	return InitAndLogicalOperatorFilterStatergy(filter...)
}
