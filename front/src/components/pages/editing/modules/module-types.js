export default {
	page_body: {
		title: "Page Body",
		content_type: "page", // Used only for creating the type "page".
		icon: "description",
		style_props: ["use_site_width"]
	},
	article_body: {
		title: "Article Body",
		content_type: "post", // Used only for creating the type "post".
		icon: "description",
		style_props: ["use_site_width"]
	},
	text: {
		title: "Text",
		icon: "text_fields",
		default_opts: {
			"delta_ops": [{"insert": "\n"}],
			"body_html": "<p><br></p>"
		},
		component_name: "richtext-editor",
		style_props: ["paragraph_margin_top", "paragraph_margin_bottom", "color", "pad_top", "pad_bottom"]
	},
	html: {
		title: "HTML",
		icon: "code",
		default_opts: {
			"content": ""
		},
		component_name: "html-editor",
		style_props: ["pad_top", "pad_bottom"]
	},
	image: {
		title: "Image",
		icon: "insert_photo",
		default_body: {
			"src_type": "",
			"content": ""
		},
		component_name: "image-editor"
	}//,
	// big_heading: {
	// 	title: "Big Heading",
	// 	icon: "title",
	// 	default_body: {
	// 		"content": ""
	// 	}
	// },
	// blurb: {
	// 	title: "Blurb",
	// 	icon: "lightbulb_outline"
	// },
}
