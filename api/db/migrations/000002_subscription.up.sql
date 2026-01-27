-- Subscription table
CREATE TABLE "subscription" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    "userId" UUID UNIQUE NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    "polarSubscriptionId" VARCHAR(255),
    tier VARCHAR(50) NOT NULL DEFAULT 'free',
    "scheduledTier" VARCHAR(50),
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    "currentPeriodEnd" TIMESTAMPTZ,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_subscription_user_id ON "subscription" ("userId");