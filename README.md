![Mascot](assets/owl.png)

# Owl Blogs

A simple web server for blogs generated from Markdown files.

**_This project is not yet stable. Expect frequent breaking changes! Only use this if you are willing to regularly adjust your project accordingly._**



## Build

```
CGO_ENABLED=1 go build -o owl ./cmd/owl
```

## Run

To run the web server use the command:

```
owl web
```

The blog will run on port 3000 (http://localhost:3000)

To create a new account:

```
owl new-author -u <name> -p <password>
```

To retrieve a list of all commands run:

```
owl -h
```