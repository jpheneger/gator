# Gator - RSS Feed Blog aggreGATOR

Backed by PostgreSQL, leveraging sqlc for db code generation, migrations handled through Goose.

## Installation Requirements:
- PostgreSQL
    - macOS: `brew install postgresql@15`
    - Linux / WSL (Debian): `sudo apt update; sudo apt install postgresql postgresql-contrib`
- GO
- goose
    - `go install github.com/pressly/goose/v3/cmd/goose@latest`
- sqlc
    - `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

## Post-Install:
- Run `psql --version` to ensure psql is installed properly.
- (Linux only) Update postgres password: `sudo passwd postgres`
- Start the Postgres server in the background
    - macOS: `brew services start postgresql@15`
    - Linux/WSL: `sudo service postgresql start`
- Connect to the server - enter the `psql` shell:
    - macOS: `psql postgres`
    - Linux: `sudo -u postgres psql`

You should see a new prompt that looks like this:
`postgres-#`

- Create a new database:
    - `CREATE DATABASE gator;`

- Connect to the new database:
    - `\c gator`

You should see a new prompt that looks like this:
`gator=#`

- (Linux only) set the user password:
    - `ALTER USER postgres PASSWORD 'postgres';`

- Query the database:
    - `SELECT version();`

All otther changes to the databse are managed by the code and/or migration scripts.