CREATE TABLE IF NOT EXISTS lis_tests (
  id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  lis_test_id   VARCHAR(50)     NOT NULL COMMENT 'LIS examinations.testId — unique identity',
  local_code    VARCHAR(50)     NOT NULL DEFAULT '' COMMENT 'LIS localCode; not unique',
  test_name     VARCHAR(200)    NOT NULL DEFAULT '' COMMENT 'LIS testName / display label',
  status        ENUM('aktif','nonaktif') NOT NULL DEFAULT 'aktif',
  created_at    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_lis_tests_lis_test_id (lis_test_id),
  KEY idx_lis_tests_local_code (local_code),
  KEY idx_lis_tests_name (test_name),
  KEY idx_lis_tests_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS lis_mapping_tests (
  id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  lis_tests_pk    BIGINT UNSIGNED NOT NULL COMMENT 'FK lis_tests.id (surrogate PK)',
  id_template     INT(11)         NOT NULL COMMENT 'template_laboratorium.id_template',
  kd_jenis_prw    VARCHAR(15)     NOT NULL COMMENT 'template_laboratorium.kd_jenis_prw',
  status          ENUM('aktif','nonaktif') NOT NULL DEFAULT 'aktif',
  created_by      VARCHAR(50)     NOT NULL DEFAULT '',
  created_at      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_mapping (lis_tests_pk, id_template, kd_jenis_prw),
  KEY idx_mapping_template (id_template),
  KEY idx_mapping_panel (kd_jenis_prw),
  KEY idx_mapping_catalog (lis_tests_pk),
  CONSTRAINT fk_mapping_lis_tests
    FOREIGN KEY (lis_tests_pk) REFERENCES lis_tests (id)
    ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
