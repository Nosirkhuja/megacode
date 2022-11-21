package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// FiberPrometheus ...
type FiberPrometheus struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	statusCode      *prometheus.CounterVec
	userWent        *prometheus.CounterVec
	userExit        *prometheus.CounterVec
	defaultURL      string
}

func create(registry prometheus.Registerer, serviceName string) *FiberPrometheus {

	counter := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Name: serviceName + "_http_request_counter",
			Help: "counter request",
		},
		[]string{"path", "method"},
	)

	userwent := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Name: serviceName + "_user_went_counter",
			Help: "counter of users went",
		},
		[]string{"path", "method"},
	)

	userexit := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Name: serviceName + "_user_exit_counter",
			Help: "counter of users exit",
		},
		[]string{"path", "method"},
	)

	histogram := promauto.With(registry).NewHistogramVec(prometheus.HistogramOpts{
		Name: serviceName + "_http_total_request_execute_time",
		Help: "execute time counter for HTTP request",
	},
		[]string{"path", "method"},
	)

	statusCode := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Name: serviceName + "_http_total_request_status_code",
			Help: "status code counter for HTTP request",
		},
		[]string{"path", "method", "code"},
	)

	return &FiberPrometheus{
		requestsTotal:   counter,
		requestDuration: histogram,
		statusCode:      statusCode,
		userExit:        userexit,
		userWent:        userwent,
		defaultURL:      "/metrics",
	}
}

// New creates a new instance of FiberPrometheus middleware
// serviceName is available as a const label
func New(serviceName string) *FiberPrometheus {
	return create(prometheus.DefaultRegisterer, serviceName)
}

// Middleware is the actual default middleware implementation
func (ps *FiberPrometheus) Middleware(ctx *fiber.Ctx) error {

	start := time.Now()
	method := ctx.Route().Method

	if ctx.Route().Path == ps.defaultURL {
		return ctx.Next()
	}

	err := ctx.Next()
	// initialize with default error code
	// https://docs.gofiber.io/guide/error-handling
	status := fiber.StatusInternalServerError
	if err != nil {
		if e, ok := err.(*fiber.Error); ok {
			// Get correct error code from fiber.Error type
			status = e.Code
		}
	} else {
		status = ctx.Response().StatusCode()
	}

	path := ctx.Route().Path

	ps.requestsTotal.WithLabelValues(path, method).Inc()
	ps.statusCode.WithLabelValues(path, method, strconv.Itoa(status)).Inc()
	elapsed := float64(time.Since(start).Nanoseconds()) / 1e9
	ps.requestDuration.WithLabelValues(path, method).Observe(elapsed)

	return err
}
