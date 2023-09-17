FROM golang:1.18

COPY . /app/

WORKDIR /app

RUN go mod download
RUN go build --o app ./src

CMD ["./app"]