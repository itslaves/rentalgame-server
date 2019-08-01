FROM golang:1.12-alpine as BUILD

WORKDIR /app
COPY . /app/
RUN apk add --no-cache --virtual .build-deps git
RUN go build -o rg-server
RUN apk del .build-deps

FROM alpine:3.10

WORKDIR /app
COPY --from=BUILD /app/rg-server /app/rg-server
COPY --from=BUILD /app/config /app/config

EXPOSE 8000
ENTRYPOINT [ "./rg-server" ]