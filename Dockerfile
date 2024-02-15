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

ADD config/api-server.yml .

## entrypoint in shell form so that we can use $PORT environment variable
ENTRYPOINT ["./data-stager-api", "-c", "api-server.yml"]

# stage 2: build image for the worker
FROM centos:7 as worker

# install irods client
RUN ( yum -y install epel-release wget curl )

# install required packages for NFS and SSSD clients
RUN ( yum install -y nfs4-acl-tools sssd-client attr acl )

# clean up temporary files created by yum install
RUN ( yum clean all && \
      rm -rf /var/cache/yum/* && \
      rm -rf /tmp/* )

# create work directory 
RUN ( mkdir -p /opt/stager-worker )
WORKDIR /opt/stager-worker

# create empty directory for icat's public certificate.
RUN ( mkdir -p /opt/irods/ssl )

VOLUME ["/project", "project_freenas", "/project_cephfs", "/home"]
COPY --from=0 /tmp/data-stager/bin/data-stager-worker .

ADD config/worker.yml .

## entrypoint in shell form so that we can use $PORT environment variable
ENTRYPOINT ["./data-stager-worker", "-c", "worker.yml"]
