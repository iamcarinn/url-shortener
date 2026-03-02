package save

import (
	"errors"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

type Request struct {
	URL string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 10

func New(log *slog.Logger, storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request)  {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),											// доп пар-р в логгер
			slog.String("request_id", middleware.GetReqID(r.Context())),	// трейсинг запросов
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("fail decode recuest body", sl.Err(err))

			render.JSON(w, r, response.Error("fail decode request"))	// возв. json с ответои клиенту
			
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// Валидируем запрос, кот. пришел
		if err := validator.New().Struct(req); err != nil {	// валидируем структуру
			validateErr := err.(validator.ValidationErrors)	// приводим ошибку к нужному типу
			
			log.Error("invalid request", sl.Err(err)) // логгируем ошибку

			render.JSON(w, r, response.ValidationError(validateErr))	// возв. json с ответои с ошибкой клиенту
			return

		}

		// Указан ли алиас в запросе
		alias := req.Alias
		if alias == "" {
			// генерируем если нет
			alias = random.NewRandomString(aliasLength)
		}
		// TODO: доделать
		// TODO: сделать проверку, если уже есть такой алиас
		id, err := urlSaver.SaveURL(req.URL, alias)

		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, response.Error("url already exists"))

			return 
		}



	}

}