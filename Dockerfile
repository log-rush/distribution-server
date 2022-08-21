FROM golang:1.18-alpine AS Builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN go build -o ./build/server ./app/main.go

FROM alpine:3.14 AS Production
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=Builder /app/build ./
CMD [ "/app/server" ]