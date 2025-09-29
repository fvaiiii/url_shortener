package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"project/internal/http-server/handlers/url/save"
	"project/internal/http-server/handlers/url/save/mocks"
	"project/internal/lib/logger/handlers/slogdiscard"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cases := []struct {
		name         string
		alias        string
		url          string
		respError    string
		mockError    error
		expectedCode int
	}{
		{
			name:         "Success",
			alias:        "test_alias",
			url:          "https://google.com",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Empty alias",
			alias:        "",
			url:          "https://google.com",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Empty URL",
			url:          "",
			alias:        "some_alias",
			respError:    "field URL is a required field",
			expectedCode: http.StatusBadRequest, // 400
		},
		{
			name:         "Invalid URL",
			url:          "some invalid URL",
			alias:        "some_alias",
			respError:    "field URL is not a valid URL",
			expectedCode: http.StatusBadRequest, // 400
		},
		{
			name:         "SaveURL Error",
			alias:        "test_alias",
			url:          "https://google.com",
			respError:    "failed to add url",
			mockError:    errors.New("unexpected error"),
			expectedCode: http.StatusInternalServerError, // 500
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if c.expectedCode == http.StatusOK || c.mockError != nil {
				urlSaverMock.On("SaveURL", c.url, mock.AnythingOfType("string")).
					Return(int64(1), c.mockError).
					Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			router := gin.New()
			router.POST("/save", handler)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, c.url, c.alias)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, c.expectedCode, rr.Code,
				"Expected status %d, got %d. Response body: %s",
				c.expectedCode, rr.Code, rr.Body.String())

			body := rr.Body.String()
			var resp save.Response

			if rr.Code == http.StatusOK || rr.Code == http.StatusBadRequest {
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				require.Equal(t, c.respError, resp.Error)
			}
		})
	}
}
