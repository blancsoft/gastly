import path from "path";

import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: [
      { find: '@', replacement: path.resolve(__dirname, 'src') },
    ],
  },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ["./src/setupTests.ts"],
  },
  build: {
    rollupOptions: {
      output: {
        assetFileNames: ({ name }) => {
          if (name.endsWith("wasm.br")) {
            return "assets/[name][extname]"
          }

          return "assets/[name]-[hash][extname]"
        }
      }
    }
  }
})
