package logger

import (
	"fmt"
	"time"

	"github.com/roninzo/log"
)

func traceErrMsg(source, sql string, rows int64, elapsed time.Duration, err error) (string, log.Map) {
	return fmt.Sprintf("%s %s", log.MesgGormTrace, err.Error()),
		log.Map{
			"source": source,
			"time":   getLatency(elapsed),
			"rows":   getRows(rows),
			"query":  sql,
		}
}

func traceWarnMsg(source, sql string, rows int64, elapsed time.Duration, slowThreshold time.Duration) (string, log.Map) {
	return fmt.Sprintf("%s SLOW SQL >= %v", log.MesgGormTrace, slowThreshold),
		log.Map{
			"source":    source,
			"time":      getLatency(elapsed),
			"rows":      getRows(rows),
			"threshold": slowThreshold.String(),
			"query":     sql,
		}
}

func traceMsg(source, sql string, rows int64, elapsed time.Duration) (string, log.Map) {
	return log.MesgGormTrace,
		log.Map{
			"source": source,
			"time":   getLatency(elapsed),
			"rows":   getRows(rows),
			"query":  sql,
		}
}

func getRows(rows int64) string {
	if rows == -1 {
		return "-"
	}
	return fmt.Sprintf("%d", rows)
}

func getLatency(elapsed time.Duration) string {
	return fmt.Sprintf("%fms", float64(elapsed.Nanoseconds())/1e6)
}
