FROM alpine:3.7 as crawler
ENV APP_ENV production
RUN apk add --no-cache ca-certificates tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir /app
COPY ./bin/crawler /app/
WORKDIR /app
CMD ["./crawler"]

FROM alpine:3.7 as mu
ENV APP_ENV production
RUN apk add --no-cache ca-certificates tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir -p /app/public
COPY ./bin/mu /app/
COPY ./public /app/public
WORKDIR /app
EXPOSE 7980
CMD ["./mu"]