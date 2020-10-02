[![Go Report Card](https://goreportcard.com/badge/github.com/Christian-Bull/slack-notify-go)](https://goreportcard.com/report/github.com/Christian-Bull/slack-notify-go)

# Slack Notify Go

This is a pretty simple app to post custom notifications to slack channels. As of now, this only pulls info from the twitch api. Notifying a slack channel when a user goes live. I've included some simple build instructions that I used to get this running on a k8s cluster.


# Info

Put your config into `/etc/slack-notify/config.json` using the structure in the config example.

Manifests are located in the ./manifests directory

# Build/deploy info:

## Clone Repo

`git clone https://github.com/Christian-Bull/slack-notify-go.git`

## Docker:

build the image

    docker build --tag <sometagname> . 

run it

    docker run -v /etc/slack-notify:/etc/slack-notify -d <sometagname>

You can also build this for multi-arch using buildx
    
    docker buildx build --platform linux/arm64,linux/amd64 -t <sometagname> .


## Kubernetes:

Set up the image on a registry (I used dockerhub)

I just retagged my image after the fact

    docker tag [image id] [registry:tag]                    # format
    docker tag 8922d588eec6 csbull55/slack-notify:initial   # what I did

Create namespace and deployment info

    kubectl create -f manifests/*

Inside the deployment, change the image (line 24) to your own registry link

You could also apply each manifest separately

    kubectl apply -f manifests/cbull-namespace.yaml
    kubectl apply -f manifests/slack-go-deploy.yaml

This should return

    pod/slack-go created

Lets see if it's running 

    kubectl get pod slack-go

    NAME           READY   STATUS    RESTARTS   AGE
    slack-go       1/1     Running   0          49m

When you mess up the initial deployment (definitely _not_ speaking from experience), k8s makes it super easy to clean up your mess.

    kubectl delete pod [name]
