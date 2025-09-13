-- +goose Up
CREATE SCHEMA IF NOT EXISTS auth;
ALTER TABLE public.users SET SCHEMA auth;
ALTER TABLE public.roles SET SCHEMA auth;
ALTER TABLE public.users_roles SET SCHEMA auth;
ALTER TABLE public.users_security_questions SET SCHEMA auth;
ALTER TABLE public.security_questions SET SCHEMA auth;


-- +goose Down
ALTER TABLE auth.security_questions SET SCHEMA public;
ALTER TABLE auth.users_security_questions SET SCHEMA public;
ALTER TABLE auth.users_roles SET SCHEMA public;
ALTER TABLE auth.roles SET SCHEMA public;
ALTER TABLE auth.users SET SCHEMA public;
DROP SCHEMA IF EXISTS auth;
