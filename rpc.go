package main

import (
	"google.golang.org/protobuf/compiler/protogen"
)

func (g *Generator) rpcGenerate(component string, plugin *protogen.Plugin) error {
	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}
		if len(file.Services) == 0 {
			continue
		}

		gname := file.GeneratedFilenamePrefix + "_micro_" + component + ".pb.go"
		gfile := plugin.NewGeneratedFile(gname, file.GoImportPath)

		gfile.P("// Code generated by protoc-gen-micro")
		gfile.P("// source: ", *file.Proto.Name)
		gfile.P("package ", file.GoPackageName)

		gfile.P()
		gfile.P("import (")
		gfile.P(`"context"`)
		gfile.P()
		gfile.P(`micro_api "github.com/unistack-org/micro/v3/api"`)
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
			generateServiceClientMethods(gfile, service, false)
			generateServiceServer(gfile, service)
			generateServiceServerMethods(gfile, service)
			generateServiceRegister(gfile, service)
		}
	}

	return nil
}
