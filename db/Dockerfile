# Dockerfile to create a postgres image with SSL key/cert
FROM postgres:latest
LABEL "Product"="Postgres (SSL)"

COPY ca/server.key /var/lib/postgresql/server.key
COPY ca/server.crt /var/lib/postgresql/server.crt
RUN chown postgres /var/lib/postgresql/server.key && \
    chmod 600 /var/lib/postgresql/server.key
RUN echo "host all  all    0.0.0.0/0  md5" >> /var/lib/postgresql/data/pg_hba.conf
RUN echo "listen_addresses='*'" >> /var/lib/postgresql/data/postgresql.conf

# commands to adjust server.key ownership prior to image creation:
# chown 999:999 server.key # set postgres as owner; 70:70 for alpine; 999:999 for debian
# chmod 600 server.key
# chown 0:70 server.key # set root as owner
# chmod 640 server.key