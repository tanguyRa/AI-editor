<script lang="ts">
    import type { LayoutProps } from "./$types";
    import { authClient, useSession } from "$lib/auth-client";

    let { children }: LayoutProps = $props();

    const session = useSession();

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
                <h2>You're already signed in!</h2>
                <p class="user-name">{$session.data.user.name}</p>
                <p class="user-email">{$session.data.user.email}</p>
                <button class="btn btn-secondary full-width" onclick={handleSignOut}>
                    Sign Out
                </button>
            </div>
        {:else}
            {@render children()}
        {/if}
    </div>
</div>

<style>
    .full-width {
        width: 100%;
    }

    .avatar {
        margin: 0 auto var(--spacing-lg);
    }
</style>
