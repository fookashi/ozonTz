CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username TEXT NOT NULL,
    roles TEXT[] NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    is_commentable BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    post_id UUID NOT NULL,
    parent_id UUID,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_roles ON users USING gin(roles);

CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments USING hash(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments USING hash(parent_id);
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments USING btree(created_at);

CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts USING hash(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts USING btree(created_at);
