package model

import (
	"fmt"
	"io"
	"strconv"
)

type Mutation struct {
}

type Query struct {
}

type Subscription struct {
}

type SortBy string

const (
	SortByNewest SortBy = "NEWEST"
	SortByOldest SortBy = "OLDEST"
	SortByTop    SortBy = "TOP"
)

var AllSortBy = []SortBy{
	SortByNewest,
	SortByOldest,
	SortByTop,
}

func (e SortBy) IsValid() bool {
	switch e {
	case SortByNewest, SortByOldest, SortByTop:
		return true
	}
	return false
}

func (e SortBy) String() string {
	return string(e)
}

func (e *SortBy) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortBy(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortBy", str)
	}
	return nil
}

func (e SortBy) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
