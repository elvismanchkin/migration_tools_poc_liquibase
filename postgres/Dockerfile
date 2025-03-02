FROM postgres:15-alpine

COPY config/postgresql.conf /etc/postgresql/postgresql.conf

COPY init-scripts/ /docker-entrypoint-initdb.d/

CMD ["postgres", "-c", "config_file=/etc/postgresql/postgresql.conf"]