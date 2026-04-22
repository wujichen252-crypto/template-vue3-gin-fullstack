/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: '#000000',
        background: '#FFFFFF',
        surface: '#F5F5F5',
        accent: '#2563EB'
      }
    },
  },
  plugins: [],
}
