
# Hướng dẫn triển khai và truy cập trang web CRM Procast

Sau đây là các bước triển khai và tiếp cận trang web CRM Procast:

## Truy cập vào thư mục dự án
Trên máy chủ **`jupyter-fall2324w20g6`**

Trỏ đến thư mục code/batman:
```bash 
cd home/joyvan/code/batman
``` 
hoặc
```bash
cd ~/code/batman
```

## Triển khai phần mềm
Thực thi phần mềm back-end
```bash
go run app/service/main.go
```

Triển khai cổng của dịch vụ để có thể tiếp cận trang web qua internet, chạy dòng lệnh trên một terminal khác hoặc ấn <kbd> Ctrl </kbd> + <kbd> C </kbd> để dừng chương trình và chạy dòng lệnh trên rồi thực thi lại chương trình
```bash
/etc/jupyter/bin/expose 8081
```

### Truy cập trang web
Truy cập trang web thông qua đường dẫn:
<http://fall2324w20g6.int3306.freeddns.org/web/>

Đăng ký tài khoản mới hoặc đăng nhập tài khoản sau để truy cập với quyền `admin`:

> ***Username/Email:*** admin@procast.com
>
> ***Password:*** admin@1234





