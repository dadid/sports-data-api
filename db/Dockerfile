FROM postgres:latest

# On Windows root will own the files, and they will have permissions 755
COPY server.key /var/lib/postgresql/server.key
COPY server.crt /var/lib/postgresql/server.crt
# update the privileges on the .key, no need to touch the .crt; user is 0:70 for alpine variant, 999:999 for debian
RUN chmod 600 /var/lib/postgresql/server.key
RUN chown postgres:postgres /var/lib/postgresql/server.key