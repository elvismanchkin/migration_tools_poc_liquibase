FROM postgres:16-alpine

COPY config/postgresql.conf /etc/postgresql/postgresql.conf

COPY init-scripts/ /docker-entrypoint-initdb.d/

CMD ["postgres", "-c", "config_file=/etc/postgresql/postgresql.conf"]