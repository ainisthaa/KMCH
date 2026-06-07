INSERT INTO doctor_room (room_id, room_name) VALUES
    ('room-001', 'Doctor Room 1'),
    ('room-002', 'Doctor Room 2'),
    ('room-003', 'Doctor Room 3'),
    ('room-004', 'Doctor Room 4'),
    ('room-005', 'Doctor Room 5')
ON DUPLICATE KEY UPDATE room_name = VALUES(room_name);
