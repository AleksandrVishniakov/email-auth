# email-auth
## How to start
* Clone this repository
  ```bash
  $ git clone https://github.com/AleksandrVishniakov/email-auth
  $ cd email-auth
  ```
* Set environment variables or create .env file in root directory:
  ```env
  DB_PASSWORD=""
  EMAIL_PASSWORD=""
  ```
* Run ```make``` or execute this commands:
  ```bash
  $ docker build -t email-auth:local .
  $ docker compose up
  ```
* Container started in docker on port 8003 (can be changed in docker-compose.yml)
## Start example
To start example page you need to install npm and execute this commands in root directory:
```bash
$ cd web/app
$ npm install
$ npm start
```
Page will be available on [http://localhost:3000](http://localhost:3000)
## API docs
### Auth user
Can be used for creating a new account or signing in into an existing one
#### Request
```http
GET /auth?email=your_email
```
Checks if user with provided email exists in database and sends an email with 6-digit verification code
#### Response
```json
{
  "isUserAuthorized": true
}
```
```isUserAutorized=true``` if user has a verified email
#### Status codes:
* ```400 Bad Request``` if an email parameter is empty
* ```500 Internal Server Error``` if error in server work occurred
### Verify email
Used for verifying an email with provided 6-digit code
#### Request
```http
GET /verify/your_email?code=XXXXXX
```
Compares verifying codes and register user if he wasn't
#### Status codes:
* ```400 Bad Request``` if an email or code parameter is empty, or if code parsing to int failed
* ```401 Unauthorized``` if provided code does not match with a required one
* ```404 Not Found``` if a provided email does not exist in database
* ```500 Internal Server Error``` if error in server work occurred
### Get user
Returns information about user with provided email
#### Request
```http
GET /user/your_email
```
#### Response
```json
{
  "id": 0,
  "email": "john.doe@example.com",
  "isEmailVerified": true,
  "createdAt": "05-06-2024 14:30"
}
```
#### Status codes
* ```400 Bad Request``` if an email parameter is empty
* ```404 Not Found``` if a provided email does not exist in database
* ```500 Internal Server Error``` if error in server work occurred
### API error
API returns an error in this format:
```json
{
  "code": 400,
  "message": "bad request",
  "timestamp": "05-06-2024 14:30"
}
```
