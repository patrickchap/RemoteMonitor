/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'selector',
  content: ['./views/**/*.templ', './node_modules/flowbite/**/*.js'],
  theme: {
    extend: {},
    colors: {
      backgrounddark: '#1a202c',
      'background-dark': '#131F39',
      'surface-dark': '#223459',
    }
  },
  plugins: [
    require('flowbite/plugin')
  ]
}

