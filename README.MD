### Setup
For testing purposes, I use docker to run a local PostgreSql server.

I use my make file cmd and a .env to pass in params. Run those cmd's to set it up(you will need to make a .env in root)

To clean your docker instance(db) use: sudo docker rm -f my-tycoon-db(or your docker name.) or make clean-docker-db

---

### DataBase tables: 

- Player:

   - player_id : BIGINT

   - total_currency : BIGINT

   - rebirths : INTEGER

   - placed_objects : JSONB(so we can index and it's faster to read, not write.)

   - current_server_id : TEXT

   - date_last_updated : TIMESTAMPTZ

---
### Environment Setup (Testing Only)

Create a `.env` file in the root directory with the following variables:

| Variable                     | Description               | Example                                                                   |
|------------------------------|---------------------------|---------------------------------------------------------------------------|
| `PORT`                       | API server port           | `8080`                                                                    |
| `DOCKER_PORT`                | PostGresql container port | `3307`                                                                    |
| `DOCKER_PASS`                | PostgreSql password       | Generate with `openssl rand -hex 32`                                      |
| `DOCKER_DB_NAME`             | Docker container name     | `my-tycoon-db`                                                            |
| `POSTGRESS_DB_NAME`          | Database name             | `tycoon_db`                                                               |
| `DATABASE_CONNECTION_STRING` | PostgreSQL connection URL | `postgres://postgres:<password>@127.0.0.1:3307/tycoon_db?sslmode=disable` |
| `APIKEY`                     | API authentication key    | Generate with `openssl rand -hex 32`                                      |

**Generate secure credentials:**
```bash
openssl rand -hex 32  # Use for DOCKER_PASS and APIKEY
```

    