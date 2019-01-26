export default {
	bkg_color: {
		name: "Background Color",
		type: "color",
		val: ""
	},
	bkg_img: {
		name: "Background Image",
		type: "image",
		val: ""
	},
	bkg_size: {
		name: "Background Image Size",
		type: "select",
		opt: [
			{
				name: "Cover",
				val: "cover"
			},
			{
				name: "Contain",
				val: "contain"
			},
			{
				name: "Auto",
				val: "auto"
			}
		],
		val: "cover"
	},
	color: {
		name: "Text Color",
		type: "color",
		val: ""
	},
	pad_top: {
		name: "Top Padding",
		type: "text",
		opt: /^[0-9]*$/,
		val: "",
		units: "px"
	},
	pad_bottom: {
		name: "Bottom Padding",
		type: "text",
		opt: /^[0-9]*$/,
		val: "",
		units: "px"
	},
	section_row_space: {
		name: "Spacing Between Rows",
		type: "select",
		opt: [
			{
				name: "No Spacing",
				val: "none"
			},
			{
				name: "Narrow",
				val: "narrow"
			},
			{
				name: "Normal",
				val: "normal"
			},
			{
				name: "Wide",
				val: "wide"
			}
		],
		val: "normal"
	},
	row_col_space: {
		name: "Spacing Between Columns",
		type: "select",
		opt: [
			{
				name: "No Spacing",
				val: "none"
			},
			{
				name: "Narrow",
				val: "narrow"
			},
			{
				name: "Normal",
				val: "normal"
			},
			{
				name: "Wide",
				val: "wide"
			}
		],
		val: "normal"
	},
	paragraph_margin_top: {
		name: "Top Margin of Paragraphs",
		type: "text",
		val: "",
		units: "px"
	},
	paragraph_margin_bottom: {
		name: "Bottom Margin of Paragraphs",
		type: "text",
		val: "",
		units: "px"
	},
	std_section_width: {
		name: "Standard Section Maximum Width",
		help: "Recommended range (px): 950&ndash;1200",
		err: "Enter only a number",
		opt: /^[0-9]+$/,
		val: "1100",
		units: "px"
	}
}
