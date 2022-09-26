# docker-easy-owntracks

:warning: This is not ready to run in production yet, please only use it for testing.

After running the Dockerfile, the image needs three things:

Port 80/tcp -> 80

Port 80/udp -> 80

Port 443/tcp -> 443

Port 443/udp -> 443

Directory /owntracks-storage -> /owntracks-storage

It should autoconfigure and by default only be accessible on `localhost`

You can change this and other stuff in the `/owntracks-storage/configuration.conf` file.
