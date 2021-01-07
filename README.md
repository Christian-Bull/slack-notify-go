[![Go Report Card](https://goreportcard.com/badge/github.com/Christian-Bull/slack-notify-go)](https://goreportcard.com/report/github.com/Christian-Bull/slack-notify-go)

# Slack Notify Go

This is a pretty simple app to post custom notifications to slack channels. As of now, this only pulls info from the twitch api. Notifying a slack channel when a user goes live. I've included some simple build instructions that I used to get this running on a k8s cluster.


# Info

Manifests are located in the ./manifests directory

# Build/deploy info:

## Clone Repo

`git clone https://github.com/Christian-Bull/slack-notify-go.git`

## Docker:

build the image

    docker build --tag <sometagname> . 

run it

     docker run -v /etc/slack-notify:/etc/slack-notify -e CONFIG_PATH="/etc/slack-notify/config.json" -d <sometagname>

You can also build this for multi-arch using buildx
    
    docker buildx build --platform linux/arm64,linux/amd64 -t [registry:tagname] --push .

This is a pretty comprehensive guide to setting up and building golang docker images with buildx https://medium.com/@artur.klauser/building-multi-architecture-docker-images-with-buildx-27d80f7e2408

## Kubernetes:

Set up the image on a registry (I used dockerhub) - if using buildx this is not needed since it's done when building the image

I just retagged my image and pushed it after the fact

    docker tag [image id] [registry:tag]                    # format
    docker tag 8922d588eec6 csbull55/slack-notify:initial   # what I did
    docker push csbull55/slack-notify:<version>

Create configmap

    kubectl create configmap -n=<namespace> slack-config --from-file=/etc/slack-notify/config.json


Create namespace and deployment info

    kubectl create -f manifests/*

Should be up and running, lets check

```cbull@cbull:~$ kubectl get pods -n=cbull
NAME                       READY   STATUS    RESTARTS   AGE
slack-go-ddf77c8d6-zjrb7   1/1     Running   0          8m1s
```

### Argocd:
```project: default
source:
  repoURL: 'https://github.com/Christian-Bull/slack-notify-go'
  path: manifests
  targetRevision: HEAD
destination:
  server: 'https://kubernetes.default.svc'
  namespace: cbull
syncPolicy:
  automated: {}
```
9 lines of yaml and my CD is setup. When I make changes to my manifests in my repo, argo automatically syncs that state with the cluster and redeploys. `syncPolicy: automated`