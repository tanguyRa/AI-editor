<script lang="ts">
    import { useSession } from "$lib/auth-client";

    const session = useSession();

    const features = [
        {
            icon: "sparkles",
            title: "AI-Powered Sorting",
            description:
                "Automatically categorize and prioritize your emails using advanced machine learning algorithms.",
        },
        {
            icon: "zap",
            title: "Smart Replies",
            description:
                "Generate contextual response suggestions that match your writing style and tone.",
        },
        {
            icon: "shield",
            title: "Privacy First",
            description:
                "Your emails stay private. All AI processing happens securely with enterprise-grade encryption.",
        },
        {
            icon: "clock",
            title: "Schedule & Snooze",
            description:
                "Send emails at the perfect time and snooze messages to resurface when you need them.",
        },
        {
            icon: "search",
            title: "Instant Search",
            description:
                "Find any email in milliseconds with powerful full-text search and smart filters.",
        },
        {
            icon: "layout",
            title: "Unified Inbox",
            description:
                "Manage all your email accounts from one beautiful, distraction-free interface.",
        },
    ];

    const integrations = [
        { name: "Gmail", logo: "G" },
        { name: "Outlook", logo: "O" },
        { name: "iCloud", logo: "iC" },
        { name: "Yahoo", logo: "Y!" },
        { name: "ProtonMail", logo: "P" },
        { name: "Fastmail", logo: "FM" },
    ];

    const faqs = [
        {
            question: "How does Flowmail's AI work?",
            answer: "Flowmail uses advanced natural language processing to understand your emails and suggest actions. The AI learns from your behavior to provide increasingly personalized suggestions while keeping your data private and secure.",
        },
        {
            question: "Is my email data secure?",
            answer: "Absolutely. We use end-to-end encryption and never store your email content on our servers. All AI processing is done securely, and we never sell or share your data with third parties.",
        },
        {
            question: "Can I use Flowmail with my existing email accounts?",
            answer: "Yes! Flowmail works with all major email providers including Gmail, Outlook, iCloud, Yahoo, and any IMAP-compatible email service.",
        },
        {
            question: "Is there a free plan available?",
            answer: "Yes, we offer a generous free tier that includes core features. Premium plans unlock advanced AI capabilities, priority support, and additional integrations.",
        },
        {
            question: "How do I migrate from my current email client?",
            answer: "Simply connect your email accounts and Flowmail will sync everything automatically. Your existing folders, labels, and settings are preserved. No manual migration needed.",
        },
    ];

    let openFaq = $state<number | null>(null);

    function toggleFaq(index: number) {
        openFaq = openFaq === index ? null : index;
    }
</script>

<div class="landing">
    <!-- Navigation -->
    <nav class="nav">
        <div class="nav-container">
            <a href="/" class="logo">
                <span class="logo-icon">
                    <svg
                        viewBox="0 0 24 24"
                        width="28"
                        height="28"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                    >
                        <path
                            d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"
                        />
                        <polyline points="22,6 12,13 2,6" />
                    </svg>
                </span>
                <span class="logo-text">flowmail</span>
            </a>
            <div class="nav-links">
                <a href="#features">Features</a>
                <a href="#integrations">Integrations</a>
                <a href="#about">About</a>
                <a href="#faq">FAQ</a>
            </div>
            <div class="nav-actions">
                {#if $session.data?.user}
                    <a href="/app" class="btn btn-secondary">Dashboard</a>
                {:else}
                    <a href="/login" class="btn btn-ghost">Sign In</a>
                    <a href="/register" class="btn btn-primary">Get Started</a>
                {/if}
            </div>
        </div>
    </nav>

    <!-- Hero Section -->
    <section class="hero">
        <div class="hero-bg"></div>
        <div class="hero-content">
            <div class="hero-badge">Introducing Flowmail 2.0</div>
            <h1>
                Email that works <span class="hero-gradient-text">for you</span>
            </h1>
            <p class="hero-subtitle">
                The smart email client that uses AI to help you focus on what
                matters. Sort, reply, and manage your inbox in a fraction of the
                time.
            </p>
            <div class="hero-actions">
                <a href="/register" class="btn btn-primary btn-lg">
                    Start for Free
                    <svg
                        viewBox="0 0 24 24"
                        width="20"
                        height="20"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                    >
                        <line x1="5" y1="12" x2="19" y2="12" />
                        <polyline points="12 5 19 12 12 19" />
                    </svg>
                </a>
                <a href="#features" class="btn btn-outline btn-lg">
                    See Features
                </a>
            </div>
            <p class="hero-note">
                Free forever for personal use. No credit card required.
            </p>
        </div>
        <div class="hero-visual">
            <div class="app-preview">
                <div class="preview-header">
                    <div class="preview-dots">
                        <span></span><span></span><span></span>
                    </div>
                    <span class="preview-title">Inbox</span>
                </div>
                <div class="preview-content">
                    <div class="preview-email priority">
                        <div class="email-avatar">JD</div>
                        <div class="email-body">
                            <div class="email-header">
                                <span class="email-sender">John Doe</span>
                                <span class="email-badge">Priority</span>
                            </div>
                            <span class="email-subject"
                                >Q4 Report - Action Required</span
                            >
                        </div>
                    </div>
                    <div class="preview-email">
                        <div class="email-avatar">AS</div>
                        <div class="email-body">
                            <span class="email-sender">Alice Smith</span>
                            <span class="email-subject"
                                >Meeting notes from today</span
                            >
                        </div>
                    </div>
                    <div class="preview-email">
                        <div class="email-avatar">TB</div>
                        <div class="email-body">
                            <span class="email-sender">Team Bot</span>
                            <span class="email-subject"
                                >Weekly summary ready</span
                            >
                        </div>
                    </div>
                    <div class="preview-ai">
                        <svg
                            viewBox="0 0 24 24"
                            width="16"
                            height="16"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                        >
                            <polygon
                                points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"
                            />
                        </svg>
                        <span>AI sorted 47 emails while you were away</span>
                    </div>
                </div>
            </div>
        </div>
    </section>

    <!-- Features Section -->
    <section id="features" class="section features-section">
        <div class="container">
            <div class="section-header">
                <h2>
                    Powerful features, <span class="gradient-text"
                        >simple experience</span
                    >
                </h2>
                <p>
                    Everything you need to take control of your inbox and
                    reclaim your time.
                </p>
            </div>
            <div class="features-grid">
                {#each features as feature}
                    <div class="feature-card card-hover">
                        <div class="feature-icon">
                            {#if feature.icon === "sparkles"}
                                <svg
                                    viewBox="0 0 24 24"
                                    width="24"
                                    height="24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                >
                                    <path
                                        d="M12 3l1.5 4.5L18 9l-4.5 1.5L12 15l-1.5-4.5L6 9l4.5-1.5L12 3z"
                                    />
                                    <path
                                        d="M5 19l1 3 1-3 3-1-3-1-1-3-1 3-3 1 3 1z"
                                    />
                                </svg>
                            {:else if feature.icon === "zap"}
                                <svg
                                    viewBox="0 0 24 24"
                                    width="24"
                                    height="24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                >
                                    <polygon
                                        points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"
                                    />
                                </svg>
                            {:else if feature.icon === "shield"}
                                <svg
                                    viewBox="0 0 24 24"
                                    width="24"
                                    height="24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                >
                                    <path
                                        d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"
                                    />
                                </svg>
                            {:else if feature.icon === "clock"}
                                <svg
                                    viewBox="0 0 24 24"
                                    width="24"
                                    height="24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                >
                                    <circle cx="12" cy="12" r="10" />
                                    <polyline points="12 6 12 12 16 14" />
                                </svg>
                            {:else if feature.icon === "search"}
                                <svg
                                    viewBox="0 0 24 24"
                                    width="24"
                                    height="24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                >
                                    <circle cx="11" cy="11" r="8" />
                                    <line
                                        x1="21"
                                        y1="21"
                                        x2="16.65"
                                        y2="16.65"
                                    />
                                </svg>
                            {:else if feature.icon === "layout"}
                                <svg
                                    viewBox="0 0 24 24"
                                    width="24"
                                    height="24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                >
                                    <rect
                                        x="3"
                                        y="3"
                                        width="18"
                                        height="18"
                                        rx="2"
                                        ry="2"
                                    />
                                    <line x1="3" y1="9" x2="21" y2="9" />
                                    <line x1="9" y1="21" x2="9" y2="9" />
                                </svg>
                            {/if}
                        </div>
                        <h3>{feature.title}</h3>
                        <p>{feature.description}</p>
                    </div>
                {/each}
            </div>
        </div>
    </section>

    <!-- Integrations Section -->
    <section id="integrations" class="section integrations-section">
        <div class="container">
            <div class="section-header">
                <h2>
                    Works with your <span class="gradient-text"
                        >favorite providers</span
                    >
                </h2>
                <p>
                    Connect all your email accounts and manage them from one
                    place.
                </p>
            </div>
            <div class="integrations-grid">
                {#each integrations as integration}
                    <div class="integration-card">
                        <div class="integration-logo">{integration.logo}</div>
                        <span class="integration-name">{integration.name}</span>
                    </div>
                {/each}
            </div>
            <p class="integrations-note">
                Plus any IMAP-compatible email service
            </p>
        </div>
    </section>

    <!-- About Section -->
    <section id="about" class="section about-section">
        <div class="container">
            <div class="about-content">
                <div class="about-text">
                    <h2>
                        Built for the way <span class="gradient-text"
                            >you work</span
                        >
                    </h2>
                    <p>
                        We believe email should be a tool that empowers you, not
                        a burden that drains your energy. Flowmail was born from
                        our frustration with cluttered inboxes and endless email
                        threads.
                    </p>
                    <p>
                        Our mission is simple: use the latest AI technology to
                        give you back the time you spend managing email, so you
                        can focus on the work that actually matters.
                    </p>
                    <div class="about-stats">
                        <div class="stat">
                            <span class="stat-value">2M+</span>
                            <span class="stat-label">Emails sorted daily</span>
                        </div>
                        <div class="stat">
                            <span class="stat-value">50K+</span>
                            <span class="stat-label">Happy users</span>
                        </div>
                        <div class="stat">
                            <span class="stat-value">3hrs</span>
                            <span class="stat-label">Saved per week</span>
                        </div>
                    </div>
                </div>
                <div class="about-visual">
                    <div class="about-card">
                        <div class="about-card-icon">
                            <svg
                                viewBox="0 0 24 24"
                                width="32"
                                height="32"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                            >
                                <path
                                    d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"
                                />
                            </svg>
                        </div>
                        <h3>Made with care</h3>
                        <p>
                            Every feature is thoughtfully designed to make your
                            email experience delightful.
                        </p>
                    </div>
                </div>
            </div>
        </div>
    </section>

    <!-- FAQ Section -->
    <section id="faq" class="section faq-section">
        <div class="container">
            <div class="section-header">
                <h2>
                    Frequently asked <span class="gradient-text">questions</span
                    >
                </h2>
                <p>Everything you need to know about Flowmail.</p>
            </div>
            <div class="faq-list">
                {#each faqs as faq, index}
                    <div class="faq-item" class:open={openFaq === index}>
                        <button
                            class="faq-question"
                            onclick={() => toggleFaq(index)}
                        >
                            <span>{faq.question}</span>
                            <svg
                                viewBox="0 0 24 24"
                                width="20"
                                height="20"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                            >
                                <polyline points="6 9 12 15 18 9" />
                            </svg>
                        </button>
                        {#if openFaq === index}
                            <div class="faq-answer">
                                <p>{faq.answer}</p>
                            </div>
                        {/if}
                    </div>
                {/each}
            </div>
        </div>
    </section>

    <!-- CTA Section -->
    <section class="cta-section">
        <div class="container">
            <div class="cta-content">
                <h2>Ready to transform your inbox?</h2>
                <p>
                    Join thousands of professionals who have already made the
                    switch to smarter email.
                </p>
                <a href="/register" class="btn btn-white btn-lg">
                    Get Started Free
                    <svg
                        viewBox="0 0 24 24"
                        width="20"
                        height="20"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                    >
                        <line x1="5" y1="12" x2="19" y2="12" />
                        <polyline points="12 5 19 12 12 19" />
                    </svg>
                </a>
            </div>
        </div>
    </section>

    <!-- Footer -->
    <footer class="footer">
        <div class="footer-container">
            <div class="footer-brand">
                <a href="/" class="logo">
                    <span class="logo-icon">
                        <svg
                            viewBox="0 0 24 24"
                            width="24"
                            height="24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                        >
                            <path
                                d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"
                            />
                            <polyline points="22,6 12,13 2,6" />
                        </svg>
                    </span>
                    <span class="logo-text">flowmail</span>
                </a>
                <p>The smart email client that works for you.</p>
            </div>
            <div class="footer-links">
                <div class="footer-column">
                    <h4>Product</h4>
                    <a href="#features">Features</a>
                    <a href="/pricing">Pricing</a>
                    <a href="#integrations">Integrations</a>
                </div>
                <div class="footer-column">
                    <h4>Company</h4>
                    <a href="#about">About</a>
                    <a href="#faq">FAQ</a>
                    <a href="/contact">Contact</a>
                </div>
                <div class="footer-column">
                    <h4>Legal</h4>
                    <a href="/privacy">Privacy</a>
                    <a href="/terms">Terms</a>
                </div>
            </div>
        </div>
        <div class="footer-bottom">
            <p>&copy; 2024 Flowmail. All rights reserved.</p>
        </div>
    </footer>
</div>

<style>
    /* Landing Page Specific Styles */
    .landing {
        overflow-x: hidden;
    }

    /* Hero Section */
    .hero {
        position: relative;
        min-height: 100vh;
        display: flex;
        align-items: center;
        padding: 8rem 2rem 4rem;
        overflow: hidden;
    }

    .hero-bg {
        position: absolute;
        inset: 0;
        background: var(--gradient-primary);
        z-index: -1;
    }

    .hero-bg::after {
        content: "";
        position: absolute;
        inset: 0;
        background: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.05'%3E%3Ccircle cx='30' cy='30' r='2'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
    }

    .hero-content {
        max-width: 600px;
        color: white;
        z-index: 1;
    }

    .hero-badge {
        display: inline-block;
        padding: var(--spacing-sm) var(--spacing-md);
        background: rgba(255, 255, 255, 0.15);
        border-radius: var(--radius-full);
        font-size: var(--font-size-sm);
        font-weight: 500;
        margin-bottom: var(--spacing-lg);
    }

    .hero h1 {
        font-size: var(--font-size-5xl);
        font-weight: 800;
        line-height: 1.1;
        margin-bottom: var(--spacing-lg);
    }

    .hero-gradient-text {
        background: var(--gradient-hero-text);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
    }

    .hero-subtitle {
        font-size: var(--font-size-xl);
        line-height: 1.6;
        opacity: 0.9;
        margin-bottom: var(--spacing-xl);
    }

    .hero-actions {
        display: flex;
        gap: var(--spacing-md);
        flex-wrap: wrap;
    }

    .hero-note {
        margin-top: var(--spacing-lg);
        font-size: var(--font-size-sm);
        opacity: 0.7;
    }

    .hero-visual {
        position: absolute;
        right: 5%;
        top: 50%;
        transform: translateY(-50%);
        z-index: 1;
    }

    /* App Preview Mockup */
    .app-preview {
        width: 380px;
        background: white;
        border-radius: var(--radius-xl);
        box-shadow: 0 25px 80px rgba(0, 0, 0, 0.3);
        overflow: hidden;
    }

    .preview-header {
        padding: var(--spacing-md);
        background: var(--color-bg-secondary);
        border-bottom: 1px solid var(--color-border);
        display: flex;
        align-items: center;
        gap: var(--spacing-md);
    }

    .preview-dots {
        display: flex;
        gap: 0.375rem;
    }

    .preview-dots span {
        width: 12px;
        height: 12px;
        border-radius: 50%;
        background: var(--color-border);
    }

    .preview-dots span:first-child {
        background: #ef4444;
    }
    .preview-dots span:nth-child(2) {
        background: #f59e0b;
    }
    .preview-dots span:last-child {
        background: #22c55e;
    }

    .preview-title {
        font-weight: 600;
        color: var(--color-text-secondary);
    }

    .preview-content {
        padding: var(--spacing-md);
    }

    .preview-email {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        padding: 0.875rem;
        border-radius: var(--radius-md);
        margin-bottom: var(--spacing-sm);
        transition: background var(--transition-base);
    }

    .preview-email:hover {
        background: var(--color-bg-secondary);
    }

    .preview-email.priority {
        background: #fef3c7;
    }

    .email-avatar {
        width: 40px;
        height: 40px;
        border-radius: 50%;
        background: var(--gradient-primary);
        color: white;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: var(--font-size-xs);
        font-weight: 600;
    }

    .email-body {
        flex: 1;
        min-width: 0;
    }

    .email-header {
        display: flex;
        align-items: center;
        gap: var(--spacing-sm);
        margin-bottom: var(--spacing-xs);
    }

    .email-sender {
        font-weight: 600;
        color: var(--color-text);
        font-size: 0.9rem;
    }

    .email-badge {
        font-size: 0.65rem;
        padding: 0.125rem 0.5rem;
        background: var(--color-warning);
        color: white;
        border-radius: var(--radius-full);
        font-weight: 600;
        text-transform: uppercase;
    }

    .email-subject {
        display: block;
        color: var(--color-text-muted);
        font-size: 0.85rem;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .preview-ai {
        display: flex;
        align-items: center;
        gap: var(--spacing-sm);
        padding: 0.75rem;
        background: var(--gradient-primary);
        color: white;
        border-radius: var(--radius-sm);
        font-size: 0.8rem;
        margin-top: var(--spacing-sm);
    }

    /* Features Section */
    .features-section {
        background: var(--color-bg-secondary);
    }

    .features-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
        gap: var(--spacing-xl);
    }

    .feature-card {
        background: white;
        padding: var(--spacing-xl);
        border-radius: var(--radius-xl);
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
    }

    .feature-icon {
        width: 48px;
        height: 48px;
        background: var(--gradient-primary);
        border-radius: var(--radius-lg);
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        margin-bottom: var(--spacing-lg);
    }

    .feature-card h3 {
        font-size: var(--font-size-xl);
        margin-bottom: 0.75rem;
    }

    .feature-card p {
        color: var(--color-text-muted);
        line-height: 1.6;
    }

    /* Integrations Section */
    .integrations-grid {
        display: flex;
        flex-wrap: wrap;
        justify-content: center;
        gap: var(--spacing-lg);
    }

    .integration-card {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 0.75rem;
        padding: var(--spacing-lg) var(--spacing-xl);
        background: var(--color-bg-secondary);
        border-radius: var(--radius-lg);
        transition:
            transform var(--transition-base),
            box-shadow var(--transition-base);
    }

    .integration-card:hover {
        transform: translateY(-2px);
        box-shadow: var(--shadow-sm);
    }

    .integration-logo {
        width: 56px;
        height: 56px;
        background: white;
        border-radius: var(--radius-lg);
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: var(--font-size-xl);
        font-weight: 700;
        color: var(--color-primary);
        box-shadow: var(--shadow-sm);
    }

    .integration-name {
        font-weight: 500;
        color: var(--color-text-secondary);
    }

    .integrations-note {
        text-align: center;
        margin-top: var(--spacing-xl);
        color: var(--color-text-light);
        font-size: 0.95rem;
    }

    /* About Section */
    .about-section {
        background: var(--color-bg-secondary);
    }

    .about-content {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: var(--spacing-3xl);
        align-items: center;
    }

    .about-text h2 {
        font-size: var(--font-size-4xl);
        margin-bottom: var(--spacing-lg);
    }

    .about-text p {
        color: var(--color-text-muted);
        line-height: 1.7;
        margin-bottom: var(--spacing-md);
    }

    .about-stats {
        display: flex;
        gap: var(--spacing-xl);
        margin-top: var(--spacing-xl);
    }

    .stat {
        display: flex;
        flex-direction: column;
    }

    .stat-value {
        font-size: var(--font-size-3xl);
        font-weight: 700;
        color: var(--color-primary);
    }

    .stat-label {
        font-size: var(--font-size-sm);
        color: var(--color-text-muted);
    }

    .about-visual {
        display: flex;
        justify-content: center;
    }

    .about-card {
        background: white;
        padding: var(--spacing-xl);
        border-radius: var(--radius-xl);
        box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08);
        max-width: 320px;
        text-align: center;
    }

    .about-card-icon {
        width: 64px;
        height: 64px;
        background: var(--gradient-primary);
        border-radius: var(--radius-xl);
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        margin: 0 auto var(--spacing-lg);
    }

    .about-card h3 {
        margin-bottom: 0.75rem;
        font-size: var(--font-size-xl);
    }

    .about-card p {
        color: var(--color-text-muted);
        line-height: 1.6;
    }

    /* FAQ Section */
    .faq-list {
        max-width: 700px;
        margin: 0 auto;
    }

    .faq-item {
        border-bottom: 1px solid var(--color-border);
    }

    .faq-question {
        width: 100%;
        padding: var(--spacing-lg) 0;
        background: none;
        border: none;
        display: flex;
        justify-content: space-between;
        align-items: center;
        cursor: pointer;
        text-align: left;
        font-family: var(--font-family);
        font-size: 1.05rem;
        font-weight: 500;
        color: var(--color-text);
    }

    .faq-question svg {
        transition: transform var(--transition-base);
        color: var(--color-text-light);
    }

    .faq-item.open .faq-question svg {
        transform: rotate(180deg);
    }

    .faq-answer {
        padding-bottom: var(--spacing-lg);
    }

    .faq-answer p {
        color: var(--color-text-muted);
        line-height: 1.7;
    }

    /* CTA Section */
    .cta-section {
        padding: var(--spacing-4xl) var(--spacing-xl);
        background: var(--gradient-primary);
    }

    .cta-content {
        text-align: center;
        color: white;
    }

    .cta-content h2 {
        font-size: var(--font-size-4xl);
        margin-bottom: var(--spacing-md);
    }

    .cta-content p {
        font-size: var(--font-size-lg);
        opacity: 0.9;
        margin-bottom: var(--spacing-xl);
    }

    /* Responsive */
    @media (max-width: 1024px) {
        .hero-visual {
            display: none;
        }

        .hero-content {
            max-width: 100%;
            text-align: center;
        }

        .hero-actions {
            justify-content: center;
        }

        .about-content {
            grid-template-columns: 1fr;
        }

        .about-visual {
            order: -1;
        }
    }

    @media (max-width: 768px) {
        .hero h1 {
            font-size: var(--font-size-4xl);
        }

        .about-stats {
            flex-wrap: wrap;
        }
    }

    @media (max-width: 480px) {
        .hero {
            padding: 6rem 1rem 3rem;
        }

        .hero h1 {
            font-size: var(--font-size-3xl);
        }

        .hero-actions {
            flex-direction: column;
        }

        .features-grid {
            grid-template-columns: 1fr;
        }
    }
</style>
