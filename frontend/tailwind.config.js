/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/**/*.{js,jsx,ts,tsx}',
    './src/components/admin/**/*.{js,jsx,ts,tsx}',
    './src/components/ui/**/*.{js,jsx,ts,tsx}',
    './public/index.html',
  ],
  darkMode: 'class',
  // Safelist important classes that might be dynamically generated
  safelist: [
    'bg-blue-50',
    'bg-green-50', 
    'bg-red-50',
    'bg-yellow-50',
    'bg-purple-50',
    'bg-amber-50',
    'text-blue-700',
    'text-green-700',
    'text-red-700',
    'text-yellow-700',
    'text-purple-700',
    'text-amber-700',
    'border-blue-200',
    'border-green-200',
    'border-red-200',
    'border-yellow-200',
    'border-purple-200',
    'border-amber-200',
    // Grid responsive classes
    'grid-cols-1',
    'sm:grid-cols-2',
    'md:grid-cols-2',
    'lg:grid-cols-4',
    'xl:grid-cols-4',
    // Responsive display classes
    'hidden',
    'sm:block',
    'md:block',
    'lg:block',
    'lg:flex',
    'lg:w-64',
    // Animation classes
    'animate-spin',
    'animate-pulse',
  ],
  theme: {
    extend: {
      colors: {
        // Shadcn/ui colors
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary))",
          foreground: "hsl(var(--secondary-foreground))",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive))",
          foreground: "hsl(var(--destructive-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
        },
        accent: {
          DEFAULT: "hsl(var(--accent))",
          foreground: "hsl(var(--accent-foreground))",
        },
        popover: {
          DEFAULT: "hsl(var(--popover))",
          foreground: "hsl(var(--popover-foreground))",
        },
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
    },
  },
  plugins: [],
}