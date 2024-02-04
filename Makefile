all: d-build d-compose

local: app/cmd/app/main.go
	go run app/cmd/app/main.go

d-build: .
	docker build -t email-auth:local .

d-compose:
	docker compose up