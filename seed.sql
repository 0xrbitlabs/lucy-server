INSERT INTO users (
  id, username, phone_number,
  password, account_type, created_at
)
VALUES
  ('1', 'admin_user', '+1234567890', 'password123', 'admin', '2025-01-02 12:00:00'),
  ('2', 'regular_user', '+0987654321', 'securePass456', 'regular', '2025-01-02 12:05:00'),
  ('3', 'seller_user', '+1122334455', 'sellerPass789', 'seller', '2025-01-02 12:10:00')
ON CONFLICT DO NOTHING;
