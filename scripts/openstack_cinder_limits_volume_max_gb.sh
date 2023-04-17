#!/bin/sh

export OS_AUTH_URL=$(cat /tmp/config | grep auth-url | awk -F'"' '{ print $2 }')
export OS_IDENTITY_API_VERSION="3"
export OS_INTERFACE="public"
export OS_PASSWORD=$(cat /tmp/config | grep password | awk -F'"' '{ print $2 }')
export OS_PROJECT_ID=$(cat /tmp/config | grep tenant-id | awk -F'"' '{ print $2 }')
export OS_REGION_NAME=$(cat /tmp/config | grep region | awk -F'"' '{ print $2 }')
export OS_USERNAME=$(cat /tmp/config | grep username | awk -F'"' '{ print $2 }')
export OS_USER_DOMAIN_NAME=$(cat /tmp/config | grep domain-name | awk -F'"' '{ print $2 }')

MAXTOTALVOLUMEGIGABYTES=$(openstack limits show --absolute -f json | jq '.[] | select(.Name=="maxTotalVolumeGigabytes") | .Value')

echo '{"labels": {"tenant_id": "'$OS_PROJECT_ID'"}, "results": {"items": '$MAXTOTALVOLUMEGIGABYTES'} }'

exit 0
