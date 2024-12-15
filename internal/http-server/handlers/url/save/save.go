package save

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const (
	aliasLength = 4
)

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlServer URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode request body")

			render.JSON(w, r, resp.Error("Failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("Invalid request body")

			render.JSON(w, r, resp.Error("Invalid request"))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlServer.SaveURL(req.URL, alias)

		if errors.Is(err, storage.ErrURLAlreadyExists) {
			log.Error("Url already exists")

			render.JSON(w, r, resp.Error("Url already exists"))

			return
		}
		if err != nil {
			log.Error("Unknown error")
			log.Error(err.Error())

			render.JSON(w, r, resp.Error("Unknown error"))

			return
		}

		log.Info("Url added", slog.Int64("id", id))

		// TODO: Transfer to function
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
