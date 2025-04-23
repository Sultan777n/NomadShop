CREATE TABLE products (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR NOT NULL,
                          price INTEGER NOT NULL,
                          description TEXT NOT NULL,
                          image TEXT NOT NULL,
                          color VARCHAR NOT NULL,
                          size VARCHAR NOT NULL,
                          category_id INTEGER NOT NULL REFERENCES categories(id),
                          stock INTEGER NOT NULL
);
