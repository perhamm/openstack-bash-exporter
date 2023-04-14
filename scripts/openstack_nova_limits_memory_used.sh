#!/bin/sh

TOTALRAMUSED=$(openstack --os-cloud=openstack limits show --absolute -f json | jq '.[] | select(.Name=="totalRAMUsed") | .Value')

TENANT_ID=$(cat clouds.yaml | shyaml get-value clouds.openstack.auth.project_id)

PROJECT_NAME=$(cat clouds.yaml | shyaml get-value clouds.openstack.auth.project_name)

echo '{"labels": {"tenant_id": "'$TENANT_ID'","project_name": "'$PROJECT_NAME'"}, "results": {"items": '$TOTALRAMUSED'} }'

exit 0
