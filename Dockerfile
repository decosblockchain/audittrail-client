FROM golang:alpine as build

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh gcc musl-dev
ENV GOROOT=/usr/local/go
RUN go get -v github.com/ethereum/go-ethereum
RUN go get -v github.com/gorilla/mux
RUN go get -v github.com/kardianos/service
COPY . /usr/local/go/src/github.com/decosblockchain/audittrail-client
WORKDIR /usr/local/go/src/github.com/decosblockchain/audittrail-client
RUN go get -v ./...
RUN go build


FROM alpine
RUN apk add --no-cache ca-certificates
RUN mkdir -p /app/bin/data
COPY --from=build /usr/local/go/src/github.com/decosblockchain/audittrail-client/audittrail-client /app/bin/audittrail-client
RUN cd /app/bin
WORKDIR /app/bin
EXPOSE 8001

CMD ["./audittrail-client", "console"]