# The Personal Finances API
Expenses is an API in Go for managing your monthly expenses, initially intended as 
project for practing web services creation in Go, although fully-featured and production-ready.

Still very much a work in progress.

## Features
- Multiple users with authentication for each
- User activation via email
- Configurable CORS handling
- Built-in rate limiter
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

Expenses uses the PostgreSQL database. Supply a value to the
the `-db-dsn` flag via the environment variable `EXPENSES_DB_DSN`
to especify to which PostgreSQL instance the API should connect 
to at startup.

Connection Pool settings can be edited using the rest of
the `-db` flags.

### SMTP

SMTP is used to send activation tokens in order to verify
an user's email address. An example of SMTP configuration can
be seen on the Makefile's `run/api` target:

```
exec bin/api \
	-db-dsn=${EXPENSES-DB-DSN} \
	-smtp-host=${EXPENSES_SMTP_HOST} \
	-smtp-port=${EXPENSES_SMTP_PORT} \
	-smtp-sender=${EXPENSES_SMTP_SENDER}
```

Note that no values are passed to the `-smtp-username` and `-smtp-password`
flags. Just as `-db-dsn`, for security reasons, they receive their values
via the environment variables `EXPENSES_SMTP_USERNAME` and `EXPENSES_SMTP_PASSWORD`.


However, it is recommended that you pass 

These SMTP credentials can be gotten from a email service such
as Mailtrap or SendGrid.


### Rate Limiting 

Rate limiting is enabled by default. To disable it, set `-limiter-enabled` to false.

### CORS Handling

Cross Origin Requests handling can be configured by the 
`-cors-trusted-origins` flag. By setting it to "*", you
especify that every origin should be trusted, which is
useful if you wish that your Expenses instance is publically 
available to every browser frontend. However, if you want
to restrict it's access to only a few selected frontends,
set `-cors-trusted-origins` to a space separated list
of the frontend origins.
