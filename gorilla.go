package main

import (
	"google.golang.org/protobuf/compiler/protogen"
)

var (
	gorillaPackageFiles map[protogen.GoPackageName]struct{}
)

func (g *Generator) gorillaGenerate(component string, plugin *protogen.Plugin) error {
	gorillaPackageFiles = make(map[protogen.GoPackageName]struct{})
	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}
		if len(file.Services) == 0 {
			continue
		}
		if _, ok := gorillaPackageFiles[file.GoPackageName]; ok {
			continue
		}

		gorillaPackageFiles[file.GoPackageName] = struct{}{}
		gname := "micro" + "_" + component + ".pb.go"
		gfile := plugin.NewGeneratedFile(gname, ".")

		gfile.P("// Code generated by protoc-gen-micro")
		gfile.P("package ", file.GoPackageName)
		gfile.P()

		gfile.P("import (")
		gfile.P(`"fmt"`)
		gfile.P(`"net/http"`)
		gfile.P(`"reflect"`)
		gfile.P(`"strings"`)
		gfile.P(`mux "github.com/gorilla/mux"`)
		gfile.P(`micro_api "github.com/unistack-org/micro/v3/api"`)
		gfile.P(")")
		gfile.P()

		gfile.P("func RegisterHandlers(r *mux.Router, h interface{}, eps []*micro_api.Endpoint) error {")
		gfile.P("v := reflect.ValueOf(h)")
		gfile.P("if v.NumMethod() < 1 {")
		gfile.P(`return fmt.Errorf("handler has no methods: %T", h)`)
		gfile.P("}")
		gfile.P("for _, ep := range eps {")
		gfile.P(`idx := strings.Index(ep.Name, ".")`)
		gfile.P("if idx < 1 || len(ep.Name) <= idx {")
		gfile.P(`return fmt.Errorf("invalid api.Endpoint name: %s", ep.Name)`)
		gfile.P("}")
		gfile.P("name := ep.Name[idx+1:]")
		gfile.P("m := v.MethodByName(name)")
		gfile.P("if !m.IsValid() || m.IsZero() {")
		gfile.P(`return fmt.Errorf("invalid handler, method %s not found", name)`)
		gfile.P("}")
		gfile.P("rh, ok := m.Interface().(func(http.ResponseWriter, *http.Request))")
		gfile.P("if !ok {")
		gfile.P(`return fmt.Errorf("invalid handler: %#+v", m.Interface())`)
		gfile.P("}")
		gfile.P("r.HandleFunc(ep.Path[0], rh).Methods(ep.Method...).Name(ep.Name)")
		gfile.P("}")
		gfile.P("return nil")
		gfile.P("}")
	}

	return nil
}
