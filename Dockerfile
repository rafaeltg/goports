FROM golang:1.21.4-alpine3.18 as builder

RUN apk add --update --no-cache \
    curl \
    git \
    make

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN make build-server
RUN make build-ingestor

FROM scratch as server
COPY --from=builder /src/build/http /http
ENTRYPOINT [ "/http" ]

FROM scratch as ingestor
COPY --from=builder /src/build/ingestor /ingestor
ENTRYPOINT [ "/ingestor" ]