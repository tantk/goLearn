FROM mysql:5.7

COPY ./data ./tmp

RUN /bin/bash -c "/usr/bin/mysqld_safe --skip-grant-tables &" && \
  sleep 5 && \
  mysql -u root -e "CREATE DATABASE courseapp" && \
  mysql -u root courseapp < /tmp/courseapp_course.sql
