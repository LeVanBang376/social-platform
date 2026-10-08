DROP INDEX IF EXISTS idx_comments_post_parent_created;

CREATE INDEX idx_follows_follower
    ON follows(follower_id);

CREATE INDEX idx_post_likes_post
    ON post_likes(post_id);

CREATE INDEX idx_comment_likes_comment
    ON comment_likes(comment_id);