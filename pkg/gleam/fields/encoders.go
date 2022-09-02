package fields

import (
	"fmt"
	"strings"

	pgs "github.com/lyft/protoc-gen-star"
)

type GenEnc struct {
	type_name            string
	func_name            string
	message_name         string
	constructors         []*Constructor
	extract_patterns     []string
	reconstruct_patterns []string
	reconstruct_vars     []string
	reconstruct_type     string
	is_enum              bool
	is_oneof             bool
}

func GenEncDecFromMessage(msg pgs.Message, gleam_type *GleamType) *GenEnc {
	extract_patterns := []string{}
	reconstruct_patterns := []string{}
	reconstruct_type := []string{"atom.Atom"}

	if len(gleam_type.Constructors) != 1 {
		panic("GenEncDecFromMessage message with len(constructors) != 1")
	}

	c := gleam_type.Constructors[0]

	extract_patterns = append(extract_patterns, c.RenderAsPatternMatch("", true, ""))
	reconstruct_patterns = append(reconstruct_patterns, c.RenderAsPatternMatch("_", false, ""))

	for _, f := range c.fields {
		if f.is_message && f.gt != nil {
                   if f.repeated {

	reconstruct_type = append(reconstruct_type, fmt.Sprintf("List(%s)", f.gt.Constructors[0].RenderAsGPBTuple()))
                      } else {

			reconstruct_type = append(reconstruct_type, fmt.Sprintf("gleam_pb.Undefined(%s)", f.gt.Constructors[0].RenderAsGPBTuple()))
                      }
		} else if f.is_oneof && f.gt != nil {
			reconstruct_type = append(reconstruct_type, fmt.Sprintf("gleam_pb.Undefined(#(atom.Atom, x_%s))", f.name.LowerSnakeCase()))
		} else if f.map_elem_is_message && f.gt != nil {
			reconstruct_type = append(reconstruct_type, fmt.Sprintf("List(#(%s, %s))", f.map_key, f.gt.Constructors[0].RenderAsGPBTuple()))
		} else if f.is_enum {
			reconstruct_type = append(reconstruct_type, "atom.Atom")
		} else {
			reconstruct_type = append(reconstruct_type, f.Render(true))
		}
	}

	return &GenEnc{
		type_name:            msg.Name().UpperCamelCase().String(),
		func_name:            format_func_name(msg.Name()),
		message_name:         format_fqn(msg.FullyQualifiedName()),
		constructors:         gleam_type.Constructors,
		extract_patterns:     extract_patterns,
		reconstruct_patterns: reconstruct_patterns,
		reconstruct_vars:     []string{},
		reconstruct_type:     strings.Join(reconstruct_type, ","),
		is_enum:              false,
		is_oneof:             false,
	}
}

func genOneOfFieldRecFunc(val_name string, type_name pgs.Name, c *Constructor, optional bool) string {

	con_name := format_constructor_name(c.name)

	if optional {
		field_rec_type := fmt.Sprintf("#(atom.Atom, gleam_pb.Undefined(%s))", c.fields.Render(true))

		return fmt.Sprintf(`let rec_%s = fn(v: %s) -> %s { 
			case v {
				#(_, gleam_pb.Undefined) -> %s(option.None)
				#(_, %s) -> %s(option.Some(%s |> gleam_pb.force_a_to_b))
			}}`, val_name, field_rec_type, type_name, con_name, val_name, con_name, val_name)
	}

	return fmt.Sprintf(`let rec_%s = fn(v: %s) -> %s { 
		case v {
			#(_, %s) -> %s(%s)
		}}`, val_name, c.RenderAsGPBTuple(), type_name, val_name, con_name, val_name)
}

func GenEncDecFromOneOf(msg pgs.Message, oo pgs.OneOf, gleam_type *GleamType) *GenEnc {
	extract_patterns := []string{}
	reconstruct_patterns := []string{}
	reconstruct_vars := []string{}

	for i, c := range gleam_type.Constructors {
		if len(c.fields) != 1 {
			panic("GenEncDecFromOneOf len(constructor.field) != 1")
		}
		val_name := oo.Fields()[i].Name().LowerSnakeCase().String()
		atom_val := fmt.Sprintf("atom.create_from_string(\"%s\")", val_name)
		atom_val_name := "atom_" + val_name

		extract_patterns = append(extract_patterns, c.RenderAsPatternMatch(atom_val, true, ""))
		reconstruct_vars = append(reconstruct_vars, fmt.Sprintf("let %s = %s", atom_val_name, atom_val))

		reconstruct_vars = append(reconstruct_vars, genOneOfFieldRecFunc(val_name, gleam_type.TypeName, c, c.fields[0].optional))

		reconstruct_patterns = append(reconstruct_patterns, fmt.Sprintf("#(at, _) if at == %s -> m |> gleam_pb.force_a_to_b |> rec_%s", atom_val_name, val_name))
	}

	return &GenEnc{
		type_name:            format_oneof_type_name(msg, oo).String(),
		func_name:            format_func_name(format_oneof_type_name(msg, oo)),
		message_name:         format_fqn(oo.FullyQualifiedName()),
		constructors:         gleam_type.Constructors,
		extract_patterns:     extract_patterns,
		reconstruct_patterns: reconstruct_patterns,
		reconstruct_vars:     reconstruct_vars,
		reconstruct_type:     "",
		is_enum:              false,
		is_oneof:             true,
	}
}

//{{ range .printers }}
//fn show_{{ .type_name }}(a : {{ .type_name }}) -> String {
//{{ range .params}}  {{.printer_name}}(a.{{.name}}),  {{ end }}
//}
//{{ end }}
func GenPrinterFromMessage(enum pgs.Message, gleam_type *GleamType) map[string]interface{} {
	constructors := []map[string]interface{}{}

	field_package := enum.File().Package()

        split := strings.Split(field_package.ProtoName().String(), ".")
        pkg := pgs.Name(split[len(split)-1]).LowerSnakeCase().String()

	//if len(c.fields) > 0 {
	//	return fmt.Sprintf("%s(%s)", format_constructor_name(c.name), c.fields.Render(false))
	//}
	//return format_constructor_name(c.name)

	for _, f := range gleam_type.Constructors[0].fields {
          printer := ""
          if f.is_message || f.is_enum {
            printer = pgs.Name(strings.ReplaceAll(f.type_name, "_", "")).LowerSnakeCase().String()
            printer = "show_" + printer 
            if f.pkg_name != "" {
            printer = f.pkg_name + "." + printer   
             }
             field_name := f.name.LowerSnakeCase().String()
             if f.repeated {
              printer = `string_builder.to_string(string_builder.from_strings(list.flatten(list.map(a.`+field_name+`, `+printer+`))))`
             } else if f.is_message {
            printer = `case a.`+field_name+` {
              option.None -> "option.None"
              option.Some(n) -> string_builder.to_string(string_builder.from_strings(`+printer+`(n)))
             }` 
           } else {
             printer = printer + "(a." +field_name+ ")"
             }
           } else if f.is_oneof {
             printer = `"unimplemented_one_of"` 
           } else {
            printer = "primitive_show_"+pgs.Name(f.type_name).LowerSnakeCase().String()+"(a."+f.name.LowerSnakeCase().String()+")"
          }
          con := map[string]interface{}{
                "printer":   printer ,
	  }
	  constructors = append(constructors, con)
        }

        foreign_type_name := pkg + "." + gleam_type.Constructors[0].name.String()

	return map[string]interface{}{
		"lowercase_type_name":  enum.Name().UpperCamelCase().LowerSnakeCase().String(),
       		"type_name":            enum.Name().UpperCamelCase().String(),
                "foreign_type_name":    foreign_type_name, 
		"params":               constructors,
	}
}
func GenPrinterFromEnum(enum pgs.Enum, gleam_type *GleamType) map[string]interface{} {
	constructors := []map[string]interface{}{}


	field_package := enum.File().Package()

        split := strings.Split(field_package.ProtoName().String(), ".")
        pkg := pgs.Name(split[len(split)-1]).LowerSnakeCase().String()


	for _, c := range gleam_type.Constructors {
          con := map[string]interface{}{
		"name":    c.Render(),
                "pkg":    pkg,
	  }
		constructors = append(constructors, con)

        }
	return map[string]interface{}{
		"type_name":            enum.Name().LowerSnakeCase().String(),
		"constructors":         constructors,
	}
}

func GenEncDecFromEnum(enum pgs.Enum, gleam_type *GleamType) *GenEnc {
	extract_patterns := []string{}
	reconstruct_patterns := []string{}
	reconstruct_vars := []string{}

	for i, c := range gleam_type.Constructors {
		gpb_enum_name := enum.Values()[i].Name()
                varname := enum.Values()[i].Name().LowerSnakeCase()

		extract_patterns = append(extract_patterns, fmt.Sprintf("%s -> atom.create_from_string(\"%s\")", c.Render(), gpb_enum_name))
		reconstruct_patterns = append(reconstruct_patterns, fmt.Sprintf("x if x == %s -> %s", varname, c.Render()))
		reconstruct_vars = append(reconstruct_vars, fmt.Sprintf("let %s = atom.create_from_string(\"%s\")", varname, gpb_enum_name))
	}

	return &GenEnc{
		type_name:            enum.Name().String(),
		func_name:            format_func_name(enum.Name()),
		message_name:         format_fqn(enum.FullyQualifiedName()),
		constructors:         gleam_type.Constructors,
		extract_patterns:     extract_patterns,
		reconstruct_patterns: reconstruct_patterns,
		reconstruct_vars:     reconstruct_vars,
		reconstruct_type:     "",
		is_enum:              true,
		is_oneof:             false,
	}
}

func (g *GenEnc) RenderAsMap() map[string]interface{} {
	return map[string]interface{}{
		"type_name":            g.type_name,
		"func_name":            g.func_name,
		"message_name":         g.message_name,
		"constructors":         g.constructors,
		"extract_patterns":     g.extract_patterns,
		"reconstruct_patterns": g.reconstruct_patterns,
		"reconstruct_vars":     g.reconstruct_vars,
		"reconstruct_type":     g.reconstruct_type,
		"is_enum":              g.is_enum,
		"is_oneof":             g.is_oneof,
	}
}
