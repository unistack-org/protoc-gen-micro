package main

import (
	"strings"

	openapi_options "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

type Error struct {
	packageName string
	types       []string
}

func (g *Generator) httpGenerate(component string, plugin *protogen.Plugin) error {
	errors := make(map[string]struct{})

	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}
		if len(file.Services) == 0 {
			continue
		}

		gname := file.GeneratedFilenamePrefix + "_micro_" + component + ".pb.go"
		gfile := plugin.NewGeneratedFile(gname, ".")

		gfile.P("// Code generated by protoc-gen-micro")
		gfile.P("// source: ", *file.Proto.Name)
		gfile.P("package ", file.GoPackageName)

		gfile.P()
		gfile.P("import (")
		gfile.P(`"context"`)
		gfile.P()
		gfile.P(`micro_api "github.com/unistack-org/micro/v3/api"`)
		gfile.P(`micro_client_http "github.com/unistack-org/micro-client-http/v3"`)
		gfile.P(`micro_client "github.com/unistack-org/micro/v3/client"`)
		gfile.P(`micro_server "github.com/unistack-org/micro/v3/server"`)
		gfile.P(")")
		gfile.P()

		gfile.P("// Reference imports to suppress errors if they are not otherwise used.")
		gfile.P("var (")
		gfile.P("_ ", "micro_api.Endpoint")
		gfile.P("_ ", "context.Context")
		gfile.P(" _ ", "micro_client.Option")
		gfile.P(" _ ", "micro_server.Option")
		gfile.P(")")
		gfile.P()

		for _, service := range file.Services {
			generateServiceClient(gfile, service)
			generateServiceClientMethods(gfile, service, true)
			generateServiceServer(gfile, service)
			generateServiceServerMethods(gfile, service)
			generateServiceRegister(gfile, service)
			if component == "http" {
				for k, v := range getErrors(service) {
					errors[k] = v
				}
			}
		}
	}

	files := make(map[string]*Error)
	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}

		err, ok := files[file.GeneratedFilenamePrefix]
		if !ok {
			err = &Error{packageName: string(file.GoPackageName)}
		}
		fok := false
		for _, message := range file.Messages {
			if _, ok := errors["."+string(message.Desc.FullName())]; ok {
				fok = true
				err.types = append(err.types, string(message.Desc.FullName()))
			}
		}
		if fok {
			files[file.GeneratedFilenamePrefix] = err
		}
	}

	for file, err := range files {
		gfile := plugin.NewGeneratedFile(file+"_micro_errors.pb.go", ".")
		generateServiceErrors(gfile, err)
	}

	return nil
}

func getErrors(service *protogen.Service) map[string]struct{} {
	errors := make(map[string]struct{})

	for _, method := range service.Methods {
		if method.Desc.Options() == nil {
			continue
		}
		if !proto.HasExtension(method.Desc.Options(), openapi_options.E_Openapiv2Operation) {
			continue
		}

		opts := proto.GetExtension(method.Desc.Options(), openapi_options.E_Openapiv2Operation)
		if opts == nil {
			continue
		}

		r := opts.(*openapi_options.Operation)
		for _, response := range r.Responses {
			if response.Schema == nil || response.Schema.JsonSchema == nil {
				continue
			}
			errors[response.Schema.JsonSchema.Ref] = struct{}{}
		}
	}

	return errors
}

func generateServiceErrors(gfile *protogen.GeneratedFile, err *Error) {
	gfile.P("package ", err.packageName)
	gfile.P("import (")
	gfile.P(`"fmt"`)
	gfile.P(")")
	for _, typ := range err.types {
		gfile.P("func (err *", typ[strings.LastIndex(typ, ".")+1:], ") Error() string {")
		gfile.P(`return fmt.Sprintf("%#v", err)`)
		gfile.P("}")
		gfile.P()
	}
}
