FROM alpine

# create app directory
RUN mkdir -p /app/static/static
#WORKDIR /app

# copy backend binary
COPY back/docker-stalker /app

ARG frontend_path="/app/static"

# copy over build files
COPY front/build/index.html ${frontend_path}
COPY front/build/static ${frontend_path}/static

# set environment variable to be frontend path
ENV APP_BUILD_FOLDER=$frontend_path

# expose server port
EXPOSE 8080

# copy over react app build files too, so that backend can serve
ENTRYPOINT ["/app/docker-stalker"]