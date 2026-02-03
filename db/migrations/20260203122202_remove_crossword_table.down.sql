-- Revert changes (simplified, data loss on crosswords structure)
CREATE TABLE crosswords (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    level_id BIGINT UNSIGNED NOT NULL,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (level_id) REFERENCES levels(id) ON DELETE CASCADE
);

ALTER TABLE crossword_questions ADD COLUMN crossword_id BIGINT UNSIGNED NOT NULL AFTER id;
ALTER TABLE crossword_questions ADD CONSTRAINT fk_crossword_questions_crossword FOREIGN KEY (crossword_id) REFERENCES crosswords(id) ON DELETE CASCADE;
ALTER TABLE crossword_questions DROP FOREIGN KEY fk_crossword_questions_level;
ALTER TABLE crossword_questions DROP COLUMN level_id;

RENAME TABLE user_level_scores TO user_crossword_scores;
ALTER TABLE user_crossword_scores ADD COLUMN crossword_id BIGINT UNSIGNED NOT NULL AFTER user_id;
ALTER TABLE user_crossword_scores ADD CONSTRAINT fk_user_crossword_scores_crossword FOREIGN KEY (crossword_id) REFERENCES crosswords(id) ON DELETE CASCADE;
ALTER TABLE user_crossword_scores DROP FOREIGN KEY fk_user_level_scores_level;
ALTER TABLE user_crossword_scores DROP COLUMN level_id;
