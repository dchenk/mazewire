const base32Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"

export const base32 = function(buffer) {
	let bits = 0
	let value = 0
	let output = ""

	const length = buffer.byteLength
	for (let i = 0; i < length; i++) {
		value = (value << 8) | buffer[i]
		bits += 8
		while (bits >= 5) {
			output += base32Alphabet[(value >>> (bits - 5)) & 31]
			bits -= 5
		}
	}

	if (bits > 0) {
		output += base32Alphabet[(value << (5 - bits)) & 31]
	}

	return output
}
