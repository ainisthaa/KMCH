CREATE TABLE IF NOT EXISTS patient_queue (
    queue_id    INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    line_id     VARCHAR(100) NOT NULL,
    event_id    INT UNSIGNED NOT NULL,
    queue       VARCHAR(30)  NOT NULL COMMENT 'Internal FIFO timestamp — never expose to client',
    room_id     VARCHAR(50)  DEFAULT NULL,
    status      VARCHAR(20)  NOT NULL DEFAULT 'waiting',
    station     VARCHAR(20)  NOT NULL DEFAULT 'queue',
    q_starttime DATETIME     DEFAULT NULL,
    q_endtime   DATETIME     DEFAULT NULL,
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status_queue (status, queue),
    INDEX idx_line_event   (line_id, event_id),
    FOREIGN KEY (line_id)  REFERENCES patient_info(line_id),
    FOREIGN KEY (event_id) REFERENCES event_info(event_id),
    FOREIGN KEY (room_id)  REFERENCES doctor_room(room_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
