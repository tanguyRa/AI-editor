DROP INDEX IF EXISTS idx_sessions_user_id;

DROP INDEX IF EXISTS idx_sessions_token;

DROP INDEX IF EXISTS idx_accounts_user_id;

DROP INDEX IF EXISTS idx_verifications_identifier;

DROP TABLE IF EXISTS "verification";

DROP TABLE IF EXISTS "account";

DROP TABLE IF EXISTS "session";

DROP TABLE IF EXISTS "user";