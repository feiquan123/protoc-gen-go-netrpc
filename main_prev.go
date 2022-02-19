// +build prev

package main

import (
	"io/ioutil"
	"os"

	_ "github.com/feiquan123/protoc-gen-go-netrpc/gengonetrpcprev"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"google.golang.org/protobuf/proto"
)

//go:generate go install -tags=prev .
func main() {
	g := generator.New()

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		g.Error(err, "reading input")
	}

	if err := proto.Unmarshal(data, g.Request); err != nil {
		g.Error(err, "parsing input proto")
	}

	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}

	g.CommandLineParameters(g.Request.GetParameter())

	// Create a wrapped version of the Descriptors and EnumDescriptors that
	// point to the file that defines them.
	g.WrapTypes()

	g.SetPackageNames()
	g.BuildTypeNameMap()
	g.GenerateAllFiles()

	// Send back the results
	data, err = proto.Marshal(g.Response)
	if err != nil {
		g.Error(err, "failed to marshal output proto")
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		g.Error(err, "failed to write output proto")
	}
}
