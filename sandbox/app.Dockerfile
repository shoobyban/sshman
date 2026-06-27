FROM node:22-alpine AS frontend

WORKDIR /src/frontend
COPY frontend/ ./
RUN npm install && npm run build

FROM golang:1.26-alpine AS builder

WORKDIR /src
COPY . ./
COPY --from=frontend /src/frontend/dist ./cmd/dist
RUN go build -o /out/sshman .

FROM alpine:3.22

RUN apk add --no-cache openssh-keygen ca-certificates
WORKDIR /app

COPY --from=builder /out/sshman /usr/local/bin/sshman
COPY sandbox/app-entrypoint.sh /app-entrypoint.sh
COPY sandbox/sshman.json /sandbox/sshman.json

RUN chmod +x /app-entrypoint.sh && mkdir -p /app/.ssh /sandbox-keys /sandbox

EXPOSE 8080

ENTRYPOINT ["/app-entrypoint.sh"]
