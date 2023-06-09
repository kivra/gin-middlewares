package metalog

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const levelINFO = "INFO"

const levelERROR = "ERROR"

func LogFormatter(param gin.LogFormatterParams) string { // nolint
	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency -= param.Latency % time.Second
	}

	logMessage := Get(param.Request)

	var level string
	if lev, ok := logMessage["level"]; ok {
		level = lev.(string)
	} else if _, ok := logMessage["error"]; ok || param.StatusCode >= 500 {
		level = levelERROR
	} else if param.ErrorMessage != "" {
		level = levelERROR
		logMessage["error"] = strings.TrimSuffix(param.ErrorMessage, "\n")
	} else {
		level = levelINFO
	}

	logMessage["client_ip"] = param.ClientIP
	logMessage["duration"] = float64(param.Latency.Nanoseconds()) / float64(TimeUnit.Nanoseconds())
	logMessage["headers"] = format(param.Request.Header)
	logMessage["level"] = LevelFormatter(level)
	logMessage["method"] = param.Method
	logMessage["module"] = "GIN"
	logMessage["path"] = param.Request.URL.Path
	logMessage["query"] = param.Request.URL.Query()
	logMessage["response_body_size"] = param.BodySize
	logMessage["status"] = param.StatusCode
	logMessage[TimeKey] = param.TimeStamp.Format(TimeFormat)

	jsonLog, err := json.Marshal(logMessage)

	if err != nil {
		return fmt.Sprintf(
			"{\"level\": \"%q\", \"module\": \"metalog\", \"message\": \"failed to convert log message to JSON\"}\n",
			LevelFormatter(levelERROR),
		)
	}

	return string(jsonLog) + "\n"
}

func format(m map[string][]string) map[string]string {
	formatted := make(map[string]string)
	for key, values := range m {
		formatted[key] = strings.Join(values, ", ")
	}
	return formatted
}
