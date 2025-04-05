CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(255) NOT NULL UNIQUE,
                       password TEXT NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE groups (
                        id SERIAL PRIMARY KEY,
                        name VARCHAR(255),
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE group_members (
                               id SERIAL PRIMARY KEY,
                               group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
                               user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE blocked_users (
                               id SERIAL PRIMARY KEY,
                               blocker_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               blocked_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE messages (
                          id SERIAL PRIMARY KEY,
                          sender_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          group_id INTEGER REFERENCES groups(id) ON DELETE CASCADE,
                          content TEXT NOT NULL,
                          scheduled_at TIMESTAMP,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);