package util

import (
	"ginorm/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"strings"
)

func GetAbsPath(relPath string) string {
	if config.Env.Platform == config.Windows {
		return config.Env.RootDir + strings.Replace(relPath, "\\", "/", -1)
	} else {
		return config.Env.RootDir + strings.Replace(relPath, "\\", "/", -1)
	}
}

// GenerateTraceID 生成一个随机的 Trace ID
func GenerateTraceID() string {
	traceID, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("failed to generate uuid: %v", err)
	}
	return traceID.String()
}

// GetTraceID 从 context 中获取 Trace ID
func GetTraceID(c *gin.Context) string {
	if traceID, ok := c.Value("trace_id").(string); ok {
		return traceID
	}
	return "unknown"
}
