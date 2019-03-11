# Mazewire
--

Mazewire is a modern content management system and framework for websites and app backends.

## Running Mazewire and bootstrapping

Use commands prepared in the Makefile (run each under the `make` command):
 * `gen`: generate Protocol Buffers code for Go.
 * `crdb_setup`: set up the database and tables for a local CockroachDB node.
 * `crdb`: start the CockroachDB database.
 * `crdb_sql`: start the Cockroach SQL client.
 * `start`: build and run the app, watching for changes and reloading.

## About Mazewire

The goal of this project is to provide a free and open source content management system and framework for
building websites and app backends. Security and reliability are the primary focus. The extensible plugin
architecture is meant to make it safe to run third-party plugins without compromising the main process.

You can run Mazewire on a single node, or you can run a cluster of stateless Mazewire instances which will
orchestrate themselves for any tasks that your system needs to do asynchronously or regularly, such as a
scheduled job.

## Technologies used

The core is written in Go, chosen for its speed, strict typing, and support for Protocol Buffers and gRPC.
Plugins can be written in any language so long as they can speak gRPC over a Unix domain socket. The admin
dashboard is a single-page app built on Vue, but there's a chance it'll be replaced with React to make it
easier for more developers to enrich the admin dash interface with plugins. External communication is done
in a plain REST-like fashion using Protocol Buffers encoding, but we'd like to add gRPC as browsers improve
their support for HTTP/2 and get better at handling binary data.

The `sql` directory describes the database schema, and the `data` directory has a Protocol Buffers message
type for for each table. You can use any kind of database for which you have an implementation of the
[`data.DB`](https://godoc.org/github.com/dchenk/mazewire/pkg/data) interface. This project prefers CockroachDB
for its scalability and ACID transaction guarantees. We're working on support for PostgreSQL and MySQL.
