![Mascot](assets/owl.png)

# Owl Blogs

Owl-blogs is a blogging software focused on simplicity with IndieWeb and Fediverse support.

# Usage

**Detailed information can be found in  [Setup](SETUP.md)**

To retrieve a list of all commands run:

```
owl -h
```


# Development

## Build

```
CGO_ENABLED=1 go build -o owl ./cmd/owl
```

For development with live reload use `air` ([has to be install first](https://github.com/air-verse/air))

## Tests

The project has two test suites; "unit tests" written in go and "end-to-end tests" written in python.

### Unit Tests

```
go test ./...
```

### End-to-End tests

- Start the docker compose setup in the `e2e_tests` directory.
- Install the python dependencies into a virtualenv
```
cd e2e_tests
python3 -m venv venv
. venv/bin/activate
pip install -r requirements.txt
```
- Run the e2e_tests with `pytest`
