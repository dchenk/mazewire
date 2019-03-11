import Vue from "vue"
import App from "./App.vue"
import router from "./router.ts"
import msgpack from "msgpack-lite/lib/browser.js"
import Toasted from "vue-toasted"
import {hexEncode} from "./hex-encode"

import "./style.css"

import "dynamic-tabs/dist/tabs.css"
import "dynamic-textfields/textfields.css"

import "@material/select/dist/mdc.select.min.css"
import "@material/drawer/dist/mdc.drawer.min.css"
import "@material/list/dist/mdc.list.min.css"

import VueCodemirror from "vue-codemirror"
import "codemirror/lib/codemirror.css"
import "codemirror/addon/selection/active-line.js"
import "codemirror/mode/css/css.js"
import "codemirror/mode/javascript/javascript.js"
import "codemirror/mode/htmlmixed/htmlmixed.js"

import "./editor-theme.css"

import {dateExtPrefix, msgpPackDate, msgpUnpackDate} from "./msgp-date.js"
// Files with the ".rawcss" extension may belong either just to the page editing iframe component or also
// also in the main app. The following mainAndIframeCSS array lists the CSS files that should be used in
// the main app. These files are not pre-processed by an CSS tool (TODO: should they be?).
import layoutsCSS from "./layouts.rawcss"

const MessagePackCodec = msgpack.createCodec({
	uint8array: true
})

// Add the time.Time encoding used in github.com/dchenk/msgp.
MessagePackCodec.addExtUnpacker(dateExtPrefix, msgpUnpackDate)
MessagePackCodec.addExtPacker(dateExtPrefix, Date, msgpPackDate)


const mainAndIframeCSS = [layoutsCSS]

mainAndIframeCSS.forEach(function(elem) {
	const sc = document.createElement("STYLE")
	sc.innerHTML = elem
	document.head.appendChild(sc)
})

Vue.use(VueCodemirror, {
	options: {
		theme: "ww",
		indentUnit: 4,
		indentWithTabs: true,
		lineNumbers: true,
		styleActiveLine: true
	}
})

Vue.config.productionTip = false

Vue.use(Toasted, {duration: 3800}) // Toast-style notifications plugin.

Vue.prototype.$showDialog = function(dialogSetup) {
	this.$root.$children[0].dialogSetup = dialogSetup
	this.$root.$children[0].dialogActive = true
}

Vue.prototype.$hideDialog = function() {
	this.$root.$children[0].dialogActive = false
}

Vue.prototype.$showErrDialog = function(errMessage) {
	this.$showDialog({
		title: "An error occurred",
		acceptText: "OK",
		rejectText: "N/A",
		innerHTML: errMessage
	})
}

const apiPath = process.env === "development" ? "http://localhost:8080/api/" : "/api/"

// Set up the active site configuration for requests, to allow requests from one site to be made about another.
// Upon initial loading, the vue object this.activeSite will be set from window.activeSite.
Vue.prototype.activeSite = {
	id: 0, // the ID of the site the user is editing; to the server the default 0 means current host
	domain: window.location.host,
	name: "",
	logo: "",
	favicon: "",
	role: ""
}

Vue.prototype.$showLoader = function() {
	this.$root.$children[0].loading = true
}

Vue.prototype.$hideLoader = function() {
	this.$root.$children[0].loading = false
}

// doingAuthStep indicates whether we are in the processing of authentication, triggered by a request
// that requires the user to be authenticated.
let doingAuthStep = false

// $req makes an API request. The handleGood argument should be a function that handles the decoded response.
// The handleError argument is optional: authentication and server errors and alerts are handled by $req.
Vue.prototype.$req = function(method, endpoint, data, handleGood, handleError = () => {
}) {

	this.$showLoader()

	const opts = {
		method: method,
		credentials: "include",
		cache: "no-store"
	}

	if (data !== null) {
		// Add the API parameter of the site the user is currently editing.
		if (this.activeSite.id !== 0) {
			data.site = this.activeSite.id
		}

		if (process.env === "development") {
			console.log("Request data:", data)
		}

		const dataEncoded = msgpack.encode(data, {codec: MessagePackCodec})
		if (method === "GET") {
			// The endpoint URL may already have a "?" if this is a retry request, which means the hex-encoded data
			// is already part of the URL.
			if (!endpoint.includes("?")) {
				endpoint += "?data=" + hexEncode(dataEncoded)
			}
		} else {
			opts.headers = new Headers({"Content-Type": "application/x-msgpack"})
			opts.body = dataEncoded
		}
	}

	// doingAuthStep should be set to true if we are already accessing the "auth" endpoint or if, as we
	// can see below, there is a 403 error returned as a response.
	if (endpoint === "auth") {
		doingAuthStep = true
	}

	fetch(apiPath + endpoint, opts).then(resp => {
		switch (resp.status) {
			case 200:
				resp.arrayBuffer().then(respData => {
					respData = msgpack.decode(new Uint8Array(respData), {codec: MessagePackCodec})
					if (respData.warn !== undefined) {
						for (let i = 0; i < respData.warn.length; i++) {
							alert(respData.warn[i]) // TODO: show each in a toast
						}
					}
					handleGood(respData.body)
				}).catch(err => {
					alert("An error occurred reading the response: " + err)
				})
				doingAuthStep = false
				return
			case 403:
				// A 403 error may indicate either that nobody is logged in (possibly timed out) or that the user does
				// not have sufficient privileges to make the request.
				// Prompt the user with a login dialog and then retry this same request.
				// The handleError callback will not be called.
				this.$showDialog({
					title: "Log in to continue",
					componentName: "login-dialog",
					acceptText: "LOGIN",
					componentOptions: {noHide: true},
					responseCallback: this.$req.bind(this, method, endpoint, data, handleGood, handleError)
				})
				doingAuthStep = true
				return
			default:
				resp.text().then(errMsg => {
					if (!doingAuthStep) {
						this.$toasted.error(errMsg)
					}
					// If we are trying to log in (or doing something else where the error handling callback is defined),
					// we call the handleError callback. If doing auth step, then the login dialog must be currently open.
					handleError(errMsg)
					doingAuthStep = false // must be inside callback
				}).catch(caughtMsg => {
					this.$toasted.error(caughtErrMessage(caughtMsg.message))
					doingAuthStep = false // must be inside callback
				})
		}
	}).catch(err => {
		// The request could not be made because the browser blocked it or there is no internet connection.
		let msg = err.message
		if (msg !== "") {
			msg = " (" + msg + ")"
		}
		this.$toasted.error("Your request could not be made. Please check your internet connection." + msg)
		doingAuthStep = false
	}).finally(this.$hideLoader.bind(this))

}

function caughtErrMessage(msg) {
	return "An error occurred. (" + msg + ")"
}

// The possible roles for registered users.
const siteRoles = ["author", "editor", "admin", "owner", "super"];

// $atLeastRole says if the user's role on activeSite is at least as high as the role given.
Vue.prototype.$atLeastRole = function(role) {
	return siteRoles.indexOf(this.activeSite.role) >= siteRoles.indexOf(role)
}

new Vue({
	el: "#app",
	router,
	render: h => h(App),
})
