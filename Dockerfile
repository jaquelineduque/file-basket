FROM golang:1.15

WORKDIR /app/server
COPY /src/. /app/server

RUN go build -o /app/bin
CMD ["/app/bin"]