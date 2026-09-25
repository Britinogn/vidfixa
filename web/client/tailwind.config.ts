import type { Config } from 'tailwindcss'

export default <Partial<Config>>{
  theme: {
    extend: {
      fontFamily: { sans: ['Inter', 'ui-sans-serif', 'system-ui'] },
      colors: {
        brand: { 50: '#fef2f2', 600: '#e53935', 700: '#c62828' },
      },
      boxShadow: { card: '0 8px 28px rgba(15, 23, 42, 0.045)' },
    },
  },
}
