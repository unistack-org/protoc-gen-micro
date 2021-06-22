package main

import (
	"google.golang.org/protobuf/compiler/protogen"
)

func (g *Generator) rpcGenerate(component string, plugin *protogen.Plugin, genClient bool, genServer bool) error {
	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}
		if len(file.Services) == 0 {
			continue
		}

		gname := file.GeneratedFilenamePrefix + "_micro_" + component + ".pb.go"
		path := file.GoImportPath
		if g.standalone {
			path = "."
		}
		gfile := plugin.NewGeneratedFile(gname, path)

		gfile.P("// Code generated by protoc-gen-go-micro. DO NOT EDIT.")
		gfile.P("// protoc-gen-go-micro version: " + versionComment)
		gfile.P("// source: ", file.Proto.GetName())
		gfile.P()
		gfile.P("package ", file.GoPackageName)
		gfile.P()

		gfile.Import(contextPackage)
		gfile.Import(microApiPackage)
		if genClient {
			gfile.Import(microClientPackage)
		}
		if genServer {
			gfile.Import(microServerPackage)
		}
		for _, service := range file.Services {
			if genClient {
				generateServiceClient(gfile, service)
				generateServiceClientMethods(gfile, service, false)
			}
			if genServer {
				generateServiceServer(gfile, service)
				generateServiceServerMethods(gfile, service)
				generateServiceRegister(gfile, service)
			}
		}
	}

	return nil
}
