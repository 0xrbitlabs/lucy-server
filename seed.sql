-- create admin user
-- password is always 'password'
insert into users (
  id, phone, username, password, account_type
)
values (
  '1', '0', 'admin1',
  '$2y$10$xBJcodAk96RUkv6tKqXLXuC0xl3/PTQxhY.oEiMirZ14PsaBksCjO',
  'admin'
) on conflict do nothing;

-- create seller account
insert into users (
  id, phone, username, password, account_type
)
values (
  '2', '1', 'seller1',
  '$2y$10$xBJcodAk96RUkv6tKqXLXuC0xl3/PTQxhY.oEiMirZ14PsaBksCjO',
  'seller'
) on conflict do nothing;
