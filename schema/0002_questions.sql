-- +goose Up
CREATE TABLE IF NOT EXISTS security_questions (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    question text,
    created_at timestamptz DEFAULT NOW()
);

-- https://ux.stackexchange.com/questions/44151/security-question-what-questions-do-you-ask
INSERT INTO security_questions (question) VALUES
('What city were you born in?'),
('What was the make and model of your first car?'),
('What is your high school mascot?'),
('What was the last name of your third grade teacher?'),
('In what city or town was your first job?'),
('What is the first name of the person that you first kissed?');

CREATE TABLE IF NOT EXISTS users_security_questions (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid NOT NULL, 
    question_id uuid NOT NULL,
    answer_hash text NOT NULL,
    salt text,
    created_at timestamptz
);

-- +goose Down

DROP TABLE IF EXISTS security_questions;
DROP TABLE IF EXISTS users_security_questions;
