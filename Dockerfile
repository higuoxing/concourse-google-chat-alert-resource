FROM golang:1.21-alpine AS build
WORKDIR /go/src/github.com/higuoxing/concourse-google-chat-alert-resource
RUN apk --no-cache add --update git

COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -o /check github.com/higuoxing/concourse-google-chat-alert-resource/check
RUN go build -o /in github.com/higuoxing/concourse-google-chat-alert-resource/in
RUN go build -o /out github.com/higuoxing/concourse-google-chat-alert-resource/out

FROM alpine:3.18
RUN apk add --no-cache ca-certificates

COPY --from=build /check /opt/resource/check
COPY --from=build /in /opt/resource/in
COPY --from=build /out /opt/resource/out
