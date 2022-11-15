FROM alpine:3.7
ENV APP_ENV prod
RUN apk add --no-cache ca-certificates tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir -p /app/bin /app/conf
COPY ./dagger/backend/commander /app/bin
EXPOSE 7970
VOLUME /app/conf
WORKDIR /app
CMD ["./bin/commander", "-c", "conf/prod.yml"]