FROM postgres:9.5.2
RUN localedef -i ru_RU -c -f UTF-8 -A /usr/share/locale/locale.alias ru_RU.UTF-8
COPY ./*.sql /docker-entrypoint-initdb.d/