#!/bin/bash

set -e

dockerRId=$(uuidgen)
k8sRId=$(uuidgen)
kafkaRId=$(uuidgen)
helmRId=$(uuidgen)

dockerTId=$(uuidgen)
k8sTId=$(uuidgen)
kafkaTId=$(uuidgen)
helmTId=$(uuidgen)
backendTId=$(uuidgen)
infraTId=$(uuidgen)
messagingTId=$(uuidgen)

psqlScript=$(
	cat <<EOF
BEGIN;


INSERT INTO app.resource VALUES 
('$dockerRId', current_timestamp, current_timestamp, 'docker desc', 'https://www.docker.com/', 'BEGINNER', 'ARTICLE'),
('$k8sRId', current_timestamp, current_timestamp, 'k8s desc', 'https://kubernetes.io/', 'BEGINNER', 'ARTICLE'),
('$kafkaRId', current_timestamp, current_timestamp, 'kafka desc', 'https://kafka.apache.org/', 'ADVANCED', 'ARTICLE'),
('$helmRId', current_timestamp, current_timestamp, 'helm desc', 'https://www.youtube.com/results?search_query=helm+', 'BEGINNER', 'VIDEO');

INSERT INTO app.tag VALUES 
('$dockerTId', current_timestamp, current_timestamp, 'docker'),
('$k8sTId', current_timestamp, current_timestamp, 'k8s'),
('$kafkaTId', current_timestamp, current_timestamp, 'kafka'),
('$helmTId', current_timestamp, current_timestamp, 'helm'),
('$backendTId', current_timestamp, current_timestamp, 'backend'),
('$infraTId', current_timestamp, current_timestamp, 'infra'),
('$messagingTId', current_timestamp, current_timestamp, 'messaging');

INSERT INTO app.resources_tags VALUES 
('$(uuidgen)', current_timestamp, current_timestamp, '$dockerRId' , '$dockerTId'),
('$(uuidgen)', current_timestamp, current_timestamp, '$dockerRId' , '$infraTId'),
('$(uuidgen)', current_timestamp, current_timestamp, '$k8sRId' , '$k8sTId'),
('$(uuidgen)', current_timestamp, current_timestamp, '$k8sRId' , '$infraTId'),
('$(uuidgen)', current_timestamp, current_timestamp, '$helmRId' , '$helmTId'),
('$(uuidgen)', current_timestamp, current_timestamp, '$helmRId' , '$infraTId'),
('$(uuidgen)', current_timestamp, current_timestamp, '$kafkaRId' , '$kafkaTId'),
('$(uuidgen)', current_timestamp, current_timestamp, '$kafkaRId' , '$backendTId'),
('$(uuidgen)', current_timestamp, current_timestamp, '$kafkaRId' , '$messagingTId');

COMMIT;

EOF
)

echo $psqlScript >seed_data.sql
PGPASSWORD=password psql -h localhost -p 5432 -U user -d app -a -f seed_data.sql
rm seed_data.sql
echo 'Done seeding the db'
