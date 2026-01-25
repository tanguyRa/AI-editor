export default defineEventHandler((event) => {
    const {
        private: { polarAccessToken, polarCheckoutSuccessUrl, polarServer },
    } = useRuntimeConfig();

    const checkoutHandler = Checkout({
        accessToken: polarAccessToken,
        successUrl: polarCheckoutSuccessUrl,
        returnUrl: "https://myapp.com", // An optional URL which renders a back-button in the Checkout
        server: polarServer as "sandbox" | "production",
        theme: "dark", // Enforces the theme - System-preferred theme will be set if left omitted
    });

    return checkoutHandler(event);
});