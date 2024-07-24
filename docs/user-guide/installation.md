# Installation

## Linux

### Binary

- Download the [latest owl-blogs release](https://github.com/H4kor/owl-blogs/releases) 
```
wget https://github.com/H4kor/owl-blogs/releases/download/v0.3.4/owl-linux-amd64 --show-progress --progress=bar -O owl`
```

- Mark the file as executable
```
chmod +x owl
```
- Start the owl. This will run the server bound to port 3000, listening on localhost.
```
./owl web
```
- Open http://localhost:3000 to confirm the server is running

### Docker

You can also use the [owl-blogs container image](https://github.com/H4kor/owl-blogs/pkgs/container/owl-blogs).

This will start the web server and expose it on port 3000, listening on localhost. The directory `./owl-data` will be mounted to the container to store your blog's data.

```
docker run -p 127.0.0.1:3000:3000/tcp -v ./owl-data:/owl --name owl --restart=unless-stopped -d ghcr.io/h4kor/owl-blogs:v0.3.4 web -b :3000
```

Open http://localhost:3000 to confirm the server is running

When you continue with the docker setup any owl command in the following guides have to be executed in the container. 
Try to get the help message with `docker exec owl owl help`.

### Reverse Proxy

Owl-blogs should be run behind a reverse proxy, for example [caddy](https://caddyserver.com/).

#### Configure Caddy as Reverse Proxy

> *[How to install caddy](https://caddyserver.com/docs/install)*


After installing Caddy, insert the following snippet into your Caddyfile.
Change `blog.example.com` to your own domain name.
You may have to restart caddy afterwards for the changes to take effect.

```
blog.example.com {
    encode {
        gzip
        zstd
    }
    reverse_proxy localhost:3000
    log {
        output file /var/log/blog.log
    }
}
```
