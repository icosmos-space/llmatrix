FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
WORKDIR /src/proxy
RUN CGO_ENABLED=0 go build -o /llmatrix-proxy ./cmd/llmatrix-proxy

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /llmatrix-proxy /usr/local/bin/llmatrix-proxy
EXPOSE 8080
ENTRYPOINT ["llmatrix-proxy"]
