<template>
	<div id="pages-posts">
		<button v-if="!trashed" class="material-bttn material-bttn-raised bkg-primary" @click="showCreateDialog">CREATE NEW {{ pType.toUpperCase() }}</button>
		<h5 v-if="ppList.length === 0">(No {{ pType }}s yet)</h5>
		<ul class="pagepost-list">
			<li v-for="pp in ppList" class="material-card">
				<div class="pp-parent">
					<router-link :to="pType+'s/edit/'+pp.content.id | linkSlug" class="edit-link">
						<h6>{{ pp.content.title }}<span class="draft-unsaved-status" v-if="pp.content.status === 'draft'"> (Draft)</span>
							<span class="draft-unsaved-status" v-else-if="pp.content.status === 'unsaved'"> (Unpublished Changes)</span>
						</h6>
					</router-link>
					<a :href="pp.content.slug | linkSlug" target="_blank" class="preview-link"><div class="material-icons">open_in_new</div></a>
				</div>
				<div class="pp-children" v-if="pp.children">
					<div v-for="ppchild in pp.children" class="pp-child">
						<router-link :to="pType+'s/edit/'+ppchild.id | linkSlug" class="edit-link">
							<h6>&nbsp;&mdash; {{ ppchild.content.title }}<span class="draft-unsaved-status" v-if="pp.content.status === 'draft'"> (Draft)</span>
								<span class="draft-unsaved-status" v-else-if="pp.content.status === 'unsaved'"> (Unpublished Changes)</span>
							</h6>
						</router-link>
						<a :href="ppchild.content.slug | linkSlug" target="_blank" class="preview-link"><div class="material-icons">open_in_new</div></a>
					</div>
				</div>
			</li>
		</ul>
		<div class="pagination" v-if="pagesTotal > 20">
			<span class="prev-page" v-if="offset">Previous Page</span> <span class="next-page">Next Page</span>
		</div>
	</div>
</template>

<script>

	import {MDCRipple} from "@material/ripple"

	export default {
		filters: {
			linkSlug(slug) {
				if (slug.substr(0, 1) === "/") { return slug }
				return "/" + slug
			}
		},
		props: {
			"pType": {
				type: String,
				default: "page"
			},
			"trashed": {
				type: Boolean,
				default: false
			}
		},
		data() { return {
			ppList: [],
			pagesTotal: 0,
			offset: 0
		}},
		watch: {
			pType: function() {
				this.refreshList()
			},
			trashed: function() {
				this.refreshList()
			},
			offset: function() {
				this.refreshList()
			}
		},
		created() {
			this.refreshList()
		},
		mounted() {
			MDCRipple.attachTo(document.querySelector("#pages-posts .material-bttn"));
		},
		methods: {
			refreshList() {
				this.$req("GET", "pagepost", {pp_type: this.pType, offset: this.offset, trashed: this.trashed},
					resp => {
						this.ppList = resp.items
					}
				)
			},
			updatePagination(total, offset) {
				this.pagesTotal = total
				this.offset = offset
			},
			showCreateDialog() {
				this.$showDialog({
					title: `Create a ${this.pType}`,
					componentName: "create-pagepost",
					componentOptions: {pType: this.pType},
					acceptText: "CREATE"
				})
			}
		},
	}

</script>

<style>

	.pagepost-list li {
		margin: 0 0 15px 0;
	}

	.pagepost-list .pp-parent, .pagepost-list .pp-child {
		display: flex;
		flex-flow: row nowrap;
		justify-content: space-between;
	}

	.pagepost-list h6 {
		margin: 3px 0 5px;
	}

	.pagepost-list a.edit-link {
	}

	.preview-link .material-icons {
		margin-top: 3px;
	}

	.draft-unsaved-status {
		font-style: italic;
	}

</style>
