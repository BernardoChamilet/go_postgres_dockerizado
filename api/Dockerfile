FROM golang:1.22

WORKDIR /usr/src/diretorioApi

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/minhaApi

EXPOSE 5000

CMD ["minhaApi"]