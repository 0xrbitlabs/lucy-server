create table admins (
  id text not null primary key,
  username text not null unique,
  password text not null
);

create table users (
  id text not null primary key,
  user_type text not null default 'seller',
  phone_number text not null unique,
  password text not null,
  verified boolean not null default false,
  name text not null,
  profile_picture text not null default 'https://picsum.photos/200/300',
  description text not null,
  country text not null,
  town text not null
);

create table categories (
  id text not null primary key,
  name text not null unique,
  scope text not null default 'public'
);

create table products (
  id text not null primary key,
  user_id text not null references users(id),
  category_id text not null references categories(id),
  label text not null,
  image text not null,
  description text not null,
  price numeric not null,
  deliverable boolean not null default true,
  sold_out boolean not null default false
);

create table verification_codes (
  id text not null primary key,
  user_id text not null references users(id),
  code text not null,
  sent_at timestamp not null,
  used boolean not null default false
);
