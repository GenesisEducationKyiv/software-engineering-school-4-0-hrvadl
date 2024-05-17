package rate

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func NewHandler(rg Getter, log *slog.Logger) *Handler {
	return &Handler{
		rg:  rg,
		log: log,
	}
}

type Getter interface {
	GetRate(ctx context.Context) (float32, error)
}

type Handler struct {
	log *slog.Logger
	rg  Getter
}

func (h *Handler) GetRate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	rat, err := h.rg.GetRate(ctx)
	if err != nil {
		h.log.Error("Failed to get rate", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid status value"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprint(rat)))
}