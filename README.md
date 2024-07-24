![Mascot](assets/owl.png)

# Owl Blogs

Owl-blogs is a blogging software focused on simplicity with IndieWeb and Fediverse support.

# Usage

Full Documentation can be found on the [owl-blogs website](https://h4kor.github.io/owl-blogs/)

- [Installation](https://h4kor.github.io/owl-blogs/user-guide/installation/)
- [Setup](https://h4kor.github.io/owl-blogs/user-guide/setup/)

# Development

## Build

```
CGO_ENABLED=1 go build -o owl ./cmd/owl
```

For development with live reload use `air` ([has to be installed first](https://github.com/air-verse/air))

## Tests

All tests are implemented in go and can be executed by using:

```
go test ./...
```

## Publishing

1. Update `OWL_VERSION` number in `config/config.go`
2. Push to main branch
3. Create Release with same version number
4. GitHub Actions build binary and add them to the release