-- ============================================
-- Users
-- ============================================

CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(254) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    date_of_birth DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- ============================================
-- Follows
-- ============================================

CREATE TABLE follows (
    follower_id UUID NOT NULL,
    following_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (follower_id, following_id),

    CONSTRAINT fk_follows_follower
        FOREIGN KEY (follower_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_follows_following
        FOREIGN KEY (following_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE,

    CONSTRAINT chk_follows_not_self
        CHECK (follower_id <> following_id)
);

-- Get users that a user is following
CREATE INDEX idx_follows_follower
    ON follows(follower_id);

-- Get followers of a user
CREATE INDEX idx_follows_following
    ON follows(following_id);


-- ============================================
-- Posts
-- ============================================

CREATE TABLE posts (
    post_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id UUID NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_posts_user
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE
);

-- User's posts / profile feed
CREATE INDEX idx_posts_user_created
    ON posts(user_id, created_at DESC);

-- Global feed / latest posts
CREATE INDEX idx_posts_created
    ON posts(created_at DESC);


-- ============================================
-- Post Likes
-- ============================================

CREATE TABLE post_likes (
    post_id BIGINT NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (post_id, user_id),

    CONSTRAINT fk_post_likes_post
        FOREIGN KEY (post_id)
        REFERENCES posts(post_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_post_likes_user
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE
);

-- Find users who liked a post
CREATE INDEX idx_post_likes_post
    ON post_likes(post_id);

-- Find posts liked by a user
CREATE INDEX idx_post_likes_user
    ON post_likes(user_id);


-- ============================================
-- Comments
-- ============================================

CREATE TABLE comments (
    comment_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    post_id BIGINT NOT NULL,
    user_id UUID NOT NULL,
    parent_comment_id BIGINT,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_comments_post
        FOREIGN KEY (post_id)
        REFERENCES posts(post_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_comments_user
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_comments_parent
        FOREIGN KEY (parent_comment_id)
        REFERENCES comments(comment_id)
        ON DELETE CASCADE
);

-- Get comments of a post
CREATE INDEX idx_comments_post_created
    ON comments(post_id, created_at DESC);

-- Get comments made by a user
CREATE INDEX idx_comments_user
    ON comments(user_id);


-- ============================================
-- Comment Likes
-- ============================================

CREATE TABLE comment_likes (
    comment_id BIGINT NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (comment_id, user_id),

    CONSTRAINT fk_comment_likes_comment
        FOREIGN KEY (comment_id)
        REFERENCES comments(comment_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_comment_likes_user
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE
);

-- Find users who liked a comment
CREATE INDEX idx_comment_likes_comment
    ON comment_likes(comment_id);

-- Find comments liked by a user
CREATE INDEX idx_comment_likes_user
    ON comment_likes(user_id);