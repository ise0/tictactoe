#!/bin/bash
createdb -U postgres -T template0 tictactoe;
pg_restore --verbose --clean -U postgres --dbname tictactoe /backup-tictactoe.sql;