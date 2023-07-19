##
## Build Container
##
FROM golang:1.20-alpine as build


RUN apk add --no-cache git

WORKDIR /tmp/owl

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/owl ./cmd/owl


##
## Run Container
##
FROM alpine
RUN apk add ca-certificates

COPY --from=build /tmp/owl/out/ /bin/

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
ENTRYPOINT ["/bin/owl"]