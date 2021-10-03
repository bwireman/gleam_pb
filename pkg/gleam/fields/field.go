package fields

import (
	"fmt"
	"strings"

	pgs "github.com/lyft/protoc-gen-star"
)

type Field struct {
	name      pgs.Name
	type_name string
	repeated  bool
	optional  bool
}

func FieldFromField(f pgs.Field) *Field {
	typeName, repeated, optional := convert(f)

	return &Field{
		name:      f.Name(),
		type_name: typeName,
		repeated:  repeated,
		optional:  optional,
	}
}

func FieldFromOneOf(msg pgs.Message, o pgs.OneOf) *Field {
	return &Field{
		name:      o.Name(),
		type_name: (msg.Name() + o.Name()).UpperCamelCase().String(),
		repeated:  false,
		optional:  false,
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

var primitives = map[pgs.ProtoType]string{
	pgs.DoubleT:  "Float",
	pgs.FloatT:   "Float",
	pgs.Int64T:   "Int",
	pgs.UInt64T:  "Int",
	pgs.Int32T:   "Int",
	pgs.Fixed64T: "Int",
	pgs.Fixed32T: "Int",
	pgs.UInt32T:  "Int",
	pgs.SFixed32: "Int",
	pgs.SFixed64: "Int",
	pgs.SInt64:   "Int",
	pgs.SInt32:   "Int",
	pgs.StringT:  "String",
	pgs.BoolT:    "Bool",
	pgs.BytesT:   "BitString",
}

var wktToTypeName = map[pgs.WellKnownType]string{
	pgs.BoolValueWKT:   fmt.Sprintf("option.Option(%s)", primitives[pgs.BoolT]),
	pgs.StringValueWKT: fmt.Sprintf("option.Option(%s)", primitives[pgs.StringT]),
	pgs.DoubleValueWKT: fmt.Sprintf("option.Option(%s)", primitives[pgs.DoubleT]),
	pgs.FloatValueWKT:  fmt.Sprintf("option.Option(%s)", primitives[pgs.FloatT]),
	pgs.BytesValueWKT:  fmt.Sprintf("option.Option(%s)", primitives[pgs.BytesT]),
	pgs.Int32ValueWKT:  fmt.Sprintf("option.Option(%s)", primitives[pgs.Int32T]),
	pgs.Int64ValueWKT:  fmt.Sprintf("option.Option(%s)", primitives[pgs.Int64T]),
	pgs.UInt32ValueWKT: fmt.Sprintf("option.Option(%s)", primitives[pgs.UInt32T]),
	pgs.UInt64ValueWKT: fmt.Sprintf("option.Option(%s)", primitives[pgs.UInt64T]),
	pgs.AnyWKT:         "gleam_pb.Any",
	pgs.DurationWKT:    "gleam_pb.Duration",
	pgs.EmptyWKT:       "gleam_pb.Empty",
	pgs.StructWKT:      "gleam_pb.Struct",
	pgs.TimestampWKT:   "gleam_pb.Timestamp",
	pgs.ValueWKT:       "gleam_pb.Value",
	pgs.ListValueWKT:   "list.List",
}

func _convert(f pgs.Field) (bool, string) {
	var embed pgs.Message
	if f.Type().IsMap() || f.Type().IsRepeated() {
		embed = f.Type().Element().Embed()
	} else {
		embed = f.Type().Embed()
	}

	if embed.IsWellKnown() {
		return false, wktToTypeName[embed.WellKnownType()]
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

func convert(f pgs.Field) (string, bool, bool) {
	pType := f.Type().ProtoType()
	repeated := f.Type().IsRepeated()
	optional := false
	pName := ""

	switch pType {
	case pgs.EnumT:
		var enum pgs.Enum

		if f.Type().IsRepeated() {
			enum = f.Type().Element().Enum()
		} else {
			enum = f.Type().Enum()
		}

		pName = enum.Name().UpperCamelCase().String()
	case pgs.GroupT, pgs.MessageT:
		if f.Type().IsMap() {
			elem_type := ""
			if f.Type().Element().IsEmbed() {
				_, elem_type = _convert(f)
			} else {
				elem_type = primitives[f.Type().Element().ProtoType()]
			}

			pName = fmt.Sprintf("map.Map(%s, %s)", primitives[f.Type().Key().ProtoType()], elem_type)
		} else {
			optional, pName = _convert(f)
		}
	default:
		pName = primitives[pType]
	}

	return pName, repeated, optional

}
