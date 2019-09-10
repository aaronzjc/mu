FROM alpine:3.7 as crawler
RUN apk add --no-cache ca-certificates tzdata
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir /app
COPY ./bin/crawler /app/
WORKDIR /app
CMD ["./crawler"]

FROM alpine:3.7 as mu
RUN apk add --no-cache ca-certificates tzdata
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir -p /app/public
COPY ./bin/mu /app/
COPY ./public /app/public
WORKDIR /app
EXPOSE 7980
CMD ["./mu"]