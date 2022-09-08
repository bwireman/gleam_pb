package fields

import pgs "github.com/lyft/protoc-gen-star"
import "strings"

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
                name := enum.Name()
                first := enum.Values()[0].Name()
                if ! strings.HasPrefix(first.String(),name.String()) {
                  return (name + first).UpperCamelCase()
                }
 
	return first.UpperCamelCase()

}
