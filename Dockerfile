FROM alpine:3.7 as crawler

RUN mkdir /app

COPY ./bin/crawler /app/

CMD ["/app/crawler"]

FROM alpine:3.7 as mu

RUN mkdir -p /app/public

COPY ./bin/mu /app/
COPY ./public /app/

EXPOSE 80

CMD ["/app/mu"]