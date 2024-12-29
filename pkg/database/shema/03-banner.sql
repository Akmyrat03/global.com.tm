CREATE TABLE banner (
    id SERIAL PRIMARY KEY,
    image_path VARCHAR(255) NOT NULL
);

CREATE TABLE banner_translate (
    id SERIAL PRIMARY KEY,
    banner_id INT REFERENCES banner (id) ON DELETE CASCADE,
    lang_id INT REFERENCES languages (id) ON DELETE CASCADE,
    title VARCHAR(255),
    description VARCHAR(255)
);