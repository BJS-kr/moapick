FROM golang:1.22.0 AS builder
# 내부 인증서 문제 때문에 일단 배포하기 위해 full blown image 사용
# 우분투에서 돌리고자 할 때:docker build -t iron2ron/moaick --platform="linux/amd64" .
# WORKDIR /build

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download 

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /moapick


# FROM scratch
# COPY --from=builder /moapick /moapick
# COPY --from=builder /build/test.env /test.env

ENTRYPOINT ["/moapick"]

