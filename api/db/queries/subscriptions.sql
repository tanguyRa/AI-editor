-- name: CreateSubscription :one
INSERT INTO
    "subscription" (
        "userId",
        "polarSubscriptionId",
        tier,
        "scheduledTier",
        status,
        "currentPeriodEnd"
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: CreateSubscriptionWithId :one
INSERT INTO
    "subscription" (
        id,
        "userId",
        "polarSubscriptionId",
        tier,
        "scheduledTier",
        status,
        "currentPeriodEnd"
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    *;

-- name: GetSubscriptionByID :one
SELECT * FROM "subscription" WHERE id = $1;

-- name: GetSubscriptionByUserID :one
SELECT * FROM "subscription" WHERE "userId" = $1;

-- name: GetSubscriptionByPolarID :one
SELECT * FROM "subscription" WHERE "polarSubscriptionId" = $1;

-- name: UpdateSubscription :one
UPDATE "subscription"
SET
    "polarSubscriptionId" = $2,
    tier = $3,
    "scheduledTier" = $4,
    status = $5,
    "currentPeriodEnd" = $6,
    "updatedAt" = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateSubscriptionByUserID :one
UPDATE "subscription"
SET
    "polarSubscriptionId" = $2,
    tier = $3,
    "scheduledTier" = $4,
    status = $5,
    "currentPeriodEnd" = $6,
    "updatedAt" = CURRENT_TIMESTAMP
WHERE
    "userId" = $1
RETURNING
    *;

-- name: DeleteSubscription :one
DELETE FROM "subscription" WHERE id = $1 RETURNING *;

-- name: DeleteSubscriptionByUserID :one
DELETE FROM "subscription" WHERE "userId" = $1 RETURNING *;