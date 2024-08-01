/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./**/*.{html,js,templ}", "./node_modules/flowbite/**/*.js"],
    theme: {
        extend: {
            colors: {
                background: '#FFFFFF',
                backgrounddark: '#1f2937',
                primary: '#6200EE',
                primaryDark: '#BB86FC',
                primaryVariant: '#3700B3',
                primaryVariantDark: '#3700B3',
                secondary: '#03DAC6',
                secondarydark: '#03DAC6',
                surface: '#FFF',
                surfaceDark: '#071952',
                error: '#B00020',
                errorDark: '#CF6679',
                onError: '#FFF',
                onErrorDark: '#000',
                onPrimary: '#FFF',
                onPrimaryDark: '#000',
                onSecondary: '#000',
                onSecondaryDark: '#000',
                onBackground: '#000',
                onBackgroundDark: '#FFF',
                onSurface: '#000',
                onSurfaceDark: '#FFF'
            },
        },
        plugins: [
            require('flowbite/plugin')
        ],
        darkMode: 'media',
        important: true,
    }
}

