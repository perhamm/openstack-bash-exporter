#!/bin/sh

MAXTOTALCORES=$(openstack limits show --absolute -f json | jq '.[] | select(.Name=="maxTotalCores") | .Value')

echo '{"labels": {"tenant_id": "'$OS_PROJECT_ID'","project_name": "'$OS_PROJECT_NAME'"}, "results": {"items": '$MAXTOTALCORES'} }'

exit 0
