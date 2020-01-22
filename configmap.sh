#! /bin/bash
# this creates the configmap used in the cluster
# super simple so I didn't bother with yaml

kubectl create configmap -n=cbull slack-config --from-file=/etc/slack-notify/config.json
