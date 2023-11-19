# Start from the latest golang base image
FROM golang:alpine AS builder

# Add Maintainer Info
LABEL maintainer="Tomeu Uris tomeu.uris.dev@gmail.com"

RUN apk --no-cache add ca-certificates gcc musl-dev

# Install curl and unzip
RUN apk add --no-cache curl unzip git

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Generate Swagger documentation
RUN swag init --parseDependency --parseInternal

# Build the Go app
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .



# Start a new stage for development
FROM alpine:latest AS development

RUN apk --no-cache add ca-certificates

# Create a new user and switch to that user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Change to non-root privilege
USER appuser

WORKDIR /home/appuser/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder --chown=appuser:appgroup --chmod=555 /app/main .

# Copy the Swagger documentation from the previous stage
COPY --from=builder --chown=appuser:appgroup --chmod=444 /app/docs ./docs

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"] 



# Start a new stage for production
FROM alpine:latest AS production

ENV ENV=prod

RUN apk --no-cache add ca-certificates

# Create a new user and switch to that user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Change to non-root privilege
USER appuser

WORKDIR /home/appuser/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder --chown=appuser:appgroup --chmod=555 /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"] 