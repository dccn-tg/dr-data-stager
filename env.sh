IRODS_HOST=icat.data.donders.ru.nl
IRODS_PORT=1247
IRODS_ZONE_NAME=/nl.ru.donders
IRODS_USER_NAME=irods
IRODS_ICAT_CERT=$(pwd)/docker/worker/icat-prod.pem
API_CONFIG=$(pwd)/config/api-server.yml
WORKER_CONFIG=$(pwd)/config/worker.yml

# volumes
HOME_VOL=/tmp/home
PROJECT_VOL=/tmp/project
PROJECT_CEPHFS_VOL=/tmp/project_cephfs
PROJECT_FREENAS_VOL=/tmp/project_freenas
