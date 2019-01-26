<template>
	<div id="new-section-dialog">
		<h6>Choose a section width type</h6>
		<div class="two-radio-options">
			<div class="bttn-holder">
				<input class="material-radio" type="radio" value="standard" id="sec-standard" v-model="type">
				<label for="sec-standard"><span class="radio-circle material-icons"/>Standard Section</label>
			</div>
			<div class="bttn-holder">
				<input class="material-radio" type="radio" value="fullwidth" id="sec-full" v-model="type">
				<label for="sec-full"><span class="radio-circle material-icons"/>Full-Width Section</label>
			</div>
		</div>
		<button class="material-bttn material-bttn-raised" @click="useDynamicSection">Use a Saved Dynamic Section</button>
	</div>
</template>

<script>

	export default {
		data() { return {
			type: "standard"
		}},
		methods: {
			acceptResponse() {
				this.$parent.responseCallback(this.type)
			},
			useDynamicSection() {
				this.$showErrDialog("This feature is not available yet!")
			}
		}
	}

</script>

<style>

	#new-section-dialog .bttn-holder {
		display: inline-flex;
	}

</style>
