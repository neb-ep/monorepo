CREATE TABLE users (
    user_id serial,
    username character varying(255) NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    email character varying(320) NOT NULL,
    passwhash character varying(255) NOT NULL,
    is_active boolean NOT NULL,
    created_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone,

    CONSTRAINT users_pkey PRIMARY KEY (user_id),

    CONSTRAINT users_username_min_length_check CHECK ((char_length(username) > 6)),
    CONSTRAINT users_email_min_length_check CHECK ((char_length(email) > 6))
);

CREATE UNIQUE INDEX users_username_key ON users (lower(username::text));
CREATE UNIQUE INDEX users_email_key ON users (lower(email::text));

COMMENT ON TABLE users IS 'Пользователи';
COMMENT ON COLUMN users.user_id IS 'ID пользователя';
COMMENT ON COLUMN users.username IS 'Никнейм';
COMMENT ON COLUMN users.first_name IS 'Имя';
COMMENT ON COLUMN users.last_name IS 'Фамилия';
COMMENT ON COLUMN users.email IS 'Email';
COMMENT ON COLUMN users.passwhash IS 'Хэш пароля';
COMMENT ON COLUMN users.is_active IS 'Флаг активности уч. ползьзователя';
COMMENT ON COLUMN users.created_at IS 'Дата создания';
COMMENT ON COLUMN users.deleted_at IS 'Дата удаления';


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
