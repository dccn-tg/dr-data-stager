# source data volumes
HOME_VOL=/tmp/home
PROJECT_VOL=/tmp/project
PROJECT_CEPHFS_VOL=/tmp/project_cephfs
PROJECT_FREENAS_VOL=/tmp/project_freenas

# OIDC authentication
AUTH_SERVER=https://login.dccn.nl
AUTH_CLIENT_ID=clientid
AUTH_CLIENT_SECRET=clientsecret

# stager task persistent store
TASK_DB_REDIS_DATA=/tmp/data

# configuration files
API_CONFIG=./config/api-server.yml
WORKER_CONFIG=./config/worker.yml
UI_CONFIG=./ui/config/default.json

# iCAT server certificate for secured iRODS communications
IRODS_ICAT_CERT=./docker/worker/icat-prod.pem

# crypto key pair for DR credential en/de-cryption
CRYPTO_RSA_PUBLIC=./config/public.pem
CRYPTO_RSA_PRIVATE=./config/keypair.pem