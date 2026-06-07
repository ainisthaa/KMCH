CREATE TABLE IF NOT EXISTS patient_check (
    check_id            INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    line_id             VARCHAR(100) NOT NULL,
    event_id            INT UNSIGNED NOT NULL,
    psyeval_form        TINYINT(1)   NOT NULL DEFAULT 0,
    is_sv               TINYINT(1)   NOT NULL DEFAULT 0,
    is_paid             TINYINT(1)   NOT NULL DEFAULT 0,
    needs_transfer      TINYINT(1)   NOT NULL DEFAULT 0,
    transfer_completed  TINYINT(1)   NOT NULL DEFAULT 0,
    needs_psychologist  TINYINT(1)   NOT NULL DEFAULT 0,
    psychologist_done   TINYINT(1)   NOT NULL DEFAULT 0,
    route_type          VARCHAR(1)   DEFAULT NULL,
    created_at          DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY idx_line_event (line_id, event_id),
    FOREIGN KEY (line_id)   REFERENCES patient_info(line_id),
    FOREIGN KEY (event_id)  REFERENCES event_info(event_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
