-- Migration: Create table account_verifications

-- Create the account_verifications table
CREATE TABLE IF NOT EXISTS account_verifications (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,                -- Primary key
    account_id BIGINT NOT NULL,                          -- Foreign key referencing accounts
    type VARCHAR(50) NOT NULL,                           -- Type of verification (e.g., 'email', 'phone', '2fa', etc.)
    token VARCHAR(255) NULL UNIQUE,                  -- Unique token or OTP (hash if sensitive)
    expires_at TIMESTAMP NOT NULL,                       -- Expiration timestamp
    verified BOOLEAN DEFAULT FALSE,                      -- Verification status
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,      -- Record creation timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- Last update timestamp
    CONSTRAINT fk_account_id FOREIGN KEY (account_id) REFERENCES accounts (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

-- Check if the account_id and type index exists and create it if it doesn't
SELECT COUNT(1) INTO @IndexExists
FROM INFORMATION_SCHEMA.STATISTICS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'account_verifications' 
  AND INDEX_NAME = 'idx_account_verifications_account_type';

SET @Query = IF(@IndexExists = 0, 'CREATE INDEX idx_account_verifications_account_type ON account_verifications (account_id, type)', NULL);
PREPARE stmt FROM @Query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Check if the token index exists and create it if it doesn't
SELECT COUNT(1) INTO @IndexExists
FROM INFORMATION_SCHEMA.STATISTICS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'account_verifications' 
  AND INDEX_NAME = 'idx_account_verifications_token';

SET @Query = IF(@IndexExists = 0, 'CREATE UNIQUE INDEX idx_account_verifications_token ON account_verifications (token)', NULL);
PREPARE stmt FROM @Query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
