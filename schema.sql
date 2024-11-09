create table if not exists users (
  id text not null primary key,
  phone text not null unique,
  username text not null,
  password text not null,
  account_type text not null default 'seller' -- can be 'seller' 'regular' 'admin'
);

create table if not exists product_categories (
  id text not null primary key,
  label text not null unique,
  active boolean not null default true
);

create table if not exists products (
  id text not null primary key,
  label text not null,
  category text not null references product_categories(id),
  description text not null,
  price numeric not null,
  listed_by text not null references users(id)
);

create table if not exists sessions (
  id text not null primary key,
  valid boolean not null default true,
  user_id text not null references users(id)
);

create table if not exists auth_codes (
  id serial not null,
  code text not null,
  used boolean not null default false,
  generated_for text not null
);
