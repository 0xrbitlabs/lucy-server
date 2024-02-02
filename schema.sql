create table admins (
  id text not null primary key,
  username text not null unique,
  password text not null
);

create table users (
  id text not null primary key,
  user_type text not null,
  phone_number text not null unique,
  password text not null,
  name text not null,
  description text not null,
  country text not null,
  town text not null
);
