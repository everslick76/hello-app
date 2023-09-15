FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./

RUN go mod download

RUN go mod tidy

COPY *.go ./

RUN go build -o /hello-app

CMD [ "/hello-app" ]
