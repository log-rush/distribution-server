package lrp

import "errors"

type LRPMessage struct {
	OPCode  LRPOpcode
	Payload []byte
}

type LRPOpcode = int8

const (
	OprSubscribe   LRPOpcode = 65 // 0b010010 // 18
	OprUnsubscribe LRPOpcode = 66 // 0b010100 // 20
	OprAlive       LRPOpcode = 67 // 0b001010 // 10
	OprStillAlive  LRPOpcode = 68 // 0b001100 // 12
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
