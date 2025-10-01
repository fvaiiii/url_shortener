package redirect_test

import (
	"net/http/httptest"
	"project/internal/http-server/handlers/get/mocks"
	"project/internal/http-server/handlers/redirect"
	"project/internal/lib/api"
	"project/internal/lib/logger/handlers/slogdiscard"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://www.google.com/",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			urlGetterMock := mocks.NewURLGetter(t)

			if c.respError == "" || c.mockError != nil {

				urlGetterMock.On("GetURL", c.alias).Return(c.url, c.mockError).Once()
			}

			r := gin.New()

			r.GET("/:alias", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedToURL, err := api.GetRedirect(c.url + "/" + c.alias)
			require.NoError(t, err)

			assert.Equal(t, c.url, redirectedToURL)

		})
	}
}
