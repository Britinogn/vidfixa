export type UserRole = 'user' | 'admin'
export type PlanTier = 'free' | 'plus' | 'pro'
export type DownloadStatus = 'queued' | 'processing' | 'completed' | 'failed' | 'cancelled'

export interface User {
  id: string
  full_name: string
  email: string
  role: UserRole
  created_at: string
}

export interface AuthResponse {
  user: User
  token?: string
}

export interface Usage {
  used: number
  limit: number
  remaining: number
}

export interface DashboardResponse {
  user: User
  plan: PlanTier
  usage: Usage
}

export interface Download {
  id: string
  user_id?: string
  anon_id?: string
  url: string
  platform: string
  status: DownloadStatus
  error?: string
  created_at: string
  completed_at?: string
}

export interface SubscriptionResponse {
  plan: PlanTier
}

export interface CheckoutResponse {
  checkout_url: string
}

export interface AdminOverview {
  total_users: number
  active_subscriptions: number
  total_downloads: number
  failed_downloads: number
  successful_payments: number
}

export interface AdminUser extends User {}

export interface AdminSubscription {
  id: string
  user_id: string
  user_email: string
  plan: PlanTier
  status: string
  provider_subscription_id?: string
  started_at?: string
  expires_at?: string
  created_at: string
}

export interface AdminPayment {
  id: string
  user_id: string
  user_email: string
  subscription_id?: string
  provider: string
  provider_reference: string
  amount: string
  currency: string
  status: string
  created_at: string
}

export interface AdminDownload {
  id: string
  user_id?: string
  user_email?: string
  anon_id?: string
  url: string
  platform: string
  status: DownloadStatus
  error?: string
  created_at: string
  completed_at?: string
}
