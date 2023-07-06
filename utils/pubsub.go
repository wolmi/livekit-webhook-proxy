package utils

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/spf13/viper"
)

type PubSub struct {
	Client *pubsub.Client
}

func (ps *PubSub) Init(ctx context.Context) {
	metadata := MetaData{}

	metadata.Init()
	project_id := viper.GetString("project-id")
	if metadata.OnGCE {
		project_id, _ = metadata.client.ProjectID()
	}

	client, err := pubsub.NewClient(ctx, project_id)
	if err != nil {
		log.Fatal(err)
	}

	ps.Client = client

}
