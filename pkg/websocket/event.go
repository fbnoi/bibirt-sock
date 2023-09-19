package websocket

import (
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var messages = map[string]proto.Message{}

func RegisterMessage(m proto.Message) {
	a, _ := anypb.New(m)
	typeUrl := string(a.GetTypeUrl())
	if _, ok := messages[typeUrl]; !ok {
		messages[typeUrl] = m
	}
}

func GetMessage(typeUrl string) (proto.Message, error) {
	if m, ok := messages[typeUrl]; ok {
		nm := proto.Clone(m)
		proto.Reset(nm)
		return nm, nil
	}

	return nil, errors.Errorf("internal.GetMessage: type %s is not registered", typeUrl)
}
