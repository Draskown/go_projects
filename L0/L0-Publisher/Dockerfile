FROM golang:1.21-alpine AS build

WORKDIR /build

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main main.go

FROM alpine:3.18

COPY --from=build main /app/main

ENTRYPOINT ["/app/main"]