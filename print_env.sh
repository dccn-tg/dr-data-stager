#!/bin/bash

echo "# version"
echo "DOCKER_IMAGE_TAG=$DOCKER_IMAGE_TAG"
echo 
echo "# docker registry endpoint"
echo "DOCKER_REGISTRY=$DOCKER_REGISTRY"
echo 
echo "# volume for home directory"
echo "HOME_VOL=$HOME_VOL"
echo
echo "# volume for project directory"
echo "PROJECT_VOL=$PROJECT_VOL"
echo
echo "# volume for project_freenas directory"
echo "PROJECT_FREENAS_VOL=$PROJECT_FREENAS_VOL"
echo
echo "# volume for project_cephfs directory"
echo "PROJECT_CEPHFS_VOL=$PROJECT_CEPHFS_VOL"
echo
echo "# OIDC auth server"
echo "AUTH_SERVER=$AUTH_SERVER"
echo
echo "# OIDC client ID"
echo "AUTH_CLIENT_ID=$AUTH_CLIENT_ID"
echo
echo "# OIDC client secret"
echo "AUTH_CLIENT_SECRET=$AUTH_CLIENT_SECRET"
echo
echo "# Stager task database persistent store"
echo "TASK_DB_REDIS_DATA=$TASK_DB_REDIS_DATA"
echo
echo "# API server configuration file"
echo "API_CONFIG=$API_CONFIG"
echo
echo "# worker configuration file"
echo "WORKER_CONFIG=$WORKER_CONFIG"
echo
echo "# iCAT certificate for SSL communication and data traffic encryption"
echo "IRODS_ICAT_CERT=$IRODS_ICAT_CERT"
echo
echo "# RSA public key for DR credential en-/decryption"
echo "CRYPTO_RSA_PUBLIC=$CRYPTO_RSA_PUBLIC"
echo
echo "# RSA private key for DR credential en-/decryption"
echo "CRYPTO_RSA_PRIVATE=$CRYPTO_RSA_PRIVATE"
