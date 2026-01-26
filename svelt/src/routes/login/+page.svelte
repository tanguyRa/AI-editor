<script lang="ts">
    import { authClient, useSession } from "$lib/auth-client";
    import { goto } from "$app/navigation";

    const session = useSession();

    let email = $state("");
    let password = $state("");
    let loading = $state(false);
    let error = $state("");

    async function handleEmailLogin() {
        if (!email || !password) {
            error = "Please fill in all fields";
            return;
        }
        loading = true;
        error = "";

        const result = await authClient.signIn.email({
            email,
            password,
        });

        if (result.error) {
            error = result.error.message || "Failed to sign in";
        } else {
            goto("/");
        }
        loading = false;
    }

    async function handleGitHubLogin() {
        loading = true;
        error = "";
        await authClient.signIn.social({
            provider: "github",
        });
    }

    async function handleSignOut() {
        await authClient.signOut();
    }
</script>

<div class="auth-container">
    <div class="auth-card">
        {#if $session.data}
            <div class="logged-in">
                <div class="avatar">
                    {$session.data.user.name?.charAt(0).toUpperCase() || "U"}
                </div>
                <h2>Welcome back!</h2>
                <p class="user-name">{$session.data.user.name}</p>
                <p class="user-email">{$session.data.user.email}</p>
                <button class="btn btn-secondary" onclick={handleSignOut}>
                    Sign Out
                </button>
            </div>
        {:else}
            <div class="auth-header">
                <h1>Welcome back</h1>
                <p>Sign in to your account to continue</p>
            </div>

            {#if error}
                <div class="error-message">{error}</div>
            {/if}

            <form onsubmit={(e) => { e.preventDefault(); handleEmailLogin(); }}>
                <div class="form-group">
                    <label for="email">Email</label>
                    <input
                        type="email"
                        id="email"
                        bind:value={email}
                        placeholder="you@example.com"
                        disabled={loading}
                    />
                </div>

                <div class="form-group">
                    <label for="password">Password</label>
                    <input
                        type="password"
                        id="password"
                        bind:value={password}
                        placeholder="Enter your password"
                        disabled={loading}
                    />
                </div>

                <button type="submit" class="btn btn-primary" disabled={loading}>
                    {loading ? "Signing in..." : "Sign In"}
                </button>
            </form>

            <div class="divider">
                <span>or continue with</span>
            </div>

            <button class="btn btn-social" onclick={handleGitHubLogin} disabled={loading}>
                <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor">
                    <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/>
                </svg>
                Continue with GitHub
            </button>

            <p class="auth-footer">
                Don't have an account? <a href="/register">Sign up</a>
            </p>
        {/if}
    </div>
</div>

<style>
    .auth-container {
        min-height: 100vh;
        display: flex;
        align-items: center;
        justify-content: center;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        padding: 1rem;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
    }

    .auth-card {
        background: white;
        border-radius: 16px;
        box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
        padding: 2.5rem;
        width: 100%;
        max-width: 400px;
    }

    .auth-header {
        text-align: center;
        margin-bottom: 2rem;
    }

    .auth-header h1 {
        margin: 0 0 0.5rem;
        font-size: 1.75rem;
        font-weight: 700;
        color: #1a1a2e;
    }

    .auth-header p {
        margin: 0;
        color: #6b7280;
        font-size: 0.95rem;
    }

    .form-group {
        margin-bottom: 1.25rem;
    }

    .form-group label {
        display: block;
        margin-bottom: 0.5rem;
        font-size: 0.875rem;
        font-weight: 500;
        color: #374151;
    }

    .form-group input {
        width: 100%;
        padding: 0.75rem 1rem;
        border: 2px solid #e5e7eb;
        border-radius: 10px;
        font-size: 1rem;
        transition: border-color 0.2s, box-shadow 0.2s;
        box-sizing: border-box;
    }

    .form-group input:focus {
        outline: none;
        border-color: #667eea;
        box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
    }

    .form-group input::placeholder {
        color: #9ca3af;
    }

    .form-group input:disabled {
        background: #f9fafb;
        cursor: not-allowed;
    }

    .btn {
        width: 100%;
        padding: 0.875rem 1.5rem;
        border: none;
        border-radius: 10px;
        font-size: 1rem;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
    }

    .btn:disabled {
        opacity: 0.7;
        cursor: not-allowed;
    }

    .btn-primary {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
    }

    .btn-primary:hover:not(:disabled) {
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
    }

    .btn-secondary {
        background: #f3f4f6;
        color: #374151;
    }

    .btn-secondary:hover:not(:disabled) {
        background: #e5e7eb;
    }

    .btn-social {
        background: #24292e;
        color: white;
        margin-top: 1rem;
    }

    .btn-social:hover:not(:disabled) {
        background: #1a1e22;
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
    }

    .divider {
        display: flex;
        align-items: center;
        margin: 1.5rem 0;
        color: #9ca3af;
        font-size: 0.875rem;
    }

    .divider::before,
    .divider::after {
        content: "";
        flex: 1;
        height: 1px;
        background: #e5e7eb;
    }

    .divider span {
        padding: 0 1rem;
    }

    .auth-footer {
        text-align: center;
        margin-top: 1.5rem;
        color: #6b7280;
        font-size: 0.95rem;
    }

    .auth-footer a {
        color: #667eea;
        text-decoration: none;
        font-weight: 600;
    }

    .auth-footer a:hover {
        text-decoration: underline;
    }

    .error-message {
        background: #fef2f2;
        border: 1px solid #fecaca;
        color: #dc2626;
        padding: 0.75rem 1rem;
        border-radius: 10px;
        margin-bottom: 1.25rem;
        font-size: 0.875rem;
    }

    .logged-in {
        text-align: center;
    }

    .avatar {
        width: 80px;
        height: 80px;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        margin: 0 auto 1.5rem;
        font-size: 2rem;
        font-weight: 700;
        color: white;
    }

    .logged-in h2 {
        margin: 0 0 0.5rem;
        font-size: 1.5rem;
        color: #1a1a2e;
    }

    .user-name {
        margin: 0 0 0.25rem;
        font-size: 1.125rem;
        font-weight: 600;
        color: #374151;
    }

    .user-email {
        margin: 0 0 1.5rem;
        color: #6b7280;
        font-size: 0.95rem;
    }
</style>
