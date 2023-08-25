package message

import "google.golang.org/protobuf/proto"

func TypeUrl(m proto.Message) string {
	return string(m.ProtoReflect().Descriptor().FullName())
}
