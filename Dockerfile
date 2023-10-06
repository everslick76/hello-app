FROM golang:1.20-alpine

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /hello-app

EXPOSE 8080

CMD [ "/hello-app" ]
