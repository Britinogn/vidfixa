<template>
  <ul class="flex flex-wrap items-center gap-x-4 gap-y-3 text-xs font-medium text-slate-500 sm:gap-x-7 sm:text-sm">
    <li>
      <a
        href="https://instagram.com"
        target="_blank"
        rel="noopener noreferrer"
        class="group flex items-center gap-2 transition hover:text-slate-800"
      >
        <i class="platform-dot instagram transition group-hover:scale-105">IG</i>
        <span>Instagram</span>
      </a>
    </li>
    <li>
      <a
        href="https://facebook.com"
        target="_blank"
        rel="noopener noreferrer"
        class="group flex items-center gap-2 transition hover:text-slate-800"
      >
        <i class="platform-dot facebook transition group-hover:scale-105">f</i>
        <span>Facebook</span>
      </a>
    </li>
    <li>
      <a
        href="https://x.com"
        target="_blank"
        rel="noopener noreferrer"
        class="group flex items-center gap-2 transition hover:text-slate-800"
      >
        <i class="platform-dot x-platform transition group-hover:scale-105">X</i>
        <span>𝕏</span>
      </a>
    </li>
    <li>
      <a
        href="https://linkedin.com"
        target="_blank"
        rel="noopener noreferrer"
        class="group flex items-center gap-2 transition hover:text-slate-800"
      >
        <i class="platform-dot linkedin transition group-hover:scale-105">in</i>
        <span>LinkedIn</span>
      </a>
    </li>
    <li>
      <a
        href="https://tiktok.com"
        target="_blank"
        rel="noopener noreferrer"
        class="group flex items-center gap-2 transition hover:text-slate-800"
      >
        <i class="platform-dot tiktok transition group-hover:scale-105">TT</i>
        <span>TikTok</span>
      </a>
    </li>
  </ul>
</template>

<style scoped>
.platform-dot {
  display: inline-grid;
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  place-items: center;
  border-radius: 5px;
  color: white;
  font-size: 9px;
  font-style: normal;
  font-weight: 700;
  letter-spacing: -0.02em;
}
.instagram { background: linear-gradient(135deg, #7c3aed, #f43f5e, #f59e0b); }
.facebook { background: #1877f2; border-radius: 50%; font-size: 12px; font-family: Georgia, serif; }
.x-platform { background: #0f172a; }
.linkedin { background: #0a66c2; border-radius: 3px; }
.tiktok { background: #010101; }
</style>