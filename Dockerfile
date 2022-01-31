   
FROM golang:1.16-alpine
RUN apk add build-base



WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main .



ENTRYPOINT ["./main"]