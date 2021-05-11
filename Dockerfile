FROM golang:alpine

RUN mkdir /app
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /main
WORKDIR /app
ENTRYPOINT ["/main"]