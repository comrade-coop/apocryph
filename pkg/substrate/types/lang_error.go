package types

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
)

type LangError struct {
	IsCouldNotReadInput bool
}

func (e LangError) Encode(encoder scale.Encoder) error {
	switch {
	case e.IsCouldNotReadInput:
		return encoder.PushByte(0)
	}
	return nil
}

func (e *LangError) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch b {
	case 0:
		e.IsCouldNotReadInput = true
	default:
		return fmt.Errorf("Unknown discriminator for encoded LangError: %d", b)
	}
	return nil
}
