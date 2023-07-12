package filters

type EmailFilter struct {
	BaseFilterTemplate
}

func InitEmailFilter(existingFilter IFilter, logicalOperator LogicalOperator, conditionalOperator ConditionalOperator, val any) IFilter {
	return &EmailFilter{
		BaseFilterTemplate: BaseFilterTemplate{
			existingFilter:      existingFilter,
			val:                 val,
			conditionalOperator: conditionalOperator,
			logicalOperator:     logicalOperator,
			key:                 "email",
		},
	}
}
