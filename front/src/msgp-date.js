import {Int64BE} from "int64-buffer"

export const dateExtPrefix = 0x05

const extensionLength = 12 // The first 8 bytes are seconds, the last 4 bytes are nanoseconds.

export function msgpUnpackDate(u8arr) {
	if (u8arr.length < extensionLength) {
		throw new Error("msgp-date: got a length less than 12 for timestamp data")
	}
	const big = new Int64BE(u8arr.slice(0, 8))
	const seconds = big.toNumber()
	const d = new Date()
	d.setTime(seconds*1000)
	return d
}

export function msgpPackDate(date) {
	const seconds = Math.ceil(date.getTime()/1000)
	const secondsBuffer = (new Int64BE(seconds)).toArrayBuffer()
	const arr = new Uint8Array(extensionLength)
	arr.set(secondsBuffer, 4)
	return arr
}
