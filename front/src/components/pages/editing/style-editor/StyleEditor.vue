<template>
	<ul id="styles-editor">
		<li v-for="(property, propIndx) in styling">
			<h6 class="prop-help" v-if="property.help !== undefined" v-html="property.help" />
			<template v-if="property.type === 'text'">
				<div class="dyn-textfield">
					<input type="text" class="dyn-textfield-input" :id="property.prop_key+op.val" v-model="styling[propIndx].val">
					<label :for="property.prop_key+op.val" class="dyn-textfield-label">{{ property.name }}</label>
				</div>
			</template>
			<template v-if="property.type === 'radio'">
				<h4>{{ property.name }}</h4>
				<div v-for="op in property.opt">
					<input type="radio" class="material-radio" :id="property.prop_key+op.val" :value="op.val" v-model="styling[propIndx].val">
					<label :for="property.prop_key+op.val"><span class="radio-circle material-icons" />{{ op.name }}</label>
				</div>
			</template>
			<template v-if="property.type === 'select'">
				<h4>{{ property.name }}</h4>
				<select v-model="styling[propIndx].val">
					<option v-for="op in property.opt" :value="op.val">
						{{ op.name }}
					</option>
				</select>
			</template>
			<div class="input-err" v-if="property.err !== undefined && property.invalid" v-html="property.err || 'Your input is invalid'" />
		</li>
	</ul>
</template>

<script>

	import allProps from "./properties"
	import moduleTypes from "../modules/module-types"

	const elemProps = {
		site: ["font_fam"],
		page: ["std_section_width"],
		section: ["pad_top", "pad_bottom", "bkg_color", "bkg_img", "bkg_size", "section_row_space"],
		row: ["pad_top", "pad_bottom", "bkg_color", "bkg_img", "bkg_size", "row_col_space"],
		module: {}
	}

	export default {
		name: "StyleEditor",
		props: {
			elemType: {
				type: String,
				required: true
			},
			elemSubType: {
				type: String,
				default: "" // TODO: or require?
			}
		},
		data() { return {
			styling: []
		}},
		watch: {
			elem_type: function() {
				this.resetStyling()
			},
			elem_sub_type: function() {
				this.resetStyling()
			}
		},
		methods: {
			getStyles() {
				let s = {}
				for (let i = 0, l = this.styling.length; i < l; i++) {
					// Into s copy over the value set for the property (ignoring all other settings).
					s[this.styling[i].prop_key] = this.styling[i].val
				}
				return s
			},
			resetStyling(styles = {}) {
				let propsToUse = []
				if (this.elem_sub_type === undefined || this.elem_sub_type === "") {
					propsToUse = elemProps[this.elem_type]
				} else {
					propsToUse = elemProps[this.elem_type][this.elem_sub_type]
				}
				// Identify the properties that apply to the element being styled.
				for (let i = 0, l = propsToUse.length; i < l; i++) {
					let theProp = allProps[propsToUse[i]]
					theProp.prop_key = propsToUse[i] // prop_key is the key of the object in the allProps object.
					theProp.invalid = false
					if (styles[theProp.prop_key] !== undefined) {
						theProp.val = styles[theProp.prop_key]
					}
					propsToUse.push(theProp)
				}
				this.styling = propsToUse
			}
		}
	}

</script>

<style>

	#styles-editor {
		margin: 0;
		padding-left: 0;
	}

	#styles-editor h6.prop-help {
		margin: 9px 0 6px;
		line-height: 1.3em;
	}

	#styles-editor .input-err {
		font-size: 16px;
		font-style: italic;
	}

</style>
