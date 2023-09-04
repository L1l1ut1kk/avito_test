FROM golang:alpine as builder


WORKDIR /avito
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./main

# Run the application
CMD ["./main"]


