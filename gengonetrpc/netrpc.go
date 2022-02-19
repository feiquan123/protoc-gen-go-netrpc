package gengonetrpc

import (
	"bytes"
	"log"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
)

const (
	fmtPackage     = protogen.GoImportPath("fmt")
	netPackage     = protogen.GoImportPath("net")
	rpcPackage     = protogen.GoImportPath("net/rpc")
	jsonrpcPackage = protogen.GoImportPath("net/rpc/jsonrpc")
	syncPackage    = protogen.GoImportPath("sync")
)

// GenerateFileContent generates the gRPC service definitions, excluding the package statement.
func GenerateFileContent(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile) {
	if len(file.Services) == 0 {
		return
	}

	// import
	g.P("var _ = ", fmtPackage.Ident("Errorf"))
	g.P("var _ = ", netPackage.Ident("Listen"))
	g.P("var _ ", rpcPackage.Ident("Client"))
	g.P("var _ = ", jsonrpcPackage.Ident("NewClientCodec"))
	g.P("var _ ", syncPackage.Ident("Once"))
	g.P()

	// interface、service、client struct
	for _, service := range file.Services {
		genService(gen, file, g, service)
	}
}

func genService(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service) {
	spec := buildServiceSpec(g, service)

	var buf bytes.Buffer
	t := template.Must(template.New("").Parse(tmpService))
	err := t.Execute(&buf, spec)
	if err != nil {
		log.Fatal(err)
	}
	g.P(buf.String())
}

func buildServiceSpec(g *protogen.GeneratedFile, service *protogen.Service) *ServiceSpec {
	spec := &ServiceSpec{
		ServiceName: service.GoName,
	}

	for _, method := range service.Methods {
		spec.MethodList = append(spec.MethodList, ServiceMethodSpec{
			MethodName:     method.GoName,
			InputTypeName:  g.QualifiedGoIdent(method.Input.GoIdent),
			OutputTypeName: g.QualifiedGoIdent(method.Output.GoIdent),
		})
	}
	return spec
}
