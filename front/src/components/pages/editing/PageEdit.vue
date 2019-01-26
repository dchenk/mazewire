<template>
	<div>
		<page-meta @changeMade="unsavedChanges = true" @saveChanges="saveChanges" @publish="publish" :unsaved-changes="unsavedChanges" :ready-to-publish="readyToPublish" ref="meta"/>
		<button @click="loadTemplate">USE TEMPLATE</button>
		<room-edit @changeMade="unsavedChanges = true" :tree-got="tree" :user-c-s-s="userCSS"/>
		<page-versions-list :versions="versions"/>
	</div>
</template>

<script>
	import PageMeta from "./PageMeta"
	import RoomEdit from "./RoomEdit"
	import PageVersionsList from "./PageVersionsList"

	export default {
		components: {
			PageMeta,
			RoomEdit,
			PageVersionsList
		},
		data() { return {
			pagepost: null,
			tree: [],
			dynData: null,
			versions: [],
			userCSS: "", // User's custom CSS (used in page meta component).

			unsavedChanges: false,
			readyToPublish: false,
		}},
		created() {
			this.$req("GET", "page-edit", {page: Number(this.$route.params.ppId)},
				resp => {
					this.pagepost = resp.pagepost
					this.tree = resp.tree
					this.dynData = resp.dyn_data
					this.versions = resp.versions
					this.userCSS = resp.user_css
					this.$refs.meta.setGotData(resp)
				}
			)
		},
		methods: {
			loadTemplate() {
				this.$showDialog({
					title: "Use a template for the page",
					componentName: "template-chooser"
				})
			},
			saveChanges() {
				if (!this.unsavedChanges) {
					this.$showErrDialog("It looks like you didn't make any changes yet.")
					return
				}
				this.$req("POST", "page-edit",
					{
						"page": Number(this.$route.params.ppId),
						"tree": this.tree,
						"pagepost": this.$refs.meta.yieldAllData(),
						"css": this.userCSS
					},
					resp => {
						this.versions.unshift({dt: resp.insert_dt, id: resp.insert_id})
						this.unsavedChanges = false
						this.readyToPublish = true
					}
				)
			},
			setNewCSS(newCSS) {
				this.userCSS = newCSS
			},
			publish() {
				if (!this.readyToPublish) {
					return
				}
				this.$req("PATCH", "page-edit", {page: Number(this.$route.params.ppId)},
					resp => {
						this.readyToPublish = false
					}
				)
			}
		}
	}
</script>

<style>
</style>