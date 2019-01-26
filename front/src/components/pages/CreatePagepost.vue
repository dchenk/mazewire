<template>
	<div id="create-options">
		<p class="err-message" v-if="errMsg !== ''">{{ errMsg }}</p>
		<div class="dyn-textfield">
			<input type="text" class="dyn-textfield-input" id="new-title" v-model="newPageTitle">
			<label for="new-title" class="dyn-textfield-label">Page Title</label>
		</div>
		<div class="dyn-textfield slug-input">
			<input type="text" class="dyn-textfield-input" id="new-slug" @input="setCleanSlug" v-model="newPageSlug">
			<label for="new-slug" class="dyn-textfield-label">URL Path Slug</label>
		</div>
		<p class="formatted-slug">https://{{ activeDomain }}/{{ newPageSlug !== '/' ? newPageSlug : '' }}</p>
	</div>
</template>

<script>

	import DynamicTextFields from "dynamic-textfields"

	export default {
		props: {
			pType: {
				type: String,
				default: "page"
			}
		},
		data() { return {
			newPageTitle: "",
			newPageSlug: "",
			slugChangedManually: false,
			activeDomain: this.activeSite.domain,
			errMsg: ""
		}},
		watch: {
			newPageTitle: function() {
				if (!this.slugChangedManually) {
					this.newPageSlug = this.cleanUpSlug(this.newPageTitle)
				}
			}
		},
		mounted() {

			let tf = new DynamicTextFields("create-options")
			tf.registerAll()

			tf.registerRefresher("new-title", "new-slug") // the slug is set automatically when title is changed

		},
		methods: {
			acceptResponse() {

				const title = this.newPageTitle.trim()
				this.newPageTitle = title
				const slug = this.cleanUpSlug(this.newPageSlug).trim()
				this.newPageSlug = slug

				if (title === "" || slug === "") {
					this.errMsg = "All fields need to be completed and be valid."
					return
				}

				// If the slug is simply "/" allow the request and let the server say if there is an error.
				if (slug.includes("/")) {
					if (slug !== "/") {
						this.errMsg = "You may not have any slashes (/) in the slug. If you're trying to create a page "+
							"for the site's home page, simply type in a single slash."
						return
					}
				}

				this.$req("POST", "pagepost", {title: title, slug: slug, pp_type: this.pType},
					resp => {
						// Redirect to edit the new page.
						this.$hideDialog()
						this.$router.push("/pages/edit/"+resp.new_id)
					}
				)

			},
			cleanUpSlug(slug) {
				// Allow slashes here in case the user is entering simply "/" for the home page.
				return slug.replace(/\s+/g, "-").replace(/[^\w+.$~*-/]/g, "").replace(/-+/g, "-").toLowerCase()
			},
			setCleanSlug(evt) {
				this.newPageSlug = this.cleanUpSlug(evt.target.value)
				this.slugChangedManually = true
			}
		}

	}

</script>

<style>
	.slug-input {
		margin: 10px 0 1px;
	}
	.formatted-slug {
		margin: 6px 0 2px;
	}
</style>
