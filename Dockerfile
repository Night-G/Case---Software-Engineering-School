FROM golang:1.22

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY ./DBmigrations /app/migrations

RUN go build -o main .

ENV ETHEREAL_EMAIL=""
ENV ETHEREAL_PASSWORD=""

EXPOSE 8080
CMD ["./main"]