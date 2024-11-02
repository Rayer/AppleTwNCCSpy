#!/bin/sh
gcloud functions deploy AppleTWNCCCrawler --gen2 --region=us-west1 --runtime=go122 --memory 128Mi --source=. --entry-point=AppleProductMonitor --trigger-topic=AppleTWNccSpyScheduler