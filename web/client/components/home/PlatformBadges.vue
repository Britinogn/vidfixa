<template>
  <div class="flex flex-wrap items-center gap-x-4 gap-y-3 text-xs font-medium text-slate-500 sm:gap-x-7 sm:text-sm">
    <span class="flex items-center gap-2"><i class="platform-dot instagram"> </i><span>Instagram</span></span>
    <span class="flex items-center gap-2"><i class="platform-dot facebook"> </i><span>Facebook</span></span>
    <span class="flex items-center gap-2"><i class="platform-dot x-platform">X</i><span>X</span></span>
    <span class="flex items-center gap-2"><i class="platform-dot linkedin">in</i><span>LinkedIn</span></span>
  </div>
</template>

<style scoped>
.platform-dot { display: inline-grid; width: 19px; height: 19px; flex-shrink: 0; place-items: center; border-radius: 5px; color: white; font-size: 10px; font-style: normal; font-weight: 700; }
.instagram { background: linear-gradient(135deg, #7c3aed, #f43f5e, #f59e0b); }
.facebook { background: #1877f2; border-radius: 50%; }
.x-platform { background: #0f172a; }
.linkedin { background: #0a66c2; border-radius: 3px; font-size: 11px; }
</style>