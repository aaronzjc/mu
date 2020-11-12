FROM alpine:3.7 as commander
ENV APP_ENV production
RUN apk add --no-cache ca-certificates tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir -p /app/bin /app/conf
COPY ./bin/commander /app/bin
COPY ./conf /app/conf
EXPOSE 7970
VOLUME /app/conf
WORKDIR /app/bin
CMD ["./commander"]

FROM alpine:3.7 as agent
ENV APP_ENV production
RUN apk add --no-cache ca-certificates tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir -p /app/bin
COPY ./bin/agent /app/bin/
WORKDIR /app/bin
EXPOSE 7990
CMD ["./agent"]

FROM alpine:3.7 as api
ENV APP_ENV production
RUN apk add --no-cache ca-certificates tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir -p /app/bin /app/conf /app/public
COPY ./bin/api /app/bin
COPY ./conf /app/conf
COPY ./public /app/public
EXPOSE 7980
VOLUME /app/conf
WORKDIR /app/bin
CMD ["./api"]

FROM nginx:stable-alpine as frontend
COPY ./public /usr/share/nginx/html
COPY ./deploy/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
VOLUME /usr/share/nginx/html
CMD ["nginx", "-g", "daemon off;"]