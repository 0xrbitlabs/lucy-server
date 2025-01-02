create table users (
  id text not null primary key,
  username text not null,
  phone_number text not null,
  password text not null,
  user_type text not null, -- can be either 'admin', 'regular' or 'seller'
  created_at timestamp not null
);

create table sessions (
  id text not null primary key,
  user_id text not null references users(id),
  valid boolean not null default true
);
