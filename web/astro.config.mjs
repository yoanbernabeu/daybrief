import { defineConfig } from 'astro/config';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  site: 'https://yoanbernabeu.github.io',
  base: '/daybrief',
  output: 'static',
  vite: {
    plugins: [tailwindcss()]
  }
});
