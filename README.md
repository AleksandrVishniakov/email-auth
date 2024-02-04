# email-auth
## How to start
* Clone this repository
  ```bash
  $ git clone https://github.com/AleksandrVishniakov/email-auth
  $ cd email-auth
  ```
* Set envirionment variables or create .env file in root directory:
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
To start example page you need to install npm and exetute this commands in root directory:
```bash
$ cd web/app
$ npm install
$ npm start
```
Page wiil be avalilable on [http://localhost:3000](http://localhost:3000)
