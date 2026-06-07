CREATE TABLE IF NOT EXISTS doctor_room (
    room_id           VARCHAR(50) PRIMARY KEY,
    room_name         VARCHAR(100) NOT NULL,
    last_assigned_at  DATETIME DEFAULT NULL,
    last_completed_at DATETIME DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
