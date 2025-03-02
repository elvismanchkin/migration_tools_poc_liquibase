FROM golang:1.20-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum files first
COPY go.mod go.sum ./

# Fix missing dependencies and download them
RUN go mod tidy
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o template-service .

# Use surnet/alpine-wkhtmltopdf as an intermediate stage for wkhtmltopdf binaries
FROM surnet/alpine-wkhtmltopdf:3.21.2-0.12.6-small AS wkhtmltopdf

# Final stage
FROM alpine:3.17

WORKDIR /app

# Install wkhtmltopdf dependencies
RUN apk add --no-cache \
    libstdc++ \
    libx11 \
    libxrender \
    libxext \
    ca-certificates \
    fontconfig \
    freetype \
    ttf-liberation \
    ttf-opensans \
    postgresql-client \
    xvfb

# Copy wkhtmltopdf from the intermediate stage
COPY --from=wkhtmltopdf /bin/wkhtmltopdf /bin/wkhtmltopdf
#COPY --from=wkhtmltopdf /bin/wkhtmltopdf /bin/
#COPY --from=wkhtmltopdf /bin/libwkhtmltox.so /bin/

# Copy templates and static files
COPY --from=builder /app/templates/ /app/templates/
COPY --from=builder /app/static/ /app/static/

# Copy the binary from the builder stage
COPY --from=builder /app/template-service /app/

# Copy entrypoint script
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Verify wkhtmltopdf installation
RUN wkhtmltopdf --version

# Expose the port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/app/entrypoint.sh"]