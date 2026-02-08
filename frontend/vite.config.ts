import { defineConfig, loadEnv, type UserConfigExport } from "vite";
import react from "@vitejs/plugin-react";
import tsconfigPaths from "vite-tsconfig-paths";
import { devtools } from "@tanstack/devtools-vite";

// https://vite.dev/config/
export default defineConfig(({ command, mode }) => {
  const env = loadEnv(mode, process.cwd());
  const commonConfig: UserConfigExport = {
    plugins: [
      devtools(),
      react({
        babel: {
          plugins: [["babel-plugin-react-compiler"]],
        },
      }),
      tsconfigPaths(),
    ],
    base: "/",
  };

  if (command === "serve") {
    return {
      ...commonConfig,
      server: {
        open: false,
        port: Number(env.VITE_PORT),
        strictPort: true,
        proxy: {
          "/api": {
            target: env.VITE_API_BASE_URL,
            changeOrigin: true,
            rewrite: (path) => {
              return path.replace(/^\/api/, "");
            },
          },
        },
      },
    };
  }

  return {
    ...commonConfig,
  };
});
