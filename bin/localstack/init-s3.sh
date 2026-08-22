#!/bin/env bash

awslocal s3 mb s3://vidtube-videos-1
awslocal s3api put-bucket-cors --bucket vidtube-videos-1 --cors-configuration file:///etc/localstack/init/cors.json
