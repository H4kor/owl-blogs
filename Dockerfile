##
## Build Container
##
FROM golang:1.21-alpine as build


RUN apk add --no-cache --update git gcc g++

WORKDIR /tmp/owl

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o ./out/owl ./cmd/owl


##
## Run Container
##
FROM alpine
RUN apk add ca-certificates

COPY --from=build /tmp/owl/out/ /bin/

# This container exposes port 8080 to the outside world
EXPOSE 3000

WORKDIR /owl

# Run the binary program produced by `go install`
ENTRYPOINT ["/bin/owl"]