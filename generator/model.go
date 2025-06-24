package generator

import (
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/descriptorpb"
)

type MessageProto[E protoimpl.ExtensionInfo] struct {
	Proto  *descriptorpb.DescriptorProto
	ExInfo *E
}
