package app

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

func newNilMarshaller() *nilMarshaler {
	return &nilMarshaler{
		Marshaler: &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true,
				UseProtoNames:   true,
			},
		},
	}
}

type nilMarshaler struct {
	runtime.Marshaler
}

func (cm *nilMarshaler) Marshal(any) ([]byte, error) {
	return nil, nil
}

func (cm *nilMarshaler) Unmarshal([]byte, any) error {
	return nil
}
