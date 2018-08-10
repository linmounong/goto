FROM golang:1.10-alpine as build

# build
# required by sqlite
RUN apk add --no-cache g++
WORKDIR /go/src/github.com/linmounong/goto
COPY . .
RUN go install github.com/linmounong/goto

# ship
FROM alpine:3.7
COPY --from=build /go/bin/goto /goto
ENTRYPOINT ["/goto"]
