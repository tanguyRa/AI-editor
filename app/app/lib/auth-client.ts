import { createAuthClient } from "better-auth/vue"

export const authClient = createAuthClient({
    //you can pass client configuration here
})


export const {
    signIn,
    signOut,
    signUp,
    useSession,
    requestPasswordReset,
    resetPassword,
} = authClient;