package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/storage/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"	//логгер для тестов
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string	// название теста
		alias     string
		url       string
		respError string	// ожидаемая ошибка
		mockError error
	}{
		// успешный
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
		// успешный без алиаса
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://google.com",
		},
		// неуспешный с пустым url
		{
			name:      "Empty URL",
			url:       "",
			alias:     "some_alias",
			respError: "field URL is a required field",
		},
		// неуспешный с невалидным url
		{
			name:      "Invalid URL",
			url:       "some invalid URL",
			alias:     "some_alias",
			respError: "field URL is not a valid URL",
		},
		// неуспешный с ошибкой из стораджа
		{
			name:      "SaveURL Error",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "fail add url",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			// параллельный запуск тестов
			t.Parallel()
			// мок стораджа для тестов
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			StorageMock := mocks.NewMockStorage(ctrl)

			// настройка мока
			if tc.respError == "" || tc.mockError != nil {
				// ожидаем, когда будет вызов SaveURL с url и любым алиасом, то вернуть ошибку
				StorageMock.EXPECT().
					SaveURL(tc.url, gomock.Any()).
					Return(tc.mockError).
					Times(1) // ожидаем вызов SaveURL 1 раз
			}
			// создаем хендлер
			handler := save.New(slogdiscard.NewDiscardLogger(), StorageMock)
			// тело запроса
			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)
			// создаем запрос
			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err) // выходим из теста, если запрос не создался

			// записываем ответ
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			
			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Response
			// проверяем результат
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}