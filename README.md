# The Personal Finances API
Expenses is an API in Go for managing your monthly expenses, initially intended as 
project for practing web services creation in Go, although fully-featured and production-ready.

Still very much a work in progress.

## Features
- Multiple users with authentication for each
- User activation via email
- Configurable CORS handling
- More soon

## Configuration

The application is configured via command-line flags.
You can see a list of all available flags by starting the
application with the -help flag.

```
make build/api
./bin/api -help
```

### Database 

Expenses uses the PostgreSQL database. By supplying the
`-db-dsn` flag, you especify to which PostgreSQL instance
the API should connect to at startup. 

Connection Pool settings can be edited using the rest of
the `-db` flags.

### SMTP

SMTP is used to send activation tokens in order to verify
an user's email address. An example of SMTP configuration can
be seen on the Makefile's `run/api` target:

```
exec bin/api \
	-db-dsn=${EXPENSES-DB-DSN} \
	-smtp-host=${EXPENSES-SMTP-HOST} \
	-smtp-port=${EXPENSES-SMTP-PORT} \
	-smtp-username=${EXPENSES-SMTP-USERNAME} \
	-smtp-password=${EXPENSES-SMTP-PASSWORD} \
	-smtp-sender=${EXPENSES-SMTP-SENDER}
```

These SMTP credentials can be gotten from a email service such
as Mailtrap or SendGrid.

### CORS Handling

Cross Origin Requests handling can be configured by the 
-cors-trusted-origins flag. By setting it to "*", you
especify that every origin should be trusted, which is
useful if you wish that your Expenses instance is publically 
available to every browser frontend. However, if you want
to restrict it's access to only a few selected frontends,
set `-cors-trusted-origins` to a space separated list
of the frontend origins.
