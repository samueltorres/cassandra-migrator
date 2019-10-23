FROM golang:1.13 as builder
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod /app
RUN go get -v all
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux GOPROXY=https://proxy.golang.org go build -o app *.go

FROM python:2-alpine
ENV CQLSH_VERSION 5.0.3
RUN pip install cqlsh==5.0.3
WORKDIR /app
COPY --from=builder /app/app .
ENTRYPOINT ["./app"]