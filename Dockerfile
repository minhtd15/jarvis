FROM redis:latest

# Sử dụng hình ảnh chứa RabbitMQ
FROM rabbitmq:latest

# Cài đặt Go 1.17
ENV GOLANG_VERSION 1.17
RUN wget -O go.tgz "https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz" \
    && tar -C /usr/local -xzf go.tgz \
    && rm go.tgz

# Thiết lập biến môi trường cho Go
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV GOBIN="/go/bin"

# Kiểm tra phiên bản Go đã cài đặt
RUN go version

# Vào thư mục làm việc
WORKDIR /go/src/app

# Copy mã nguồn ứng dụng Go vào thư mục làm việc
COPY . .

# Chạy các lệnh khác bạn muốn thực thi, ví dụ: go build, go test, ...
RUN go run app/service/main.go
# Mở cổng ứng dụng của RabbitMQ (mặc định là 5672)
EXPOSE 5672
