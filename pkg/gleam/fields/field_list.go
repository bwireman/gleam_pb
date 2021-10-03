package fields

import "strings"

type FieldList []*Field

func NewFieldList() FieldList {
	return []*Field{}
}

func (f FieldList) Render() string {
	rendered := []string{}

	for _, field := range f {
		rendered = append(rendered, field.Render())
	}

	return strings.Join(rendered, ", ")
}
