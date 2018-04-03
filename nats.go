package utils

type Filter map[string]interface{}
type Updates map[string]interface{}
type Sort []string
type Limit int
type Max int

type FindOptions struct {
	Filter Filter
	Sort   Sort
	Limit  Limit
	Max    Max
}

type UpdateOptions struct {
	Filter  Filter
	Updates Updates
}

type DeleteOptions struct {
	Filter Filter
}

type HasOptions struct {
	Filter Filter
}
