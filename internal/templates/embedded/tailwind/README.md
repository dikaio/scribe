# Scribe Site with Tailwind CSS 4.1

This site was created using Scribe static site generator with modern Tailwind CSS 4.1.

## Tailwind CSS Approach

This site uses a utility-first CSS approach with Tailwind CSS:

- All styling is done with Tailwind utility classes directly in HTML templates
- Built-in prose classes for beautiful typography (no plugin needed)
- No custom CSS needed in most cases
- Rapid development with composable utility classes
- Optimized build process produces minimal CSS

## Development

To start developing your site:

### Initial Setup

1. Install dependencies:
   ```
   npm install
   ```

2. Start the Tailwind CSS compiler:
   ```
   npm run dev
   ```

3. In another terminal, start the Scribe development server:
   ```
   scribe serve
   ```

4. View your site at http://localhost:8080

Or for convenience, use the all-in-one command:
   ```
   scribe run
   ```

### Building for Production

To build your site for production:

1. Compile and minify the CSS:
   ```
   npm run build
   ```

2. Build the site:
   ```
   scribe build
   ```

The built site will be in the `public` directory.

### Modern Tailwind CSS Setup

This site uses the modern Tailwind CSS 4.1 approach:
- Simplified setup with just `@import "tailwindcss"` in the CSS
- Comes with both tailwindcss and @tailwindcss/cli for processing
- Built-in prose classes for beautiful typography (no plugin needed)
- No configuration files needed (no tailwind.config.js)
- Customization available through `@layer` directives if needed
- Focus on using utility classes in HTML rather than writing custom CSS