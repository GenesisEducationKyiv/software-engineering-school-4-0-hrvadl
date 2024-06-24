package app

import (
	"errors"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const successMessage = "OK"

func newMarshaller() *marshaller {
	return &marshaller{
		Marshaler: &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true,
				UseProtoNames:   true,
			},
		},
	}
}

type marshaller struct {
	runtime.Marshaler
}

func (cm *marshaller) Marshal(v any) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return nil, errors.New("failed to cast value to proto msg interface")
	}

	return cm.Marshaler.Marshal(newSuccessResponse(successMessage, msg))
}
