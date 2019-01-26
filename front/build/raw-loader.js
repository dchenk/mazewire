module.exports = function rawLoader(source) {
	return "module.exports = " + JSON.stringify(source).replace(/\u2028/g, "\\u2028").replace(/\u2029/g, "\\u2029")
}
