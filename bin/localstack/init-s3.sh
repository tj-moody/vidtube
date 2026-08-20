#!/bin/env bash

awslocal s3 mb s3://vidtube-videos
awslocal s3api put-bucket-cors --bucket vidtube-videos --cors-configuration file:///etc/localstack/init/cors.json
