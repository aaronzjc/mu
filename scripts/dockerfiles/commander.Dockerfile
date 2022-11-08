FROM alpine:3.7
ENV APP_ENV production
RUN apk add --no-cache ca-certificates tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir -p /app/bin /app/conf
COPY ./dagger/backend/commander /app/bin
COPY ./conf /app/conf
EXPOSE 7970
VOLUME /app/conf
WORKDIR /app/bin
CMD ["./commander"]