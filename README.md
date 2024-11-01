# car-simulator

[![GoTemplate](https://img.shields.io/badge/go/template-black?logo=go)](https://github.com/SchwarzIT/go-template)

Exploring concurrency in golang by simulating the electronics behind a car dashboard.

---

This simulation uses a three-layered architecture (presentation, application and data).
Small controllers are responsible for implementing a small part of the application logic

Input to the application layer controllers is provided only through events on the shared event bus.
There are two sources of events:

- Clicking a control in the user interface
- Controllers

![architecture](.github/architecture.png)

```bash
make help
```

## Setup

To get your setup up and running the only thing you have to do is

```bash
make all
```

This will initialize a git repo, download the dependencies in the latest versions and install all needed tools.
If needed code generation will be triggered in this target as well.

## Test & lint

Run linting

```bash
make lint
```

Run tests

```bash
make test
```

## Credits

<a href="https://www.vecteezy.com/free-vector/car-dashboard-icons">Car Dashboard Icons Vectors by Vecteezy</a>
