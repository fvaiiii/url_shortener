package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	resp "project/internal/lib/api/response"
	"project/internal/lib/logger/sl"
	"project/internal/storage"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.redirect.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", c.GetString("request_id")),
		)

		alias := c.Param("alias")
		if alias == "" {
			log.Info("alias is empty")
			c.JSON(http.StatusBadRequest, resp.Error("invalid request"))
			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			c.JSON(http.StatusBadRequest, resp.Error("invalid request"))
			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))
			c.JSON(http.StatusInternalServerError, resp.Error("internal error"))
			return
		}

		log.Info("got url", slog.String("url", resURL))

		c.Redirect(http.StatusFound, resURL)
	}
}
