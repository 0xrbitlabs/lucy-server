create table admins (
  id text primary key,
  username text not null unique,
  password text not null,
  is_super boolean not null default false
);

insert into admins (id, username, password, is_super)
values ('1','super', '$2y$08$csvIkjxk6fCR9CGVp0tCpeGFyVRt0iE9PFIVyKkXPX0iyo0XFUFZW', true);

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

create table categories (
  id text not null primary key,
  label text not null unique,
  description text not null,
  enabled bool not null default false
);
