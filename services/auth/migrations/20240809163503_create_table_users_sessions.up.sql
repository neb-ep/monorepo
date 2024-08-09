CREATE TABLE users_sessions (
    users_sessions_id serial,
    user_id integer NOT NULL,
    token character varying(32) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    is_active boolean NOT NULL,
    used_at timestamp with time zone,

    CONSTRAINT users_sessions_pkey PRIMARY KEY (users_sessions_id),
    CONSTRAINT users_sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(user_id)
);

COMMENT ON TABLE users_sessions IS 'Сессии пользователей';
COMMENT ON COLUMN users_sessions.users_sessions_id IS 'Ид сессии';
COMMENT ON COLUMN users_sessions.user_id IS 'Ид пользователя';
COMMENT ON COLUMN users_sessions.token IS 'Токен сессии';
COMMENT ON COLUMN users_sessions.created_at IS 'Дата и время создания сессии';
COMMENT ON COLUMN users_sessions.is_active IS 'Флаг что сессия активна';
COMMENT ON COLUMN users_sessions.used_at IS 'Дата и время использование токена сессии';
