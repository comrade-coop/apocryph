package proto

import (
	"encoding/json"

	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bufbuild/protoyaml-go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"gopkg.in/yaml.v2"
)

type Format struct {
	unmarshal func(b []byte, m protoreflect.ProtoMessage) error
	marshal   func(m protoreflect.ProtoMessage) ([]byte, error)
}

func (f Format) Unmarshal(b []byte, m protoreflect.ProtoMessage) error {
	return f.unmarshal(b, m)
}

func (f Format) Marshal(m protoreflect.ProtoMessage) ([]byte, error) {
	return f.marshal(m)
}

func MarshalYaml(message proto.Message) ([]byte, error) { // Copied from protoyaml to use yaml.v2 (due to https://github.com/go-yaml/yaml/issues/720)
	data, err := protojson.Marshal(message)
	if err != nil {
		return nil, err
	}
	var jsonVal interface{}
	if err := json.Unmarshal(data, &jsonVal); err != nil {
		return nil, err
	}
	return yaml.Marshal(jsonVal)
}

var Formats = map[string]Format{
	"json":   {protojson.Unmarshal, protojson.Marshal},
	"yaml":   {protoyaml.Unmarshal, MarshalYaml},
	"yml":    {protoyaml.Unmarshal, MarshalYaml},
	"pb":     {proto.Unmarshal, proto.Marshal},
	"pbtext": {prototext.Unmarshal, prototext.Marshal},
}

var FormatNames []string

func init() {
	FormatNames := make([]string, 0, len(Formats))
	for name := range Formats {
		FormatNames = append(FormatNames, name)
	}
}

func Unmarshal(format string, bytes []byte, m protoreflect.ProtoMessage) error {
	f, ok := Formats[format]
	if !ok {
		return errors.New("Unknown format: " + format)
	}
	err := f.Unmarshal(bytes, m)
	if err != nil {
		return fmt.Errorf("Failed unmarshalling as %s: %w", format, err)
	}
	return nil
}

func Marshal(format string, m protoreflect.ProtoMessage) ([]byte, error) {
	f, ok := Formats[format]
	if !ok {
		return nil, errors.New("Unknown format: " + format)
	}
	bytes, err := f.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("Failed unmarshalling as %s: %w", format, err)
	}
	return bytes, nil
}

func DetectFormat(path string) string {
	for name := range Formats {
		if strings.HasSuffix(path, "."+name) {
			return name
		}
	}
	return path
}

func UnmarshalFile(path string, format string, m protoreflect.ProtoMessage) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if format == "" {
		format = DetectFormat(path)
	}

	return Unmarshal(format, bytes, m)
}

func MarshalFile(path string, format string, m protoreflect.ProtoMessage) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	if format == "" {
		format = DetectFormat(path)
	}

	bytes, err := Marshal(format, m)
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)
	return err
}
