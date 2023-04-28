FROM golang:latest AS builder

ADD backend/ /app
WORKDIR /app

RUN go mod download
RUN go build -o /main .

FROM golang:latest AS builderTwo

ADD cli/ /app
WORKDIR /app

RUN go mod download
RUN go build -o /cli .

FROM ubuntu:latest
RUN apt-get -y update
RUN apt-get -y install git

COPY --from=builder /main ./
COPY --from=builderTwo /cli ./
ENTRYPOINT ["./main"]
EXPOSE 8000