# LinkStowr API

API for the [LinkStowr obsidian plugin](https://github.com/joelseq/obsidian-linkstowr), [web app](https://github.com/joelseq/linkstowr-web), and [chrome extension](https://github.com/joelseq/linkstowr-extension).

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests

```bash
make all
```

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Live reload the application:

```bash
make watch
```

Run sqlc to generate Go code from the schema and query:

```bash
make db
```

Run the test suite:

```bash
make test
```

Clean up binary from the last build:

```bash
make clean
```
