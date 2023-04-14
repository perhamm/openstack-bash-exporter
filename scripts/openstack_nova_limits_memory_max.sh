#!/bin/sh

MAXTOTALRAMSIZE=$(openstack limits show --absolute -f json | jq '.[] | select(.Name=="maxTotalRAMSize") | .Value')

echo '{"labels": {"tenant_id": "'$OS_PROJECT_ID'","project_name": "'$OS_PROJECT_NAME'"}, "results": {"items": '$MAXTOTALRAMSIZE'} }'

exit 0
