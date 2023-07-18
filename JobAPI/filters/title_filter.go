package filters

type TitleFilter struct {
	BaseFilterTemplate
}

func InitTitleFilter(existingFilter IFilter, logicalOperator LogicalOperator, conditionalOperator ConditionalOperator, val any) IFilter {
	return &UsernameFilter{
		BaseFilterTemplate: BaseFilterTemplate{
			existingFilter:      existingFilter,
			val:                 val,
			conditionalOperator: conditionalOperator,
			logicalOperator:     logicalOperator,
			key:                 "title",
		},
	}

}
