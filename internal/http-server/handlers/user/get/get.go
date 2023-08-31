package get

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Dmitriy770/user-segmentation-service/internal/lib/logger/api/response"
	"github.com/Dmitriy770/user-segmentation-service/internal/lib/logger/sl"
	"github.com/Dmitriy770/user-segmentation-service/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	UserId   int      `json:"user_id"`
	Segments []string `json:"segments"`
	response.Response
}

type UserGeter interface {
	GetUser(id int) (*models.User, error)
}

func New(log *slog.Logger, userGeter UserGeter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.update"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userIdStr := chi.URLParam(r, "userId")
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			log.Info("failed to get user", sl.Err(err))
			render.JSON(w, r, response.Error("not valid user id"))
			return
		}

		user, err := userGeter.GetUser(int(userId))
		if err != nil {
			log.Info("failed to get user", sl.Err(err))
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		log.Info("get user", slog.Int("user id", userId))

		render.JSON(w, r, Response{
			UserId:   userId,
			Segments: user.Segments,
			Response: response.OK(),
		})
	}
}
