# Use golang:alpine as the builder stage
FROM golang:alpine as builder

LABEL maintainer="fredsu@stud.ntnu.no"
LABEL deployedBy="Fredrik Sundt-Hansen"

RUN apk update && apk add --no-cache git

# Setup the working directory
RUN mkdir /app
WORKDIR /app

# Copy the go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application's source code
COPY . .

# Build the applications
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o reception cmd/reception/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rest cmd/rest/main.go

# Set backend.sh script executable
COPY backend.sh ./backend.sh
RUN chmod +x backend.sh

# Start a new stage from scratch for the final image
FROM alpine:latest

# Copy the built executables and the backend.sh script from the builder stage
COPY --from=builder /app/reception /app/reception
COPY --from=builder /app/rest /app/rest
COPY --from=builder /app/backend.sh /app/backend.sh
COPY --from=builder /app/.env /app/.env

# Set the working directory
WORKDIR /app

# Ensure the script is executable
RUN chmod +x backend.sh

# Expose the application's port
EXPOSE 8080

# Set the entrypoint to the backend.sh script
ENTRYPOINT ["/bin/sh", "/app/backend.sh"]