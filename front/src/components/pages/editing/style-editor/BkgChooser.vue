<template>
	<div class="background-chooser">
		<div class="bkg-type-option">
			<div>
				<input class="material-radio" type="radio" value="color" id="bkgtype-color" name="bkg_type" v-model="bkgType">
				<label for="bkgtype-color"><span class="radio-circle material-icons"/>Color</label>
			</div>
			<div>
				<div class="dyn-textfield">
					<input type="text" class="dyn-textfield-input" id="bkg-color" v-model="bkgColor">
					<label for="bkg-color" class="dyn-textfield-label">Background Color</label>
				</div>
				<div class="input-helper-info" v-if="bkgErr">The value should include the unit type: #/rgb/rgba/hsl/hsla</div>
			</div>
		</div>
		<div class="bkg-type-option">
			<div>
				<input class="material-radio" type="radio" value="img" id="bkgtype-img" name="bkg_type" v-model="bkgType">
				<label for="bkgtype-img"><span class="radio-circle material-icons"/>Image</label>
			</div>
			<div>
				<div class="dyn-textfield">
					<input type="text" class="dyn-textfield-input" id="bkg-img" v-model="bkgImg">
					<label for="bkg-img" class="dyn-textfield-label">Background Image URL</label>
				</div>
				<label for="img-bkg-size">Image Size:</label>
				<select id="img-bkg-size" v-model="bkgSize">
					<option value="cover">Cover</option>
					<option value="contain">Contain Inside</option>
					<option value="auto">Auto</option>
				</select>
			</div>
		</div>
	</div>
</template>

<script>

	export default {
		name: "BkgChooser",
		props: {
			"styles": {
				type: Object,
				default: () => {}
			}
		},
		data() { return {
			bkgType: "color",
			bkgColor: "",
			bkgImg: "",
			bkgSize: "",
			bkgErr: ""
		}},
		watch: {
			styles: {
				handler: function(newVal) {
					if (newVal.bkg_type !== undefined && newVal.bkg_type !== "") {
						this.bkgType = newVal.bkg_type
					}
					if (newVal.bkg_color !== undefined && newVal.bkg_color !== "") {
						this.bkgColor = newVal.bkg_color
					}
					if (newVal.bkg_img !== undefined && newVal.bkg_img !== "") {
						this.bkgImg = newVal.bkg_img
					}
					if (newVal.bkg_size !== undefined && newVal.bkg_size !== "") {
						this.bkgSize = newVal.bkg_size
					}
				},
				deep: true
			},
			bkgColor: function(newBkg) {
				this.bkgErr = newBkg !== undefined && newBkg !== "" && newBkg.indexOf("#") === -1 && newBkg.indexOf("rgb") === -1 && newBkg.indexOf("rgba") === -1 && newBkg.indexOf("hsl") === -1 && newBkg.indexOf("hsla") === -1
			}
		},
		methods: {
			getStyles() {
			}
		}
	}

</script>

<style>

	.bkg-type-option {
		display: flex;
		align-items: center;
		border: 1px solid gray;
		margin-bottom: 8px;
		border-radius: 3px;
		padding: 10px 13px;
	}

	.bkg-type-option > div:first-child {
		width: 95px;
	}

	.bkg-type-option > div:last-child {
		flex-grow: 1;
	}

</style>
