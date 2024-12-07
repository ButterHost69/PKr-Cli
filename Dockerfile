# Just for trying doesnt work [orginally for automating testing]
FROM golang:1.23

WORKDIR /app

EXPOSE 9000
EXPOSE 9001

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# BUILD PKr-CLI
RUN go build -o cli .

# BUILD PKr-Base
RUN cd PKr-base && go build -o ../base

# CMD [ "./base" ]