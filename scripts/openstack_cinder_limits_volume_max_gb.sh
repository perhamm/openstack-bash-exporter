#!/bin/sh

MAXTOTALVOLUMEGIGABYTES=$(openstack limits show --absolute -f json | jq '.[] | select(.Name=="maxTotalVolumeGigabytes") | .Value')

echo '{"labels": {"tenant_id": "'$OS_PROJECT_ID'","project_name": "'$OS_PROJECT_NAME'"}, "results": {"items": '$MAXTOTALVOLUMEGIGABYTES'} }'

exit 0
