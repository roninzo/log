package logger

import (
	"io"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// Config defines the config for middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// // Format defines the logging tags
	// //
	// // Optional. Default: [${time}] ${status} - ${latency} ${method} ${path}\n
	// Format string

	// TimeFormat https://programming.guide/go/format-parse-string-time-date-example.html
	//
	// Optional. Default: 15:04:05
	TimeFormat string

	// TimeZone can be specified, such as "UTC" and "America/New_York" and "Asia/Chongqing", etc
	//
	// Optional. Default: "Local"
	TimeZone string

	// TimeInterval is the delay before the timestamp is updated
	//
	// Optional. Default: 500 * time.Millisecond
	TimeInterval time.Duration

	// Output is a writter where logs are written
	//
	// Default: os.Stderr
	Output io.Writer

	// Flag defines what to log.
	//
	// Default: LstdFlags (Lpid | LIP | Luid | Lstatus | LbytesSent | Llatency | Lmethod | Lpath)
	Flag int

	// Context Key for User ID
	//
	// Default: "userid"
	ContextKeyUID string

	// Context Key for Request ID
	//
	// Default: "requestid"
	ContextKeyRID string

	// enableColors     bool
	// enableLatency    bool
	timeZoneLocation *time.Location
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	// Format:       "[${time}] ${status} - ${latency} ${method} ${path}\n",
	// enableColors: true,
	Next:          nil,
	TimeFormat:    "15:04:05",
	TimeZone:      "Local",
	ContextKeyUID: "userid",
	ContextKeyRID: requestid.ConfigDefault.ContextKey, // "requestid",
	TimeInterval:  500 * time.Millisecond,
	Output:        os.Stderr,
	Flag:          LstdFlags,
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// // Enable colors if no custom format or output is given
	// if validCustomFormat(cfg.Format) && cfg.Output == nil {
	// 	cfg.enableColors = true
	// }

	// Set default values
	// if cfg.Format == "" {
	// 	cfg.Format = ConfigDefault.Format
	// }
	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}
	if cfg.TimeZone == "" {
		cfg.TimeZone = ConfigDefault.TimeZone
	}
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = ConfigDefault.TimeFormat
	}
	if int(cfg.TimeInterval) <= 0 {
		cfg.TimeInterval = ConfigDefault.TimeInterval
	}
	if cfg.Flag == 0 {
		cfg.Flag = ConfigDefault.Flag
	}
	if cfg.ContextKeyUID == "" {
		cfg.TimeFormat = ConfigDefault.ContextKeyUID
	}
	if cfg.ContextKeyRID == "" {
		cfg.TimeFormat = ConfigDefault.ContextKeyRID
	}
	if cfg.Output == nil {
		cfg.Output = ConfigDefault.Output
	}
	return cfg
}

// // Function to check if the logger format is compatible for coloring
// func validCustomFormat(format string) bool {
// 	validTemplates := []string{"${status}", "${method}"}
// 	if format == "" {
// 		return true
// 	}
// 	for _, template := range validTemplates {
// 		if !strings.Contains(format, template) {
// 			return false
// 		}
// 	}
// 	return true
// }
