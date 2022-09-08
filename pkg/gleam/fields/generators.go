package fields

import (
	"strings"

	pgs "github.com/lyft/protoc-gen-star"
)

type GeneratorFn struct {
	type_name pgs.Name
	types     []*GleamPrimitiveOrValue
}

func (g *GeneratorFn) RenderAsMap() map[string]interface{} {

	fields := []string{}
	for _, t := range g.types {
		if t.Primitive != Unknown {
			fields = append(fields, GleamPrimitiveDefaultValues[t.Primitive])
		} else {
			fields = append(fields, t.Value)
		}

	}

	return map[string]interface{}{
		"type_name": g.type_name.UpperCamelCase().String(),
		"func_name": "new_" + g.type_name.LowerSnakeCase().String(),
		"fields":    strings.Join(fields, ", "),
		"has_fields":    len(fields) > 0,
	}
}

func GeneratorFnFromGleamType(t *GleamType) *GeneratorFn {
	if len(t.Constructors) != 1 {
		return nil
	}

	g := &GeneratorFn{
		type_name: t.TypeName,
		types:     []*GleamPrimitiveOrValue{},
	}

	for _, f := range t.Constructors[0].fields {
		g.types = append(g.types, f.gleam_primitive)
	}

	return g
}
