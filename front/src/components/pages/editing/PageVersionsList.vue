<template>
	<div id="page-versions">
		<div v-if="versions.length > 0">
			<h5>Page Versions</h5>
			<ul>
				<li v-for="ver in versions" @click="showOldVersion(ver.id)">
					{{ ver.ts.toLocaleDateString() + " " + ver.ts.toLocaleTimeString() }}
				</li>
			</ul>
		</div>
	</div>
</template>

<script>
	export default {
		props: {
			versions: {
				type: Array,
				default: () => { return [] } // the saved versions of this page
			}
		},
		data() { return {
		}},
		methods: {
			showOldVersion(verID) {
				this.$showErrDialog("show old version "+verID)
			}
		}
	}
</script>

<style>
	#page-versions {
		margin-top: 15px;
	}

	#page-versions ul {
		margin: 0;
	}

	#page-versions li {
		cursor: pointer;
	}
</style>