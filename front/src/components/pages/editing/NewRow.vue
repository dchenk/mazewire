<template>
	<div id="new-row-dialog">
		<h6>Choose the kinds of columns this row will have</h6>

		<ul id="col-types-list">
			<li v-for="(columns, rowType) in rowTypes" class="room-row" :class="['row-'+rowType, {'selected-type': type === rowType}]" @click="type = rowType">
				<div v-for="col in columns" class="column">
					{{ col }}
				</div>
			</li>
		</ul>
	</div>
</template>

<script>
	import rowTypes from "./row-types.js"

	export default {
		data() { return {
			type: "full",
			rowTypes: Object.freeze(rowTypes)
		}},
		methods: {
			acceptResponse() {
				this.$parent.responseCallback(this.type)
			}
		}
	}

</script>

<style>

	#col-types-list {
		display: flex;
		flex-wrap: wrap;
	}

	#col-types-list li {
		width: 33.333%;
		margin: 4px;
		padding: 2px;
		box-sizing: border-box;
		border: 3px solid white;
		cursor: pointer;
		border-radius: 3px;
	}

	#col-types-list li.selected-type {
		border-color: #0277BD;
	}

	#col-types-list .column {
		color: #fff;
		background-color: #00BCD4;
		padding: 11px 2px;
		font-weight: 600;
		text-align: center;
	}

	@media screen and (max-width: 575px) {
		#col-types-list li {
			width: 50%;
		}
	}

</style>
