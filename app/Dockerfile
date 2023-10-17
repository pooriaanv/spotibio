FROM golang:alpine

WORKDIR /app

COPY spotibio.go .
COPY .env .

RUN go mod init pooriaanv.com/spotibio

RUN go mod tidy

RUN go get

RUN go build -o spotibio spotibio.go

CMD ["./spotibio"]