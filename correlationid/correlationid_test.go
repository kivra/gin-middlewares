package correlationid

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetCorrelationIdOnResponse(t *testing.T) {
	req, res, engine := setup()
	engine.ServeHTTP(res, req)

	if res.Header().Get(Header) == "" {
		t.Fatalf("expected non-empty correlation header")
	}
}

func TestDisableCorrelationIdOnResponse(t *testing.T) {
	SetOnResponse = false
	req, res, engine := setup()
	engine.ServeHTTP(res, req)

	correlationID := res.Header().Get(Header)
	if correlationID != "" {
		t.Fatalf("got correlation ID %s, but expected none", correlationID)
	}
}

func TestUseExistingCorrelationIdIfProvided(t *testing.T) {
	req, res, engine := setup()
	req.Header.Set(Header, "hello")
	engine.ServeHTTP(res, req)

	correlationID := res.Header().Get(Header)
	if correlationID != "hello" {
		t.Fatalf("expected existing correlation ID to be returned, got %s", correlationID)
	}
}

func setup() (*http.Request, *httptest.ResponseRecorder, *gin.Engine) {
	engine := gin.New()
	engine.GET("/test", Middleware, func(c *gin.Context) { c.Status(200) })

	req := httptest.NewRequest(http.MethodGet, "/test", new(bytes.Buffer))
	res := httptest.NewRecorder()
	return req, res, engine
}
