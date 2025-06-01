package pagination
type Pagination[T any] struct {
	Page int	`json:"page"`
	LastPage int `json:"last_page"`
	Count	int	`json:"count"`
	Items []T	`json:"items"`
}