FROM golang:1.22.4

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /static

EXPOSE 6996

CMD ["/static"]