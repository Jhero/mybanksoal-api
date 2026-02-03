DROP TABLE IF EXISTS user_crossword_scores;
DROP TABLE IF EXISTS crossword_questions;
DROP TABLE IF EXISTS crossword_levels;

CREATE TABLE crossword_levels (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE crossword_questions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    level_id BIGINT UNSIGNED NOT NULL,
    number INT NOT NULL,
    clue TEXT NOT NULL,
    answer VARCHAR(255) NOT NULL,
    is_across BOOLEAN NOT NULL,
    `row` INT NOT NULL,
    col INT NOT NULL,
    FOREIGN KEY (level_id) REFERENCES crossword_levels(id) ON DELETE CASCADE
);

CREATE TABLE user_crossword_scores (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    level_id BIGINT UNSIGNED NOT NULL,
    score INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE(user_id, level_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (level_id) REFERENCES crossword_levels(id) ON DELETE CASCADE
);

-- Seed Level 1
INSERT INTO crossword_levels (id, title) VALUES (1, 'Level 1 - Dasar');

INSERT INTO crossword_questions (level_id, number, clue, answer, is_across, `row`, col) VALUES 
(1, 1, 'Ibukota Indonesia', 'JAKARTA', true, 0, 0),
(1, 2, 'Bahasa pemrograman Flutter', 'DART', false, 0, 2);

-- Seed Level 2
INSERT INTO crossword_levels (id, title) VALUES (2, 'Level 2 - Geografi');

INSERT INTO crossword_questions (level_id, number, clue, answer, is_across, `row`, col) VALUES 
(2, 1, 'Gunung tertinggi di dunia', 'EVEREST', true, 2, 0);
