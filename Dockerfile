FROM golang:1.21

WORKDIR /usr/app/

COPY ./vendor /usr/app/vendor
COPY ./cmd /usr/app/cmd
COPY ./internal /usr/app/internal
COPY ./go.mod /usr/app
COPY ./go.sum /usr/app

RUN go build -mod vendor -o ./bin/timetracker ./cmd/main.go

COPY ./app.env /usr/app
COPY ./auth.env /usr/app

STOPSIGNAL SIGINT

CMD [ "/usr/app/bin/timetracker" ]
