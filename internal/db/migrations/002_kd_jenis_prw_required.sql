UPDATE lis_mapping_tests m
INNER JOIN template_laboratorium t ON t.id_template = m.id_template
SET m.kd_jenis_prw = t.kd_jenis_prw
WHERE TRIM(m.kd_jenis_prw) = '';

ALTER TABLE lis_mapping_tests
  MODIFY COLUMN kd_jenis_prw VARCHAR(15) NOT NULL COMMENT 'template_laboratorium.kd_jenis_prw';
