### Presentation

This program is a quick example about how to handle a websocket and a queue:

- Integrates events via RabbitMQ
- Offers an api for subscribing to these events in websockets (with several possible clients subscribing)
- Offers a replay-policy, allowing N past events to be replayed when a client connects
- Store these events in a postgres DB, so you can replay the last N events, even if the server reboots.

### Getting started

A docker compose file is included to launch a PG database and a RabbitMQ event bus.

As well, 4 binaries are included:

- `cmd/migrate/main.go` to do the DB migration (table creation,...)
- `cmd/core/main.go` the main program
- `cmd/RQ/main.go` a helper that sends messages into the RQ (it reads STDIN).
   Messages should have the following format: ` {"message":"YO"}`
- `cmd/WS/main.go` a client that reads and displays every message sent through the websocket. 

