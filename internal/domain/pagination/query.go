package pagination

import (
	"strconv"
	"strings"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
)

const (
	ErrInvalidQueryPage   = validation.ValidationError("query 'page' must be greater than or equal to 0")
	ErrInvalidQuerySize   = validation.ValidationError("query 'size' must be greater than 1 and less than or equal to 50")
	ErrInvalidQuerySort   = validation.ValidationError("query 'sort' must be 'asc' or 'desc'")
	ErrInvalidQuerySearch = validation.ValidationError("query 'search' length must be less than or equal to 50")
)

type Query struct {
	page   int
	size   int
	sort   string
	search string
}

func NewQuery(page, size, sort, search string) (*Query, error) {
	if page == "" {
		page = "0"
	}

	if size == "" {
		size = "10"
	}

	if sort == "" {
		sort = "asc"
	}

	pg, err := strconv.Atoi(page)
	if err != nil || pg < 0 {
		return nil, ErrInvalidQueryPage
	}

	sz, err := strconv.Atoi(size)
	if err != nil || sz < 1 || sz > 50 {
		return nil, ErrInvalidQuerySize
	}

	st := strings.ToUpper(sort)
	if st != "ASC" && st != "DESC" {
		return nil, ErrInvalidQuerySort
	}

	sh := strings.ToUpper(search)
	if len(sh) > 50 {
		return nil, ErrInvalidQuerySearch
	}

	return &Query{
		page:   pg,
		size:   sz,
		sort:   st,
		search: sh,
	}, nil
}

func (q *Query) Page() int {
	return q.page
}

func (q *Query) Size() int {
	return q.size
}

func (q *Query) Sort() string {
	return q.sort
}

func (q *Query) Search() string {
	return q.search
}
