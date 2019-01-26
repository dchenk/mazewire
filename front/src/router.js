import Vue from "vue"
import Router from "vue-router"

import Dashboard from "./components/Dashboard.vue"

import PagesPosts from "./components/pages/PagesPosts.vue"

import PageEdit from "./components/pages/editing/PageEdit.vue"

// import PostSettings from "./components/posts/Settings.vue"

// import Media from "./components/media/Media.vue"
// import MediaAlbums from "./components/media/Albums.vue"

// import Design from "./components/design/Design.vue"

// import Settings from "./components/settings/Settings.vue"
// import SettingsAdvanced from "./components/settings/Advanced.vue"

// import Account from "./components/account/Account.vue"
// import AccountPayment from "./components/account/Payment.vue"
// import AccountMessages from "./components/account/Messages.vue"

// import Help from "./components/help/Help.vue"

// Extra settings for super admins
import AdminSrcs from "./components/super-settings/AdminSrcs.vue"

Vue.use(Router);

export default new Router({
	routes: [
		{
			path: "/",
			component: Dashboard
		},
		{
			path: "/pages",
			component: PagesPosts,
			props: {pType: "page"}
		},
		{
			path: "/pages/trash",
			component: PagesPosts,
			props: {pType: "page", trashed: true}
		},
		{
			path: "/pages/edit/:ppId",
			component: PageEdit,
			name: "page_editor"
		},
		{
			path: "/posts",
			component: PagesPosts,
			props: {pType: "post"}
		},
		{
			path: "/posts/trash",
			component: PagesPosts,
			props: {pType: "post", trashed: true}
		},
		// {
		// 	path: "/posts/settings",
		// 	component: PostSettings
		// },
		{
			path: "/posts/edit/:ppId",
			component: PageEdit,
			name: "post_editor"
		},
		// {
		// 	path: "/media",
		// 	component: Media
		// },
		// {
		// 	path: "/media/albums",
		// 	component: MediaAlbums
		// },
		// {
		// 	path: "/design",
		// 	component: Design
		// },
		// {
		// 	path: "/settings",
		// 	component: Settings
		// },
		// {
		// 	path: "/settings/advanced",
		// 	component: SettingsAdvanced
		// },
		// {
		// 	path: "/account",
		// 	component: Account
		// },
		// {
		// 	path: "/account/payment",
		// 	component: AccountPayment
		// },
		// {
		// 	path: "/account/messages",
		// 	component: AccountMessages
		// },
		// {
		// 	path: "/help",
		// 	component: Help
		// },
		{
			path: "/srcs",
			component: AdminSrcs
		}
	],
	mode: "history",
	base: "/admin/"
})
