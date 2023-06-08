# Primary database (PgSQL) migration files

Please, follow the next rules for your migrations:

- Only one SQL query per file
- Rollback (down) migration must be written for each migration file
- Changes should only concern the scheme, **not** the data

## Docs

- [golang-migrate/migrate/postgres](https://github.com/golang-migrate/migrate/tree/master/database/postgres)

