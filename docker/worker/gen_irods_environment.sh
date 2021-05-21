#!/bin/bash

## This script generates irods_environment.json file in /opt/irods with
## settings specified with environment variables.

## template
template=$(cat << EOF
{
    "irods_host": "icat.data.donders.ru.nl",
    "irods_port": 1247,
    "irods_user_name": "irods",
    "irods_default_resource": "",
    "irods_home": "/rdm/di",
    "irods_zone_name": "rdm",
    "irods_client_server_negotiation": "request_server_negotiation",
    "irods_client_server_policy": "CS_NEG_REFUSE",
    "irods_encryption_key_size": 32,
    "irods_encryption_salt_size": 8,
    "irods_encryption_num_hash_rounds": 16,
    "irods_encryption_algorithm": "AES-256-CBC",
    "irods_default_hash_scheme": "SHA256",
    "irods_match_hash_policy": "compatible",
    "irods_ssl_ca_certificate_file": "/opt/irods/ssl/icat.pem",
    "irods_ssl_verify_server": "cert",
    "irods_authentication_scheme": "PAM",
    "irods_authentication_file": "/opt/irods/.irodsA"
}
EOF
)

echo "${template}" | jq -M \
	--arg irods_host ${IRODS_HOST} \
	--arg irods_port ${IRODS_PORT} \
	--arg irods_zone_name ${IRODS_ZONE_NAME} \
	--arg irods_user_name ${IRODS_USER_NAME} \
	--arg irods_home ${IRODS_ZONE_NAME}/home/${IRODS_USER_NAME} \
	'.irods_host=$irods_host | .irods_port=($irods_port|tonumber) | .irods_zone_name=$irods_zone_name | .irods_user_name=$irods_user_name | .irods_home=$irods_home' \
	> ${IRODS_ENVIRONMENT_FILE}
