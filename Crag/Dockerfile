FROM golang:1.22.4

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /crag

EXPOSE 6969

CMD ["/crag"]