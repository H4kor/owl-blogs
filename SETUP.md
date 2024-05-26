# Running owl-blogs

Owl-blogs uses a single SQLite database to store all data and configuration of a blog.
This database is automatically created in the current directory, named `owlblogs.db`, when running any command.


## Setup

Configuration of the blog is done via the UI.
To access the configuration you need an author account which can be created with the command:

```
owl new-author -u <name> -p <password>
```

Afterwards start the web server with:

```
owl web
```

This command starts the webserver listening on the address **:3000**.
Owl-blogs should be run behind a reverse proxy, for example [caddy](https://caddyserver.com/).

The main configuration can be found in the [admin menu](http://localhost:3000) (use the link in the footer of the blog) as "Site Settings".
Set the *"Full Url"* value to your domain, including protocol (e.g. `https://blog.example.com`).
Naming and appearance of the blog can be controlled here.
Additional raw HTML can be added to the `<head>` (e.g. additional CSS or scripts) and to the `<footer>`. 


### Configure Caddy as Reverse Proxy

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


