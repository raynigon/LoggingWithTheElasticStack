#!/bin/bash
set -e

docker-compose build logstash
docker-compose up -d
echo "Prepare to configure Elasticsearch..."
host="localhost:9200"

docker-compose stop logstash
docker-compose stop filebeat

until $(curl --output /dev/null --silent --head --fail --max-time 5 "$host"); do
    printf '.'
    sleep 1
done
echo ""
# First wait for ES to start...
response="404"
until [ "$response" = "200" ]; do
    response=$(curl --write-out %{http_code} --silent --output /dev/null --max-time 5 "$host")
    echo "Elastic Search is unavailable - sleeping"
    sleep 1
done

# next wait for ES status to turn to Green
health="red"

until [ "$health" = 'green' ]; do
    health="$(curl --max-time 5 -fsSL "$host/_cat/health?h=status")"
    health="$(echo "$health" | sed -E 's/^[[:space:]]+|[[:space:]]+$//g')" # trim whitespace (otherwise we'll have "green ")
    echo "Elastic Search is unavailable - sleeping"
    sleep 1
done

echo "Elastic Search is up"

curl -X PUT "$host/_ilm/policy/elastictalk_policy" -H 'Content-Type: application/json' -d "@data/lifecycle-policy.json"
echo ""
curl -X PUT "$host/_template/elastictalk_template" -H 'Content-Type: application/json' -d "@data/index-template.json"
echo ""

# Configure Kibana
echo "Waiting for Kibana"

until $(curl --output /dev/null --silent --head --fail --max-time 5 "localhost:5601"); do
    printf '.'
    sleep 1
done
echo ""

curl -X POST "localhost:5601/api/saved_objects/index-pattern/elastictalk-*" -d "@data/index-pattern.json" -H 'Content-Type: application/json' -H 'kbn-xsrf: anything'
echo ""
curl -X POST "localhost:5601/api/saved_objects/visualization/access-logs-by-status" -d "@data/access-logs-by-status-graph.json" -H 'Content-Type: application/json' -H 'kbn-xsrf: anything'
echo ""
curl -X POST "localhost:5601/api/saved_objects/visualization/searchterms-frequency" -d "@data/search-terms-chart.json" -H 'Content-Type: application/json' -H 'kbn-xsrf: anything'
echo ""
curl -X POST "localhost:5601/api/saved_objects/visualization/browsers-chart" -d "@data/browsers-chart.json" -H 'Content-Type: application/json' -H 'kbn-xsrf: anything'
echo ""
curl -X POST "localhost:5601/api/saved_objects/visualization/users-map" -d "@data/users-map.json" -H 'Content-Type: application/json' -H 'kbn-xsrf: anything'
echo ""
curl -X POST "localhost:5601/api/saved_objects/dashboard/shop-dashboard" -d "@data/shop-dashboard.json" -H 'Content-Type: application/json' -H 'kbn-xsrf: anything'
echo ""

sleep 15
echo "Starting Logstash and Filebeat"
docker-compose start logstash
docker-compose start filebeat