FROM golang:1.21-alpine as build

WORKDIR /go/src/email-auth

COPY app ./app/

COPY go.mod go.sum ./
#RUN go mod tidy
RUN go mod download

RUN go build -o ../../bin/app ./app/cmd/app/main.go

FROM alpine
WORKDIR /go

COPY configs ./configs/
COPY ./.env .
COPY ./web ./web/

COPY --from=build /go/bin/app /bin/app

CMD ["app"]