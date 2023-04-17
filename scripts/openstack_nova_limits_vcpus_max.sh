#!/bin/sh

if [ ! -f /tmp/config ]; then kubectl -n d8-cloud-provider-openstack get secrets  cloud-controller-manager -o json | jq -r '.data."cloud-config"' | base64 -d > /tmp/config; fi

export OS_AUTH_URL=$(cat /tmp/config | grep auth-url | awk -F'"' '{ print $2 }')
export OS_IDENTITY_API_VERSION="3"
export OS_INTERFACE="public"
export OS_PASSWORD=$(cat /tmp/config | grep password | awk -F'"' '{ print $2 }')
export OS_PROJECT_ID=$(cat /tmp/config | grep tenant-id | awk -F'"' '{ print $2 }')
export OS_REGION_NAME=$(cat /tmp/config | grep region | awk -F'"' '{ print $2 }')
export OS_USERNAME=$(cat /tmp/config | grep username | awk -F'"' '{ print $2 }')
export OS_USER_DOMAIN_NAME=$(cat /tmp/config | grep domain-name | awk -F'"' '{ print $2 }')

MAXTOTALCORES=$(openstack limits show --absolute -f json | jq '.[] | select(.Name=="maxTotalCores") | .Value')

echo '{"labels": {"tenant_id": "'$OS_PROJECT_ID'"}, "results": {"items": '$MAXTOTALCORES'} }'

exit 0
