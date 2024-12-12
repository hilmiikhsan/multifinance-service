-- +goose Up
-- +goose StatementBegin
CREATE TABLE credit_limits (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    tenor_month INT NOT NULL,
    limit_amount DECIMAL(15,2) NOT NULL,
    UNIQUE KEY unique_customer_tenor (customer_id, tenor_month),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_credit_limits_customer_id ON credit_limits (customer_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS credit_limits;
-- +goose StatementEnd
