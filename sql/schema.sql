CREATE TYPE user_type AS ENUM ('super_admin', 'admin', 'seller', 'regular');

create table users (
  id text not null primary key,
  type user_type not null,
  phone_number text not null unique,
  password text not null,
  username text not null,
  description text not null,
  country text not null,
  town text not null
);

create table categories (
  id text not null primary key,
  label text not null unique,
  description text not null,
  enabled bool not null default false
);
