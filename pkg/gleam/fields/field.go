package fields

import (
	"fmt"
	"strings"

	pgs "github.com/lyft/protoc-gen-star"
)

type Field struct {
	name            pgs.Name
	gleam_primitive *GleamPrimitiveOrValue
	type_name       string
	repeated        bool
	optional        bool
}

func FieldFromField(f pgs.Field) *Field {
	gleam_primitive, type_name, repeated, optional := convert(f)

	return &Field{
		name:            f.Name(),
		type_name:       type_name,
		gleam_primitive: gleam_primitive,
		repeated:        repeated,
		optional:        optional,
	}
}

func FieldFromOneOf(msg pgs.Message, o pgs.OneOf) *Field {
	type_name := msg.Name().UpperCamelCase().String() + o.Name().UpperCamelCase().String()

	return &Field{
		name:            o.Name(),
		gleam_primitive: Option.AsPrimitiveOrValue(),
		type_name:       type_name,
		repeated:        false,
		optional:        true,
	}
}

func (f *Field) Render() string {
	typeName := f.type_name

	// check repeated first because a `List(Option(_))` is dumb
	if f.repeated {
		typeName = fmt.Sprintf("list.List(%s)", typeName)
	} else if f.optional {
		typeName = fmt.Sprintf("option.Option(%s)", typeName)
	}
	return fmt.Sprintf("%s: %s", f.name.LowerSnakeCase(), typeName)
}

func _convert(f pgs.Field) (bool, string) {
	var embed pgs.Message
	if f.Type().IsMap() || f.Type().IsRepeated() {
		embed = f.Type().Element().Embed()
	} else {
		embed = f.Type().Embed()
	}

	if embed.IsWellKnown() {
		return false, WKTToTypeName[embed.WellKnownType()]
	}

	type_package := embed.Package()
	field_package := f.File().Package()
	pName := embed.Name().UpperCamelCase().String()

	if type_package != field_package {
		split := strings.Split(type_package.ProtoName().String(), ".")
		pName = split[len(split)-1] + "." + pName
	}

	return true, pName
}

func convert(f pgs.Field) (*GleamPrimitiveOrValue, string, bool, bool) {
	p_type := f.Type().ProtoType()
	repeated := f.Type().IsRepeated()
	optional := false
	p_name := ""
	gleam_primitive_or_value := Unknown.AsPrimitiveOrValue()

	switch p_type {
	case pgs.EnumT:
		var enum pgs.Enum

		if repeated {
			enum = f.Type().Element().Enum()
		} else {
			enum = f.Type().Enum()
		}

		gleam_primitive_or_value.Value = (enum.Name() + enum.Values()[0].Name()).UpperCamelCase().String()
		p_name = enum.Name().UpperCamelCase().String()

	case pgs.GroupT, pgs.MessageT:
		if f.Type().IsMap() {
			elem_type := ""
			if f.Type().Element().IsEmbed() {
				_, elem_type = _convert(f)
			} else {
				elem_type = ProtoTypeToPrimitives[f.Type().Element().ProtoType()].Render()
			}

			p_name = fmt.Sprintf("map.Map(%s, %s)", ProtoTypeToPrimitives[f.Type().Key().ProtoType()].Render(), elem_type)
			gleam_primitive_or_value.Primitive = Map
		} else {
			optional, p_name = _convert(f)
		}
	default:
		gleam_primitive_or_value.Primitive = ProtoTypeToPrimitives[p_type]
		p_name = gleam_primitive_or_value.Primitive.Render()
	}

	if repeated {
		gleam_primitive_or_value.Primitive = List
	} else if optional {
		gleam_primitive_or_value.Primitive = Option
	}

	return gleam_primitive_or_value, p_name, repeated, optional
}
