-- Create table if it does not exist
CREATE TABLE IF NOT EXISTS tenants (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,        -- Primary key with auto-increment
    name VARCHAR(255) NOT NULL,                  -- Tenant's name, cannot be null
    address VARCHAR(255),                        -- Optional, maximum 255 characters
    email VARCHAR(255) NOT NULL,                 -- Unique and indexed email
    phone VARCHAR(20) NOT NULL,                  -- Unique and indexed phone number
    timezone VARCHAR(100) NOT NULL,              -- Required timezone
    opening_hours TIME,                          -- Optional opening hours
    closing_hours TIME,                          -- Optional closing hours
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Automatically set when the record is created
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- Automatically updated
);

-- Check if the email index exists and create it if it doesn't
SELECT COUNT(1) INTO @IndexExists
FROM INFORMATION_SCHEMA.STATISTICS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'tenants' 
  AND INDEX_NAME = 'idx_email';

SET @Query = IF(@IndexExists = 0, 'CREATE UNIQUE INDEX idx_email ON tenants (email)', NULL);
PREPARE stmt FROM @Query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Check if the phone index exists and create it if it doesn't
SELECT COUNT(1) INTO @IndexExists
FROM INFORMATION_SCHEMA.STATISTICS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'tenants' 
  AND INDEX_NAME = 'idx_phone';

SET @Query = IF(@IndexExists = 0, 'CREATE UNIQUE INDEX idx_phone ON tenants (phone)', NULL);
PREPARE stmt FROM @Query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
