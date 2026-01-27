import { createAuthClient } from "better-auth/svelte"
import { polarClient } from "@polar-sh/better-auth";

export const authClient = createAuthClient({
    //you can pass client configuration here
    plugins: [
        polarClient()
    ]
})


export const {
    signIn,
    signOut,
    signUp,
    useSession,
    requestPasswordReset,
    resetPassword,
    changeEmail,
    changePassword,
    updateUser,
    checkout, usage, customer,
} = authClient;