FROM alpine:3.7
ENV APP_ENV prod
RUN apk add --no-cache ca-certificates tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir -p /app/bin /app/conf /app/public
COPY ./dagger/backend/api /app/bin
COPY ./dagger/frontend /app/public
EXPOSE 7980
VOLUME /app/conf
WORKDIR /app
CMD ["./bin/api", "-c", "conf/prod.yml"]