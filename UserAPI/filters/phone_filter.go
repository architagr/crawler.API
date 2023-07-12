package filters

type PhoneFilter struct {
	BaseFilterTemplate
}

func InitPhoneFilter(existingFilter IFilter, logicalOperator LogicalOperator, conditionalOperator ConditionalOperator, val any) IFilter {
	return &PhoneFilter{
		BaseFilterTemplate: BaseFilterTemplate{
			existingFilter:      existingFilter,
			val:                 val,
			conditionalOperator: conditionalOperator,
			logicalOperator:     logicalOperator,
			key:                 "phone",
		},
	}

}
