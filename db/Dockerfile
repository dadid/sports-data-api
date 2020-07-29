# Dockerfile to create a postgres image with SSL key/cert
FROM postgres:latest

# On Windows root will own the files, and they will have permissions 755
COPY privatekey.key /var/lib/postgresql/privatekey.key
COPY certificate.crt /var/lib/postgresql/certificate.crt
# update the privileges on the .key, no need to touch the .crt
RUN chown 999:999 /var/lib/postgresql/privatekey.key 
RUN chmod 600 /var/lib/postgresql/privatekey.key
RUN chown postgres:postgres /var/lib/postgresql/privatekey.key