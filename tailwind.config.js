/** @type {import('tailwindcss').Config} */
module.exports = {
    darkMode: 'class',
    content: [
        "./ui/templates/**/*.html",
        "./internal/**/*.go",
    ],
    theme: {
        extend: {
            colors: {
                primary: {
                    DEFAULT: '#FF8A00',
                    dark: '#E07900',
                },
                secondary: '#005D3A',
                background: {
                    light: '#F8F9FA',
                    dark: '#1C1C1E',
                },
                surface: {
                    light: '#FFFFFF',
                    dark: '#2C2C2E',
                },
                text: {
                    main: '#2C3E50',
                    sub: '#5A6C7D',
                },
                earth: {
                    clay: '#A0522D',
                    ochre: '#CC7722',
                    'ochre-light': '#E09540',
                    dark: '#1A120E',
                    sand: '#F4EBD0',
                    cream: '#FAF8F1',
                    accent: '#F58608'
                }
            },
            fontFamily: {
                display: ['"Inter"', 'sans-serif'],
                sans: ['"Inter"', 'sans-serif'],
                serif: ['"Playfair Display"', 'serif'],
            },
            borderRadius: {
                "DEFAULT": "12px",
                "lg": "20px",
                "xl": "32px",
                "full": "9999px"
            },
            boxShadow: {
                'soft': '0 4px 20px -2px rgba(0, 0, 0, 0.05)',
                'juicy': '0 8px 30px -4px rgba(255, 138, 0, 0.15)',
                'lifted': '0 16px 40px -8px rgba(255, 138, 0, 0.2)',
            }
        },
    },
    plugins: [
        require('@tailwindcss/forms'),
        require('@tailwindcss/container-queries'),
    ],
}
