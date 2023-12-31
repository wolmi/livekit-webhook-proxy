# livekit-webhook-proxy
This is a proxy to allow sending Livekit webhook events to a GCP PubSub topic

# config.yaml example

```yaml
port: 8080
topic: <pubsub-topic>
project-id: <gcp-project-id>
```

# How to use it

The simplest way to authenticate to GCP PubSub is to use a service account with the proper permission.
You can use the environment variable `GOOGLE_APPLICATION_CREDENTIALS` to specify where the service account json file is located, inside a GKE cluster or GCE instance it will be detected atomatically.

```bash
export GOOGLE_APPLICATION_CREDENTIALS=sa.json
livekit-webhook-proxy --config config.yaml
```
# Configure your Livekit server

Once you have the proxy running and accessible by your Livekit server you have to activate the webhook and send it to the `/publish` endpoint

```yaml
webhooks:
  api_key: <key>
  urls:
    - "https://webhook-proxy:8080/publish"
```

# Test with docker-compose

Fist you have to set the right config in the [livekit.yaml](livekit.yaml) config file to set the `key:secret` and also create a [config.yaml](#configyaml-example)

```bash
docker-compose up
```

You will need to generate a valid token with your `key:secret` and you hace test it in the [React example](https://example.livekit.io/#/) from Livekit.

In the environment variables of the proxy it assumes you have `gcloud` properly configured with de `application_default_credentials.json` located in the default home folder. If you want to use a different service account or it's not in the right location plase eddit the [docker-compose.yaml](docker-compose.yaml)