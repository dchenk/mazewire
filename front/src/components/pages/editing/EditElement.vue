<template>
	<div id="edit-section-dialog">
		<elem-name ref="customName" />

		<style-editor elem_type="section" ref="styling" />

		<template v-if="elemType === 'section'">
			<h6>Section Width Type</h6>
			<div class="two-radio-options">
				<div class="bttn-holder">
					<input class="material-radio" type="radio" value="standard" id="sec-standard" name="sec_type" v-model="secType">
					<label for="sec-standard"><span class="radio-circle material-icons" />Standard Section</label>
				</div>

				<div class="bttn-holder">
					<input class="material-radio" type="radio" value="fullwidth" id="sec-fullwidth" name="sec_type" v-model="secType">
					<label for="sec-fullwidth"><span class="radio-circle material-icons" />Full-Width Section</label>
				</div>
			</div>
		</template>

		<hr>

		<h6>Custom CSS Styling</h6>
		<div class="ace-editor-container">
			<div id="elem-css-editor" class="ace-editor-inner" />
		</div>

		<hr>

		<h4>List of Custom Element Classes</h4>
		<div class="dyn-textfield">
			<input type="text" class="dyn-textfield-input" id="sec-css-classes" v-model="cssClasses">
			<label for="sec-css-classes" class="dyn-textfield-label">CSS Classes (Space-separated)</label>
		</div>

		<h4>Custom Element Id Attribute</h4>
		<div class="dyn-textfield">
			<input type="text" class="dyn-textfield-input" id="sec-css-id" v-model="cssID">
			<label for="sec-css-id" class="dyn-textfield-label">CSS ID</label>
		</div>
	</div>
</template>

<script>

	import DynamicTextFields from "dynamic-textfields"
	import ElemName from "./ElemName"
	import StyleEditor from "./style-editor/StyleEditor"

	export default {
		components: {
			ElemName,
			StyleEditor
		},
		props: {
			sectionID: {
				type: [Number, String],
				default: 0
			},
			settings: {
				type: Object,
				default: () => {}
			},
		},
		data() { return {
			secType: "",
			cssEditor: null, // Ace editor
			dtf: null // DynamicTextFields object
		}},
		watch: { // necessary in case the old dialog component is also an "edit-section" and oldType and newType are already set
			settings: {
				handler: function() {
					this.resetSettings()
				},
				deep: true
			},
			classes: function(newVal) {
				if (newVal.indexOf(".") !== -1) {
					this.classes = newVal.replace(".", "")
				}
			}
		},
		mounted() {

			this.dtf = new DynamicTextFields("edit-section-dialog")
			this.dtf.registerAll()

			this.cssEditor = this.ace.edit("section-css-editor")
			this.cssEditor.setTheme("brace/theme/tomorrow");
			this.cssEditor.session.setMode("brace/mode/css");

			this.resetOptions()

		},
		methods: {
			acceptResponse() {
				this.$parent.responseCallback(this.sectionID,
					{
						type: this.secType,
						name: this.$refs.customName.getCustomName(),
						styles: this.$refs.styling.getStyles(),
						classes: this.classes,
						id_attr: this.idAttr
					}
				)
			},
			resetSettings() {

				this.secType = this.settings.type

				this.$refs.customName.resetCustomName(this.settings.name)

				this.styling.resetStyling(this.settings.styles)

				if (this.settings.classes) {
					this.classes = this.settings.classes
				} else {
					this.classes = ""
				}

				if (this.settings.id_attr) {
					this.idAttr = this.settings.id_attr
				} else {
					this.idAttr = ""
				}

				this.$nextTick(() => {
					this.dtf.resetAllStyles()
				})

			}
		}
	}
</script>

<style>
	#edit-section-dialog .ace-editor-container {
		height: 200px;
	}
</style>
