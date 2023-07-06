package utils

import (
	"net/http"

	"cloud.google.com/go/compute/metadata"
)

type MetaData struct {
	httpClient *http.Client
	client     *metadata.Client
	OnGCE      bool
}

func (m *MetaData) Init() {
	m.httpClient = http.DefaultClient
	m.client = metadata.NewClient(m.httpClient)
	m.OnGCE = metadata.OnGCE()
}
