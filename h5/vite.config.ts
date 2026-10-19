import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { fileURLToPath, URL } from 'node:url';

export default defineConfig({
  base: './',
  plugins: [react()],
  resolve: {
    alias: [
      { find: /^@ds(?=\/|$)/, replacement: fileURLToPath(new URL('./src/ds', import.meta.url)) },
      { find: /^@(?=\/)/, replacement: fileURLToPath(new URL('./src', import.meta.url)) },
    ],
  },
});
