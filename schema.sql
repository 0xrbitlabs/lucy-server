create table if not exists users (
  id text not null primary key,
  phone text not null unique,
  username text not null,
  password text not null,
  account_type text not null default 'seller' -- can be 'seller' 'regular' 'admin'
);

create table if not exists product_category (
  id text not null primary key,
  label text not null,
  active boolean not null default true
);

create table if not exists product (
  id text not null primary key,
  label text not null,
  description text not null,
  price numerical not null,
  listed_by text not null references users(id)
);
