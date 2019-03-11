<template>
	<div id="app">
		<nav class="mdc-drawer mdc-drawer--persistent" id="navdrawer">
			<nav class="mdc-drawer__drawer">
				<div class="adjusted-for-toolbar">
					<div id="current-site-switch">
						<em>Currently editing:</em><b>{{ this.activeSite.name }}</b>
						<button id="switch-site" class="material-bttn material-bttn-raised smaller" @click="switchSite">
							<span>Switch</span><span class="material-icons">arrow_drop_down</span>
						</button>
					</div>
					<hr>
					<router-link v-for="link in navLinks.slice(0, navLinks.length-2)" :key="link.path" v-if="link.role === undefined || $atLeastRole(link.role)" :to="'/'+link.path" @click.native="handleMainNavClick" :class="{'has-tabs': link.tabs !== undefined}">
						<i class="material-icons" aria-hidden="true">{{ link.icon }}</i>{{ link.name }}
					</router-link>
					<hr>
					<router-link v-for="link in navLinks.slice(navLinks.length-2)" :key="link.path" :to="'/'+link.path" @click.native="handleMainNavClick" :class="{'has-tabs': link.tabs !== undefined}">
						<i class="material-icons" aria-hidden="true">{{ link.icon }}</i>{{ link.name }}
					</router-link>
				</div>
			</nav>
		</nav>
		<header id="main-toolbar">
			<div>
				<div class="material-icons" id="main-menu-icon" @click="toggleTheNav">
					menu
				</div>
				<h4>{{ activePosition.name }}</h4>
				<span v-if="loading" v-html="loadingIcon" />
			</div>
			<!--<button id="dash-info" @click="showDashInfo"><img src="./assets/dash-logo-icon.svg"></button>--><!-- TODO: there needs to be a logo and a link to mazewire.com -->
			<div id="user-info">
				<span id="user-fname" @click="showProfileOptions">{{ userInfo.fname }}</span>
				<button @click="showProfileOptions">
					<img :src="userImg">
				</button>
			</div>
		</header>
		<div id="main">
			<div class="adjusted-for-toolbar main-body">
				<div id="section-tabs" class="dynamic-tabs" v-show="onTabbingPath">
					<div class="dynamic-tab-arrow arrow-left">
						<div class="material-icons" aria-label="scroll left">
							navigate_before
						</div>
					</div>
					<div class="dynamic-tabs-framer">
						<nav class="dynamic-tabs-list">
							<router-link v-for="tab in allTabs" :key="tab.id" :to="'/'+tab.parentPath+(tab.path ? '/'+tab.path : '')" :id="tab.id" class="dynamic-tab">
								{{ tab.name }}
							</router-link>
						</nav>
						<div class="dynamic-tabs-indicator">
							<div class="dt-indicator-bar" />
						</div>
					</div>
					<div class="dynamic-tab-arrow arrow-right">
						<div class="material-icons" aria-label="scroll right">
							navigate_next
						</div>
					</div>
				</div>
				<router-view />
			</div>
			<footer id="bottom-footer">
				Copyright &copy; 2017 Mazewire LLC
			</footer>
		</div>
		<div id="user-popover" class="material-card" v-show="userInfoShowing">
			<img :src="userImg">
			<b>{{ userInfo.fname }} {{ userInfo.lname }}</b>
			<router-link to="/account">
				Account Details
			</router-link>
			<button @click="signOut" class="material-bttn">
				Sign Out
			</button>
		</div>
		<DialogBox v-show="dialogActive" v-bind="dialogSetup" />
	</div>
</template>

<script>

	import {MDCPersistentDrawer} from "@material/drawer"

	import router from "./router"
	import mainNavLinks from "./nav-links"
	import DialogBox from "./components/Dialog"
	import DynamicTabs from "dynamic-tabs"

	// Iframe component with slots.
	import "./components/v-iframe.js"

	import portraitDefault from "./assets/svg-data-url/portrait-default.svg"
	import loadingIcon from "./assets/svg-inline/loading.svg"

	// noTabRouteNames lists the names of the routes where tabs should not be displayed even though the path's first
	// part has tabs defined.
	const noTabRouteNames = ["page_editor", "post_editor"]

	// findNavLink returns the first element of mainNavLinks whose path property matches the path argument.
	function findNavLink(path) {
		for (let i = 0, l = mainNavLinks.length; i < l; i++) {
			if (mainNavLinks[i].path === path) {
				return mainNavLinks[i]
			}
		}
		return {}
	}

	// firstRoute is the split path (following "admin") upon initial load of the page.
	const firstRoute = window.location.pathname.split("/").slice(2)

	export default {
		components: {
			DialogBox
		},
		data() { return {
			activePosition: findNavLink(firstRoute[0]), // a copy of the first-level element from mainNavLinks
			navObject: null,
			tabsObj: null,
			navLinks: mainNavLinks,
			allTabs: [], // allTabs contains all of the tabs in the app as a single flat array; it is used for rendering the entire list of tabs with v-for
			sectionTabs: {}, // sectionTabs contains all main nav sections and each section's tabs; it is used for registering tabs when a section is switched to
			onTabbingPath: true,
			userInfo: {id: 0, fname: "", lname: "", img: ""},

			userInfoShowing: false,

			// dialogSetup specifies how to set up the user (pop-up) dialog at any particular moment.
			// When the value of dialogSetup is changed, all properties are optional, but either componentName or innerHTML should be set.
			dialogSetup: {
				title:            "", // the <h5> title to display
				componentName:    "", // the name of the component to create inside the dialog
				innerHTML:        "", // the HTML to insert inside the dialog body
				componentOptions: {}, // data to pass through to the child component created within the dialog
				responseCallback: function(){}, // a callback to be called from within the acceptResponse method in the child component
				acceptText:       "", // the text in the "accept" button
				rejectText:       ""  // the text in the "reject" button, or "N/A" to not display such a button
			},

			// dialogActive indicates whether the dialog is currently visible
			dialogActive: false,

			// loading indicates whether someting on the page is currently loading
			loading: true,
			loadingIcon
		}},
		computed: {
			userImg: function() {
				return this.userInfo.img ? this.userInfo.img : portraitDefault
			}
		},
		created() {
			mainNavLinks.forEach(navLink => {
				if (navLink.tabs !== undefined) {
					this.sectionTabs[navLink.path] = []
					navLink.tabs.forEach(tab => {
						const tabID = "t-" + navLink.path + tab.path;
						this.allTabs.push({
							id: tabID,
							name: tab.name,
							path: tab.path,
							parentPath: navLink.path
						})
						this.sectionTabs[navLink.path].push(tabID)
					})
				}
			})

			this.activeSite = window.activeSite
			this.userInfo = window.userInfo
		},
		mounted() {

			// create main nav drawer and set it to open if on desktop or tablet landscape screen
			this.navObject = new MDCPersistentDrawer(document.getElementById("navdrawer"));
			if (window.innerWidth > 600) {
				this.navObject.open = true;
			}

			this.checkIfOnTabbingPath(this.activePosition)

			// if user landed on a tab that's not first, highlight the main nav section link
			if (firstRoute.length > 1 && firstRoute[1] !== "") {
				document.getElementsByClassName("has-tabs router-link-active")[0].classList.add("router-link-exact-active")
			}

			// Attach the tabs element.
			this.tabsObj = new DynamicTabs("section-tabs");
			if (this.activePosition.tabs !== undefined) {
				// We need to switch to the non-first tab.
				this.tabsObj.registerTabs(this.sectionTabs[this.activePosition.path])
				if (firstRoute.length > 1 && firstRoute[1] !== "") {
					const tabIndex = this.sectionTabs[firstRoute[0]].indexOf("t-" + firstRoute[0] + firstRoute[1]);
					if (tabIndex > 0) {
						this.tabsObj.setActiveTabIndex(tabIndex);
					}
				}
			}

			// Retrieve the messages pertaining to this site and user.
			// todo

			// After each router-link click, adjust active section link and toolbar title.
			router.afterEach((to, from) => {

				// remove from the paths the first part (empty "")
				const toPath = to.path.split("/").slice(1)
				const fromPath = from.path.split("/").slice(1)

//				console.log("fromPath", fromPath)
//				console.log("toPath", toPath)

				if (toPath[0] !== fromPath[0]) { // the section is changed (main nav link clicked)
					const newActivePos = findNavLink(toPath[0])
					this.tabsObj.deregisterAllTabs()
					this.checkIfOnTabbingPath(newActivePos) // a child component will set onTabbingPath to false if it needs to
					this.activePosition = newActivePos // setting it here because DOM is not updated until all data changes are flushed
					this.$nextTick(() => { // update tabs after new data is flushed
						if (this.onTabbingPath) { // if the new section has tabs
							this.tabsObj.registerTabs(this.sectionTabs[toPath[0]])
						}
					})
				} else { // stayed on the same section, but tab changed
					this.checkIfOnTabbingPath(this.activePosition) // a child component will set onTabbingPath to false if it needs to
					if (toPath.length > 1) { // we don't want to check here if a main nav link was clicked (to switch to first tab) because the first tab would be activated twice
						// set active main nav link class (router-link-exact-active) if the main nav link has-tabs
						this.$nextTick(() => {
							document.getElementsByClassName("has-tabs router-link-active")[0].classList.add("router-link-exact-active") // $nextTick because Vue router messes with classList
						})
					}
				}

			})

			// If the user is not logged in, prompt them to log in.
			if (!this.$atLeastRole("author")) {
				this.$showDialog({
					title: "Log in to continue",
					componentName: "login-dialog",
					acceptText: "LOGIN",
					componentOptions: {noHide: true},
					responseCallback: this.getAllMessages.bind(this)
				})
				return // Do not get user or site messages.
			}

			this.getAllMessages()

		},
		methods: {
			toggleTheNav() {
				this.navObject.open = !this.navObject.open
			},
			handleMainNavClick() {
				this.tabsObj.setActiveTabIndex(0)
				if (window.innerWidth < 601) {
					this.navObject.open = false;
				}
			},
			checkIfOnTabbingPath(activePosition) {
				this.onTabbingPath = activePosition.tabs !== undefined && noTabRouteNames.indexOf(this.$route.name) < 0
			},

			// showDashInfo shows a dialog telling the user about the sites admin dashboard and offering a link to go to the main website.
			showDashInfo() {
				this.$showDialog({
					//title: "The Mazewire Site Dashboard",
					componentName: "admin-dash-info",
					acceptText: "CLOSE",
					rejectText: "N/A",
				})
			},
			switchSite() {
				this.$showDialog({
					title: "Your Web Sites",
					componentName: "sites",
					acceptText: "CLOSE",
					rejectText: "N/A",
				})
			},
			showProfileOptions() {
				this.userInfoShowing = !this.userInfoShowing
			},
			getAllMessages() {
				// TODO: display them each in a snackbar
			},
			getUserMessages() {
				// TODO: display them each in a snackbar
			},
			getSiteMessages() {
				// TODO: display them each in a snackbar
			},
			signOut() {
				this.$showDialog({
					innerHTML: "Are you sure you want to sign out?",
					responseCallback: function() {
						alert("sign out!")
					},
					acceptText: "YES"
				})
			}
		}
	}
</script>

<style>

	header#main-toolbar {
		position: fixed;
		top: 0;
		left: 0;
		width: 100vw;
		z-index: 9999;
		display: flex;
		align-items: center;
		justify-content: space-between;
		box-shadow: 0 2px 4px -1px rgba(0, 0, 0, 0.2), 0 4px 5px 0 rgba(0, 0, 0, 0.14), 0 1px 10px 0 rgba(0, 0, 0, 0.12);
		background-color: var(--theme-primary);
		color: rgba(255, 255, 255, 0.9);
	}

	#main-toolbar > div {
		display: flex;
		align-items: center;
	}

	#main-toolbar > div span {
		margin-left: 17px;
	}

	#main-toolbar > div span svg {
		vertical-align: top;
	}

	#main-toolbar h4 {
		margin: 0;
		font-size: 21px;
	}

	#main-toolbar button {
		border: none;
		background: none;
	}

	#main-menu-icon {
		padding: 15px 24px;
		margin-right: 20px;
		user-select: none;
	}

	#user-info {
		color: #fff;
		margin-right: 26px;
		display: flex;
		align-items: center;
	}

	#user-info button {
		padding: 0;
		margin: 0;
	}

	#user-fname {
		margin-right: 12px;
		font-size: 14px;
		font-weight: 600;
		cursor: pointer;
	}

	#user-info img {
		height: 30px;
		width: 30px;
		border-radius: 50%;
	}

	/*#toolbar-top button img {*/
		/*height: 22px;*/
		/*max-width: 22px;*/
	/*}*/

	/*#toolbar-top button#switch-site {*/
		/*margin-left: 11px;*/
		/*color: #2d2d2d;*/
		/*background-color: rgba(255, 255, 255, 0.9);*/
		/*font-size: 15px;*/
		/*font-family: inherit;*/
		/*border-radius: 3px;*/
		/*padding: 0 8px 0 12px;*/
		/*display: flex;*/
		/*align-items: center;*/
	/*}*/

	/*#toolbar-top button#switch-site:hover {*/
		/*background-color: #ffffff;*/
	/*}*/

	button#switch-site span {
		line-height: 22px;
	}

	#switch-site span.material-icons {
		font-size: 22px;
	}

	#main-toolbar button img {
		height: 22px;
		max-width: 22px;
	}

	#current-site-switch em, #current-site-switch b {
		display: block;
	}

	#current-site-switch b {
		font-weight: 600;
		margin: 2px 0 10px;
	}

	button#switch-site {
		width: 100%;
		display: flex;
		align-items: center;
		padding-top: 8px;
		padding-bottom: 8px;
		justify-content: center;
	}

	#navdrawer a, #current-site-switch {
		padding: 14px 13px 14px 15px;
	}

	#navdrawer a {
		color: #757575;
		display: flex;
		font-size: 16px;
		font-weight: 600;
		align-items: center;
	}

	#navdrawer a i {
		margin-right: 11px;
	}

	#navdrawer hr {
		height: 1px;
		margin: 4px 0;
		border: none;
		background-color: rgba(0,0,0,.12);
	}

	#user-popover {
		z-index: 9999;
		position: absolute;
		right: 26px;
		top: 84px;
		display: grid;
		background: #fff;
		grid-template-areas: "user_img user_name" "user_img acct_details" "signout signout";
		grid-template-columns: 96px auto;
		padding: 19px;
		width: 260px;
	}

	#user-popover img {
		grid-area: user_img;
		width: 86px;
	}

	#user-popover b {
		grid-area: user_name;
		font-size: 15px;
		font-weight: 600;
		margin: 0 0 6px;
	}

	#user-popover a {
		grid-area: acct_details;
		font-weight: 400;
	}

	#user-popover button {
		grid-area: signout;
		margin-top: 16px;
		background-color: var(--theme-red);
		color: #fff;
		font-weight: 600;
	}

	#navdrawer .router-link-exact-active, #navdrawer .router-link-exact-active > i {
		color: var(--theme-primary);
	}

	nav.mdc-drawer__drawer {
		overflow-y: auto;
	}

	.adjusted-for-toolbar {
		margin-top: 70px;
	}

	#section-tabs {
		margin: 7px 0 9px;
		border-bottom: 1px solid #9E9E9E;
	}

	@media (max-width: 600px) {
		#navdrawer {
			position: absolute;
			z-index: 3;
			height: 100vh;
		}
		#user-fname {
			display: none;
		}
	}

	@media screen and (max-width: 415px) {

		#main-menu-icon {
			margin-right: 15px;
		}

		#main-toolbar h4 {
			font-size: 20px;
		}

	}

</style>
