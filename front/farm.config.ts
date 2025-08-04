import path from "node:path";
import { defineConfig } from "@farmfe/core";
import postcss from "@farmfe/js-plugin-postcss";

export default defineConfig({
	plugins: [
		[
			"@farmfe/plugin-react",
			{
				runtime: "automatic",
			},
		],
		postcss(),
	],
	compilation: {
		resolve: {
			alias: {
				"@": path.resolve(__dirname, "./src"),
			},
		},
	},
	server: {
		cors: true,
		proxy: {
			"/api/anime": {
				target: "https://api.jikan.moe/v4",
				changeOrigin: true,
				rewrite: (path) => path.replace(/^\/api/, ""),
			},
		},
	},
});
