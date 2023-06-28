package core

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func jsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			log := make(map[string]interface{})

			log["status"] = params.StatusCode
			log["path"] = params.Path
			log["method"] = params.Method
			//log["start_time"] = params.TimeStamp.Format("2006/01/02 - 15:04:05")
			//log["remote_addr"] = params.ClientIP
			log["latency"] = params.Latency.String()

			s, _ := json.Marshal(log)
			return "[Request] - " + string(s) + "\n"
		},
	)
}
