package save

import (
	"errors"
	"log/slog"
	"net/http"
	resp "project/internal/lib/api/response"
	"project/internal/lib/logger/sl"
	"project/internal/lib/random"
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

// TODO: moveto config (ot to db)
const aliasLength = 6

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.url.save.New"

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

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			c.JSON(http.StatusConflict, resp.Error("url already exists"))
			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))
			c.JSON(http.StatusInternalServerError, resp.Error("failed to add url"))
			return
		}

		log.Info("url added", slog.Int64("id", id))
		responseOK(c, alias)
	}
}

func responseOK(c *gin.Context, alias string) {
	c.JSON(http.StatusOK, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
