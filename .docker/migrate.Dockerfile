ARG BUILDER_IMAGE=golang:1.21.1-alpine
ARG ALPINE_IMAGE=alpine:3.15
################################################################################
# Build stage
################################################################################
FROM $BUILDER_IMAGE AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .


RUN go mod download

COPY . .

#RUN go mod download
RUN go build -o /bin/migrate ./cmd/migrate/main.go

FROM $ALPINE_IMAGE

WORKDIR /root/
COPY --from=build /bin/migrate .
COPY --from=build /app/dev/local.env dev/local.env
COPY --from=build /app/migrations migrations

CMD ["./migrate"]
