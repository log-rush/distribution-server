package lrp

type LRPDecoder interface {
	Decode(msg []byte) (LRPMessage, error)
}

type LRPDecoderImpl struct{}

func NewDecoder() LRPDecoder {
	return &LRPDecoderImpl{}
}

func (e *LRPDecoderImpl) Decode(msg []byte) (LRPMessage, error) {
	if len(msg) == 0 {
		return LRPMessage{}, nil
	}
	return NewMesssage(int8(msg[0]), msg[1:]), nil
}
