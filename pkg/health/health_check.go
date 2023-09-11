package health

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/hellofresh/health-go/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

const serviceName = "chat-api"
const serviceVersion = "v1.0.0"

const dbTimeout = time.Second * 5
const brokerTimeout = time.Second * 5

type Health interface {
	Handler() http.Handler
}

type HealthCheck struct {
	health *health.Health
}

func NewHealthCheck(db *sql.DB, conn *amqp.Connection) *HealthCheck {
	h, _ := health.New(health.WithComponent(health.Component{
		Name:    serviceName,
		Version: serviceVersion,
	}))

	h.Register(health.Config{
		Name:    "database",
		Timeout: dbTimeout,
		Check: func(context.Context) error {
			if err := db.Ping(); err != nil {
				return errors.New("connection refused")
			}

			return nil
		},
	})

	h.Register(health.Config{
		Name:    "broker",
		Timeout: brokerTimeout,
		Check: func(context.Context) error {
			if conn.IsClosed() {
				return errors.New("connection closed")
			}

			return nil
		},
	})

	return &HealthCheck{
		health: h,
	}
}

func (h *HealthCheck) Handler() http.Handler {
	return h.health.Handler()
}
