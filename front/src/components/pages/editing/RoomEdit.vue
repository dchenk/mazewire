<template>
	<div id="page-editor" :class="{'preview-mode': previewing}">
		<div id="page-edit-preview"><button @click="previewMode" v-if="!previewing" v-html="Preview"/><button @click="exitPreviewMode" v-if="previewing" v-html="Close"/></div>
		<div id="editor-wrap">
			<i-frame name="editing-frame" id="editing-frame" :node-name="'div'" node-i-d="page-content" :styles="editorIframeStyles" @rendered="treeRendered">
				<room-section v-for="(section, sectionIndx) in tree" :key="sectionIndx" v-bind="section" :tree-index="sectionIndx" :previewing="previewing" @positionSection="positionSectionControls"/>
				<div v-if="tree.length === 0 && !previewing" id="no-room-elements">
					<strong>The page is empty</strong><button @click="addSection(0)" class="add-section-empty">ADD SECTION</button>
				</div>
			</i-frame>
		</div>
	</div>
</template>

<script>

	//import draggable from "vuedraggable"
	import RoomSection from "./RoomSection.vue"
	import moduleTypes from "./modules/module-types"
	import {thingsAreEqual} from "../../../things-are-equal"

	import Sortable from "sortablejs"

	import Preview from "../../../assets/svg-inline/preview.svg"
	import Close from "../../../assets/svg-inline/close.svg"
	import EditPencil from "../../../assets/svg-inline/edit-pencil.svg"
	import Add from "../../../assets/svg-inline/add.svg"

	import editorIframeStyles from "./iframe-inline-styles.rawcss"
	import layoutsCSS from "../../../layouts.rawcss"

	const changeMade = "changeMade"

	// rowTypes gives the number of columns each type of row has.
	import rowTypes from "./row-types.js"

	export default {
		components: {
			RoomSection
		},
		props: {
			treeGot: {
				type: Array,
				default: () => { return [] }
			},
			userCSS: {
				type: String,
				default: ""
			}
		},
		data() { return {
			tree: [],
			dynData: {},

			newElIndex: 0, // the index in the tree where the next new section, row, or module will be added
			newElColumnIndex: 0, // the index of the column into which the new module will be added

			moduleTypes: moduleTypes, // Copy to component for use in the template.

			editorIframeStyles: [editorIframeStyles, layoutsCSS, this.userCSS],

			iframeElement: null,
			iframeWindow: null,

			sortable: null,

			sectionAddAboveButton: null,
			sectionAddBelowButton: null,

			previewing: false,
			Preview,
			Close
		}},
		watch: {
			treeGot: function(newTree) {
				if (Array.isArray(newTree)) {
					this.tree = newTree
				}
			},
			userCSS: function(newCSS) {
				this.editorIframeStyles = [editorIframeStyles, layoutsCSS, newCSS]
			},
			template: function(tmpl) {
				if (this.tree.length > 0) {
					// TODO:
					this.$showDialog({
						innerHTML: "There are currently items in the page. Do you want to remove everything in the current design and "+
							"load the new template?",
						acceptText: "Yes, Use New Design"
					})
				}
				this.tree = tmpl // TODO: do some validations if the tree is not empty
			}
		},
		methods: {
			positionSectionControls(rect, sectionIndex) {
				if (!this.sectionAddBelowButton) { // TODO: determine if this conditional can be removed safely
					this.createControlButtons(this.iframeWindow.document)
				}
				const iframeHeight = Number(window.getComputedStyle(this.iframeElement).height.replace("px", ""))
				const left = rect.left + rect.width/2 - 15 + "px"

				let topA = rect.top - 15
				if (topA <= 0) {
					topA += 16
				}
				this.sectionAddAboveButton.style.top = topA + "px"
				this.sectionAddAboveButton.style.left = left
				this.sectionAddAboveButton.style.display = ""

				let topB = rect.bottom - 15
				if (rect.bottom >= iframeHeight) {
					topB -= 16
				}
				this.sectionAddBelowButton.style.top = topB + "px"
				this.sectionAddBelowButton.style.left = left
				this.sectionAddBelowButton.style.display = ""

				this.newElIndex = sectionIndex
			},
			positionRowControls(rect) { // TODO: just like positionSectionControls
			},
			treeRendered() {
				// Here we must re-compute window.frames upon every update (replacement) of the iframe.
				const w = window.frames["editing-frame"]
				this.iframeWindow = w
				if (this.tree.length > 0) {
					this.createControlButtons(w.document)
				}
				this.sortable = new Sortable(w.document.getElementById("page-content"), {
					group: "section",
					handle: ".move-section",
					animation: 200,
					onStart: () => {
						this.sectionAddBelowButton.style.display = "none"
					},
					onChoose: () => {
						this.sectionAddBelowButton.style.display = "none"
					}
				})
				this.iframeElement = document.getElementById("editing-frame")
				const h = w.document.documentElement.getBoundingClientRect().height
				this.iframeElement.style.height = h + "px"
			},
			createControlButtons(iframeDoc) {
				if (this.previewing) { return }
				const a = iframeDoc.createElement("BUTTON")
				a.innerHTML = Add
				a.classList.add("room-edit-control", "room-add-button")
				a.style.display = "none" // Start by not showing it anywhere until mouse over an editable element.
				a.addEventListener("click", () => {
					this.addSection(this.newElIndex) // newElIndex changes dynamically
				})
				iframeDoc.body.appendChild(a)
				this.sectionAddAboveButton = a

				const b = iframeDoc.createElement("BUTTON")
				b.innerHTML = Add
				b.classList.add("room-edit-control", "room-add-button")
				b.style.display = "none" // Start by not showing it anywhere until mouse over an editable element.
				b.addEventListener("click", () => {
					this.addSection(this.newElIndex+1) // newElIndex changes dynamically
				})
				iframeDoc.body.appendChild(b)
				this.sectionAddBelowButton = b
			},
			previewMode() {
				this.previewing = true
				const bttns = this.iframeWindow.document.getElementsByClassName("room-add-button")
				for (let i = bttns.length; i < 0; i--) {
					bttns[i].remove()
				}
			},
			exitPreviewMode() {
				this.previewing = false
				this.createControlButtons(this.iframeWindow.document)
			},
			addSection(intoWhichIndex) {
				this.newElIndex = intoWhichIndex
				this.$showDialog({
					title: "Add a section",
					componentName: "new-section",
					responseCallback: this.insertSection
				})
			},
			insertSection(secType) {
				this.$hideDialog()
				if (secType !== "standard" && secType !== "fullwidth") { // validate section class
					this.$showErrDialog("There was an error inserting the section")
					return
				}
				const ns = {
					type: secType,
					name: "",
					rows: [],
					styles: {}
				}
				this.tree.splice(this.newElIndex, 0, ns)
				this.$emit(changeMade) // activate button to save changes
			},
			editSection(sectionID) {
				this.$showDialog({
					title: "Edit section",
					componentName: "edit-section",
					componentOptions: {
						sectionID: sectionID,
						settings: this.content[sectionID]
					},
					responseCallback: this.updateSection
				})
			},
			updateSection(originalId, newData) {

				this.$hideDialog()

				console.log(newData)

				// Check if the updated section was actually changed.
				if (thingsAreEqual(this.content[originalId], newData)) {
					console.log("sections are equal")
					return
				}
				console.log("sections are not equal")

			},
			duplicateSection(section, secIndex) {
				// TODO
			},
			addRow(intoWhichSectionID) {
				this.newElAfterWhich = intoWhichSectionID
				this.$showDialog({
					title: "Add a row",
					componentName: "new-row",
					responseCallback: this.insertRow
				})
			},
			insertRow(rowType) {
				this.$hideDialog()
				const rowConfig = rowTypes.find(function(elem) {
					return elem.type = rowType
				})
				if (!rowConfig) {
					this.$showErrDialog("There was an error inserting the row")
					return
				}
				const nr = {
					type: rowType,
					styles: {},
					modules: []
				}
				for (let i = 0; i < rowConfig.cols.length; i++) {
					// Create an empty array for each column.
					nr.modules.push([])
				}
				if (this.tree.length === 0) {
					this.tree[0].rows.push(nr)
				} else {
					for (let i = 0; i < this.tree.length; i++) { // loop over sections
						if (this.tree[i].id === this.newElAfterWhich) { // newElAfterWhich is the section ID into which the row is being pushed
							this.tree[i].rows.push(newTreeItem)
							break
						}
					}
				}
				this.$emit(changeMade) // activate button to save changes
			},
			editRow(rowID) {
				this.$showDialog({
					title: "Edit row",
					componentName: "edit-row",
					componentOptions: {
						rowID: rowID,
						data: this.content[rowID]
					},
					responseCallback: this.updateRow
				})
			},
			updateRow(newData) {

				this.$hideDialog()

				console.log("orig", this.content[newData.originalID])
				console.log("new", newData.body)

				newData.body.type = this.content[newData.originalID].type // cannot change row type

				// Check if the updated row is different from before change.
				if (thingsAreEqual(this.content[newData.originalID], newData.body)) {
					console.log("things are equal")
					return
				}
				console.log("things are not equal")

				// handle updates if row was already edited
				if (typeof newData.originalID === "string" && newData.originalID.indexOf("n") !== -1) {
					// replace edited row body (don't create a new one)
					this.content[newData.originalID] = newData.body
					this.$emit(changeMade)
					return
				}

				console.log("got this far")

				// handle first updates to row; this updates all rows with the old id
				const nr = {}
				for (let i = 0; i < this.tree.length; i++) { // loop over sections
					for (let j = 0; j < this.tree[i].rows.length; j++) { // loop over rows
//						console.log("this.tree[i].rows[j].id", this.tree[i].rows[j].id)
						if (this.tree[i].rows[j].id === newData.originalID) {
							// TODO
							// console.log("this.tree[i].rows[j].id OLDID", this.tree[i].rows[j].id)
							// this.tree[i].rows[j].id = newID
						}
					}
				}

				this.$emit(changeMade)

			},
			duplicateRow(row, rowIndex) {
				// TODO
			},
			addModule(rowID, columnIndex) {
				this.newElAfterWhich = rowID
				this.newElColumnIndex = columnIndex
				this.$showDialog({
					title: "Add a module",
					componentName: "new-module",
					responseCallback: this.insertModule
				})
			},
			insertModule(newData) {
				this.$hideDialog()
				if (!moduleTypes.hasOwnProperty(newData.type)) { // validate row type
					this.$showErrDialog("There was an error inserting the module")
					return
				}
				if (this.tree.length === 0) {
					this.$showErrDialog("You must create a row first and then add a module")
					return
				}
				const nm = {}
				let entry = newData.defaultBody
				entry.type = newData.type
				entry.custom_name = ""
				for (let i = 0; i < this.tree.length; i++) { // Loop over sections.
					for (let j = 0; j < this.tree[i].rows.length; j++) { // Loop over rows.
						if (this.tree[i].rows[j].id === this.newElIndex) { // newElIndex is the row into which the module is being pushed.
							for (let col = 0; col < this.tree[i].rows[j].modules.length; col++) {
								if (col === this.newElColumnIndex) {
									this.tree[i].rows[j].modules[col].push(nm) // TODO
									break
								}
							}
						}
					}
				}
				this.$emit(changeMade) // activate button to save changes
				this.editModule(newID)
			},
			editModule(module) {
				this.$showDialog({
					title: "Edit "+moduleTypes[module.type].title+" module",
					componentName: moduleTypes[module.type].componentName,
					componentOptions: {
						data: module
					},
					responseCallback: this.updateModule
				})
			},
			updateModule(newData) {

				this.$hideDialog()

				// TODO: the stuff below is outdated -- there needs to be work done only if DYNAMIC module data was changed
				if (!newData.dyn) {
					return
				}

				console.log("updating module", JSON.parse(JSON.stringify(newData)))

				// set type once here, not in every module component; cannot update module type
				newData.body.type = this.content[newData.originalID].type

				// check if updated module is different from before change
				if (thingsAreEqual(this.content[newData.originalID], newData.body)) {
					return
				}
				console.log("things are not equal")

				// handle updates if module was already edited
				if (typeof newData.originalID === "string" && newData.originalID.indexOf("n") !== -1) {
					// replace edited module body (don't create a new one)
					this.content[newData.originalID] = newData.body
					this.$emit(changeMade)
					return
				}

				// handle first updates to module; this updates all modules with the old id
				const newID = "n"+this.newContent.length
				this.$set(this.content, newID, newData.body) // must use Vue set method to make the property reactive
				this.newContent.push({id: newID, role: "module"})
				for (let i = 0; i < this.tree.length; i++) { // loop over sections
					for (let j = 0; j < this.tree[i].rows.length; j++) { // loop over rows
						for (let k = 0; k < this.tree[i].rows[j].modules.length; k++) { // loop over columns
							for (let l = 0; l < this.tree[i].rows[j].modules[k].length; l++) { // loop over modules
								if (this.tree[i].rows[j].modules[k][l] === newData.originalID) {
									this.tree[i].rows[j].modules[k].splice(l, 1, newID)
								}
							}
						}
					}
				}

				this.$emit(changeMade)

			},
			duplicateModule(module, moduleIndex) {
				// TODO
			}
		}
	}

</script>

<style>

	#page-editor.preview-mode {
		border: 1px solid gray;
		position: fixed;
		width: 100%;
		box-sizing: border-box;
		top: 0;
		left: 0;
		height: 100vh;
		padding: 40px;
		background-color: #EEEEEE;
		z-index: 99999;
	}

	#editor-wrap {
		border: 1px solid gray;
	}

	#page-editor.preview-mode #editor-wrap {
		border: none;
	}

	iframe#editing-frame {
		width: 100%;
		min-height: 200px;
		border: none;
		vertical-align: top;
		background: #fff;
	}

	#page-editor.preview-mode iframe#editing-frame {
		box-shadow: 0 1px 20px 15px rgba(0, 0, 0, 0.2), 0 4px 12px 0 rgba(0, 0, 0, 0.14);
	}

	#page-edit-preview {
		text-align: right;
	}

	#page-edit-preview button {
		background: none;
		border: none;
		padding: 2px 2px 2px 3px;
		margin-bottom: 1px;
	}

	#page-edit-preview button svg {
		vertical-align: top;
	}

	#page-editor.preview-mode #page-edit-preview button {
		position: absolute;
		top: 0;
		right: 0;
		padding: 6px 6px 4px 4px;
	}

	/*#page-tree section.jmsection {*/
		/*-webkit-user-select: none; !* Apply dragging styles to section (applied also to child rows and modules). *!*/
		/*-moz-user-select: none;*/
		/*-ms-user-select: none;*/
		/*user-select: none;*/
		/*cursor: move;*/
	/*}*/

	/* draggable setup (start) */

	.dragarea:empty {
		background-color: #DCEDC8;
		border: 1px solid #C5E1A5;
		border-radius: 3px;
		margin-bottom: 8px;
		padding: 4px 8px;
	}

	.dragarea.row-holder:empty:before {
		content: "Drop rows here...";
	}

	.dragarea.module-holder:empty:before {
		content: "Drop modules here...";
	}

	/* draggable setup (end) */

</style>
