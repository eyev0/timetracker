FROM golang:1.21

WORKDIR /usr/app/

COPY . /usr/app

RUN make bin/timetracker

STOPSIGNAL SIGINT

CMD [ "/usr/app/bin/timetracker" ]
