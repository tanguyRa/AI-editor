-- Project table
CREATE TABLE "project" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    "userId" UUID NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE ("userId", slug)
);

-- Indexes
CREATE INDEX idx_project_user_id ON "project" ("userId");