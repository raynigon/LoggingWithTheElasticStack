package utils

import (
	"elastic-talk-search/config"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	// Field Configuration
	zerolog.TimestampFieldName = "@timestamp"
	zerolog.MessageFieldName = "msg"
	zerolog.LevelFieldName = "loglevel"
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.999Z"
	// Log Level Configuration
	config := config.Get()
	logLevelStr := config.Log.Level
	switch {
	case logLevelStr == "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		break
	case logLevelStr == "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		break
	case logLevelStr == "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
		break
	case logLevelStr == "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		break
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

// NewLogger creates a new Logger Pointer
func NewLogger() *zerolog.Logger {
	config := config.Get()
	sublogger := zerolog.New(os.Stdout).
		With().
		Stack().
		Str("service", config.Application.Name).
		Str("application_type", "service").
		Logger()
	return &sublogger
}

// NewApplicationLogger creates a new Application Logger Pointer
func NewApplicationLogger(logger *zerolog.Logger) *zerolog.Logger {
	sublogger := logger.
		With().
		Str("log_type", "application").
		Logger()
	return &sublogger
}

// NewLoggerWithCorrelationID creates a new Logger Pointer with the given correlaction id
func NewLoggerWithCorrelationID(correlationID string) *zerolog.Logger {
	sublogger := NewLogger().
		With().
		Str("correlation-id", correlationID).
		Logger()
	return &sublogger
}

// NewRecoveryHandler creates a panic recovery error logger
func NewRecoveryHandler(sublog *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				sublog.Error().
					Int("status", 500).
					Str("request_method", c.Request.Method).
					Str("uri", c.Request.URL.Path).
					Str("client_ip", c.ClientIP()).
					Str("remote_address", c.ClientIP()).
					Str("agent", c.Request.UserAgent()).
					Str("correlation-id", c.GetHeader("correlation-id")).
					Str("log_type", "application").
					AnErr("payload", err.(error)).
					Msgf("Panic Occured: %v", err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// NewAccessLogger creates a new Access Logger which will log all requests to stdout
func NewAccessLogger(sublog *zerolog.Logger) gin.HandlerFunc {
	skip := make(map[string]struct{}, 2)
	skip["/admin/health"] = struct{}{}
	skip["/admin/healthcheck"] = struct{}{}
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		c.Next()
		track := true

		if _, ok := skip[path]; ok {
			track = false
		}
		if track {
			end := time.Now().UTC()
			latency := end.Sub(start)

			msg := "Request"
			if len(c.Errors) > 0 {
				msg = c.Errors.String()
			}

			dumplogger := sublog.With().
				Int("status", c.Writer.Status()).
				Str("request_method", c.Request.Method).
				Str("uri", path).
				Str("client_ip", c.ClientIP()).
				Str("remote_address", c.ClientIP()).
				Dur("response_time", latency).
				Str("agent", c.Request.UserAgent()).
				Str("correlation-id", c.GetHeader("correlation-id")).
				Str("log_type", "access").
				Logger()
			switch {
			case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
				{
					dumplogger.Warn().
						Msg(msg)
				}
			case c.Writer.Status() >= http.StatusInternalServerError:
				{
					dumplogger.Error().
						Msg(msg)
				}
			default:
				dumplogger.Info().
					Msg(msg)
			}
		}
	}
}
