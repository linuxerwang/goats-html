goog.provide('goats.runtime.filters');

goog.require('goog.object');
goog.require('goog.math');
goog.require('goog.string');
goog.require('goog.string.format');


/**
 * Outputs debug information to console.
 *
 * @param {*} input The input.
 */
goats.runtime.filters.debug = function(input) {
	window.console.log(input);
};

/**
 * Returns the length of the input.
 *
 * @param {Array|Object} input The input array or object.
 * @returns {number} the length of the input.
 */
goats.runtime.filters.length = function(input) {
	if (input) {
		if (input.length) {
			return input.length;
		} else {
			return goog.object.getCount(input);
		}
	}
	return 0;
};

/**
 * An alias of length().
 */
goats.runtime.filters.len = goats.runtime.filters.length;

/**
 * Converts the input to a title (capitalize the frist letter of each word).
 *
 * @param {string} input The input.
 * @returns {string} The title.
 */
goats.runtime.filters.title = function(input) {
	return goog.string.toTitleCase(input);
};

/**
 * Centers the input string with the given width.
 *
 * @param {string} input The input.
 * @param {number} width The width.
 * @returns {string} The result.
 */
goats.runtime.filters.center = function(input, width) {
	if (input.length >= width) {
		return input;
	}
	var extra = input.length;
	var left = extra / 2;
	var right = extra - left;

	return goats.runtime.filters.fillSpaces(left) + input +
		goats.runtime.filters.fillSpaces(right);
};

/**
 * Justifies the string to the left.
 *
 * @param {string} input The input.
 * @param {number} width The width.
 * @returns {string} The result.
 */
goats.runtime.filters.ljust = function(input, width) {
	return goats.runtime.filters.just(input, width, true /* left */);
};

/**
 * Justifies the string to the right.
 *
 * @param {string} input The input.
 * @param {number} width The width.
 * @returns {string} The result.
 */
goats.runtime.filters.rjust = function(input, width) {
	return goats.runtime.filters.just(input, width, false /* left */);
};

/**
 * Cuts a substring from the input string.
 *
 * @param {string} input The input.
 * @param {string} removed The string to remove.
 * @returns {string} The result.
 */
goats.runtime.filters.cut = function(input, removed) {
	return goog.string.remove(input, removed);
};

/**
 * Joins an array of strings.
 *
 * @param {Array.<string>} input The array of strings.
 * @param {string} separator The separator.
 * @returns {string} The result.
 */
goats.runtime.filters.join = function(input, separator) {
	return input.join(separator);
};

/**
 * Formats the float number to string.
 *
 * @param {number} input The float number.
 * @param {number} precision The precision.
 * @returns {string} The result.
 */
goats.runtime.filters.floatformat = function(input, precision) {
	if (precision < 0) {
		precision = -precision;
	}
	return input.toFixed(precision);
};

/**
 * Quotes the input string.
 *
 * @param {string} input The input.
 * @returns {string} The result.
 */
goats.runtime.filters.quote = function(input) {
	return goog.string.quote(input);
};

/**
 * Formats the input.
 *
 * @param {string} f The format string.
 * @param {...string|number} var_args Values formatString is to be filled with.
 * @returns {string} The result.
 */
goats.runtime.filters.format = function(f, var_args) {
	var args = Array.prototype.slice.call(arguments);
	return goog.string.format.apply(null, args);
};

/**
 * Parse integer from string or integer.
 *
 * @param {string} val the string or integer.
 * @returns {string} The result.
 */
goats.runtime.filters.parseInt = function(val) {
	if (goog.math.isInt(val)) {
		return val;
	}
        return goog.string.parseInt(val);
};

// ================ utility functions ================


/**
 * @const
 * @private
 */
var SPACES = '                                                                                                    ';

/**
 * Generate space string with the given width.
 *
 * @param {number} width The width.
 * @returns {number|string} The result.
 */
goats.runtime.filters.fillSpaces = function(width) {
	if (width == SPACES.length) {
		return SPACES;
	} else if (width < SPACES.length) {
		return SPACES.slice(0, width);
	}

	var n = width / SPACES.length;
	var r = width % SPACES.length;
	var s = "";
	for (var i = 0; i < n; i++) {
		// It's not optimal in complexity, but in practice it's very rare
		// to have to create such a long empty string.
		s += SPACES;
	}
	if (r > 0) {
		s =+ goats.runtime.filters.fillSpaces(r);
	}
	return s;
};

/**
 * Justifies the string to left or right.
 *
 * @param {string} input The input.
 * @param {number} width The width.
 * @param {boolean} left Whether justify to left.
 * @returns {string} The result.
 */
goats.runtime.filters.just = function(input, width, left) {
	if (input.length < width) {
		var extra = width - input.length;
		if (left) {
			input += goats.runtime.filters.fillSpaces(extra);
		} else {
			input = goats.runtime.filters.fillSpaces(extra) + input;
		}
	}
	return input;
};
