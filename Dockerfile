#Temel Image
FROM golang:1.24.5

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o goshort

CMD ["./goshort"]