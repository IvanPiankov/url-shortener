package delete

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

type UrlDeleter interface {
	DeleteUrl(ctx context.Context, alias string) error
}

func New(log *slog.Logger, UrlDeleter UrlDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("Alias is empty")
			render.JSON(w, r, resp.Error("Alias is empty"))
			return
		}

		err := UrlDeleter.DeleteUrl(r.Context(), alias)

		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("Url by alias not found")
			render.JSON(w, r, resp.Error("Not found url"))
			return
		}

		if err != nil {
			log.Info("Unknown error")
			render.JSON(w, r, resp.Error("Unknown error"))
			return
		}

		render.JSON(w, r, resp.OK())
	}

}
