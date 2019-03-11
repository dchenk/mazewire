<template>
	<div>
		<h4>Current Versions</h4>
		<b>CSS:</b> {{ currentSrcs["admin-style"] }}<br>
		<b>JS Manifest:</b> {{ currentSrcs["admin-manifest"] }}<br>
		<b>JS App:</b> {{ currentSrcs["admin-app"] }}<br>
		<b>JS Vendor:</b> {{ currentSrcs["admin-vendor"] }}<br>
		<h4>New Versions</h4>
		<textarea v-model="newSrcs" /><br>
		<button type="button" @click="checkIfFilesExist">
			Check if files exist
		</button>
		<button type="button" @click="update">
			Update
		</button>
	</div>
</template>

<script>

	export default {
		data() { return {
			currentSrcs: {},
			newSrcs: "",
			newSrcsJSON: {}
		}},
		created() {
			this.$req("GET", "admin", null,
				srcs => {
					this.currentSrcs = srcs
				},
				err => {
					this.$showErrDialog(err)
				}
			)
		},
		methods: {
			update() {
				if (!this.parseJSON()) {
					return
				}
				this.$req("POST", "admin", {srcs: this.newSrcsJSON},
					resp => {
						this.$showDialog({
							innerHTML: "GOOD. " + resp + " files updated"
						})
					},
					err => {
						this.$showErrDialog(err)
					}
				)
			},
			checkIfFilesExist() { // check that they're actually on the GCS bucket
				if (!this.parseJSON()) {
					return
				}
				for (let prop in this.newSrcsJSON) {
					if (this.newSrcsJSON.hasOwnProperty(prop)) {
						let url = "https://storage.googleapis.com/mazewire_resources/" +
							(prop === "css" ? "app." : prop+".") + this.newSrcsJSON[prop] +
							(prop === "css" ? ".css" : ".js")
						window.open(url)
					}
				}
			},
			parseJSON() {
				try {
					this.newSrcsJSON = JSON.parse(this.newSrcs)
					return true
				} catch (e) {
					this.$showErrDialog("Error parsing input JSON object: " + e)
				}
				return false
			}
		}
	}

</script>

<style>
	.big-input {
		padding: 5px 8px;
		width: 100%;
	}
</style>
