-- name: CreateUser :exec
INSERT into users
(id, name, email, password, is_guest, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetUser :one
SELECT * FROM users
WHERE users.id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE users.email = $1;

-- name: UpdateUserVerification :exec
UPDATE users SET
(id, is_verified) =
($1, $2) WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE users.id = $1;

-- name: CreateBoard :exec
INSERT INTO boards 
(id, name, description, user_id, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetBoard :one
SELECT * FROM boards
WHERE boards.id = $1;

-- name: GetBoardAndUsers :many
SELECT sqlc.embed(boards), sqlc.embed(users), sqlc.embed(board_memberships) FROM boards
INNER JOIN board_memberships on board_memberships.board_id = boards.id
INNER JOIN users on board_memberships.user_id = users.id
WHERE boards.id = $1
ORDER BY boards.created_at DESC;

-- name: ListOwnedBoards :many
SELECT * FROM boards
WHERE boards.user_id = $1
ORDER BY boards.created_at DESC;

-- name: ListOwnedBoardAndUsers :many
SELECT sqlc.embed(boards), sqlc.embed(users), sqlc.embed(board_memberships) FROM boards
INNER JOIN board_memberships on board_memberships.board_id = boards.id
INNER JOIN users on board_memberships.user_id = users.id
WHERE boards.user_id = $1
ORDER BY boards.created_at DESC;

-- name: ListSharedBoardAndUsers :many
SELECT sqlc.embed(boards), sqlc.embed(users), sqlc.embed(board_memberships) FROM boards
INNER JOIN board_memberships on board_memberships.board_id = boards.id
INNER JOIN users on board_memberships.user_id = users.id
WHERE board_memberships.user_id = $1
AND board_memberships.role = 'MEMBER'
ORDER BY board_memberships.created_at DESC;

-- name: CreateMembership :exec
INSERT INTO board_memberships 
(id, user_id, board_id, role, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6);

-- name: DeleteBoard :exec
DELETE from boards WHERE id = $1;

-- name: CreatePost :exec
INSERT INTO posts
(id, user_id, content, color, height, created_at, updated_at, post_order, post_group_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: CreatePostGroup :exec
INSERT INTO post_groups
(id, board_id, title, pos_x, pos_y, z_index, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetPost :one
SELECT * FROM posts
WHERE posts.id = $1;

-- name: ListPostGroups :many
SELECT sqlc.embed(post_groups), sqlc.embed(posts) FROM post_groups
INNER JOIN posts on posts.post_group_id = post_groups.id
WHERE post_groups.board_id = $1
ORDER BY posts.post_order ASC;

-- name: UpdatePost :exec
UPDATE posts SET
(id, user_id, content, color, height, created_at, updated_at, post_order, post_group_id) =
($1, $2, $3, $4, $5, $6, $7, $8, $9) WHERE id = $1;

-- name: DeletePost :exec
DELETE from posts WHERE id = $1;

-- name: GetPostGroup :one
SELECT * FROM post_groups
WHERE post_groups.id = $1;

-- name: UpdatePostGroup :exec
UPDATE post_groups SET
(id, board_id, title, pos_x, pos_y, z_index, created_at, updated_at) =
($1, $2, $3, $4, $5, $6, $7, $8) WHERE id = $1;

-- name: DeletePostGroup :exec
DELETE from post_groups WHERE id = $1;

-- name: ListUsersByFuzzyEmail :many
SELECT * FROM users
ORDER BY levenshtein(users.email, $1) LIMIT 10;

-- name: CreateInvite :exec
INSERT INTO board_invites
(id, board_id, sender_id, receiver_id, status, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetInvite :one
SELECT sqlc.embed(board_invites), sqlc.embed(s), sqlc.embed(r) FROM board_invites
JOIN users s on s.id = board_invites.sender_id
JOIN users r on r.id = board_invites.receiver_id
WHERE board_invites.id = $1;

-- name: UpdateInvite :exec
UPDATE board_invites SET
(id, board_id, sender_id, receiver_id, status, created_at, updated_at) =
($1, $2, $3, $4, $5, $6, $7) WHERE id = $1;

-- name: ListInvitesByBoard :many
SELECT sqlc.embed(board_invites), sqlc.embed(users) FROM board_invites
INNER JOIN users on users.id = board_invites.receiver_id
WHERE board_invites.board_id = sqlc.arg('board_id') AND
(status = sqlc.narg('status') OR sqlc.narg('status') IS NULL)
ORDER BY board_invites.updated_at DESC;

-- name: ListInvitesByReceiver :many
SELECT sqlc.embed(board_invites), sqlc.embed(users), sqlc.embed(boards) FROM board_invites
INNER JOIN boards on boards.id = board_invites.board_id
INNER JOIN users on users.id = board_invites.sender_id 
WHERE board_invites.receiver_id = sqlc.arg('receiver_id') AND
(status = sqlc.narg('status') OR sqlc.narg('status') IS NULL)
ORDER BY board_invites.updated_at DESC;

-- name: CreateEmailVerification :exec
INSERT INTO email_verifications
(id, code, user_id, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateEmailVerification :exec
UPDATE email_verifications SET
(user_id, is_verified) =
($1, $2) WHERE user_id = $1 AND is_verified IS NULL;

-- name: GetEmailVerification :one
SELECT * FROM email_verifications WHERE user_id = $1 AND is_verified IS NULL
ORDER BY created_at DESC LIMIT 1;
