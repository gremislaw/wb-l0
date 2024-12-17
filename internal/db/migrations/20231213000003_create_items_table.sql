-- Таблица для хранения информации о товарах
-- 20231213000003_create_items_table.sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255),
    chrt_id INT,
    track_number VARCHAR(255),
    price DECIMAL(10, 2),
    rid VARCHAR(255),
    name VARCHAR(100),
    sale INT,
    size VARCHAR(10),
    total_price DECIMAL(10, 2),
    nm_id INT,
    brand VARCHAR(100),
    status INT,
    FOREIGN KEY (order_uid) REFERENCES orders(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE items;