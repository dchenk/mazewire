import Vue from "vue"

Vue.component("i-frame", {
	props: {
		nodeName: {
			type: String,
			default: "div"
		},
		nodeID: {
			type: String,
			required: true
		},
		styles: {
			type: Array,
			default: () => { return [] }
		}
	},
	methods: {
		renderChildren() {
			const slotsChildren = this.$slots.default
			const doc = this.$el.contentDocument
			const el = doc.createElement("DIV")
			doc.body.appendChild(el)

			for (let i = 0; i < this.styles.length; i++) {
				const sc = doc.createElement("STYLE")
				sc.innerHTML = this.styles[i]
				doc.head.appendChild(sc)
			}

			// Must evaluate within this context.
			const nodeName = this.nodeName
			const nodeID = this.nodeID

			const iFrameApp = new Vue({
				data: {
					children: Object.freeze(slotsChildren) // Freeze to prevent unnecessary reactification of vNodes.
				},
				render(h) {
					return h(nodeName, {attrs: {id: nodeID}}, this.children)
				}
			})
			iFrameApp.$mount(el)

			this.$emit("rendered")
		}
	},
	render(h) {
		return h("iframe", {
			on: { load: this.renderChildren },
			domProps: {
				src: "about:blank"
			}
		})
	}
})
