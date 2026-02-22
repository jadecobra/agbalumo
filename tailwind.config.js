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
                "primary": "#FF8A00",           // ripe agbalumo skin
                "primary-variant": "#F57C00",
                "secondary": "#689F38",         // fresh leaf green
                "background-light": "#FFFBF5",  // warm cream background
                "background-dark": "#23160f",   // dark warm brown
                "surface-light": "#FFF8F0",     // creamy white flesh
                "surface-dark": "#2f221c",
                "text-main": "#3E2723",         // deep Nigerian earth brown
                "text-sub": "#6d4c41",
                "accent-star": "#C2185B",       // star-seed magenta
            },
            fontFamily: {
                "display": ["Lexend", "sans-serif"],
                "serif": ["Playfair Display", "serif"]
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
