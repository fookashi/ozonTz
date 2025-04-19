FROM postgres:16-alpine

# Копируем конфигурацию PostgreSQL
COPY config/postgresql.conf /etc/postgresql/postgresql.conf

# Копируем скрипт инициализации
COPY build/postgres/init.sh /docker-entrypoint-initdb.d/
RUN chmod +x /docker-entrypoint-initdb.d/init.sh

# Указываем конфигурационный файл
CMD ["postgres", "-c", "config_file=/etc/postgresql/postgresql.conf"]