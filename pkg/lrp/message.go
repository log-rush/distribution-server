package lrp

import "errors"

type LRPMessage struct {
	OPCode  LRPOpcode
	Payload []byte
}

type LRPOpcode = int8

const (
	OprSubscribe   LRPOpcode = 0b010010
	OprUnsubscribe LRPOpcode = 0b010100
	OprAlive       LRPOpcode = 0b001010
	OprStillAlive  LRPOpcode = 0b001100
)

var (
	ErrMessageEmpty = errors.New("LRP: Cannot decode empty message")
)

func NewMesssage(operation LRPOpcode, payload []byte) LRPMessage {
	return LRPMessage{
		OPCode:  operation,
		Payload: payload,
	}
}
