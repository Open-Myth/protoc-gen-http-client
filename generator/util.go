package generator

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/modern-go/reflect2"
	options "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/descriptorpb"
)

// getOption retrieves an option from a protocol buffer message.
// It takes in a message of type protoreflect.Message and an extension of type *protoimpl.ExtensionInfo.
// It returns a proto.Message and a bool indicating whether the option was found or not.
func getOption(message protoreflect.Message, extension *protoimpl.ExtensionInfo) (proto.Message, bool) {
	ref := message.Get(extension.TypeDescriptor())
	if !ref.IsValid() {
		return nil, false
	}
	inter := ref.Message().Interface()
	if inter == nil {
		return nil, false
	}
	if reflect2.IsNil(inter) {
		return nil, false
	}

	// TODO: use reflection to get the option
	// ext := proto.GetExtension(meth.Options, options.E_Http)
	// opts, ok := ext.(*options.HttpRule)

	return inter, true
}

func getGoPackageName(file *descriptorpb.FileDescriptorProto) (impPath string, pkg string, ok bool) {
	opt := file.GetOptions().GetGoPackage()
	if opt == "" {
		impPath = file.GetPackage()
		pkg = impPath
		return impPath, pkg, false
	}
	// A semicolon-delimited suffix delimits the import path and package name.
	sc := strings.Index(opt, ";")
	if sc >= 0 {
		return opt[:sc], opt[sc+1:], true
	}
	// The presence of a slash implies there's an import path.
	slash := strings.LastIndex(opt, "/")
	if slash >= 0 {
		return opt, opt[slash+1:], true
	}
	return "", opt, true
}

type BuildInfoReaderFunc = func() (*debug.BuildInfo, bool)

// GetVersion description of the Go function.
func GetVersion(infoReader BuildInfoReaderFunc) string {
	info, ok := infoReader()
	if !ok || info.Main.Version == "" {
		return "(unknown)"
	}

	return info.Main.Version
}

func GetPathFromHTTPRule(opt *options.HttpRule) (string, string) {
	switch opt.Pattern.(type) {
	case *options.HttpRule_Get:
		rule := opt.Pattern.(*options.HttpRule_Get)
		return http.MethodGet, rule.Get
	case *options.HttpRule_Post:
		rule := opt.Pattern.(*options.HttpRule_Post)
		return http.MethodPost, rule.Post
	case *options.HttpRule_Put:
		rule := opt.Pattern.(*options.HttpRule_Put)
		return http.MethodPut, rule.Put
	case *options.HttpRule_Delete:
		rule := opt.Pattern.(*options.HttpRule_Delete)
		return http.MethodDelete, rule.Delete
	case *options.HttpRule_Patch:
		rule := opt.Pattern.(*options.HttpRule_Patch)
		return http.MethodDelete, rule.Patch
	}

	return "", ""
}

// ExtractAPIOptions retrieves the HTTP options for a given method descriptor.
func ExtractAPIOptions(meth *descriptorpb.MethodDescriptorProto) (*options.HttpRule, error) {
	if meth.Options == nil {
		return nil, nil
	}
	if !proto.HasExtension(meth.Options, options.E_Http) {
		return nil, nil
	}
	ext := proto.GetExtension(meth.Options, options.E_Http)
	opts, ok := ext.(*options.HttpRule)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want an HttpRule", ext)
	}

	return opts, nil
}
