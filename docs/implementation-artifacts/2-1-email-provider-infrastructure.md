# Story 2.1: Email Provider Infrastructure

Status: ready-for-dev

## Story

As a **system administrator**,
I want to **configure email delivery through SendGrid or Resend**,
So that **the application can send transactional emails reliably**.

## Acceptance Criteria

1. **Given** the application starts with `EMAIL_PROVIDER=sendgrid` and valid SendGrid API key **When** the configuration is loaded **Then** the SendGrid email provider is initialized **And** the provider is available for dependency injection into handlers

2. **Given** the application starts with `EMAIL_PROVIDER=resend` and valid Resend API key **When** the configuration is loaded **Then** the Resend email provider is initialized

3. **Given** the application starts with missing or invalid email configuration **When** the configuration is validated **Then** the application fails fast with a clear error message (NFR19) **And** the error indicates which configuration is missing

4. **Given** an email provider is configured **When** an email send is requested **Then** the provider makes an HTTP POST to the provider's API **And** the request includes proper authentication headers

## Tasks / Subtasks

- [ ] Task 1: Define EmailProvider interface (AC: 1, 2, 4)
  - [ ] 1.1: Create `internal/email/provider.go`
  - [ ] 1.2: Define `Email` struct with To, From, Subject, HTMLBody, TextBody fields
  - [ ] 1.3: Define `EmailProvider` interface with `Send(ctx, email) error` method
  - [ ] 1.4: Define `EmailConfig` struct for provider configuration

- [ ] Task 2: Implement SendGrid provider (AC: 1, 4)
  - [ ] 2.1: Create `internal/email/sendgrid.go`
  - [ ] 2.2: Implement SendGrid API client using stdlib net/http
  - [ ] 2.3: Use SendGrid v3 Mail Send API: `POST https://api.sendgrid.com/v3/mail/send`
  - [ ] 2.4: Set Authorization header with Bearer token
  - [ ] 2.5: Format request body per SendGrid API spec
  - [ ] 2.6: Handle API errors and return meaningful errors

- [ ] Task 3: Implement Resend provider (AC: 2, 4)
  - [ ] 3.1: Create `internal/email/resend.go`
  - [ ] 3.2: Implement Resend API client using stdlib net/http
  - [ ] 3.3: Use Resend API: `POST https://api.resend.com/emails`
  - [ ] 3.4: Set Authorization header with Bearer token
  - [ ] 3.5: Format request body per Resend API spec
  - [ ] 3.6: Handle API errors and return meaningful errors

- [ ] Task 4: Add email configuration to config.go (AC: 1, 2, 3)
  - [ ] 4.1: Add EMAIL_PROVIDER env var (sendgrid | resend)
  - [ ] 4.2: Add SENDGRID_API_KEY env var
  - [ ] 4.3: Add RESEND_API_KEY env var
  - [ ] 4.4: Add EMAIL_FROM env var (default sender address)
  - [ ] 4.5: Add FRONTEND_URL env var (for email links)

- [ ] Task 5: Implement config validation (AC: 3)
  - [ ] 5.1: Validate EMAIL_PROVIDER is set and valid
  - [ ] 5.2: Validate required API key exists for selected provider
  - [ ] 5.3: Validate EMAIL_FROM is set
  - [ ] 5.4: Fail fast on startup with clear error messages

- [ ] Task 6: Create email provider factory (AC: 1, 2)
  - [ ] 6.1: Create `NewEmailProvider(config) (EmailProvider, error)` function
  - [ ] 6.2: Return SendGrid or Resend provider based on config
  - [ ] 6.3: Wire into server startup in cmd/server/main.go

- [ ] Task 7: Write tests (AC: 1, 2, 3, 4)
  - [ ] 7.1: Test SendGrid provider formats request correctly
  - [ ] 7.2: Test Resend provider formats request correctly
  - [ ] 7.3: Test config validation fails on missing provider
  - [ ] 7.4: Test config validation fails on missing API key
  - [ ] 7.5: Test factory creates correct provider type

## Dev Notes

### Architecture Requirements

**No External SDKs:**
- Use stdlib `net/http` for API calls
- No SendGrid SDK, no Resend SDK
- Keep dependencies minimal

**Fire-and-Forget Pattern (for callers):**
```go
// How handlers will use this (in later stories)
go func() {
    if err := emailService.Send(ctx, email); err != nil {
        slog.Error("email send failed", "error", err, "to", email.To)
    }
}()
```

### EmailProvider Interface

```go
// internal/email/provider.go

package email

import "context"

type Email struct {
    To       string
    From     string
    Subject  string
    HTMLBody string
    TextBody string
}

type EmailProvider interface {
    Send(ctx context.Context, email Email) error
}
```

### SendGrid Implementation

**API Endpoint:** `POST https://api.sendgrid.com/v3/mail/send`

**Request Format:**
```json
{
  "personalizations": [
    {
      "to": [{"email": "recipient@example.com"}]
    }
  ],
  "from": {"email": "sender@example.com"},
  "subject": "Subject line",
  "content": [
    {"type": "text/plain", "value": "Plain text body"},
    {"type": "text/html", "value": "<p>HTML body</p>"}
  ]
}
```

**Headers:**
```
Authorization: Bearer <SENDGRID_API_KEY>
Content-Type: application/json
```

**Implementation:**
```go
// internal/email/sendgrid.go

package email

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
)

type SendGridProvider struct {
    apiKey string
    client *http.Client
}

func NewSendGridProvider(apiKey string) *SendGridProvider {
    return &SendGridProvider{
        apiKey: apiKey,
        client: &http.Client{Timeout: 10 * time.Second},
    }
}

func (p *SendGridProvider) Send(ctx context.Context, email Email) error {
    payload := map[string]any{
        "personalizations": []map[string]any{
            {"to": []map[string]string{{"email": email.To}}},
        },
        "from":    map[string]string{"email": email.From},
        "subject": email.Subject,
        "content": []map[string]string{
            {"type": "text/plain", "value": email.TextBody},
            {"type": "text/html", "value": email.HTMLBody},
        },
    }

    body, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("marshal sendgrid payload: %w", err)
    }

    req, err := http.NewRequestWithContext(ctx, "POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewReader(body))
    if err != nil {
        return fmt.Errorf("create sendgrid request: %w", err)
    }

    req.Header.Set("Authorization", "Bearer "+p.apiKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := p.client.Do(req)
    if err != nil {
        return fmt.Errorf("sendgrid request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        return fmt.Errorf("sendgrid returned status %d", resp.StatusCode)
    }

    return nil
}
```

### Resend Implementation

**API Endpoint:** `POST https://api.resend.com/emails`

**Request Format:**
```json
{
  "from": "sender@example.com",
  "to": ["recipient@example.com"],
  "subject": "Subject line",
  "html": "<p>HTML body</p>",
  "text": "Plain text body"
}
```

**Headers:**
```
Authorization: Bearer <RESEND_API_KEY>
Content-Type: application/json
```

**Implementation:**
```go
// internal/email/resend.go

package email

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type ResendProvider struct {
    apiKey string
    client *http.Client
}

func NewResendProvider(apiKey string) *ResendProvider {
    return &ResendProvider{
        apiKey: apiKey,
        client: &http.Client{Timeout: 10 * time.Second},
    }
}

func (p *ResendProvider) Send(ctx context.Context, email Email) error {
    payload := map[string]any{
        "from":    email.From,
        "to":      []string{email.To},
        "subject": email.Subject,
        "html":    email.HTMLBody,
        "text":    email.TextBody,
    }

    body, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("marshal resend payload: %w", err)
    }

    req, err := http.NewRequestWithContext(ctx, "POST", "https://api.resend.com/emails", bytes.NewReader(body))
    if err != nil {
        return fmt.Errorf("create resend request: %w", err)
    }

    req.Header.Set("Authorization", "Bearer "+p.apiKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := p.client.Do(req)
    if err != nil {
        return fmt.Errorf("resend request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        return fmt.Errorf("resend returned status %d", resp.StatusCode)
    }

    return nil
}
```

### Configuration

**Environment Variables:**
```bash
EMAIL_PROVIDER=sendgrid          # or "resend"
SENDGRID_API_KEY=SG.xxx          # required if provider=sendgrid
RESEND_API_KEY=re_xxx            # required if provider=resend
EMAIL_FROM=noreply@example.com   # sender address
FRONTEND_URL=http://localhost:3000  # for verification links
```

**Config Struct Addition:**
```go
// internal/config/config.go

type Config struct {
    // ... existing fields ...

    // Email configuration
    EmailProvider  string // "sendgrid" or "resend"
    SendGridAPIKey string
    ResendAPIKey   string
    EmailFrom      string
    FrontendURL    string
}

func (c *Config) Validate() error {
    // ... existing validation ...

    if c.EmailProvider == "" {
        return fmt.Errorf("EMAIL_PROVIDER is required")
    }
    if c.EmailProvider != "sendgrid" && c.EmailProvider != "resend" {
        return fmt.Errorf("EMAIL_PROVIDER must be 'sendgrid' or 'resend'")
    }
    if c.EmailProvider == "sendgrid" && c.SendGridAPIKey == "" {
        return fmt.Errorf("SENDGRID_API_KEY is required when EMAIL_PROVIDER=sendgrid")
    }
    if c.EmailProvider == "resend" && c.ResendAPIKey == "" {
        return fmt.Errorf("RESEND_API_KEY is required when EMAIL_PROVIDER=resend")
    }
    if c.EmailFrom == "" {
        return fmt.Errorf("EMAIL_FROM is required")
    }
    return nil
}
```

### Provider Factory

```go
// internal/email/provider.go

func NewEmailProvider(cfg *config.Config) (EmailProvider, error) {
    switch cfg.EmailProvider {
    case "sendgrid":
        return NewSendGridProvider(cfg.SendGridAPIKey), nil
    case "resend":
        return NewResendProvider(cfg.ResendAPIKey), nil
    default:
        return nil, fmt.Errorf("unknown email provider: %s", cfg.EmailProvider)
    }
}
```

### Testing Approach

Use httptest to mock API responses:
```go
func TestSendGridProvider_Send(t *testing.T) {
    // Create mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Verify request format
        assert.Equal(t, "POST", r.Method)
        assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))

        var body map[string]any
        json.NewDecoder(r.Body).Decode(&body)
        assert.Equal(t, "Test Subject", body["subject"])

        w.WriteHeader(http.StatusAccepted)
    }))
    defer server.Close()

    // Test with mock server URL
    provider := &SendGridProvider{
        apiKey:  "test-key",
        client:  server.Client(),
        baseURL: server.URL, // Add baseURL field for testing
    }

    err := provider.Send(context.Background(), Email{
        To:      "test@example.com",
        Subject: "Test Subject",
    })
    assert.NoError(t, err)
}
```

### References

- [Source: docs/planning-artifacts/architecture.md#Email Provider Architecture]
- [Source: docs/planning-artifacts/epics.md#Story 2.1]
- [SendGrid v3 API Docs](https://docs.sendgrid.com/api-reference/mail-send/mail-send)
- [Resend API Docs](https://resend.com/docs/api-reference/emails/send-email)

## Dev Agent Record

### Agent Model Used

(To be filled by dev agent)

### Debug Log References

(To be filled during implementation)

### Completion Notes List

(To be filled during implementation)

### Change Log

(To be filled during implementation)

### File List

**Expected New Files:**
- api/internal/email/provider.go
- api/internal/email/sendgrid.go
- api/internal/email/resend.go
- api/internal/email/provider_test.go
- api/internal/email/sendgrid_test.go
- api/internal/email/resend_test.go

**Expected Modified Files:**
- api/internal/config/config.go (add email config fields + validation)
- api/cmd/server/main.go (wire email provider)
