#! /bin/bash -e

# Push the new docker image
docker push eu.gcr.io/the-treep-api-1525507752734/covidtracker

# Replace deployment image
kubectl set image deployment/covidtracker covidtracker=eu.gcr.io/the-treep-api-1525507752734/covidtracker:$CIRCLE_SHA1
kubectl set image deployment/billing-cron billing-cron=eu.gcr.io/the-treep-api-1525507752734/billing-cron:$CIRCLE_SHA1
