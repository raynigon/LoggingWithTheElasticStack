FROM docker.elastic.co/logstash/logstash:7.3.0

# RUN bin/logstash-plugin install logstash-filter-json_encode

COPY logstash.yml config/
COPY pipeline/ pipeline/