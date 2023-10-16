# # base image
# FROM golang:1.18-alpine as builder

# RUN mkdir /app

# COPY . /app

# WORKDIR /app

# RUN CGO_ENABLED=0 go build -o mailApp ./cmd/api

# RUN chmod +x /app/mailApp

# build tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY mailApp /app
COPY templates /templates


CMD ["/app/mailApp"]