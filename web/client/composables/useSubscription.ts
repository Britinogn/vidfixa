import { useMutation, useQuery } from '@tanstack/vue-query'
import type { CheckoutResponse, PlanTier, SubscriptionResponse } from '~/types/api'

export function useSubscription() {
  const api = useApi()
  const auth = useAuthStore()
  const subscription = useQuery({
    queryKey: ['subscription'],
    queryFn: () => api<SubscriptionResponse>('/subscription/'),
    enabled: computed(() => auth.isAuthenticated),
  })
  const checkout = useMutation({
    mutationFn: (tier: Extract<PlanTier, 'plus' | 'pro'>) =>
      api<CheckoutResponse>('/subscription/checkout', { method: 'POST', body: { tier } }),
    onSuccess: ({ checkout_url }) => {
      window.location.assign(checkout_url)
    },
  })
  return { subscription, checkout }
}
