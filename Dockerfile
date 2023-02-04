ARG REPO_LOCATION=registry.symper.vn
ARG BASE_VERSION=1.0
FROM golang:1.16.7
RUN mkdir /app
WORKDIR /app
ADD ./ ./

RUN go mod download
RUN go build -o data-connector
RUN chmod -R 777 /app/log
RUN chmod -R 777 /app
EXPOSE 1323
WORKDIR /src

CMD [ "/app/data-connector" ]