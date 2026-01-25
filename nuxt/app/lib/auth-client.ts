import { createAuthClient } from "better-auth/vue"
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
    checkout, usage, customer
} = authClient;