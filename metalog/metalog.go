package metalog

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type key struct {
	namespace string
}

// Namespace in gin context reserved for metalog objects
const Namespace = "gin/metalog"

// TimeFormat used to format time in logs
var TimeFormat = time.RFC3339

// TimeKey that holds the log time in logs
var TimeKey = "timestamp"

// LevelFormatter to format log level strings like "info" and "error"
var LevelFormatter = strings.ToUpper

// TimeUnit for request durations
var TimeUnit = time.Second

// Add metadata to a request context
func Add(moreData map[string]interface{}, c *gin.Context) {
	metaData := Get(c.Request)
	for key, value := range moreData {
		metaData[key] = value
	}
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), newKey(), metaData))
}

// Add a message to the metalog context
func AddMessage(message interface{}, c *gin.Context) {
	Add(gin.H{"message": message}, c)
}

// Add an error to the metalog context
func AddError(err interface{}, c *gin.Context) {
	Add(gin.H{"error": err}, c)
}

// Get metadata from the request context
func Get(req *http.Request) map[string]interface{} {
	metaData, ok := req.Context().Value(newKey()).(map[string]interface{})
	if ok {
		return metaData
	}
	return make(map[string]interface{})
}

func newKey() key {
	return key{namespace: Namespace}
}
