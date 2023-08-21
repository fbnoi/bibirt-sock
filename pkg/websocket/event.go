package websocket

import (
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

var factories = map[string]*sync.Pool{}

func RegisterMessage(m proto.Message) {
	typeUrl := string(m.ProtoReflect().Descriptor().FullName())
	if _, ok := factories[typeUrl]; !ok {
		factories[typeUrl] = &sync.Pool{
			New: func() any {
				nm := proto.Clone(m)
				proto.Reset(nm)
				return nm
			},
		}
	}
}

func GetMessage(typeUrl string) (proto.Message, error) {
	if factory, ok := factories[typeUrl]; ok {
		m := factory.Get()
		return m.(proto.Message), nil
	}

	return nil, errors.Errorf("internal.GetMessage: type %s is not registered", typeUrl)
}

func PutMessage(m proto.Message) error {
	typeUrl := string(m.ProtoReflect().Descriptor().FullName())
	if factory, ok := factories[typeUrl]; ok {
		factory.Put(m)
		return nil
	}

	return errors.Errorf("internal.GetMessage: type %s is not registered", typeUrl)
}
