export const thingsAreEqual = function(a, b) {

	const easyTypes = ["string", "number", "boolean", "undefined"]

	if (easyTypes.includes(typeof a) || easyTypes.includes(typeof b)) {
		return a === b
	}

	// If we got this far, then both things have type "object".

	// null will never be strictly equal to anything but null.
	if (a === null || b === null) {
		return a === b
	}

	// Create an array of property names.
	let aKeys = Object.keys(a);

	// Check the number of properties in each thing.
	if (aKeys.length !== Object.keys(b).length) {
		return false;
	}

	for (let i = 0, l = aKeys.length; i < l; i++) {
		if (!thingsAreEqual(a[aKeys[i]], b[aKeys[i]])) {
			return false
		}
	}

	return true;

}
