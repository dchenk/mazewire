<template>
	<ul>
		<li v-for="tmpl in templates" @click="chosenTemplate = tmpl.name">
			<h6>{{ tmpl.name }}</h6>
		</li>
	</ul>
</template>

<script>
	import templateA from "./templates/a.js"
	const templates = [
		{
			name: "A",
			tmpl: templateA
		}
	]

	export default {
		data() { return {
			templates,
			chosenTemplate: ""
		}}
	}
</script>

<style>

</style>