CREATE TABLE languages (
    id SERIAL PRIMARY KEY,
    language VARCHAR(100) NOT NULL
);

INSERT INTO languages (language) 
VALUES 
('tkm'),
('eng'),
('rus');