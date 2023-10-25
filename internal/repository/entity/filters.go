package entity

import (
	"fmt"
	"strings"
)

type Filter struct {
	Limit      int
	Offset     int
	Conditions []Condition
}

type Condition struct {
	Name  string
	Value string
}

func newCondition(name, value string) Condition {
	return Condition{
		Name:  name,
		Value: value,
	}
}

func NewFilter(offset, limit int, params map[string]string) *Filter {
	filter := &Filter{
		Limit:  limit,
		Offset: offset,
	}

	for key, value := range params {
		filter.Conditions = append(filter.Conditions, newCondition(key, value))
	}

	return filter
}

func (f Filter) Build() string {
	var conditions []string

	for _, condition := range f.Conditions {
		conditions = append(conditions, condition.Build())
	}

	condition := strings.Join(conditions, " AND ")
	if len(condition) != 0 {
		condition = " WHERE " + condition
	}

	return fmt.Sprintf("%s LIMIT %d OFFSET %d", condition, f.Limit, f.Offset)
}

func (f Condition) Build() string {
	return fmt.Sprintf("%s = '%s'", f.Name, f.Value)
}
