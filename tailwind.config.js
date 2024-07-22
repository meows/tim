/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    'ui/template/**/*.templ'
  ],
  darkMode: 'class', // or 'media' or 'class'
  theme: {
    extend: {
      fontFamily: ['Courier Prime', 'monospace'],
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
    require('@tailwindcss/line-clamp'),
    require('@tailwindcss/aspect-ratio'),
    require('daisyui'),
  ],
  corePlugins: { preflight: true }
}

