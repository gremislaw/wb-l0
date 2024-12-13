-- Таблица для хранения информации о заказах
-- 20231213000000_create_orders_table.sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255),
    entry VARCHAR(50),
    locale VARCHAR(10),
    internal_signature TEXT,
    customer_id VARCHAR(50),
    delivery_service VARCHAR(50),
    shardkey VARCHAR(10),
    sm_id INT,
    date_created TIMESTAMP,
    oof_shard VARCHAR(10)
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE orders;
