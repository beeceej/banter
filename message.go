package banter

import (
	"github.com/beeceej/banter/pb"
)

func MsgPing(origin *pb.Peer) *pb.Message {
	return &pb.Message{Msg: pb.Msg_PING, Origin: origin}
}

func MsgPong(origin *pb.Peer) *pb.Message {
	return &pb.Message{Msg: pb.Msg_PONG, Origin: origin}
}
