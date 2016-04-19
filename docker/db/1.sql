CREATE ROLE browsercall LOGIN UNENCRYPTED PASSWORD 'browsercall';
CREATE DATABASE browsercall WITH ENCODING 'UTF8' LC_COLLATE 'ru_RU.utf8'  LC_CTYPE 'ru_RU.utf8' OWNER browsercall TEMPLATE template0;

\c browsercall;

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(256),
  password VARCHAR(256)
);

CREATE UNIQUE INDEX ON users(username);