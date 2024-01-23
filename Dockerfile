FROM golang:1.21.6-alpine3.19 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download 

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /moapick

FROM scratch
COPY --from=builder /moapick /moapick
COPY --from=builder /build/test.env /test.env

ENTRYPOINT ["/moapick"]

