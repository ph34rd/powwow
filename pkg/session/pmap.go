package session

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	"github.com/ph34rd/powwow/pb"
)

const (
	_ = iota
	pbServerHandshake
	pbClientHandshake
	pbWoWRequest
	pbWoWResponse
)

var typeToIndex = map[string]uint16{
	proto.MessageName(&pb.ServerHandshake{}): pbServerHandshake,
	proto.MessageName(&pb.ClientHandshake{}): pbClientHandshake,
	proto.MessageName(&pb.WoWRequest{}):      pbWoWRequest,
	proto.MessageName(&pb.WoWResponse{}):     pbWoWResponse,
}

func pbMessageType(m proto.Message) (uint16, error) {
	name := proto.MessageName(m)
	v, ok := typeToIndex[name]
	if !ok {
		return 0, fmt.Errorf("proto message not found: %s", name)
	}
	return v, nil
}
