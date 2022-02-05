package fields

func (p GleamPrimitive) AsPrimitiveOrValue() *GleamPrimitiveOrValue {
	return &GleamPrimitiveOrValue{
		Primitive: p,
		Value:     "",
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
		return "List"
	case Map:
		return "List"
	case Option:
		return "option.Option"
	case Bool:
		return "Bool"
	case BitString:
		return "BitString"
	default:
		return "Nil"
	}
}
