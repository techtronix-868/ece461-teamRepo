FROM golang:latest AS builder

ADD . /app
WORKDIR /app

RUN go mod download
RUN go build -o /main .

FROM ubuntu:latest
COPY --from=builder /main ./
ENTRYPOINT ["./main"]
EXPOSE 8000