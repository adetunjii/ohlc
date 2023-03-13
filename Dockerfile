FROM golang:1.19-alpine AS build_stage

RUN apk add --no-cache git

# set the current working directory inside the container
WORKDIR /tmp/ohlc

COPY go.* ./

RUN go mod tidy

COPY . .

# install curl 
RUN apk add curl

# download the golang-migrate package and unzip
# RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz


# build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o bin/ohlc_app .


FROM alpine:3.9
RUN apk add ca-certificates

WORKDIR /app

COPY --from=build_stage /tmp/ohlc/bin/ohlc_app /app/ohlc_app
# COPY --from=build_stage /tmp/ohlc/migrate /app/migrate


COPY app.env .
COPY start.sh .
COPY wait-for.sh .

COPY db/migration ./migration

# ARG db_url
# ENV DB_URL = $db_url

EXPOSE 8080

CMD ["./ohlc_app"]

# ENTRYPOINT [ "/app/start.sh"]