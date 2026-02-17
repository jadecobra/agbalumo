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
                "primary": "#FF5E0E", // Agbalumo Orange
                "secondary": "#2D5A27", // Palm Leaf Green
                "background-light": "#FFF2EB", // Pale Agbalumo Orange Tint
                "background-dark": "#23160f", // Dark Warm Brown
                "surface-light": "#ffffff",
                "surface-dark": "#2f221c",
                "text-main": "#181310",
                "text-sub": "#6d4c41",
            },
            fontFamily: {
                "display": ["Lexend", "sans-serif"]
            },
            borderRadius: {
                "DEFAULT": "0.5rem",
                "lg": "1rem",
                "xl": "1.5rem",
                "full": "9999px"
            },
            boxShadow: {
                'soft': '0 4px 20px -2px rgba(0, 0, 0, 0.05)',
            }
        },
    },
    plugins: [
        require('@tailwindcss/forms'),
        require('@tailwindcss/container-queries'),
    ],
}
