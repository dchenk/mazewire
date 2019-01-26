<template>
	<div>
		insert a dynamic element
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
		created() {
			// this.$req("page-edit/dyn-"+this.componentOptions.elType)
		}
	}
</script>

<style>

</style>
