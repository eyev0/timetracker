FROM golang:1.21

WORKDIR /usr/app/

COPY . /usr/app

RUN go build -mod vendor -o ./bin/timetracker ./cmd/main.go

STOPSIGNAL SIGINT

CMD [ "/usr/app/bin/timetracker" ]
