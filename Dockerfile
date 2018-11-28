FROM alpine

WORKDIR /app

COPY gamelink-fcm  ./

RUN apk update && apk add --no-cache ca-certificates

ENTRYPOINT [ "./gamelink-fcm" ]