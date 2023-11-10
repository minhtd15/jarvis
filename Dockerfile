# Use the official Golang image as the base image
FROM golang:1.18.2-alpine
WORKDIR /app/service
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o /out/main ./
ENTRYPOINT ["/app/service/main"]
