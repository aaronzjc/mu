FROM nginx:stable-alpine
COPY ./dagger/frontend /usr/share/nginx/html
COPY ./scripts/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
VOLUME /usr/share/nginx/html
CMD ["nginx", "-g", "daemon off;"]