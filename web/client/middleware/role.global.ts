export default defineNuxtRouteMiddleware(async (to) => {
  if (import.meta.server) return

  const auth = useAuthStore()
  await auth.initialize()

  const isLogin = to.path === '/login' || to.path === '/register'
  const isDashboard = to.path === '/dashboard' || to.path.startsWith('/dashboard/')
  const isAdmin = to.path === '/admin' || to.path.startsWith('/admin/')
  const protectedPath = isDashboard || isAdmin

  if (protectedPath && !auth.user) {
    return navigateTo({ path: '/login', query: { redirect: to.fullPath } })
  }
  if (isAdmin && auth.user?.role !== 'admin') return navigateTo('/dashboard')
  if (isDashboard && auth.user?.role === 'admin') return navigateTo('/admin')
  if (isLogin && auth.user) return navigateTo(auth.user.role === 'admin' ? '/admin' : '/dashboard')
})
