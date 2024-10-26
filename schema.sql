create table if not exists users (
  id text not null primary key,
  phone text not null unique,
  username text not null,
  password text not null,
  account_type text not null default 'seller' -- can be 'seller' 'regular' 'admin'
);
