-- name: CreateProject :one
INSERT INTO
    "project" (
        "userId",
        name,
        slug,
        description
    )
VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: CreateProjectWithId :one
INSERT INTO
    "project" (
        id,
        "userId",
        name,
        slug,
        description
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: GetProjectByID :one
SELECT * FROM "project" WHERE id = $1;

-- name: GetProjectByUserIDAndSlug :one
SELECT * FROM "project" WHERE "userId" = $1 AND slug = $2;

-- name: ListProjectsByUserID :many
SELECT *
FROM "project"
WHERE
    "userId" = $1
ORDER BY "createdAt" DESC;

-- name: UpdateProject :one
UPDATE "project"
SET
    name = $2,
    slug = $3,
    description = $4,
    "updatedAt" = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteProject :one
DELETE FROM "project" WHERE id = $1 RETURNING *;

-- name: DeleteProjectsByUserID :exec
DELETE FROM "project" WHERE "userId" = $1;