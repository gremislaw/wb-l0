-- Таблица для хранения информации о доставке
-- 20231213000001_create_deliveries_table.sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE deliveries (
    order_uid VARCHAR(255),
    name VARCHAR(100),
    phone VARCHAR(20),
    zip VARCHAR(20),
    city VARCHAR(100),
    address VARCHAR(255),
    region VARCHAR(100),
    email VARCHAR(100),
    FOREIGN KEY (order_uid) REFERENCES orders(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE deliveries;