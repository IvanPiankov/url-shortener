package redirect

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage"
)

type UrlGetter interface {
	GetUrl(ctx context.Context, alias string) (string, error)
}

func New(log *slog.Logger, urlGetter UrlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("Alias is empty")

			render.JSON(w, r, resp.Error("Invalid Request"))

			return
		}
		resUrl, err := urlGetter.GetUrl(r.Context(), alias)

		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("Url with alias not found")

			render.JSON(w, r, resp.Error("Not Found"))
			return
		}

		if err != nil {
			log.Error(err.Error())
			log.Info("Failed to get url")

			render.JSON(w, r, resp.Error("Not Found"))
			return
		}

		http.Redirect(w, r, resUrl, http.StatusFound)
	}

}
