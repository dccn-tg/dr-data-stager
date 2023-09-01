# stage 0: compile go program
FROM golang:1.20
RUN mkdir -p /tmp/data-stager
WORKDIR /tmp/data-stager
ADD internal ./internal
ADD pkg ./pkg
ADD go.mod .
ADD go.sum .
ADD Makefile .
RUN GOOS=linux make

# stage 1: build image for the api-server
FROM centos:7 as api-server
RUN yum install -y nfs4-acl-tools sssd-client attr acl && yum clean all && rm -rf /var/cache/yum/*
WORKDIR /root
EXPOSE 8080
VOLUME ["/project", "project_freenas", "/project_cephfs", "/home"]
COPY --from=0 /tmp/data-stager/bin/data-stager-api .

## entrypoint in shell form so that we can use $PORT environment variable
ENTRYPOINT ["./data-stager-api"]

# stage 2: build image for the worker
FROM centos:7 as worker

# install irods client
RUN ( yum -y install epel-release wget curl )

RUN ( rpm --import https://packages.irods.org/irods-signing-key.asc && \
      wget -qO - https://packages.irods.org/renci-irods.yum.repo | \
      tee /etc/yum.repos.d/renci-irods.yum.repo )

RUN ( yum -y install irods-icommands )

# install required packages for NFS and SSSD clients
RUN ( yum install -y nfs4-acl-tools sssd-client attr acl )

# clean up temporary files created by yum install
RUN ( yum clean all && \
      rm -rf /var/cache/yum/* && \
      rm -rf /tmp/* )

# install jq
RUN ( curl -L -o /usr/local/bin/jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 && chmod +x /usr/local/bin/jq )

# create work directory 
RUN ( mkdir -p /opt/irods )
WORKDIR /opt/irods

# create empty directory for icat's public certificate.
RUN ( mkdir -p /opt/irods/ssl )

# define expected env. variables
ENV IRODS_HOST=${IRODS_HOST:-icat.data.donders.ru.nl}
ENV IRODS_PORT=${IRODS_PORT:-1247}
ENV IRODS_ZONE=${IRODS_PORT:-/nl.ru.donders}
ENV IRODS_USER_NAME=${IRODS_USER_NAME:-irods}

# overwrite default location of the irods_environemnts.json file
ENV IRODS_ENVIRONMENT_FILE=/opt/irods/irods_environments.json

# copy script for generating irods_environments.json file
COPY docker/worker/gen_irods_environment.sh .
RUN chmod +x ./gen_irods_environment.sh

# copy entrypoint script
COPY docker/worker/entrypoint.sh .
RUN chmod +x ./entrypoint.sh

VOLUME ["/project", "project_freenas", "/project_cephfs", "/home"]
COPY --from=0 /tmp/data-stager/bin/data-stager-worker .

# program entrypoint
ENTRYPOINT ["./entrypoint.sh"]
