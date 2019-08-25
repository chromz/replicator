FROM golang:latest

LABEL maintainer="Rodrigo Custodio <jorocuva@gmail.com>"

WORKDIR /app

COPY . /app

RUN make

CMD ["./fastfood"]
