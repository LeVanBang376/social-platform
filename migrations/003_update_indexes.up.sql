CREATE INDEX idx_comments_post_parent_created
    ON comments(post_id, parent_comment_id, created_at DESC);

DROP INDEX IF EXISTS idx_follows_follower;

DROP INDEX IF EXISTS idx_post_likes_post;

DROP INDEX IF EXISTS idx_comment_likes_comment;