#!/bin/env bash

READY_DIR=/etc/localstack/init/ready.d

awslocal s3 mb s3://vidtube-videos-1
awslocal s3api put-bucket-cors \
  --bucket vidtube-videos-1 \
  --cors-configuration file://$READY_DIR/cors.json

QUEUE_URL=$(awslocal sqs create-queue --queue-name vidtube-uploads --query QueueUrl --output text)
echo "created queue: $QUEUE_URL"

QUEUE_ATTRIBUTES=$(python3 -c '
import json, sys
with open(sys.argv[1]) as f:
    print(json.dumps({"Policy": json.dumps(json.load(f))}))
' "$READY_DIR/queue-policy.json")

awslocal sqs set-queue-attributes \
  --queue-url "$QUEUE_URL" \
  --attributes "$QUEUE_ATTRIBUTES"

awslocal s3api put-bucket-notification-configuration \
  --bucket vidtube-videos-1 \
  --notification-configuration file://$READY_DIR/bucket-notifications.json
