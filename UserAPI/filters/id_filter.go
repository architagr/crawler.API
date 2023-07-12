package filters

type IdFilter struct {
	BaseFilterTemplate
}

func InitIdFilter(existingFilter IFilter, logicalOperator LogicalOperator, conditionalOperator ConditionalOperator, val any) IFilter {
	return &IdFilter{
		BaseFilterTemplate: BaseFilterTemplate{
			existingFilter:      existingFilter,
			val:                 val,
			conditionalOperator: conditionalOperator,
			logicalOperator:     logicalOperator,
			key:                 "_id",
		},
	}
}
