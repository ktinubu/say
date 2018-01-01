FROM alpine

RUN apk update && apk add flite
ADD app /app
ADD client /client
ENTRYPOINT ["/app"]