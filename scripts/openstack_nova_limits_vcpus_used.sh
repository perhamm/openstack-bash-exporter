#!/bin/sh

TOTALCORESUSED=$(openstack limits show --absolute -f json | jq '.[] | select(.Name=="totalCoresUsed") | .Value')

echo '{"labels": {"tenant_id": "'$OS_PROJECT_ID'","project_name": "'$OS_PROJECT_NAME'"}, "results": {"items": '$TOTALCORESUSED'} }'

exit 0
