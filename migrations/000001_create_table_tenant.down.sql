-- Drop the phone index if it exists
SELECT COUNT(1) INTO @IndexExists
FROM INFORMATION_SCHEMA.STATISTICS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'tenants' 
  AND INDEX_NAME = 'idx_phone';

SET @Query = IF(@IndexExists > 0, 'ALTER TABLE tenants DROP INDEX idx_phone', NULL);
PREPARE stmt FROM @Query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Drop the email index if it exists
SELECT COUNT(1) INTO @IndexExists
FROM INFORMATION_SCHEMA.STATISTICS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'tenants' 
  AND INDEX_NAME = 'idx_email';

SET @Query = IF(@IndexExists > 0, 'ALTER TABLE tenants DROP INDEX idx_email', NULL);
PREPARE stmt FROM @Query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Drop the tenants table if it exists
DROP TABLE IF EXISTS tenants;
