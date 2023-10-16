package proto

import (
	"errors"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var UnmarshalFormats = map[string]func(b []byte, m protoreflect.ProtoMessage) error{
	"json": protojson.Unmarshal,
	"pb":   proto.Unmarshal,
	"text": prototext.Unmarshal,
}

var UnmarshalFormatNames []string

func init() {
	UnmarshalFormatNames := make([]string, 0, len(UnmarshalFormats))
	for name := range UnmarshalFormats {
		UnmarshalFormatNames = append(UnmarshalFormatNames, name)
	}
}

func Unmarshal(format string, bytes []byte, m protoreflect.ProtoMessage) error {
	Unmarshal := UnmarshalFormats[format]
	if Unmarshal == nil {
		return errors.New("Unknown format: " + format)
	}
	return Unmarshal(bytes, m)
}
