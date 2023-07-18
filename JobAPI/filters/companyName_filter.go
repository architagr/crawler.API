package filters

type CompanynameFilter struct {
	BaseFilterTemplate
}

func InitCompanynameFilter(existingFilter IFilter, logicalOperator LogicalOperator, conditionalOperator ConditionalOperator, val any) IFilter {
	return &UsernameFilter{
		BaseFilterTemplate: BaseFilterTemplate{
			existingFilter:      existingFilter,
			val:                 val,
			conditionalOperator: conditionalOperator,
			logicalOperator:     logicalOperator,
			key:                 "companyName",
		},
	}

}
