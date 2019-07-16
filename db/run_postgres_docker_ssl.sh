# create ssl key and cert; -nodes creates a key without a passphrase
openssl req -x509 -nodes -days 365 -sub '/CN=localhost' -newkey rsa:4096 -keyout server.key -out server.crt

# set postgres user as owner of the server.key and set permissions to 600; user is 0:70 for alpine variant, 999:999 for debian
chown 999:999 server.key
chmod 600 server.key

container_name=localpostgres1
pass=Hotdog10!
volume_name=postgresdata

# only works on linux
docker run -d --name $container_name \
    -p 5432:5432 \
    -e POSTGRES_PASSWORD=$pass \
    -v $PWD/postgres/data:/var/lib/postgresql/data \
    -v $PWD/server.crt:/var/lib/postgresql/server.crt:ro \
    -v $PWD/server.key:/var/lib/postgresql/server.key:ro \
  postgres:latest \
    -c ssl=on \
    -c ssl_cert_file=/var/lib/postgresql/server.crt \
    -c ssl_key_file=/var/lib/postgresql/server.key

# for windows - create SSL certs; use dockerfile to create postgres image that copies the certs; run command below *DOESN'T WORK*
# docker run -d --name $container_name \
#     -p 5432:5432 \
#     -e POSTGRES_PASSWORD=$pass \
#     --volume $volume_name:/var/lib/postgresql/data \
#   postgres:latest \
#     -c ssl=on \
#     -c ssl_cert_file=/var/lib/postgresql/server.crt \
#     -c ssl_key_file=/var/lib/postgresql/server.key