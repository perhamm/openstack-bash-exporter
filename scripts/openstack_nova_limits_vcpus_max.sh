#!/bin/sh

MAXTOTALCORES=$(openstack --os-cloud=openstack limits show --absolute -f json | jq '.[] | select(.Name=="maxTotalCores") | .Value')

TENANTID=$(cat clouds.yaml | shyaml get-value clouds.openstack.auth.project_id)

PROJECTNAME=$(cat clouds.yaml | shyaml get-value clouds.openstack.auth.project_name)

echo '{"labels": {"tenantid": "'$TENANTID'","projectname": "'$PROJECTNAME'"}, "results": {"items": '$MAXTOTALCORES'} }'

exit 0
