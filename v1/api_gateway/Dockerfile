FROM golang:1.19-alpine

RUN go version
WORKDIR /api-gateway
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ ./

RUN go build -o api-gateway main.go
# EXPOSE 8081
CMD ["./api-gateway"]