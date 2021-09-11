import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import VitePluginElementPlus from 'vite-plugin-element-plus'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  return {
    plugins: [
      vue(),
      VitePluginElementPlus({
        // if you need to use the *.scss source file, you need to uncomment this comment
        // useSource: true
        format: mode === 'development' ? 'esm' : 'cjs',
      }),
    ],
  }
})