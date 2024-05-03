FROM golang:1.22.2 AS build-stage
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /migrate ./cmd/migrate/main.go

FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM alpine AS build-release-stage
WORKDIR /

COPY --from=build-stage /api /api
COPY --from=build-stage /migrate /migrate

RUN mkdir -p /cmd/migrate/migrations
COPY ./cmd/migrate/migrations /cmd/migrate/migrations

COPY ./cmd/scripts/run.sh /run.sh
RUN chmod +x run.sh

EXPOSE 8080

ENTRYPOINT ["/bin/sh", "/run.sh"]