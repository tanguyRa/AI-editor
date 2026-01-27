# repository

```tree
repository/
├── README.md
├── accounts.sql.go
│   ├── type CreateAccountParams {UserId: uuid.UUID, AccountId: string, ProviderId: string, Password: *string}
│   ├── type CreateAccountWithIdParams {ID: uuid.UUID, UserId: uuid.UUID, AccountId: string, ProviderId: string, Password: *string}
│   ├── type GetAccountByUserIdAndProviderParams {UserId: uuid.UUID, ProviderId: string}
│   ├── func (*Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
│   ├── func (*Queries) CreateAccountWithId(ctx context.Context, arg CreateAccountWithIdParams) (Account, error)
│   ├── func (*Queries) DeleteAccount(ctx context.Context, id uuid.UUID) (Account, error)
│   ├── func (*Queries) GetAccountById(ctx context.Context, id uuid.UUID) (Account, error)
│   └── func (*Queries) GetAccountByUserIdAndProvider(ctx context.Context, arg GetAccountByUserIdAndProviderParams) (Account, error)
├── db.go
│   ├── type DBTX interface{}
│   ├── type Queries {db: DBTX}
│   ├── func New(db DBTX) *Queries
│   └── func (*Queries) WithTx(tx pgx.Tx) *Queries
├── events.sql.go
│   ├── type CreateEventParams {UserId: uuid.UUID, Type: string, Data: []byte}
│   ├── type CreateEventWithIdParams {ID: uuid.UUID, UserId: uuid.UUID, Type: string, Data: []byte}
│   ├── type GetEventByUserIDAndTypeParams {UserId: uuid.UUID, Type: string}
│   ├── func (*Queries) CreateEvent(ctx context.Context, arg CreateEventParams) (Event, error)
│   ├── func (*Queries) CreateEventWithId(ctx context.Context, arg CreateEventWithIdParams) (Event, error)
│   ├── func (*Queries) GetEventByID(ctx context.Context, id uuid.UUID) (Event, error)
│   ├── func (*Queries) GetEventByUserIDAndType(ctx context.Context, arg GetEventByUserIDAndTypeParams) (Event, error)
│   └── func (*Queries) ListEventsByUserID(ctx context.Context, userid uuid.UUID) ([]Event, error)
├── jwks.sql.go
│   ├── type GetJwksSetsRow {ID: uuid.UUID, PublicKey: string, CreatedAt: time.Time, ExpiresAt: *time.Time}
│   └── func (*Queries) GetJwksSets(ctx context.Context) ([]GetJwksSetsRow, error)
├── models.go
│   ├── type Account {ID: uuid.UUID, UserId: uuid.UUID, AccountId: string, ProviderId: string, AccessToken: *string, RefreshToken: *string, AccessTokenExpiresAt: *time.Time, RefreshTokenExpiresAt: *time.Time, Scope: *string, IdToken: *string, Password: *string, CreatedAt: time.Time, UpdatedAt: time.Time}
│   ├── type Event {ID: uuid.UUID, UserId: uuid.UUID, Data: []byte, Type: string, CreatedAt: time.Time, UpdatedAt: time.Time}
│   ├── type Jwk {ID: uuid.UUID, PublicKey: string, PrivateKey: string, CreatedAt: time.Time, ExpiresAt: *time.Time}
│   ├── type Project {ID: uuid.UUID, UserId: uuid.UUID, Name: string, Slug: string, Description: *string, CreatedAt: time.Time, UpdatedAt: time.Time}
│   ├── type Session {ID: uuid.UUID, UserId: uuid.UUID, Token: string, ExpiresAt: time.Time, IpAddress: *string, UserAgent: *string, CreatedAt: time.Time, UpdatedAt: time.Time}
│   ├── type Subscription {ID: uuid.UUID, UserId: uuid.UUID, PolarSubscriptionId: *string, Tier: string, ScheduledTier: *string, Status: string, CurrentPeriodEnd: *time.Time, CreatedAt: time.Time, UpdatedAt: time.Time}
│   ├── type User {ID: uuid.UUID, Name: string, Email: string, EmailVerified: bool, Image: *string, CreatedAt: time.Time, UpdatedAt: time.Time}
│   └── type Verification {ID: uuid.UUID, Identifier: string, Value: string, ExpiresAt: time.Time, CreatedAt: time.Time, UpdatedAt: time.Time}
├── projects.sql.go
│   ├── type CreateProjectParams {UserId: uuid.UUID, Name: string, Slug: string, Description: *string}
│   ├── type CreateProjectWithIdParams {ID: uuid.UUID, UserId: uuid.UUID, Name: string, Slug: string, Description: *string}
│   ├── type GetProjectByUserIDAndSlugParams {UserId: uuid.UUID, Slug: string}
│   ├── type UpdateProjectParams {ID: uuid.UUID, Name: string, Slug: string, Description: *string}
│   ├── func (*Queries) CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error)
│   ├── func (*Queries) CreateProjectWithId(ctx context.Context, arg CreateProjectWithIdParams) (Project, error)
│   ├── func (*Queries) DeleteProject(ctx context.Context, id uuid.UUID) (Project, error)
│   ├── func (*Queries) DeleteProjectsByUserID(ctx context.Context, userid uuid.UUID) error
│   ├── func (*Queries) GetProjectByID(ctx context.Context, id uuid.UUID) (Project, error)
│   ├── func (*Queries) GetProjectByUserIDAndSlug(ctx context.Context, arg GetProjectByUserIDAndSlugParams) (Project, error)
│   ├── func (*Queries) ListProjectsByUserID(ctx context.Context, userid uuid.UUID) ([]Project, error)
│   └── func (*Queries) UpdateProject(ctx context.Context, arg UpdateProjectParams) (Project, error)
├── sessions.sql.go
│   ├── type CreateSessionParams {Token: string, UserId: uuid.UUID, ExpiresAt: time.Time, IpAddress: *string, UserAgent: *string}
│   ├── type CreateSessionWithIdParams {ID: uuid.UUID, Token: string, UserId: uuid.UUID, ExpiresAt: time.Time, IpAddress: *string, UserAgent: *string}
│   ├── type UpdateSessionParams {ID: uuid.UUID, Token: string, UserId: uuid.UUID, ExpiresAt: time.Time, IpAddress: *string, UserAgent: *string}
│   ├── func (*Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
│   ├── func (*Queries) CreateSessionWithId(ctx context.Context, arg CreateSessionWithIdParams) (Session, error)
│   ├── func (*Queries) DeleteSession(ctx context.Context, id uuid.UUID) (Session, error)
│   ├── func (*Queries) DeleteUserSessions(ctx context.Context, userid uuid.UUID) error
│   ├── func (*Queries) GetSession(ctx context.Context, id uuid.UUID) (Session, error)
│   ├── func (*Queries) GetSessionByToken(ctx context.Context, token string) (Session, error)
│   └── func (*Queries) UpdateSession(ctx context.Context, arg UpdateSessionParams) (Session, error)
├── subscriptions.sql.go
│   ├── type CreateSubscriptionParams {UserId: uuid.UUID, PolarSubscriptionId: *string, Tier: string, ScheduledTier: *string, Status: string, CurrentPeriodEnd: *time.Time}
│   ├── type CreateSubscriptionWithIdParams {ID: uuid.UUID, UserId: uuid.UUID, PolarSubscriptionId: *string, Tier: string, ScheduledTier: *string, Status: string, CurrentPeriodEnd: *time.Time}
│   ├── type UpdateSubscriptionParams {ID: uuid.UUID, PolarSubscriptionId: *string, Tier: string, ScheduledTier: *string, Status: string, CurrentPeriodEnd: *time.Time}
│   ├── type UpdateSubscriptionByUserIDParams {UserId: uuid.UUID, PolarSubscriptionId: *string, Tier: string, ScheduledTier: *string, Status: string, CurrentPeriodEnd: *time.Time}
│   ├── func (*Queries) CreateSubscription(ctx context.Context, arg CreateSubscriptionParams) (Subscription, error)
│   ├── func (*Queries) CreateSubscriptionWithId(ctx context.Context, arg CreateSubscriptionWithIdParams) (Subscription, error)
│   ├── func (*Queries) DeleteSubscription(ctx context.Context, id uuid.UUID) (Subscription, error)
│   ├── func (*Queries) DeleteSubscriptionByUserID(ctx context.Context, userid uuid.UUID) (Subscription, error)
│   ├── func (*Queries) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (Subscription, error)
│   ├── func (*Queries) GetSubscriptionByPolarID(ctx context.Context, polarsubscriptionid *string) (Subscription, error)
│   ├── func (*Queries) GetSubscriptionByUserID(ctx context.Context, userid uuid.UUID) (Subscription, error)
│   ├── func (*Queries) UpdateSubscription(ctx context.Context, arg UpdateSubscriptionParams) (Subscription, error)
│   └── func (*Queries) UpdateSubscriptionByUserID(ctx context.Context, arg UpdateSubscriptionByUserIDParams) (Subscription, error)
└── users.sql.go
    ├── type CreateUserParams {Name: string, Email: string, EmailVerified: bool, Image: *string}
    ├── type CreateUserWithIdParams {ID: uuid.UUID, Name: string, Email: string, EmailVerified: bool, Image: *string}
    ├── type UpdateUserParams {ID: uuid.UUID, Name: string, Email: string, EmailVerified: bool, Image: *string}
    ├── func (*Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
    ├── func (*Queries) CreateUserWithId(ctx context.Context, arg CreateUserWithIdParams) (User, error)
    ├── func (*Queries) DeleteUser(ctx context.Context, id uuid.UUID) (User, error)
    ├── func (*Queries) GetUserByEmail(ctx context.Context, email string) (User, error)
    ├── func (*Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
    ├── func (*Queries) ListUsers(ctx context.Context) ([]User, error)
    └── func (*Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
```
