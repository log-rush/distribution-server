package lrp

import "errors"

type LRPMessage struct {
	OPCode  LRPOpcode
	Payload []byte
}

type LRPOpcode = int8

const (
	OprSubscribe   LRPOpcode = 0b010010 // 18
	OprUnsubscribe LRPOpcode = 0b010100 // 20
	OprAlive       LRPOpcode = 0b001010 // 10
	OprStillAlive  LRPOpcode = 0b100010 // 34
	OprLog         LRPOpcode = 0b100100 // 36
	OprErr         LRPOpcode = 0b100110 // 38
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
