<template>
	<div class="dialog-overlay">
		<div class="dialog-inner" id="main-dialog">
			<h5 v-show="title !== ''">{{ title }}</h5>
			<component ref="dialogChild" :is="componentName" v-if="componentName !== ''" v-bind="componentOptions"/>
			<div v-else-if="innerHTML !== ''" v-html="innerHTML"/>
			<div class="dialog-buttons">
				<button class="dialog-reject material-bttn" v-if="rejectText !== 'N/A'" @click="hideAndCancel">{{ rejectText || "CANCEL" }}</button>
				<button class="dialog-accept material-bttn" @click="accept">{{ acceptText || "ACCEPT" }}</button>
			</div>
		</div>
	</div>
</template>

<script>

	// All available dialog inner components are imported here.

	import AdminDashInfo from "./AdminDashInfo.vue"

	import Sites from "./sites/Sites.vue"

	import CreatePagepost from "./pages/CreatePagepost.vue"
	import LoginDialog from "./LoginDialog.vue"

	// Page editor
	import NewSection from "./pages/editing/AddSection.vue"
	import NewRow from "./pages/editing/NewRow.vue"
	import NewModule from "./pages/editing/NewModule.vue"
	import EditElement from "./pages/editing/EditElement.vue"
	import TemplateChooser from "./pages/editing/TemplateChooser.vue"

	// site settings
	// import NavlinkOptions from "./design/NavlinkOptions.vue"

	export default {
		components: {
			AdminDashInfo,
			Sites,
			LoginDialog,
			CreatePagepost,
			NewSection,
			NewRow,
			NewModule,
			EditElement,
			TemplateChooser
			// NavlinkOptions
		},
		props: {
			title: {
				type: String,
				default: ""
			},
			componentName: {
				type: String,
				default: ""
			},
			innerHTML: {
				type: String,
				default: ""
			},
			componentOptions: {
				type: [String, Number, Boolean, Object, Array],
				default: false
			},
			responseCallback: {
				type: Function,
				default: () => {}
			},
			acceptText: {
				type: String,
				default: ""
			},
			rejectText: {
				type: String,
				default: ""
			}
		},
		methods: {
			accept() {
				// If the dialog consists of a component which has an acceptResponse method defined, that function is
				// called and nothing else is done here.
				if (this.$refs.dialogChild !== undefined && this.$refs.dialogChild.acceptResponse !== undefined) {
					this.$refs.dialogChild.acceptResponse()
					return
				}
				// The special "noHide" property, when true, makes the dialog stay open after the accept button was clicked.
				if (!this.componentOptions || !this.componentOptions.noHide) {
					this.$hideDialog()
				}
			},
			hideAndCancel() {
				this.$parent.dialogActive = false
				if (this.$refs.dialogChild !== undefined && this.$refs.dialogChild.cancelCallback) {
					this.$refs.dialogChild.cancelCallback()
				}
			}
		}
	}

</script>

<style>

	.dialog-overlay {
		position: fixed;
		top: 0;
		left: 0;
		z-index: 9999;
		width: 100vw;
		height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		background-color: rgba(0, 0, 0, 0.18);
	}

	.dialog-inner {
		width: 800px;
		max-width: 94%;
		max-height: 95%;
		overflow-y: auto;
		box-sizing: border-box;
		transform: translateY(-1%);
		background-color: #fff;
		box-shadow: 0 11px 15px -7px rgba(0, 0, 0, 0.2), 0 24px 38px 3px rgba(0, 0, 0, 0.14), 0 9px 46px 8px rgba(0, 0, 0, 0.12);
		border-radius: 2px;
		padding: 11px 24px 7px;
	}

	.dialog-inner h5 {
		margin-top: 11px;
	}

	.dialog-buttons {
		display: flex;
		margin-top: 10px;
		float: right;
	}

	.dialog-buttons .dialog-accept {
		margin-left: 7px;
	}

	.dialog-inner hr {
		margin: 18px 5px;
		border: 2px solid #009688;
		border-radius: 2px;
	}

</style>
