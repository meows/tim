/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    'ui/template/**/*.templ',
    'cmd/web/admin_handlers.go',
  ],
  darkMode: 'class',
  theme: {
    extend: {
      fontFamily: {
        sans: ['Courier Prime', 'monospace'],
      },
      colors: {
        lavender1: '#a8a3c0',
        lavender2: '#a9a4c1',
        lavender3: '#a9a5bf',
        lavender4: '#a7a2bd',
        lavender5: '#a7a2bf',
        goBlue: '#5999b6',
      },
      height: {
        '2/3-screen': '66.67vh',
      },
      maxHeight: {
        '2/3-screen': '66.67vh',
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
    require('@tailwindcss/line-clamp'),
    require('@tailwindcss/aspect-ratio'),
    require('daisyui'),
  ],
  corePlugins: { preflight: true },
  safelist: [
    'max-w-screen-sm',
    'max-w-1/3',
    'mx-auto',
    'mb-6',
    'error-message',
    'container',
  ],
  daisyui: {
    themes: [
      "nord",
      {
        mytheme: {
          "primary": "#4a6572",
          "secondary": "#38b2ac",
          "accent": "#5999b6",
          "neutral": "#f5f5f5",
          "base-100": "#a9a5bf",
          "info": "#3b82f6",
          "success": "#10b981",
          "warning": "#f59e0b",
          "error": "#ef4444",
        },
      },
    ],
  },
}
