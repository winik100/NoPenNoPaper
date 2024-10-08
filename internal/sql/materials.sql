CREATE TABLE IF NOT EXISTS materials (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(50) NOT NULL,
    file_name VARCHAR(100) NOT NULL,
    uploaded_by INTEGER NOT NULL,
    CONSTRAINT fk_users_materials_id FOREIGN KEY (uploaded_by) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_filename_user UNIQUE (file_name, uploaded_by)
);