// Borrowed heavily from:
// https://github.com/gofiber/fiber/tree/master/middleware/logger

package logger

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/roninzo/log"
)

// Fiber logger via composition.
type Logger struct {
	log.Logger
	Config
}

// Fiber logging handler.
func New(lgr log.Logger, config ...Config) func(c *fiber.Ctx) error {
	// Set default config
	cfg := configDefault(config...)

	// Current logger.
	l := &Logger{lgr, cfg}

	// Get timezone location
	tz, err := time.LoadLocation(l.TimeZone)
	if err != nil || tz == nil {
		l.timeZoneLocation = time.Local
	} else {
		l.timeZoneLocation = tz
	}

	// Set PID once
	pid := strconv.Itoa(os.Getpid())

	// Tag overrides from Config.
	Tags[Lrid] = l.ContextKeyRID
	Tags[Luid] = l.ContextKeyUID

	// Set variables
	var (
		once       sync.Once
		errHandler fiber.ErrorHandler
	)

	// Return new handler
	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if l.Next != nil && l.Next(c) {
			return c.Next()
		}

		// Set error handler once
		once.Do(func() {
			// override error handler
			errHandler = c.App().Config().ErrorHandler
		})

		var start, stop time.Time

		// Set latency start time
		if l.Flag&Llatency == 0 {
			start = time.Now()
		}

		// Handle request, store err for logging
		chainErr := c.Next()

		// Manually call error handler
		if chainErr != nil {
			if err := errHandler(c, chainErr); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		// Set latency stop time
		if l.Flag&Llatency == 0 {
			stop = time.Now()
		}

		// This is our custom list of arguments to log
		fields := log.Map{}
		{
			var key string

			// LstdFlags
			{
				if l.Flag&Lpid == 0 {
					key = Tags[Lpid]
					fields[key] = pid
				}
				if l.Flag&LIP == 0 {
					key = Tags[LIP]
					fields[key] = c.IP()
				}
				if l.Flag&Luid == 0 {
					key = Tags[Luid]
					fields[key] = c.Locals(key)
				}
				if l.Flag&Lstatus == 0 {
					key = Tags[Lstatus]
					fields[key] = c.Response().StatusCode()
				}
				if l.Flag&LbytesSent == 0 {
					key = Tags[LbytesSent]
					fields[key] = len(c.Response().Body())
				}
				if l.Flag&Llatency == 0 {
					key = Tags[Llatency]
					fields[key] = stop.Sub(start).String()
				}
				if l.Flag&Lmethod == 0 {
					key = Tags[Lmethod]
					fields[key] = c.Method()
				}
				if l.Flag&Lpath == 0 {
					key = Tags[Lpath]
					fields[key] = c.Path()
				}
			}

			// Non-LstdFlags
			{
				if l.Flag&Ltime == 0 {
					key = Tags[Ltime]
					fields[key] = time.Now().In(l.timeZoneLocation).Format(l.TimeFormat)
				}
				if l.Flag&Lreferer == 0 {
					key = Tags[Lreferer]
					fields[key] = c.Get(fiber.HeaderReferer)
				}
				if l.Flag&Lprotocol == 0 {
					key = Tags[Lprotocol]
					fields[key] = c.Protocol()
				}
				if l.Flag&Lport == 0 {
					key = Tags[Lport]
					fields[key] = c.Port()
				}
				if l.Flag&LIPs == 0 {
					key = Tags[LIPs]
					fields[key] = c.Get(fiber.HeaderXForwardedFor)
				}
				if l.Flag&Lhost == 0 {
					key = Tags[Lhost]
					fields[key] = c.Hostname()
				}
				if l.Flag&LURL == 0 {
					key = Tags[LURL]
					fields[key] = c.OriginalURL()
				}
				if l.Flag&LUA == 0 {
					key = Tags[LUA]
					fields[key] = c.Get(fiber.HeaderUserAgent)
				}
				if l.Flag&LresBody == 0 {
					key = Tags[LresBody]
					fields[key] = c.Response().Body()
				}
				if l.Flag&LqueryStringParams == 0 {
					key = Tags[LqueryStringParams]
					fields[key] = c.Request().URI().QueryArgs().String()
				}
				if l.Flag&Lbody == 0 {
					key = Tags[Lbody]
					fields[key] = c.Body()
				}
				if l.Flag&LbytesReceived == 0 {
					key = Tags[LbytesReceived]
					fields[key] = len(c.Request().Body())
				}
				if l.Flag&Lroute == 0 {
					key = Tags[Lroute]
					fields[key] = c.Route().Path
				}
				if l.Flag&Lerror == 0 {
					key = Tags[Lerror]
					fields[key] = "-"

					// Error if exist
					if chainErr != nil {
						fields[key] = chainErr.Error()
					}
				}
				if l.Flag&Lrid == 0 {
					key = Tags[Lrid]
					fields[key] = c.Locals(key)
				}
			}
		}

		l.Info(log.MesgFiberLogger, fields)
		return nil
	}
}
