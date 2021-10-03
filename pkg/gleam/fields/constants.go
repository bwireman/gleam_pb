package fields

import (
	"fmt"

	pgs "github.com/lyft/protoc-gen-star"
)

type GleamPrimitive int

type GleamPrimitiveOrValue struct {
	Primitive GleamPrimitive
	Value     string
}

const (
	Unknown GleamPrimitive = iota
	Int
	Float
	String
	List
	Map
	Option
	Bool
	BitString
	Any
	Duration
	Timestamp
	Value
	Empty
)

var ProtoTypeToPrimitives = map[pgs.ProtoType]GleamPrimitive{
	pgs.DoubleT:  Float,
	pgs.FloatT:   Float,
	pgs.Int64T:   Int,
	pgs.UInt64T:  Int,
	pgs.Int32T:   Int,
	pgs.Fixed64T: Int,
	pgs.Fixed32T: Int,
	pgs.UInt32T:  Int,
	pgs.SFixed32: Int,
	pgs.SFixed64: Int,
	pgs.SInt64:   Int,
	pgs.SInt32:   Int,
	pgs.StringT:  String,
	pgs.BoolT:    Bool,
	pgs.BytesT:   BitString,
}

var GleamPrimitiveDefaultValues = map[GleamPrimitive]string{
	Int:       "0",
	Float:     "0.0",
	String:    "\"\"",
	List:      "list.new()",
	Map:       "map.new()",
	Option:    "option.None",
	Bool:      "false",
	BitString: "<<>>",
	Any:       "Nil",
	Duration:  "option.None",
	Timestamp: "option.None",
	Value:     "Nil",
	Empty:     "#()",
}

var WKTToTypeName = map[pgs.WellKnownType]string{
	pgs.BoolValueWKT:   fmt.Sprintf("option.Option(%s)", ProtoTypeToPrimitives[pgs.BoolT].Render()),
	pgs.StringValueWKT: fmt.Sprintf("option.Option(%s)", ProtoTypeToPrimitives[pgs.StringT].Render()),
	pgs.DoubleValueWKT: fmt.Sprintf("option.Option(%s)", ProtoTypeToPrimitives[pgs.DoubleT].Render()),
	pgs.FloatValueWKT:  fmt.Sprintf("option.Option(%s)", ProtoTypeToPrimitives[pgs.FloatT].Render()),
	pgs.BytesValueWKT:  fmt.Sprintf("option.Option(%s)", ProtoTypeToPrimitives[pgs.BytesT].Render()),
	pgs.Int32ValueWKT:  fmt.Sprintf("option.Option(%s)", ProtoTypeToPrimitives[pgs.Int32T].Render()),
	pgs.Int64ValueWKT:  fmt.Sprintf("option.Option(%s)", ProtoTypeToPrimitives[pgs.Int64T].Render()),
	pgs.UInt32ValueWKT: fmt.Sprintf("option.Option(%s)", ProtoTypeToPrimitives[pgs.UInt32T].Render()),
	pgs.UInt64ValueWKT: fmt.Sprintf("option.Option(%s)", ProtoTypeToPrimitives[pgs.UInt64T].Render()),
	pgs.AnyWKT:         "gleam_pb.Any",
	pgs.DurationWKT:    "option.Option(gleam_pb.Duration)",
	pgs.EmptyWKT:       "gleam_pb.Empty",
	pgs.StructWKT:      "gleam_pb.Struct",
	pgs.TimestampWKT:   "option.Option(gleam_pb.Timestamp)",
	pgs.ValueWKT:       "gleam_pb.Value",
	pgs.ListValueWKT:   "list.List",
}
