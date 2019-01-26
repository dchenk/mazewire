<template>
	<div>
		<div id="title-slug">
			<div class="meta-value">
				<div class="dyn-textfield">
					<input type="text" class="dyn-textfield-input required notempty" id="pg-title" v-model="pagepost.content.title">
					<label for="pg-title" class="dyn-textfield-label">Page title</label>
				</div>
			</div>
			<div class="meta-value">
				<div class="dyn-textfield">
					<input type="text" id="pg-slug" class="dyn-textfield-input required notempty" :class="{'not-set':!slugOk}" :disabled="isHomePage" v-model="pagepost.content.slug">
					<label for="pg-slug" class="dyn-textfield-label">URL slug</label>
					<div id="homepage-options" v-if="isHomePage"><span>(Homepage)</span> <button @click="showHomepageOptions">Change Homepage</button></div>
				</div>
			</div>
		</div>
		<div class="meta-value">
			<button class="material-bttn material-bttn-raised has-icon" @click="moreMetaOptions"><span class="material-icons">menu</span>More Page Options</button>
		</div>
		<h5>Page status: {{ pagepost.content.status }}</h5>
		<button class="material-bttn material-bttn-raised" :disabled="!unsavedChanges" @click="$emit('saveChanges')">Save Changes</button>
		<button class="material-bttn material-bttn-raised" :disabled="!readyToPublish" @click="$emit('publish')">Publish Page</button>
		<div class="page-css-editing">
			<button class="material-bttn material-bttn-raised smaller" v-show="!showingCssEditor" @click="showPageCssEditor">Show CSS Editor</button>
			<button class="material-bttn material-bttn-raised smaller" v-show="showingCssEditor" @click="showingCssEditor = false">Hide CSS Editor</button>
			<codemirror ref="cssEditor" :class="{'css-editor-expanded': showingCssEditor}" :value="userCSS" :options="{mode: 'text/css'}" @input="getChangedCSS"/>
		</div>
	</div>
</template>

<script>

	import DynamicTextFields from "dynamic-textfields"

	const changeMade = "changeMade"

	let tf // the DynamicTextFields object (must be accessible by component methods)

	export default {
		props: {
			unsavedChanges: {
				type: Boolean,
				default: false
			},
			readyToPublish: {
				type: Boolean,
				default: false
			}
		},
		data() { return {
			pagepost: {
				img: "",
				content: {
					slug: "",
					title: "",
					status: "",
				},
				meta: null
			},

			userCSS: "",
			showingCssEditor: false,

			isHomePage: false,
			validSlug: /^\w[\w+.$~*-]{0,78}\w$/,
			slugOk: true
		}},
		watch: {
			slug: function() {
				if (!this.isHomePage) {
					this.slug = this.slug.replace(/\s+/g, "-").replace(/-+/g, "-").replace(/[^\w]$/, "").toLowerCase()
					this.slugOk = this.validSlug.test(this.slug)
				}
				this.$emit(changeMade)
			},
			title: function() {
				this.title = this.title.replace(/\s+/g, " ")
				this.$emit(changeMade)
			}
		},
		mounted() {
			// Set up the styling of the page title and slug fields.
			tf = new DynamicTextFields("title-slug")
			tf.registerAll()
		},
		methods: {
			// setGotData takes the data retrieved by the API call and enter it into the page.
			setGotData(resp) {
				this.pagepost = resp.pagepost
				this.isHomePage = this.pagepost.content.slug === "/"
				this.$nextTick(() => {
					this.refreshTextfields()
				})
				this.userCSS = resp.user_css
			},
			refreshTextfields() {
				if (tf) {
					tf.resetAllStyles()
				}
			},
			showPageCssEditor() {
				this.showingCssEditor = true
				// cssEditor.clearSelection()
				// setTimeout(() => {
					// Resize only after the container's height has expanded fully.
					// cssEditor.resize()
				// }, 124)
			},
			getChangedCSS(newCSS) {
				this.$parent.setNewCSS(newCSS)
				this.$emit(changeMade)
			},
			moreMetaOptions() {
				// TODO
			},
			showHomepageOptions() {
				this.$showErrDialog("homepage options not editable yet")
			},

			// yieldAllData returns an object containing all the changed meta key->value pairs that the user is saving and other
			// configuration settings for the pagepost.
			yieldAllData() {
				return this.pagepost
			}
		}
	}
</script>

<style>
	.meta-value {
		margin: 10px 0;
	}

	#homepage-options {
		margin-top: -25px;
	}

	#homepage-options > span {
		padding-left: 11px;
		color: #2c3e50c2;
	}

	.page-css-editing {
		margin-bottom: 16px;
	}

	.page-css-editing .vue-codemirror {
		visibility: hidden;
		transition: height .12s ease-in-out;
		height: 0;
	}

	.vue-codemirror.css-editor-expanded {
		visibility: visible;
		height: 300px;
	}

	.page-css-editing button.material-bttn {
		margin-bottom: 15px;
	}
</style>