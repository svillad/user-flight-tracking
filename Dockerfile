FROM golang:1.20-bullseye as builder

RUN go install golang.org/dl/go1.20@latest \
    && go1.20 download

# Instalar mockgen
RUN go install github.com/golang/mock/mockgen@latest

# Instalar golangci-lint
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

WORKDIR /build

COPY go.mod .
COPY go.sum .
COPY . .

# Compilar el microservicio
RUN go build -o ./app ./server

# Exponer el puerto del microservicio
EXPOSE 8080

# Comando para ejecutar el microservicio
CMD ["./app"]