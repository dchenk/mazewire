<template>
	<div class="dyn-textfield" id="element-name">
		<input type="text" class="dyn-textfield-input" id="elem-name" v-model="customName">
		<label for="elem-name" class="dyn-textfield-label">Name (for editor only)</label>
	</div>
</template>

<script>

	export default {
		data() { return {
			customName: ""
		}},
		methods: {
			getCustomName() {
				return this.newCustomName.trim()
			},
			resetCustomName(name = "") {
				this.customName = name
			}
		}
	}

</script>

<style>
	#element-name {
		margin: -3px 0 9px;
	}
</style>
