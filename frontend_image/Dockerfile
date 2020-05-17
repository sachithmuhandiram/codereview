FROM httpd:alpine

#RUN rc-service apache2 start
#RUN apk add openrc
COPY ./views /usr/local/apache2/htdocs/

COPY httpd.conf /usr/local/apache2/conf/httpd.conf

#ENV APACHE_DOCUMENT_ROOT=/var/www/html
