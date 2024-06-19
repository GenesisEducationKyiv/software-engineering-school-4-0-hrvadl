package rate

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/gw/internal/transport/http/handlers"
)

const getRateTimeout = 6 * time.Second

func NewHandler(rg Getter, log *slog.Logger) *Handler {
	return &Handler{
		rg:  rg,
		log: log,
	}
}

//go:generate mockgen -destination=./mocks/mock_getter.go -package=mocks . Getter
type Getter interface {
	GetRate(ctx context.Context) (float32, error)
}

type Handler struct {
	log *slog.Logger
	rg  Getter
}

// GetRate godoc
// @Summary      Get USD -> UAH exchange rate
// @Tags         Rate
// @Produce      json
// @Success      200  {object}  handlers.Response[float32]
// @Failure      400  {object}  handlers.ErrorResponse
// @Router       /api/rate [get]
func (h *Handler) GetRate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), getRateTimeout)
	defer cancel()

	exchangeRate, err := h.rg.GetRate(ctx)
	if err != nil {
		h.httpFail(w, http.StatusBadRequest, err)
		return
	}

	h.httpSuccess(w, http.StatusOK, exchangeRate, "successfully got rate")
}

func (h *Handler) httpSuccess(w http.ResponseWriter, code int, data any, msg string) {
	res, err := handlers.NewSuccessResponse(msg, data)
	if err != nil {
		h.log.Error("Failed to construct success response", slog.Any("err", err))
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(code)
	if _, err := w.Write(res); err != nil {
		h.log.Error("Failed to write response", slog.Any("err", err))
	}
}

func (h *Handler) httpFail(w http.ResponseWriter, code int, err error) {
	res, err := handlers.NewErrResponse(err)
	if err != nil {
		h.log.Error("Failed to construct error response", slog.Any("err", err))
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(code)
	if _, err := w.Write(res); err != nil {
		h.log.Error("Failed to write response", slog.Any("err", err))
	}
}
