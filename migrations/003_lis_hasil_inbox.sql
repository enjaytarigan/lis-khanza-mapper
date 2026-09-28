CREATE TABLE IF NOT EXISTS lis_hasil_inbox (
  id               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  no_order         VARCHAR(50)     NOT NULL COMMENT 'SIMRS permintaan_lab.noorder / MedQLab noOrder',
  no_laboratorium  VARCHAR(50)     NOT NULL DEFAULT '',
  status           ENUM('accepted','posted','failed') NOT NULL DEFAULT 'accepted',
  mapped_count     INT             NOT NULL DEFAULT 0,
  unmapped_count   INT             NOT NULL DEFAULT 0,
  detail_written   INT             NOT NULL DEFAULT 0,
  error_message    TEXT            NULL,
  payload_json     LONGTEXT        NULL,
  created_at       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_lis_hasil_inbox_order (no_order),
  KEY idx_lis_hasil_inbox_status (status),
  KEY idx_lis_hasil_inbox_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
