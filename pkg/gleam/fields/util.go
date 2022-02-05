package fields

import pgs "github.com/lyft/protoc-gen-star"

func format_fqn(n string) string {
	return n[1:]
}

func format_func_name(n pgs.Name) string {
	return n.LowerSnakeCase().String()
}

func format_oneof_type_name(msg pgs.Message, o pgs.OneOf) pgs.Name {
	return msg.Name().UpperCamelCase() + o.Name().UpperCamelCase()
}

func format_enum_name(enum pgs.Enum) pgs.Name {
	return (enum.Name() + enum.Values()[0].Name()).UpperCamelCase()
}
