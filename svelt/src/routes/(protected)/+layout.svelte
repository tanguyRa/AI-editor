<script lang="ts">
    import type { LayoutProps } from "./$types";
    import { useSession, signOut } from "$lib/auth-client";
    import { goto } from "$app/navigation";

    let { children }: LayoutProps = $props();

    const session = useSession();

    let sidebarCollapsed = $state(false);

    // Load collapsed state from localStorage on mount
    $effect(() => {
        if (typeof window !== "undefined") {
            const stored = localStorage.getItem("sidebar-collapsed");
            if (stored !== null) {
                sidebarCollapsed = stored === "true";
            }
        }
    });

    // Persist collapsed state to localStorage
    function toggleSidebar() {
        sidebarCollapsed = !sidebarCollapsed;
        if (typeof window !== "undefined") {
            localStorage.setItem("sidebar-collapsed", String(sidebarCollapsed));
        }
    }

    async function handleLogout() {
        await signOut();
        goto("/");
    }
</script>

{#if $session.isPending}
    <div class="loading-container">
        <div class="spinner spinner-dark"></div>
    </div>
{:else if !$session.data}
    <div class="forbidden-container">
        <div class="forbidden-card">
            <div class="forbidden-icon">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="64"
                    height="64"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <circle cx="12" cy="12" r="10"></circle>
                    <line x1="4.93" y1="4.93" x2="19.07" y2="19.07"></line>
                </svg>
            </div>
            <h1>403</h1>
            <h2>Access Denied</h2>
            <p>You must be logged in to access this page.</p>
            <div class="forbidden-actions">
                <a href="/login" class="btn btn-primary">Log In</a>
                <a href="/" class="btn btn-secondary">Go Home</a>
            </div>
        </div>
    </div>
{:else}
    <div class="app-layout" class:collapsed={sidebarCollapsed}>
        <aside class="sidebar">
            <div class="sidebar-header">
                {#if !sidebarCollapsed}
                    <a href="/app" class="logo">
                        <svg
                            class="logo-icon"
                            xmlns="http://www.w3.org/2000/svg"
                            width="24"
                            height="24"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        >
                            <path
                                d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"
                            ></path>
                            <polyline points="22,6 12,13 2,6"></polyline>
                        </svg>
                        <span class="logo-text">Flowmail</span>
                    </a>
                {/if}
                <button
                    class="toggle-btn"
                    onclick={toggleSidebar}
                    aria-label="Toggle sidebar"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    >
                        {#if sidebarCollapsed}
                            <polyline points="9 18 15 12 9 6"></polyline>
                        {:else}
                            <polyline points="15 18 9 12 15 6"></polyline>
                        {/if}
                    </svg>
                </button>
            </div>

            <nav class="sidebar-nav">
                <div class="nav-section">
                    <a href="/app" class="nav-link">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            width="20"
                            height="20"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        >
                            <path
                                d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"
                            ></path>
                        </svg>
                        {#if !sidebarCollapsed}
                            <span>Projects</span>
                        {/if}
                    </a>
                </div>

                <div class="nav-section nav-section-bottom">
                    <span class="nav-section-label"
                        >{sidebarCollapsed ? "" : "Settings"}</span
                    >
                    <a href="/settings/profile" class="nav-link">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            width="20"
                            height="20"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        >
                            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"
                            ></path>
                            <circle cx="12" cy="7" r="4"></circle>
                        </svg>
                        {#if !sidebarCollapsed}
                            <span>Profile</span>
                        {/if}
                    </a>
                    <a href="/settings/billing" class="nav-link">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            width="20"
                            height="20"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        >
                            <rect
                                x="1"
                                y="4"
                                width="22"
                                height="16"
                                rx="2"
                                ry="2"
                            ></rect>
                            <line x1="1" y1="10" x2="23" y2="10"></line>
                        </svg>
                        {#if !sidebarCollapsed}
                            <span>Billing</span>
                        {/if}
                    </a>
                    <button class="nav-link logout-btn" onclick={handleLogout}>
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            width="20"
                            height="20"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        >
                            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"
                            ></path>
                            <polyline points="16 17 21 12 16 7"></polyline>
                            <line x1="21" y1="12" x2="9" y2="12"></line>
                        </svg>
                        {#if !sidebarCollapsed}
                            <span>Logout</span>
                        {/if}
                    </button>
                </div>
            </nav>
        </aside>

        <main class="main-content">
            {@render children()}
        </main>
    </div>
{/if}

<style>
    /* Loading State */
    .loading-container {
        min-height: 100vh;
        display: flex;
        align-items: center;
        justify-content: center;
        background: var(--color-bg-secondary);
    }

    /* Forbidden State */
    .forbidden-container {
        min-height: 100vh;
        display: flex;
        align-items: center;
        justify-content: center;
        background: var(--gradient-primary);
        padding: var(--spacing-md);
    }

    .forbidden-card {
        background: white;
        border-radius: var(--radius-xl);
        box-shadow: var(--shadow-xl);
        padding: var(--spacing-3xl);
        text-align: center;
        max-width: 400px;
        width: 100%;
    }

    .forbidden-icon {
        color: var(--color-error);
        margin-bottom: var(--spacing-lg);
    }

    .forbidden-card h1 {
        font-size: 4rem;
        color: var(--color-error);
        margin-bottom: var(--spacing-xs);
        line-height: 1;
    }

    .forbidden-card h2 {
        font-size: var(--font-size-2xl);
        color: var(--color-text);
        margin-bottom: var(--spacing-md);
    }

    .forbidden-card p {
        color: var(--color-text-muted);
        margin-bottom: var(--spacing-xl);
    }

    .forbidden-actions {
        display: flex;
        gap: var(--spacing-md);
        justify-content: center;
    }

    /* App Layout */
    .app-layout {
        display: flex;
        min-height: 100vh;
    }

    /* Sidebar */
    .sidebar {
        width: 240px;
        background: var(--color-bg);
        border-right: 1px solid var(--color-border);
        display: flex;
        flex-direction: column;
        transition: width var(--transition-slow);
        position: fixed;
        top: 0;
        left: 0;
        bottom: 0;
        z-index: var(--z-sticky);
    }

    .app-layout.collapsed .sidebar {
        width: 72px;
    }

    .sidebar-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: var(--spacing-lg);
        border-bottom: 1px solid var(--color-border);
    }

    .sidebar-header .logo {
        display: flex;
        align-items: center;
        gap: var(--spacing-sm);
        text-decoration: none;
        color: var(--color-text);
        overflow: hidden;
    }

    .sidebar-header .logo-icon {
        color: var(--color-primary);
        flex-shrink: 0;
    }

    .sidebar-header .logo-text {
        font-size: var(--font-size-lg);
        font-weight: 700;
        white-space: nowrap;
    }

    .toggle-btn {
        background: none;
        border: none;
        padding: var(--spacing-sm);
        border-radius: var(--radius-sm);
        cursor: pointer;
        color: var(--color-text-muted);
        transition: all var(--transition-fast);
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
    }

    .toggle-btn:hover {
        background: var(--color-bg-tertiary);
        color: var(--color-text);
    }

    .app-layout.collapsed .toggle-btn {
        margin-left: auto;
        margin-right: auto;
    }

    .app-layout.collapsed .sidebar-header {
        flex-direction: column;
        gap: var(--spacing-sm);
    }

    /* Sidebar Navigation */
    .sidebar-nav {
        flex: 1;
        display: flex;
        flex-direction: column;
        padding: var(--spacing-md);
        overflow-y: auto;
    }

    .nav-section {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-xs);
    }

    .nav-section-bottom {
        margin-top: auto;
        border-top: 1px solid var(--color-border);
        padding-top: var(--spacing-md);
    }

    .nav-section-label {
        font-size: var(--font-size-xs);
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.5px;
        color: var(--color-text-light);
        padding: var(--spacing-sm) var(--spacing-sm);
        min-height: 24px;
    }

    .nav-link {
        display: flex;
        align-items: center;
        gap: var(--spacing-md);
        padding: var(--spacing-sm) var(--spacing-md);
        border-radius: var(--radius-md);
        color: var(--color-text-secondary);
        text-decoration: none;
        font-weight: 500;
        transition: all var(--transition-fast);
        background: none;
        border: none;
        width: 100%;
        cursor: pointer;
        font-size: var(--font-size-base);
        font-family: var(--font-family);
    }

    .nav-link:hover {
        background: var(--color-bg-tertiary);
        color: var(--color-text);
    }

    .nav-link svg {
        flex-shrink: 0;
    }

    .nav-link span {
        white-space: nowrap;
        overflow: hidden;
    }

    .app-layout.collapsed .nav-link {
        justify-content: center;
        padding: var(--spacing-sm);
    }

    .logout-btn {
        color: var(--color-error);
    }

    .logout-btn:hover {
        background: var(--color-error-bg);
        color: var(--color-error);
    }

    /* Main Content */
    .main-content {
        flex: 1;
        margin-left: 240px;
        background: var(--color-bg-secondary);
        min-height: 100vh;
        transition: margin-left var(--transition-slow);
    }

    .app-layout.collapsed .main-content {
        margin-left: 72px;
    }

    /* Responsive */
    @media (max-width: 768px) {
        .sidebar {
            transform: translateX(-100%);
            width: 240px !important;
        }

        .sidebar.open {
            transform: translateX(0);
        }

        .main-content {
            margin-left: 0 !important;
        }
    }
</style>
