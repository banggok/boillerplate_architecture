-- Create table if it does not exist
CREATE TABLE IF NOT EXISTS accounts (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,         -- Primary key with auto-increment
    name VARCHAR(100) NOT NULL,                   -- Account name, cannot be null
    email VARCHAR(255) NOT NULL UNIQUE,           -- Unique and indexed email
    phone VARCHAR(20) NOT NULL UNIQUE,            -- Unique and indexed phone number
    password TEXT NOT NULL,                       -- hashed password
    tenant_id BIGINT,                             -- Foreign key referencing tenants (nullable for SET NULL)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Automatically set when the record is created
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- Automatically updated
    CONSTRAINT fk_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE                        -- Foreign key constraints
);

-- Check if the email index exists and create it if it doesn't
SELECT COUNT(1) INTO @IndexExists
FROM INFORMATION_SCHEMA.STATISTICS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'accounts' 
  AND INDEX_NAME = 'idx_account_email';

SET @Query = IF(@IndexExists = 0, 'CREATE UNIQUE INDEX idx_account_email ON accounts (email)', NULL);
PREPARE stmt FROM @Query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Check if the phone index exists and create it if it doesn't
SELECT COUNT(1) INTO @IndexExists
FROM INFORMATION_SCHEMA.STATISTICS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'accounts' 
  AND INDEX_NAME = 'idx_account_phone';

SET @Query = IF(@IndexExists = 0, 'CREATE UNIQUE INDEX idx_account_phone ON accounts (phone)', NULL);
PREPARE stmt FROM @Query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Check if the tenant_id index exists and create it if it doesn't
SELECT COUNT(1) INTO @IndexExists
FROM INFORMATION_SCHEMA.STATISTICS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'accounts' 
  AND INDEX_NAME = 'idx_tenant_id';

SET @Query = IF(@IndexExists = 0, 'CREATE INDEX idx_tenant_id ON accounts (tenant_id)', NULL);
PREPARE stmt FROM @Query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
