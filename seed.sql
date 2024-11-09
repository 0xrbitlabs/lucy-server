-- create admin user
-- password is 'password'
insert into users (
  id, phone, username, password, account_type
)
values (
  '1', '0', 'admin1',
  '$2y$10$xBJcodAk96RUkv6tKqXLXuC0xl3/PTQxhY.oEiMirZ14PsaBksCjO',
  'admin'
);
