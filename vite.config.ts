import path from "path";
import { execSync } from "node:child_process"


import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import type { PluginOption } from "vite";

const go = (): PluginOption => {
  return {
    name: 'go-build',
    enforce: 'pre',
    buildStart() {
      execSync("make build")
    },
  }
}

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [go(), react()],
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
