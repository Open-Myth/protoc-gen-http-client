FROM golang:1.24.4 AS install-stage

ENV GOBIN=/usr/local/bin

WORKDIR /app
COPY . .

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

RUN go mod download github.com/googleapis/googleapis@v0.0.0-20221209211743-f7f499371afa

RUN go install ./

ENV MOD=$GOPATH/pkg/mod

RUN mkdir -p /collect/validate
RUN cp -r $MOD/github.com/googleapis/googleapis@v0.0.0-20221209211743-f7f499371afa/google /collect/.

FROM alpine:latest AS generate-stage
RUN apk add curl unzip libc6-compat

ENV PROTOC_VERSION=3.14.0
ENV GRPC_WEB_VERSION=1.2.1
ENV BUFBUILD_VERSION=0.24.0

RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_VERSION/protoc-$PROTOC_VERSION-linux-x86_64.zip
RUN unzip protoc-$PROTOC_VERSION-linux-x86_64.zip -d protoc3
RUN mv protoc3/bin/* /usr/local/bin/
RUN mv protoc3/include/ /usr/local/include/

# copy binary files
COPY --from=install-stage /usr/local/bin/ /usr/local/bin/
COPY --from=install-stage /bin/protoc-gen-* /usr/local/bin/

# # copy default proto files
COPY --from=install-stage /usr/local/include/* /usr/local/include/

COPY --from=install-stage /collect/google/ /usr/local/include/google/
# COPY --from=install-stage /collect/validate /usr/local/include/

COPY --from=install-stage /app/proto /proto