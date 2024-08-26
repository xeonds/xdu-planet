import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from "unplugin-auto-import/vite";
import Components from "unplugin-vue-components/vite";
import { ElementPlusResolver } from "unplugin-vue-components/resolvers";
import path from 'path'

const pathSrc = path.resolve(__dirname, 'src')
export default defineConfig({
  plugins: [
    vue(),
    AutoImport({
      resolvers: [ElementPlusResolver()],
    }),
    Components({
      resolvers: [ElementPlusResolver()],
    }),
  ],
      resolve: {
        alias: {
            '@': pathSrc,
        },
    },
  server: {
    proxy: {
      "/api/v1/": {
        target: "http://localhost:8192/",
        changeOrigin: true,
      },
      "/db/": {
        target: "http://localhost:8192/",
        changeOrigin: true,
      }
    },
  },
  base: "./",
})
