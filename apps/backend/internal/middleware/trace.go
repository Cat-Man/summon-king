package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
)

type contextKey string

const (
	traceContextKey contextKey = "trace_id"
	traceIDHeader              = "X-Trace-ID"
)

func InjectTraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader(traceIDHeader)
		if traceID == "" {
			traceID = newTraceID()
		}

        c.Set(string(traceContextKey), traceID)
        c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), traceContextKey, traceID))
        c.Writer.Header().Set(traceIDHeader, traceID)
        c.Next()
	}
}

func GetTraceID(c *gin.Context) string {
    if traceID, _ := c.Request.Context().Value(traceContextKey).(string); traceID != "" {
        return traceID
    }
    if value, ok := c.Get(string(traceContextKey)); ok {
        if traceID, isString := value.(string); isString {
            return traceID
        }
	}
	return ""
}

func newTraceID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err == nil {
		return hex.EncodeToString(bytes[:])
	}
	return hex.EncodeToString([]byte(time.Now().Format("20060102150405.000000000")))
}
