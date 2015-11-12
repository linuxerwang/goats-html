goog.provide('goats.runtime.TagAttrs');

goog.require('goog.string');
goog.require('goog.object');


/**
 * Tag attrs.
 *
 * @constructor
 */
goats.runtime.TagAttrs = function() {
	/**
	 * The tag attrs.
	 * @private {Object}
	 */
	this.attrs_ = {};
};

/**
 * Add a tag attr.
 *
 * @param {string} name The name of the attr.
 * @param {string} value The value of the attr.
 */
goats.runtime.TagAttrs.prototype.add = function(name, value) {
	var v = this.attrs_[name];
	if (v) {
		if (name == "class") {
			this.attrs_[name] += " " + value;
		} else if (name == "style") {
			this.attrs_[name] += "; " + value;
		} else {
			this.attrs_[name] = value;
		}
	} else {
		this.attrs_[name] = value;
	}
};

/**
 * Merge tag attrs from another object.
 *
 * @param {Object} fromAttrs The object to merge attrs from.
 */
goats.runtime.TagAttrs.prototype.mergeFrom = function(fromAttrs) {
	goog.object.forEach(fromAttrs, function(element, index, obj) {
		this.add(index, element);
	}, this);
};

/**
 * Get the attrs.
 *
 * @return {Object} the attrs object.
 */
goats.runtime.TagAttrs.prototype.get = function() {
	return this.attrs_;
};
