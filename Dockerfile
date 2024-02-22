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

# stage 1: build image for the required packages
FROM centos:7 as base
RUN yum install -y nfs4-acl-tools sssd-client attr acl && yum clean all && rm -rf /var/cache/yum/*
# clean up temporary files created by yum install
RUN ( yum clean all && \
      rm -rf /var/cache/yum/* && \
      rm -rf /tmp/* )
# create work directory 
RUN ( mkdir -p /opt/stager )
# create configuration directory 
RUN ( mkdir -p /etc/stager )
# expected data sources
VOLUME ["/project", "/project_freenas", "/project_cephfs", "/home"]

# stage 2: build image for api-server
FROM base as api-server
WORKDIR /opt/stager
EXPOSE 8080
COPY --from=0 /tmp/data-stager/bin/data-stager-api .
ENTRYPOINT ["./data-stager-api"]

# stage 3: build image for the worker
FROM base as worker
WORKDIR /opt/stager
RUN ( mkdir -p /opt/irods/ssl )
COPY --from=0 /tmp/data-stager/bin/data-stager-worker .
ENTRYPOINT ["./data-stager-worker"]
