## Otp App BackEnd Api Service

### generate swagger docs

`make swagger`

## clone repo

`git clone https://github.com/ppeymann/top-app.git`

### build docker image

`sudo docker build --tag otpapp:latest`

## 1. ثبت‌نام (Sign Up)

ایجاد کاربر جدید با شماره موبایل.  
**Endpoint:**  
POST /api/v1/user/signup

**Request Body:**

```json
{
  "mobile": "09123456789"
}
```

Response:
{
"result": {
"accessToken": "JWT_TOKEN",
"refreshToken": "JWT_REFRESH"
},
"errors": [],
"status": 200
}

POST /api/v1/user/signin

Request Body:
{
"mobile": "09123456789"
}

POST /api/v1/user/otp

Request Body:
{
"mobile": "09123456789",
"otp": "123456"
}

GET /api/v1/user/
Headers:
Authorization: Bearer <ACCESS_TOKEN>

GET /api/v1/user/{offset}/{page}
Headers:
Authorization: Bearer <ACCESS_TOKEN>

Response:
{
"result": [
{
"id": 1,
"mobile": "09123456789",
"createdAt": "2025-08-12T12:34:56Z"
}
],
"errors": [],
"status": 200
}
