FROM alpine:3.7 as crawler
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -g app app
RUN mkdir /app
COPY ./bin/crawler /app/
RUN chown -R app:app /app
USER app
CMD ["/app/crawler"]

FROM alpine:3.7 as mu
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -g app app
RUN mkdir -p /app/public
COPY ./bin/mu /app/
COPY ./public /app/
RUN chown -R app:app /app
USER app
EXPOSE 80
CMD ["/app/mu"]