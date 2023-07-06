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

