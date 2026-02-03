-- Modify crossword_questions
ALTER TABLE crossword_questions ADD COLUMN level_id BIGINT UNSIGNED NOT NULL AFTER id;

-- Migrate data for crossword_questions
UPDATE crossword_questions cq
JOIN crosswords c ON cq.crossword_id = c.id
SET cq.level_id = c.level_id;

-- Add FK for crossword_questions
ALTER TABLE crossword_questions ADD CONSTRAINT fk_crossword_questions_level FOREIGN KEY (level_id) REFERENCES levels(id) ON DELETE CASCADE;

-- Drop crossword_id from crossword_questions
ALTER TABLE crossword_questions DROP FOREIGN KEY crossword_questions_ibfk_1; -- Assuming standard naming or based on previous file
ALTER TABLE crossword_questions DROP COLUMN crossword_id;

-- Modify user_crossword_scores -> user_level_scores
RENAME TABLE user_crossword_scores TO user_level_scores;
ALTER TABLE user_level_scores ADD COLUMN level_id BIGINT UNSIGNED NOT NULL AFTER user_id;

-- Migrate data for scores (if any)
UPDATE user_level_scores uls
JOIN crosswords c ON uls.crossword_id = c.id
SET uls.level_id = c.level_id;

-- Add FK for user_level_scores
ALTER TABLE user_level_scores ADD CONSTRAINT fk_user_level_scores_level FOREIGN KEY (level_id) REFERENCES levels(id) ON DELETE CASCADE;

-- Drop crossword_id and old unique constraint
ALTER TABLE user_level_scores DROP FOREIGN KEY user_crossword_scores_ibfk_2; -- Assuming standard naming
ALTER TABLE user_level_scores DROP INDEX user_id; -- Unique constraint usually creates an index
ALTER TABLE user_level_scores DROP COLUMN crossword_id;
ALTER TABLE user_level_scores ADD UNIQUE(user_id, level_id);

-- Drop crosswords table
DROP TABLE crosswords;
