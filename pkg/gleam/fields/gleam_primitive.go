package fields

func (p GleamPrimitive) AsPrimitiveOrValue() *GleamPrimitiveOrValue {
	return &GleamPrimitiveOrValue{
		Primitive: p,
		Value:    "",
	}
}

func (p GleamPrimitive) Render() string {
	switch p {
	case Int:
		return "Int"
	case Float:
		return "Float"
	case String:
		return "String"
	case List:
		return "list.List"
	case Map:
		return "map.Map"
	case Option:
		return "option.Option"
	case Bool:
		return "Bool"
	case BitString:
		return "BitString"
	case Any:
		return "gleam_pb.Any"
	case Duration:
		return "gleam_pb.Duration"
	case Timestamp:
		return "gleam_pb.Timestamp"
	case Value:
		return "gleam_pb.Value"
	case Empty:
		return "gleam_pb.Empty"
	default:
		return "Nil"
	}
}
