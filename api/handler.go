package api

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/samarthasthan/luganodes-task/internal/store/controller"
	"github.com/samarthasthan/luganodes-task/pkg/logger"
	"github.com/sirupsen/logrus"
)

var (
	// Request counters
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_count",
			Help: "Number of requests received",
		},
		[]string{"path"},
	)

	// Request duration histograms
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

type (
	Handlers struct {
		*echo.Echo
		controller *controller.Controller
		log        *logger.Logger
		tracer     *zipkin.Tracer
	}
)

func init() {
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)
}

// NewHandler creates a new handler
func NewHandler(c *controller.Controller, l *logger.Logger, t *zipkin.Tracer) *Handlers {
	return &Handlers{Echo: echo.New(), controller: c, log: l, tracer: t}
}

// Handle handles the routes
func (h *Handlers) Handle() {
	h.Use(echoprometheus.NewMiddleware("apigateway")) // adds middleware to gather metrics
	h.GET("/metrics", echoprometheus.NewHandler())
	h.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			h.log.WithFields(logrus.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("request")

			return nil
		},
	})) // adds route to serve gathered metrics
	h.GET("/deposits", h.controller.GetDeposit)
}
