#!/bin/sh

TOTALRAMUSED=$(openstack limits show --absolute -f json | jq '.[] | select(.Name=="totalRAMUsed") | .Value')

echo '{"labels": {"tenant_id": "'$OS_PROJECT_ID'","project_name": "'$OS_PROJECT_NAME'"}, "results": {"items": '$TOTALRAMUSED'} }'

exit 0
