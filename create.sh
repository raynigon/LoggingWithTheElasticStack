#!/bin/bash
cd elastic-stack/
./build-and-run.sh
sleep 15 # Wait Until Logstash and Filebeat are up and running
cd ../shop
docker-compose up