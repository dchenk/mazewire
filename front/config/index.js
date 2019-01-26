"use strict"
// Template version: 1.3.1
// see http://vuejs-templates.github.io/webpack for documentation.

const path = require("path")

module.exports = {
	dev: {
		assetsSubDirectory: "src/assets",
		assetsPublicPath: "/",
		proxyTable: {},

		// Development server settings
		host: "localhost",
		port: 8082, // can be overwritten by process.env.PORT, if port is in use, a free one will be determined
		autoOpenBrowser: false,
		errorOverlay: true,
		notifyOnErrors: true,
		poll: false, // https://webpack.js.org/configuration/dev-server/#devserver-watchoptions-

		// If true, code will be linted during bundling, and linting errors and warnings will be shown in the console.
		useEslint: false,
		// If useEslint above is true and this is true, then eslint errors and warnings will also be shown in an overlay
		// in the browser.
		showEslintErrorsInOverlay: false,

		devtool: "cheap-module-eval-source-map", // https://webpack.js.org/configuration/devtool/#development

		// If you have problems debugging vue-files in devtools, it may help to set this to false.
		// https://vue-loader.vuejs.org/en/options.html#cachebusting
		cacheBusting: true,

		cssSourceMap: true
	},
	build: {
		index: path.resolve(__dirname, "../dist/index.html"),
		assetsRoot: path.resolve(__dirname, "../dist"),
		assetsSubDirectory: "",
		assetsPublicPath: "https://cdn.mazewire.com/",
		productionSourceMap: false,

		devtool: "#source-map", // https://webpack.js.org/configuration/devtool/#production

		// Gzip off by default as many popular CDNs already gzip all static assets for you.
		// Before enabling, make sure to: npm install --save-dev compression-webpack-plugin
		productionGzip: false,
		productionGzipExtensions: ["js", "css"],

		// Run the build command with an extra argument to view the bundle analyzer report after build finishes:
		// `npm run build --report`
		// Set to `true` or `false` to always turn it on or off.
		bundleAnalyzerReport: process.env.npm_config_report
	}
}
