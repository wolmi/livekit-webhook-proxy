package types

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"livekit-webhook-proxy/utils"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/gommon/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type Proxy struct {
	pubsub  utils.PubSub
	Server  *echo.Echo
	Metrics struct {
		EventsReceived  *prometheus.CounterVec
		EventsPublished *prometheus.CounterVec
	}
}

func (p *Proxy) Init(ctx context.Context, metrics bool, port int) {
	p.Server = echo.New()
	p.pubsub.Init(ctx)

	p.Server.HideBanner = true

	if metrics {
		p.Metrics.EventsReceived = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "livekit_webhook_proxy_event_received_total",
			Help: "The total number of events received",
		}, []string{"event", "room"})
		p.Metrics.EventsPublished = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "livekit_webhook_proxy_event_published_total",
			Help: "The total number of events published",
		}, []string{"event", "room"})

		prometheus.MustRegister(p.Metrics.EventsReceived)
		prometheus.MustRegister(p.Metrics.EventsPublished)

		prometheus.Unregister(collectors.NewGoCollector())

		go func() {
			metrics := echo.New() // this Echo will run on separate port 8081
			metrics.HideBanner = true
			metrics.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics
			if err := metrics.Start(fmt.Sprintf(":%d", port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}()

	}
	// get log level from flag
	logLevel := viper.GetBool("debug")
	if logLevel {
		p.Server.Logger.SetLevel(log.DEBUG)
	} else {
		p.Server.Logger.SetLevel(log.INFO)
	}

	p.Server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	p.Server.POST("/publish", p.publish)
	p.Server.Logger.Fatal(p.Server.Start(fmt.Sprintf(":%d", viper.GetInt("port"))))

}

func (p *Proxy) publish(c echo.Context) error {
	var payload map[string]interface{}
	if err := json.NewDecoder(c.Request().Body).Decode(&payload); err != nil {
		p.Server.Logger.Errorf("could not decode body: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "could not bind body").SetInternal(err)
	}
	p.Server.Logger.Infof("livekit event %s in room %s", payload["event"], payload["room"].(map[string]interface{})["name"])

	p.Metrics.EventsReceived.With(prometheus.Labels{"event": fmt.Sprintf("%v", payload["event"]), "room": fmt.Sprintf("%v", payload["room"].(map[string]interface{})["name"])}).Inc()

	jsonPayload, _ := json.Marshal(payload)
	p.Server.Logger.Debugf("event payload data: %s", jsonPayload)

	topic := p.pubsub.Client.Topic(viper.GetString("topic"))
	res := topic.Publish(c.Request().Context(), &pubsub.Message{
		Data: jsonPayload,
	})
	msgID, err := res.Get(c.Request().Context())
	if err != nil {
		p.Server.Logger.Fatal(err)
	}
	p.Server.Logger.Debugf("event published with msgID %v", msgID)

	p.Metrics.EventsPublished.With(prometheus.Labels{"event": fmt.Sprintf("%v", payload["event"]), "room": fmt.Sprintf("%v", payload["room"].(map[string]interface{})["name"])}).Inc()

	return c.JSON(http.StatusOK, payload)
}
