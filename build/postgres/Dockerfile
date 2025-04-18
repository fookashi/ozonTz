FROM postgres:16-alpine

COPY config/postgresql.conf /etc/postgresql/postgresql.conf
COPY build/postgres/init.sh /docker-entrypoint-initdb.d/init.sh
RUN chmod +x /docker-entrypoint-initdb.d/init.sh

CMD ["postgres", "-c", "config_file=/etc/postgresql/postgresql.conf"]