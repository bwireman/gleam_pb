package fields

import (
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
		return fmt.Sprintf("%s%s(%s)", format_name(c.type_name), format_name(c.name), c.fields.Render())
	}
	return format_name(c.name)
}

type GleamType struct {
	type_name    pgs.Name
	constructors []*Constructor
}

func (g *GleamType) RenderAsMap() map[string]interface{} {
	cons := []string{}
	for _, con := range g.constructors {
		cons = append(cons, con.Render())
	}

	return map[string]interface{}{
		"type_name":    format_name(g.type_name),
		"constructors": cons,
	}
}

func GleamTypeFromMessage(msg pgs.Message) *GleamType {
	fields := NewFieldList()
	for _, field := range msg.Fields() {
		if !field.InOneOf() {
			fields = append(fields, FieldFromField(field))
		}
	}

	for _, oneof := range msg.OneOfs() {
		fields = append(fields, FieldFromOneOf(msg, oneof))
	}

	return &GleamType{
		type_name: msg.Name(),
		constructors: []*Constructor{
			{
				name:   msg.Name(),
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
		type_name:    containing_message.Name().UpperCamelCase() + oneof.Name().UpperCamelCase(),
		constructors: cons,
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
		type_name:    enum.Name(),
		constructors: cons,
	}
}

func format_name(name pgs.Name) string {
	return name.UpperCamelCase().String()
}
