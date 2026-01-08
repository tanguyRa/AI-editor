# repository

```tree
repository/
├── README.md
├── db.go
│   ├── type DBTX interface{}
│   ├── type Queries {db: DBTX}
│   ├── func New(db DBTX) *Queries
│   └── func (*Queries) WithTx(tx pgx.Tx) *Queries
├── models.go
│   ├── type Account {ID: uuid.UUID, UserId: uuid.UUID, AccountId: string, ProviderId: string, AccessToken: *string, RefreshToken: *string, AccessTokenExpiresAt: *time.Time, RefreshTokenExpiresAt: *time.Time, Scope: *string, IdToken: *string, Password: *string, CreatedAt: time.Time, UpdatedAt: time.Time}
│   ├── type Session {ID: uuid.UUID, UserId: uuid.UUID, Token: string, ExpiresAt: time.Time, IpAddress: *string, UserAgent: *string, CreatedAt: time.Time, UpdatedAt: time.Time}
│   ├── type User {ID: uuid.UUID, Name: string, Email: string, EmailVerified: bool, Image: *string, CreatedAt: time.Time, UpdatedAt: time.Time}
│   └── type Verification {ID: uuid.UUID, Identifier: string, Value: string, ExpiresAt: time.Time, CreatedAt: time.Time, UpdatedAt: time.Time}
├── sessions.sql.go
│   ├── type CreateSessionParams {ID: uuid.UUID, Token: string, UserId: uuid.UUID, ExpiresAt: time.Time, IpAddress: *string, UserAgent: *string}
│   ├── type UpdateSessionParams {ID: uuid.UUID, Token: string, UserId: uuid.UUID, ExpiresAt: time.Time, IpAddress: *string, UserAgent: *string}
│   ├── func (*Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
│   ├── func (*Queries) DeleteSession(ctx context.Context, id uuid.UUID) (Session, error)
│   ├── func (*Queries) DeleteUserSessions(ctx context.Context, userid uuid.UUID) error
│   ├── func (*Queries) GetSession(ctx context.Context, id uuid.UUID) (Session, error)
│   └── func (*Queries) UpdateSession(ctx context.Context, arg UpdateSessionParams) (Session, error)
└── users.sql.go
    ├── type CreateUserParams {ID: uuid.UUID, Name: string, Email: string, EmailVerified: bool, Image: *string}
    ├── type UpdateUserParams {ID: uuid.UUID, Name: string, Email: string, EmailVerified: bool, Image: *string}
    ├── func (*Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
    ├── func (*Queries) DeleteUser(ctx context.Context, id uuid.UUID) (User, error)
    ├── func (*Queries) GetUserByEmail(ctx context.Context, email string) (User, error)
    ├── func (*Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
    ├── func (*Queries) ListUsers(ctx context.Context) ([]User, error)
    └── func (*Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
```
