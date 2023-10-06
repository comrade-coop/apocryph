package types

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
)

type Option[Value interface{}] struct {
	HasValue bool
	Value    Value
}

func (o Option[Value]) Encode(encoder scale.Encoder) error {
	return encoder.EncodeOption(o.HasValue, o.Value)
}

func (o *Option[Value]) Decode(decoder scale.Decoder) error {
	return decoder.DecodeOption(&o.HasValue, &o.Value)
}

type Result[Value interface{}, Error interface{}] struct {
	IsError bool
	Value   Value
	Error   Error
}

func (r Result[Value, Error]) Encode(encoder scale.Encoder) error {
	if r.IsError {
		err := encoder.PushByte(1)
		if err != nil {
			return err
		}
		return encoder.Encode(r.Error)
	} else {
		err := encoder.PushByte(0)
		if err != nil {
			return err
		}
		return encoder.Encode(r.Value)
	}
}

func (r *Result[Value, Error]) Decode(decoder scale.Decoder) error {
	b, _ := decoder.ReadOneByte()
	switch b {
	case 0:
		r.IsError = false
		return decoder.Decode(&r.Value)
	case 1:
		r.IsError = true
		return decoder.Decode(&r.Error)
	default:
		return fmt.Errorf("Unknown byte prefix for encoded Result<,>: %d", b)
	}
}
