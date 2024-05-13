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

INSERT INTO users (id, username, phone, password, account_type, description, country, town) 
VALUES 
('1', 'admin_user', '123456789', '$2y$10$FtxxyQXJjwwST2tmKppQf.CArH84ekTLG3Ko0VXeH3vBOQC.2yrU2', 'admin', 'Administrator account', 'USA', 'New York');
