FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o tp-security ./cmd/main.go

CMD ["./tp-security"]