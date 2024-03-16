FROM golang:1.21.6

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go *.env ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /crag

EXPOSE 6969

CMD ["/crag"]