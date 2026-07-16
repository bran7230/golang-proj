## Setup
For testing purposes, i use docker to run a local MYSQP server.

I use my make file cmds and a .env to pass in params. Run those cmd's to set it up(you will need to make a .env in root)

To clean your docker instance(db) use: sudo docker rm -f my-tycoon-db(or your docker name.) or make clean-docker-db

### DataBase tables: 

- Player:

   - player_id : BIGINT

   - total_currency : BIGINT

   - rebirths : INTEGER

   - placed_objects : JSONB(so we can index and it's faster to read, not write.)

   - current_server_id : TEXT

   - date_last_updated : TIMESTAMPTZ
    