package fields

import (
        "strings"
	"fmt"

	pgs "github.com/lyft/protoc-gen-star"
)

type Constructor struct {
	type_name pgs.Name
	name      pgs.Name
	fields    FieldList
}

func (c *Constructor) Render() string {
	if len(c.fields) > 0 {
		return fmt.Sprintf("%s(%s)", format_constructor_name(c.name), c.fields.Render(false))
	}
	return format_constructor_name(c.name)
}

func (c *Constructor) RenderAsGPBTuple() string {
	return fmt.Sprintf("#(atom.Atom, %s)", c.fields.RenderAsGPBTuple(true))
}

func (c *Constructor) RenderAsPatternMatch(overwriteName string, isExtract bool, guard string) string {
	//if len(c.fields) > 0 {
		pattern := c.fields.RenderAsPatternMatch(false, isExtract)
		result := c.fields.RenderAsPatternMatch(true, isExtract)

		if len(overwriteName) == 0 {
			overwriteName = "reserved__struct_name"
		}

		gleam := fmt.Sprintf("%s(%s)", format_constructor_name(c.name), pattern)
		gpb := fmt.Sprintf("#(%s, %s)", overwriteName, result)
		if !isExtract {
			gleam = fmt.Sprintf("%s(%s)", format_constructor_name(c.name), result)
			gpb = fmt.Sprintf("#(%s, %s)", overwriteName, pattern)
		}

                if len(c.fields) == 0 {
                   gleam = fmt.Sprintf("%s", format_constructor_name(c.name))
                }


		mid := " -> "
		if guard != "" {
			mid = guard + mid
		}

		if isExtract {
			return gleam + mid + gpb
		} else {
			return gpb + mid + gleam
		}

	//}
	//return "BROKEN_DEFAULT" + format_constructor_name(c.name)
}

type GleamType struct {
	TypeName     pgs.Name
	Constructors []*Constructor
	IsEnum       bool
}

func (g *GleamType) RenderAsMap() map[string]interface{} {
	cons := []string{}

	for _, con := range g.Constructors {
		cons = append(cons, con.Render())
	}

	return map[string]interface{}{
		"type_name":    format_constructor_name(g.TypeName),
		"constructors": cons,
	}
}

func GleamTypeFromMessage(msg pgs.Message) *GleamType {
	fields := NewFieldList()

	oneOfs := map[string]interface{}{}

	for _, field := range msg.Fields() {
		oo := field.OneOf()
		if !field.InOneOf() {
			fields = append(fields, FieldFromField(field))
		} else if _, ok := oneOfs[oo.FullyQualifiedName()]; !ok && oo != nil {
			oneOfs[oo.FullyQualifiedName()] = nil
			fields = append(fields, FieldFromOneOf(msg, oo))
		}
	}
        
        fixedname := pgs.Name(strings.Replace(msg.Name().String(), "_", "", -1))
	return &GleamType{
		TypeName: fixedname,
		Constructors: []*Constructor{
			{
				name:   fixedname,
				fields: fields,
			},
		},
	}
}

func GleamTypeFromOnoeOf(containing_message pgs.Message, oneof pgs.OneOf) *GleamType {
	cons := []*Constructor{}

	for _, oneof_field := range oneof.Fields() {
		cons = append(cons, &Constructor{
			type_name: containing_message.Name(),
			name:      oneof.Name() + oneof_field.Name(),
			fields:    []*Field{FieldFromField(oneof_field)},
		})
	}

	return &GleamType{
		TypeName:     containing_message.Name().UpperCamelCase() + oneof.Name().UpperCamelCase(),
		Constructors: cons,
	}
}


func GleamTypeFromEnum(enum pgs.Enum) *GleamType {
	cons := []*Constructor{}

	for _, enum_val := range enum.Values() {
		cons = append(cons, &Constructor{
			name:   enum.Name() + enum_val.Name(),
			fields: nil,
		})
	}

	return &GleamType{
		TypeName:     enum.Name(),
		Constructors: cons,
		IsEnum:       true,
	}
}

func format_constructor_name(name pgs.Name) string {
	return name.UpperCamelCase().String()
}
