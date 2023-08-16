package search

type Page[T any] struct {
	Page  int   `json:"page"`
	Size  int   `json:"size"`
	Total int64 `json:"total"`
	Items []T
}

func NewPage[T any](page int, size int, total int64, items []T) *Page[T] {
	return &Page[T]{
		Page:  page,
		Size:  size,
		Total: total,
		Items: items,
	}
}

func MapPage[T any, K any](page *Page[T], mapper func(T) K) *Page[K] {
	result := &Page[K]{
		Page:  page.Page,
		Size:  page.Size,
		Total: page.Total,
		Items: make([]K, len(page.Items)),
	}

	for i := 0; i < len(page.Items); i++ {
		result.Items[i] = mapper(page.Items[i])
	}

	return result
}
