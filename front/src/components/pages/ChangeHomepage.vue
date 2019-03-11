<template>
	<div id="change-homepage">
		<p class="err-message" v-if="errMsg !== ''">
			{{ errMsg }}
		</p>
		<div class="dyn-textfield">
			<input type="text" class="dyn-textfield-input" id="new-title" v-model="newPageTitle">
			<label for="new-title" class="dyn-textfield-label">Page Title</label>
		</div>
		<div class="dyn-textfield slug-input">
			<input type="text" class="dyn-textfield-input" id="new-slug" @input="setCleanSlug" v-model="newPageSlug">
			<label for="new-slug" class="dyn-textfield-label">URL Path Slug</label>
		</div>
		<div class="formatted-slug">
			https://{{ activeDomain }}/{{ newPageSlug }}
		</div>
	</div>
</template>

<script>

	import DynamicTextFields from "dynamic-textfields"

	export default {
		data() { return {
			newPageTitle: "",
			newPageSlug: "",
			slugChangedManually: false,
			activeDomain: window.activeSite.domain,
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

			let tf = new DynamicTextFields("change-homepage")
			tf.registerAll()
			this.$nextTick(() => {
				tf.resetAllStyles()
			})

		},
		methods: {
			validateAndCreate() {

				const title = this.newPageTitle.trim()
				this.newPageTitle = title
				const slug = this.cleanUpSlug(this.newPageSlug)
				this.newPageSlug = slug

				if (title === "" || slug === "") {
					this.errMsg = "All fields need to be completed and be valid"
					return
				}

				this.$req("pagepost/create", {title: title, slug: slug},
					resp => {
						// redirect to new page
						this.$router.push("/pages/edit/"+resp)
					}
				)

			},
			updatePagination(total, offset) {
				this.pagesTotal = total
				this.offset = offset
			},
			cleanUpSlug(slug) {
				return slug.replace(/\s+/g, "-").replace(/-+/g, "-").toLowerCase()
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
</style>
