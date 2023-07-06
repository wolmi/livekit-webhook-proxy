package types

import (
	"context"
	"encoding/json"
	"fmt"
	"livekit-webhook-proxy/utils"
	"log"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type Proxy struct {
	pubsub utils.PubSub
	echo   *echo.Echo
}

func (p *Proxy) Init(ctx context.Context) {
	p.echo = echo.New()
	p.pubsub.Init(ctx)

	p.echo.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	p.echo.POST("/publish", p.publish)
	p.echo.Logger.Fatal(p.echo.Start(fmt.Sprintf(":%d", viper.GetInt("port"))))
}

func (p *Proxy) publish(c echo.Context) error {
	var payload map[string]interface{}
	if err := json.NewDecoder(c.Request().Body).Decode(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "could not bind body").SetInternal(err)
	}
	log.Printf("[INFO] livekit event %s in room %s", payload["event"], payload["room"].(map[string]interface{})["name"])
	jsonPayload, _ := json.Marshal(payload)
	log.Printf("[DEBUG] event payload data: %s", jsonPayload)

	topic := p.pubsub.Client.Topic(viper.GetString("topic"))
	res := topic.Publish(c.Request().Context(), &pubsub.Message{
		Data: jsonPayload,
	})
	msgID, err := res.Get(c.Request().Context())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[DEBUG] event published with msgID %v", msgID)
	return c.JSON(http.StatusOK, payload)
}