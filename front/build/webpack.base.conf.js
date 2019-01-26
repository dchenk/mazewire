"use strict"
const path = require("path")
const utils = require("./utils")
const config = require("../config")
const vueLoaderConfig = require("./vue-loader.conf")

function resolve(dir) {
	return path.join(__dirname, "..", dir)
}

module.exports = {
	context: path.resolve(__dirname, "../"),
	entry: {
		app: "./src/main.js"
	},
	output: {
		path: config.build.assetsRoot,
		filename: "[name].js",
		publicPath: process.env.NODE_ENV === "production"
			? config.build.assetsPublicPath
			: config.dev.assetsPublicPath
	},
	resolve: {
		extensions: [".js", ".vue", ".json"],
		alias: {
			"@": resolve("src"),
		}
	},
	module: {
		rules: [
			// The eslint rule will either be inserted or result in nothing to evaluate at all.
			...(config.dev.useEslint ? [{
				test: /\.(js|vue)$/,
				loader: "eslint-loader",
				enforce: "pre",
				include: [resolve("src"), resolve("test")],
				options: {
					formatter: require("eslint-friendly-formatter"),
					emitWarning: !config.dev.showEslintErrorsInOverlay
				}
			}] : []),
			{
				test: /\.vue$/,
				loader: "vue-loader",
				options: vueLoaderConfig
			},
			{
				test: /\.js$/,
				loader: "babel-loader",
				include: [resolve("src"), resolve("test"), resolve("node_modules/webpack-dev-server/client")]
			},
			{
				test: /.*svg-data-url\/.+\.svg$/,
				loader: "url-loader"
			},
			{
				test: /.*svg-inline\/.+\.svg$/,
				use: [{
					loader: resolve("build/raw-loader.js")
				}]
			},
			{
				test: /\.(png|jpe?g|gif)(\?.*)?$/,
				loader: "url-loader",
				options: {
					limit: 4096, // Assets up to 4 KB in size are base64 encoded inline.
					name: utils.assetsPath("[name].[hash:7].[ext]")
				}
			},
			{
				// The special extension "rawcss" is used to load CSS for an iframe.
				test: /.*\.rawcss$/,
				loader: resolve("build/raw-loader.js")
			},
			{
				test: /\.(mp4|webm|ogg|mp3|wav|flac|aac)(\?.*)?$/,
				loader: "url-loader",
				options: {
					limit: 1024,
					name: utils.assetsPath("[name].[hash:7].[ext]")
				}
			},
			{
				test: /\.(woff2?|eot|ttf|otf)(\?.*)?$/,
				loader: "url-loader",
				options: {
					limit: 1024,
					name: utils.assetsPath("[name].[hash:7].[ext]")
				}
			}
		]
	},
	node: {
		// Prevent webpack from injecting useless setImmediate polyfill because Vue source contains it (although only uses it if it's native).
		setImmediate: false,
		// Prevent webpack from injecting mocks to Node native modules that does not make sense for the client.
		dgram: "empty",
		fs: "empty",
		net: "empty",
		tls: "empty",
		child_process: "empty"
	}
}
