FROM golang:1.21-alpine3.19 as builder-go

WORKDIR /src

COPY ../go.mod go.sum ./
RUN go mod download

ARG VERSION
COPY .. .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o srv ./cmd


FROM alpine:3.19 as runner

WORKDIR /cmd
COPY --from=builder-go /src/srv .

CMD ["/cmd/srv"]
