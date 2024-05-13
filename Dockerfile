# build stage
FROM golang:1.21.5 as builder

WORKDIR /app
COPY . .
COPY .gitconfig /etc/gitconfig

RUN go mod download

RUN make build

# production stage
FROM alpine:latest
COPY --from=builder /app/api ./api
COPY --from=builder /app/build .
COPY --from=builder /app/config/prod.yaml ./config/

RUN apk update && \
    apk add --no-cache ca-certificates

RUN update-ca-certificates

EXPOSE 3000

CMD ["./app", "s"]
