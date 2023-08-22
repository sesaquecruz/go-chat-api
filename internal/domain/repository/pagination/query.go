package pagination

import (
	"strconv"
	"strings"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
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
		return nil, validation.ErrInvalidQueryPage
	}

	sz, err := strconv.Atoi(size)
	if err != nil || sz < 1 || sz > 50 {
		return nil, validation.ErrInvalidQuerySize
	}

	st := strings.ToUpper(sort)
	if st != "ASC" && st != "DESC" {
		return nil, validation.ErrInvalidQuerySort
	}

	sh := strings.ToUpper(search)
	if len(sh) > 50 {
		return nil, validation.ErrMaxQuerySearch
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
