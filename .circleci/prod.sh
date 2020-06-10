#! /bin/bash -e

# Push the new docker image
docker push eu.gcr.io/the-treep-api-1525507752734/covidtracker
docker push eu.gcr.io/the-treep-api-1525507752734/covidtracker-refresher

# Replace deployment image
kubectl set image deployment/covidtracker covidtracker=eu.gcr.io/the-treep-api-1525507752734/covidtracker:$CIRCLE_SHA1
kubectl set image deployment/covidtracker-refresher covidtracker-refresher=eu.gcr.io/the-treep-api-1525507752734/covidtracker-refresher:$CIRCLE_SHA1
