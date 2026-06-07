CREATE TABLE IF NOT EXISTS patient_info (
    line_id       VARCHAR(100) PRIMARY KEY,
    first_name    VARCHAR(100),
    last_name     VARCHAR(100),
    tel_no        VARCHAR(20),
    national_id   VARCHAR(20),
    passport_id   VARCHAR(30),
    register_date DATETIME,
    student_id    VARCHAR(30),
    employee_id   VARCHAR(30)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
