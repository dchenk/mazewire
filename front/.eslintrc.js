// https://eslint.org/docs/user-guide/configuring

module.exports = {
	root: true,
	env: {browser: true},

	// Docs: https://github.com/vuejs/eslint-plugin-vue
	extends: ["plugin:vue/strongly-recommended"],

	// Lint *.vue files.
	plugins: ["vue"],

	// Check if imports actually resolve.
	settings: {
		"import/resolver": {
			webpack: {config: "build/webpack.base.conf.js"}
		}
	},

	// Custom rules here.
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
