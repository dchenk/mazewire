<template>
	<div id="dashboard">
		<h1>Welcome to your site&rsquo;s admin dashboard!</h1>
		<h2 v-if="numberOfSites > 1">
			You are currently editing the website {{ activeHost }}
		</h2>

		TODO: Change active editing site if user is associated with more than one site as at least
		author...
	</div>
</template>

<script>
	import Vue from "vue"

	export default Vue.extend({
		data() {
			return {
				activeHost: "",
				numberOfSites: 1
			}
		},
		mounted() {
		},
		method: {
			changeActiveSite(newID, newHost) {
				this.activeHost = newHost
				window.activeSite.id = newID
				window.activeSite.id = newHost
			}
		}
	})
</script>
