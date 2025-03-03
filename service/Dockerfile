FROM golang:1.20-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o template-service .

FROM surnet/alpine-wkhtmltopdf:3.21.2-0.12.6-small AS wkhtmltopdf

FROM alpine:3.17
WORKDIR /app
RUN apk add --no-cache \
    libstdc++ libx11 libxrender libxext \
    ca-certificates fontconfig freetype \
    ttf-liberation ttf-opensans \
    postgresql-client xvfb

COPY --from=wkhtmltopdf /bin/wkhtmltopdf /bin/wkhtmltopdf
COPY --from=builder /app/templates/ /app/templates/
COPY --from=builder /app/static/ /app/static/
COPY --from=builder /app/template-service /app/template-service
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

EXPOSE 8080
ENTRYPOINT ["/app/entrypoint.sh"]