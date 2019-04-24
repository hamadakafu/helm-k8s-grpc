FROM golang

WORKDIR /go/src/app

ENV GO111MODULE=on

COPY go.mod go.sum ./

RUN go mod download

COPY . .    

RUN go build -o /go/src/app/app

CMD [ "./app" ] 