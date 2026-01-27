<script lang="ts">
    import type { LayoutProps } from "./$types";
    import { useSession } from "$lib/auth-client";
    import { goto } from "$app/navigation";

    let { children }: LayoutProps = $props();

    const session = useSession();

    $effect(() => {
        if ($session.data) {
            goto("/app");
        }
    });
</script>

{#if !$session.data}
    <div class="auth-container">
        <div class="auth-card">
            {@render children()}
        </div>
    </div>
{/if}

<style>
    .full-width {
        width: 100%;
    }

    .avatar {
        margin: 0 auto var(--spacing-lg);
    }
</style>
