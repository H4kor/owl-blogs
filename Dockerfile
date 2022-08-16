FROM golang:1.19-alpine as build


RUN apk add --no-cache git

WORKDIR /tmp/owl

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/owl-web ./cmd/owl-web
RUN go build -o ./out/owl-cli ./cmd/owl-cli

FROM alpine:3.9 
RUN apk add ca-certificates

COPY --from=build /tmp/owl/out/ /bin/

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["/bin/owl-web"]