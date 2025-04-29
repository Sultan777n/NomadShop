CREATE TABLE favorite_items (
                                id SERIAL PRIMARY KEY,
                                user_id INT REFERENCES users(id) ON DELETE CASCADE,
                                product_id INT REFERENCES products(id),
                                UNIQUE (user_id, product_id)
);
