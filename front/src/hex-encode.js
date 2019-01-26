export function hexEncode(u8arr) {
	let hex = ""
	u8arr.forEach(x => {
		hex += ("0" + x.toString(16)).slice(-2)
	})
	return hex
}