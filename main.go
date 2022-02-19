// +build !prev

package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/feiquan123/protoc-gen-go-netrpc/gengonetrpc"
	gengo "google.golang.org/protobuf/cmd/protoc-gen-go/internal_gengo"
	"google.golang.org/protobuf/compiler/protogen"
)

//go:generate go install .
func main() {
	// 新版使用默认生成的 protogen.Plugin  对象进行操作
	var (
		flags        flag.FlagSet
		plugins      = flags.String("plugins", "", "list of plugins to enable (supported values: netrpc)")
		importPrefix = flags.String("import_prefix", "", "prefix to prepend to import paths")
	)
	importRewriteFunc := func(importPath protogen.GoImportPath) protogen.GoImportPath {
		switch importPath {
		case "context", "fmt", "math":
			return importPath
		}
		if *importPrefix != "" {
			return protogen.GoImportPath(*importPrefix) + importPath
		}
		return importPath
	}
	protogen.Options{
		ParamFunc:         flags.Set,
		ImportRewriteFunc: importRewriteFunc,
	}.Run(func(gen *protogen.Plugin) error {
		netrpc := false
		for _, plugin := range strings.Split(*plugins, ",") {
			switch plugin {
			case "netrpc":
				netrpc = true
			case "":
			default:
				return fmt.Errorf("protoc-gen-go-netrpc: unknown plugin %q", plugin)
			}
		}
		_ = netrpc
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			g := gengo.GenerateFile(gen, f)
			if netrpc {
				gengonetrpc.GenerateFileContent(gen, f, g)
			}
		}
		gen.SupportedFeatures = gengo.SupportedFeatures
		return nil
	})
}
