FROM golang:1.16-buster

WORKDIR /app

COPY . .

RUN go build -o source_app . && chmod +x source_app

EXPOSE 8080

ENTRYPOINT [ "/app/source_app" ]
