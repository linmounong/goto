FROM golang:1.9.2-alpine3.7 as build

# build
RUN apk add --no-cache g++
WORKDIR /go/src/github.com/linmounong/goto
COPY . .
RUN go install github.com/linmounong/goto

# ship
FROM mm.taou.com:5000/algo/alpine:3.7
COPY --from=build /go/bin/goto /goto
ENTRYPOINT ["/goto"]
