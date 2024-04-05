FROM golang:latest

LABEL maintainer="fredsu@stud.ntnu.no"
LABEL deployedBy="Fredrik Sundt-Hansen"

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o reception cmd/reception/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rest cmd/rest/main.go

WORKDIR /app

COPY backend.sh ./app/backend.sh
RUN chmod +x ./app/backend.sh

EXPOSE 8080

ENTRYPOINT ["/bin/sh", "./app/backend.sh"]