FROM golang:1.16.0 AS build-stage
ENV GOOS linux
ENV GOARCH amd64
COPY . ./src/github.com/ralpioxxcs/go-onedrive-cli
WORKDIR /go/src/github.com/ralpioxxcs/go-onedrive-cli
RUN go build

FROM busybox:1.33
COPY --from=build-stage /go/src/github.com/ralpioxxcs/go-onedrive-cli /usr/local/bin/go-onedrive-cli

ENTRYPOINT ["go-onedrive-cli"]