<template>
  <div>
    <UPageHero
      title="Email, reimagined."
      description="Connect your Gmail, Outlook, or iCloud accounts to a unified inbox powered by AI. Smart filters, multiple inboxes, and intelligent agents work together to make email a pleasure again."
      :links="[{
        label: 'Get early access',
        to: '/pricing',
        trailingIcon: 'i-lucide-arrow-right',
        size: 'xl'
      }, {
        label: 'See how it works',
        to: '#features',
        icon: 'i-lucide-play-circle',
        size: 'xl',
        color: 'neutral',
        variant: 'subtle'
      }]"
    >
      <template #headline>
        <UBadge color="primary" variant="subtle" size="lg">
          Now in private beta
        </UBadge>
      </template>
    </UPageHero>

    <UPageSection
      id="integrations"
      title="Connect all your accounts"
      description="Seamlessly sync with the email providers you already use. One inbox to rule them all."
    >
      <div class="flex justify-center gap-12 flex-wrap">
        <div class="flex flex-col items-center gap-2">
          <UIcon name="i-simple-icons-gmail" class="w-12 h-12 text-red-500" />
          <span class="text-sm text-gray-500 dark:text-gray-400">Gmail</span>
        </div>
        <div class="flex flex-col items-center gap-2">
          <UIcon name="i-simple-icons-microsoftoutlook" class="w-12 h-12 text-blue-500" />
          <span class="text-sm text-gray-500 dark:text-gray-400">Outlook</span>
        </div>
        <div class="flex flex-col items-center gap-2">
          <UIcon name="i-simple-icons-apple" class="w-12 h-12 text-gray-700 dark:text-gray-300" />
          <span class="text-sm text-gray-500 dark:text-gray-400">iCloud</span>
        </div>
        <div class="flex flex-col items-center gap-2">
          <UIcon name="i-lucide-mail" class="w-12 h-12 text-purple-500" />
          <span class="text-sm text-gray-500 dark:text-gray-400">IMAP/SMTP</span>
        </div>
      </div>
    </UPageSection>

    <UPageSection
      id="features"
      title="Finally, email that works for you"
      description="We've rethought every aspect of email to help you focus on what matters and let AI handle the rest."
      :features="[{
        icon: 'i-lucide-brain',
        title: 'AI-powered triage',
        description: 'Our intelligent agents automatically categorize, prioritize, and surface what needs your attention. Never miss an important email again.'
      }, {
        icon: 'i-lucide-inbox',
        title: 'Multiple smart inboxes',
        description: 'Create custom inboxes for projects, clients, or topics. Emails flow automatically to where they belong based on your rules and AI learning.'
      }, {
        icon: 'i-lucide-filter',
        title: 'Smart filters that learn',
        description: 'Filters that adapt to your behavior. The more you use them, the smarter they get. No more manual rule creation.'
      }, {
        icon: 'i-lucide-bot',
        title: 'AI agents at your service',
        description: 'Delegate tasks to AI agents: draft replies, schedule follow-ups, summarize threads, or extract action items automatically.'
      }, {
        icon: 'i-lucide-bell-off',
        title: 'Intelligent notifications',
        description: 'Only get notified for what truly matters. Our AI understands urgency and importance, so you can stay focused.'
      }, {
        icon: 'i-lucide-zap',
        title: 'Instant actions',
        description: 'One-click actions powered by AI: archive newsletters, snooze non-urgent items, or batch process similar emails effortlessly.'
      }]"
    />

    <UPageSection
      title="Your AI email assistant"
      description="Meet your new inbox companion. It reads, understands, and helps you manage your email like never before."
    >
      <UCard class="max-w-3xl mx-auto">
        <div class="space-y-4">
          <div class="flex items-start gap-4">
            <UIcon name="i-lucide-sparkles" class="w-6 h-6 text-primary-500 mt-1 shrink-0" />
            <div>
              <p class="font-medium">Summarize long threads</p>
              <p class="text-sm text-gray-500 dark:text-gray-400">Get the key points from 50-email threads in seconds</p>
            </div>
          </div>
          <div class="flex items-start gap-4">
            <UIcon name="i-lucide-pen-line" class="w-6 h-6 text-primary-500 mt-1 shrink-0" />
            <div>
              <p class="font-medium">Draft contextual replies</p>
              <p class="text-sm text-gray-500 dark:text-gray-400">AI writes drafts that sound like you, based on your style</p>
            </div>
          </div>
          <div class="flex items-start gap-4">
            <UIcon name="i-lucide-calendar-check" class="w-6 h-6 text-primary-500 mt-1 shrink-0" />
            <div>
              <p class="font-medium">Extract action items</p>
              <p class="text-sm text-gray-500 dark:text-gray-400">Automatically create tasks and calendar events from emails</p>
            </div>
          </div>
          <div class="flex items-start gap-4">
            <UIcon name="i-lucide-search" class="w-6 h-6 text-primary-500 mt-1 shrink-0" />
            <div>
              <p class="font-medium">Semantic search</p>
              <p class="text-sm text-gray-500 dark:text-gray-400">Find emails by meaning, not just keywords</p>
            </div>
          </div>
        </div>
      </UCard>
    </UPageSection>

    <UPageSection
      title="How it works"
      description="Get started in three simple steps"
    >
      <div class="grid md:grid-cols-3 gap-8 max-w-4xl mx-auto">
        <div class="text-center">
          <div class="w-12 h-12 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center mx-auto mb-4">
            <span class="text-xl font-bold text-primary-600 dark:text-primary-400">1</span>
          </div>
          <h3 class="font-semibold mb-2">Connect your accounts</h3>
          <p class="text-sm text-gray-500 dark:text-gray-400">Link your Gmail, Outlook, iCloud, or any IMAP account securely</p>
        </div>
        <div class="text-center">
          <div class="w-12 h-12 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center mx-auto mb-4">
            <span class="text-xl font-bold text-primary-600 dark:text-primary-400">2</span>
          </div>
          <h3 class="font-semibold mb-2">AI learns your patterns</h3>
          <p class="text-sm text-gray-500 dark:text-gray-400">Our AI analyzes your email habits and preferences to personalize your experience</p>
        </div>
        <div class="text-center">
          <div class="w-12 h-12 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center mx-auto mb-4">
            <span class="text-xl font-bold text-primary-600 dark:text-primary-400">3</span>
          </div>
          <h3 class="font-semibold mb-2">Enjoy inbox zero</h3>
          <p class="text-sm text-gray-500 dark:text-gray-400">Watch as your inbox transforms into an organized, stress-free space</p>
        </div>
      </div>
    </UPageSection>

    <UPageSection>
      <UPageCTA
        title="Ready to transform your inbox?"
        description="Join thousands on the waitlist for early access. Be among the first to experience email that actually works for you."
        variant="subtle"
        :links="[{
          label: 'Join the waitlist',
          to: '/pricing',
          trailingIcon: 'i-lucide-arrow-right',
          color: 'neutral'
        }]"
      />
    </UPageSection>
  </div>
</template>
