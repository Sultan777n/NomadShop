CREATE TABLE orders (
                        id SERIAL PRIMARY KEY,
                        user_id INT NOT NULL REFERENCES users(id),
                        order_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        status VARCHAR(50) NOT NULL DEFAULT 'pending',
                        total DECIMAL(10,2) NOT NULL
);
