ALTER TABLE crossword_questions
ADD COLUMN questions_id BIGINT UNSIGNED DEFAULT NULL,
ADD CONSTRAINT fk_crossword_questions_questions
FOREIGN KEY (questions_id) REFERENCES questions(id)
ON DELETE SET NULL
ON UPDATE CASCADE;
