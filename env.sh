# source data volumes
HOME_VOL=/tmp/home
PROJECT_VOL=/tmp/project
PROJECT_CEPHFS_VOL=/tmp/project_cephfs
PROJECT_FREENAS_VOL=/tmp/project_freenas

# stager task persistent store
TASK_DB_REDIS_DATA=/tmp/data

# configuration files
API_CONFIG=$(pwd)/config/api-server.yml
WORKER_CONFIG=$(pwd)/config/worker.yml

# iCAT server certificate for secured iRODS communications
IRODS_ICAT_CERT=$(pwd)/docker/worker/icat-prod.pem