package iface

import (
	"github.com/golang/protobuf/proto"
)

type IProtoMessage interface {
	proto.Message
	Marshal() ([]byte, error)
	MarshalTo([]byte) (int, error)
	Unmarshal([]byte) error
	Size() int
}
