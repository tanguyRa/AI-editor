# Project Documentation

SvelteKit app with better-auth, Polar.sh payments, and Neon DB.

---

## Route Structure

```
src/routes/
├── (marketing)/          # Public pages
│   ├── +page.svelte      # Landing (/)
│   ├── pricing/
│   └── +layout.svelte
│
├── (auth)/               # Auth pages
│   ├── login/
│   ├── register/
│   └── +layout.svelte
│
├── (protected)/          # Requires auth
│   ├── +layout.svelte    # Shared sidebar layout
│   ├── +layout.server.ts # Auth guard → 403 if not logged in
│   ├── app/
│   │   ├── +page.svelte  # Dashboard (project list)
│   │   └── [slug]/
│   └── settings/
│
└── api/auth/[...all]/    # better-auth catch-all
```

**Auth guard:** Returns 403 page with links to `/login` and `/`.

**Sidebar:** Collapsible (icons only / expanded). State stored in localStorage.

---

## Polar Integration

- Polar plugin auto-creates Polar customer on user registration
- Polar uses `external_id` referencing our `user.id` — no need to store `polar_customer_id`
- All tier/status updates come from Polar webhooks

### Tier Mapping

Map Polar product IDs to tier values in config:

| Polar Product ID | Tier |
|------------------|------|
| prod_xxx | premium_1 |
| prod_yyy | premium_2 |

---

## Onboarding Flows

### Flow A: Landing CTA → Free

```
Landing → CTA → /register → account created → subscription (tier=free) → /app
```

### Flow B: Pricing → Paid

```
/pricing → select plan → /register?plan=X → account created (tier=free)
    → Polar checkout → payment complete → webhook updates tier
    → redirect to /app?checkout_id={ID}

If checkout cancelled/abandoned → redirect to /
```

**Key:** Always create as `free` first. Webhook upgrades tier after payment confirmed.

---

## Subscription Lifecycle

### Upgrade

```
Webhook: subscription.created
→ Set tier from product ID, status = 'active', current_period_end
```

### Scheduled Downgrade/Cancel

```
Webhook: subscription.updated (with pending change)
→ Set scheduled_tier to upcoming tier (or 'free' if cancelling)
→ tier unchanged until period ends
```

### Period End (Downgrade Executes)

```
Webhook: subscription.updated or subscription.canceled
→ tier = scheduled_tier (or 'free')
→ scheduled_tier = null
```

### Payment Failure

```
Webhook: subscription.updated (status = 'past_due')
→ status = 'past_due', tier unchanged
→ Grace period: 2 days (Polar-managed)
→ If not resolved: Polar sends cancellation webhook → tier = 'free'
```

---

## Project Limits

| Tier | Max Projects |
|------|--------------|
| free | 1 |
| premium_1 | unlimited |
| premium_2 | unlimited |

### Free User with Multiple Projects (Downgrade)

Projects remain in DB but only the most recently updated is visible.

**Query logic:**

```
SELECT * FROM project WHERE user_id = X ORDER BY updated_at DESC
If tier = 'free': return first only
Else: return all
```

---

## UI States

### Account/Billing Page

| State | Display |
|-------|---------|
| `tier=free` | "Free plan — [Upgrade]" |
| `tier=premium_x`, `scheduled_tier=null` | "Premium — renews {date} — [Cancel]" |
| `tier=premium_x`, `scheduled_tier=free` | "Premium until {date}. You won't be charged. — [Resubscribe]" |
| `tier=premium_x`, `scheduled_tier=premium_1` | "Switching to Premium 1 on {date}" |
| `status=past_due` | "Payment failed — [Update payment]" |

### Checkout Return

Landing on `/app?checkout_id={ID}`:

| Webhook arrived? | Action |
|------------------|--------|
| Yes (tier updated) | Show success, premium access |
| No (race condition) | Brief loading/polling, then refresh |
