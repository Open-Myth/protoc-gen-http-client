package main

import (
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/Open-Myth/protoc-gen-http-client/generator"
)

const (
	extensionName = "http-client"
)

func main() {

	opt := protogen.Options{}
	opt.Run(func(p *protogen.Plugin) error {
		p.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		for _, f := range p.Files {
			if !f.Generate {
				continue
			}
			isGen := false
			gen := generator.NewGenerator(p, f, extensionName)

			svcInfos := gen.GetServiceInfos()
			svcDescInfos := gen.GetServiceDescInfos()
			for i, svc := range svcInfos {
				stmt := generateHTTPClientByService(svc, svcDescInfos[i])
				if stmt != nil {
					gen.GetJen().Add(stmt)
					isGen = true
				}
			}
			if isGen {
				if err := gen.Render(); err != nil {
					grpclog.Errorf("unable to generate %s: %v", extensionName, err)
				}
			}

		}

		return nil
	})

}
