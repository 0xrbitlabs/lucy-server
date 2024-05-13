create table if not exists users (
  id text not null primary key,
  username text not null,
  phone text not null unique,
  password text not null,
  account_type text not null,
  description text not null,
  country text not null default '',
  town text not null default ''
);

create table if not exists categories (
  id text not null primary key,
  label text not null unique,
  description text not null,
  enabled boolean not null
);

create table if not exists products (
  id  text not null primary key,
  category_id text not null references categories(id),
  label text not null,
  description text not null,
  price numeric not null,
  image text not null
);
