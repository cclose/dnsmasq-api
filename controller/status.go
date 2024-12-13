package controller

import (
	"bytes"
	"github.com/cclose/dnsmasq-api/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

type IStatusController interface {
	GetStatus(c echo.Context) error
	Register(e *echo.Echo)
}

type StatusController struct {
	BuildInfo model.BuildInfo

	startTime time.Time
}

func NewStatusController(bi model.BuildInfo) *StatusController {
	return &StatusController{
		BuildInfo: bi,
		startTime: time.Now(),
	}
}

func (sc *StatusController) Register(e *echo.Echo) {
	e.GET("/statusz", sc.GetStatus)
	e.GET("/metricz", sc.GetMetrics)
}

func (sc *StatusController) GetStatus(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"status":     "running",
		"time":       time.Now().Format(time.RFC3339),
		"uptime":     time.Since(sc.startTime).String(),
		"version":    sc.BuildInfo.Version,
		"commit":     sc.BuildInfo.Commit,
		"build_time": sc.BuildInfo.BuildTime.Format(time.RFC3339),
	})
}

func MetricsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqData := map[string]string{
				"method": c.Request().Method,
				"path":   c.Path(),
			}

			durCount := metrics.GetOrCreateHistogram("http_requests_duration" + buildLabels(reqData))

			start := time.Now()
			err := next(c)
			durCount.UpdateDuration(start)

			// Add status code to req data
			reqData["status"] = strconv.Itoa(c.Response().Status)

			// Increment the total request counter
			totCount := metrics.GetOrCreateCounter("http_requests_total" + buildLabels(reqData))
			totCount.Inc()

			return err
		}
	}
}

func (sc *StatusController) GetMetrics(c echo.Context) error {
	buf := new(bytes.Buffer)
	metrics.WritePrometheus(buf, true)
	return c.String(http.StatusOK, buf.String())
}

// Adapted from a Library by tstuyfzand
func buildLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteByte('{')
	for k, v := range labels {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		b.WriteString(k)
		b.WriteString(`="`)
		b.WriteString(strings.Replace(v, "\"", "\\", -1))
		b.WriteByte('"')
	}
	b.WriteString(`}`)
	return b.String()
}
