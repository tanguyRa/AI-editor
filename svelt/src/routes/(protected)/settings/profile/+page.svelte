<script lang="ts">
    import {
        useSession,
        updateUser,
        changeEmail,
        changePassword,
    } from "$lib/auth-client";

    const session = useSession();

    let name = $state("");
    let email = $state("");
    let currentPassword = $state("");
    let newPassword = $state("");
    let confirmPassword = $state("");

    let profileLoading = $state(false);
    let emailLoading = $state(false);
    let passwordLoading = $state(false);

    let profileMessage = $state<{
        type: "success" | "error";
        text: string;
    } | null>(null);
    let emailMessage = $state<{
        type: "success" | "error";
        text: string;
    } | null>(null);
    let passwordMessage = $state<{
        type: "success" | "error";
        text: string;
    } | null>(null);

    // Initialize form values when session loads
    $effect(() => {
        if ($session.data?.user) {
            name = $session.data.user.name || "";
            email = $session.data.user.email || "";
        }
    });

    async function handleUpdateProfile() {
        profileMessage = null;
        profileLoading = true;

        try {
            const { error } = await updateUser({ name });
            if (error) {
                profileMessage = {
                    type: "error",
                    text: error.message || "Failed to update profile",
                };
            } else {
                profileMessage = {
                    type: "success",
                    text: "Profile updated successfully",
                };
            }
        } catch {
            profileMessage = {
                type: "error",
                text: "An unexpected error occurred",
            };
        } finally {
            profileLoading = false;
        }
    }

    async function handleUpdateEmail() {
        emailMessage = null;
        emailLoading = true;

        try {
            const { error } = await changeEmail({ newEmail: email });
            if (error) {
                emailMessage = {
                    type: "error",
                    text: error.message || "Failed to update email",
                };
            } else {
                emailMessage = {
                    type: "success",
                    text: "Verification email sent. Please check your inbox.",
                };
            }
        } catch {
            emailMessage = {
                type: "error",
                text: "An unexpected error occurred",
            };
        } finally {
            emailLoading = false;
        }
    }

    async function handleUpdatePassword() {
        passwordMessage = null;

        if (newPassword !== confirmPassword) {
            passwordMessage = { type: "error", text: "Passwords do not match" };
            return;
        }

        if (newPassword.length < 8) {
            passwordMessage = {
                type: "error",
                text: "Password must be at least 8 characters",
            };
            return;
        }

        passwordLoading = true;

        try {
            const { error } = await changePassword({
                currentPassword,
                newPassword,
            });
            if (error) {
                passwordMessage = {
                    type: "error",
                    text: error.message || "Failed to update password",
                };
            } else {
                passwordMessage = {
                    type: "success",
                    text: "Password updated successfully",
                };
                currentPassword = "";
                newPassword = "";
                confirmPassword = "";
            }
        } catch {
            passwordMessage = {
                type: "error",
                text: "An unexpected error occurred",
            };
        } finally {
            passwordLoading = false;
        }
    }
</script>

<div class="settings-page">
    <header class="settings-header">
        <h1>Profile Settings</h1>
        <p>Manage your account information</p>
    </header>

    {#if $session.data?.user}
        <div class="settings-sections">
            <!-- Profile Section -->
            <section class="settings-section">
                <div class="section-header">
                    <h2>Profile Information</h2>
                    <p>Update your personal details</p>
                </div>

                <form
                    class="settings-form"
                    onsubmit={(e) => {
                        e.preventDefault();
                        handleUpdateProfile();
                    }}
                >
                    <div class="form-group">
                        <label for="name">Name</label>
                        <input
                            type="text"
                            id="name"
                            bind:value={name}
                            placeholder="Your name"
                            required
                        />
                    </div>

                    {#if profileMessage}
                        <div
                            class="message"
                            class:success={profileMessage.type === "success"}
                            class:error={profileMessage.type === "error"}
                        >
                            {profileMessage.text}
                        </div>
                    {/if}

                    <button
                        type="submit"
                        class="btn btn-primary"
                        disabled={profileLoading}
                    >
                        {#if profileLoading}
                            <span class="spinner spinner-sm"></span>
                        {/if}
                        Save Changes
                    </button>
                </form>
            </section>

            <!-- Email Section -->
            <section class="settings-section">
                <div class="section-header">
                    <h2>Email Address</h2>
                    <p>Change your email address</p>
                </div>

                <form
                    class="settings-form"
                    onsubmit={(e) => {
                        e.preventDefault();
                        handleUpdateEmail();
                    }}
                >
                    <div class="form-group">
                        <label for="email">Email</label>
                        <input
                            type="email"
                            id="email"
                            bind:value={email}
                            placeholder="your@email.com"
                            required
                        />
                        <span class="form-hint"
                            >A verification email will be sent to confirm the
                            change</span
                        >
                    </div>

                    {#if emailMessage}
                        <div
                            class="message"
                            class:success={emailMessage.type === "success"}
                            class:error={emailMessage.type === "error"}
                        >
                            {emailMessage.text}
                        </div>
                    {/if}

                    <button
                        type="submit"
                        class="btn btn-primary"
                        disabled={emailLoading ||
                            email === $session.data.user.email}
                    >
                        {#if emailLoading}
                            <span class="spinner spinner-sm"></span>
                        {/if}
                        Update Email
                    </button>
                </form>
            </section>

            <!-- Password Section -->
            <section class="settings-section">
                <div class="section-header">
                    <h2>Password</h2>
                    <p>Change your password</p>
                </div>

                <form
                    class="settings-form"
                    onsubmit={(e) => {
                        e.preventDefault();
                        handleUpdatePassword();
                    }}
                >
                    <div class="form-group">
                        <label for="current-password">Current Password</label>
                        <input
                            type="password"
                            id="current-password"
                            bind:value={currentPassword}
                            placeholder="Enter current password"
                            required
                        />
                    </div>

                    <div class="form-group">
                        <label for="new-password">New Password</label>
                        <input
                            type="password"
                            id="new-password"
                            bind:value={newPassword}
                            placeholder="Enter new password"
                            required
                            minlength="8"
                        />
                        <span class="form-hint">Minimum 8 characters</span>
                    </div>

                    <div class="form-group">
                        <label for="confirm-password"
                            >Confirm New Password</label
                        >
                        <input
                            type="password"
                            id="confirm-password"
                            bind:value={confirmPassword}
                            placeholder="Confirm new password"
                            required
                        />
                    </div>

                    {#if passwordMessage}
                        <div
                            class="message"
                            class:success={passwordMessage.type === "success"}
                            class:error={passwordMessage.type === "error"}
                        >
                            {passwordMessage.text}
                        </div>
                    {/if}

                    <button
                        type="submit"
                        class="btn btn-primary"
                        disabled={passwordLoading ||
                            !currentPassword ||
                            !newPassword ||
                            !confirmPassword}
                    >
                        {#if passwordLoading}
                            <span class="spinner spinner-sm"></span>
                        {/if}
                        Change Password
                    </button>
                </form>
            </section>
        </div>
    {:else}
        <div class="loading-state">
            <div class="spinner spinner-dark"></div>
        </div>
    {/if}
</div>

<style>
    .settings-page {
        padding: var(--spacing-xl);
        max-width: 800px;
    }

    .settings-header {
        margin-bottom: var(--spacing-2xl);
    }

    .settings-header h1 {
        font-size: var(--font-size-3xl);
        color: var(--color-text);
        margin-bottom: var(--spacing-sm);
    }

    .settings-header p {
        color: var(--color-text-muted);
        font-size: var(--font-size-lg);
    }

    .settings-sections {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-xl);
    }

    .settings-section {
        background: var(--color-bg);
        border-radius: var(--radius-lg);
        padding: var(--spacing-xl);
        box-shadow: var(--shadow-sm);
    }

    .section-header {
        margin-bottom: var(--spacing-lg);
        padding-bottom: var(--spacing-md);
        border-bottom: 1px solid var(--color-border);
    }

    .section-header h2 {
        font-size: var(--font-size-xl);
        color: var(--color-text);
        margin-bottom: var(--spacing-xs);
    }

    .section-header p {
        color: var(--color-text-muted);
        font-size: var(--font-size-sm);
    }

    .settings-form {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-md);
    }

    .settings-form .btn {
        align-self: flex-start;
        margin-top: var(--spacing-sm);
    }

    .message {
        padding: var(--spacing-sm) var(--spacing-md);
        border-radius: var(--radius-md);
        font-size: var(--font-size-sm);
    }

    .message.success {
        background: #ecfdf5;
        color: var(--color-success);
        border: 1px solid #a7f3d0;
    }

    .message.error {
        background: var(--color-error-bg);
        color: var(--color-error);
        border: 1px solid var(--color-error-border);
    }

    .loading-state {
        display: flex;
        justify-content: center;
        padding: var(--spacing-3xl);
    }

    @media (max-width: 768px) {
        .settings-page {
            padding: var(--spacing-md);
        }

        .settings-section {
            padding: var(--spacing-lg);
        }
    }
</style>
