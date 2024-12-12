CREATE TABLE IF NOT EXISTS customers (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nik VARCHAR(16) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    legal_name VARCHAR(255) NOT NULL,
    birth_place VARCHAR(100),
    birth_date DATE,
    salary DECIMAL(15,2),
    ktp_photo_path VARCHAR(255),
    selfie_photo_path VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE credit_limits (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    tenor_month INT NOT NULL,
    limit_amount DECIMAL(15,2) NOT NULL,
    UNIQUE KEY unique_customer_tenor (customer_id, tenor_month),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS transactions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    contract_number VARCHAR(50) NOT NULL UNIQUE,
    on_the_road_price DECIMAL(15,2),
    admin_fee DECIMAL(15,2),
    installment_amount DECIMAL(15,2),
    interest_amount DECIMAL(15,2),
    asset_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX idx_customers_nik ON customers (nik);
CREATE INDEX idx_credit_limits_customer_id ON credit_limits (customer_id);
CREATE INDEX idx_transactions_customer_id ON transactions (customer_id);
