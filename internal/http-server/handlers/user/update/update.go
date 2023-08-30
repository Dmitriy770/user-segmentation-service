package update

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Dmitriy770/user-segmentation-service/internal/lib/logger/api/response"
	"github.com/Dmitriy770/user-segmentation-service/internal/lib/logger/sl"
	"github.com/Dmitriy770/user-segmentation-service/internal/serevices/users"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	UserId         int      `json:"user_id" validate:"required"`
	AddSegments    []string `json:"add_segments"`
	DeleteSegments []string `json:"delete_segments"`
}

type Response struct {
	response.Response
}

type SegmentUpdater interface {
	UpdateUser(userId int, slugsForAdd []string, slugsForDelete []string) error
}

func New(log *slog.Logger, segmentUpdater SegmentUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segment.update"

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

		err = segmentUpdater.UpdateUser(req.UserId, req.AddSegments, req.DeleteSegments)
		if errors.Is(err, users.ErrUserHaveSegment) {
			log.Info("failed to update user`s segments", sl.Err(err))
			render.JSON(w, r, response.Error("user have one of this segment"))
			return
		}
		if errors.Is(err, users.ErrUserDoesntHaveSegment) {
			log.Info("failed to update user`s segments", sl.Err(err))
			render.JSON(w, r, response.Error("user doesn`t have one of this segment"))
			return
		}
		if errors.Is(err, users.ErrSegmentDoesNotExist) {
			log.Info("failed to update user`s segments", sl.Err(err))
			render.JSON(w, r, response.Error("user segment doesn`t exist"))
			return
		}
		if err != nil {
			log.Info("failed to update user`s segments", sl.Err(err))
			render.JSON(w, r, response.Error("some error"))
			return
		}

		// err = segmentCreater.DeleteSegment(req.Slug)
		// if errors.Is(err, segments.ErrSlugNotFound) {
		// 	log.Info("segment already exists", slog.String("slug", req.Slug))
		// 	render.JSON(w, r, response.Error("segment alreadt exists"))
		// 	return
		// }
		// if err != nil {
		// 	log.Error("failed to delete segment", sl.Err(err))
		// 	render.JSON(w, r, response.Error("failed to delete segment"))
		// 	return
		// }

		// log.Info("segment deleted", slog.String("slug", req.Slug))

		render.JSON(w, r, response.OK())
	}
}
