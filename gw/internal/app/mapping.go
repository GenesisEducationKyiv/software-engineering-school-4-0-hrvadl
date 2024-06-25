package app

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

func newResponseMapper(log *slog.Logger) *responseMapper {
	return &responseMapper{log}
}

type responseMapper struct {
	log *slog.Logger
}

func (rm *responseMapper) mapGRPCErr(
	_ context.Context,
	_ *runtime.ServeMux,
	_ runtime.Marshaler,
	w http.ResponseWriter,
	_ *http.Request,
	err error,
) {
	s := status.Convert(err)
	res, err := json.Marshal(newErrResponse(s.Message()))
	if err != nil {
		rm.log.Error("Failed to convert marshall err response", slog.Any("err", err))
	}

	w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
	if _, err := w.Write(res); err != nil {
		rm.log.Error("Failed to write error", slog.Any("err", err))
	}
}
