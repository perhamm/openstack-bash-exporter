#!/bin/sh

TOTALGIGABYTESUSED=$(openstack limits show --absolute -f json | jq '.[] | select(.Name=="totalGigabytesUsed") | .Value')

echo '{"labels": {"tenant_id": "'$OS_PROJECT_ID'","project_name": "'$OS_PROJECT_NAME'"}, "results": {"items": '$TOTALGIGABYTESUSED'} }'

exit 0

