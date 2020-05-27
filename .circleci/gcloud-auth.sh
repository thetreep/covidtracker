#! /bin/bash -e

echo ${GOOGLE_AUTH} > ${HOME}/gcp-key.json
gcloud auth activate-service-account --key-file ${HOME}/gcp-key.json
gcloud --quiet config set project the-treep-api-1525507752734
gcloud config set compute/zone europe-west1-c
gcloud auth configure-docker --quiet

# Setup basic gcloud config
gcloud --quiet config set container/cluster thetreep-api-cluster
gcloud --quiet container clusters get-credentials thetreep-api-cluster

# setting current k8s context namespace
kubectl config set-context $(kubectl config current-context)

# Display k8s config
kubectl config view
kubectl config current-context
