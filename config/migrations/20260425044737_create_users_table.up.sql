CREATE TABLE products (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(150) NOT NULL,
                          price NUMERIC(10,2) NOT NULL,
                          stock INT DEFAULT 0,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);