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
RUN mkdir -p /app/bin /app/conf app/public
COPY ./bin/mu /app/bin
COPY ./conf /app/conf
COPY ./public /app/public
EXPOSE 7980
VOLUME /app/conf
WORKDIR /app/bin
CMD ["./mu"]

FROM nginx:stable-alpine as mu-frontend
COPY --from=mu /app/public /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]