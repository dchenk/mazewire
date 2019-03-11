<template>
	<div id="module-rich-text" />
</template>

<script>

	import Quill from "../quill-setup"

	export default {
		props: {
			delta: {
				type: Array,
				default: () => [{"insert":"\n"}]
			},
		},
		data() { return {
			editor: {}
		}},
		watch: {
			delta_ops: function() {
				this.resetContent() // reset the component data without destroying the editor in the DOM
			}
		},
		mounted() {

			this.editor = new Quill("#module-rich-text", {
				modules: {
					toolbar: [
						[{header: [1, 2, false]}],
						["bold", "italic", "underline", {"align": []}, {"color": []}, "link", "code-block"],
						[{"list": "ordered"}, {"list": "bullet"}],
						[{"indent": "-1"}, {"indent": "+1"}],
						[{"script": "sub"}, {"script": "super"}],
						["clean"]
					]
				},
				theme: "snow"
			})

			this.resetContent()

		},
		methods: {
			acceptResponse() {
				this.$parent.responseCallback({
					originalID: this.componentOptions.moduleID,
					body: {
						delta_ops: JSON.stringify(this.editor.getContents().ops) // stored as JSON in the DB
					}
				})
			},
			resetContent() {
				let ops = this.delta_ops
				if (typeof ops === "string" && ops !== "") {
					ops = JSON.parse(ops)
				}
				if (Array.isArray(ops) && ops.length > 0) {
					this.editor.setContents(ops)
				} else {
					this.editor.setContents([{"insert":"\n"}])
				}
			},
			cancelCallback() {
				this.editor.setContents([{"insert":"\n"}])
			}
		}
	}

</script>

<style>

	#module-rich-text .ql-editor { /* TODO: check this */
		min-height: 80px;
	}

</style>
