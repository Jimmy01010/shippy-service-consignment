FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache git

RUN mkdir /app
WORKDIR /app

ENV GO111MODULE=on

COPY . .

WORKDIR consignment-cli

RUN  go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shippy-cli-consignment


FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
ADD consignment-cli/consignment.json /app/consignment.json
COPY --from=builder /app/consignment-cli/shippy-cli-consignment .

CMD ["./shippy-cli-consignment"]