import { authClient } from "~/lib/auth-client";
export default defineNuxtRouteMiddleware(async (to, from) => {
    const { data: session } = await authClient.useSession(useFetch);

    console.log("Current session in middleware:", session.value);

    if (!session.value) {
        return navigateTo("/");
    }
});