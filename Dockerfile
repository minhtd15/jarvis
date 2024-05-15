## Use the RabbitMQ base image
#FROM golang:1.17.12 AS builder
#
#WORKDIR /edu
#
## Copy source code into the container
#COPY . .
#
## Compile the application
##RUN go mod vendor
#RUN go build -o ./bin/app ./app/service/main.go
#
#
## Use the Alpine base image for the final application
#FROM alpine:latest
#
## Create the directory for the binary
#RUN mkdir -p /usr/local/bin
#
## Copy the binary from the builder stage to the final image
#COPY --from=builder /edu/bin/app/service /usr/local/bin
#
## Set executable permissions on the binary
#RUN chmod +x /usr/local/bin/service
#
## Run the application when the container starts
#CMD ["/usr/local/bin/service"]
#

# syntax=docker/dockerfile:1

# ==============================# ==============================# ==============================# ==============================# ==============================


FROM golang:1.17.12

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /batman1 ./app/service/main.go
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /docker-gs-ping ./app/service/main.go
# Optional:
# To bind to a TCP port, runtime parametememrs must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 80

# Run
CMD ["/batman1"]




