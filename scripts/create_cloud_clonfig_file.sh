#!/bin/sh

if [ ! -f /tmp/config ]; then kubectl -n d8-cloud-provider-openstack get secrets  cloud-controller-manager -o json | jq -r '.data."cloud-config"' | base64 -d > /tmp/config; fi
