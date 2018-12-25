FROM golang:alpine

# setup the environment
ENV TZ=Europe/Berlin

# install dependencies
RUN apk --update --no-cache add git gcc musl-dev tzdata && \
    apk add zbar-dev --update-cache --repository https://nl.alpinelinux.org/alpine/edge/testing
RUN rm /usr/lib/libzbar.so
WORKDIR /work
ADD ./ ./

# build the go binary
RUN go build -ldflags \
        '-X "main.BuildTime='$(date -Iminutes)'" \
         -X "main.GitCommit='$(git rev-parse --short HEAD)'" \
         -X "main.GitBranch='$(git rev-parse --abbrev-ref HEAD)'" \
         -X "main.BuildNumber='$CI_BUILDNR'" \
         -s -w' \
         -v -o /work/uprun_linux .
