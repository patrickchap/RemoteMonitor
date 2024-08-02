/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'selector',
  content: ['./views/**/*.templ','./node_modules/flowbite/**/*.js',],
  theme: {
    extend: {},
    colors: {
      backgrounddark: '#1a202c',
    }
  },
    plugins: [
        require('flowbite/plugin')
    ]
}

