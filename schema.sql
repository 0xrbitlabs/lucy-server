create table users (
  id text not null primary key,
  phone_number text not null unique,
  full_name text not null,
  profile_picture text not null default ''
);
