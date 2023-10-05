package types

import (
	"context"
	"encoding/json"
	"fmt"
	"livekit-webhook-proxy/utils"
	"net/http"

	"github.com/labstack/gommon/log"
	"github.com/prometheus/client_golang/prometheus"

	"cloud.google.com/go/pubsub"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type Proxy struct {
	pubsub  utils.PubSub
	Server  *echo.Echo
	Metrics struct {
		eventReceived  *prometheus.CounterVec
		eventPublished *prometheus.CounterVec
	}
}

func (p *Proxy) Init(ctx context.Context) {
	p.Server = echo.New()
	p.pubsub.Init(ctx)

	p.Server.HideBanner = true

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

	p.Metrics.eventReceived = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "livekit_webhook_proxy_event_received_total",
		Help: "The total number of events received",
	}, []string{"event", "room"})
	p.Metrics.eventPublished = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "livekit_webhook_proxy_event_published_total",
		Help: "The total number of events published",
	}, []string{"event", "room"})
}

func (p *Proxy) publish(c echo.Context) error {
	var payload map[string]interface{}
	if err := json.NewDecoder(c.Request().Body).Decode(&payload); err != nil {
		p.Server.Logger.Errorf("could not decode body: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "could not bind body").SetInternal(err)
	}
	p.Server.Logger.Infof("livekit event %s in room %s", payload["event"], payload["room"].(map[string]interface{})["name"])

	p.Metrics.eventReceived.WithLabelValues(payload["event"].(string), payload["room"].(map[string]interface{})["name"].(string)).Inc()

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

	p.Metrics.eventPublished.WithLabelValues(payload["event"].(string), payload["room"].(map[string]interface{})["name"].(string)).Inc()

	return c.JSON(http.StatusOK, payload)
}
