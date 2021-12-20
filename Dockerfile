FROM golang:1.16-buster as builder

WORKDIR /build

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o app . && chmod +x app

FROM debian:buster-slim as runner

RUN apt-get update && apt-get install -y ca-certificates openssl

COPY --from=builder /build/app /

ENTRYPOINT [ "/app" ]

CMD [ "5" ]