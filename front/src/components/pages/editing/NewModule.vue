<template>
	<ul id="module-types-list">
		<li v-for="(moduleInfo, moduleType) in moduleTypes" v-if="!moduleInfo.content_type || contentType === moduleInfo.content_type" :class="{'selected-type': type === moduleType}" @click="type = moduleType">
			<span class="material-icons">{{ moduleInfo.icon }}</span>{{ moduleInfo.title }}
		</li>
	</ul>
</template>

<script>

	import moduleTypes from "./modules/module-types"

	export default {
		props: {
			contentType: {
				type: String,
				default: ""
			}
		}, // List all the modules that have no "type" specified, along with the modules that have the right type if contentType prop is set.
		data() { return {
			moduleTypes: moduleTypes
		}},
		methods: {
			acceptResponse() {
				this.$parent.responseCallback({
					type: this.type,
					default_body: this.moduleTypes[this.type].default_body
				})
			}
		}
	}

</script>

<style>

	#module-types-list {
		display: flex;
		flex-wrap: wrap;
	}

	#module-types-list li {
		width: 33.333%;
		box-sizing: border-box;
		padding: 6px 10px;
		margin-bottom: 4px;
		display: flex;
		border: 3px solid white;
		border-radius: 3px;
		cursor: pointer;
	}

	#module-types-list  li.selected-type {
		border-color: #0277BD;
	}

	#module-types-list li span:first-child {
		margin-right: 6px;
	}

	@media screen and (max-width: 575px) {
		#module-types-list li {
			width: 50%;
		}
	}

</style>
