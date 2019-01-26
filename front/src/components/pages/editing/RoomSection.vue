<template>
	<section :class="['room-section', 'section-'+type]">
		type: {{ type }} <br>Rows: {{ rows }}
		<button class="move-section room-edit-control" v-if="!previewing" v-html="Drag"/>
	</section>
</template>

<script>
	// import Sortable from "sortablejs"
	import Drag from "../../../assets/svg-inline/drag.svg"

	const positionSection = "positionSection"

	export default {
		name: "RoomSection",
		props: {
			type: {
				type: String,
				default: "standard"
			},
			name: {
				type: String,
				default: ""
			},
			rows: {
				type: Array,
				default: () => { return [] }
			},
			styles: {
				type: Object,
				default: () => { return {} }
			},
			treeIndex: {
				type: Number,
				required: true
			},
			previewing: {
				type: Boolean,
				default: false
			}
		},
		data() { return	{
			sortable: null,

			Drag
		}},
		mounted() {
			this.$el.addEventListener("mouseenter", () => {
				this.$emit(positionSection, this.$el.getBoundingClientRect(), this.treeIndex)
			})
			// this.$el.addEventListener("mouseleave", () => {
			// 	// this.$emit(positionSection, this.$el.getBoundingClientRect(), this.treeIndex)
			// 	// this.hovering = false
			// })
			// this.sortable = new Sortable(this.$el)
		},
		methods: {
			// Attach listeners to the addAbove and addBelow buttons to re-position them on page scroll or resize.
			// attachPositionListeners(iframeWindow) {
			// 	iframeWindow.addEventListener("move_elem", this.positionInsertButtons.bind(this))
			// },
			// positionInsertButtons() {
			// 	const sectionRect = this.$el.getBoundingClientRect()
			// 	console.log("rect:", this.$el, sectionRect)
			// }
		}
	}
</script>

<style>

</style>