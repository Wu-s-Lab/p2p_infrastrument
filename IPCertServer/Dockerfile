FROM golang:latest as builder 

WORKDIR /go/src/IPCertServer/
COPY . /go/src/IPCertServer
 
RUN go get github.com/mattn/go-sqlite3\
    && go get gopkg.in/yaml.v2\
    && go get github.com/pkg/errors\
    && go build -o IPCertServer -ldflags '-linkmode "external" -extldflags "-static"' .

FROM alpine:latest
WORKDIR /app/
COPY --from=builder  /go/src/IPCertServer/IPCertServer .
COPY --from=builder  /go/src/IPCertServer/config.yaml .

ENTRYPOINT ["./IPCertServer"]
