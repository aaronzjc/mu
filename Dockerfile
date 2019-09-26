FROM alpine:3.7 as crawler
ENV APP_ENV production
RUN apk add --no-cache ca-certificates tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir -p /app/bin
COPY ./bin/crawler /app/bin/
WORKDIR /app/bin
EXPOSE 7990
CMD ["./crawler"]

FROM alpine:3.7 as mu
ENV APP_ENV production
RUN apk add --no-cache ca-certificates tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir -p /app/public /app/bin
COPY ./bin/mu /app/bin
COPY ./public /app/public
WORKDIR /app/bin
EXPOSE 7980
CMD ["./mu"]