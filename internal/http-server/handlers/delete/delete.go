package delete

import (
	"errors"
	"log/slog"
	"net/http"
	resp "project/internal/lib/api/response"
	"project/internal/lib/logger/sl"
	"project/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}
type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=URLDelete
type URLDelete interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDelete URLDelete) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.delete.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", c.GetString("request_id")),
		)

		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			if errors.Is(err, errors.New("EOF")) {

				log.Error("request body is empty")
				c.JSON(http.StatusBadRequest, resp.Error("empty request"))
				return
			}
			log.Error("failed to decode request body", sl.Err(err))
			c.JSON(http.StatusBadRequest, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))
			c.JSON(http.StatusBadRequest, resp.ValidationError(validateErr))
			return
		}

		deleteAlias := req.Alias
		if deleteAlias == "" {
			log.Error("alias is required for deletion")
			c.JSON(http.StatusBadRequest, resp.Error("alias is required for deletion"))
			return
		}

		err := urlDelete.DeleteURL(deleteAlias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", deleteAlias))
			c.JSON(http.StatusConflict, resp.Error("url not found"))
			return
		}

		if err != nil {
			log.Error("failed to delete url", sl.Err(err))
			c.JSON(http.StatusInternalServerError, resp.Error("failed to delete url"))
			return
		}

		log.Info("url deleted successfully", slog.String("alias", deleteAlias))

		responseOK(c, deleteAlias)
	}
}

func responseOK(c *gin.Context, deleteAlias string) {
	c.JSON(http.StatusOK, Response{
		Response: resp.OK(),
		Alias:    deleteAlias,
	})
}
