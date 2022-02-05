package fields

import "strings"

type FieldList []*Field

func NewFieldList() FieldList {
	return []*Field{}
}

func (f FieldList) Render(asPattern bool) string {
	rendered := []string{}

	for _, field := range f {
		rendered = append(rendered, field.Render(asPattern))
	}

	return strings.Join(rendered, ", ")
}

func (f FieldList) RenderAsPatternMatch(rightSide bool, isExtract bool) string {
	rendered := []string{}

	for _, field := range f {
		rendered = append(rendered, field.RenderAsPatternMatch(rightSide, isExtract))
	}

	return strings.Join(rendered, ", ")
}
