FROM golang:1.13-buster as build

WORKDIR /go/src/nms-controller
ADD . .

RUN go mod download
RUN go build -o /go/main

# FROM gcr.io/distroless/base-debian10
FROM golang:1.13-buster
WORKDIR /go/
COPY --from=build /go/main .
COPY .env .

# ENV POSTGRES_URL "host=61.219.26.42 port=5432 user=postgres password=4ziw9jh70b3v0yydbk48 dbname=nms sslmode=disable"

# EXPOSE 8080

CMD ["./main"]
