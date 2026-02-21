package models

// Summary holds the aggregated data for the dashboard
type Summary struct {
	TotalIncome  int64 `json:"total_income"`
	TotalExpense int64 `json:"total_expense"`
	NetBalance   int64 `json:"net_balance"`
}
