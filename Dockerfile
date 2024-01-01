FROM golang:1.21-alpine
WORKDIR /app/src/cmd/cutelet

RUN go install github.com/cosmtrek/air@latest

CMD ["air", "-c", ".air.toml"]