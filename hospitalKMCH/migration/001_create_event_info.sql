CREATE TABLE IF NOT EXISTS event_info (
    event_id       INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    event_name     VARCHAR(255) NOT NULL,
    event_date_from DATETIME    NOT NULL,
    event_date_to   DATETIME    NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO event_info (event_name, event_date_from, event_date_to)
VALUES ('Annual Health Check 2026', '2026-06-01 08:00:00', '2026-06-30 17:00:00');
