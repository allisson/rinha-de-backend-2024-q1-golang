#### development stage
FROM golang:1.22 AS build-env

# set envvar
ENV CGO_ENABLED=0
ENV GOOS=linux

# set workdir
WORKDIR /code

# get project dependencies
COPY go.mod /code/
RUN go mod download

# copy files
COPY . /code

# generate binary
RUN go build -tags=go_json -ldflags="-s -w" -o ./rinha ./cmd/rinha

#### final stage
FROM gcr.io/distroless/base:nonroot
COPY --from=build-env /code/rinha /
ENTRYPOINT ["/rinha"]
