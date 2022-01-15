package gleam

import (
	"strings"
	"text/template"

	"github.com/bwireman/gleam_pb/pkg/gleam/fields"
	pgs "github.com/lyft/protoc-gen-star"
)

type GleamModule struct {
	*pgs.ModuleBase
	tpl *template.Template
}

func Gleam() *GleamModule { return &GleamModule{ModuleBase: &pgs.ModuleBase{}} }

func (g *GleamModule) InitContext(c pgs.BuildContext) {
	g.ModuleBase.InitContext(c)

	g.tpl = template.Must(template.New("gleam-package-template").Parse(fields.GleamTemplate))
}

func (g *GleamModule) Name() string { return "gleam" }

func (g *GleamModule) Execute(_targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	g.AddCustomFile(g.OutputPath()+"/src/gleam_pb.gleam", fields.GleamPB, 0644)

	for _, p := range pkgs {
		allMessages := []pgs.Message{}
		allEnums := []pgs.Enum{}
		imports := []string{"gleam/option", "gleam/list", "gleam/pair", "gleam/dynamic", "gleam/erlang/atom", "gleam_pb"}

		for _, file := range p.Files() {
			allMessages = append(allMessages, file.AllMessages()...)
			allEnums = append(allEnums, file.AllEnums()...)

			for _, imp := range file.Imports() {
				if imp.Package().ProtoName() != p.ProtoName() {
					new_import := strings.ReplaceAll(imp.Package().ProtoName().String(), ".", "/")

					imports = append(imports, new_import)
				}

			}
		}
		g.generate(allMessages, allEnums, p, imports)
	}

	return g.Artifacts()
}

func (g *GleamModule) generate(all_messages []pgs.Message, all_enums []pgs.Enum, pkg pgs.Package, imports []string) {
	gleam_types_map := []map[string]interface{}{}
	enc_dec := []map[string]interface{}{}
	generators := []map[string]interface{}{}

	for _, enum := range all_enums {
		// don't need enum generators
		enum_gleam_type := fields.GleamTypeFromEnum(enum)
		gleam_types_map = append(gleam_types_map, enum_gleam_type.RenderAsMap())
		enc_dec = append(enc_dec, fields.GenEncDecFromEnum(enum, enum_gleam_type).RenderAsMap())
	}

	for _, msg := range all_messages {
		for _, oo := range msg.OneOfs() {
			// don't need oneof generators
			oo_gleam_type := fields.GleamTypeFromOnoeOf(msg, oo)
			oo_gleam_type_map := oo_gleam_type.RenderAsMap()
			gleam_types_map = append(gleam_types_map, oo_gleam_type_map)
			enc_dec = append(enc_dec, fields.GenEncDecFromOneOf(msg, oo, oo_gleam_type).RenderAsMap())
		}

		msg_gleam_type := fields.GleamTypeFromMessage(msg)

		if generator := fields.GeneratorFnFromGleamType(msg_gleam_type); generator != nil {
			generators = append(generators, generator.RenderAsMap())
		}

		msg_gleam_type_map := msg_gleam_type.RenderAsMap()
		gleam_types_map = append(gleam_types_map, msg_gleam_type_map)
		enc_dec = append(enc_dec, fields.GenEncDecFromMessage(msg, msg_gleam_type).RenderAsMap())
	}

	g.AddGeneratorTemplateFile(strings.Replace(pkg.ProtoName().String(), ".", "/", -1)+".gleam", g.tpl, map[string]interface{}{
		"imports":    imports,
		"package":    pkg.ProtoName().LowerSnakeCase().String(),
		"messages":   gleam_types_map,
		"generators": generators,
		"enc_dec":    enc_dec,
	})
}
