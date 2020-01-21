#! /bin/bash
# this creates the configmap used in the cluster

kube create configmap -n=cbull slack-config --from-file=assets/config.json
