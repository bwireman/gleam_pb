package fields

import "strings"
import "fmt"

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

func (f FieldList) RenderAsGPBTuple(asPattern bool) string {
	rendered := []string{}

	for _, field := range f {
            if field.is_message && field.gt != nil {
              t := "" 
              if field.repeated {
                t = fmt.Sprintf("List(%s)",  field.gt.Constructors[0].RenderAsGPBTuple())
              } else {
                t = fmt.Sprintf("gleam_pb.Undefined(%s)", field.gt.Constructors[0].RenderAsGPBTuple())
              }
       	      rendered = append(rendered, t)
            } else {
	      rendered = append(rendered, field.RenderAsGPB(asPattern))
            }
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
