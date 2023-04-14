#!/bin/sh

TOTALRAMUSED=$(openstack --os-cloud=openstack limits show --absolute -f json | jq '.[] | select(.Name=="totalRAMUsed") | .Value')

TENANTID=$(cat clouds.yaml | shyaml get-value clouds.openstack.auth.project_id)

PROJECTNAME=$(cat clouds.yaml | shyaml get-value clouds.openstack.auth.project_name)

echo '{"labels": {"tenantid": "'$TENANTID'","projectname": "'$PROJECTNAME'"}, "results": {"items": '$TOTALRAMUSED'} }'

exit 0
