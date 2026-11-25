CREATE UNIQUE INDEX IF NOT EXISTS idx_cart_items_cart_id_product_id
ON cart_items(cart_id, product_id);
