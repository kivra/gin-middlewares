package correlationid

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-contrib/uuid"
)

// Header that holds the Correlation ID
var Header = "X-Correlation-Id"

// SetOnRequest controls whether to set Correlation IDs on request objects
var SetOnRequest = true

// SetOnResponse controls whether to set Correlation IDs on response objects
var SetOnResponse = true

// Middleware to set Correlation IDs on request and response objects
func Middleware(c *gin.Context) {
	correlationID := c.Request.Header.Get(Header)
	if correlationID == "" {
		correlationID = strings.ToUpper(uuid.NewV4().String())
	}
	if SetOnRequest {
		c.Request.Header.Set(Header, correlationID)
	}
	if SetOnResponse {
		c.Header(Header, correlationID)
	}
	c.Next()
}
