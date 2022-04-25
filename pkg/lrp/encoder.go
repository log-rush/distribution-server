package lrp

type LRPEncoder interface {
	Encode(msg LRPMessage) []byte
}

type LRPEncoderImpl struct{}

func NewEncoder() LRPEncoder {
	return &LRPEncoderImpl{}
}

func (e *LRPEncoderImpl) Encode(msg LRPMessage) []byte {
	return append([]byte{byte(msg.OPCode)}, msg.Payload...)
}
