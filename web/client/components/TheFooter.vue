<template>
  <footer class="border-t border-slate-200 bg-white">
    <div class="mx-auto w-full max-w-7xl px-4 py-8 sm:px-5 md:px-8">
      <div
        class="flex flex-col gap-8 md:flex-row md:items-start md:justify-between"
      >
        <!-- Brand -->
        <div class="max-w-xs">
          <NuxtLink to="/" aria-label="VidFixa home">
            <BrandMark />
          </NuxtLink>

          <p class="mt-3 text-sm leading-6 text-slate-500">
            A simple way to save supported videos from the web.
          </p>
        </div>

        <!-- Navigation -->
        <div class="grid grid-cols-2 gap-x-10 gap-y-8 sm:grid-cols-4">
          <!-- Product -->
          <div>
            <h2
              class="text-xs font-semibold uppercase tracking-wider text-slate-900"
            >
              Product
            </h2>

            <nav class="mt-3 flex flex-col gap-2.5 text-sm">
              <NuxtLink
                to="/#features"
                class="text-slate-500 transition hover:text-slate-900"
              >
                Features
              </NuxtLink>

              <NuxtLink
                to="/#pricing"
                class="text-slate-500 transition hover:text-slate-900"
              >
                Pricing
              </NuxtLink>

              <NuxtLink
                to="/#how-it-works"
                class="text-slate-500 transition hover:text-slate-900"
              >
                How it works
              </NuxtLink>
            </nav>
          </div>

          <!-- Company -->
          <div>
            <h2
              class="text-xs font-semibold uppercase tracking-wider text-slate-900"
            >
              Company
            </h2>

            <nav class="mt-3 flex flex-col gap-2.5 text-sm">
              <NuxtLink
                to="/about"
                class="text-slate-500 transition hover:text-slate-900"
              >
                About
              </NuxtLink>

              <!-- <NuxtLink
                to="/contact"
                class="text-slate-500 transition hover:text-slate-900"
              >
                Contact
              </NuxtLink> -->
              <a
                href="mailto:brightonwuemeri@gmail.com"
                class="text-slate-500 transition hover:text-slate-900"
              >
                Contact
              </a>
            </nav>
          </div>

          <!-- Legal -->
          <!-- <div>
            <h2
              class="text-xs font-semibold uppercase tracking-wider text-slate-900"
            >
              Legal
            </h2>

            <nav class="mt-3 flex flex-col gap-2.5 text-sm">
              <NuxtLink
                to="/privacy"
                class="text-slate-500 transition hover:text-slate-900"
              >
                Privacy Policy
              </NuxtLink>

              <NuxtLink
                to="/terms"
                class="text-slate-500 transition hover:text-slate-900"
              >
                Terms of Service
              </NuxtLink>
            </nav>
          </div> -->

          <!-- Builder -->
          <div>
            <h2
              class="text-xs font-semibold uppercase tracking-wider text-slate-900"
            >
              Builder
            </h2>

            <nav class="mt-3 flex flex-col gap-2.5 text-sm">
              <NuxtLink
                to="/about"
                class="text-slate-500 transition hover:text-slate-900"
              >
                Tino Ctemz
              </NuxtLink>

              <a
                href="https://x.com/Britinogn"
                target="_blank"
                rel="noopener noreferrer"
                class="text-slate-500 transition hover:text-slate-900"
              >
                X / @Britinogn
              </a>
            </nav>
          </div>
        </div>
      </div>

      <!-- Bottom -->
      <div
        class="mt-8 flex flex-col gap-2 border-t border-slate-100 pt-6 text-xs text-slate-400 sm:flex-row sm:items-center sm:justify-between"
      >
        <span>
          © {{ new Date().getFullYear() }} VidFixa. All rights reserved.
        </span>

        <span>
          Built by
          <NuxtLink
            to="/about"
            class="font-medium text-slate-600 transition hover:text-brand-600"
          >
            Tino Ctemz
          </NuxtLink>
        </span>
      </div>
    </div>
  </footer>
</template>