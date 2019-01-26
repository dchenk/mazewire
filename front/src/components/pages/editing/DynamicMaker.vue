<template>
	<div class="component-dynamic-maker">
		this component gives the options to make the section, row, or module dynamic
	</div>
</template>

<script>

	export default {
		props: {
			componentOptions: {
				type: [String, Number, Boolean, Object, Array],
				default: false
			}
		},
		data() { return {
		}},
		watch: {
			componentOptions: {
				handler: function() {this.resetContent()}, // Reset the component data without destroying the editor in the DOM.
				deep: true
			}
		}
	}

</script>
