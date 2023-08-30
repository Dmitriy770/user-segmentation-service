package add

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Dmitriy770/user-segmentation-service/internal/lib/logger/api/response"
	"github.com/Dmitriy770/user-segmentation-service/internal/lib/logger/sl"
	"github.com/Dmitriy770/user-segmentation-service/internal/serevices/segments"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Slug string `json:"slug" validate:"required"`
}

type Response struct {
	response.Response
}

type SegmentCreater interface {
	AddSegment(slug string) error
}

func New(log *slog.Logger, segmentCreater SegmentCreater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segment.add"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode")
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validatorErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.Error("invalid request"))
			render.JSON(w, r, response.ValidatiobError(validatorErr))

			return
		}

		err = segmentCreater.AddSegment(req.Slug)
		if errors.Is(err, segments.ErrSlugBusy) {
			log.Info("segment already exists", slog.String("slug", req.Slug))
			render.JSON(w, r, response.Error("segment already exists"))
			return
		}
		if err != nil {
			log.Error("failed to add segment", sl.Err(err))
			render.JSON(w, r, response.Error("failed to add segment"))
			return
		}

		log.Info("segment added", slog.String("slug", req.Slug))

		render.JSON(w, r, response.OK())
	}
}
