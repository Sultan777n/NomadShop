CREATE TABLE user_roles (
                            id SERIAL PRIMARY KEY,
                            user_id INT REFERENCES users(id) ON DELETE CASCADE,
                            role_id INT REFERENCES roles(id) ON DELETE CASCADE,
                            UNIQUE(user_id, role_id)
);
