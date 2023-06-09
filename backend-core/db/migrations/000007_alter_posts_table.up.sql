ALTER TABLE IF EXISTS posts
DROP COLUMN board_id,
DROP COLUMN pos_x,
DROP COLUMN pos_y,
DROP COLUMN z_index;

ALTER TABLE IF EXISTS posts
ADD COLUMN post_order FLOAT;

ALTER TABLE IF EXISTS posts
ADD COLUMN post_group_id UUID NOT NULL,
ADD CONSTRAINT fk_post_group
    FOREIGN KEY (post_group_id)
    REFERENCES post_groups (id)
    ON DELETE CASCADE;
