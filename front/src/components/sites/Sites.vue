<template>
	<div>
		<h6>Switch to another site</h6>
		<ul id="sites-list">
			<li v-for="us in userSites"><a :href="us.domain" target="_blank"><img :src="siteIcon(us)">{{ us.name }}</a> (Role: {{ us.role }})</li>
		</ul>
		<h6>Create a new site</h6>
		<button class="material-bttn material-bttn-raised bkg-primary" @click="createSite">CREATE SITE</button>
	</div>
</template>

<script>
	import defaultSiteIcon from "../../assets/svg-data-url/web-icon.svg"

	export default {
		name: "Sites",
		data() { return  {
			userSites: []
		}},
		created() {
			this.$req("GET", "user/sites", null,
			resp => {
				this.userSites = resp
			})
		},
		methods: {
			createSite() {
				// todo
			},
			siteIcon(site) {
				if (site.favicon) {
					return site.favicon
				}
				if (site.logo) {
					return site.logo
				}
				return defaultSiteIcon
			}
		},
	}
</script>

<style>
	ul#sites-list {}
	#sites-list li {
		display: flex;
		align-items: center;
	}
	#sites-list li a {
		display: inline-flex;
		align-items: center;
		margin-right: 12px;
	}
	#sites-list li img {
		margin-right: 12px;
	}
</style>