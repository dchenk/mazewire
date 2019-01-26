<template>
	<div>
		<button class="page-edit-icon room-edit-control" v-html="EditPencil"/>
		<button class="page-edit-icon room-edit-control" v-html="Drag"/>
	</div>
</template>

<script>
	import EditPencil from "../../../assets/svg-inline/edit-pencil.svg"
	import Drag from "../../../assets/svg-inline/drag.svg"

	export default {
		data() { return {
			EditPencil,
			Drag
		}},
		mounted() {
		}
	}
</script>