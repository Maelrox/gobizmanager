/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#f0f7ff',
          100: '#e6f2ff',
          200: '#c9e2ff',
          300: '#a5d0ff',
          400: '#7eb8ff',
          500: '#5a9eff',
          600: '#3d84ff',
          700: '#2b6bff',
          800: '#1a52ff',
          900: '#0a3aff',
        },
        secondary: {
          50: '#f5f7fa',
          100: '#e9edf2',
          200: '#d1d9e6',
          300: '#b8c4d9',
          400: '#9fafcc',
          500: '#869abf',
          600: '#6d85b2',
          700: '#5470a5',
          800: '#3b5b98',
          900: '#22468b',
        },
        accent: {
          50: '#fef2f2',
          100: '#fee2e2',
          200: '#fecaca',
          300: '#fca5a5',
          400: '#f87171',
          500: '#ef4444',
          600: '#dc2626',
          700: '#b91c1c',
          800: '#991b1b',
          900: '#7f1d1d',
        },
        neutral: {
          50: '#fafafa',
          100: '#f5f5f5',
          200: '#e5e5e5',
          300: '#d4d4d4',
          400: '#a3a3a3',
          500: '#737373',
          600: '#525252',
          700: '#404040',
          800: '#262626',
          900: '#171717',
        },
      },
      fontFamily: {
        sans: ['Inter var', 'system-ui', 'sans-serif'],
        logo: ['"Plus Jakarta Sans"', 'sans-serif'],
      },
      boxShadow: {
        'soft': '0 2px 15px -3px rgba(0, 0, 0, 0.07), 0 10px 20px -2px rgba(0, 0, 0, 0.04)',
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
} 