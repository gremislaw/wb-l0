-- Таблица для хранения информации о платеже
-- 20231213000002_create_payments_table.sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE payments (
    order_uid VARCHAR(255),
    transaction VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(10),
    provider VARCHAR(50),
    amount DECIMAL(10, 2),
    payment_dt INT,
    bank VARCHAR(50),
    delivery_cost DECIMAL(10, 2),
    goods_total DECIMAL(10, 2),
    custom_fee DECIMAL(10, 2),
    FOREIGN KEY (order_uid) REFERENCES orders(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE payments;