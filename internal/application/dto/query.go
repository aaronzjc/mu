package dto

type Query struct {
	Query string
	Args  []interface{}
	Order string
	Limit int
}
