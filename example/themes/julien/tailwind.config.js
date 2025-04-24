/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./views/**/*.html", "./layouts/**/*.html"],
    theme: {
        extend: {},
    },
    plugins: [
        require('@tailwindcss/typography')
    ],
}

