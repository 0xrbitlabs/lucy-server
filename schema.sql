create table if not exists users (
  id text not null primary key,
  username text not null,
  phone_number text not null,
  password text not null,
  account_type text not null, -- can be either 'admin', 'regular' or 'seller'
  created_at timestamp not null
);

create table if not exists sessions (
  id text not null primary key,
  user_id text not null references users(id),
  valid boolean not null default true
);

create table if not exists product_categories (
  id text not null primary key,
  label text not null unique,
  description text not null,
  created_at timestamp not null
);
