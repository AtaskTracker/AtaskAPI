FROM golang:alpine@sha256:4dd403b2e7a689adc5b7110ba9cd5da43d216cfcfccfbe2b35680effcf336c7e

RUN mkdir /app
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /main
WORKDIR /app
ENTRYPOINT ["/main"]