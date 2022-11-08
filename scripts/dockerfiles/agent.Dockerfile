FROM alpine:3.7
ENV APP_ENV production
RUN apk add --no-cache ca-certificates tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir -p /app/bin
COPY ./dagger/backend/agent /app/bin/
WORKDIR /app/bin
EXPOSE 7990
CMD ["./agent"]