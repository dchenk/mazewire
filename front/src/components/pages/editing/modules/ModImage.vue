<template>
	<div id="module-image">
		<div class="dyn-textfield">
			<input type="text" class="dyn-textfield-input" id="image-url-field" v-model="imageURL">
			<label for="image-url-field" class="dyn-textfield-label">Image URL</label>
		</div>
	</div>
</template>

<script>

	export default {
		props: {
			settings: {
				type: Object,
				default: () => {}
			}
		},
		data() { return {
			styles: {},
			imageURL: ""
		}},
		watch: {
			componentOptions: {
				handler: function() {this.resetContent()}, // reset the component data without destroying the editor in the DOM
				deep: true
			}
		},
		methods: {
			acceptResponse() {
			},
			resetContent() {
				this.$refs.moduleTitle.resetName(this.settings.name)
			}
		}
	}

</script>
