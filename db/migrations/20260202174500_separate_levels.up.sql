DROP TABLE IF EXISTS user_crossword_scores;
DROP TABLE IF EXISTS crossword_questions;
DROP TABLE IF EXISTS crossword_levels;

CREATE TABLE levels (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE crosswords (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    level_id BIGINT UNSIGNED NOT NULL,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (level_id) REFERENCES levels(id) ON DELETE CASCADE
);

CREATE TABLE crossword_questions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    crossword_id BIGINT UNSIGNED NOT NULL,
    number INT NOT NULL,
    clue TEXT NOT NULL,
    answer VARCHAR(255) NOT NULL,
    is_across BOOLEAN NOT NULL,
    `row` INT NOT NULL,
    col INT NOT NULL,
    FOREIGN KEY (crossword_id) REFERENCES crosswords(id) ON DELETE CASCADE
);

CREATE TABLE user_crossword_scores (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    crossword_id BIGINT UNSIGNED NOT NULL,
    score INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE(user_id, crossword_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (crossword_id) REFERENCES crosswords(id) ON DELETE CASCADE
);

-- Seed Data
INSERT INTO levels (id, name) VALUES (1, 'Beginner');
INSERT INTO levels (id, name) VALUES (2, 'Intermediate');

INSERT INTO crosswords (id, level_id, title) VALUES (1, 1, 'Dasar');
INSERT INTO crosswords (id, level_id, title) VALUES (2, 2, 'Geografi');

INSERT INTO crossword_questions (crossword_id, number, clue, answer, is_across, `row`, col) VALUES 
(1, 1, 'Ibukota Indonesia', 'JAKARTA', true, 0, 0),
(1, 2, 'Bahasa pemrograman Flutter', 'DART', false, 0, 2),
(2, 1, 'Gunung tertinggi di dunia', 'EVEREST', true, 2, 0);
