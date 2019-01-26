<template>
	<div id="module-html-editor">
		<div class="ace-editor-container"><div id="html-editor" class="ace-editor-inner"/></div>
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
			editor: {},
			styles: {}
		}},
		watch: {
			componentOptions: {
				handler: function() {this.resetContent()}, // reset the component data without destroying the editor in the DOM
				deep: true
			}
		},
		mounted() {

			this.editor = this.ace.edit("html-editor");
			this.editor.setTheme("brace/theme/tomorrow");
			this.editor.session.setMode("brace/mode/html");
			this.editor.$blockScrolling = Infinity;
			this.editor.setOptions({
				showPrintMargin: false,
				enableBasicAutocompletion: true,
				enableSnippets: true,
				enableLiveAutocompletion: false
			});
			this.resetContent()
			this.editor.resize()
		},
		methods: {
			acceptResponse() {
				this.$parent.responseCallback({
					originalID: this.componentOptions.moduleID,
					body: {
						content: this.editor.getValue()
					}
				})
			},
			resetContent() {
				this.$refs.moduleTitle.resetName(this.componentOptions.data.name)
				if (this.componentOptions.data.content !== undefined) {
					this.editor.setValue(this.componentOptions.data.content)
				} else {
					this.editor.setValue("")
				}
				this.editor.focus()
				this.editor.clearSelection()
				this.styles = this.componentOptions.data.styles
			},
			cancelCallback() {
				this.editor.setValue("")
				this.$refs.moduleTitle.resetName()
				this.styles = {}
			}
		}
	}

</script>

<style>

	#module-html-editor .ace-editor-container {
		height: 150px;
	}

</style>
