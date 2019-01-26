export default [
	{
		name: "Dashboard",
		path: "",
		icon: "dashboard"
	},
	{
		name: "Pages",
		path: "pages",
		icon: "web",
		tabs: [
			{
				name: "PUBLISHED & DRAFTS",
				path: ""
			},
			{
				name: "TRASH",
				path: "trash"
			}
		]
	},
	{
		name: "Blog Posts",
		path: "posts",
		icon: "description",
		tabs: [
			{
				name: "PUBLISHED & DRAFTS",
				path: ""
			},
			{
				name: "TRASH",
				path: "trash"
			},
			{
				name: "SETTINGS",
				path: "settings",
				role: "admin"
			}
		]
	},
	{
		name: "Media",
		path: "media",
		icon: "photo_library",
		tabs: [
			{
				name: "ALL MEDIA",
				path: ""
			},
			{
				name: "ALBUMS",
				path: "albums"
			}
		]
	},
	{
		name: "Design",
		path: "design",
		icon: "palette"
	},
	{
		name: "Site Settings",
		path: "settings",
		icon: "build",
		tabs: [
			{
				name: "MAIN SETTINGS",
				path: ""
			},
			{
				name: "ADVANCED",
				path: "advanced"
			}
		],
		role: "admin"
	},
	{
		name: "Admin Source Files",
		path: "srcs",
		icon: "code",
		role: "super"
	},
	{
		name: "Account",
		path: "account",
		icon: "account_box",
		tabs: [
			{
				name: "MAIN ACCOUNT OPTIONS",
				path: ""
			},
			{
				name: "PAYMENT",
				path: "payment"
			},
			{
				name: "MESSAGES",
				path: "messages"
			}
		]
	},
	{
		name: "Help & Info",
		path: "help",
		icon: "help"
	}
]
