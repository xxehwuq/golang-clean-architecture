# golang-clean-architecture

Golang Clean Architecture Implementation

[API Documentation](https://documenter.getpostman.com/view/20867991/2sA2rAzhbA)

### Setup

1. Create `.env` file and set values:

```dotenv
LOG_LEVEL=<debug>

POSTGRES_URL=<postgres://jack:secret@pg.example.com:5432/mydb>

TOKENS_SIGNING_KEY=<secret-key>

PASSWORD_SALT=<secret-salt>

REDIS_URL=<redis://user:pass@localhost:6379/db>
```

2. Run:

```shell
make
```
