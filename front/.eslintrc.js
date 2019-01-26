// https://eslint.org/docs/user-guide/configuring

module.exports = {
	root: true,
	parserOptions: {
		parser: "babel-eslint"
	},
	env: { browser: true },
	// Docs: https://github.com/vuejs/eslint-plugin-vue
	// Use plugin:vue/strongly-recommended or plugin:vue/recommended for stricter rules than with plugin:vue/essential.
	extends: ["plugin:vue/strongly-recommended"],
	plugins: ["vue"], // Lint *.vue files.
	// Check if imports actually resolve.
	settings: {
		"import/resolver": {
			webpack: { config: "build/webpack.base.conf.js" }
		}
	},
	// Add custom rules here.
	rules: {
		"quotes": ["error", "double"],
		// Vue-specific rules.
		"vue/max-attributes-per-line": "off",
		"vue/html-indent": ["error", "tab", {"attribute": 1}],
		"vue/require-v-for-key": "warn",
		"vue/order-in-components": "warn",
		// Allow debugger during development.
		"no-debugger": process.env.NODE_ENV === "production" ? "error" : "off"
	}
}
