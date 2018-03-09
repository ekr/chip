// Lightly modified from Mint

package minq

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// We use this environment variable to control logging.  It should be a
// comma-separated list of log tags (see below) or "*" to enable all logging.
const logConfigVar = "MINQ_LOG"

// Pre-defined log types
const (
	logTypeAead        = "aead"
	logTypeCodec       = "codec"
	logTypeConnBuffer  = "connbuffer"
	logTypeConnection  = "connection"
	logTypeAck         = "ack"
	logTypeFrame       = "frame"
	logTypeHandshake   = "handshake"
	logTypeTls         = "tls"
	logTypeTrace       = "trace"
	logTypeServer      = "server"
	logTypeUdp         = "udp"
	logTypeStream      = "stream"
	logTypeFlowControl = "flow"
	logTypePacket      = "packet" // Just send notes on which packets are sent and received
	logTypeCongestion  = "congestion"
)

var (
	logFunction = log.Printf
	logAll      = false
	logSettings = map[string]bool{}
)

func init() {
	parseLogEnv(os.Environ())
}

func parseLogEnv(env []string) {
	for _, stmt := range env {
		if strings.HasPrefix(stmt, logConfigVar+"=") {
			val := stmt[len(logConfigVar)+1:]

			if val == "*" {
				logAll = true
			} else {
				for _, t := range strings.Split(val, ",") {
					logSettings[t] = true
				}
			}
		}
	}
}

func logf(tag string, format string, args ...interface{}) {
	if logAll || logSettings[tag] {
		fullFormat := fmt.Sprintf("[%s] %s", tag, format)
		logFunction(fullFormat, args...)
	}
}

type loggingFunction func(string, string, ...interface{})

func SetLogOutput(f func(string, ...interface{})) {
	logFunction = f
}

func newConnectionLogger(c *Connection) loggingFunction {
	return func(tag string, format string, args ...interface{}) {
		if logAll || logSettings[tag] {
			logf(tag, c.String()+": "+format, args...)
		}
	}
}

func newStreamLogger(id uint64, dir string, f loggingFunction) loggingFunction {
	extra := fmt.Sprintf("%s stream %d: ", dir, id)
	return func(tag string, format string, args ...interface{}) {
		f(tag, extra+format, args...)
	}
}
