FROM node:21-alpine

# application metadata
LABEL donders.ru.nl.app_name "stager-ui"
LABEL donders.ru.nl.app_maintainer "h.lee@donders.ru.nl"
LABEL donders.ru.nl.app_code_repository "https://github.com/Donders-Institute/data-stager"

# install the stager-ui nodejs applicaiton
RUN ( mkdir -p /opt/stager-ui/bin )
WORKDIR /opt/stager-ui
COPY package.json package.json
COPY package-lock.json package-lock.json
RUN ( npm install )

COPY *.js ./
COPY bin ./bin
COPY lib ./lib
COPY routes ./routes
COPY views ./views
COPY public ./public
RUN ( mkdir config )
COPY config/default.json ./config

# COPY start_stager-ui.sh start_stager-ui.sh
# RUN ( chmod +x start_stager-ui.sh )

# run a stager-ui
# temporarily disable check on self-signed certificate until
# webdav-fs module supports it, see issue: 
# https://github.com/perry-mitchell/webdav-fs/issues/54
#
# TODO: replace it with trusting the self-signed certificate
ENV NODE_TLS_REJECT_UNAUTHORIZED=0
EXPOSE 3080
CMD [ "npm", "start" ]
