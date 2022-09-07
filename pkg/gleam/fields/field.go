package fields

import (
	"fmt"
	"strings"

	pgs "github.com/lyft/protoc-gen-star"
)

type Field struct {
	name             pgs.Name
	gleam_primitive  *GleamPrimitiveOrValue
	type_name        string
        type_name_no_pkg string
	atom_name        string
	extract_name     string
	reconstruct_name string
        pkg_name         string
	repeated         bool
	optional         bool
	is_enum          bool
	is_message       bool
	is_oneof         bool
	// list
	list_elem_is_message bool
	//map
	map_elem_is_message bool
	map_key             string
	map_elem            string

	gt *GleamType
}

type conv struct {
        pkg_name            string
	type_name           string
        type_name_no_pkg    string
	gleam_primitive     *GleamPrimitiveOrValue
	atom_name           string
	extract_name        string
	reconstruct_name    string
	repeated            bool
	optional            bool
	is_message          bool
	map_elem_is_message bool
	map_key             string
	map_elem            string
	gt                  *GleamType
}

func FieldFromField(f pgs.Field) *Field {
	c := convert(f)

	return &Field{
		name:                 f.Name(),
                pkg_name:             c.pkg_name, 
		gleam_primitive:      c.gleam_primitive,
		type_name:            c.type_name,
                type_name_no_pkg:      c.type_name_no_pkg, 
		atom_name:            c.atom_name,
		extract_name:         c.extract_name,
		reconstruct_name:     c.reconstruct_name,
		repeated:             c.repeated,
		optional:             c.optional,
		is_enum:              f.Type().IsEnum(),
		is_message:           c.is_message,
		is_oneof:             false,
		map_elem_is_message:  c.map_elem_is_message,
		list_elem_is_message: c.is_message && c.repeated,
		map_key:              c.map_key,
		map_elem:             c.map_elem,
		gt:                   c.gt,
	}
}

func FieldFromOneOf(msg pgs.Message, o pgs.OneOf) *Field {
	type_name := format_oneof_type_name(msg, o)
	func_name := format_func_name(type_name)

	return &Field{
		name:                 o.Name(),
                type_name_no_pkg:     o.Name().String(),
                pkg_name:             "", 
		gleam_primitive:      Option.AsPrimitiveOrValue(),
		type_name:            type_name.String(),
		atom_name:            format_fqn(o.FullyQualifiedName()),
		extract_name:         "extract_" + func_name,
		reconstruct_name:     "reconstruct_" + func_name,
		repeated:             false,
		optional:             true,
		is_enum:              false,
		is_message:           false,
		is_oneof:             true,
		list_elem_is_message: false,
		map_elem_is_message:  false,
		map_key:              "",
		map_elem:             "",
		gt:                   GleamTypeFromOnoeOf(msg, o),
	}
}

func (f *Field) func_name() string {
	return pgs.Name(f.type_name).LowerSnakeCase().String()
}

func (f *Field) Render(asPattern bool) string {
	typeName := f.type_name

	if f.repeated {
		typeName = fmt.Sprintf("List(%s)", typeName)
	} else if f.optional && !asPattern {
		typeName = fmt.Sprintf("option.Option(%s)", typeName)
	}

	if !asPattern {
		return fmt.Sprintf("%s: %s", f.name.LowerSnakeCase(), typeName)
	}

	return typeName
}

func (f *Field) RenderAsGPB(asPattern bool) string {
	typeName := f.type_name

        if f.is_oneof {
		typeName = "gleam_pb.Undefined(#(atom.Atom, x))" 
        } else if f.is_message && f.gt != nil {
		typeName = f.gt.Constructors[0].RenderAsGPBTuple()
        } else if f.is_enum {
		typeName = "atom.Atom" 
        } else if f.repeated {
		typeName = fmt.Sprintf("List(%s)", typeName)
	} else if f.optional && !asPattern {
		typeName = fmt.Sprintf("option.Option(%s)", typeName)
	}

	if !asPattern {
		return fmt.Sprintf("%s: %s", f.name.LowerSnakeCase(), typeName)
	}

	return typeName
}


func (f *Field) RenderAsPatternMatch(rightSide bool, isExtract bool) string {
	v := f.name.LowerSnakeCase().String()

	if rightSide {
		func_name := f.reconstruct_name
		if isExtract {
			func_name = f.extract_name
		}
		atom_name := fmt.Sprintf("atom.create_from_string(\"%s\")", f.atom_name)

		if f.optional && !f.repeated {
			if f.is_message {
				if isExtract {
					v = fmt.Sprintf(`case %s {
							option.Some(x) -> %s(%s, x)
							option.None -> gleam_pb.Undefined |> dynamic.from
						}`, v, func_name, atom_name)
				} else {
					v = fmt.Sprintf(`case %s {
							gleam_pb.Undefined -> option.None
							x -> x |> gleam_pb.force_a_to_b |> %s |> option.Some 
						}`, v, func_name)
				}
			} else {
				if isExtract {
					v = fmt.Sprintf(`case %s {
						option.Some(x) -> %s(x)
						option.None -> gleam_pb.Undefined |> dynamic.from
					}`, v, func_name)
				} else {
					v = fmt.Sprintf("%s(%s)", func_name, v)
				}
			}
		} else if f.is_enum {
			v = fmt.Sprintf("%s(%s)", func_name, v)
		} else if f.map_elem_is_message {
			lambda := func_name
			if isExtract {
				lambda = fmt.Sprintf("fn (x) { %s(%s, x)}", func_name, atom_name)
			}
			v = fmt.Sprintf("list.map(%s, pair.map_second(_, %s))", v, lambda)
		} else if f.list_elem_is_message {
			lambda := func_name
			if isExtract {
				lambda = fmt.Sprintf("fn (x) { %s(%s, x)}", func_name, atom_name)
			}

			v = fmt.Sprintf("list.map(%s, %s)", v, lambda)
		}
	}

	return v
}

func _convert(f pgs.Field) (pName string, pName1 string, atom_name string, extract_name string, reconstruct_name string, gt *GleamType) {
	var embed pgs.Message
	if f.Type().IsMap() || f.Type().IsRepeated() {
		embed = f.Type().Element().Embed()
	} else {
		embed = f.Type().Embed()
	}

	gt = GleamTypeFromMessage(embed)

        fix_fname := pgs.Name(strings.Replace(embed.Name().String(), "_", "", -1))
	embed_func_name := format_func_name(fix_fname)
	extract_name = "extract_" + embed_func_name
	reconstruct_name = "reconstruct_" + embed_func_name
	atom_name = format_fqn(embed.FullyQualifiedName())

	type_package := embed.Package()
	field_package := f.File().Package()
	pName = strings.Replace(embed.Name().UpperCamelCase().String(), "_", "", -1) 
        pName1 = pName 
	if type_package != field_package {
		split := strings.Split(type_package.ProtoName().String(), ".")
		pkg := pgs.Name(split[len(split)-1]).LowerSnakeCase().String()

		extract_name = pkg + ".extract_" + embed.Name().LowerSnakeCase().String()
		reconstruct_name = pkg + ".reconstruct_" + embed.Name().LowerSnakeCase().String()

		pName = pkg + "." + pName
	}

	return pName, pName1, atom_name, extract_name, reconstruct_name, gt
}

func convert(f pgs.Field) *conv {
	cv := &conv{}
	cv.repeated = f.Type().IsRepeated()
	cv.extract_name = "extract_"
	cv.reconstruct_name = "reconstruct_"

	gleam_primitive_or_value := Unknown.AsPrimitiveOrValue()
	p_type := f.Type().ProtoType()
	 
        if f.Imports() != nil {
          new_import := strings.ReplaceAll(
                                                  f.Imports()[0].Package().ProtoName().LowerSnakeCase().String(), ".", "/")
          cv.pkg_name = new_import


        }
 
	switch p_type {
	case pgs.EnumT:


		var enum pgs.Enum

		if cv.repeated {
			enum = f.Type().Element().Enum()
		} else {
			enum = f.Type().Enum()
		}
		cv.atom_name = format_fqn(enum.FullyQualifiedName())
		gleam_primitive_or_value.Value = format_enum_name(enum).String()

                cv.type_name_no_pkg = enum.Name().UpperCamelCase().String()
                cv.type_name = enum.Name().UpperCamelCase().String()
                cv.extract_name = "extract_" + pgs.Name(cv.type_name).LowerSnakeCase().String()
		cv.reconstruct_name = "reconstruct_" + pgs.Name(cv.type_name).LowerSnakeCase().String()
 
                if f.Imports() != nil {
		
                  pkg := strings.ReplaceAll(
                                                  f.Imports()[0].Package().ProtoName().LowerSnakeCase().String(), ".", "/")

                  gleam_primitive_or_value.Value = pkg + "." + gleam_primitive_or_value.Value
                  cv.extract_name = pkg + ".extract_" + pgs.Name(cv.type_name).LowerSnakeCase().String()
		  cv.reconstruct_name = pkg + ".reconstruct_" + pgs.Name(cv.type_name).LowerSnakeCase().String()
                  cv.type_name = pkg + "." + cv.type_name
 


                }

		

	case pgs.GroupT, pgs.MessageT:

		if f.Type().IsMap() {
			if f.Type().Element().IsEmbed() {
				cv.map_elem_is_message = true
				cv.map_elem, cv.type_name_no_pkg, cv.atom_name, cv.extract_name, cv.reconstruct_name, cv.gt = _convert(f)
			} else {
				cv.map_elem = ProtoTypeToPrimitives[f.Type().Element().ProtoType()].Render()
			}

			cv.map_key = ProtoTypeToPrimitives[f.Type().Key().ProtoType()].Render()
			cv.type_name = fmt.Sprintf("List(#(%s, %s))", cv.map_key, cv.map_elem)
			gleam_primitive_or_value.Primitive = Map
		} else {
			cv.optional = true
			cv.is_message = true
			cv.type_name, cv.type_name_no_pkg, cv.atom_name, cv.extract_name, cv.reconstruct_name, cv.gt = _convert(f)
			gleam_primitive_or_value.Primitive = Option
		}
	default:

		gleam_primitive_or_value.Primitive = ProtoTypeToPrimitives[p_type]
		cv.type_name = gleam_primitive_or_value.Primitive.Render()
                if f.Imports() != nil {
                  new_import := strings.ReplaceAll(
                                                  f.Imports()[0].Package().ProtoName().LowerSnakeCase().String(), ".", "/")

                  cv.type_name = new_import  + "."  + cv.type_name
                }

	}

	if cv.repeated {
		gleam_primitive_or_value.Primitive = List
	} else if cv.optional {
		gleam_primitive_or_value.Primitive = Option
	}

	cv.gleam_primitive = gleam_primitive_or_value


        
        //if f.Imports() != nil {
	//  new_import := strings.ReplaceAll(
        //                                  f.Imports()[0].Package().ProtoName().LowerSnakeCase().String(), ".", "/")

        //  cv.type_name = new_import  + "."  + cv.type_name
        //}

	return cv
}
