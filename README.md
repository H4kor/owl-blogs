![Mascot](assets/owl.png)

# Owl Blogs

Owl-blogs is a blogging software focused on simplicity with IndieWeb and Fediverse support.

# Usage

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

# Development

## Build

```
CGO_ENABLED=1 go build -o owl ./cmd/owl
```

For development with live reload use `air` ([has to be install first](https://github.com/cosmtrek/air))

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
