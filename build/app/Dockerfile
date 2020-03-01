FROM golang:1.13

ENV JONGME_PATH /go/src/jongme 

RUN mkdir -p ${JONGME_PATH}/app

WORKDIR ${JONGME_PATH}/app

COPY app ${JONGME_PATH}/app

RUN go get -v
RUN go get github.com/githubnemo/CompileDaemon

EXPOSE 8083

ENTRYPOINT CompileDaemon --log-prefix=false -build="go build -o ../build/app/jongme" -command="../build/app/jongme"
