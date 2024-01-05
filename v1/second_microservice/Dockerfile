FROM golang:1.19-alpine

RUN go version
WORKDIR /second-microservice
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ ./

RUN go build -o second-microservice main.go
# EXPOSE 8081
CMD ["./second-microservice"]