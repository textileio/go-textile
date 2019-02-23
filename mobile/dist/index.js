/*eslint-disable block-scoped-var, id-length, no-control-regex, no-magic-numbers, no-prototype-builtins, no-redeclare, no-shadow, no-var, sort-vars*/
import * as $protobuf from "protobufjs/minimal";

// Common aliases
const $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;

// Exported root namespace
const $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});

export const CafeChallenge = $root.CafeChallenge = (() => {

    /**
     * Properties of a CafeChallenge.
     * @exports ICafeChallenge
     * @interface ICafeChallenge
     * @property {string} address CafeChallenge address
     */

    /**
     * Constructs a new CafeChallenge.
     * @exports CafeChallenge
     * @classdesc Represents a CafeChallenge.
     * @implements ICafeChallenge
     * @constructor
     * @param {ICafeChallenge=} [properties] Properties to set
     */
    function CafeChallenge(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeChallenge address.
     * @member {string} address
     * @memberof CafeChallenge
     * @instance
     */
    CafeChallenge.prototype.address = "";

    /**
     * Creates a new CafeChallenge instance using the specified properties.
     * @function create
     * @memberof CafeChallenge
     * @static
     * @param {ICafeChallenge=} [properties] Properties to set
     * @returns {CafeChallenge} CafeChallenge instance
     */
    CafeChallenge.create = function create(properties) {
        return new CafeChallenge(properties);
    };

    /**
     * Encodes the specified CafeChallenge message. Does not implicitly {@link CafeChallenge.verify|verify} messages.
     * @function encode
     * @memberof CafeChallenge
     * @static
     * @param {ICafeChallenge} message CafeChallenge message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeChallenge.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.address != null && message.hasOwnProperty("address"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.address);
        return writer;
    };

    /**
     * Encodes the specified CafeChallenge message, length delimited. Does not implicitly {@link CafeChallenge.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeChallenge
     * @static
     * @param {ICafeChallenge} message CafeChallenge message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeChallenge.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeChallenge message from the specified reader or buffer.
     * @function decode
     * @memberof CafeChallenge
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeChallenge} CafeChallenge
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeChallenge.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeChallenge();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.address = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeChallenge message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeChallenge
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeChallenge} CafeChallenge
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeChallenge.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeChallenge message.
     * @function verify
     * @memberof CafeChallenge
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeChallenge.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.address != null && message.hasOwnProperty("address"))
            if (!$util.isString(message.address))
                return "address: string expected";
        return null;
    };

    /**
     * Creates a CafeChallenge message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeChallenge
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeChallenge} CafeChallenge
     */
    CafeChallenge.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeChallenge)
            return object;
        let message = new $root.CafeChallenge();
        if (object.address != null)
            message.address = String(object.address);
        return message;
    };

    /**
     * Creates a plain object from a CafeChallenge message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeChallenge
     * @static
     * @param {CafeChallenge} message CafeChallenge
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeChallenge.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults)
            object.address = "";
        if (message.address != null && message.hasOwnProperty("address"))
            object.address = message.address;
        return object;
    };

    /**
     * Converts this CafeChallenge to JSON.
     * @function toJSON
     * @memberof CafeChallenge
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeChallenge.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeChallenge;
})();

export const CafeNonce = $root.CafeNonce = (() => {

    /**
     * Properties of a CafeNonce.
     * @exports ICafeNonce
     * @interface ICafeNonce
     * @property {string} value CafeNonce value
     */

    /**
     * Constructs a new CafeNonce.
     * @exports CafeNonce
     * @classdesc Represents a CafeNonce.
     * @implements ICafeNonce
     * @constructor
     * @param {ICafeNonce=} [properties] Properties to set
     */
    function CafeNonce(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeNonce value.
     * @member {string} value
     * @memberof CafeNonce
     * @instance
     */
    CafeNonce.prototype.value = "";

    /**
     * Creates a new CafeNonce instance using the specified properties.
     * @function create
     * @memberof CafeNonce
     * @static
     * @param {ICafeNonce=} [properties] Properties to set
     * @returns {CafeNonce} CafeNonce instance
     */
    CafeNonce.create = function create(properties) {
        return new CafeNonce(properties);
    };

    /**
     * Encodes the specified CafeNonce message. Does not implicitly {@link CafeNonce.verify|verify} messages.
     * @function encode
     * @memberof CafeNonce
     * @static
     * @param {ICafeNonce} message CafeNonce message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeNonce.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.value != null && message.hasOwnProperty("value"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.value);
        return writer;
    };

    /**
     * Encodes the specified CafeNonce message, length delimited. Does not implicitly {@link CafeNonce.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeNonce
     * @static
     * @param {ICafeNonce} message CafeNonce message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeNonce.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeNonce message from the specified reader or buffer.
     * @function decode
     * @memberof CafeNonce
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeNonce} CafeNonce
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeNonce.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeNonce();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.value = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeNonce message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeNonce
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeNonce} CafeNonce
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeNonce.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeNonce message.
     * @function verify
     * @memberof CafeNonce
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeNonce.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.value != null && message.hasOwnProperty("value"))
            if (!$util.isString(message.value))
                return "value: string expected";
        return null;
    };

    /**
     * Creates a CafeNonce message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeNonce
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeNonce} CafeNonce
     */
    CafeNonce.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeNonce)
            return object;
        let message = new $root.CafeNonce();
        if (object.value != null)
            message.value = String(object.value);
        return message;
    };

    /**
     * Creates a plain object from a CafeNonce message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeNonce
     * @static
     * @param {CafeNonce} message CafeNonce
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeNonce.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults)
            object.value = "";
        if (message.value != null && message.hasOwnProperty("value"))
            object.value = message.value;
        return object;
    };

    /**
     * Converts this CafeNonce to JSON.
     * @function toJSON
     * @memberof CafeNonce
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeNonce.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeNonce;
})();

export const CafeRegistration = $root.CafeRegistration = (() => {

    /**
     * Properties of a CafeRegistration.
     * @exports ICafeRegistration
     * @interface ICafeRegistration
     * @property {string} address CafeRegistration address
     * @property {string} value CafeRegistration value
     * @property {string} nonce CafeRegistration nonce
     * @property {Uint8Array} sig CafeRegistration sig
     * @property {string} token CafeRegistration token
     */

    /**
     * Constructs a new CafeRegistration.
     * @exports CafeRegistration
     * @classdesc Represents a CafeRegistration.
     * @implements ICafeRegistration
     * @constructor
     * @param {ICafeRegistration=} [properties] Properties to set
     */
    function CafeRegistration(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeRegistration address.
     * @member {string} address
     * @memberof CafeRegistration
     * @instance
     */
    CafeRegistration.prototype.address = "";

    /**
     * CafeRegistration value.
     * @member {string} value
     * @memberof CafeRegistration
     * @instance
     */
    CafeRegistration.prototype.value = "";

    /**
     * CafeRegistration nonce.
     * @member {string} nonce
     * @memberof CafeRegistration
     * @instance
     */
    CafeRegistration.prototype.nonce = "";

    /**
     * CafeRegistration sig.
     * @member {Uint8Array} sig
     * @memberof CafeRegistration
     * @instance
     */
    CafeRegistration.prototype.sig = $util.newBuffer([]);

    /**
     * CafeRegistration token.
     * @member {string} token
     * @memberof CafeRegistration
     * @instance
     */
    CafeRegistration.prototype.token = "";

    /**
     * Creates a new CafeRegistration instance using the specified properties.
     * @function create
     * @memberof CafeRegistration
     * @static
     * @param {ICafeRegistration=} [properties] Properties to set
     * @returns {CafeRegistration} CafeRegistration instance
     */
    CafeRegistration.create = function create(properties) {
        return new CafeRegistration(properties);
    };

    /**
     * Encodes the specified CafeRegistration message. Does not implicitly {@link CafeRegistration.verify|verify} messages.
     * @function encode
     * @memberof CafeRegistration
     * @static
     * @param {ICafeRegistration} message CafeRegistration message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeRegistration.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.address != null && message.hasOwnProperty("address"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.address);
        if (message.value != null && message.hasOwnProperty("value"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.value);
        if (message.nonce != null && message.hasOwnProperty("nonce"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.nonce);
        if (message.sig != null && message.hasOwnProperty("sig"))
            writer.uint32(/* id 4, wireType 2 =*/34).bytes(message.sig);
        if (message.token != null && message.hasOwnProperty("token"))
            writer.uint32(/* id 5, wireType 2 =*/42).string(message.token);
        return writer;
    };

    /**
     * Encodes the specified CafeRegistration message, length delimited. Does not implicitly {@link CafeRegistration.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeRegistration
     * @static
     * @param {ICafeRegistration} message CafeRegistration message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeRegistration.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeRegistration message from the specified reader or buffer.
     * @function decode
     * @memberof CafeRegistration
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeRegistration} CafeRegistration
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeRegistration.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeRegistration();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.address = reader.string();
                break;
            case 2:
                message.value = reader.string();
                break;
            case 3:
                message.nonce = reader.string();
                break;
            case 4:
                message.sig = reader.bytes();
                break;
            case 5:
                message.token = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeRegistration message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeRegistration
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeRegistration} CafeRegistration
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeRegistration.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeRegistration message.
     * @function verify
     * @memberof CafeRegistration
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeRegistration.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.address != null && message.hasOwnProperty("address"))
            if (!$util.isString(message.address))
                return "address: string expected";
        if (message.value != null && message.hasOwnProperty("value"))
            if (!$util.isString(message.value))
                return "value: string expected";
        if (message.nonce != null && message.hasOwnProperty("nonce"))
            if (!$util.isString(message.nonce))
                return "nonce: string expected";
        if (message.sig != null && message.hasOwnProperty("sig"))
            if (!(message.sig && typeof message.sig.length === "number" || $util.isString(message.sig)))
                return "sig: buffer expected";
        if (message.token != null && message.hasOwnProperty("token"))
            if (!$util.isString(message.token))
                return "token: string expected";
        return null;
    };

    /**
     * Creates a CafeRegistration message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeRegistration
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeRegistration} CafeRegistration
     */
    CafeRegistration.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeRegistration)
            return object;
        let message = new $root.CafeRegistration();
        if (object.address != null)
            message.address = String(object.address);
        if (object.value != null)
            message.value = String(object.value);
        if (object.nonce != null)
            message.nonce = String(object.nonce);
        if (object.sig != null)
            if (typeof object.sig === "string")
                $util.base64.decode(object.sig, message.sig = $util.newBuffer($util.base64.length(object.sig)), 0);
            else if (object.sig.length)
                message.sig = object.sig;
        if (object.token != null)
            message.token = String(object.token);
        return message;
    };

    /**
     * Creates a plain object from a CafeRegistration message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeRegistration
     * @static
     * @param {CafeRegistration} message CafeRegistration
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeRegistration.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.address = "";
            object.value = "";
            object.nonce = "";
            if (options.bytes === String)
                object.sig = "";
            else {
                object.sig = [];
                if (options.bytes !== Array)
                    object.sig = $util.newBuffer(object.sig);
            }
            object.token = "";
        }
        if (message.address != null && message.hasOwnProperty("address"))
            object.address = message.address;
        if (message.value != null && message.hasOwnProperty("value"))
            object.value = message.value;
        if (message.nonce != null && message.hasOwnProperty("nonce"))
            object.nonce = message.nonce;
        if (message.sig != null && message.hasOwnProperty("sig"))
            object.sig = options.bytes === String ? $util.base64.encode(message.sig, 0, message.sig.length) : options.bytes === Array ? Array.prototype.slice.call(message.sig) : message.sig;
        if (message.token != null && message.hasOwnProperty("token"))
            object.token = message.token;
        return object;
    };

    /**
     * Converts this CafeRegistration to JSON.
     * @function toJSON
     * @memberof CafeRegistration
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeRegistration.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeRegistration;
})();

export const CafeRefreshSession = $root.CafeRefreshSession = (() => {

    /**
     * Properties of a CafeRefreshSession.
     * @exports ICafeRefreshSession
     * @interface ICafeRefreshSession
     * @property {string} access CafeRefreshSession access
     * @property {string} refresh CafeRefreshSession refresh
     */

    /**
     * Constructs a new CafeRefreshSession.
     * @exports CafeRefreshSession
     * @classdesc Represents a CafeRefreshSession.
     * @implements ICafeRefreshSession
     * @constructor
     * @param {ICafeRefreshSession=} [properties] Properties to set
     */
    function CafeRefreshSession(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeRefreshSession access.
     * @member {string} access
     * @memberof CafeRefreshSession
     * @instance
     */
    CafeRefreshSession.prototype.access = "";

    /**
     * CafeRefreshSession refresh.
     * @member {string} refresh
     * @memberof CafeRefreshSession
     * @instance
     */
    CafeRefreshSession.prototype.refresh = "";

    /**
     * Creates a new CafeRefreshSession instance using the specified properties.
     * @function create
     * @memberof CafeRefreshSession
     * @static
     * @param {ICafeRefreshSession=} [properties] Properties to set
     * @returns {CafeRefreshSession} CafeRefreshSession instance
     */
    CafeRefreshSession.create = function create(properties) {
        return new CafeRefreshSession(properties);
    };

    /**
     * Encodes the specified CafeRefreshSession message. Does not implicitly {@link CafeRefreshSession.verify|verify} messages.
     * @function encode
     * @memberof CafeRefreshSession
     * @static
     * @param {ICafeRefreshSession} message CafeRefreshSession message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeRefreshSession.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.access != null && message.hasOwnProperty("access"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.access);
        if (message.refresh != null && message.hasOwnProperty("refresh"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.refresh);
        return writer;
    };

    /**
     * Encodes the specified CafeRefreshSession message, length delimited. Does not implicitly {@link CafeRefreshSession.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeRefreshSession
     * @static
     * @param {ICafeRefreshSession} message CafeRefreshSession message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeRefreshSession.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeRefreshSession message from the specified reader or buffer.
     * @function decode
     * @memberof CafeRefreshSession
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeRefreshSession} CafeRefreshSession
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeRefreshSession.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeRefreshSession();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.access = reader.string();
                break;
            case 2:
                message.refresh = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeRefreshSession message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeRefreshSession
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeRefreshSession} CafeRefreshSession
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeRefreshSession.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeRefreshSession message.
     * @function verify
     * @memberof CafeRefreshSession
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeRefreshSession.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.access != null && message.hasOwnProperty("access"))
            if (!$util.isString(message.access))
                return "access: string expected";
        if (message.refresh != null && message.hasOwnProperty("refresh"))
            if (!$util.isString(message.refresh))
                return "refresh: string expected";
        return null;
    };

    /**
     * Creates a CafeRefreshSession message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeRefreshSession
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeRefreshSession} CafeRefreshSession
     */
    CafeRefreshSession.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeRefreshSession)
            return object;
        let message = new $root.CafeRefreshSession();
        if (object.access != null)
            message.access = String(object.access);
        if (object.refresh != null)
            message.refresh = String(object.refresh);
        return message;
    };

    /**
     * Creates a plain object from a CafeRefreshSession message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeRefreshSession
     * @static
     * @param {CafeRefreshSession} message CafeRefreshSession
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeRefreshSession.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.access = "";
            object.refresh = "";
        }
        if (message.access != null && message.hasOwnProperty("access"))
            object.access = message.access;
        if (message.refresh != null && message.hasOwnProperty("refresh"))
            object.refresh = message.refresh;
        return object;
    };

    /**
     * Converts this CafeRefreshSession to JSON.
     * @function toJSON
     * @memberof CafeRefreshSession
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeRefreshSession.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeRefreshSession;
})();

export const CafePublishContact = $root.CafePublishContact = (() => {

    /**
     * Properties of a CafePublishContact.
     * @exports ICafePublishContact
     * @interface ICafePublishContact
     * @property {string} token CafePublishContact token
     * @property {IContact} contact CafePublishContact contact
     */

    /**
     * Constructs a new CafePublishContact.
     * @exports CafePublishContact
     * @classdesc Represents a CafePublishContact.
     * @implements ICafePublishContact
     * @constructor
     * @param {ICafePublishContact=} [properties] Properties to set
     */
    function CafePublishContact(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafePublishContact token.
     * @member {string} token
     * @memberof CafePublishContact
     * @instance
     */
    CafePublishContact.prototype.token = "";

    /**
     * CafePublishContact contact.
     * @member {IContact} contact
     * @memberof CafePublishContact
     * @instance
     */
    CafePublishContact.prototype.contact = null;

    /**
     * Creates a new CafePublishContact instance using the specified properties.
     * @function create
     * @memberof CafePublishContact
     * @static
     * @param {ICafePublishContact=} [properties] Properties to set
     * @returns {CafePublishContact} CafePublishContact instance
     */
    CafePublishContact.create = function create(properties) {
        return new CafePublishContact(properties);
    };

    /**
     * Encodes the specified CafePublishContact message. Does not implicitly {@link CafePublishContact.verify|verify} messages.
     * @function encode
     * @memberof CafePublishContact
     * @static
     * @param {ICafePublishContact} message CafePublishContact message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafePublishContact.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.token != null && message.hasOwnProperty("token"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.token);
        if (message.contact != null && message.hasOwnProperty("contact"))
            $root.Contact.encode(message.contact, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified CafePublishContact message, length delimited. Does not implicitly {@link CafePublishContact.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafePublishContact
     * @static
     * @param {ICafePublishContact} message CafePublishContact message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafePublishContact.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafePublishContact message from the specified reader or buffer.
     * @function decode
     * @memberof CafePublishContact
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafePublishContact} CafePublishContact
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafePublishContact.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafePublishContact();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.token = reader.string();
                break;
            case 2:
                message.contact = $root.Contact.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafePublishContact message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafePublishContact
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafePublishContact} CafePublishContact
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafePublishContact.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafePublishContact message.
     * @function verify
     * @memberof CafePublishContact
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafePublishContact.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.token != null && message.hasOwnProperty("token"))
            if (!$util.isString(message.token))
                return "token: string expected";
        if (message.contact != null && message.hasOwnProperty("contact")) {
            let error = $root.Contact.verify(message.contact);
            if (error)
                return "contact." + error;
        }
        return null;
    };

    /**
     * Creates a CafePublishContact message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafePublishContact
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafePublishContact} CafePublishContact
     */
    CafePublishContact.fromObject = function fromObject(object) {
        if (object instanceof $root.CafePublishContact)
            return object;
        let message = new $root.CafePublishContact();
        if (object.token != null)
            message.token = String(object.token);
        if (object.contact != null) {
            if (typeof object.contact !== "object")
                throw TypeError(".CafePublishContact.contact: object expected");
            message.contact = $root.Contact.fromObject(object.contact);
        }
        return message;
    };

    /**
     * Creates a plain object from a CafePublishContact message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafePublishContact
     * @static
     * @param {CafePublishContact} message CafePublishContact
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafePublishContact.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.token = "";
            object.contact = null;
        }
        if (message.token != null && message.hasOwnProperty("token"))
            object.token = message.token;
        if (message.contact != null && message.hasOwnProperty("contact"))
            object.contact = $root.Contact.toObject(message.contact, options);
        return object;
    };

    /**
     * Converts this CafePublishContact to JSON.
     * @function toJSON
     * @memberof CafePublishContact
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafePublishContact.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafePublishContact;
})();

export const CafePublishContactAck = $root.CafePublishContactAck = (() => {

    /**
     * Properties of a CafePublishContactAck.
     * @exports ICafePublishContactAck
     * @interface ICafePublishContactAck
     * @property {string} id CafePublishContactAck id
     */

    /**
     * Constructs a new CafePublishContactAck.
     * @exports CafePublishContactAck
     * @classdesc Represents a CafePublishContactAck.
     * @implements ICafePublishContactAck
     * @constructor
     * @param {ICafePublishContactAck=} [properties] Properties to set
     */
    function CafePublishContactAck(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafePublishContactAck id.
     * @member {string} id
     * @memberof CafePublishContactAck
     * @instance
     */
    CafePublishContactAck.prototype.id = "";

    /**
     * Creates a new CafePublishContactAck instance using the specified properties.
     * @function create
     * @memberof CafePublishContactAck
     * @static
     * @param {ICafePublishContactAck=} [properties] Properties to set
     * @returns {CafePublishContactAck} CafePublishContactAck instance
     */
    CafePublishContactAck.create = function create(properties) {
        return new CafePublishContactAck(properties);
    };

    /**
     * Encodes the specified CafePublishContactAck message. Does not implicitly {@link CafePublishContactAck.verify|verify} messages.
     * @function encode
     * @memberof CafePublishContactAck
     * @static
     * @param {ICafePublishContactAck} message CafePublishContactAck message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafePublishContactAck.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        return writer;
    };

    /**
     * Encodes the specified CafePublishContactAck message, length delimited. Does not implicitly {@link CafePublishContactAck.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafePublishContactAck
     * @static
     * @param {ICafePublishContactAck} message CafePublishContactAck message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafePublishContactAck.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafePublishContactAck message from the specified reader or buffer.
     * @function decode
     * @memberof CafePublishContactAck
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafePublishContactAck} CafePublishContactAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafePublishContactAck.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafePublishContactAck();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafePublishContactAck message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafePublishContactAck
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafePublishContactAck} CafePublishContactAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafePublishContactAck.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafePublishContactAck message.
     * @function verify
     * @memberof CafePublishContactAck
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafePublishContactAck.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        return null;
    };

    /**
     * Creates a CafePublishContactAck message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafePublishContactAck
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafePublishContactAck} CafePublishContactAck
     */
    CafePublishContactAck.fromObject = function fromObject(object) {
        if (object instanceof $root.CafePublishContactAck)
            return object;
        let message = new $root.CafePublishContactAck();
        if (object.id != null)
            message.id = String(object.id);
        return message;
    };

    /**
     * Creates a plain object from a CafePublishContactAck message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafePublishContactAck
     * @static
     * @param {CafePublishContactAck} message CafePublishContactAck
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafePublishContactAck.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults)
            object.id = "";
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        return object;
    };

    /**
     * Converts this CafePublishContactAck to JSON.
     * @function toJSON
     * @memberof CafePublishContactAck
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafePublishContactAck.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafePublishContactAck;
})();

export const CafeStore = $root.CafeStore = (() => {

    /**
     * Properties of a CafeStore.
     * @exports ICafeStore
     * @interface ICafeStore
     * @property {string} token CafeStore token
     * @property {Array.<string>} cids CafeStore cids
     */

    /**
     * Constructs a new CafeStore.
     * @exports CafeStore
     * @classdesc Represents a CafeStore.
     * @implements ICafeStore
     * @constructor
     * @param {ICafeStore=} [properties] Properties to set
     */
    function CafeStore(properties) {
        this.cids = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeStore token.
     * @member {string} token
     * @memberof CafeStore
     * @instance
     */
    CafeStore.prototype.token = "";

    /**
     * CafeStore cids.
     * @member {Array.<string>} cids
     * @memberof CafeStore
     * @instance
     */
    CafeStore.prototype.cids = $util.emptyArray;

    /**
     * Creates a new CafeStore instance using the specified properties.
     * @function create
     * @memberof CafeStore
     * @static
     * @param {ICafeStore=} [properties] Properties to set
     * @returns {CafeStore} CafeStore instance
     */
    CafeStore.create = function create(properties) {
        return new CafeStore(properties);
    };

    /**
     * Encodes the specified CafeStore message. Does not implicitly {@link CafeStore.verify|verify} messages.
     * @function encode
     * @memberof CafeStore
     * @static
     * @param {ICafeStore} message CafeStore message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeStore.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.token != null && message.hasOwnProperty("token"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.token);
        if (message.cids != null && message.cids.length)
            for (let i = 0; i < message.cids.length; ++i)
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.cids[i]);
        return writer;
    };

    /**
     * Encodes the specified CafeStore message, length delimited. Does not implicitly {@link CafeStore.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeStore
     * @static
     * @param {ICafeStore} message CafeStore message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeStore.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeStore message from the specified reader or buffer.
     * @function decode
     * @memberof CafeStore
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeStore} CafeStore
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeStore.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeStore();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.token = reader.string();
                break;
            case 2:
                if (!(message.cids && message.cids.length))
                    message.cids = [];
                message.cids.push(reader.string());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeStore message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeStore
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeStore} CafeStore
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeStore.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeStore message.
     * @function verify
     * @memberof CafeStore
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeStore.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.token != null && message.hasOwnProperty("token"))
            if (!$util.isString(message.token))
                return "token: string expected";
        if (message.cids != null && message.hasOwnProperty("cids")) {
            if (!Array.isArray(message.cids))
                return "cids: array expected";
            for (let i = 0; i < message.cids.length; ++i)
                if (!$util.isString(message.cids[i]))
                    return "cids: string[] expected";
        }
        return null;
    };

    /**
     * Creates a CafeStore message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeStore
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeStore} CafeStore
     */
    CafeStore.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeStore)
            return object;
        let message = new $root.CafeStore();
        if (object.token != null)
            message.token = String(object.token);
        if (object.cids) {
            if (!Array.isArray(object.cids))
                throw TypeError(".CafeStore.cids: array expected");
            message.cids = [];
            for (let i = 0; i < object.cids.length; ++i)
                message.cids[i] = String(object.cids[i]);
        }
        return message;
    };

    /**
     * Creates a plain object from a CafeStore message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeStore
     * @static
     * @param {CafeStore} message CafeStore
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeStore.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.cids = [];
        if (options.defaults)
            object.token = "";
        if (message.token != null && message.hasOwnProperty("token"))
            object.token = message.token;
        if (message.cids && message.cids.length) {
            object.cids = [];
            for (let j = 0; j < message.cids.length; ++j)
                object.cids[j] = message.cids[j];
        }
        return object;
    };

    /**
     * Converts this CafeStore to JSON.
     * @function toJSON
     * @memberof CafeStore
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeStore.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeStore;
})();

export const CafeObjectList = $root.CafeObjectList = (() => {

    /**
     * Properties of a CafeObjectList.
     * @exports ICafeObjectList
     * @interface ICafeObjectList
     * @property {Array.<string>} cids CafeObjectList cids
     */

    /**
     * Constructs a new CafeObjectList.
     * @exports CafeObjectList
     * @classdesc Represents a CafeObjectList.
     * @implements ICafeObjectList
     * @constructor
     * @param {ICafeObjectList=} [properties] Properties to set
     */
    function CafeObjectList(properties) {
        this.cids = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeObjectList cids.
     * @member {Array.<string>} cids
     * @memberof CafeObjectList
     * @instance
     */
    CafeObjectList.prototype.cids = $util.emptyArray;

    /**
     * Creates a new CafeObjectList instance using the specified properties.
     * @function create
     * @memberof CafeObjectList
     * @static
     * @param {ICafeObjectList=} [properties] Properties to set
     * @returns {CafeObjectList} CafeObjectList instance
     */
    CafeObjectList.create = function create(properties) {
        return new CafeObjectList(properties);
    };

    /**
     * Encodes the specified CafeObjectList message. Does not implicitly {@link CafeObjectList.verify|verify} messages.
     * @function encode
     * @memberof CafeObjectList
     * @static
     * @param {ICafeObjectList} message CafeObjectList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeObjectList.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.cids != null && message.cids.length)
            for (let i = 0; i < message.cids.length; ++i)
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.cids[i]);
        return writer;
    };

    /**
     * Encodes the specified CafeObjectList message, length delimited. Does not implicitly {@link CafeObjectList.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeObjectList
     * @static
     * @param {ICafeObjectList} message CafeObjectList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeObjectList.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeObjectList message from the specified reader or buffer.
     * @function decode
     * @memberof CafeObjectList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeObjectList} CafeObjectList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeObjectList.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeObjectList();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                if (!(message.cids && message.cids.length))
                    message.cids = [];
                message.cids.push(reader.string());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeObjectList message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeObjectList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeObjectList} CafeObjectList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeObjectList.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeObjectList message.
     * @function verify
     * @memberof CafeObjectList
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeObjectList.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.cids != null && message.hasOwnProperty("cids")) {
            if (!Array.isArray(message.cids))
                return "cids: array expected";
            for (let i = 0; i < message.cids.length; ++i)
                if (!$util.isString(message.cids[i]))
                    return "cids: string[] expected";
        }
        return null;
    };

    /**
     * Creates a CafeObjectList message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeObjectList
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeObjectList} CafeObjectList
     */
    CafeObjectList.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeObjectList)
            return object;
        let message = new $root.CafeObjectList();
        if (object.cids) {
            if (!Array.isArray(object.cids))
                throw TypeError(".CafeObjectList.cids: array expected");
            message.cids = [];
            for (let i = 0; i < object.cids.length; ++i)
                message.cids[i] = String(object.cids[i]);
        }
        return message;
    };

    /**
     * Creates a plain object from a CafeObjectList message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeObjectList
     * @static
     * @param {CafeObjectList} message CafeObjectList
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeObjectList.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.cids = [];
        if (message.cids && message.cids.length) {
            object.cids = [];
            for (let j = 0; j < message.cids.length; ++j)
                object.cids[j] = message.cids[j];
        }
        return object;
    };

    /**
     * Converts this CafeObjectList to JSON.
     * @function toJSON
     * @memberof CafeObjectList
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeObjectList.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeObjectList;
})();

export const CafeObject = $root.CafeObject = (() => {

    /**
     * Properties of a CafeObject.
     * @exports ICafeObject
     * @interface ICafeObject
     * @property {string} token CafeObject token
     * @property {string} cid CafeObject cid
     * @property {Uint8Array} data CafeObject data
     * @property {Uint8Array} node CafeObject node
     */

    /**
     * Constructs a new CafeObject.
     * @exports CafeObject
     * @classdesc Represents a CafeObject.
     * @implements ICafeObject
     * @constructor
     * @param {ICafeObject=} [properties] Properties to set
     */
    function CafeObject(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeObject token.
     * @member {string} token
     * @memberof CafeObject
     * @instance
     */
    CafeObject.prototype.token = "";

    /**
     * CafeObject cid.
     * @member {string} cid
     * @memberof CafeObject
     * @instance
     */
    CafeObject.prototype.cid = "";

    /**
     * CafeObject data.
     * @member {Uint8Array} data
     * @memberof CafeObject
     * @instance
     */
    CafeObject.prototype.data = $util.newBuffer([]);

    /**
     * CafeObject node.
     * @member {Uint8Array} node
     * @memberof CafeObject
     * @instance
     */
    CafeObject.prototype.node = $util.newBuffer([]);

    /**
     * Creates a new CafeObject instance using the specified properties.
     * @function create
     * @memberof CafeObject
     * @static
     * @param {ICafeObject=} [properties] Properties to set
     * @returns {CafeObject} CafeObject instance
     */
    CafeObject.create = function create(properties) {
        return new CafeObject(properties);
    };

    /**
     * Encodes the specified CafeObject message. Does not implicitly {@link CafeObject.verify|verify} messages.
     * @function encode
     * @memberof CafeObject
     * @static
     * @param {ICafeObject} message CafeObject message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeObject.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.token != null && message.hasOwnProperty("token"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.token);
        if (message.cid != null && message.hasOwnProperty("cid"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.cid);
        if (message.data != null && message.hasOwnProperty("data"))
            writer.uint32(/* id 3, wireType 2 =*/26).bytes(message.data);
        if (message.node != null && message.hasOwnProperty("node"))
            writer.uint32(/* id 4, wireType 2 =*/34).bytes(message.node);
        return writer;
    };

    /**
     * Encodes the specified CafeObject message, length delimited. Does not implicitly {@link CafeObject.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeObject
     * @static
     * @param {ICafeObject} message CafeObject message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeObject.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeObject message from the specified reader or buffer.
     * @function decode
     * @memberof CafeObject
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeObject} CafeObject
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeObject.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeObject();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.token = reader.string();
                break;
            case 2:
                message.cid = reader.string();
                break;
            case 3:
                message.data = reader.bytes();
                break;
            case 4:
                message.node = reader.bytes();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeObject message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeObject
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeObject} CafeObject
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeObject.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeObject message.
     * @function verify
     * @memberof CafeObject
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeObject.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.token != null && message.hasOwnProperty("token"))
            if (!$util.isString(message.token))
                return "token: string expected";
        if (message.cid != null && message.hasOwnProperty("cid"))
            if (!$util.isString(message.cid))
                return "cid: string expected";
        if (message.data != null && message.hasOwnProperty("data"))
            if (!(message.data && typeof message.data.length === "number" || $util.isString(message.data)))
                return "data: buffer expected";
        if (message.node != null && message.hasOwnProperty("node"))
            if (!(message.node && typeof message.node.length === "number" || $util.isString(message.node)))
                return "node: buffer expected";
        return null;
    };

    /**
     * Creates a CafeObject message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeObject
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeObject} CafeObject
     */
    CafeObject.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeObject)
            return object;
        let message = new $root.CafeObject();
        if (object.token != null)
            message.token = String(object.token);
        if (object.cid != null)
            message.cid = String(object.cid);
        if (object.data != null)
            if (typeof object.data === "string")
                $util.base64.decode(object.data, message.data = $util.newBuffer($util.base64.length(object.data)), 0);
            else if (object.data.length)
                message.data = object.data;
        if (object.node != null)
            if (typeof object.node === "string")
                $util.base64.decode(object.node, message.node = $util.newBuffer($util.base64.length(object.node)), 0);
            else if (object.node.length)
                message.node = object.node;
        return message;
    };

    /**
     * Creates a plain object from a CafeObject message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeObject
     * @static
     * @param {CafeObject} message CafeObject
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeObject.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.token = "";
            object.cid = "";
            if (options.bytes === String)
                object.data = "";
            else {
                object.data = [];
                if (options.bytes !== Array)
                    object.data = $util.newBuffer(object.data);
            }
            if (options.bytes === String)
                object.node = "";
            else {
                object.node = [];
                if (options.bytes !== Array)
                    object.node = $util.newBuffer(object.node);
            }
        }
        if (message.token != null && message.hasOwnProperty("token"))
            object.token = message.token;
        if (message.cid != null && message.hasOwnProperty("cid"))
            object.cid = message.cid;
        if (message.data != null && message.hasOwnProperty("data"))
            object.data = options.bytes === String ? $util.base64.encode(message.data, 0, message.data.length) : options.bytes === Array ? Array.prototype.slice.call(message.data) : message.data;
        if (message.node != null && message.hasOwnProperty("node"))
            object.node = options.bytes === String ? $util.base64.encode(message.node, 0, message.node.length) : options.bytes === Array ? Array.prototype.slice.call(message.node) : message.node;
        return object;
    };

    /**
     * Converts this CafeObject to JSON.
     * @function toJSON
     * @memberof CafeObject
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeObject.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeObject;
})();

export const CafeStoreThread = $root.CafeStoreThread = (() => {

    /**
     * Properties of a CafeStoreThread.
     * @exports ICafeStoreThread
     * @interface ICafeStoreThread
     * @property {string} token CafeStoreThread token
     * @property {string} id CafeStoreThread id
     * @property {Uint8Array} ciphertext CafeStoreThread ciphertext
     */

    /**
     * Constructs a new CafeStoreThread.
     * @exports CafeStoreThread
     * @classdesc Represents a CafeStoreThread.
     * @implements ICafeStoreThread
     * @constructor
     * @param {ICafeStoreThread=} [properties] Properties to set
     */
    function CafeStoreThread(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeStoreThread token.
     * @member {string} token
     * @memberof CafeStoreThread
     * @instance
     */
    CafeStoreThread.prototype.token = "";

    /**
     * CafeStoreThread id.
     * @member {string} id
     * @memberof CafeStoreThread
     * @instance
     */
    CafeStoreThread.prototype.id = "";

    /**
     * CafeStoreThread ciphertext.
     * @member {Uint8Array} ciphertext
     * @memberof CafeStoreThread
     * @instance
     */
    CafeStoreThread.prototype.ciphertext = $util.newBuffer([]);

    /**
     * Creates a new CafeStoreThread instance using the specified properties.
     * @function create
     * @memberof CafeStoreThread
     * @static
     * @param {ICafeStoreThread=} [properties] Properties to set
     * @returns {CafeStoreThread} CafeStoreThread instance
     */
    CafeStoreThread.create = function create(properties) {
        return new CafeStoreThread(properties);
    };

    /**
     * Encodes the specified CafeStoreThread message. Does not implicitly {@link CafeStoreThread.verify|verify} messages.
     * @function encode
     * @memberof CafeStoreThread
     * @static
     * @param {ICafeStoreThread} message CafeStoreThread message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeStoreThread.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.token != null && message.hasOwnProperty("token"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.token);
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.id);
        if (message.ciphertext != null && message.hasOwnProperty("ciphertext"))
            writer.uint32(/* id 3, wireType 2 =*/26).bytes(message.ciphertext);
        return writer;
    };

    /**
     * Encodes the specified CafeStoreThread message, length delimited. Does not implicitly {@link CafeStoreThread.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeStoreThread
     * @static
     * @param {ICafeStoreThread} message CafeStoreThread message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeStoreThread.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeStoreThread message from the specified reader or buffer.
     * @function decode
     * @memberof CafeStoreThread
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeStoreThread} CafeStoreThread
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeStoreThread.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeStoreThread();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.token = reader.string();
                break;
            case 2:
                message.id = reader.string();
                break;
            case 3:
                message.ciphertext = reader.bytes();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeStoreThread message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeStoreThread
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeStoreThread} CafeStoreThread
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeStoreThread.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeStoreThread message.
     * @function verify
     * @memberof CafeStoreThread
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeStoreThread.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.token != null && message.hasOwnProperty("token"))
            if (!$util.isString(message.token))
                return "token: string expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.ciphertext != null && message.hasOwnProperty("ciphertext"))
            if (!(message.ciphertext && typeof message.ciphertext.length === "number" || $util.isString(message.ciphertext)))
                return "ciphertext: buffer expected";
        return null;
    };

    /**
     * Creates a CafeStoreThread message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeStoreThread
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeStoreThread} CafeStoreThread
     */
    CafeStoreThread.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeStoreThread)
            return object;
        let message = new $root.CafeStoreThread();
        if (object.token != null)
            message.token = String(object.token);
        if (object.id != null)
            message.id = String(object.id);
        if (object.ciphertext != null)
            if (typeof object.ciphertext === "string")
                $util.base64.decode(object.ciphertext, message.ciphertext = $util.newBuffer($util.base64.length(object.ciphertext)), 0);
            else if (object.ciphertext.length)
                message.ciphertext = object.ciphertext;
        return message;
    };

    /**
     * Creates a plain object from a CafeStoreThread message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeStoreThread
     * @static
     * @param {CafeStoreThread} message CafeStoreThread
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeStoreThread.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.token = "";
            object.id = "";
            if (options.bytes === String)
                object.ciphertext = "";
            else {
                object.ciphertext = [];
                if (options.bytes !== Array)
                    object.ciphertext = $util.newBuffer(object.ciphertext);
            }
        }
        if (message.token != null && message.hasOwnProperty("token"))
            object.token = message.token;
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.ciphertext != null && message.hasOwnProperty("ciphertext"))
            object.ciphertext = options.bytes === String ? $util.base64.encode(message.ciphertext, 0, message.ciphertext.length) : options.bytes === Array ? Array.prototype.slice.call(message.ciphertext) : message.ciphertext;
        return object;
    };

    /**
     * Converts this CafeStoreThread to JSON.
     * @function toJSON
     * @memberof CafeStoreThread
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeStoreThread.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeStoreThread;
})();

export const CafeStored = $root.CafeStored = (() => {

    /**
     * Properties of a CafeStored.
     * @exports ICafeStored
     * @interface ICafeStored
     * @property {string} id CafeStored id
     */

    /**
     * Constructs a new CafeStored.
     * @exports CafeStored
     * @classdesc Represents a CafeStored.
     * @implements ICafeStored
     * @constructor
     * @param {ICafeStored=} [properties] Properties to set
     */
    function CafeStored(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeStored id.
     * @member {string} id
     * @memberof CafeStored
     * @instance
     */
    CafeStored.prototype.id = "";

    /**
     * Creates a new CafeStored instance using the specified properties.
     * @function create
     * @memberof CafeStored
     * @static
     * @param {ICafeStored=} [properties] Properties to set
     * @returns {CafeStored} CafeStored instance
     */
    CafeStored.create = function create(properties) {
        return new CafeStored(properties);
    };

    /**
     * Encodes the specified CafeStored message. Does not implicitly {@link CafeStored.verify|verify} messages.
     * @function encode
     * @memberof CafeStored
     * @static
     * @param {ICafeStored} message CafeStored message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeStored.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        return writer;
    };

    /**
     * Encodes the specified CafeStored message, length delimited. Does not implicitly {@link CafeStored.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeStored
     * @static
     * @param {ICafeStored} message CafeStored message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeStored.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeStored message from the specified reader or buffer.
     * @function decode
     * @memberof CafeStored
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeStored} CafeStored
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeStored.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeStored();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeStored message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeStored
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeStored} CafeStored
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeStored.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeStored message.
     * @function verify
     * @memberof CafeStored
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeStored.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        return null;
    };

    /**
     * Creates a CafeStored message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeStored
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeStored} CafeStored
     */
    CafeStored.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeStored)
            return object;
        let message = new $root.CafeStored();
        if (object.id != null)
            message.id = String(object.id);
        return message;
    };

    /**
     * Creates a plain object from a CafeStored message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeStored
     * @static
     * @param {CafeStored} message CafeStored
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeStored.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults)
            object.id = "";
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        return object;
    };

    /**
     * Converts this CafeStored to JSON.
     * @function toJSON
     * @memberof CafeStored
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeStored.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeStored;
})();

export const CafeDeliverMessage = $root.CafeDeliverMessage = (() => {

    /**
     * Properties of a CafeDeliverMessage.
     * @exports ICafeDeliverMessage
     * @interface ICafeDeliverMessage
     * @property {string} id CafeDeliverMessage id
     * @property {string} client CafeDeliverMessage client
     */

    /**
     * Constructs a new CafeDeliverMessage.
     * @exports CafeDeliverMessage
     * @classdesc Represents a CafeDeliverMessage.
     * @implements ICafeDeliverMessage
     * @constructor
     * @param {ICafeDeliverMessage=} [properties] Properties to set
     */
    function CafeDeliverMessage(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeDeliverMessage id.
     * @member {string} id
     * @memberof CafeDeliverMessage
     * @instance
     */
    CafeDeliverMessage.prototype.id = "";

    /**
     * CafeDeliverMessage client.
     * @member {string} client
     * @memberof CafeDeliverMessage
     * @instance
     */
    CafeDeliverMessage.prototype.client = "";

    /**
     * Creates a new CafeDeliverMessage instance using the specified properties.
     * @function create
     * @memberof CafeDeliverMessage
     * @static
     * @param {ICafeDeliverMessage=} [properties] Properties to set
     * @returns {CafeDeliverMessage} CafeDeliverMessage instance
     */
    CafeDeliverMessage.create = function create(properties) {
        return new CafeDeliverMessage(properties);
    };

    /**
     * Encodes the specified CafeDeliverMessage message. Does not implicitly {@link CafeDeliverMessage.verify|verify} messages.
     * @function encode
     * @memberof CafeDeliverMessage
     * @static
     * @param {ICafeDeliverMessage} message CafeDeliverMessage message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeDeliverMessage.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.client != null && message.hasOwnProperty("client"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.client);
        return writer;
    };

    /**
     * Encodes the specified CafeDeliverMessage message, length delimited. Does not implicitly {@link CafeDeliverMessage.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeDeliverMessage
     * @static
     * @param {ICafeDeliverMessage} message CafeDeliverMessage message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeDeliverMessage.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeDeliverMessage message from the specified reader or buffer.
     * @function decode
     * @memberof CafeDeliverMessage
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeDeliverMessage} CafeDeliverMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeDeliverMessage.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeDeliverMessage();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.client = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeDeliverMessage message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeDeliverMessage
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeDeliverMessage} CafeDeliverMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeDeliverMessage.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeDeliverMessage message.
     * @function verify
     * @memberof CafeDeliverMessage
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeDeliverMessage.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.client != null && message.hasOwnProperty("client"))
            if (!$util.isString(message.client))
                return "client: string expected";
        return null;
    };

    /**
     * Creates a CafeDeliverMessage message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeDeliverMessage
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeDeliverMessage} CafeDeliverMessage
     */
    CafeDeliverMessage.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeDeliverMessage)
            return object;
        let message = new $root.CafeDeliverMessage();
        if (object.id != null)
            message.id = String(object.id);
        if (object.client != null)
            message.client = String(object.client);
        return message;
    };

    /**
     * Creates a plain object from a CafeDeliverMessage message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeDeliverMessage
     * @static
     * @param {CafeDeliverMessage} message CafeDeliverMessage
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeDeliverMessage.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            object.client = "";
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.client != null && message.hasOwnProperty("client"))
            object.client = message.client;
        return object;
    };

    /**
     * Converts this CafeDeliverMessage to JSON.
     * @function toJSON
     * @memberof CafeDeliverMessage
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeDeliverMessage.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeDeliverMessage;
})();

export const CafeCheckMessages = $root.CafeCheckMessages = (() => {

    /**
     * Properties of a CafeCheckMessages.
     * @exports ICafeCheckMessages
     * @interface ICafeCheckMessages
     * @property {string} token CafeCheckMessages token
     */

    /**
     * Constructs a new CafeCheckMessages.
     * @exports CafeCheckMessages
     * @classdesc Represents a CafeCheckMessages.
     * @implements ICafeCheckMessages
     * @constructor
     * @param {ICafeCheckMessages=} [properties] Properties to set
     */
    function CafeCheckMessages(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeCheckMessages token.
     * @member {string} token
     * @memberof CafeCheckMessages
     * @instance
     */
    CafeCheckMessages.prototype.token = "";

    /**
     * Creates a new CafeCheckMessages instance using the specified properties.
     * @function create
     * @memberof CafeCheckMessages
     * @static
     * @param {ICafeCheckMessages=} [properties] Properties to set
     * @returns {CafeCheckMessages} CafeCheckMessages instance
     */
    CafeCheckMessages.create = function create(properties) {
        return new CafeCheckMessages(properties);
    };

    /**
     * Encodes the specified CafeCheckMessages message. Does not implicitly {@link CafeCheckMessages.verify|verify} messages.
     * @function encode
     * @memberof CafeCheckMessages
     * @static
     * @param {ICafeCheckMessages} message CafeCheckMessages message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeCheckMessages.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.token != null && message.hasOwnProperty("token"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.token);
        return writer;
    };

    /**
     * Encodes the specified CafeCheckMessages message, length delimited. Does not implicitly {@link CafeCheckMessages.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeCheckMessages
     * @static
     * @param {ICafeCheckMessages} message CafeCheckMessages message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeCheckMessages.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeCheckMessages message from the specified reader or buffer.
     * @function decode
     * @memberof CafeCheckMessages
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeCheckMessages} CafeCheckMessages
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeCheckMessages.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeCheckMessages();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.token = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeCheckMessages message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeCheckMessages
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeCheckMessages} CafeCheckMessages
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeCheckMessages.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeCheckMessages message.
     * @function verify
     * @memberof CafeCheckMessages
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeCheckMessages.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.token != null && message.hasOwnProperty("token"))
            if (!$util.isString(message.token))
                return "token: string expected";
        return null;
    };

    /**
     * Creates a CafeCheckMessages message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeCheckMessages
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeCheckMessages} CafeCheckMessages
     */
    CafeCheckMessages.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeCheckMessages)
            return object;
        let message = new $root.CafeCheckMessages();
        if (object.token != null)
            message.token = String(object.token);
        return message;
    };

    /**
     * Creates a plain object from a CafeCheckMessages message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeCheckMessages
     * @static
     * @param {CafeCheckMessages} message CafeCheckMessages
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeCheckMessages.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults)
            object.token = "";
        if (message.token != null && message.hasOwnProperty("token"))
            object.token = message.token;
        return object;
    };

    /**
     * Converts this CafeCheckMessages to JSON.
     * @function toJSON
     * @memberof CafeCheckMessages
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeCheckMessages.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeCheckMessages;
})();

export const CafeMessages = $root.CafeMessages = (() => {

    /**
     * Properties of a CafeMessages.
     * @exports ICafeMessages
     * @interface ICafeMessages
     * @property {Array.<ICafeMessage>} messages CafeMessages messages
     */

    /**
     * Constructs a new CafeMessages.
     * @exports CafeMessages
     * @classdesc Represents a CafeMessages.
     * @implements ICafeMessages
     * @constructor
     * @param {ICafeMessages=} [properties] Properties to set
     */
    function CafeMessages(properties) {
        this.messages = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeMessages messages.
     * @member {Array.<ICafeMessage>} messages
     * @memberof CafeMessages
     * @instance
     */
    CafeMessages.prototype.messages = $util.emptyArray;

    /**
     * Creates a new CafeMessages instance using the specified properties.
     * @function create
     * @memberof CafeMessages
     * @static
     * @param {ICafeMessages=} [properties] Properties to set
     * @returns {CafeMessages} CafeMessages instance
     */
    CafeMessages.create = function create(properties) {
        return new CafeMessages(properties);
    };

    /**
     * Encodes the specified CafeMessages message. Does not implicitly {@link CafeMessages.verify|verify} messages.
     * @function encode
     * @memberof CafeMessages
     * @static
     * @param {ICafeMessages} message CafeMessages message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeMessages.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.messages != null && message.messages.length)
            for (let i = 0; i < message.messages.length; ++i)
                $root.CafeMessage.encode(message.messages[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified CafeMessages message, length delimited. Does not implicitly {@link CafeMessages.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeMessages
     * @static
     * @param {ICafeMessages} message CafeMessages message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeMessages.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeMessages message from the specified reader or buffer.
     * @function decode
     * @memberof CafeMessages
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeMessages} CafeMessages
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeMessages.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeMessages();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                if (!(message.messages && message.messages.length))
                    message.messages = [];
                message.messages.push($root.CafeMessage.decode(reader, reader.uint32()));
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeMessages message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeMessages
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeMessages} CafeMessages
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeMessages.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeMessages message.
     * @function verify
     * @memberof CafeMessages
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeMessages.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.messages != null && message.hasOwnProperty("messages")) {
            if (!Array.isArray(message.messages))
                return "messages: array expected";
            for (let i = 0; i < message.messages.length; ++i) {
                let error = $root.CafeMessage.verify(message.messages[i]);
                if (error)
                    return "messages." + error;
            }
        }
        return null;
    };

    /**
     * Creates a CafeMessages message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeMessages
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeMessages} CafeMessages
     */
    CafeMessages.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeMessages)
            return object;
        let message = new $root.CafeMessages();
        if (object.messages) {
            if (!Array.isArray(object.messages))
                throw TypeError(".CafeMessages.messages: array expected");
            message.messages = [];
            for (let i = 0; i < object.messages.length; ++i) {
                if (typeof object.messages[i] !== "object")
                    throw TypeError(".CafeMessages.messages: object expected");
                message.messages[i] = $root.CafeMessage.fromObject(object.messages[i]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from a CafeMessages message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeMessages
     * @static
     * @param {CafeMessages} message CafeMessages
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeMessages.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.messages = [];
        if (message.messages && message.messages.length) {
            object.messages = [];
            for (let j = 0; j < message.messages.length; ++j)
                object.messages[j] = $root.CafeMessage.toObject(message.messages[j], options);
        }
        return object;
    };

    /**
     * Converts this CafeMessages to JSON.
     * @function toJSON
     * @memberof CafeMessages
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeMessages.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeMessages;
})();

export const CafeDeleteMessages = $root.CafeDeleteMessages = (() => {

    /**
     * Properties of a CafeDeleteMessages.
     * @exports ICafeDeleteMessages
     * @interface ICafeDeleteMessages
     * @property {string} token CafeDeleteMessages token
     */

    /**
     * Constructs a new CafeDeleteMessages.
     * @exports CafeDeleteMessages
     * @classdesc Represents a CafeDeleteMessages.
     * @implements ICafeDeleteMessages
     * @constructor
     * @param {ICafeDeleteMessages=} [properties] Properties to set
     */
    function CafeDeleteMessages(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeDeleteMessages token.
     * @member {string} token
     * @memberof CafeDeleteMessages
     * @instance
     */
    CafeDeleteMessages.prototype.token = "";

    /**
     * Creates a new CafeDeleteMessages instance using the specified properties.
     * @function create
     * @memberof CafeDeleteMessages
     * @static
     * @param {ICafeDeleteMessages=} [properties] Properties to set
     * @returns {CafeDeleteMessages} CafeDeleteMessages instance
     */
    CafeDeleteMessages.create = function create(properties) {
        return new CafeDeleteMessages(properties);
    };

    /**
     * Encodes the specified CafeDeleteMessages message. Does not implicitly {@link CafeDeleteMessages.verify|verify} messages.
     * @function encode
     * @memberof CafeDeleteMessages
     * @static
     * @param {ICafeDeleteMessages} message CafeDeleteMessages message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeDeleteMessages.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.token != null && message.hasOwnProperty("token"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.token);
        return writer;
    };

    /**
     * Encodes the specified CafeDeleteMessages message, length delimited. Does not implicitly {@link CafeDeleteMessages.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeDeleteMessages
     * @static
     * @param {ICafeDeleteMessages} message CafeDeleteMessages message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeDeleteMessages.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeDeleteMessages message from the specified reader or buffer.
     * @function decode
     * @memberof CafeDeleteMessages
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeDeleteMessages} CafeDeleteMessages
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeDeleteMessages.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeDeleteMessages();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.token = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeDeleteMessages message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeDeleteMessages
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeDeleteMessages} CafeDeleteMessages
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeDeleteMessages.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeDeleteMessages message.
     * @function verify
     * @memberof CafeDeleteMessages
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeDeleteMessages.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.token != null && message.hasOwnProperty("token"))
            if (!$util.isString(message.token))
                return "token: string expected";
        return null;
    };

    /**
     * Creates a CafeDeleteMessages message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeDeleteMessages
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeDeleteMessages} CafeDeleteMessages
     */
    CafeDeleteMessages.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeDeleteMessages)
            return object;
        let message = new $root.CafeDeleteMessages();
        if (object.token != null)
            message.token = String(object.token);
        return message;
    };

    /**
     * Creates a plain object from a CafeDeleteMessages message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeDeleteMessages
     * @static
     * @param {CafeDeleteMessages} message CafeDeleteMessages
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeDeleteMessages.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults)
            object.token = "";
        if (message.token != null && message.hasOwnProperty("token"))
            object.token = message.token;
        return object;
    };

    /**
     * Converts this CafeDeleteMessages to JSON.
     * @function toJSON
     * @memberof CafeDeleteMessages
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeDeleteMessages.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeDeleteMessages;
})();

export const CafeDeleteMessagesAck = $root.CafeDeleteMessagesAck = (() => {

    /**
     * Properties of a CafeDeleteMessagesAck.
     * @exports ICafeDeleteMessagesAck
     * @interface ICafeDeleteMessagesAck
     * @property {boolean} more CafeDeleteMessagesAck more
     */

    /**
     * Constructs a new CafeDeleteMessagesAck.
     * @exports CafeDeleteMessagesAck
     * @classdesc Represents a CafeDeleteMessagesAck.
     * @implements ICafeDeleteMessagesAck
     * @constructor
     * @param {ICafeDeleteMessagesAck=} [properties] Properties to set
     */
    function CafeDeleteMessagesAck(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeDeleteMessagesAck more.
     * @member {boolean} more
     * @memberof CafeDeleteMessagesAck
     * @instance
     */
    CafeDeleteMessagesAck.prototype.more = false;

    /**
     * Creates a new CafeDeleteMessagesAck instance using the specified properties.
     * @function create
     * @memberof CafeDeleteMessagesAck
     * @static
     * @param {ICafeDeleteMessagesAck=} [properties] Properties to set
     * @returns {CafeDeleteMessagesAck} CafeDeleteMessagesAck instance
     */
    CafeDeleteMessagesAck.create = function create(properties) {
        return new CafeDeleteMessagesAck(properties);
    };

    /**
     * Encodes the specified CafeDeleteMessagesAck message. Does not implicitly {@link CafeDeleteMessagesAck.verify|verify} messages.
     * @function encode
     * @memberof CafeDeleteMessagesAck
     * @static
     * @param {ICafeDeleteMessagesAck} message CafeDeleteMessagesAck message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeDeleteMessagesAck.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.more != null && message.hasOwnProperty("more"))
            writer.uint32(/* id 1, wireType 0 =*/8).bool(message.more);
        return writer;
    };

    /**
     * Encodes the specified CafeDeleteMessagesAck message, length delimited. Does not implicitly {@link CafeDeleteMessagesAck.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeDeleteMessagesAck
     * @static
     * @param {ICafeDeleteMessagesAck} message CafeDeleteMessagesAck message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeDeleteMessagesAck.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeDeleteMessagesAck message from the specified reader or buffer.
     * @function decode
     * @memberof CafeDeleteMessagesAck
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeDeleteMessagesAck} CafeDeleteMessagesAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeDeleteMessagesAck.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeDeleteMessagesAck();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.more = reader.bool();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeDeleteMessagesAck message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeDeleteMessagesAck
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeDeleteMessagesAck} CafeDeleteMessagesAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeDeleteMessagesAck.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeDeleteMessagesAck message.
     * @function verify
     * @memberof CafeDeleteMessagesAck
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeDeleteMessagesAck.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.more != null && message.hasOwnProperty("more"))
            if (typeof message.more !== "boolean")
                return "more: boolean expected";
        return null;
    };

    /**
     * Creates a CafeDeleteMessagesAck message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeDeleteMessagesAck
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeDeleteMessagesAck} CafeDeleteMessagesAck
     */
    CafeDeleteMessagesAck.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeDeleteMessagesAck)
            return object;
        let message = new $root.CafeDeleteMessagesAck();
        if (object.more != null)
            message.more = Boolean(object.more);
        return message;
    };

    /**
     * Creates a plain object from a CafeDeleteMessagesAck message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeDeleteMessagesAck
     * @static
     * @param {CafeDeleteMessagesAck} message CafeDeleteMessagesAck
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeDeleteMessagesAck.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults)
            object.more = false;
        if (message.more != null && message.hasOwnProperty("more"))
            object.more = message.more;
        return object;
    };

    /**
     * Converts this CafeDeleteMessagesAck to JSON.
     * @function toJSON
     * @memberof CafeDeleteMessagesAck
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeDeleteMessagesAck.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeDeleteMessagesAck;
})();

export const Contact = $root.Contact = (() => {

    /**
     * Properties of a Contact.
     * @exports IContact
     * @interface IContact
     * @property {string} id Contact id
     * @property {string} address Contact address
     * @property {string} username Contact username
     * @property {string} avatar Contact avatar
     * @property {Array.<ICafe>} inboxes Contact inboxes
     * @property {google.protobuf.ITimestamp} created Contact created
     * @property {google.protobuf.ITimestamp} updated Contact updated
     * @property {Array.<string>} threads Contact threads
     */

    /**
     * Constructs a new Contact.
     * @exports Contact
     * @classdesc Represents a Contact.
     * @implements IContact
     * @constructor
     * @param {IContact=} [properties] Properties to set
     */
    function Contact(properties) {
        this.inboxes = [];
        this.threads = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Contact id.
     * @member {string} id
     * @memberof Contact
     * @instance
     */
    Contact.prototype.id = "";

    /**
     * Contact address.
     * @member {string} address
     * @memberof Contact
     * @instance
     */
    Contact.prototype.address = "";

    /**
     * Contact username.
     * @member {string} username
     * @memberof Contact
     * @instance
     */
    Contact.prototype.username = "";

    /**
     * Contact avatar.
     * @member {string} avatar
     * @memberof Contact
     * @instance
     */
    Contact.prototype.avatar = "";

    /**
     * Contact inboxes.
     * @member {Array.<ICafe>} inboxes
     * @memberof Contact
     * @instance
     */
    Contact.prototype.inboxes = $util.emptyArray;

    /**
     * Contact created.
     * @member {google.protobuf.ITimestamp} created
     * @memberof Contact
     * @instance
     */
    Contact.prototype.created = null;

    /**
     * Contact updated.
     * @member {google.protobuf.ITimestamp} updated
     * @memberof Contact
     * @instance
     */
    Contact.prototype.updated = null;

    /**
     * Contact threads.
     * @member {Array.<string>} threads
     * @memberof Contact
     * @instance
     */
    Contact.prototype.threads = $util.emptyArray;

    /**
     * Creates a new Contact instance using the specified properties.
     * @function create
     * @memberof Contact
     * @static
     * @param {IContact=} [properties] Properties to set
     * @returns {Contact} Contact instance
     */
    Contact.create = function create(properties) {
        return new Contact(properties);
    };

    /**
     * Encodes the specified Contact message. Does not implicitly {@link Contact.verify|verify} messages.
     * @function encode
     * @memberof Contact
     * @static
     * @param {IContact} message Contact message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Contact.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.address != null && message.hasOwnProperty("address"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.address);
        if (message.username != null && message.hasOwnProperty("username"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.username);
        if (message.avatar != null && message.hasOwnProperty("avatar"))
            writer.uint32(/* id 4, wireType 2 =*/34).string(message.avatar);
        if (message.inboxes != null && message.inboxes.length)
            for (let i = 0; i < message.inboxes.length; ++i)
                $root.Cafe.encode(message.inboxes[i], writer.uint32(/* id 5, wireType 2 =*/42).fork()).ldelim();
        if (message.created != null && message.hasOwnProperty("created"))
            $root.google.protobuf.Timestamp.encode(message.created, writer.uint32(/* id 6, wireType 2 =*/50).fork()).ldelim();
        if (message.updated != null && message.hasOwnProperty("updated"))
            $root.google.protobuf.Timestamp.encode(message.updated, writer.uint32(/* id 7, wireType 2 =*/58).fork()).ldelim();
        if (message.threads != null && message.threads.length)
            for (let i = 0; i < message.threads.length; ++i)
                writer.uint32(/* id 8, wireType 2 =*/66).string(message.threads[i]);
        return writer;
    };

    /**
     * Encodes the specified Contact message, length delimited. Does not implicitly {@link Contact.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Contact
     * @static
     * @param {IContact} message Contact message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Contact.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Contact message from the specified reader or buffer.
     * @function decode
     * @memberof Contact
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Contact} Contact
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Contact.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Contact();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.address = reader.string();
                break;
            case 3:
                message.username = reader.string();
                break;
            case 4:
                message.avatar = reader.string();
                break;
            case 5:
                if (!(message.inboxes && message.inboxes.length))
                    message.inboxes = [];
                message.inboxes.push($root.Cafe.decode(reader, reader.uint32()));
                break;
            case 6:
                message.created = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 7:
                message.updated = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 8:
                if (!(message.threads && message.threads.length))
                    message.threads = [];
                message.threads.push(reader.string());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Contact message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Contact
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Contact} Contact
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Contact.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Contact message.
     * @function verify
     * @memberof Contact
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Contact.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.address != null && message.hasOwnProperty("address"))
            if (!$util.isString(message.address))
                return "address: string expected";
        if (message.username != null && message.hasOwnProperty("username"))
            if (!$util.isString(message.username))
                return "username: string expected";
        if (message.avatar != null && message.hasOwnProperty("avatar"))
            if (!$util.isString(message.avatar))
                return "avatar: string expected";
        if (message.inboxes != null && message.hasOwnProperty("inboxes")) {
            if (!Array.isArray(message.inboxes))
                return "inboxes: array expected";
            for (let i = 0; i < message.inboxes.length; ++i) {
                let error = $root.Cafe.verify(message.inboxes[i]);
                if (error)
                    return "inboxes." + error;
            }
        }
        if (message.created != null && message.hasOwnProperty("created")) {
            let error = $root.google.protobuf.Timestamp.verify(message.created);
            if (error)
                return "created." + error;
        }
        if (message.updated != null && message.hasOwnProperty("updated")) {
            let error = $root.google.protobuf.Timestamp.verify(message.updated);
            if (error)
                return "updated." + error;
        }
        if (message.threads != null && message.hasOwnProperty("threads")) {
            if (!Array.isArray(message.threads))
                return "threads: array expected";
            for (let i = 0; i < message.threads.length; ++i)
                if (!$util.isString(message.threads[i]))
                    return "threads: string[] expected";
        }
        return null;
    };

    /**
     * Creates a Contact message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Contact
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Contact} Contact
     */
    Contact.fromObject = function fromObject(object) {
        if (object instanceof $root.Contact)
            return object;
        let message = new $root.Contact();
        if (object.id != null)
            message.id = String(object.id);
        if (object.address != null)
            message.address = String(object.address);
        if (object.username != null)
            message.username = String(object.username);
        if (object.avatar != null)
            message.avatar = String(object.avatar);
        if (object.inboxes) {
            if (!Array.isArray(object.inboxes))
                throw TypeError(".Contact.inboxes: array expected");
            message.inboxes = [];
            for (let i = 0; i < object.inboxes.length; ++i) {
                if (typeof object.inboxes[i] !== "object")
                    throw TypeError(".Contact.inboxes: object expected");
                message.inboxes[i] = $root.Cafe.fromObject(object.inboxes[i]);
            }
        }
        if (object.created != null) {
            if (typeof object.created !== "object")
                throw TypeError(".Contact.created: object expected");
            message.created = $root.google.protobuf.Timestamp.fromObject(object.created);
        }
        if (object.updated != null) {
            if (typeof object.updated !== "object")
                throw TypeError(".Contact.updated: object expected");
            message.updated = $root.google.protobuf.Timestamp.fromObject(object.updated);
        }
        if (object.threads) {
            if (!Array.isArray(object.threads))
                throw TypeError(".Contact.threads: array expected");
            message.threads = [];
            for (let i = 0; i < object.threads.length; ++i)
                message.threads[i] = String(object.threads[i]);
        }
        return message;
    };

    /**
     * Creates a plain object from a Contact message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Contact
     * @static
     * @param {Contact} message Contact
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Contact.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults) {
            object.inboxes = [];
            object.threads = [];
        }
        if (options.defaults) {
            object.id = "";
            object.address = "";
            object.username = "";
            object.avatar = "";
            object.created = null;
            object.updated = null;
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.address != null && message.hasOwnProperty("address"))
            object.address = message.address;
        if (message.username != null && message.hasOwnProperty("username"))
            object.username = message.username;
        if (message.avatar != null && message.hasOwnProperty("avatar"))
            object.avatar = message.avatar;
        if (message.inboxes && message.inboxes.length) {
            object.inboxes = [];
            for (let j = 0; j < message.inboxes.length; ++j)
                object.inboxes[j] = $root.Cafe.toObject(message.inboxes[j], options);
        }
        if (message.created != null && message.hasOwnProperty("created"))
            object.created = $root.google.protobuf.Timestamp.toObject(message.created, options);
        if (message.updated != null && message.hasOwnProperty("updated"))
            object.updated = $root.google.protobuf.Timestamp.toObject(message.updated, options);
        if (message.threads && message.threads.length) {
            object.threads = [];
            for (let j = 0; j < message.threads.length; ++j)
                object.threads[j] = message.threads[j];
        }
        return object;
    };

    /**
     * Converts this Contact to JSON.
     * @function toJSON
     * @memberof Contact
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Contact.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Contact;
})();

export const ContactList = $root.ContactList = (() => {

    /**
     * Properties of a ContactList.
     * @exports IContactList
     * @interface IContactList
     * @property {Array.<IContact>} items ContactList items
     */

    /**
     * Constructs a new ContactList.
     * @exports ContactList
     * @classdesc Represents a ContactList.
     * @implements IContactList
     * @constructor
     * @param {IContactList=} [properties] Properties to set
     */
    function ContactList(properties) {
        this.items = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ContactList items.
     * @member {Array.<IContact>} items
     * @memberof ContactList
     * @instance
     */
    ContactList.prototype.items = $util.emptyArray;

    /**
     * Creates a new ContactList instance using the specified properties.
     * @function create
     * @memberof ContactList
     * @static
     * @param {IContactList=} [properties] Properties to set
     * @returns {ContactList} ContactList instance
     */
    ContactList.create = function create(properties) {
        return new ContactList(properties);
    };

    /**
     * Encodes the specified ContactList message. Does not implicitly {@link ContactList.verify|verify} messages.
     * @function encode
     * @memberof ContactList
     * @static
     * @param {IContactList} message ContactList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ContactList.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.items != null && message.items.length)
            for (let i = 0; i < message.items.length; ++i)
                $root.Contact.encode(message.items[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified ContactList message, length delimited. Does not implicitly {@link ContactList.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ContactList
     * @static
     * @param {IContactList} message ContactList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ContactList.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ContactList message from the specified reader or buffer.
     * @function decode
     * @memberof ContactList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ContactList} ContactList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ContactList.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ContactList();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                if (!(message.items && message.items.length))
                    message.items = [];
                message.items.push($root.Contact.decode(reader, reader.uint32()));
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ContactList message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ContactList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ContactList} ContactList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ContactList.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ContactList message.
     * @function verify
     * @memberof ContactList
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ContactList.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.items != null && message.hasOwnProperty("items")) {
            if (!Array.isArray(message.items))
                return "items: array expected";
            for (let i = 0; i < message.items.length; ++i) {
                let error = $root.Contact.verify(message.items[i]);
                if (error)
                    return "items." + error;
            }
        }
        return null;
    };

    /**
     * Creates a ContactList message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ContactList
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ContactList} ContactList
     */
    ContactList.fromObject = function fromObject(object) {
        if (object instanceof $root.ContactList)
            return object;
        let message = new $root.ContactList();
        if (object.items) {
            if (!Array.isArray(object.items))
                throw TypeError(".ContactList.items: array expected");
            message.items = [];
            for (let i = 0; i < object.items.length; ++i) {
                if (typeof object.items[i] !== "object")
                    throw TypeError(".ContactList.items: object expected");
                message.items[i] = $root.Contact.fromObject(object.items[i]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from a ContactList message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ContactList
     * @static
     * @param {ContactList} message ContactList
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ContactList.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.items = [];
        if (message.items && message.items.length) {
            object.items = [];
            for (let j = 0; j < message.items.length; ++j)
                object.items[j] = $root.Contact.toObject(message.items[j], options);
        }
        return object;
    };

    /**
     * Converts this ContactList to JSON.
     * @function toJSON
     * @memberof ContactList
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ContactList.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ContactList;
})();

export const User = $root.User = (() => {

    /**
     * Properties of a User.
     * @exports IUser
     * @interface IUser
     * @property {string} address User address
     * @property {string} name User name
     * @property {string} avatar User avatar
     */

    /**
     * Constructs a new User.
     * @exports User
     * @classdesc Represents a User.
     * @implements IUser
     * @constructor
     * @param {IUser=} [properties] Properties to set
     */
    function User(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * User address.
     * @member {string} address
     * @memberof User
     * @instance
     */
    User.prototype.address = "";

    /**
     * User name.
     * @member {string} name
     * @memberof User
     * @instance
     */
    User.prototype.name = "";

    /**
     * User avatar.
     * @member {string} avatar
     * @memberof User
     * @instance
     */
    User.prototype.avatar = "";

    /**
     * Creates a new User instance using the specified properties.
     * @function create
     * @memberof User
     * @static
     * @param {IUser=} [properties] Properties to set
     * @returns {User} User instance
     */
    User.create = function create(properties) {
        return new User(properties);
    };

    /**
     * Encodes the specified User message. Does not implicitly {@link User.verify|verify} messages.
     * @function encode
     * @memberof User
     * @static
     * @param {IUser} message User message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    User.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.address != null && message.hasOwnProperty("address"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.address);
        if (message.name != null && message.hasOwnProperty("name"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.name);
        if (message.avatar != null && message.hasOwnProperty("avatar"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.avatar);
        return writer;
    };

    /**
     * Encodes the specified User message, length delimited. Does not implicitly {@link User.verify|verify} messages.
     * @function encodeDelimited
     * @memberof User
     * @static
     * @param {IUser} message User message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    User.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a User message from the specified reader or buffer.
     * @function decode
     * @memberof User
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {User} User
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    User.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.User();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.address = reader.string();
                break;
            case 2:
                message.name = reader.string();
                break;
            case 3:
                message.avatar = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a User message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof User
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {User} User
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    User.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a User message.
     * @function verify
     * @memberof User
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    User.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.address != null && message.hasOwnProperty("address"))
            if (!$util.isString(message.address))
                return "address: string expected";
        if (message.name != null && message.hasOwnProperty("name"))
            if (!$util.isString(message.name))
                return "name: string expected";
        if (message.avatar != null && message.hasOwnProperty("avatar"))
            if (!$util.isString(message.avatar))
                return "avatar: string expected";
        return null;
    };

    /**
     * Creates a User message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof User
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {User} User
     */
    User.fromObject = function fromObject(object) {
        if (object instanceof $root.User)
            return object;
        let message = new $root.User();
        if (object.address != null)
            message.address = String(object.address);
        if (object.name != null)
            message.name = String(object.name);
        if (object.avatar != null)
            message.avatar = String(object.avatar);
        return message;
    };

    /**
     * Creates a plain object from a User message. Also converts values to other types if specified.
     * @function toObject
     * @memberof User
     * @static
     * @param {User} message User
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    User.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.address = "";
            object.name = "";
            object.avatar = "";
        }
        if (message.address != null && message.hasOwnProperty("address"))
            object.address = message.address;
        if (message.name != null && message.hasOwnProperty("name"))
            object.name = message.name;
        if (message.avatar != null && message.hasOwnProperty("avatar"))
            object.avatar = message.avatar;
        return object;
    };

    /**
     * Converts this User to JSON.
     * @function toJSON
     * @memberof User
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    User.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return User;
})();

export const Thread = $root.Thread = (() => {

    /**
     * Properties of a Thread.
     * @exports IThread
     * @interface IThread
     * @property {string} id Thread id
     * @property {string} key Thread key
     * @property {Uint8Array} sk Thread sk
     * @property {string} name Thread name
     * @property {string} schema Thread schema
     * @property {string} initiator Thread initiator
     * @property {Thread.Type} type Thread type
     * @property {Thread.Sharing} sharing Thread sharing
     * @property {Array.<string>} members Thread members
     * @property {Thread.State} state Thread state
     * @property {string} head Thread head
     * @property {IBlock} headBlock Thread headBlock
     * @property {INode} schemaNode Thread schemaNode
     * @property {number} blockCount Thread blockCount
     * @property {number} peerCount Thread peerCount
     */

    /**
     * Constructs a new Thread.
     * @exports Thread
     * @classdesc Represents a Thread.
     * @implements IThread
     * @constructor
     * @param {IThread=} [properties] Properties to set
     */
    function Thread(properties) {
        this.members = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Thread id.
     * @member {string} id
     * @memberof Thread
     * @instance
     */
    Thread.prototype.id = "";

    /**
     * Thread key.
     * @member {string} key
     * @memberof Thread
     * @instance
     */
    Thread.prototype.key = "";

    /**
     * Thread sk.
     * @member {Uint8Array} sk
     * @memberof Thread
     * @instance
     */
    Thread.prototype.sk = $util.newBuffer([]);

    /**
     * Thread name.
     * @member {string} name
     * @memberof Thread
     * @instance
     */
    Thread.prototype.name = "";

    /**
     * Thread schema.
     * @member {string} schema
     * @memberof Thread
     * @instance
     */
    Thread.prototype.schema = "";

    /**
     * Thread initiator.
     * @member {string} initiator
     * @memberof Thread
     * @instance
     */
    Thread.prototype.initiator = "";

    /**
     * Thread type.
     * @member {Thread.Type} type
     * @memberof Thread
     * @instance
     */
    Thread.prototype.type = 0;

    /**
     * Thread sharing.
     * @member {Thread.Sharing} sharing
     * @memberof Thread
     * @instance
     */
    Thread.prototype.sharing = 0;

    /**
     * Thread members.
     * @member {Array.<string>} members
     * @memberof Thread
     * @instance
     */
    Thread.prototype.members = $util.emptyArray;

    /**
     * Thread state.
     * @member {Thread.State} state
     * @memberof Thread
     * @instance
     */
    Thread.prototype.state = 0;

    /**
     * Thread head.
     * @member {string} head
     * @memberof Thread
     * @instance
     */
    Thread.prototype.head = "";

    /**
     * Thread headBlock.
     * @member {IBlock} headBlock
     * @memberof Thread
     * @instance
     */
    Thread.prototype.headBlock = null;

    /**
     * Thread schemaNode.
     * @member {INode} schemaNode
     * @memberof Thread
     * @instance
     */
    Thread.prototype.schemaNode = null;

    /**
     * Thread blockCount.
     * @member {number} blockCount
     * @memberof Thread
     * @instance
     */
    Thread.prototype.blockCount = 0;

    /**
     * Thread peerCount.
     * @member {number} peerCount
     * @memberof Thread
     * @instance
     */
    Thread.prototype.peerCount = 0;

    /**
     * Creates a new Thread instance using the specified properties.
     * @function create
     * @memberof Thread
     * @static
     * @param {IThread=} [properties] Properties to set
     * @returns {Thread} Thread instance
     */
    Thread.create = function create(properties) {
        return new Thread(properties);
    };

    /**
     * Encodes the specified Thread message. Does not implicitly {@link Thread.verify|verify} messages.
     * @function encode
     * @memberof Thread
     * @static
     * @param {IThread} message Thread message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Thread.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.key != null && message.hasOwnProperty("key"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.key);
        if (message.sk != null && message.hasOwnProperty("sk"))
            writer.uint32(/* id 3, wireType 2 =*/26).bytes(message.sk);
        if (message.name != null && message.hasOwnProperty("name"))
            writer.uint32(/* id 4, wireType 2 =*/34).string(message.name);
        if (message.schema != null && message.hasOwnProperty("schema"))
            writer.uint32(/* id 5, wireType 2 =*/42).string(message.schema);
        if (message.initiator != null && message.hasOwnProperty("initiator"))
            writer.uint32(/* id 6, wireType 2 =*/50).string(message.initiator);
        if (message.type != null && message.hasOwnProperty("type"))
            writer.uint32(/* id 7, wireType 0 =*/56).int32(message.type);
        if (message.sharing != null && message.hasOwnProperty("sharing"))
            writer.uint32(/* id 8, wireType 0 =*/64).int32(message.sharing);
        if (message.members != null && message.members.length)
            for (let i = 0; i < message.members.length; ++i)
                writer.uint32(/* id 9, wireType 2 =*/74).string(message.members[i]);
        if (message.state != null && message.hasOwnProperty("state"))
            writer.uint32(/* id 10, wireType 0 =*/80).int32(message.state);
        if (message.head != null && message.hasOwnProperty("head"))
            writer.uint32(/* id 11, wireType 2 =*/90).string(message.head);
        if (message.headBlock != null && message.hasOwnProperty("headBlock"))
            $root.Block.encode(message.headBlock, writer.uint32(/* id 101, wireType 2 =*/810).fork()).ldelim();
        if (message.schemaNode != null && message.hasOwnProperty("schemaNode"))
            $root.Node.encode(message.schemaNode, writer.uint32(/* id 102, wireType 2 =*/818).fork()).ldelim();
        if (message.blockCount != null && message.hasOwnProperty("blockCount"))
            writer.uint32(/* id 103, wireType 0 =*/824).int32(message.blockCount);
        if (message.peerCount != null && message.hasOwnProperty("peerCount"))
            writer.uint32(/* id 104, wireType 0 =*/832).int32(message.peerCount);
        return writer;
    };

    /**
     * Encodes the specified Thread message, length delimited. Does not implicitly {@link Thread.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Thread
     * @static
     * @param {IThread} message Thread message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Thread.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Thread message from the specified reader or buffer.
     * @function decode
     * @memberof Thread
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Thread} Thread
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Thread.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Thread();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.key = reader.string();
                break;
            case 3:
                message.sk = reader.bytes();
                break;
            case 4:
                message.name = reader.string();
                break;
            case 5:
                message.schema = reader.string();
                break;
            case 6:
                message.initiator = reader.string();
                break;
            case 7:
                message.type = reader.int32();
                break;
            case 8:
                message.sharing = reader.int32();
                break;
            case 9:
                if (!(message.members && message.members.length))
                    message.members = [];
                message.members.push(reader.string());
                break;
            case 10:
                message.state = reader.int32();
                break;
            case 11:
                message.head = reader.string();
                break;
            case 101:
                message.headBlock = $root.Block.decode(reader, reader.uint32());
                break;
            case 102:
                message.schemaNode = $root.Node.decode(reader, reader.uint32());
                break;
            case 103:
                message.blockCount = reader.int32();
                break;
            case 104:
                message.peerCount = reader.int32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Thread message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Thread
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Thread} Thread
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Thread.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Thread message.
     * @function verify
     * @memberof Thread
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Thread.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.key != null && message.hasOwnProperty("key"))
            if (!$util.isString(message.key))
                return "key: string expected";
        if (message.sk != null && message.hasOwnProperty("sk"))
            if (!(message.sk && typeof message.sk.length === "number" || $util.isString(message.sk)))
                return "sk: buffer expected";
        if (message.name != null && message.hasOwnProperty("name"))
            if (!$util.isString(message.name))
                return "name: string expected";
        if (message.schema != null && message.hasOwnProperty("schema"))
            if (!$util.isString(message.schema))
                return "schema: string expected";
        if (message.initiator != null && message.hasOwnProperty("initiator"))
            if (!$util.isString(message.initiator))
                return "initiator: string expected";
        if (message.type != null && message.hasOwnProperty("type"))
            switch (message.type) {
            default:
                return "type: enum value expected";
            case 0:
            case 1:
            case 2:
            case 3:
                break;
            }
        if (message.sharing != null && message.hasOwnProperty("sharing"))
            switch (message.sharing) {
            default:
                return "sharing: enum value expected";
            case 0:
            case 1:
            case 2:
                break;
            }
        if (message.members != null && message.hasOwnProperty("members")) {
            if (!Array.isArray(message.members))
                return "members: array expected";
            for (let i = 0; i < message.members.length; ++i)
                if (!$util.isString(message.members[i]))
                    return "members: string[] expected";
        }
        if (message.state != null && message.hasOwnProperty("state"))
            switch (message.state) {
            default:
                return "state: enum value expected";
            case 0:
            case 1:
            case 2:
                break;
            }
        if (message.head != null && message.hasOwnProperty("head"))
            if (!$util.isString(message.head))
                return "head: string expected";
        if (message.headBlock != null && message.hasOwnProperty("headBlock")) {
            let error = $root.Block.verify(message.headBlock);
            if (error)
                return "headBlock." + error;
        }
        if (message.schemaNode != null && message.hasOwnProperty("schemaNode")) {
            let error = $root.Node.verify(message.schemaNode);
            if (error)
                return "schemaNode." + error;
        }
        if (message.blockCount != null && message.hasOwnProperty("blockCount"))
            if (!$util.isInteger(message.blockCount))
                return "blockCount: integer expected";
        if (message.peerCount != null && message.hasOwnProperty("peerCount"))
            if (!$util.isInteger(message.peerCount))
                return "peerCount: integer expected";
        return null;
    };

    /**
     * Creates a Thread message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Thread
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Thread} Thread
     */
    Thread.fromObject = function fromObject(object) {
        if (object instanceof $root.Thread)
            return object;
        let message = new $root.Thread();
        if (object.id != null)
            message.id = String(object.id);
        if (object.key != null)
            message.key = String(object.key);
        if (object.sk != null)
            if (typeof object.sk === "string")
                $util.base64.decode(object.sk, message.sk = $util.newBuffer($util.base64.length(object.sk)), 0);
            else if (object.sk.length)
                message.sk = object.sk;
        if (object.name != null)
            message.name = String(object.name);
        if (object.schema != null)
            message.schema = String(object.schema);
        if (object.initiator != null)
            message.initiator = String(object.initiator);
        switch (object.type) {
        case "Private":
        case 0:
            message.type = 0;
            break;
        case "ReadOnly":
        case 1:
            message.type = 1;
            break;
        case "Public":
        case 2:
            message.type = 2;
            break;
        case "Open":
        case 3:
            message.type = 3;
            break;
        }
        switch (object.sharing) {
        case "NotShared":
        case 0:
            message.sharing = 0;
            break;
        case "InviteOnly":
        case 1:
            message.sharing = 1;
            break;
        case "Shared":
        case 2:
            message.sharing = 2;
            break;
        }
        if (object.members) {
            if (!Array.isArray(object.members))
                throw TypeError(".Thread.members: array expected");
            message.members = [];
            for (let i = 0; i < object.members.length; ++i)
                message.members[i] = String(object.members[i]);
        }
        switch (object.state) {
        case "LoadingBehind":
        case 0:
            message.state = 0;
            break;
        case "Loaded":
        case 1:
            message.state = 1;
            break;
        case "LoadingAhead":
        case 2:
            message.state = 2;
            break;
        }
        if (object.head != null)
            message.head = String(object.head);
        if (object.headBlock != null) {
            if (typeof object.headBlock !== "object")
                throw TypeError(".Thread.headBlock: object expected");
            message.headBlock = $root.Block.fromObject(object.headBlock);
        }
        if (object.schemaNode != null) {
            if (typeof object.schemaNode !== "object")
                throw TypeError(".Thread.schemaNode: object expected");
            message.schemaNode = $root.Node.fromObject(object.schemaNode);
        }
        if (object.blockCount != null)
            message.blockCount = object.blockCount | 0;
        if (object.peerCount != null)
            message.peerCount = object.peerCount | 0;
        return message;
    };

    /**
     * Creates a plain object from a Thread message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Thread
     * @static
     * @param {Thread} message Thread
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Thread.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.members = [];
        if (options.defaults) {
            object.id = "";
            object.key = "";
            if (options.bytes === String)
                object.sk = "";
            else {
                object.sk = [];
                if (options.bytes !== Array)
                    object.sk = $util.newBuffer(object.sk);
            }
            object.name = "";
            object.schema = "";
            object.initiator = "";
            object.type = options.enums === String ? "Private" : 0;
            object.sharing = options.enums === String ? "NotShared" : 0;
            object.state = options.enums === String ? "LoadingBehind" : 0;
            object.head = "";
            object.headBlock = null;
            object.schemaNode = null;
            object.blockCount = 0;
            object.peerCount = 0;
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.key != null && message.hasOwnProperty("key"))
            object.key = message.key;
        if (message.sk != null && message.hasOwnProperty("sk"))
            object.sk = options.bytes === String ? $util.base64.encode(message.sk, 0, message.sk.length) : options.bytes === Array ? Array.prototype.slice.call(message.sk) : message.sk;
        if (message.name != null && message.hasOwnProperty("name"))
            object.name = message.name;
        if (message.schema != null && message.hasOwnProperty("schema"))
            object.schema = message.schema;
        if (message.initiator != null && message.hasOwnProperty("initiator"))
            object.initiator = message.initiator;
        if (message.type != null && message.hasOwnProperty("type"))
            object.type = options.enums === String ? $root.Thread.Type[message.type] : message.type;
        if (message.sharing != null && message.hasOwnProperty("sharing"))
            object.sharing = options.enums === String ? $root.Thread.Sharing[message.sharing] : message.sharing;
        if (message.members && message.members.length) {
            object.members = [];
            for (let j = 0; j < message.members.length; ++j)
                object.members[j] = message.members[j];
        }
        if (message.state != null && message.hasOwnProperty("state"))
            object.state = options.enums === String ? $root.Thread.State[message.state] : message.state;
        if (message.head != null && message.hasOwnProperty("head"))
            object.head = message.head;
        if (message.headBlock != null && message.hasOwnProperty("headBlock"))
            object.headBlock = $root.Block.toObject(message.headBlock, options);
        if (message.schemaNode != null && message.hasOwnProperty("schemaNode"))
            object.schemaNode = $root.Node.toObject(message.schemaNode, options);
        if (message.blockCount != null && message.hasOwnProperty("blockCount"))
            object.blockCount = message.blockCount;
        if (message.peerCount != null && message.hasOwnProperty("peerCount"))
            object.peerCount = message.peerCount;
        return object;
    };

    /**
     * Converts this Thread to JSON.
     * @function toJSON
     * @memberof Thread
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Thread.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    /**
     * Type enum.
     * @name Thread.Type
     * @enum {string}
     * @property {number} Private=0 Private value
     * @property {number} ReadOnly=1 ReadOnly value
     * @property {number} Public=2 Public value
     * @property {number} Open=3 Open value
     */
    Thread.Type = (function() {
        const valuesById = {}, values = Object.create(valuesById);
        values[valuesById[0] = "Private"] = 0;
        values[valuesById[1] = "ReadOnly"] = 1;
        values[valuesById[2] = "Public"] = 2;
        values[valuesById[3] = "Open"] = 3;
        return values;
    })();

    /**
     * Sharing enum.
     * @name Thread.Sharing
     * @enum {string}
     * @property {number} NotShared=0 NotShared value
     * @property {number} InviteOnly=1 InviteOnly value
     * @property {number} Shared=2 Shared value
     */
    Thread.Sharing = (function() {
        const valuesById = {}, values = Object.create(valuesById);
        values[valuesById[0] = "NotShared"] = 0;
        values[valuesById[1] = "InviteOnly"] = 1;
        values[valuesById[2] = "Shared"] = 2;
        return values;
    })();

    /**
     * State enum.
     * @name Thread.State
     * @enum {string}
     * @property {number} LoadingBehind=0 LoadingBehind value
     * @property {number} Loaded=1 Loaded value
     * @property {number} LoadingAhead=2 LoadingAhead value
     */
    Thread.State = (function() {
        const valuesById = {}, values = Object.create(valuesById);
        values[valuesById[0] = "LoadingBehind"] = 0;
        values[valuesById[1] = "Loaded"] = 1;
        values[valuesById[2] = "LoadingAhead"] = 2;
        return values;
    })();

    return Thread;
})();

export const ThreadList = $root.ThreadList = (() => {

    /**
     * Properties of a ThreadList.
     * @exports IThreadList
     * @interface IThreadList
     * @property {Array.<IThread>} items ThreadList items
     */

    /**
     * Constructs a new ThreadList.
     * @exports ThreadList
     * @classdesc Represents a ThreadList.
     * @implements IThreadList
     * @constructor
     * @param {IThreadList=} [properties] Properties to set
     */
    function ThreadList(properties) {
        this.items = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ThreadList items.
     * @member {Array.<IThread>} items
     * @memberof ThreadList
     * @instance
     */
    ThreadList.prototype.items = $util.emptyArray;

    /**
     * Creates a new ThreadList instance using the specified properties.
     * @function create
     * @memberof ThreadList
     * @static
     * @param {IThreadList=} [properties] Properties to set
     * @returns {ThreadList} ThreadList instance
     */
    ThreadList.create = function create(properties) {
        return new ThreadList(properties);
    };

    /**
     * Encodes the specified ThreadList message. Does not implicitly {@link ThreadList.verify|verify} messages.
     * @function encode
     * @memberof ThreadList
     * @static
     * @param {IThreadList} message ThreadList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadList.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.items != null && message.items.length)
            for (let i = 0; i < message.items.length; ++i)
                $root.Thread.encode(message.items[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified ThreadList message, length delimited. Does not implicitly {@link ThreadList.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ThreadList
     * @static
     * @param {IThreadList} message ThreadList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadList.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ThreadList message from the specified reader or buffer.
     * @function decode
     * @memberof ThreadList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ThreadList} ThreadList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadList.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ThreadList();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                if (!(message.items && message.items.length))
                    message.items = [];
                message.items.push($root.Thread.decode(reader, reader.uint32()));
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ThreadList message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ThreadList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ThreadList} ThreadList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadList.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ThreadList message.
     * @function verify
     * @memberof ThreadList
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ThreadList.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.items != null && message.hasOwnProperty("items")) {
            if (!Array.isArray(message.items))
                return "items: array expected";
            for (let i = 0; i < message.items.length; ++i) {
                let error = $root.Thread.verify(message.items[i]);
                if (error)
                    return "items." + error;
            }
        }
        return null;
    };

    /**
     * Creates a ThreadList message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ThreadList
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ThreadList} ThreadList
     */
    ThreadList.fromObject = function fromObject(object) {
        if (object instanceof $root.ThreadList)
            return object;
        let message = new $root.ThreadList();
        if (object.items) {
            if (!Array.isArray(object.items))
                throw TypeError(".ThreadList.items: array expected");
            message.items = [];
            for (let i = 0; i < object.items.length; ++i) {
                if (typeof object.items[i] !== "object")
                    throw TypeError(".ThreadList.items: object expected");
                message.items[i] = $root.Thread.fromObject(object.items[i]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from a ThreadList message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ThreadList
     * @static
     * @param {ThreadList} message ThreadList
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ThreadList.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.items = [];
        if (message.items && message.items.length) {
            object.items = [];
            for (let j = 0; j < message.items.length; ++j)
                object.items[j] = $root.Thread.toObject(message.items[j], options);
        }
        return object;
    };

    /**
     * Converts this ThreadList to JSON.
     * @function toJSON
     * @memberof ThreadList
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ThreadList.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ThreadList;
})();

export const ThreadPeer = $root.ThreadPeer = (() => {

    /**
     * Properties of a ThreadPeer.
     * @exports IThreadPeer
     * @interface IThreadPeer
     * @property {string} id ThreadPeer id
     * @property {string} thread ThreadPeer thread
     * @property {boolean} welcomed ThreadPeer welcomed
     */

    /**
     * Constructs a new ThreadPeer.
     * @exports ThreadPeer
     * @classdesc Represents a ThreadPeer.
     * @implements IThreadPeer
     * @constructor
     * @param {IThreadPeer=} [properties] Properties to set
     */
    function ThreadPeer(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ThreadPeer id.
     * @member {string} id
     * @memberof ThreadPeer
     * @instance
     */
    ThreadPeer.prototype.id = "";

    /**
     * ThreadPeer thread.
     * @member {string} thread
     * @memberof ThreadPeer
     * @instance
     */
    ThreadPeer.prototype.thread = "";

    /**
     * ThreadPeer welcomed.
     * @member {boolean} welcomed
     * @memberof ThreadPeer
     * @instance
     */
    ThreadPeer.prototype.welcomed = false;

    /**
     * Creates a new ThreadPeer instance using the specified properties.
     * @function create
     * @memberof ThreadPeer
     * @static
     * @param {IThreadPeer=} [properties] Properties to set
     * @returns {ThreadPeer} ThreadPeer instance
     */
    ThreadPeer.create = function create(properties) {
        return new ThreadPeer(properties);
    };

    /**
     * Encodes the specified ThreadPeer message. Does not implicitly {@link ThreadPeer.verify|verify} messages.
     * @function encode
     * @memberof ThreadPeer
     * @static
     * @param {IThreadPeer} message ThreadPeer message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadPeer.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.thread != null && message.hasOwnProperty("thread"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.thread);
        if (message.welcomed != null && message.hasOwnProperty("welcomed"))
            writer.uint32(/* id 3, wireType 0 =*/24).bool(message.welcomed);
        return writer;
    };

    /**
     * Encodes the specified ThreadPeer message, length delimited. Does not implicitly {@link ThreadPeer.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ThreadPeer
     * @static
     * @param {IThreadPeer} message ThreadPeer message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadPeer.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ThreadPeer message from the specified reader or buffer.
     * @function decode
     * @memberof ThreadPeer
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ThreadPeer} ThreadPeer
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadPeer.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ThreadPeer();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.thread = reader.string();
                break;
            case 3:
                message.welcomed = reader.bool();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ThreadPeer message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ThreadPeer
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ThreadPeer} ThreadPeer
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadPeer.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ThreadPeer message.
     * @function verify
     * @memberof ThreadPeer
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ThreadPeer.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.thread != null && message.hasOwnProperty("thread"))
            if (!$util.isString(message.thread))
                return "thread: string expected";
        if (message.welcomed != null && message.hasOwnProperty("welcomed"))
            if (typeof message.welcomed !== "boolean")
                return "welcomed: boolean expected";
        return null;
    };

    /**
     * Creates a ThreadPeer message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ThreadPeer
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ThreadPeer} ThreadPeer
     */
    ThreadPeer.fromObject = function fromObject(object) {
        if (object instanceof $root.ThreadPeer)
            return object;
        let message = new $root.ThreadPeer();
        if (object.id != null)
            message.id = String(object.id);
        if (object.thread != null)
            message.thread = String(object.thread);
        if (object.welcomed != null)
            message.welcomed = Boolean(object.welcomed);
        return message;
    };

    /**
     * Creates a plain object from a ThreadPeer message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ThreadPeer
     * @static
     * @param {ThreadPeer} message ThreadPeer
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ThreadPeer.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            object.thread = "";
            object.welcomed = false;
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.thread != null && message.hasOwnProperty("thread"))
            object.thread = message.thread;
        if (message.welcomed != null && message.hasOwnProperty("welcomed"))
            object.welcomed = message.welcomed;
        return object;
    };

    /**
     * Converts this ThreadPeer to JSON.
     * @function toJSON
     * @memberof ThreadPeer
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ThreadPeer.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ThreadPeer;
})();

export const Block = $root.Block = (() => {

    /**
     * Properties of a Block.
     * @exports IBlock
     * @interface IBlock
     * @property {string} id Block id
     * @property {string} thread Block thread
     * @property {string} author Block author
     * @property {Block.BlockType} type Block type
     * @property {google.protobuf.ITimestamp} date Block date
     * @property {Array.<string>} parents Block parents
     * @property {string} target Block target
     * @property {string} body Block body
     * @property {IUser} user Block user
     */

    /**
     * Constructs a new Block.
     * @exports Block
     * @classdesc Represents a Block.
     * @implements IBlock
     * @constructor
     * @param {IBlock=} [properties] Properties to set
     */
    function Block(properties) {
        this.parents = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Block id.
     * @member {string} id
     * @memberof Block
     * @instance
     */
    Block.prototype.id = "";

    /**
     * Block thread.
     * @member {string} thread
     * @memberof Block
     * @instance
     */
    Block.prototype.thread = "";

    /**
     * Block author.
     * @member {string} author
     * @memberof Block
     * @instance
     */
    Block.prototype.author = "";

    /**
     * Block type.
     * @member {Block.BlockType} type
     * @memberof Block
     * @instance
     */
    Block.prototype.type = 0;

    /**
     * Block date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof Block
     * @instance
     */
    Block.prototype.date = null;

    /**
     * Block parents.
     * @member {Array.<string>} parents
     * @memberof Block
     * @instance
     */
    Block.prototype.parents = $util.emptyArray;

    /**
     * Block target.
     * @member {string} target
     * @memberof Block
     * @instance
     */
    Block.prototype.target = "";

    /**
     * Block body.
     * @member {string} body
     * @memberof Block
     * @instance
     */
    Block.prototype.body = "";

    /**
     * Block user.
     * @member {IUser} user
     * @memberof Block
     * @instance
     */
    Block.prototype.user = null;

    /**
     * Creates a new Block instance using the specified properties.
     * @function create
     * @memberof Block
     * @static
     * @param {IBlock=} [properties] Properties to set
     * @returns {Block} Block instance
     */
    Block.create = function create(properties) {
        return new Block(properties);
    };

    /**
     * Encodes the specified Block message. Does not implicitly {@link Block.verify|verify} messages.
     * @function encode
     * @memberof Block
     * @static
     * @param {IBlock} message Block message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Block.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.thread != null && message.hasOwnProperty("thread"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.thread);
        if (message.author != null && message.hasOwnProperty("author"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.author);
        if (message.type != null && message.hasOwnProperty("type"))
            writer.uint32(/* id 4, wireType 0 =*/32).int32(message.type);
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 5, wireType 2 =*/42).fork()).ldelim();
        if (message.parents != null && message.parents.length)
            for (let i = 0; i < message.parents.length; ++i)
                writer.uint32(/* id 6, wireType 2 =*/50).string(message.parents[i]);
        if (message.target != null && message.hasOwnProperty("target"))
            writer.uint32(/* id 7, wireType 2 =*/58).string(message.target);
        if (message.body != null && message.hasOwnProperty("body"))
            writer.uint32(/* id 8, wireType 2 =*/66).string(message.body);
        if (message.user != null && message.hasOwnProperty("user"))
            $root.User.encode(message.user, writer.uint32(/* id 101, wireType 2 =*/810).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified Block message, length delimited. Does not implicitly {@link Block.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Block
     * @static
     * @param {IBlock} message Block message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Block.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Block message from the specified reader or buffer.
     * @function decode
     * @memberof Block
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Block} Block
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Block.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Block();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.thread = reader.string();
                break;
            case 3:
                message.author = reader.string();
                break;
            case 4:
                message.type = reader.int32();
                break;
            case 5:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 6:
                if (!(message.parents && message.parents.length))
                    message.parents = [];
                message.parents.push(reader.string());
                break;
            case 7:
                message.target = reader.string();
                break;
            case 8:
                message.body = reader.string();
                break;
            case 101:
                message.user = $root.User.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Block message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Block
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Block} Block
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Block.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Block message.
     * @function verify
     * @memberof Block
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Block.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.thread != null && message.hasOwnProperty("thread"))
            if (!$util.isString(message.thread))
                return "thread: string expected";
        if (message.author != null && message.hasOwnProperty("author"))
            if (!$util.isString(message.author))
                return "author: string expected";
        if (message.type != null && message.hasOwnProperty("type"))
            switch (message.type) {
            default:
                return "type: enum value expected";
            case 0:
            case 1:
            case 2:
            case 3:
            case 4:
            case 5:
            case 6:
            case 7:
            case 8:
            case 9:
            case 50:
                break;
            }
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        if (message.parents != null && message.hasOwnProperty("parents")) {
            if (!Array.isArray(message.parents))
                return "parents: array expected";
            for (let i = 0; i < message.parents.length; ++i)
                if (!$util.isString(message.parents[i]))
                    return "parents: string[] expected";
        }
        if (message.target != null && message.hasOwnProperty("target"))
            if (!$util.isString(message.target))
                return "target: string expected";
        if (message.body != null && message.hasOwnProperty("body"))
            if (!$util.isString(message.body))
                return "body: string expected";
        if (message.user != null && message.hasOwnProperty("user")) {
            let error = $root.User.verify(message.user);
            if (error)
                return "user." + error;
        }
        return null;
    };

    /**
     * Creates a Block message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Block
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Block} Block
     */
    Block.fromObject = function fromObject(object) {
        if (object instanceof $root.Block)
            return object;
        let message = new $root.Block();
        if (object.id != null)
            message.id = String(object.id);
        if (object.thread != null)
            message.thread = String(object.thread);
        if (object.author != null)
            message.author = String(object.author);
        switch (object.type) {
        case "MERGE":
        case 0:
            message.type = 0;
            break;
        case "IGNORE":
        case 1:
            message.type = 1;
            break;
        case "FLAG":
        case 2:
            message.type = 2;
            break;
        case "JOIN":
        case 3:
            message.type = 3;
            break;
        case "ANNOUNCE":
        case 4:
            message.type = 4;
            break;
        case "LEAVE":
        case 5:
            message.type = 5;
            break;
        case "MESSAGE":
        case 6:
            message.type = 6;
            break;
        case "FILES":
        case 7:
            message.type = 7;
            break;
        case "COMMENT":
        case 8:
            message.type = 8;
            break;
        case "LIKE":
        case 9:
            message.type = 9;
            break;
        case "INVITE":
        case 50:
            message.type = 50;
            break;
        }
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".Block.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        if (object.parents) {
            if (!Array.isArray(object.parents))
                throw TypeError(".Block.parents: array expected");
            message.parents = [];
            for (let i = 0; i < object.parents.length; ++i)
                message.parents[i] = String(object.parents[i]);
        }
        if (object.target != null)
            message.target = String(object.target);
        if (object.body != null)
            message.body = String(object.body);
        if (object.user != null) {
            if (typeof object.user !== "object")
                throw TypeError(".Block.user: object expected");
            message.user = $root.User.fromObject(object.user);
        }
        return message;
    };

    /**
     * Creates a plain object from a Block message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Block
     * @static
     * @param {Block} message Block
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Block.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.parents = [];
        if (options.defaults) {
            object.id = "";
            object.thread = "";
            object.author = "";
            object.type = options.enums === String ? "MERGE" : 0;
            object.date = null;
            object.target = "";
            object.body = "";
            object.user = null;
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.thread != null && message.hasOwnProperty("thread"))
            object.thread = message.thread;
        if (message.author != null && message.hasOwnProperty("author"))
            object.author = message.author;
        if (message.type != null && message.hasOwnProperty("type"))
            object.type = options.enums === String ? $root.Block.BlockType[message.type] : message.type;
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        if (message.parents && message.parents.length) {
            object.parents = [];
            for (let j = 0; j < message.parents.length; ++j)
                object.parents[j] = message.parents[j];
        }
        if (message.target != null && message.hasOwnProperty("target"))
            object.target = message.target;
        if (message.body != null && message.hasOwnProperty("body"))
            object.body = message.body;
        if (message.user != null && message.hasOwnProperty("user"))
            object.user = $root.User.toObject(message.user, options);
        return object;
    };

    /**
     * Converts this Block to JSON.
     * @function toJSON
     * @memberof Block
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Block.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    /**
     * BlockType enum.
     * @name Block.BlockType
     * @enum {string}
     * @property {number} MERGE=0 MERGE value
     * @property {number} IGNORE=1 IGNORE value
     * @property {number} FLAG=2 FLAG value
     * @property {number} JOIN=3 JOIN value
     * @property {number} ANNOUNCE=4 ANNOUNCE value
     * @property {number} LEAVE=5 LEAVE value
     * @property {number} MESSAGE=6 MESSAGE value
     * @property {number} FILES=7 FILES value
     * @property {number} COMMENT=8 COMMENT value
     * @property {number} LIKE=9 LIKE value
     * @property {number} INVITE=50 INVITE value
     */
    Block.BlockType = (function() {
        const valuesById = {}, values = Object.create(valuesById);
        values[valuesById[0] = "MERGE"] = 0;
        values[valuesById[1] = "IGNORE"] = 1;
        values[valuesById[2] = "FLAG"] = 2;
        values[valuesById[3] = "JOIN"] = 3;
        values[valuesById[4] = "ANNOUNCE"] = 4;
        values[valuesById[5] = "LEAVE"] = 5;
        values[valuesById[6] = "MESSAGE"] = 6;
        values[valuesById[7] = "FILES"] = 7;
        values[valuesById[8] = "COMMENT"] = 8;
        values[valuesById[9] = "LIKE"] = 9;
        values[valuesById[50] = "INVITE"] = 50;
        return values;
    })();

    return Block;
})();

export const BlockList = $root.BlockList = (() => {

    /**
     * Properties of a BlockList.
     * @exports IBlockList
     * @interface IBlockList
     * @property {Array.<IBlock>} items BlockList items
     */

    /**
     * Constructs a new BlockList.
     * @exports BlockList
     * @classdesc Represents a BlockList.
     * @implements IBlockList
     * @constructor
     * @param {IBlockList=} [properties] Properties to set
     */
    function BlockList(properties) {
        this.items = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * BlockList items.
     * @member {Array.<IBlock>} items
     * @memberof BlockList
     * @instance
     */
    BlockList.prototype.items = $util.emptyArray;

    /**
     * Creates a new BlockList instance using the specified properties.
     * @function create
     * @memberof BlockList
     * @static
     * @param {IBlockList=} [properties] Properties to set
     * @returns {BlockList} BlockList instance
     */
    BlockList.create = function create(properties) {
        return new BlockList(properties);
    };

    /**
     * Encodes the specified BlockList message. Does not implicitly {@link BlockList.verify|verify} messages.
     * @function encode
     * @memberof BlockList
     * @static
     * @param {IBlockList} message BlockList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    BlockList.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.items != null && message.items.length)
            for (let i = 0; i < message.items.length; ++i)
                $root.Block.encode(message.items[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified BlockList message, length delimited. Does not implicitly {@link BlockList.verify|verify} messages.
     * @function encodeDelimited
     * @memberof BlockList
     * @static
     * @param {IBlockList} message BlockList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    BlockList.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a BlockList message from the specified reader or buffer.
     * @function decode
     * @memberof BlockList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {BlockList} BlockList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    BlockList.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.BlockList();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                if (!(message.items && message.items.length))
                    message.items = [];
                message.items.push($root.Block.decode(reader, reader.uint32()));
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a BlockList message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof BlockList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {BlockList} BlockList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    BlockList.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a BlockList message.
     * @function verify
     * @memberof BlockList
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    BlockList.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.items != null && message.hasOwnProperty("items")) {
            if (!Array.isArray(message.items))
                return "items: array expected";
            for (let i = 0; i < message.items.length; ++i) {
                let error = $root.Block.verify(message.items[i]);
                if (error)
                    return "items." + error;
            }
        }
        return null;
    };

    /**
     * Creates a BlockList message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof BlockList
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {BlockList} BlockList
     */
    BlockList.fromObject = function fromObject(object) {
        if (object instanceof $root.BlockList)
            return object;
        let message = new $root.BlockList();
        if (object.items) {
            if (!Array.isArray(object.items))
                throw TypeError(".BlockList.items: array expected");
            message.items = [];
            for (let i = 0; i < object.items.length; ++i) {
                if (typeof object.items[i] !== "object")
                    throw TypeError(".BlockList.items: object expected");
                message.items[i] = $root.Block.fromObject(object.items[i]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from a BlockList message. Also converts values to other types if specified.
     * @function toObject
     * @memberof BlockList
     * @static
     * @param {BlockList} message BlockList
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    BlockList.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.items = [];
        if (message.items && message.items.length) {
            object.items = [];
            for (let j = 0; j < message.items.length; ++j)
                object.items[j] = $root.Block.toObject(message.items[j], options);
        }
        return object;
    };

    /**
     * Converts this BlockList to JSON.
     * @function toJSON
     * @memberof BlockList
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    BlockList.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return BlockList;
})();

export const BlockMessage = $root.BlockMessage = (() => {

    /**
     * Properties of a BlockMessage.
     * @exports IBlockMessage
     * @interface IBlockMessage
     * @property {string} id BlockMessage id
     * @property {string} thread BlockMessage thread
     * @property {IEnvelope} env BlockMessage env
     * @property {google.protobuf.ITimestamp} date BlockMessage date
     */

    /**
     * Constructs a new BlockMessage.
     * @exports BlockMessage
     * @classdesc Represents a BlockMessage.
     * @implements IBlockMessage
     * @constructor
     * @param {IBlockMessage=} [properties] Properties to set
     */
    function BlockMessage(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * BlockMessage id.
     * @member {string} id
     * @memberof BlockMessage
     * @instance
     */
    BlockMessage.prototype.id = "";

    /**
     * BlockMessage thread.
     * @member {string} thread
     * @memberof BlockMessage
     * @instance
     */
    BlockMessage.prototype.thread = "";

    /**
     * BlockMessage env.
     * @member {IEnvelope} env
     * @memberof BlockMessage
     * @instance
     */
    BlockMessage.prototype.env = null;

    /**
     * BlockMessage date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof BlockMessage
     * @instance
     */
    BlockMessage.prototype.date = null;

    /**
     * Creates a new BlockMessage instance using the specified properties.
     * @function create
     * @memberof BlockMessage
     * @static
     * @param {IBlockMessage=} [properties] Properties to set
     * @returns {BlockMessage} BlockMessage instance
     */
    BlockMessage.create = function create(properties) {
        return new BlockMessage(properties);
    };

    /**
     * Encodes the specified BlockMessage message. Does not implicitly {@link BlockMessage.verify|verify} messages.
     * @function encode
     * @memberof BlockMessage
     * @static
     * @param {IBlockMessage} message BlockMessage message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    BlockMessage.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.thread != null && message.hasOwnProperty("thread"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.thread);
        if (message.env != null && message.hasOwnProperty("env"))
            $root.Envelope.encode(message.env, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified BlockMessage message, length delimited. Does not implicitly {@link BlockMessage.verify|verify} messages.
     * @function encodeDelimited
     * @memberof BlockMessage
     * @static
     * @param {IBlockMessage} message BlockMessage message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    BlockMessage.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a BlockMessage message from the specified reader or buffer.
     * @function decode
     * @memberof BlockMessage
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {BlockMessage} BlockMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    BlockMessage.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.BlockMessage();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.thread = reader.string();
                break;
            case 3:
                message.env = $root.Envelope.decode(reader, reader.uint32());
                break;
            case 4:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a BlockMessage message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof BlockMessage
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {BlockMessage} BlockMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    BlockMessage.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a BlockMessage message.
     * @function verify
     * @memberof BlockMessage
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    BlockMessage.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.thread != null && message.hasOwnProperty("thread"))
            if (!$util.isString(message.thread))
                return "thread: string expected";
        if (message.env != null && message.hasOwnProperty("env")) {
            let error = $root.Envelope.verify(message.env);
            if (error)
                return "env." + error;
        }
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        return null;
    };

    /**
     * Creates a BlockMessage message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof BlockMessage
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {BlockMessage} BlockMessage
     */
    BlockMessage.fromObject = function fromObject(object) {
        if (object instanceof $root.BlockMessage)
            return object;
        let message = new $root.BlockMessage();
        if (object.id != null)
            message.id = String(object.id);
        if (object.thread != null)
            message.thread = String(object.thread);
        if (object.env != null) {
            if (typeof object.env !== "object")
                throw TypeError(".BlockMessage.env: object expected");
            message.env = $root.Envelope.fromObject(object.env);
        }
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".BlockMessage.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        return message;
    };

    /**
     * Creates a plain object from a BlockMessage message. Also converts values to other types if specified.
     * @function toObject
     * @memberof BlockMessage
     * @static
     * @param {BlockMessage} message BlockMessage
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    BlockMessage.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            object.thread = "";
            object.env = null;
            object.date = null;
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.thread != null && message.hasOwnProperty("thread"))
            object.thread = message.thread;
        if (message.env != null && message.hasOwnProperty("env"))
            object.env = $root.Envelope.toObject(message.env, options);
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        return object;
    };

    /**
     * Converts this BlockMessage to JSON.
     * @function toJSON
     * @memberof BlockMessage
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    BlockMessage.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return BlockMessage;
})();

export const Invite = $root.Invite = (() => {

    /**
     * Properties of an Invite.
     * @exports IInvite
     * @interface IInvite
     * @property {string} id Invite id
     * @property {Uint8Array} block Invite block
     * @property {string} name Invite name
     * @property {IContact} inviter Invite inviter
     * @property {google.protobuf.ITimestamp} date Invite date
     */

    /**
     * Constructs a new Invite.
     * @exports Invite
     * @classdesc Represents an Invite.
     * @implements IInvite
     * @constructor
     * @param {IInvite=} [properties] Properties to set
     */
    function Invite(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Invite id.
     * @member {string} id
     * @memberof Invite
     * @instance
     */
    Invite.prototype.id = "";

    /**
     * Invite block.
     * @member {Uint8Array} block
     * @memberof Invite
     * @instance
     */
    Invite.prototype.block = $util.newBuffer([]);

    /**
     * Invite name.
     * @member {string} name
     * @memberof Invite
     * @instance
     */
    Invite.prototype.name = "";

    /**
     * Invite inviter.
     * @member {IContact} inviter
     * @memberof Invite
     * @instance
     */
    Invite.prototype.inviter = null;

    /**
     * Invite date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof Invite
     * @instance
     */
    Invite.prototype.date = null;

    /**
     * Creates a new Invite instance using the specified properties.
     * @function create
     * @memberof Invite
     * @static
     * @param {IInvite=} [properties] Properties to set
     * @returns {Invite} Invite instance
     */
    Invite.create = function create(properties) {
        return new Invite(properties);
    };

    /**
     * Encodes the specified Invite message. Does not implicitly {@link Invite.verify|verify} messages.
     * @function encode
     * @memberof Invite
     * @static
     * @param {IInvite} message Invite message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Invite.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.block != null && message.hasOwnProperty("block"))
            writer.uint32(/* id 2, wireType 2 =*/18).bytes(message.block);
        if (message.name != null && message.hasOwnProperty("name"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.name);
        if (message.inviter != null && message.hasOwnProperty("inviter"))
            $root.Contact.encode(message.inviter, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 5, wireType 2 =*/42).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified Invite message, length delimited. Does not implicitly {@link Invite.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Invite
     * @static
     * @param {IInvite} message Invite message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Invite.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes an Invite message from the specified reader or buffer.
     * @function decode
     * @memberof Invite
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Invite} Invite
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Invite.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Invite();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.block = reader.bytes();
                break;
            case 3:
                message.name = reader.string();
                break;
            case 4:
                message.inviter = $root.Contact.decode(reader, reader.uint32());
                break;
            case 5:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes an Invite message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Invite
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Invite} Invite
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Invite.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies an Invite message.
     * @function verify
     * @memberof Invite
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Invite.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.block != null && message.hasOwnProperty("block"))
            if (!(message.block && typeof message.block.length === "number" || $util.isString(message.block)))
                return "block: buffer expected";
        if (message.name != null && message.hasOwnProperty("name"))
            if (!$util.isString(message.name))
                return "name: string expected";
        if (message.inviter != null && message.hasOwnProperty("inviter")) {
            let error = $root.Contact.verify(message.inviter);
            if (error)
                return "inviter." + error;
        }
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        return null;
    };

    /**
     * Creates an Invite message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Invite
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Invite} Invite
     */
    Invite.fromObject = function fromObject(object) {
        if (object instanceof $root.Invite)
            return object;
        let message = new $root.Invite();
        if (object.id != null)
            message.id = String(object.id);
        if (object.block != null)
            if (typeof object.block === "string")
                $util.base64.decode(object.block, message.block = $util.newBuffer($util.base64.length(object.block)), 0);
            else if (object.block.length)
                message.block = object.block;
        if (object.name != null)
            message.name = String(object.name);
        if (object.inviter != null) {
            if (typeof object.inviter !== "object")
                throw TypeError(".Invite.inviter: object expected");
            message.inviter = $root.Contact.fromObject(object.inviter);
        }
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".Invite.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        return message;
    };

    /**
     * Creates a plain object from an Invite message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Invite
     * @static
     * @param {Invite} message Invite
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Invite.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            if (options.bytes === String)
                object.block = "";
            else {
                object.block = [];
                if (options.bytes !== Array)
                    object.block = $util.newBuffer(object.block);
            }
            object.name = "";
            object.inviter = null;
            object.date = null;
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.block != null && message.hasOwnProperty("block"))
            object.block = options.bytes === String ? $util.base64.encode(message.block, 0, message.block.length) : options.bytes === Array ? Array.prototype.slice.call(message.block) : message.block;
        if (message.name != null && message.hasOwnProperty("name"))
            object.name = message.name;
        if (message.inviter != null && message.hasOwnProperty("inviter"))
            object.inviter = $root.Contact.toObject(message.inviter, options);
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        return object;
    };

    /**
     * Converts this Invite to JSON.
     * @function toJSON
     * @memberof Invite
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Invite.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Invite;
})();

export const InviteList = $root.InviteList = (() => {

    /**
     * Properties of an InviteList.
     * @exports IInviteList
     * @interface IInviteList
     * @property {Array.<IInvite>} items InviteList items
     */

    /**
     * Constructs a new InviteList.
     * @exports InviteList
     * @classdesc Represents an InviteList.
     * @implements IInviteList
     * @constructor
     * @param {IInviteList=} [properties] Properties to set
     */
    function InviteList(properties) {
        this.items = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * InviteList items.
     * @member {Array.<IInvite>} items
     * @memberof InviteList
     * @instance
     */
    InviteList.prototype.items = $util.emptyArray;

    /**
     * Creates a new InviteList instance using the specified properties.
     * @function create
     * @memberof InviteList
     * @static
     * @param {IInviteList=} [properties] Properties to set
     * @returns {InviteList} InviteList instance
     */
    InviteList.create = function create(properties) {
        return new InviteList(properties);
    };

    /**
     * Encodes the specified InviteList message. Does not implicitly {@link InviteList.verify|verify} messages.
     * @function encode
     * @memberof InviteList
     * @static
     * @param {IInviteList} message InviteList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    InviteList.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.items != null && message.items.length)
            for (let i = 0; i < message.items.length; ++i)
                $root.Invite.encode(message.items[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified InviteList message, length delimited. Does not implicitly {@link InviteList.verify|verify} messages.
     * @function encodeDelimited
     * @memberof InviteList
     * @static
     * @param {IInviteList} message InviteList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    InviteList.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes an InviteList message from the specified reader or buffer.
     * @function decode
     * @memberof InviteList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {InviteList} InviteList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    InviteList.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.InviteList();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                if (!(message.items && message.items.length))
                    message.items = [];
                message.items.push($root.Invite.decode(reader, reader.uint32()));
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes an InviteList message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof InviteList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {InviteList} InviteList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    InviteList.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies an InviteList message.
     * @function verify
     * @memberof InviteList
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    InviteList.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.items != null && message.hasOwnProperty("items")) {
            if (!Array.isArray(message.items))
                return "items: array expected";
            for (let i = 0; i < message.items.length; ++i) {
                let error = $root.Invite.verify(message.items[i]);
                if (error)
                    return "items." + error;
            }
        }
        return null;
    };

    /**
     * Creates an InviteList message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof InviteList
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {InviteList} InviteList
     */
    InviteList.fromObject = function fromObject(object) {
        if (object instanceof $root.InviteList)
            return object;
        let message = new $root.InviteList();
        if (object.items) {
            if (!Array.isArray(object.items))
                throw TypeError(".InviteList.items: array expected");
            message.items = [];
            for (let i = 0; i < object.items.length; ++i) {
                if (typeof object.items[i] !== "object")
                    throw TypeError(".InviteList.items: object expected");
                message.items[i] = $root.Invite.fromObject(object.items[i]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from an InviteList message. Also converts values to other types if specified.
     * @function toObject
     * @memberof InviteList
     * @static
     * @param {InviteList} message InviteList
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    InviteList.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.items = [];
        if (message.items && message.items.length) {
            object.items = [];
            for (let j = 0; j < message.items.length; ++j)
                object.items[j] = $root.Invite.toObject(message.items[j], options);
        }
        return object;
    };

    /**
     * Converts this InviteList to JSON.
     * @function toJSON
     * @memberof InviteList
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    InviteList.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return InviteList;
})();

export const FileIndex = $root.FileIndex = (() => {

    /**
     * Properties of a FileIndex.
     * @exports IFileIndex
     * @interface IFileIndex
     * @property {string} mill FileIndex mill
     * @property {string} checksum FileIndex checksum
     * @property {string} source FileIndex source
     * @property {string} opts FileIndex opts
     * @property {string} hash FileIndex hash
     * @property {string} key FileIndex key
     * @property {string} media FileIndex media
     * @property {string} name FileIndex name
     * @property {number|Long} size FileIndex size
     * @property {google.protobuf.ITimestamp} added FileIndex added
     * @property {google.protobuf.IStruct} meta FileIndex meta
     * @property {Array.<string>} targets FileIndex targets
     */

    /**
     * Constructs a new FileIndex.
     * @exports FileIndex
     * @classdesc Represents a FileIndex.
     * @implements IFileIndex
     * @constructor
     * @param {IFileIndex=} [properties] Properties to set
     */
    function FileIndex(properties) {
        this.targets = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * FileIndex mill.
     * @member {string} mill
     * @memberof FileIndex
     * @instance
     */
    FileIndex.prototype.mill = "";

    /**
     * FileIndex checksum.
     * @member {string} checksum
     * @memberof FileIndex
     * @instance
     */
    FileIndex.prototype.checksum = "";

    /**
     * FileIndex source.
     * @member {string} source
     * @memberof FileIndex
     * @instance
     */
    FileIndex.prototype.source = "";

    /**
     * FileIndex opts.
     * @member {string} opts
     * @memberof FileIndex
     * @instance
     */
    FileIndex.prototype.opts = "";

    /**
     * FileIndex hash.
     * @member {string} hash
     * @memberof FileIndex
     * @instance
     */
    FileIndex.prototype.hash = "";

    /**
     * FileIndex key.
     * @member {string} key
     * @memberof FileIndex
     * @instance
     */
    FileIndex.prototype.key = "";

    /**
     * FileIndex media.
     * @member {string} media
     * @memberof FileIndex
     * @instance
     */
    FileIndex.prototype.media = "";

    /**
     * FileIndex name.
     * @member {string} name
     * @memberof FileIndex
     * @instance
     */
    FileIndex.prototype.name = "";

    /**
     * FileIndex size.
     * @member {number|Long} size
     * @memberof FileIndex
     * @instance
     */
    FileIndex.prototype.size = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

    /**
     * FileIndex added.
     * @member {google.protobuf.ITimestamp} added
     * @memberof FileIndex
     * @instance
     */
    FileIndex.prototype.added = null;

    /**
     * FileIndex meta.
     * @member {google.protobuf.IStruct} meta
     * @memberof FileIndex
     * @instance
     */
    FileIndex.prototype.meta = null;

    /**
     * FileIndex targets.
     * @member {Array.<string>} targets
     * @memberof FileIndex
     * @instance
     */
    FileIndex.prototype.targets = $util.emptyArray;

    /**
     * Creates a new FileIndex instance using the specified properties.
     * @function create
     * @memberof FileIndex
     * @static
     * @param {IFileIndex=} [properties] Properties to set
     * @returns {FileIndex} FileIndex instance
     */
    FileIndex.create = function create(properties) {
        return new FileIndex(properties);
    };

    /**
     * Encodes the specified FileIndex message. Does not implicitly {@link FileIndex.verify|verify} messages.
     * @function encode
     * @memberof FileIndex
     * @static
     * @param {IFileIndex} message FileIndex message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    FileIndex.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.mill != null && message.hasOwnProperty("mill"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.mill);
        if (message.checksum != null && message.hasOwnProperty("checksum"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.checksum);
        if (message.source != null && message.hasOwnProperty("source"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.source);
        if (message.opts != null && message.hasOwnProperty("opts"))
            writer.uint32(/* id 4, wireType 2 =*/34).string(message.opts);
        if (message.hash != null && message.hasOwnProperty("hash"))
            writer.uint32(/* id 5, wireType 2 =*/42).string(message.hash);
        if (message.key != null && message.hasOwnProperty("key"))
            writer.uint32(/* id 6, wireType 2 =*/50).string(message.key);
        if (message.media != null && message.hasOwnProperty("media"))
            writer.uint32(/* id 7, wireType 2 =*/58).string(message.media);
        if (message.name != null && message.hasOwnProperty("name"))
            writer.uint32(/* id 8, wireType 2 =*/66).string(message.name);
        if (message.size != null && message.hasOwnProperty("size"))
            writer.uint32(/* id 9, wireType 0 =*/72).int64(message.size);
        if (message.added != null && message.hasOwnProperty("added"))
            $root.google.protobuf.Timestamp.encode(message.added, writer.uint32(/* id 10, wireType 2 =*/82).fork()).ldelim();
        if (message.meta != null && message.hasOwnProperty("meta"))
            $root.google.protobuf.Struct.encode(message.meta, writer.uint32(/* id 11, wireType 2 =*/90).fork()).ldelim();
        if (message.targets != null && message.targets.length)
            for (let i = 0; i < message.targets.length; ++i)
                writer.uint32(/* id 12, wireType 2 =*/98).string(message.targets[i]);
        return writer;
    };

    /**
     * Encodes the specified FileIndex message, length delimited. Does not implicitly {@link FileIndex.verify|verify} messages.
     * @function encodeDelimited
     * @memberof FileIndex
     * @static
     * @param {IFileIndex} message FileIndex message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    FileIndex.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a FileIndex message from the specified reader or buffer.
     * @function decode
     * @memberof FileIndex
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {FileIndex} FileIndex
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    FileIndex.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.FileIndex();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.mill = reader.string();
                break;
            case 2:
                message.checksum = reader.string();
                break;
            case 3:
                message.source = reader.string();
                break;
            case 4:
                message.opts = reader.string();
                break;
            case 5:
                message.hash = reader.string();
                break;
            case 6:
                message.key = reader.string();
                break;
            case 7:
                message.media = reader.string();
                break;
            case 8:
                message.name = reader.string();
                break;
            case 9:
                message.size = reader.int64();
                break;
            case 10:
                message.added = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 11:
                message.meta = $root.google.protobuf.Struct.decode(reader, reader.uint32());
                break;
            case 12:
                if (!(message.targets && message.targets.length))
                    message.targets = [];
                message.targets.push(reader.string());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a FileIndex message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof FileIndex
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {FileIndex} FileIndex
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    FileIndex.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a FileIndex message.
     * @function verify
     * @memberof FileIndex
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    FileIndex.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.mill != null && message.hasOwnProperty("mill"))
            if (!$util.isString(message.mill))
                return "mill: string expected";
        if (message.checksum != null && message.hasOwnProperty("checksum"))
            if (!$util.isString(message.checksum))
                return "checksum: string expected";
        if (message.source != null && message.hasOwnProperty("source"))
            if (!$util.isString(message.source))
                return "source: string expected";
        if (message.opts != null && message.hasOwnProperty("opts"))
            if (!$util.isString(message.opts))
                return "opts: string expected";
        if (message.hash != null && message.hasOwnProperty("hash"))
            if (!$util.isString(message.hash))
                return "hash: string expected";
        if (message.key != null && message.hasOwnProperty("key"))
            if (!$util.isString(message.key))
                return "key: string expected";
        if (message.media != null && message.hasOwnProperty("media"))
            if (!$util.isString(message.media))
                return "media: string expected";
        if (message.name != null && message.hasOwnProperty("name"))
            if (!$util.isString(message.name))
                return "name: string expected";
        if (message.size != null && message.hasOwnProperty("size"))
            if (!$util.isInteger(message.size) && !(message.size && $util.isInteger(message.size.low) && $util.isInteger(message.size.high)))
                return "size: integer|Long expected";
        if (message.added != null && message.hasOwnProperty("added")) {
            let error = $root.google.protobuf.Timestamp.verify(message.added);
            if (error)
                return "added." + error;
        }
        if (message.meta != null && message.hasOwnProperty("meta")) {
            let error = $root.google.protobuf.Struct.verify(message.meta);
            if (error)
                return "meta." + error;
        }
        if (message.targets != null && message.hasOwnProperty("targets")) {
            if (!Array.isArray(message.targets))
                return "targets: array expected";
            for (let i = 0; i < message.targets.length; ++i)
                if (!$util.isString(message.targets[i]))
                    return "targets: string[] expected";
        }
        return null;
    };

    /**
     * Creates a FileIndex message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof FileIndex
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {FileIndex} FileIndex
     */
    FileIndex.fromObject = function fromObject(object) {
        if (object instanceof $root.FileIndex)
            return object;
        let message = new $root.FileIndex();
        if (object.mill != null)
            message.mill = String(object.mill);
        if (object.checksum != null)
            message.checksum = String(object.checksum);
        if (object.source != null)
            message.source = String(object.source);
        if (object.opts != null)
            message.opts = String(object.opts);
        if (object.hash != null)
            message.hash = String(object.hash);
        if (object.key != null)
            message.key = String(object.key);
        if (object.media != null)
            message.media = String(object.media);
        if (object.name != null)
            message.name = String(object.name);
        if (object.size != null)
            if ($util.Long)
                (message.size = $util.Long.fromValue(object.size)).unsigned = false;
            else if (typeof object.size === "string")
                message.size = parseInt(object.size, 10);
            else if (typeof object.size === "number")
                message.size = object.size;
            else if (typeof object.size === "object")
                message.size = new $util.LongBits(object.size.low >>> 0, object.size.high >>> 0).toNumber();
        if (object.added != null) {
            if (typeof object.added !== "object")
                throw TypeError(".FileIndex.added: object expected");
            message.added = $root.google.protobuf.Timestamp.fromObject(object.added);
        }
        if (object.meta != null) {
            if (typeof object.meta !== "object")
                throw TypeError(".FileIndex.meta: object expected");
            message.meta = $root.google.protobuf.Struct.fromObject(object.meta);
        }
        if (object.targets) {
            if (!Array.isArray(object.targets))
                throw TypeError(".FileIndex.targets: array expected");
            message.targets = [];
            for (let i = 0; i < object.targets.length; ++i)
                message.targets[i] = String(object.targets[i]);
        }
        return message;
    };

    /**
     * Creates a plain object from a FileIndex message. Also converts values to other types if specified.
     * @function toObject
     * @memberof FileIndex
     * @static
     * @param {FileIndex} message FileIndex
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    FileIndex.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.targets = [];
        if (options.defaults) {
            object.mill = "";
            object.checksum = "";
            object.source = "";
            object.opts = "";
            object.hash = "";
            object.key = "";
            object.media = "";
            object.name = "";
            if ($util.Long) {
                let long = new $util.Long(0, 0, false);
                object.size = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
            } else
                object.size = options.longs === String ? "0" : 0;
            object.added = null;
            object.meta = null;
        }
        if (message.mill != null && message.hasOwnProperty("mill"))
            object.mill = message.mill;
        if (message.checksum != null && message.hasOwnProperty("checksum"))
            object.checksum = message.checksum;
        if (message.source != null && message.hasOwnProperty("source"))
            object.source = message.source;
        if (message.opts != null && message.hasOwnProperty("opts"))
            object.opts = message.opts;
        if (message.hash != null && message.hasOwnProperty("hash"))
            object.hash = message.hash;
        if (message.key != null && message.hasOwnProperty("key"))
            object.key = message.key;
        if (message.media != null && message.hasOwnProperty("media"))
            object.media = message.media;
        if (message.name != null && message.hasOwnProperty("name"))
            object.name = message.name;
        if (message.size != null && message.hasOwnProperty("size"))
            if (typeof message.size === "number")
                object.size = options.longs === String ? String(message.size) : message.size;
            else
                object.size = options.longs === String ? $util.Long.prototype.toString.call(message.size) : options.longs === Number ? new $util.LongBits(message.size.low >>> 0, message.size.high >>> 0).toNumber() : message.size;
        if (message.added != null && message.hasOwnProperty("added"))
            object.added = $root.google.protobuf.Timestamp.toObject(message.added, options);
        if (message.meta != null && message.hasOwnProperty("meta"))
            object.meta = $root.google.protobuf.Struct.toObject(message.meta, options);
        if (message.targets && message.targets.length) {
            object.targets = [];
            for (let j = 0; j < message.targets.length; ++j)
                object.targets[j] = message.targets[j];
        }
        return object;
    };

    /**
     * Converts this FileIndex to JSON.
     * @function toJSON
     * @memberof FileIndex
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    FileIndex.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return FileIndex;
})();

export const Node = $root.Node = (() => {

    /**
     * Properties of a Node.
     * @exports INode
     * @interface INode
     * @property {string} name Node name
     * @property {boolean} pin Node pin
     * @property {boolean} plaintext Node plaintext
     * @property {string} mill Node mill
     * @property {Object.<string,string>} opts Node opts
     * @property {google.protobuf.IStruct} jsonSchema Node jsonSchema
     * @property {Object.<string,ILink>} links Node links
     */

    /**
     * Constructs a new Node.
     * @exports Node
     * @classdesc Represents a Node.
     * @implements INode
     * @constructor
     * @param {INode=} [properties] Properties to set
     */
    function Node(properties) {
        this.opts = {};
        this.links = {};
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Node name.
     * @member {string} name
     * @memberof Node
     * @instance
     */
    Node.prototype.name = "";

    /**
     * Node pin.
     * @member {boolean} pin
     * @memberof Node
     * @instance
     */
    Node.prototype.pin = false;

    /**
     * Node plaintext.
     * @member {boolean} plaintext
     * @memberof Node
     * @instance
     */
    Node.prototype.plaintext = false;

    /**
     * Node mill.
     * @member {string} mill
     * @memberof Node
     * @instance
     */
    Node.prototype.mill = "";

    /**
     * Node opts.
     * @member {Object.<string,string>} opts
     * @memberof Node
     * @instance
     */
    Node.prototype.opts = $util.emptyObject;

    /**
     * Node jsonSchema.
     * @member {google.protobuf.IStruct} jsonSchema
     * @memberof Node
     * @instance
     */
    Node.prototype.jsonSchema = null;

    /**
     * Node links.
     * @member {Object.<string,ILink>} links
     * @memberof Node
     * @instance
     */
    Node.prototype.links = $util.emptyObject;

    /**
     * Creates a new Node instance using the specified properties.
     * @function create
     * @memberof Node
     * @static
     * @param {INode=} [properties] Properties to set
     * @returns {Node} Node instance
     */
    Node.create = function create(properties) {
        return new Node(properties);
    };

    /**
     * Encodes the specified Node message. Does not implicitly {@link Node.verify|verify} messages.
     * @function encode
     * @memberof Node
     * @static
     * @param {INode} message Node message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Node.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.name != null && message.hasOwnProperty("name"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.name);
        if (message.pin != null && message.hasOwnProperty("pin"))
            writer.uint32(/* id 2, wireType 0 =*/16).bool(message.pin);
        if (message.plaintext != null && message.hasOwnProperty("plaintext"))
            writer.uint32(/* id 3, wireType 0 =*/24).bool(message.plaintext);
        if (message.mill != null && message.hasOwnProperty("mill"))
            writer.uint32(/* id 4, wireType 2 =*/34).string(message.mill);
        if (message.opts != null && message.hasOwnProperty("opts"))
            for (let keys = Object.keys(message.opts), i = 0; i < keys.length; ++i)
                writer.uint32(/* id 5, wireType 2 =*/42).fork().uint32(/* id 1, wireType 2 =*/10).string(keys[i]).uint32(/* id 2, wireType 2 =*/18).string(message.opts[keys[i]]).ldelim();
        if (message.jsonSchema != null && message.hasOwnProperty("jsonSchema"))
            $root.google.protobuf.Struct.encode(message.jsonSchema, writer.uint32(/* id 6, wireType 2 =*/50).fork()).ldelim();
        if (message.links != null && message.hasOwnProperty("links"))
            for (let keys = Object.keys(message.links), i = 0; i < keys.length; ++i) {
                writer.uint32(/* id 8, wireType 2 =*/66).fork().uint32(/* id 1, wireType 2 =*/10).string(keys[i]);
                $root.Link.encode(message.links[keys[i]], writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim().ldelim();
            }
        return writer;
    };

    /**
     * Encodes the specified Node message, length delimited. Does not implicitly {@link Node.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Node
     * @static
     * @param {INode} message Node message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Node.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Node message from the specified reader or buffer.
     * @function decode
     * @memberof Node
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Node} Node
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Node.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Node(), key;
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.name = reader.string();
                break;
            case 2:
                message.pin = reader.bool();
                break;
            case 3:
                message.plaintext = reader.bool();
                break;
            case 4:
                message.mill = reader.string();
                break;
            case 5:
                reader.skip().pos++;
                if (message.opts === $util.emptyObject)
                    message.opts = {};
                key = reader.string();
                reader.pos++;
                message.opts[key] = reader.string();
                break;
            case 6:
                message.jsonSchema = $root.google.protobuf.Struct.decode(reader, reader.uint32());
                break;
            case 8:
                reader.skip().pos++;
                if (message.links === $util.emptyObject)
                    message.links = {};
                key = reader.string();
                reader.pos++;
                message.links[key] = $root.Link.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Node message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Node
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Node} Node
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Node.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Node message.
     * @function verify
     * @memberof Node
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Node.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.name != null && message.hasOwnProperty("name"))
            if (!$util.isString(message.name))
                return "name: string expected";
        if (message.pin != null && message.hasOwnProperty("pin"))
            if (typeof message.pin !== "boolean")
                return "pin: boolean expected";
        if (message.plaintext != null && message.hasOwnProperty("plaintext"))
            if (typeof message.plaintext !== "boolean")
                return "plaintext: boolean expected";
        if (message.mill != null && message.hasOwnProperty("mill"))
            if (!$util.isString(message.mill))
                return "mill: string expected";
        if (message.opts != null && message.hasOwnProperty("opts")) {
            if (!$util.isObject(message.opts))
                return "opts: object expected";
            let key = Object.keys(message.opts);
            for (let i = 0; i < key.length; ++i)
                if (!$util.isString(message.opts[key[i]]))
                    return "opts: string{k:string} expected";
        }
        if (message.jsonSchema != null && message.hasOwnProperty("jsonSchema")) {
            let error = $root.google.protobuf.Struct.verify(message.jsonSchema);
            if (error)
                return "jsonSchema." + error;
        }
        if (message.links != null && message.hasOwnProperty("links")) {
            if (!$util.isObject(message.links))
                return "links: object expected";
            let key = Object.keys(message.links);
            for (let i = 0; i < key.length; ++i) {
                let error = $root.Link.verify(message.links[key[i]]);
                if (error)
                    return "links." + error;
            }
        }
        return null;
    };

    /**
     * Creates a Node message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Node
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Node} Node
     */
    Node.fromObject = function fromObject(object) {
        if (object instanceof $root.Node)
            return object;
        let message = new $root.Node();
        if (object.name != null)
            message.name = String(object.name);
        if (object.pin != null)
            message.pin = Boolean(object.pin);
        if (object.plaintext != null)
            message.plaintext = Boolean(object.plaintext);
        if (object.mill != null)
            message.mill = String(object.mill);
        if (object.opts) {
            if (typeof object.opts !== "object")
                throw TypeError(".Node.opts: object expected");
            message.opts = {};
            for (let keys = Object.keys(object.opts), i = 0; i < keys.length; ++i)
                message.opts[keys[i]] = String(object.opts[keys[i]]);
        }
        if (object.jsonSchema != null) {
            if (typeof object.jsonSchema !== "object")
                throw TypeError(".Node.jsonSchema: object expected");
            message.jsonSchema = $root.google.protobuf.Struct.fromObject(object.jsonSchema);
        }
        if (object.links) {
            if (typeof object.links !== "object")
                throw TypeError(".Node.links: object expected");
            message.links = {};
            for (let keys = Object.keys(object.links), i = 0; i < keys.length; ++i) {
                if (typeof object.links[keys[i]] !== "object")
                    throw TypeError(".Node.links: object expected");
                message.links[keys[i]] = $root.Link.fromObject(object.links[keys[i]]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from a Node message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Node
     * @static
     * @param {Node} message Node
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Node.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.objects || options.defaults) {
            object.opts = {};
            object.links = {};
        }
        if (options.defaults) {
            object.name = "";
            object.pin = false;
            object.plaintext = false;
            object.mill = "";
            object.jsonSchema = null;
        }
        if (message.name != null && message.hasOwnProperty("name"))
            object.name = message.name;
        if (message.pin != null && message.hasOwnProperty("pin"))
            object.pin = message.pin;
        if (message.plaintext != null && message.hasOwnProperty("plaintext"))
            object.plaintext = message.plaintext;
        if (message.mill != null && message.hasOwnProperty("mill"))
            object.mill = message.mill;
        let keys2;
        if (message.opts && (keys2 = Object.keys(message.opts)).length) {
            object.opts = {};
            for (let j = 0; j < keys2.length; ++j)
                object.opts[keys2[j]] = message.opts[keys2[j]];
        }
        if (message.jsonSchema != null && message.hasOwnProperty("jsonSchema"))
            object.jsonSchema = $root.google.protobuf.Struct.toObject(message.jsonSchema, options);
        if (message.links && (keys2 = Object.keys(message.links)).length) {
            object.links = {};
            for (let j = 0; j < keys2.length; ++j)
                object.links[keys2[j]] = $root.Link.toObject(message.links[keys2[j]], options);
        }
        return object;
    };

    /**
     * Converts this Node to JSON.
     * @function toJSON
     * @memberof Node
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Node.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Node;
})();

export const Link = $root.Link = (() => {

    /**
     * Properties of a Link.
     * @exports ILink
     * @interface ILink
     * @property {string} use Link use
     * @property {boolean} pin Link pin
     * @property {boolean} plaintext Link plaintext
     * @property {string} mill Link mill
     * @property {Object.<string,string>} opts Link opts
     * @property {google.protobuf.IStruct} jsonSchema Link jsonSchema
     */

    /**
     * Constructs a new Link.
     * @exports Link
     * @classdesc Represents a Link.
     * @implements ILink
     * @constructor
     * @param {ILink=} [properties] Properties to set
     */
    function Link(properties) {
        this.opts = {};
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Link use.
     * @member {string} use
     * @memberof Link
     * @instance
     */
    Link.prototype.use = "";

    /**
     * Link pin.
     * @member {boolean} pin
     * @memberof Link
     * @instance
     */
    Link.prototype.pin = false;

    /**
     * Link plaintext.
     * @member {boolean} plaintext
     * @memberof Link
     * @instance
     */
    Link.prototype.plaintext = false;

    /**
     * Link mill.
     * @member {string} mill
     * @memberof Link
     * @instance
     */
    Link.prototype.mill = "";

    /**
     * Link opts.
     * @member {Object.<string,string>} opts
     * @memberof Link
     * @instance
     */
    Link.prototype.opts = $util.emptyObject;

    /**
     * Link jsonSchema.
     * @member {google.protobuf.IStruct} jsonSchema
     * @memberof Link
     * @instance
     */
    Link.prototype.jsonSchema = null;

    /**
     * Creates a new Link instance using the specified properties.
     * @function create
     * @memberof Link
     * @static
     * @param {ILink=} [properties] Properties to set
     * @returns {Link} Link instance
     */
    Link.create = function create(properties) {
        return new Link(properties);
    };

    /**
     * Encodes the specified Link message. Does not implicitly {@link Link.verify|verify} messages.
     * @function encode
     * @memberof Link
     * @static
     * @param {ILink} message Link message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Link.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.use != null && message.hasOwnProperty("use"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.use);
        if (message.pin != null && message.hasOwnProperty("pin"))
            writer.uint32(/* id 2, wireType 0 =*/16).bool(message.pin);
        if (message.plaintext != null && message.hasOwnProperty("plaintext"))
            writer.uint32(/* id 3, wireType 0 =*/24).bool(message.plaintext);
        if (message.mill != null && message.hasOwnProperty("mill"))
            writer.uint32(/* id 4, wireType 2 =*/34).string(message.mill);
        if (message.opts != null && message.hasOwnProperty("opts"))
            for (let keys = Object.keys(message.opts), i = 0; i < keys.length; ++i)
                writer.uint32(/* id 5, wireType 2 =*/42).fork().uint32(/* id 1, wireType 2 =*/10).string(keys[i]).uint32(/* id 2, wireType 2 =*/18).string(message.opts[keys[i]]).ldelim();
        if (message.jsonSchema != null && message.hasOwnProperty("jsonSchema"))
            $root.google.protobuf.Struct.encode(message.jsonSchema, writer.uint32(/* id 6, wireType 2 =*/50).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified Link message, length delimited. Does not implicitly {@link Link.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Link
     * @static
     * @param {ILink} message Link message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Link.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Link message from the specified reader or buffer.
     * @function decode
     * @memberof Link
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Link} Link
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Link.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Link(), key;
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.use = reader.string();
                break;
            case 2:
                message.pin = reader.bool();
                break;
            case 3:
                message.plaintext = reader.bool();
                break;
            case 4:
                message.mill = reader.string();
                break;
            case 5:
                reader.skip().pos++;
                if (message.opts === $util.emptyObject)
                    message.opts = {};
                key = reader.string();
                reader.pos++;
                message.opts[key] = reader.string();
                break;
            case 6:
                message.jsonSchema = $root.google.protobuf.Struct.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Link message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Link
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Link} Link
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Link.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Link message.
     * @function verify
     * @memberof Link
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Link.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.use != null && message.hasOwnProperty("use"))
            if (!$util.isString(message.use))
                return "use: string expected";
        if (message.pin != null && message.hasOwnProperty("pin"))
            if (typeof message.pin !== "boolean")
                return "pin: boolean expected";
        if (message.plaintext != null && message.hasOwnProperty("plaintext"))
            if (typeof message.plaintext !== "boolean")
                return "plaintext: boolean expected";
        if (message.mill != null && message.hasOwnProperty("mill"))
            if (!$util.isString(message.mill))
                return "mill: string expected";
        if (message.opts != null && message.hasOwnProperty("opts")) {
            if (!$util.isObject(message.opts))
                return "opts: object expected";
            let key = Object.keys(message.opts);
            for (let i = 0; i < key.length; ++i)
                if (!$util.isString(message.opts[key[i]]))
                    return "opts: string{k:string} expected";
        }
        if (message.jsonSchema != null && message.hasOwnProperty("jsonSchema")) {
            let error = $root.google.protobuf.Struct.verify(message.jsonSchema);
            if (error)
                return "jsonSchema." + error;
        }
        return null;
    };

    /**
     * Creates a Link message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Link
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Link} Link
     */
    Link.fromObject = function fromObject(object) {
        if (object instanceof $root.Link)
            return object;
        let message = new $root.Link();
        if (object.use != null)
            message.use = String(object.use);
        if (object.pin != null)
            message.pin = Boolean(object.pin);
        if (object.plaintext != null)
            message.plaintext = Boolean(object.plaintext);
        if (object.mill != null)
            message.mill = String(object.mill);
        if (object.opts) {
            if (typeof object.opts !== "object")
                throw TypeError(".Link.opts: object expected");
            message.opts = {};
            for (let keys = Object.keys(object.opts), i = 0; i < keys.length; ++i)
                message.opts[keys[i]] = String(object.opts[keys[i]]);
        }
        if (object.jsonSchema != null) {
            if (typeof object.jsonSchema !== "object")
                throw TypeError(".Link.jsonSchema: object expected");
            message.jsonSchema = $root.google.protobuf.Struct.fromObject(object.jsonSchema);
        }
        return message;
    };

    /**
     * Creates a plain object from a Link message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Link
     * @static
     * @param {Link} message Link
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Link.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.objects || options.defaults)
            object.opts = {};
        if (options.defaults) {
            object.use = "";
            object.pin = false;
            object.plaintext = false;
            object.mill = "";
            object.jsonSchema = null;
        }
        if (message.use != null && message.hasOwnProperty("use"))
            object.use = message.use;
        if (message.pin != null && message.hasOwnProperty("pin"))
            object.pin = message.pin;
        if (message.plaintext != null && message.hasOwnProperty("plaintext"))
            object.plaintext = message.plaintext;
        if (message.mill != null && message.hasOwnProperty("mill"))
            object.mill = message.mill;
        let keys2;
        if (message.opts && (keys2 = Object.keys(message.opts)).length) {
            object.opts = {};
            for (let j = 0; j < keys2.length; ++j)
                object.opts[keys2[j]] = message.opts[keys2[j]];
        }
        if (message.jsonSchema != null && message.hasOwnProperty("jsonSchema"))
            object.jsonSchema = $root.google.protobuf.Struct.toObject(message.jsonSchema, options);
        return object;
    };

    /**
     * Converts this Link to JSON.
     * @function toJSON
     * @memberof Link
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Link.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Link;
})();

export const Notification = $root.Notification = (() => {

    /**
     * Properties of a Notification.
     * @exports INotification
     * @interface INotification
     * @property {string} id Notification id
     * @property {google.protobuf.ITimestamp} date Notification date
     * @property {string} actor Notification actor
     * @property {string} subject Notification subject
     * @property {string} subjectDesc Notification subjectDesc
     * @property {string} block Notification block
     * @property {string} target Notification target
     * @property {Notification.Type} type Notification type
     * @property {string} body Notification body
     * @property {boolean} read Notification read
     */

    /**
     * Constructs a new Notification.
     * @exports Notification
     * @classdesc Represents a Notification.
     * @implements INotification
     * @constructor
     * @param {INotification=} [properties] Properties to set
     */
    function Notification(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Notification id.
     * @member {string} id
     * @memberof Notification
     * @instance
     */
    Notification.prototype.id = "";

    /**
     * Notification date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof Notification
     * @instance
     */
    Notification.prototype.date = null;

    /**
     * Notification actor.
     * @member {string} actor
     * @memberof Notification
     * @instance
     */
    Notification.prototype.actor = "";

    /**
     * Notification subject.
     * @member {string} subject
     * @memberof Notification
     * @instance
     */
    Notification.prototype.subject = "";

    /**
     * Notification subjectDesc.
     * @member {string} subjectDesc
     * @memberof Notification
     * @instance
     */
    Notification.prototype.subjectDesc = "";

    /**
     * Notification block.
     * @member {string} block
     * @memberof Notification
     * @instance
     */
    Notification.prototype.block = "";

    /**
     * Notification target.
     * @member {string} target
     * @memberof Notification
     * @instance
     */
    Notification.prototype.target = "";

    /**
     * Notification type.
     * @member {Notification.Type} type
     * @memberof Notification
     * @instance
     */
    Notification.prototype.type = 0;

    /**
     * Notification body.
     * @member {string} body
     * @memberof Notification
     * @instance
     */
    Notification.prototype.body = "";

    /**
     * Notification read.
     * @member {boolean} read
     * @memberof Notification
     * @instance
     */
    Notification.prototype.read = false;

    /**
     * Creates a new Notification instance using the specified properties.
     * @function create
     * @memberof Notification
     * @static
     * @param {INotification=} [properties] Properties to set
     * @returns {Notification} Notification instance
     */
    Notification.create = function create(properties) {
        return new Notification(properties);
    };

    /**
     * Encodes the specified Notification message. Does not implicitly {@link Notification.verify|verify} messages.
     * @function encode
     * @memberof Notification
     * @static
     * @param {INotification} message Notification message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Notification.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        if (message.actor != null && message.hasOwnProperty("actor"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.actor);
        if (message.subject != null && message.hasOwnProperty("subject"))
            writer.uint32(/* id 4, wireType 2 =*/34).string(message.subject);
        if (message.subjectDesc != null && message.hasOwnProperty("subjectDesc"))
            writer.uint32(/* id 5, wireType 2 =*/42).string(message.subjectDesc);
        if (message.block != null && message.hasOwnProperty("block"))
            writer.uint32(/* id 6, wireType 2 =*/50).string(message.block);
        if (message.target != null && message.hasOwnProperty("target"))
            writer.uint32(/* id 7, wireType 2 =*/58).string(message.target);
        if (message.type != null && message.hasOwnProperty("type"))
            writer.uint32(/* id 8, wireType 0 =*/64).int32(message.type);
        if (message.body != null && message.hasOwnProperty("body"))
            writer.uint32(/* id 9, wireType 2 =*/74).string(message.body);
        if (message.read != null && message.hasOwnProperty("read"))
            writer.uint32(/* id 10, wireType 0 =*/80).bool(message.read);
        return writer;
    };

    /**
     * Encodes the specified Notification message, length delimited. Does not implicitly {@link Notification.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Notification
     * @static
     * @param {INotification} message Notification message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Notification.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Notification message from the specified reader or buffer.
     * @function decode
     * @memberof Notification
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Notification} Notification
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Notification.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Notification();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 3:
                message.actor = reader.string();
                break;
            case 4:
                message.subject = reader.string();
                break;
            case 5:
                message.subjectDesc = reader.string();
                break;
            case 6:
                message.block = reader.string();
                break;
            case 7:
                message.target = reader.string();
                break;
            case 8:
                message.type = reader.int32();
                break;
            case 9:
                message.body = reader.string();
                break;
            case 10:
                message.read = reader.bool();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Notification message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Notification
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Notification} Notification
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Notification.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Notification message.
     * @function verify
     * @memberof Notification
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Notification.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        if (message.actor != null && message.hasOwnProperty("actor"))
            if (!$util.isString(message.actor))
                return "actor: string expected";
        if (message.subject != null && message.hasOwnProperty("subject"))
            if (!$util.isString(message.subject))
                return "subject: string expected";
        if (message.subjectDesc != null && message.hasOwnProperty("subjectDesc"))
            if (!$util.isString(message.subjectDesc))
                return "subjectDesc: string expected";
        if (message.block != null && message.hasOwnProperty("block"))
            if (!$util.isString(message.block))
                return "block: string expected";
        if (message.target != null && message.hasOwnProperty("target"))
            if (!$util.isString(message.target))
                return "target: string expected";
        if (message.type != null && message.hasOwnProperty("type"))
            switch (message.type) {
            default:
                return "type: enum value expected";
            case 0:
            case 1:
            case 2:
            case 3:
            case 4:
            case 5:
            case 6:
            case 7:
                break;
            }
        if (message.body != null && message.hasOwnProperty("body"))
            if (!$util.isString(message.body))
                return "body: string expected";
        if (message.read != null && message.hasOwnProperty("read"))
            if (typeof message.read !== "boolean")
                return "read: boolean expected";
        return null;
    };

    /**
     * Creates a Notification message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Notification
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Notification} Notification
     */
    Notification.fromObject = function fromObject(object) {
        if (object instanceof $root.Notification)
            return object;
        let message = new $root.Notification();
        if (object.id != null)
            message.id = String(object.id);
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".Notification.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        if (object.actor != null)
            message.actor = String(object.actor);
        if (object.subject != null)
            message.subject = String(object.subject);
        if (object.subjectDesc != null)
            message.subjectDesc = String(object.subjectDesc);
        if (object.block != null)
            message.block = String(object.block);
        if (object.target != null)
            message.target = String(object.target);
        switch (object.type) {
        case "INVITE_RECEIVED":
        case 0:
            message.type = 0;
            break;
        case "ACCOUNT_PEER_JOINED":
        case 1:
            message.type = 1;
            break;
        case "PEER_JOINED":
        case 2:
            message.type = 2;
            break;
        case "PEER_LEFT":
        case 3:
            message.type = 3;
            break;
        case "MESSAGE_ADDED":
        case 4:
            message.type = 4;
            break;
        case "FILES_ADDED":
        case 5:
            message.type = 5;
            break;
        case "COMMENT_ADDED":
        case 6:
            message.type = 6;
            break;
        case "LIKE_ADDED":
        case 7:
            message.type = 7;
            break;
        }
        if (object.body != null)
            message.body = String(object.body);
        if (object.read != null)
            message.read = Boolean(object.read);
        return message;
    };

    /**
     * Creates a plain object from a Notification message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Notification
     * @static
     * @param {Notification} message Notification
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Notification.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            object.date = null;
            object.actor = "";
            object.subject = "";
            object.subjectDesc = "";
            object.block = "";
            object.target = "";
            object.type = options.enums === String ? "INVITE_RECEIVED" : 0;
            object.body = "";
            object.read = false;
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        if (message.actor != null && message.hasOwnProperty("actor"))
            object.actor = message.actor;
        if (message.subject != null && message.hasOwnProperty("subject"))
            object.subject = message.subject;
        if (message.subjectDesc != null && message.hasOwnProperty("subjectDesc"))
            object.subjectDesc = message.subjectDesc;
        if (message.block != null && message.hasOwnProperty("block"))
            object.block = message.block;
        if (message.target != null && message.hasOwnProperty("target"))
            object.target = message.target;
        if (message.type != null && message.hasOwnProperty("type"))
            object.type = options.enums === String ? $root.Notification.Type[message.type] : message.type;
        if (message.body != null && message.hasOwnProperty("body"))
            object.body = message.body;
        if (message.read != null && message.hasOwnProperty("read"))
            object.read = message.read;
        return object;
    };

    /**
     * Converts this Notification to JSON.
     * @function toJSON
     * @memberof Notification
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Notification.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    /**
     * Type enum.
     * @name Notification.Type
     * @enum {string}
     * @property {number} INVITE_RECEIVED=0 INVITE_RECEIVED value
     * @property {number} ACCOUNT_PEER_JOINED=1 ACCOUNT_PEER_JOINED value
     * @property {number} PEER_JOINED=2 PEER_JOINED value
     * @property {number} PEER_LEFT=3 PEER_LEFT value
     * @property {number} MESSAGE_ADDED=4 MESSAGE_ADDED value
     * @property {number} FILES_ADDED=5 FILES_ADDED value
     * @property {number} COMMENT_ADDED=6 COMMENT_ADDED value
     * @property {number} LIKE_ADDED=7 LIKE_ADDED value
     */
    Notification.Type = (function() {
        const valuesById = {}, values = Object.create(valuesById);
        values[valuesById[0] = "INVITE_RECEIVED"] = 0;
        values[valuesById[1] = "ACCOUNT_PEER_JOINED"] = 1;
        values[valuesById[2] = "PEER_JOINED"] = 2;
        values[valuesById[3] = "PEER_LEFT"] = 3;
        values[valuesById[4] = "MESSAGE_ADDED"] = 4;
        values[valuesById[5] = "FILES_ADDED"] = 5;
        values[valuesById[6] = "COMMENT_ADDED"] = 6;
        values[valuesById[7] = "LIKE_ADDED"] = 7;
        return values;
    })();

    return Notification;
})();

export const Cafe = $root.Cafe = (() => {

    /**
     * Properties of a Cafe.
     * @exports ICafe
     * @interface ICafe
     * @property {string} peer Cafe peer
     * @property {string} address Cafe address
     * @property {string} api Cafe api
     * @property {string} protocol Cafe protocol
     * @property {string} node Cafe node
     * @property {string} url Cafe url
     * @property {Array.<string>} swarm Cafe swarm
     */

    /**
     * Constructs a new Cafe.
     * @exports Cafe
     * @classdesc Represents a Cafe.
     * @implements ICafe
     * @constructor
     * @param {ICafe=} [properties] Properties to set
     */
    function Cafe(properties) {
        this.swarm = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Cafe peer.
     * @member {string} peer
     * @memberof Cafe
     * @instance
     */
    Cafe.prototype.peer = "";

    /**
     * Cafe address.
     * @member {string} address
     * @memberof Cafe
     * @instance
     */
    Cafe.prototype.address = "";

    /**
     * Cafe api.
     * @member {string} api
     * @memberof Cafe
     * @instance
     */
    Cafe.prototype.api = "";

    /**
     * Cafe protocol.
     * @member {string} protocol
     * @memberof Cafe
     * @instance
     */
    Cafe.prototype.protocol = "";

    /**
     * Cafe node.
     * @member {string} node
     * @memberof Cafe
     * @instance
     */
    Cafe.prototype.node = "";

    /**
     * Cafe url.
     * @member {string} url
     * @memberof Cafe
     * @instance
     */
    Cafe.prototype.url = "";

    /**
     * Cafe swarm.
     * @member {Array.<string>} swarm
     * @memberof Cafe
     * @instance
     */
    Cafe.prototype.swarm = $util.emptyArray;

    /**
     * Creates a new Cafe instance using the specified properties.
     * @function create
     * @memberof Cafe
     * @static
     * @param {ICafe=} [properties] Properties to set
     * @returns {Cafe} Cafe instance
     */
    Cafe.create = function create(properties) {
        return new Cafe(properties);
    };

    /**
     * Encodes the specified Cafe message. Does not implicitly {@link Cafe.verify|verify} messages.
     * @function encode
     * @memberof Cafe
     * @static
     * @param {ICafe} message Cafe message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Cafe.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.peer != null && message.hasOwnProperty("peer"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.peer);
        if (message.address != null && message.hasOwnProperty("address"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.address);
        if (message.api != null && message.hasOwnProperty("api"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.api);
        if (message.protocol != null && message.hasOwnProperty("protocol"))
            writer.uint32(/* id 4, wireType 2 =*/34).string(message.protocol);
        if (message.node != null && message.hasOwnProperty("node"))
            writer.uint32(/* id 5, wireType 2 =*/42).string(message.node);
        if (message.url != null && message.hasOwnProperty("url"))
            writer.uint32(/* id 6, wireType 2 =*/50).string(message.url);
        if (message.swarm != null && message.swarm.length)
            for (let i = 0; i < message.swarm.length; ++i)
                writer.uint32(/* id 7, wireType 2 =*/58).string(message.swarm[i]);
        return writer;
    };

    /**
     * Encodes the specified Cafe message, length delimited. Does not implicitly {@link Cafe.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Cafe
     * @static
     * @param {ICafe} message Cafe message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Cafe.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Cafe message from the specified reader or buffer.
     * @function decode
     * @memberof Cafe
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Cafe} Cafe
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Cafe.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Cafe();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.peer = reader.string();
                break;
            case 2:
                message.address = reader.string();
                break;
            case 3:
                message.api = reader.string();
                break;
            case 4:
                message.protocol = reader.string();
                break;
            case 5:
                message.node = reader.string();
                break;
            case 6:
                message.url = reader.string();
                break;
            case 7:
                if (!(message.swarm && message.swarm.length))
                    message.swarm = [];
                message.swarm.push(reader.string());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Cafe message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Cafe
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Cafe} Cafe
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Cafe.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Cafe message.
     * @function verify
     * @memberof Cafe
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Cafe.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.peer != null && message.hasOwnProperty("peer"))
            if (!$util.isString(message.peer))
                return "peer: string expected";
        if (message.address != null && message.hasOwnProperty("address"))
            if (!$util.isString(message.address))
                return "address: string expected";
        if (message.api != null && message.hasOwnProperty("api"))
            if (!$util.isString(message.api))
                return "api: string expected";
        if (message.protocol != null && message.hasOwnProperty("protocol"))
            if (!$util.isString(message.protocol))
                return "protocol: string expected";
        if (message.node != null && message.hasOwnProperty("node"))
            if (!$util.isString(message.node))
                return "node: string expected";
        if (message.url != null && message.hasOwnProperty("url"))
            if (!$util.isString(message.url))
                return "url: string expected";
        if (message.swarm != null && message.hasOwnProperty("swarm")) {
            if (!Array.isArray(message.swarm))
                return "swarm: array expected";
            for (let i = 0; i < message.swarm.length; ++i)
                if (!$util.isString(message.swarm[i]))
                    return "swarm: string[] expected";
        }
        return null;
    };

    /**
     * Creates a Cafe message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Cafe
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Cafe} Cafe
     */
    Cafe.fromObject = function fromObject(object) {
        if (object instanceof $root.Cafe)
            return object;
        let message = new $root.Cafe();
        if (object.peer != null)
            message.peer = String(object.peer);
        if (object.address != null)
            message.address = String(object.address);
        if (object.api != null)
            message.api = String(object.api);
        if (object.protocol != null)
            message.protocol = String(object.protocol);
        if (object.node != null)
            message.node = String(object.node);
        if (object.url != null)
            message.url = String(object.url);
        if (object.swarm) {
            if (!Array.isArray(object.swarm))
                throw TypeError(".Cafe.swarm: array expected");
            message.swarm = [];
            for (let i = 0; i < object.swarm.length; ++i)
                message.swarm[i] = String(object.swarm[i]);
        }
        return message;
    };

    /**
     * Creates a plain object from a Cafe message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Cafe
     * @static
     * @param {Cafe} message Cafe
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Cafe.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.swarm = [];
        if (options.defaults) {
            object.peer = "";
            object.address = "";
            object.api = "";
            object.protocol = "";
            object.node = "";
            object.url = "";
        }
        if (message.peer != null && message.hasOwnProperty("peer"))
            object.peer = message.peer;
        if (message.address != null && message.hasOwnProperty("address"))
            object.address = message.address;
        if (message.api != null && message.hasOwnProperty("api"))
            object.api = message.api;
        if (message.protocol != null && message.hasOwnProperty("protocol"))
            object.protocol = message.protocol;
        if (message.node != null && message.hasOwnProperty("node"))
            object.node = message.node;
        if (message.url != null && message.hasOwnProperty("url"))
            object.url = message.url;
        if (message.swarm && message.swarm.length) {
            object.swarm = [];
            for (let j = 0; j < message.swarm.length; ++j)
                object.swarm[j] = message.swarm[j];
        }
        return object;
    };

    /**
     * Converts this Cafe to JSON.
     * @function toJSON
     * @memberof Cafe
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Cafe.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Cafe;
})();

export const CafeSession = $root.CafeSession = (() => {

    /**
     * Properties of a CafeSession.
     * @exports ICafeSession
     * @interface ICafeSession
     * @property {string} id CafeSession id
     * @property {string} access CafeSession access
     * @property {google.protobuf.ITimestamp} exp CafeSession exp
     * @property {string} refresh CafeSession refresh
     * @property {google.protobuf.ITimestamp} rexp CafeSession rexp
     * @property {string} subject CafeSession subject
     * @property {string} type CafeSession type
     * @property {ICafe} cafe CafeSession cafe
     */

    /**
     * Constructs a new CafeSession.
     * @exports CafeSession
     * @classdesc Represents a CafeSession.
     * @implements ICafeSession
     * @constructor
     * @param {ICafeSession=} [properties] Properties to set
     */
    function CafeSession(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeSession id.
     * @member {string} id
     * @memberof CafeSession
     * @instance
     */
    CafeSession.prototype.id = "";

    /**
     * CafeSession access.
     * @member {string} access
     * @memberof CafeSession
     * @instance
     */
    CafeSession.prototype.access = "";

    /**
     * CafeSession exp.
     * @member {google.protobuf.ITimestamp} exp
     * @memberof CafeSession
     * @instance
     */
    CafeSession.prototype.exp = null;

    /**
     * CafeSession refresh.
     * @member {string} refresh
     * @memberof CafeSession
     * @instance
     */
    CafeSession.prototype.refresh = "";

    /**
     * CafeSession rexp.
     * @member {google.protobuf.ITimestamp} rexp
     * @memberof CafeSession
     * @instance
     */
    CafeSession.prototype.rexp = null;

    /**
     * CafeSession subject.
     * @member {string} subject
     * @memberof CafeSession
     * @instance
     */
    CafeSession.prototype.subject = "";

    /**
     * CafeSession type.
     * @member {string} type
     * @memberof CafeSession
     * @instance
     */
    CafeSession.prototype.type = "";

    /**
     * CafeSession cafe.
     * @member {ICafe} cafe
     * @memberof CafeSession
     * @instance
     */
    CafeSession.prototype.cafe = null;

    /**
     * Creates a new CafeSession instance using the specified properties.
     * @function create
     * @memberof CafeSession
     * @static
     * @param {ICafeSession=} [properties] Properties to set
     * @returns {CafeSession} CafeSession instance
     */
    CafeSession.create = function create(properties) {
        return new CafeSession(properties);
    };

    /**
     * Encodes the specified CafeSession message. Does not implicitly {@link CafeSession.verify|verify} messages.
     * @function encode
     * @memberof CafeSession
     * @static
     * @param {ICafeSession} message CafeSession message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeSession.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.access != null && message.hasOwnProperty("access"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.access);
        if (message.exp != null && message.hasOwnProperty("exp"))
            $root.google.protobuf.Timestamp.encode(message.exp, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        if (message.refresh != null && message.hasOwnProperty("refresh"))
            writer.uint32(/* id 4, wireType 2 =*/34).string(message.refresh);
        if (message.rexp != null && message.hasOwnProperty("rexp"))
            $root.google.protobuf.Timestamp.encode(message.rexp, writer.uint32(/* id 5, wireType 2 =*/42).fork()).ldelim();
        if (message.subject != null && message.hasOwnProperty("subject"))
            writer.uint32(/* id 6, wireType 2 =*/50).string(message.subject);
        if (message.type != null && message.hasOwnProperty("type"))
            writer.uint32(/* id 7, wireType 2 =*/58).string(message.type);
        if (message.cafe != null && message.hasOwnProperty("cafe"))
            $root.Cafe.encode(message.cafe, writer.uint32(/* id 8, wireType 2 =*/66).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified CafeSession message, length delimited. Does not implicitly {@link CafeSession.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeSession
     * @static
     * @param {ICafeSession} message CafeSession message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeSession.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeSession message from the specified reader or buffer.
     * @function decode
     * @memberof CafeSession
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeSession} CafeSession
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeSession.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeSession();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.access = reader.string();
                break;
            case 3:
                message.exp = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 4:
                message.refresh = reader.string();
                break;
            case 5:
                message.rexp = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 6:
                message.subject = reader.string();
                break;
            case 7:
                message.type = reader.string();
                break;
            case 8:
                message.cafe = $root.Cafe.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeSession message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeSession
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeSession} CafeSession
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeSession.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeSession message.
     * @function verify
     * @memberof CafeSession
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeSession.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.access != null && message.hasOwnProperty("access"))
            if (!$util.isString(message.access))
                return "access: string expected";
        if (message.exp != null && message.hasOwnProperty("exp")) {
            let error = $root.google.protobuf.Timestamp.verify(message.exp);
            if (error)
                return "exp." + error;
        }
        if (message.refresh != null && message.hasOwnProperty("refresh"))
            if (!$util.isString(message.refresh))
                return "refresh: string expected";
        if (message.rexp != null && message.hasOwnProperty("rexp")) {
            let error = $root.google.protobuf.Timestamp.verify(message.rexp);
            if (error)
                return "rexp." + error;
        }
        if (message.subject != null && message.hasOwnProperty("subject"))
            if (!$util.isString(message.subject))
                return "subject: string expected";
        if (message.type != null && message.hasOwnProperty("type"))
            if (!$util.isString(message.type))
                return "type: string expected";
        if (message.cafe != null && message.hasOwnProperty("cafe")) {
            let error = $root.Cafe.verify(message.cafe);
            if (error)
                return "cafe." + error;
        }
        return null;
    };

    /**
     * Creates a CafeSession message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeSession
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeSession} CafeSession
     */
    CafeSession.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeSession)
            return object;
        let message = new $root.CafeSession();
        if (object.id != null)
            message.id = String(object.id);
        if (object.access != null)
            message.access = String(object.access);
        if (object.exp != null) {
            if (typeof object.exp !== "object")
                throw TypeError(".CafeSession.exp: object expected");
            message.exp = $root.google.protobuf.Timestamp.fromObject(object.exp);
        }
        if (object.refresh != null)
            message.refresh = String(object.refresh);
        if (object.rexp != null) {
            if (typeof object.rexp !== "object")
                throw TypeError(".CafeSession.rexp: object expected");
            message.rexp = $root.google.protobuf.Timestamp.fromObject(object.rexp);
        }
        if (object.subject != null)
            message.subject = String(object.subject);
        if (object.type != null)
            message.type = String(object.type);
        if (object.cafe != null) {
            if (typeof object.cafe !== "object")
                throw TypeError(".CafeSession.cafe: object expected");
            message.cafe = $root.Cafe.fromObject(object.cafe);
        }
        return message;
    };

    /**
     * Creates a plain object from a CafeSession message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeSession
     * @static
     * @param {CafeSession} message CafeSession
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeSession.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            object.access = "";
            object.exp = null;
            object.refresh = "";
            object.rexp = null;
            object.subject = "";
            object.type = "";
            object.cafe = null;
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.access != null && message.hasOwnProperty("access"))
            object.access = message.access;
        if (message.exp != null && message.hasOwnProperty("exp"))
            object.exp = $root.google.protobuf.Timestamp.toObject(message.exp, options);
        if (message.refresh != null && message.hasOwnProperty("refresh"))
            object.refresh = message.refresh;
        if (message.rexp != null && message.hasOwnProperty("rexp"))
            object.rexp = $root.google.protobuf.Timestamp.toObject(message.rexp, options);
        if (message.subject != null && message.hasOwnProperty("subject"))
            object.subject = message.subject;
        if (message.type != null && message.hasOwnProperty("type"))
            object.type = message.type;
        if (message.cafe != null && message.hasOwnProperty("cafe"))
            object.cafe = $root.Cafe.toObject(message.cafe, options);
        return object;
    };

    /**
     * Converts this CafeSession to JSON.
     * @function toJSON
     * @memberof CafeSession
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeSession.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeSession;
})();

export const CafeSessionList = $root.CafeSessionList = (() => {

    /**
     * Properties of a CafeSessionList.
     * @exports ICafeSessionList
     * @interface ICafeSessionList
     * @property {Array.<ICafeSession>} items CafeSessionList items
     */

    /**
     * Constructs a new CafeSessionList.
     * @exports CafeSessionList
     * @classdesc Represents a CafeSessionList.
     * @implements ICafeSessionList
     * @constructor
     * @param {ICafeSessionList=} [properties] Properties to set
     */
    function CafeSessionList(properties) {
        this.items = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeSessionList items.
     * @member {Array.<ICafeSession>} items
     * @memberof CafeSessionList
     * @instance
     */
    CafeSessionList.prototype.items = $util.emptyArray;

    /**
     * Creates a new CafeSessionList instance using the specified properties.
     * @function create
     * @memberof CafeSessionList
     * @static
     * @param {ICafeSessionList=} [properties] Properties to set
     * @returns {CafeSessionList} CafeSessionList instance
     */
    CafeSessionList.create = function create(properties) {
        return new CafeSessionList(properties);
    };

    /**
     * Encodes the specified CafeSessionList message. Does not implicitly {@link CafeSessionList.verify|verify} messages.
     * @function encode
     * @memberof CafeSessionList
     * @static
     * @param {ICafeSessionList} message CafeSessionList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeSessionList.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.items != null && message.items.length)
            for (let i = 0; i < message.items.length; ++i)
                $root.CafeSession.encode(message.items[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified CafeSessionList message, length delimited. Does not implicitly {@link CafeSessionList.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeSessionList
     * @static
     * @param {ICafeSessionList} message CafeSessionList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeSessionList.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeSessionList message from the specified reader or buffer.
     * @function decode
     * @memberof CafeSessionList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeSessionList} CafeSessionList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeSessionList.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeSessionList();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                if (!(message.items && message.items.length))
                    message.items = [];
                message.items.push($root.CafeSession.decode(reader, reader.uint32()));
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeSessionList message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeSessionList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeSessionList} CafeSessionList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeSessionList.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeSessionList message.
     * @function verify
     * @memberof CafeSessionList
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeSessionList.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.items != null && message.hasOwnProperty("items")) {
            if (!Array.isArray(message.items))
                return "items: array expected";
            for (let i = 0; i < message.items.length; ++i) {
                let error = $root.CafeSession.verify(message.items[i]);
                if (error)
                    return "items." + error;
            }
        }
        return null;
    };

    /**
     * Creates a CafeSessionList message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeSessionList
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeSessionList} CafeSessionList
     */
    CafeSessionList.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeSessionList)
            return object;
        let message = new $root.CafeSessionList();
        if (object.items) {
            if (!Array.isArray(object.items))
                throw TypeError(".CafeSessionList.items: array expected");
            message.items = [];
            for (let i = 0; i < object.items.length; ++i) {
                if (typeof object.items[i] !== "object")
                    throw TypeError(".CafeSessionList.items: object expected");
                message.items[i] = $root.CafeSession.fromObject(object.items[i]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from a CafeSessionList message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeSessionList
     * @static
     * @param {CafeSessionList} message CafeSessionList
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeSessionList.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.items = [];
        if (message.items && message.items.length) {
            object.items = [];
            for (let j = 0; j < message.items.length; ++j)
                object.items[j] = $root.CafeSession.toObject(message.items[j], options);
        }
        return object;
    };

    /**
     * Converts this CafeSessionList to JSON.
     * @function toJSON
     * @memberof CafeSessionList
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeSessionList.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeSessionList;
})();

export const CafeRequest = $root.CafeRequest = (() => {

    /**
     * Properties of a CafeRequest.
     * @exports ICafeRequest
     * @interface ICafeRequest
     * @property {string} id CafeRequest id
     * @property {string} peer CafeRequest peer
     * @property {string} target CafeRequest target
     * @property {ICafe} cafe CafeRequest cafe
     * @property {CafeRequest.Type} type CafeRequest type
     * @property {google.protobuf.ITimestamp} date CafeRequest date
     */

    /**
     * Constructs a new CafeRequest.
     * @exports CafeRequest
     * @classdesc Represents a CafeRequest.
     * @implements ICafeRequest
     * @constructor
     * @param {ICafeRequest=} [properties] Properties to set
     */
    function CafeRequest(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeRequest id.
     * @member {string} id
     * @memberof CafeRequest
     * @instance
     */
    CafeRequest.prototype.id = "";

    /**
     * CafeRequest peer.
     * @member {string} peer
     * @memberof CafeRequest
     * @instance
     */
    CafeRequest.prototype.peer = "";

    /**
     * CafeRequest target.
     * @member {string} target
     * @memberof CafeRequest
     * @instance
     */
    CafeRequest.prototype.target = "";

    /**
     * CafeRequest cafe.
     * @member {ICafe} cafe
     * @memberof CafeRequest
     * @instance
     */
    CafeRequest.prototype.cafe = null;

    /**
     * CafeRequest type.
     * @member {CafeRequest.Type} type
     * @memberof CafeRequest
     * @instance
     */
    CafeRequest.prototype.type = 0;

    /**
     * CafeRequest date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof CafeRequest
     * @instance
     */
    CafeRequest.prototype.date = null;

    /**
     * Creates a new CafeRequest instance using the specified properties.
     * @function create
     * @memberof CafeRequest
     * @static
     * @param {ICafeRequest=} [properties] Properties to set
     * @returns {CafeRequest} CafeRequest instance
     */
    CafeRequest.create = function create(properties) {
        return new CafeRequest(properties);
    };

    /**
     * Encodes the specified CafeRequest message. Does not implicitly {@link CafeRequest.verify|verify} messages.
     * @function encode
     * @memberof CafeRequest
     * @static
     * @param {ICafeRequest} message CafeRequest message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeRequest.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.peer != null && message.hasOwnProperty("peer"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.peer);
        if (message.target != null && message.hasOwnProperty("target"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.target);
        if (message.cafe != null && message.hasOwnProperty("cafe"))
            $root.Cafe.encode(message.cafe, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
        if (message.type != null && message.hasOwnProperty("type"))
            writer.uint32(/* id 5, wireType 0 =*/40).int32(message.type);
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 6, wireType 2 =*/50).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified CafeRequest message, length delimited. Does not implicitly {@link CafeRequest.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeRequest
     * @static
     * @param {ICafeRequest} message CafeRequest message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeRequest.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeRequest message from the specified reader or buffer.
     * @function decode
     * @memberof CafeRequest
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeRequest} CafeRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeRequest.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeRequest();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.peer = reader.string();
                break;
            case 3:
                message.target = reader.string();
                break;
            case 4:
                message.cafe = $root.Cafe.decode(reader, reader.uint32());
                break;
            case 5:
                message.type = reader.int32();
                break;
            case 6:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeRequest message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeRequest
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeRequest} CafeRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeRequest.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeRequest message.
     * @function verify
     * @memberof CafeRequest
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeRequest.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.peer != null && message.hasOwnProperty("peer"))
            if (!$util.isString(message.peer))
                return "peer: string expected";
        if (message.target != null && message.hasOwnProperty("target"))
            if (!$util.isString(message.target))
                return "target: string expected";
        if (message.cafe != null && message.hasOwnProperty("cafe")) {
            let error = $root.Cafe.verify(message.cafe);
            if (error)
                return "cafe." + error;
        }
        if (message.type != null && message.hasOwnProperty("type"))
            switch (message.type) {
            default:
                return "type: enum value expected";
            case 0:
            case 1:
            case 2:
                break;
            }
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        return null;
    };

    /**
     * Creates a CafeRequest message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeRequest
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeRequest} CafeRequest
     */
    CafeRequest.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeRequest)
            return object;
        let message = new $root.CafeRequest();
        if (object.id != null)
            message.id = String(object.id);
        if (object.peer != null)
            message.peer = String(object.peer);
        if (object.target != null)
            message.target = String(object.target);
        if (object.cafe != null) {
            if (typeof object.cafe !== "object")
                throw TypeError(".CafeRequest.cafe: object expected");
            message.cafe = $root.Cafe.fromObject(object.cafe);
        }
        switch (object.type) {
        case "STORE":
        case 0:
            message.type = 0;
            break;
        case "STORE_THREAD":
        case 1:
            message.type = 1;
            break;
        case "INBOX":
        case 2:
            message.type = 2;
            break;
        }
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".CafeRequest.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        return message;
    };

    /**
     * Creates a plain object from a CafeRequest message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeRequest
     * @static
     * @param {CafeRequest} message CafeRequest
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeRequest.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            object.peer = "";
            object.target = "";
            object.cafe = null;
            object.type = options.enums === String ? "STORE" : 0;
            object.date = null;
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.peer != null && message.hasOwnProperty("peer"))
            object.peer = message.peer;
        if (message.target != null && message.hasOwnProperty("target"))
            object.target = message.target;
        if (message.cafe != null && message.hasOwnProperty("cafe"))
            object.cafe = $root.Cafe.toObject(message.cafe, options);
        if (message.type != null && message.hasOwnProperty("type"))
            object.type = options.enums === String ? $root.CafeRequest.Type[message.type] : message.type;
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        return object;
    };

    /**
     * Converts this CafeRequest to JSON.
     * @function toJSON
     * @memberof CafeRequest
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeRequest.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    /**
     * Type enum.
     * @name CafeRequest.Type
     * @enum {string}
     * @property {number} STORE=0 STORE value
     * @property {number} STORE_THREAD=1 STORE_THREAD value
     * @property {number} INBOX=2 INBOX value
     */
    CafeRequest.Type = (function() {
        const valuesById = {}, values = Object.create(valuesById);
        values[valuesById[0] = "STORE"] = 0;
        values[valuesById[1] = "STORE_THREAD"] = 1;
        values[valuesById[2] = "INBOX"] = 2;
        return values;
    })();

    return CafeRequest;
})();

export const CafeMessage = $root.CafeMessage = (() => {

    /**
     * Properties of a CafeMessage.
     * @exports ICafeMessage
     * @interface ICafeMessage
     * @property {string} id CafeMessage id
     * @property {string} peer CafeMessage peer
     * @property {google.protobuf.ITimestamp} date CafeMessage date
     * @property {number} attempts CafeMessage attempts
     */

    /**
     * Constructs a new CafeMessage.
     * @exports CafeMessage
     * @classdesc Represents a CafeMessage.
     * @implements ICafeMessage
     * @constructor
     * @param {ICafeMessage=} [properties] Properties to set
     */
    function CafeMessage(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeMessage id.
     * @member {string} id
     * @memberof CafeMessage
     * @instance
     */
    CafeMessage.prototype.id = "";

    /**
     * CafeMessage peer.
     * @member {string} peer
     * @memberof CafeMessage
     * @instance
     */
    CafeMessage.prototype.peer = "";

    /**
     * CafeMessage date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof CafeMessage
     * @instance
     */
    CafeMessage.prototype.date = null;

    /**
     * CafeMessage attempts.
     * @member {number} attempts
     * @memberof CafeMessage
     * @instance
     */
    CafeMessage.prototype.attempts = 0;

    /**
     * Creates a new CafeMessage instance using the specified properties.
     * @function create
     * @memberof CafeMessage
     * @static
     * @param {ICafeMessage=} [properties] Properties to set
     * @returns {CafeMessage} CafeMessage instance
     */
    CafeMessage.create = function create(properties) {
        return new CafeMessage(properties);
    };

    /**
     * Encodes the specified CafeMessage message. Does not implicitly {@link CafeMessage.verify|verify} messages.
     * @function encode
     * @memberof CafeMessage
     * @static
     * @param {ICafeMessage} message CafeMessage message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeMessage.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.peer != null && message.hasOwnProperty("peer"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.peer);
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        if (message.attempts != null && message.hasOwnProperty("attempts"))
            writer.uint32(/* id 4, wireType 0 =*/32).int32(message.attempts);
        return writer;
    };

    /**
     * Encodes the specified CafeMessage message, length delimited. Does not implicitly {@link CafeMessage.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeMessage
     * @static
     * @param {ICafeMessage} message CafeMessage message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeMessage.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeMessage message from the specified reader or buffer.
     * @function decode
     * @memberof CafeMessage
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeMessage} CafeMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeMessage.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeMessage();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.peer = reader.string();
                break;
            case 3:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 4:
                message.attempts = reader.int32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeMessage message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeMessage
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeMessage} CafeMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeMessage.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeMessage message.
     * @function verify
     * @memberof CafeMessage
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeMessage.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.peer != null && message.hasOwnProperty("peer"))
            if (!$util.isString(message.peer))
                return "peer: string expected";
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        if (message.attempts != null && message.hasOwnProperty("attempts"))
            if (!$util.isInteger(message.attempts))
                return "attempts: integer expected";
        return null;
    };

    /**
     * Creates a CafeMessage message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeMessage
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeMessage} CafeMessage
     */
    CafeMessage.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeMessage)
            return object;
        let message = new $root.CafeMessage();
        if (object.id != null)
            message.id = String(object.id);
        if (object.peer != null)
            message.peer = String(object.peer);
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".CafeMessage.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        if (object.attempts != null)
            message.attempts = object.attempts | 0;
        return message;
    };

    /**
     * Creates a plain object from a CafeMessage message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeMessage
     * @static
     * @param {CafeMessage} message CafeMessage
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeMessage.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            object.peer = "";
            object.date = null;
            object.attempts = 0;
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.peer != null && message.hasOwnProperty("peer"))
            object.peer = message.peer;
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        if (message.attempts != null && message.hasOwnProperty("attempts"))
            object.attempts = message.attempts;
        return object;
    };

    /**
     * Converts this CafeMessage to JSON.
     * @function toJSON
     * @memberof CafeMessage
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeMessage.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeMessage;
})();

export const CafeClientNonce = $root.CafeClientNonce = (() => {

    /**
     * Properties of a CafeClientNonce.
     * @exports ICafeClientNonce
     * @interface ICafeClientNonce
     * @property {string} value CafeClientNonce value
     * @property {string} address CafeClientNonce address
     * @property {google.protobuf.ITimestamp} date CafeClientNonce date
     */

    /**
     * Constructs a new CafeClientNonce.
     * @exports CafeClientNonce
     * @classdesc Represents a CafeClientNonce.
     * @implements ICafeClientNonce
     * @constructor
     * @param {ICafeClientNonce=} [properties] Properties to set
     */
    function CafeClientNonce(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeClientNonce value.
     * @member {string} value
     * @memberof CafeClientNonce
     * @instance
     */
    CafeClientNonce.prototype.value = "";

    /**
     * CafeClientNonce address.
     * @member {string} address
     * @memberof CafeClientNonce
     * @instance
     */
    CafeClientNonce.prototype.address = "";

    /**
     * CafeClientNonce date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof CafeClientNonce
     * @instance
     */
    CafeClientNonce.prototype.date = null;

    /**
     * Creates a new CafeClientNonce instance using the specified properties.
     * @function create
     * @memberof CafeClientNonce
     * @static
     * @param {ICafeClientNonce=} [properties] Properties to set
     * @returns {CafeClientNonce} CafeClientNonce instance
     */
    CafeClientNonce.create = function create(properties) {
        return new CafeClientNonce(properties);
    };

    /**
     * Encodes the specified CafeClientNonce message. Does not implicitly {@link CafeClientNonce.verify|verify} messages.
     * @function encode
     * @memberof CafeClientNonce
     * @static
     * @param {ICafeClientNonce} message CafeClientNonce message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeClientNonce.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.value != null && message.hasOwnProperty("value"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.value);
        if (message.address != null && message.hasOwnProperty("address"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.address);
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified CafeClientNonce message, length delimited. Does not implicitly {@link CafeClientNonce.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeClientNonce
     * @static
     * @param {ICafeClientNonce} message CafeClientNonce message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeClientNonce.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeClientNonce message from the specified reader or buffer.
     * @function decode
     * @memberof CafeClientNonce
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeClientNonce} CafeClientNonce
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeClientNonce.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeClientNonce();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.value = reader.string();
                break;
            case 2:
                message.address = reader.string();
                break;
            case 3:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeClientNonce message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeClientNonce
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeClientNonce} CafeClientNonce
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeClientNonce.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeClientNonce message.
     * @function verify
     * @memberof CafeClientNonce
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeClientNonce.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.value != null && message.hasOwnProperty("value"))
            if (!$util.isString(message.value))
                return "value: string expected";
        if (message.address != null && message.hasOwnProperty("address"))
            if (!$util.isString(message.address))
                return "address: string expected";
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        return null;
    };

    /**
     * Creates a CafeClientNonce message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeClientNonce
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeClientNonce} CafeClientNonce
     */
    CafeClientNonce.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeClientNonce)
            return object;
        let message = new $root.CafeClientNonce();
        if (object.value != null)
            message.value = String(object.value);
        if (object.address != null)
            message.address = String(object.address);
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".CafeClientNonce.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        return message;
    };

    /**
     * Creates a plain object from a CafeClientNonce message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeClientNonce
     * @static
     * @param {CafeClientNonce} message CafeClientNonce
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeClientNonce.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.value = "";
            object.address = "";
            object.date = null;
        }
        if (message.value != null && message.hasOwnProperty("value"))
            object.value = message.value;
        if (message.address != null && message.hasOwnProperty("address"))
            object.address = message.address;
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        return object;
    };

    /**
     * Converts this CafeClientNonce to JSON.
     * @function toJSON
     * @memberof CafeClientNonce
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeClientNonce.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeClientNonce;
})();

export const CafeClient = $root.CafeClient = (() => {

    /**
     * Properties of a CafeClient.
     * @exports ICafeClient
     * @interface ICafeClient
     * @property {string} id CafeClient id
     * @property {string} address CafeClient address
     * @property {google.protobuf.ITimestamp} created CafeClient created
     * @property {google.protobuf.ITimestamp} seen CafeClient seen
     * @property {string} token CafeClient token
     */

    /**
     * Constructs a new CafeClient.
     * @exports CafeClient
     * @classdesc Represents a CafeClient.
     * @implements ICafeClient
     * @constructor
     * @param {ICafeClient=} [properties] Properties to set
     */
    function CafeClient(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeClient id.
     * @member {string} id
     * @memberof CafeClient
     * @instance
     */
    CafeClient.prototype.id = "";

    /**
     * CafeClient address.
     * @member {string} address
     * @memberof CafeClient
     * @instance
     */
    CafeClient.prototype.address = "";

    /**
     * CafeClient created.
     * @member {google.protobuf.ITimestamp} created
     * @memberof CafeClient
     * @instance
     */
    CafeClient.prototype.created = null;

    /**
     * CafeClient seen.
     * @member {google.protobuf.ITimestamp} seen
     * @memberof CafeClient
     * @instance
     */
    CafeClient.prototype.seen = null;

    /**
     * CafeClient token.
     * @member {string} token
     * @memberof CafeClient
     * @instance
     */
    CafeClient.prototype.token = "";

    /**
     * Creates a new CafeClient instance using the specified properties.
     * @function create
     * @memberof CafeClient
     * @static
     * @param {ICafeClient=} [properties] Properties to set
     * @returns {CafeClient} CafeClient instance
     */
    CafeClient.create = function create(properties) {
        return new CafeClient(properties);
    };

    /**
     * Encodes the specified CafeClient message. Does not implicitly {@link CafeClient.verify|verify} messages.
     * @function encode
     * @memberof CafeClient
     * @static
     * @param {ICafeClient} message CafeClient message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeClient.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.address != null && message.hasOwnProperty("address"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.address);
        if (message.created != null && message.hasOwnProperty("created"))
            $root.google.protobuf.Timestamp.encode(message.created, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        if (message.seen != null && message.hasOwnProperty("seen"))
            $root.google.protobuf.Timestamp.encode(message.seen, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
        if (message.token != null && message.hasOwnProperty("token"))
            writer.uint32(/* id 5, wireType 2 =*/42).string(message.token);
        return writer;
    };

    /**
     * Encodes the specified CafeClient message, length delimited. Does not implicitly {@link CafeClient.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeClient
     * @static
     * @param {ICafeClient} message CafeClient message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeClient.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeClient message from the specified reader or buffer.
     * @function decode
     * @memberof CafeClient
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeClient} CafeClient
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeClient.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeClient();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.address = reader.string();
                break;
            case 3:
                message.created = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 4:
                message.seen = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 5:
                message.token = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeClient message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeClient
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeClient} CafeClient
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeClient.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeClient message.
     * @function verify
     * @memberof CafeClient
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeClient.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.address != null && message.hasOwnProperty("address"))
            if (!$util.isString(message.address))
                return "address: string expected";
        if (message.created != null && message.hasOwnProperty("created")) {
            let error = $root.google.protobuf.Timestamp.verify(message.created);
            if (error)
                return "created." + error;
        }
        if (message.seen != null && message.hasOwnProperty("seen")) {
            let error = $root.google.protobuf.Timestamp.verify(message.seen);
            if (error)
                return "seen." + error;
        }
        if (message.token != null && message.hasOwnProperty("token"))
            if (!$util.isString(message.token))
                return "token: string expected";
        return null;
    };

    /**
     * Creates a CafeClient message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeClient
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeClient} CafeClient
     */
    CafeClient.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeClient)
            return object;
        let message = new $root.CafeClient();
        if (object.id != null)
            message.id = String(object.id);
        if (object.address != null)
            message.address = String(object.address);
        if (object.created != null) {
            if (typeof object.created !== "object")
                throw TypeError(".CafeClient.created: object expected");
            message.created = $root.google.protobuf.Timestamp.fromObject(object.created);
        }
        if (object.seen != null) {
            if (typeof object.seen !== "object")
                throw TypeError(".CafeClient.seen: object expected");
            message.seen = $root.google.protobuf.Timestamp.fromObject(object.seen);
        }
        if (object.token != null)
            message.token = String(object.token);
        return message;
    };

    /**
     * Creates a plain object from a CafeClient message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeClient
     * @static
     * @param {CafeClient} message CafeClient
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeClient.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            object.address = "";
            object.created = null;
            object.seen = null;
            object.token = "";
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.address != null && message.hasOwnProperty("address"))
            object.address = message.address;
        if (message.created != null && message.hasOwnProperty("created"))
            object.created = $root.google.protobuf.Timestamp.toObject(message.created, options);
        if (message.seen != null && message.hasOwnProperty("seen"))
            object.seen = $root.google.protobuf.Timestamp.toObject(message.seen, options);
        if (message.token != null && message.hasOwnProperty("token"))
            object.token = message.token;
        return object;
    };

    /**
     * Converts this CafeClient to JSON.
     * @function toJSON
     * @memberof CafeClient
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeClient.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeClient;
})();

export const CafeToken = $root.CafeToken = (() => {

    /**
     * Properties of a CafeToken.
     * @exports ICafeToken
     * @interface ICafeToken
     * @property {string} id CafeToken id
     * @property {Uint8Array} value CafeToken value
     * @property {google.protobuf.ITimestamp} date CafeToken date
     */

    /**
     * Constructs a new CafeToken.
     * @exports CafeToken
     * @classdesc Represents a CafeToken.
     * @implements ICafeToken
     * @constructor
     * @param {ICafeToken=} [properties] Properties to set
     */
    function CafeToken(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeToken id.
     * @member {string} id
     * @memberof CafeToken
     * @instance
     */
    CafeToken.prototype.id = "";

    /**
     * CafeToken value.
     * @member {Uint8Array} value
     * @memberof CafeToken
     * @instance
     */
    CafeToken.prototype.value = $util.newBuffer([]);

    /**
     * CafeToken date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof CafeToken
     * @instance
     */
    CafeToken.prototype.date = null;

    /**
     * Creates a new CafeToken instance using the specified properties.
     * @function create
     * @memberof CafeToken
     * @static
     * @param {ICafeToken=} [properties] Properties to set
     * @returns {CafeToken} CafeToken instance
     */
    CafeToken.create = function create(properties) {
        return new CafeToken(properties);
    };

    /**
     * Encodes the specified CafeToken message. Does not implicitly {@link CafeToken.verify|verify} messages.
     * @function encode
     * @memberof CafeToken
     * @static
     * @param {ICafeToken} message CafeToken message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeToken.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.value != null && message.hasOwnProperty("value"))
            writer.uint32(/* id 2, wireType 2 =*/18).bytes(message.value);
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified CafeToken message, length delimited. Does not implicitly {@link CafeToken.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeToken
     * @static
     * @param {ICafeToken} message CafeToken message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeToken.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeToken message from the specified reader or buffer.
     * @function decode
     * @memberof CafeToken
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeToken} CafeToken
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeToken.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeToken();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.value = reader.bytes();
                break;
            case 3:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeToken message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeToken
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeToken} CafeToken
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeToken.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeToken message.
     * @function verify
     * @memberof CafeToken
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeToken.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.value != null && message.hasOwnProperty("value"))
            if (!(message.value && typeof message.value.length === "number" || $util.isString(message.value)))
                return "value: buffer expected";
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        return null;
    };

    /**
     * Creates a CafeToken message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeToken
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeToken} CafeToken
     */
    CafeToken.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeToken)
            return object;
        let message = new $root.CafeToken();
        if (object.id != null)
            message.id = String(object.id);
        if (object.value != null)
            if (typeof object.value === "string")
                $util.base64.decode(object.value, message.value = $util.newBuffer($util.base64.length(object.value)), 0);
            else if (object.value.length)
                message.value = object.value;
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".CafeToken.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        return message;
    };

    /**
     * Creates a plain object from a CafeToken message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeToken
     * @static
     * @param {CafeToken} message CafeToken
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeToken.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            if (options.bytes === String)
                object.value = "";
            else {
                object.value = [];
                if (options.bytes !== Array)
                    object.value = $util.newBuffer(object.value);
            }
            object.date = null;
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.value != null && message.hasOwnProperty("value"))
            object.value = options.bytes === String ? $util.base64.encode(message.value, 0, message.value.length) : options.bytes === Array ? Array.prototype.slice.call(message.value) : message.value;
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        return object;
    };

    /**
     * Converts this CafeToken to JSON.
     * @function toJSON
     * @memberof CafeToken
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeToken.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeToken;
})();

export const CafeClientThread = $root.CafeClientThread = (() => {

    /**
     * Properties of a CafeClientThread.
     * @exports ICafeClientThread
     * @interface ICafeClientThread
     * @property {string} id CafeClientThread id
     * @property {string} client CafeClientThread client
     * @property {Uint8Array} ciphertext CafeClientThread ciphertext
     */

    /**
     * Constructs a new CafeClientThread.
     * @exports CafeClientThread
     * @classdesc Represents a CafeClientThread.
     * @implements ICafeClientThread
     * @constructor
     * @param {ICafeClientThread=} [properties] Properties to set
     */
    function CafeClientThread(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeClientThread id.
     * @member {string} id
     * @memberof CafeClientThread
     * @instance
     */
    CafeClientThread.prototype.id = "";

    /**
     * CafeClientThread client.
     * @member {string} client
     * @memberof CafeClientThread
     * @instance
     */
    CafeClientThread.prototype.client = "";

    /**
     * CafeClientThread ciphertext.
     * @member {Uint8Array} ciphertext
     * @memberof CafeClientThread
     * @instance
     */
    CafeClientThread.prototype.ciphertext = $util.newBuffer([]);

    /**
     * Creates a new CafeClientThread instance using the specified properties.
     * @function create
     * @memberof CafeClientThread
     * @static
     * @param {ICafeClientThread=} [properties] Properties to set
     * @returns {CafeClientThread} CafeClientThread instance
     */
    CafeClientThread.create = function create(properties) {
        return new CafeClientThread(properties);
    };

    /**
     * Encodes the specified CafeClientThread message. Does not implicitly {@link CafeClientThread.verify|verify} messages.
     * @function encode
     * @memberof CafeClientThread
     * @static
     * @param {ICafeClientThread} message CafeClientThread message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeClientThread.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.client != null && message.hasOwnProperty("client"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.client);
        if (message.ciphertext != null && message.hasOwnProperty("ciphertext"))
            writer.uint32(/* id 3, wireType 2 =*/26).bytes(message.ciphertext);
        return writer;
    };

    /**
     * Encodes the specified CafeClientThread message, length delimited. Does not implicitly {@link CafeClientThread.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeClientThread
     * @static
     * @param {ICafeClientThread} message CafeClientThread message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeClientThread.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeClientThread message from the specified reader or buffer.
     * @function decode
     * @memberof CafeClientThread
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeClientThread} CafeClientThread
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeClientThread.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeClientThread();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.client = reader.string();
                break;
            case 3:
                message.ciphertext = reader.bytes();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeClientThread message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeClientThread
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeClientThread} CafeClientThread
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeClientThread.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeClientThread message.
     * @function verify
     * @memberof CafeClientThread
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeClientThread.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.client != null && message.hasOwnProperty("client"))
            if (!$util.isString(message.client))
                return "client: string expected";
        if (message.ciphertext != null && message.hasOwnProperty("ciphertext"))
            if (!(message.ciphertext && typeof message.ciphertext.length === "number" || $util.isString(message.ciphertext)))
                return "ciphertext: buffer expected";
        return null;
    };

    /**
     * Creates a CafeClientThread message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeClientThread
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeClientThread} CafeClientThread
     */
    CafeClientThread.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeClientThread)
            return object;
        let message = new $root.CafeClientThread();
        if (object.id != null)
            message.id = String(object.id);
        if (object.client != null)
            message.client = String(object.client);
        if (object.ciphertext != null)
            if (typeof object.ciphertext === "string")
                $util.base64.decode(object.ciphertext, message.ciphertext = $util.newBuffer($util.base64.length(object.ciphertext)), 0);
            else if (object.ciphertext.length)
                message.ciphertext = object.ciphertext;
        return message;
    };

    /**
     * Creates a plain object from a CafeClientThread message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeClientThread
     * @static
     * @param {CafeClientThread} message CafeClientThread
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeClientThread.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            object.client = "";
            if (options.bytes === String)
                object.ciphertext = "";
            else {
                object.ciphertext = [];
                if (options.bytes !== Array)
                    object.ciphertext = $util.newBuffer(object.ciphertext);
            }
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.client != null && message.hasOwnProperty("client"))
            object.client = message.client;
        if (message.ciphertext != null && message.hasOwnProperty("ciphertext"))
            object.ciphertext = options.bytes === String ? $util.base64.encode(message.ciphertext, 0, message.ciphertext.length) : options.bytes === Array ? Array.prototype.slice.call(message.ciphertext) : message.ciphertext;
        return object;
    };

    /**
     * Converts this CafeClientThread to JSON.
     * @function toJSON
     * @memberof CafeClientThread
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeClientThread.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeClientThread;
})();

export const CafeClientMessage = $root.CafeClientMessage = (() => {

    /**
     * Properties of a CafeClientMessage.
     * @exports ICafeClientMessage
     * @interface ICafeClientMessage
     * @property {string} id CafeClientMessage id
     * @property {string} peer CafeClientMessage peer
     * @property {string} client CafeClientMessage client
     * @property {google.protobuf.ITimestamp} date CafeClientMessage date
     */

    /**
     * Constructs a new CafeClientMessage.
     * @exports CafeClientMessage
     * @classdesc Represents a CafeClientMessage.
     * @implements ICafeClientMessage
     * @constructor
     * @param {ICafeClientMessage=} [properties] Properties to set
     */
    function CafeClientMessage(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CafeClientMessage id.
     * @member {string} id
     * @memberof CafeClientMessage
     * @instance
     */
    CafeClientMessage.prototype.id = "";

    /**
     * CafeClientMessage peer.
     * @member {string} peer
     * @memberof CafeClientMessage
     * @instance
     */
    CafeClientMessage.prototype.peer = "";

    /**
     * CafeClientMessage client.
     * @member {string} client
     * @memberof CafeClientMessage
     * @instance
     */
    CafeClientMessage.prototype.client = "";

    /**
     * CafeClientMessage date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof CafeClientMessage
     * @instance
     */
    CafeClientMessage.prototype.date = null;

    /**
     * Creates a new CafeClientMessage instance using the specified properties.
     * @function create
     * @memberof CafeClientMessage
     * @static
     * @param {ICafeClientMessage=} [properties] Properties to set
     * @returns {CafeClientMessage} CafeClientMessage instance
     */
    CafeClientMessage.create = function create(properties) {
        return new CafeClientMessage(properties);
    };

    /**
     * Encodes the specified CafeClientMessage message. Does not implicitly {@link CafeClientMessage.verify|verify} messages.
     * @function encode
     * @memberof CafeClientMessage
     * @static
     * @param {ICafeClientMessage} message CafeClientMessage message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeClientMessage.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.peer != null && message.hasOwnProperty("peer"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.peer);
        if (message.client != null && message.hasOwnProperty("client"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.client);
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified CafeClientMessage message, length delimited. Does not implicitly {@link CafeClientMessage.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CafeClientMessage
     * @static
     * @param {ICafeClientMessage} message CafeClientMessage message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CafeClientMessage.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CafeClientMessage message from the specified reader or buffer.
     * @function decode
     * @memberof CafeClientMessage
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CafeClientMessage} CafeClientMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeClientMessage.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CafeClientMessage();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.peer = reader.string();
                break;
            case 3:
                message.client = reader.string();
                break;
            case 4:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CafeClientMessage message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CafeClientMessage
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CafeClientMessage} CafeClientMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CafeClientMessage.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CafeClientMessage message.
     * @function verify
     * @memberof CafeClientMessage
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CafeClientMessage.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.peer != null && message.hasOwnProperty("peer"))
            if (!$util.isString(message.peer))
                return "peer: string expected";
        if (message.client != null && message.hasOwnProperty("client"))
            if (!$util.isString(message.client))
                return "client: string expected";
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        return null;
    };

    /**
     * Creates a CafeClientMessage message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CafeClientMessage
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CafeClientMessage} CafeClientMessage
     */
    CafeClientMessage.fromObject = function fromObject(object) {
        if (object instanceof $root.CafeClientMessage)
            return object;
        let message = new $root.CafeClientMessage();
        if (object.id != null)
            message.id = String(object.id);
        if (object.peer != null)
            message.peer = String(object.peer);
        if (object.client != null)
            message.client = String(object.client);
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".CafeClientMessage.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        return message;
    };

    /**
     * Creates a plain object from a CafeClientMessage message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CafeClientMessage
     * @static
     * @param {CafeClientMessage} message CafeClientMessage
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CafeClientMessage.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            object.peer = "";
            object.client = "";
            object.date = null;
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.peer != null && message.hasOwnProperty("peer"))
            object.peer = message.peer;
        if (message.client != null && message.hasOwnProperty("client"))
            object.client = message.client;
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        return object;
    };

    /**
     * Converts this CafeClientMessage to JSON.
     * @function toJSON
     * @memberof CafeClientMessage
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CafeClientMessage.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CafeClientMessage;
})();

export const Message = $root.Message = (() => {

    /**
     * Properties of a Message.
     * @exports IMessage
     * @interface IMessage
     * @property {Message.Type} type Message type
     * @property {google.protobuf.IAny} payload Message payload
     * @property {number} requestId Message requestId
     * @property {boolean} isResponse Message isResponse
     */

    /**
     * Constructs a new Message.
     * @exports Message
     * @classdesc Represents a Message.
     * @implements IMessage
     * @constructor
     * @param {IMessage=} [properties] Properties to set
     */
    function Message(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Message type.
     * @member {Message.Type} type
     * @memberof Message
     * @instance
     */
    Message.prototype.type = 0;

    /**
     * Message payload.
     * @member {google.protobuf.IAny} payload
     * @memberof Message
     * @instance
     */
    Message.prototype.payload = null;

    /**
     * Message requestId.
     * @member {number} requestId
     * @memberof Message
     * @instance
     */
    Message.prototype.requestId = 0;

    /**
     * Message isResponse.
     * @member {boolean} isResponse
     * @memberof Message
     * @instance
     */
    Message.prototype.isResponse = false;

    /**
     * Creates a new Message instance using the specified properties.
     * @function create
     * @memberof Message
     * @static
     * @param {IMessage=} [properties] Properties to set
     * @returns {Message} Message instance
     */
    Message.create = function create(properties) {
        return new Message(properties);
    };

    /**
     * Encodes the specified Message message. Does not implicitly {@link Message.verify|verify} messages.
     * @function encode
     * @memberof Message
     * @static
     * @param {IMessage} message Message message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Message.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.type != null && message.hasOwnProperty("type"))
            writer.uint32(/* id 1, wireType 0 =*/8).int32(message.type);
        if (message.payload != null && message.hasOwnProperty("payload"))
            $root.google.protobuf.Any.encode(message.payload, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        if (message.requestId != null && message.hasOwnProperty("requestId"))
            writer.uint32(/* id 3, wireType 0 =*/24).int32(message.requestId);
        if (message.isResponse != null && message.hasOwnProperty("isResponse"))
            writer.uint32(/* id 4, wireType 0 =*/32).bool(message.isResponse);
        return writer;
    };

    /**
     * Encodes the specified Message message, length delimited. Does not implicitly {@link Message.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Message
     * @static
     * @param {IMessage} message Message message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Message.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Message message from the specified reader or buffer.
     * @function decode
     * @memberof Message
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Message} Message
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Message.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Message();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.type = reader.int32();
                break;
            case 2:
                message.payload = $root.google.protobuf.Any.decode(reader, reader.uint32());
                break;
            case 3:
                message.requestId = reader.int32();
                break;
            case 4:
                message.isResponse = reader.bool();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Message message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Message
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Message} Message
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Message.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Message message.
     * @function verify
     * @memberof Message
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Message.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.type != null && message.hasOwnProperty("type"))
            switch (message.type) {
            default:
                return "type: enum value expected";
            case 0:
            case 1:
            case 10:
            case 50:
            case 51:
            case 52:
            case 53:
            case 54:
            case 55:
            case 56:
            case 57:
            case 58:
            case 59:
            case 60:
            case 61:
            case 62:
            case 63:
            case 64:
            case 65:
            case 66:
            case 67:
            case 70:
            case 71:
            case 102:
            case 103:
            case 500:
            case 68:
            case 69:
            case 100:
            case 101:
                break;
            }
        if (message.payload != null && message.hasOwnProperty("payload")) {
            let error = $root.google.protobuf.Any.verify(message.payload);
            if (error)
                return "payload." + error;
        }
        if (message.requestId != null && message.hasOwnProperty("requestId"))
            if (!$util.isInteger(message.requestId))
                return "requestId: integer expected";
        if (message.isResponse != null && message.hasOwnProperty("isResponse"))
            if (typeof message.isResponse !== "boolean")
                return "isResponse: boolean expected";
        return null;
    };

    /**
     * Creates a Message message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Message
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Message} Message
     */
    Message.fromObject = function fromObject(object) {
        if (object instanceof $root.Message)
            return object;
        let message = new $root.Message();
        switch (object.type) {
        case "PING":
        case 0:
            message.type = 0;
            break;
        case "PONG":
        case 1:
            message.type = 1;
            break;
        case "THREAD_ENVELOPE":
        case 10:
            message.type = 10;
            break;
        case "CAFE_CHALLENGE":
        case 50:
            message.type = 50;
            break;
        case "CAFE_NONCE":
        case 51:
            message.type = 51;
            break;
        case "CAFE_REGISTRATION":
        case 52:
            message.type = 52;
            break;
        case "CAFE_SESSION":
        case 53:
            message.type = 53;
            break;
        case "CAFE_REFRESH_SESSION":
        case 54:
            message.type = 54;
            break;
        case "CAFE_STORE":
        case 55:
            message.type = 55;
            break;
        case "CAFE_OBJECT":
        case 56:
            message.type = 56;
            break;
        case "CAFE_OBJECT_LIST":
        case 57:
            message.type = 57;
            break;
        case "CAFE_STORE_THREAD":
        case 58:
            message.type = 58;
            break;
        case "CAFE_STORED":
        case 59:
            message.type = 59;
            break;
        case "CAFE_DELIVER_MESSAGE":
        case 60:
            message.type = 60;
            break;
        case "CAFE_CHECK_MESSAGES":
        case 61:
            message.type = 61;
            break;
        case "CAFE_MESSAGES":
        case 62:
            message.type = 62;
            break;
        case "CAFE_DELETE_MESSAGES":
        case 63:
            message.type = 63;
            break;
        case "CAFE_DELETE_MESSAGES_ACK":
        case 64:
            message.type = 64;
            break;
        case "CAFE_YOU_HAVE_MAIL":
        case 65:
            message.type = 65;
            break;
        case "CAFE_PUBLISH_CONTACT":
        case 66:
            message.type = 66;
            break;
        case "CAFE_PUBLISH_CONTACT_ACK":
        case 67:
            message.type = 67;
            break;
        case "CAFE_QUERY":
        case 70:
            message.type = 70;
            break;
        case "CAFE_QUERY_RES":
        case 71:
            message.type = 71;
            break;
        case "CAFE_PUBSUB_QUERY":
        case 102:
            message.type = 102;
            break;
        case "CAFE_PUBSUB_QUERY_RES":
        case 103:
            message.type = 103;
            break;
        case "ERROR":
        case 500:
            message.type = 500;
            break;
        case "CAFE_CONTACT_QUERY":
        case 68:
            message.type = 68;
            break;
        case "CAFE_CONTACT_QUERY_RES":
        case 69:
            message.type = 69;
            break;
        case "CAFE_PUBSUB_CONTACT_QUERY":
        case 100:
            message.type = 100;
            break;
        case "CAFE_PUBSUB_CONTACT_QUERY_RES":
        case 101:
            message.type = 101;
            break;
        }
        if (object.payload != null) {
            if (typeof object.payload !== "object")
                throw TypeError(".Message.payload: object expected");
            message.payload = $root.google.protobuf.Any.fromObject(object.payload);
        }
        if (object.requestId != null)
            message.requestId = object.requestId | 0;
        if (object.isResponse != null)
            message.isResponse = Boolean(object.isResponse);
        return message;
    };

    /**
     * Creates a plain object from a Message message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Message
     * @static
     * @param {Message} message Message
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Message.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.type = options.enums === String ? "PING" : 0;
            object.payload = null;
            object.requestId = 0;
            object.isResponse = false;
        }
        if (message.type != null && message.hasOwnProperty("type"))
            object.type = options.enums === String ? $root.Message.Type[message.type] : message.type;
        if (message.payload != null && message.hasOwnProperty("payload"))
            object.payload = $root.google.protobuf.Any.toObject(message.payload, options);
        if (message.requestId != null && message.hasOwnProperty("requestId"))
            object.requestId = message.requestId;
        if (message.isResponse != null && message.hasOwnProperty("isResponse"))
            object.isResponse = message.isResponse;
        return object;
    };

    /**
     * Converts this Message to JSON.
     * @function toJSON
     * @memberof Message
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Message.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    /**
     * Type enum.
     * @name Message.Type
     * @enum {string}
     * @property {number} PING=0 PING value
     * @property {number} PONG=1 PONG value
     * @property {number} THREAD_ENVELOPE=10 THREAD_ENVELOPE value
     * @property {number} CAFE_CHALLENGE=50 CAFE_CHALLENGE value
     * @property {number} CAFE_NONCE=51 CAFE_NONCE value
     * @property {number} CAFE_REGISTRATION=52 CAFE_REGISTRATION value
     * @property {number} CAFE_SESSION=53 CAFE_SESSION value
     * @property {number} CAFE_REFRESH_SESSION=54 CAFE_REFRESH_SESSION value
     * @property {number} CAFE_STORE=55 CAFE_STORE value
     * @property {number} CAFE_OBJECT=56 CAFE_OBJECT value
     * @property {number} CAFE_OBJECT_LIST=57 CAFE_OBJECT_LIST value
     * @property {number} CAFE_STORE_THREAD=58 CAFE_STORE_THREAD value
     * @property {number} CAFE_STORED=59 CAFE_STORED value
     * @property {number} CAFE_DELIVER_MESSAGE=60 CAFE_DELIVER_MESSAGE value
     * @property {number} CAFE_CHECK_MESSAGES=61 CAFE_CHECK_MESSAGES value
     * @property {number} CAFE_MESSAGES=62 CAFE_MESSAGES value
     * @property {number} CAFE_DELETE_MESSAGES=63 CAFE_DELETE_MESSAGES value
     * @property {number} CAFE_DELETE_MESSAGES_ACK=64 CAFE_DELETE_MESSAGES_ACK value
     * @property {number} CAFE_YOU_HAVE_MAIL=65 CAFE_YOU_HAVE_MAIL value
     * @property {number} CAFE_PUBLISH_CONTACT=66 CAFE_PUBLISH_CONTACT value
     * @property {number} CAFE_PUBLISH_CONTACT_ACK=67 CAFE_PUBLISH_CONTACT_ACK value
     * @property {number} CAFE_QUERY=70 CAFE_QUERY value
     * @property {number} CAFE_QUERY_RES=71 CAFE_QUERY_RES value
     * @property {number} CAFE_PUBSUB_QUERY=102 CAFE_PUBSUB_QUERY value
     * @property {number} CAFE_PUBSUB_QUERY_RES=103 CAFE_PUBSUB_QUERY_RES value
     * @property {number} ERROR=500 ERROR value
     * @property {number} CAFE_CONTACT_QUERY=68 CAFE_CONTACT_QUERY value
     * @property {number} CAFE_CONTACT_QUERY_RES=69 CAFE_CONTACT_QUERY_RES value
     * @property {number} CAFE_PUBSUB_CONTACT_QUERY=100 CAFE_PUBSUB_CONTACT_QUERY value
     * @property {number} CAFE_PUBSUB_CONTACT_QUERY_RES=101 CAFE_PUBSUB_CONTACT_QUERY_RES value
     */
    Message.Type = (function() {
        const valuesById = {}, values = Object.create(valuesById);
        values[valuesById[0] = "PING"] = 0;
        values[valuesById[1] = "PONG"] = 1;
        values[valuesById[10] = "THREAD_ENVELOPE"] = 10;
        values[valuesById[50] = "CAFE_CHALLENGE"] = 50;
        values[valuesById[51] = "CAFE_NONCE"] = 51;
        values[valuesById[52] = "CAFE_REGISTRATION"] = 52;
        values[valuesById[53] = "CAFE_SESSION"] = 53;
        values[valuesById[54] = "CAFE_REFRESH_SESSION"] = 54;
        values[valuesById[55] = "CAFE_STORE"] = 55;
        values[valuesById[56] = "CAFE_OBJECT"] = 56;
        values[valuesById[57] = "CAFE_OBJECT_LIST"] = 57;
        values[valuesById[58] = "CAFE_STORE_THREAD"] = 58;
        values[valuesById[59] = "CAFE_STORED"] = 59;
        values[valuesById[60] = "CAFE_DELIVER_MESSAGE"] = 60;
        values[valuesById[61] = "CAFE_CHECK_MESSAGES"] = 61;
        values[valuesById[62] = "CAFE_MESSAGES"] = 62;
        values[valuesById[63] = "CAFE_DELETE_MESSAGES"] = 63;
        values[valuesById[64] = "CAFE_DELETE_MESSAGES_ACK"] = 64;
        values[valuesById[65] = "CAFE_YOU_HAVE_MAIL"] = 65;
        values[valuesById[66] = "CAFE_PUBLISH_CONTACT"] = 66;
        values[valuesById[67] = "CAFE_PUBLISH_CONTACT_ACK"] = 67;
        values[valuesById[70] = "CAFE_QUERY"] = 70;
        values[valuesById[71] = "CAFE_QUERY_RES"] = 71;
        values[valuesById[102] = "CAFE_PUBSUB_QUERY"] = 102;
        values[valuesById[103] = "CAFE_PUBSUB_QUERY_RES"] = 103;
        values[valuesById[500] = "ERROR"] = 500;
        values[valuesById[68] = "CAFE_CONTACT_QUERY"] = 68;
        values[valuesById[69] = "CAFE_CONTACT_QUERY_RES"] = 69;
        values[valuesById[100] = "CAFE_PUBSUB_CONTACT_QUERY"] = 100;
        values[valuesById[101] = "CAFE_PUBSUB_CONTACT_QUERY_RES"] = 101;
        return values;
    })();

    return Message;
})();

export const Envelope = $root.Envelope = (() => {

    /**
     * Properties of an Envelope.
     * @exports IEnvelope
     * @interface IEnvelope
     * @property {IMessage} message Envelope message
     * @property {Uint8Array} sig Envelope sig
     */

    /**
     * Constructs a new Envelope.
     * @exports Envelope
     * @classdesc Represents an Envelope.
     * @implements IEnvelope
     * @constructor
     * @param {IEnvelope=} [properties] Properties to set
     */
    function Envelope(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Envelope message.
     * @member {IMessage} message
     * @memberof Envelope
     * @instance
     */
    Envelope.prototype.message = null;

    /**
     * Envelope sig.
     * @member {Uint8Array} sig
     * @memberof Envelope
     * @instance
     */
    Envelope.prototype.sig = $util.newBuffer([]);

    /**
     * Creates a new Envelope instance using the specified properties.
     * @function create
     * @memberof Envelope
     * @static
     * @param {IEnvelope=} [properties] Properties to set
     * @returns {Envelope} Envelope instance
     */
    Envelope.create = function create(properties) {
        return new Envelope(properties);
    };

    /**
     * Encodes the specified Envelope message. Does not implicitly {@link Envelope.verify|verify} messages.
     * @function encode
     * @memberof Envelope
     * @static
     * @param {IEnvelope} message Envelope message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Envelope.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.message != null && message.hasOwnProperty("message"))
            $root.Message.encode(message.message, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        if (message.sig != null && message.hasOwnProperty("sig"))
            writer.uint32(/* id 2, wireType 2 =*/18).bytes(message.sig);
        return writer;
    };

    /**
     * Encodes the specified Envelope message, length delimited. Does not implicitly {@link Envelope.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Envelope
     * @static
     * @param {IEnvelope} message Envelope message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Envelope.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes an Envelope message from the specified reader or buffer.
     * @function decode
     * @memberof Envelope
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Envelope} Envelope
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Envelope.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Envelope();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.message = $root.Message.decode(reader, reader.uint32());
                break;
            case 2:
                message.sig = reader.bytes();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes an Envelope message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Envelope
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Envelope} Envelope
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Envelope.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies an Envelope message.
     * @function verify
     * @memberof Envelope
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Envelope.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.message != null && message.hasOwnProperty("message")) {
            let error = $root.Message.verify(message.message);
            if (error)
                return "message." + error;
        }
        if (message.sig != null && message.hasOwnProperty("sig"))
            if (!(message.sig && typeof message.sig.length === "number" || $util.isString(message.sig)))
                return "sig: buffer expected";
        return null;
    };

    /**
     * Creates an Envelope message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Envelope
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Envelope} Envelope
     */
    Envelope.fromObject = function fromObject(object) {
        if (object instanceof $root.Envelope)
            return object;
        let message = new $root.Envelope();
        if (object.message != null) {
            if (typeof object.message !== "object")
                throw TypeError(".Envelope.message: object expected");
            message.message = $root.Message.fromObject(object.message);
        }
        if (object.sig != null)
            if (typeof object.sig === "string")
                $util.base64.decode(object.sig, message.sig = $util.newBuffer($util.base64.length(object.sig)), 0);
            else if (object.sig.length)
                message.sig = object.sig;
        return message;
    };

    /**
     * Creates a plain object from an Envelope message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Envelope
     * @static
     * @param {Envelope} message Envelope
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Envelope.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.message = null;
            if (options.bytes === String)
                object.sig = "";
            else {
                object.sig = [];
                if (options.bytes !== Array)
                    object.sig = $util.newBuffer(object.sig);
            }
        }
        if (message.message != null && message.hasOwnProperty("message"))
            object.message = $root.Message.toObject(message.message, options);
        if (message.sig != null && message.hasOwnProperty("sig"))
            object.sig = options.bytes === String ? $util.base64.encode(message.sig, 0, message.sig.length) : options.bytes === Array ? Array.prototype.slice.call(message.sig) : message.sig;
        return object;
    };

    /**
     * Converts this Envelope to JSON.
     * @function toJSON
     * @memberof Envelope
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Envelope.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Envelope;
})();

export const Error = $root.Error = (() => {

    /**
     * Properties of an Error.
     * @exports IError
     * @interface IError
     * @property {number} code Error code
     * @property {string} message Error message
     */

    /**
     * Constructs a new Error.
     * @exports Error
     * @classdesc Represents an Error.
     * @implements IError
     * @constructor
     * @param {IError=} [properties] Properties to set
     */
    function Error(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Error code.
     * @member {number} code
     * @memberof Error
     * @instance
     */
    Error.prototype.code = 0;

    /**
     * Error message.
     * @member {string} message
     * @memberof Error
     * @instance
     */
    Error.prototype.message = "";

    /**
     * Creates a new Error instance using the specified properties.
     * @function create
     * @memberof Error
     * @static
     * @param {IError=} [properties] Properties to set
     * @returns {Error} Error instance
     */
    Error.create = function create(properties) {
        return new Error(properties);
    };

    /**
     * Encodes the specified Error message. Does not implicitly {@link Error.verify|verify} messages.
     * @function encode
     * @memberof Error
     * @static
     * @param {IError} message Error message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Error.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.code != null && message.hasOwnProperty("code"))
            writer.uint32(/* id 1, wireType 0 =*/8).uint32(message.code);
        if (message.message != null && message.hasOwnProperty("message"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.message);
        return writer;
    };

    /**
     * Encodes the specified Error message, length delimited. Does not implicitly {@link Error.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Error
     * @static
     * @param {IError} message Error message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Error.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes an Error message from the specified reader or buffer.
     * @function decode
     * @memberof Error
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Error} Error
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Error.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Error();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.code = reader.uint32();
                break;
            case 2:
                message.message = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes an Error message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Error
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Error} Error
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Error.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies an Error message.
     * @function verify
     * @memberof Error
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Error.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.code != null && message.hasOwnProperty("code"))
            if (!$util.isInteger(message.code))
                return "code: integer expected";
        if (message.message != null && message.hasOwnProperty("message"))
            if (!$util.isString(message.message))
                return "message: string expected";
        return null;
    };

    /**
     * Creates an Error message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Error
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Error} Error
     */
    Error.fromObject = function fromObject(object) {
        if (object instanceof $root.Error)
            return object;
        let message = new $root.Error();
        if (object.code != null)
            message.code = object.code >>> 0;
        if (object.message != null)
            message.message = String(object.message);
        return message;
    };

    /**
     * Creates a plain object from an Error message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Error
     * @static
     * @param {Error} message Error
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Error.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.code = 0;
            object.message = "";
        }
        if (message.code != null && message.hasOwnProperty("code"))
            object.code = message.code;
        if (message.message != null && message.hasOwnProperty("message"))
            object.message = message.message;
        return object;
    };

    /**
     * Converts this Error to JSON.
     * @function toJSON
     * @memberof Error
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Error.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Error;
})();

export const google = $root.google = (() => {

    /**
     * Namespace google.
     * @exports google
     * @namespace
     */
    const google = {};

    google.protobuf = (function() {

        /**
         * Namespace protobuf.
         * @memberof google
         * @namespace
         */
        const protobuf = {};

        protobuf.Timestamp = (function() {

            /**
             * Properties of a Timestamp.
             * @memberof google.protobuf
             * @interface ITimestamp
             * @property {number|Long} seconds Timestamp seconds
             * @property {number} nanos Timestamp nanos
             */

            /**
             * Constructs a new Timestamp.
             * @memberof google.protobuf
             * @classdesc Represents a Timestamp.
             * @implements ITimestamp
             * @constructor
             * @param {google.protobuf.ITimestamp=} [properties] Properties to set
             */
            function Timestamp(properties) {
                if (properties)
                    for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * Timestamp seconds.
             * @member {number|Long} seconds
             * @memberof google.protobuf.Timestamp
             * @instance
             */
            Timestamp.prototype.seconds = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

            /**
             * Timestamp nanos.
             * @member {number} nanos
             * @memberof google.protobuf.Timestamp
             * @instance
             */
            Timestamp.prototype.nanos = 0;

            /**
             * Creates a new Timestamp instance using the specified properties.
             * @function create
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {google.protobuf.ITimestamp=} [properties] Properties to set
             * @returns {google.protobuf.Timestamp} Timestamp instance
             */
            Timestamp.create = function create(properties) {
                return new Timestamp(properties);
            };

            /**
             * Encodes the specified Timestamp message. Does not implicitly {@link google.protobuf.Timestamp.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {google.protobuf.ITimestamp} message Timestamp message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Timestamp.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.seconds != null && message.hasOwnProperty("seconds"))
                    writer.uint32(/* id 1, wireType 0 =*/8).int64(message.seconds);
                if (message.nanos != null && message.hasOwnProperty("nanos"))
                    writer.uint32(/* id 2, wireType 0 =*/16).int32(message.nanos);
                return writer;
            };

            /**
             * Encodes the specified Timestamp message, length delimited. Does not implicitly {@link google.protobuf.Timestamp.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {google.protobuf.ITimestamp} message Timestamp message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Timestamp.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes a Timestamp message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.Timestamp} Timestamp
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Timestamp.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                let end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.Timestamp();
                while (reader.pos < end) {
                    let tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        message.seconds = reader.int64();
                        break;
                    case 2:
                        message.nanos = reader.int32();
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes a Timestamp message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.Timestamp} Timestamp
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Timestamp.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a Timestamp message.
             * @function verify
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            Timestamp.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.seconds != null && message.hasOwnProperty("seconds"))
                    if (!$util.isInteger(message.seconds) && !(message.seconds && $util.isInteger(message.seconds.low) && $util.isInteger(message.seconds.high)))
                        return "seconds: integer|Long expected";
                if (message.nanos != null && message.hasOwnProperty("nanos"))
                    if (!$util.isInteger(message.nanos))
                        return "nanos: integer expected";
                return null;
            };

            /**
             * Creates a Timestamp message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.Timestamp} Timestamp
             */
            Timestamp.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.Timestamp)
                    return object;
                let message = new $root.google.protobuf.Timestamp();
                if (object.seconds != null)
                    if ($util.Long)
                        (message.seconds = $util.Long.fromValue(object.seconds)).unsigned = false;
                    else if (typeof object.seconds === "string")
                        message.seconds = parseInt(object.seconds, 10);
                    else if (typeof object.seconds === "number")
                        message.seconds = object.seconds;
                    else if (typeof object.seconds === "object")
                        message.seconds = new $util.LongBits(object.seconds.low >>> 0, object.seconds.high >>> 0).toNumber();
                if (object.nanos != null)
                    message.nanos = object.nanos | 0;
                return message;
            };

            /**
             * Creates a plain object from a Timestamp message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {google.protobuf.Timestamp} message Timestamp
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            Timestamp.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                let object = {};
                if (options.defaults) {
                    if ($util.Long) {
                        let long = new $util.Long(0, 0, false);
                        object.seconds = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                    } else
                        object.seconds = options.longs === String ? "0" : 0;
                    object.nanos = 0;
                }
                if (message.seconds != null && message.hasOwnProperty("seconds"))
                    if (typeof message.seconds === "number")
                        object.seconds = options.longs === String ? String(message.seconds) : message.seconds;
                    else
                        object.seconds = options.longs === String ? $util.Long.prototype.toString.call(message.seconds) : options.longs === Number ? new $util.LongBits(message.seconds.low >>> 0, message.seconds.high >>> 0).toNumber() : message.seconds;
                if (message.nanos != null && message.hasOwnProperty("nanos"))
                    object.nanos = message.nanos;
                return object;
            };

            /**
             * Converts this Timestamp to JSON.
             * @function toJSON
             * @memberof google.protobuf.Timestamp
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            Timestamp.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            return Timestamp;
        })();

        protobuf.Struct = (function() {

            /**
             * Properties of a Struct.
             * @memberof google.protobuf
             * @interface IStruct
             * @property {Object.<string,google.protobuf.IValue>} fields Struct fields
             */

            /**
             * Constructs a new Struct.
             * @memberof google.protobuf
             * @classdesc Represents a Struct.
             * @implements IStruct
             * @constructor
             * @param {google.protobuf.IStruct=} [properties] Properties to set
             */
            function Struct(properties) {
                this.fields = {};
                if (properties)
                    for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * Struct fields.
             * @member {Object.<string,google.protobuf.IValue>} fields
             * @memberof google.protobuf.Struct
             * @instance
             */
            Struct.prototype.fields = $util.emptyObject;

            /**
             * Creates a new Struct instance using the specified properties.
             * @function create
             * @memberof google.protobuf.Struct
             * @static
             * @param {google.protobuf.IStruct=} [properties] Properties to set
             * @returns {google.protobuf.Struct} Struct instance
             */
            Struct.create = function create(properties) {
                return new Struct(properties);
            };

            /**
             * Encodes the specified Struct message. Does not implicitly {@link google.protobuf.Struct.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.Struct
             * @static
             * @param {google.protobuf.IStruct} message Struct message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Struct.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.fields != null && message.hasOwnProperty("fields"))
                    for (let keys = Object.keys(message.fields), i = 0; i < keys.length; ++i) {
                        writer.uint32(/* id 1, wireType 2 =*/10).fork().uint32(/* id 1, wireType 2 =*/10).string(keys[i]);
                        $root.google.protobuf.Value.encode(message.fields[keys[i]], writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim().ldelim();
                    }
                return writer;
            };

            /**
             * Encodes the specified Struct message, length delimited. Does not implicitly {@link google.protobuf.Struct.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.Struct
             * @static
             * @param {google.protobuf.IStruct} message Struct message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Struct.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes a Struct message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.Struct
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.Struct} Struct
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Struct.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                let end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.Struct(), key;
                while (reader.pos < end) {
                    let tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        reader.skip().pos++;
                        if (message.fields === $util.emptyObject)
                            message.fields = {};
                        key = reader.string();
                        reader.pos++;
                        message.fields[key] = $root.google.protobuf.Value.decode(reader, reader.uint32());
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes a Struct message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.Struct
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.Struct} Struct
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Struct.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a Struct message.
             * @function verify
             * @memberof google.protobuf.Struct
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            Struct.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.fields != null && message.hasOwnProperty("fields")) {
                    if (!$util.isObject(message.fields))
                        return "fields: object expected";
                    let key = Object.keys(message.fields);
                    for (let i = 0; i < key.length; ++i) {
                        let error = $root.google.protobuf.Value.verify(message.fields[key[i]]);
                        if (error)
                            return "fields." + error;
                    }
                }
                return null;
            };

            /**
             * Creates a Struct message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.Struct
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.Struct} Struct
             */
            Struct.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.Struct)
                    return object;
                let message = new $root.google.protobuf.Struct();
                if (object.fields) {
                    if (typeof object.fields !== "object")
                        throw TypeError(".google.protobuf.Struct.fields: object expected");
                    message.fields = {};
                    for (let keys = Object.keys(object.fields), i = 0; i < keys.length; ++i) {
                        if (typeof object.fields[keys[i]] !== "object")
                            throw TypeError(".google.protobuf.Struct.fields: object expected");
                        message.fields[keys[i]] = $root.google.protobuf.Value.fromObject(object.fields[keys[i]]);
                    }
                }
                return message;
            };

            /**
             * Creates a plain object from a Struct message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.Struct
             * @static
             * @param {google.protobuf.Struct} message Struct
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            Struct.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                let object = {};
                if (options.objects || options.defaults)
                    object.fields = {};
                let keys2;
                if (message.fields && (keys2 = Object.keys(message.fields)).length) {
                    object.fields = {};
                    for (let j = 0; j < keys2.length; ++j)
                        object.fields[keys2[j]] = $root.google.protobuf.Value.toObject(message.fields[keys2[j]], options);
                }
                return object;
            };

            /**
             * Converts this Struct to JSON.
             * @function toJSON
             * @memberof google.protobuf.Struct
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            Struct.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            return Struct;
        })();

        protobuf.Value = (function() {

            /**
             * Properties of a Value.
             * @memberof google.protobuf
             * @interface IValue
             * @property {google.protobuf.NullValue} nullValue Value nullValue
             * @property {number} numberValue Value numberValue
             * @property {string} stringValue Value stringValue
             * @property {boolean} boolValue Value boolValue
             * @property {google.protobuf.IStruct} structValue Value structValue
             * @property {google.protobuf.IListValue} listValue Value listValue
             */

            /**
             * Constructs a new Value.
             * @memberof google.protobuf
             * @classdesc Represents a Value.
             * @implements IValue
             * @constructor
             * @param {google.protobuf.IValue=} [properties] Properties to set
             */
            function Value(properties) {
                if (properties)
                    for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * Value nullValue.
             * @member {google.protobuf.NullValue} nullValue
             * @memberof google.protobuf.Value
             * @instance
             */
            Value.prototype.nullValue = 0;

            /**
             * Value numberValue.
             * @member {number} numberValue
             * @memberof google.protobuf.Value
             * @instance
             */
            Value.prototype.numberValue = 0;

            /**
             * Value stringValue.
             * @member {string} stringValue
             * @memberof google.protobuf.Value
             * @instance
             */
            Value.prototype.stringValue = "";

            /**
             * Value boolValue.
             * @member {boolean} boolValue
             * @memberof google.protobuf.Value
             * @instance
             */
            Value.prototype.boolValue = false;

            /**
             * Value structValue.
             * @member {google.protobuf.IStruct} structValue
             * @memberof google.protobuf.Value
             * @instance
             */
            Value.prototype.structValue = null;

            /**
             * Value listValue.
             * @member {google.protobuf.IListValue} listValue
             * @memberof google.protobuf.Value
             * @instance
             */
            Value.prototype.listValue = null;

            // OneOf field names bound to virtual getters and setters
            let $oneOfFields;

            /**
             * Value kind.
             * @member {"nullValue"|"numberValue"|"stringValue"|"boolValue"|"structValue"|"listValue"|undefined} kind
             * @memberof google.protobuf.Value
             * @instance
             */
            Object.defineProperty(Value.prototype, "kind", {
                get: $util.oneOfGetter($oneOfFields = ["nullValue", "numberValue", "stringValue", "boolValue", "structValue", "listValue"]),
                set: $util.oneOfSetter($oneOfFields)
            });

            /**
             * Creates a new Value instance using the specified properties.
             * @function create
             * @memberof google.protobuf.Value
             * @static
             * @param {google.protobuf.IValue=} [properties] Properties to set
             * @returns {google.protobuf.Value} Value instance
             */
            Value.create = function create(properties) {
                return new Value(properties);
            };

            /**
             * Encodes the specified Value message. Does not implicitly {@link google.protobuf.Value.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.Value
             * @static
             * @param {google.protobuf.IValue} message Value message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Value.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.nullValue != null && message.hasOwnProperty("nullValue"))
                    writer.uint32(/* id 1, wireType 0 =*/8).int32(message.nullValue);
                if (message.numberValue != null && message.hasOwnProperty("numberValue"))
                    writer.uint32(/* id 2, wireType 1 =*/17).double(message.numberValue);
                if (message.stringValue != null && message.hasOwnProperty("stringValue"))
                    writer.uint32(/* id 3, wireType 2 =*/26).string(message.stringValue);
                if (message.boolValue != null && message.hasOwnProperty("boolValue"))
                    writer.uint32(/* id 4, wireType 0 =*/32).bool(message.boolValue);
                if (message.structValue != null && message.hasOwnProperty("structValue"))
                    $root.google.protobuf.Struct.encode(message.structValue, writer.uint32(/* id 5, wireType 2 =*/42).fork()).ldelim();
                if (message.listValue != null && message.hasOwnProperty("listValue"))
                    $root.google.protobuf.ListValue.encode(message.listValue, writer.uint32(/* id 6, wireType 2 =*/50).fork()).ldelim();
                return writer;
            };

            /**
             * Encodes the specified Value message, length delimited. Does not implicitly {@link google.protobuf.Value.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.Value
             * @static
             * @param {google.protobuf.IValue} message Value message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Value.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes a Value message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.Value
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.Value} Value
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Value.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                let end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.Value();
                while (reader.pos < end) {
                    let tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        message.nullValue = reader.int32();
                        break;
                    case 2:
                        message.numberValue = reader.double();
                        break;
                    case 3:
                        message.stringValue = reader.string();
                        break;
                    case 4:
                        message.boolValue = reader.bool();
                        break;
                    case 5:
                        message.structValue = $root.google.protobuf.Struct.decode(reader, reader.uint32());
                        break;
                    case 6:
                        message.listValue = $root.google.protobuf.ListValue.decode(reader, reader.uint32());
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes a Value message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.Value
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.Value} Value
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Value.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a Value message.
             * @function verify
             * @memberof google.protobuf.Value
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            Value.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                let properties = {};
                if (message.nullValue != null && message.hasOwnProperty("nullValue")) {
                    properties.kind = 1;
                    switch (message.nullValue) {
                    default:
                        return "nullValue: enum value expected";
                    case 0:
                        break;
                    }
                }
                if (message.numberValue != null && message.hasOwnProperty("numberValue")) {
                    if (properties.kind === 1)
                        return "kind: multiple values";
                    properties.kind = 1;
                    if (typeof message.numberValue !== "number")
                        return "numberValue: number expected";
                }
                if (message.stringValue != null && message.hasOwnProperty("stringValue")) {
                    if (properties.kind === 1)
                        return "kind: multiple values";
                    properties.kind = 1;
                    if (!$util.isString(message.stringValue))
                        return "stringValue: string expected";
                }
                if (message.boolValue != null && message.hasOwnProperty("boolValue")) {
                    if (properties.kind === 1)
                        return "kind: multiple values";
                    properties.kind = 1;
                    if (typeof message.boolValue !== "boolean")
                        return "boolValue: boolean expected";
                }
                if (message.structValue != null && message.hasOwnProperty("structValue")) {
                    if (properties.kind === 1)
                        return "kind: multiple values";
                    properties.kind = 1;
                    {
                        let error = $root.google.protobuf.Struct.verify(message.structValue);
                        if (error)
                            return "structValue." + error;
                    }
                }
                if (message.listValue != null && message.hasOwnProperty("listValue")) {
                    if (properties.kind === 1)
                        return "kind: multiple values";
                    properties.kind = 1;
                    {
                        let error = $root.google.protobuf.ListValue.verify(message.listValue);
                        if (error)
                            return "listValue." + error;
                    }
                }
                return null;
            };

            /**
             * Creates a Value message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.Value
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.Value} Value
             */
            Value.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.Value)
                    return object;
                let message = new $root.google.protobuf.Value();
                switch (object.nullValue) {
                case "NULL_VALUE":
                case 0:
                    message.nullValue = 0;
                    break;
                }
                if (object.numberValue != null)
                    message.numberValue = Number(object.numberValue);
                if (object.stringValue != null)
                    message.stringValue = String(object.stringValue);
                if (object.boolValue != null)
                    message.boolValue = Boolean(object.boolValue);
                if (object.structValue != null) {
                    if (typeof object.structValue !== "object")
                        throw TypeError(".google.protobuf.Value.structValue: object expected");
                    message.structValue = $root.google.protobuf.Struct.fromObject(object.structValue);
                }
                if (object.listValue != null) {
                    if (typeof object.listValue !== "object")
                        throw TypeError(".google.protobuf.Value.listValue: object expected");
                    message.listValue = $root.google.protobuf.ListValue.fromObject(object.listValue);
                }
                return message;
            };

            /**
             * Creates a plain object from a Value message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.Value
             * @static
             * @param {google.protobuf.Value} message Value
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            Value.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                let object = {};
                if (message.nullValue != null && message.hasOwnProperty("nullValue")) {
                    object.nullValue = options.enums === String ? $root.google.protobuf.NullValue[message.nullValue] : message.nullValue;
                    if (options.oneofs)
                        object.kind = "nullValue";
                }
                if (message.numberValue != null && message.hasOwnProperty("numberValue")) {
                    object.numberValue = options.json && !isFinite(message.numberValue) ? String(message.numberValue) : message.numberValue;
                    if (options.oneofs)
                        object.kind = "numberValue";
                }
                if (message.stringValue != null && message.hasOwnProperty("stringValue")) {
                    object.stringValue = message.stringValue;
                    if (options.oneofs)
                        object.kind = "stringValue";
                }
                if (message.boolValue != null && message.hasOwnProperty("boolValue")) {
                    object.boolValue = message.boolValue;
                    if (options.oneofs)
                        object.kind = "boolValue";
                }
                if (message.structValue != null && message.hasOwnProperty("structValue")) {
                    object.structValue = $root.google.protobuf.Struct.toObject(message.structValue, options);
                    if (options.oneofs)
                        object.kind = "structValue";
                }
                if (message.listValue != null && message.hasOwnProperty("listValue")) {
                    object.listValue = $root.google.protobuf.ListValue.toObject(message.listValue, options);
                    if (options.oneofs)
                        object.kind = "listValue";
                }
                return object;
            };

            /**
             * Converts this Value to JSON.
             * @function toJSON
             * @memberof google.protobuf.Value
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            Value.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            return Value;
        })();

        /**
         * NullValue enum.
         * @name google.protobuf.NullValue
         * @enum {string}
         * @property {number} NULL_VALUE=0 NULL_VALUE value
         */
        protobuf.NullValue = (function() {
            const valuesById = {}, values = Object.create(valuesById);
            values[valuesById[0] = "NULL_VALUE"] = 0;
            return values;
        })();

        protobuf.ListValue = (function() {

            /**
             * Properties of a ListValue.
             * @memberof google.protobuf
             * @interface IListValue
             * @property {Array.<google.protobuf.IValue>} values ListValue values
             */

            /**
             * Constructs a new ListValue.
             * @memberof google.protobuf
             * @classdesc Represents a ListValue.
             * @implements IListValue
             * @constructor
             * @param {google.protobuf.IListValue=} [properties] Properties to set
             */
            function ListValue(properties) {
                this.values = [];
                if (properties)
                    for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * ListValue values.
             * @member {Array.<google.protobuf.IValue>} values
             * @memberof google.protobuf.ListValue
             * @instance
             */
            ListValue.prototype.values = $util.emptyArray;

            /**
             * Creates a new ListValue instance using the specified properties.
             * @function create
             * @memberof google.protobuf.ListValue
             * @static
             * @param {google.protobuf.IListValue=} [properties] Properties to set
             * @returns {google.protobuf.ListValue} ListValue instance
             */
            ListValue.create = function create(properties) {
                return new ListValue(properties);
            };

            /**
             * Encodes the specified ListValue message. Does not implicitly {@link google.protobuf.ListValue.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.ListValue
             * @static
             * @param {google.protobuf.IListValue} message ListValue message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            ListValue.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.values != null && message.values.length)
                    for (let i = 0; i < message.values.length; ++i)
                        $root.google.protobuf.Value.encode(message.values[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
                return writer;
            };

            /**
             * Encodes the specified ListValue message, length delimited. Does not implicitly {@link google.protobuf.ListValue.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.ListValue
             * @static
             * @param {google.protobuf.IListValue} message ListValue message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            ListValue.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes a ListValue message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.ListValue
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.ListValue} ListValue
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            ListValue.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                let end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.ListValue();
                while (reader.pos < end) {
                    let tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        if (!(message.values && message.values.length))
                            message.values = [];
                        message.values.push($root.google.protobuf.Value.decode(reader, reader.uint32()));
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes a ListValue message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.ListValue
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.ListValue} ListValue
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            ListValue.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a ListValue message.
             * @function verify
             * @memberof google.protobuf.ListValue
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            ListValue.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.values != null && message.hasOwnProperty("values")) {
                    if (!Array.isArray(message.values))
                        return "values: array expected";
                    for (let i = 0; i < message.values.length; ++i) {
                        let error = $root.google.protobuf.Value.verify(message.values[i]);
                        if (error)
                            return "values." + error;
                    }
                }
                return null;
            };

            /**
             * Creates a ListValue message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.ListValue
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.ListValue} ListValue
             */
            ListValue.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.ListValue)
                    return object;
                let message = new $root.google.protobuf.ListValue();
                if (object.values) {
                    if (!Array.isArray(object.values))
                        throw TypeError(".google.protobuf.ListValue.values: array expected");
                    message.values = [];
                    for (let i = 0; i < object.values.length; ++i) {
                        if (typeof object.values[i] !== "object")
                            throw TypeError(".google.protobuf.ListValue.values: object expected");
                        message.values[i] = $root.google.protobuf.Value.fromObject(object.values[i]);
                    }
                }
                return message;
            };

            /**
             * Creates a plain object from a ListValue message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.ListValue
             * @static
             * @param {google.protobuf.ListValue} message ListValue
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            ListValue.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                let object = {};
                if (options.arrays || options.defaults)
                    object.values = [];
                if (message.values && message.values.length) {
                    object.values = [];
                    for (let j = 0; j < message.values.length; ++j)
                        object.values[j] = $root.google.protobuf.Value.toObject(message.values[j], options);
                }
                return object;
            };

            /**
             * Converts this ListValue to JSON.
             * @function toJSON
             * @memberof google.protobuf.ListValue
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            ListValue.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            return ListValue;
        })();

        protobuf.Any = (function() {

            /**
             * Properties of an Any.
             * @memberof google.protobuf
             * @interface IAny
             * @property {string} type_url Any type_url
             * @property {Uint8Array} value Any value
             */

            /**
             * Constructs a new Any.
             * @memberof google.protobuf
             * @classdesc Represents an Any.
             * @implements IAny
             * @constructor
             * @param {google.protobuf.IAny=} [properties] Properties to set
             */
            function Any(properties) {
                if (properties)
                    for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * Any type_url.
             * @member {string} type_url
             * @memberof google.protobuf.Any
             * @instance
             */
            Any.prototype.type_url = "";

            /**
             * Any value.
             * @member {Uint8Array} value
             * @memberof google.protobuf.Any
             * @instance
             */
            Any.prototype.value = $util.newBuffer([]);

            /**
             * Creates a new Any instance using the specified properties.
             * @function create
             * @memberof google.protobuf.Any
             * @static
             * @param {google.protobuf.IAny=} [properties] Properties to set
             * @returns {google.protobuf.Any} Any instance
             */
            Any.create = function create(properties) {
                return new Any(properties);
            };

            /**
             * Encodes the specified Any message. Does not implicitly {@link google.protobuf.Any.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.Any
             * @static
             * @param {google.protobuf.IAny} message Any message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Any.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.type_url != null && message.hasOwnProperty("type_url"))
                    writer.uint32(/* id 1, wireType 2 =*/10).string(message.type_url);
                if (message.value != null && message.hasOwnProperty("value"))
                    writer.uint32(/* id 2, wireType 2 =*/18).bytes(message.value);
                return writer;
            };

            /**
             * Encodes the specified Any message, length delimited. Does not implicitly {@link google.protobuf.Any.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.Any
             * @static
             * @param {google.protobuf.IAny} message Any message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Any.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes an Any message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.Any
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.Any} Any
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Any.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                let end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.Any();
                while (reader.pos < end) {
                    let tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        message.type_url = reader.string();
                        break;
                    case 2:
                        message.value = reader.bytes();
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes an Any message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.Any
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.Any} Any
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Any.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies an Any message.
             * @function verify
             * @memberof google.protobuf.Any
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            Any.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.type_url != null && message.hasOwnProperty("type_url"))
                    if (!$util.isString(message.type_url))
                        return "type_url: string expected";
                if (message.value != null && message.hasOwnProperty("value"))
                    if (!(message.value && typeof message.value.length === "number" || $util.isString(message.value)))
                        return "value: buffer expected";
                return null;
            };

            /**
             * Creates an Any message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.Any
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.Any} Any
             */
            Any.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.Any)
                    return object;
                let message = new $root.google.protobuf.Any();
                if (object.type_url != null)
                    message.type_url = String(object.type_url);
                if (object.value != null)
                    if (typeof object.value === "string")
                        $util.base64.decode(object.value, message.value = $util.newBuffer($util.base64.length(object.value)), 0);
                    else if (object.value.length)
                        message.value = object.value;
                return message;
            };

            /**
             * Creates a plain object from an Any message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.Any
             * @static
             * @param {google.protobuf.Any} message Any
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            Any.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                let object = {};
                if (options.defaults) {
                    object.type_url = "";
                    if (options.bytes === String)
                        object.value = "";
                    else {
                        object.value = [];
                        if (options.bytes !== Array)
                            object.value = $util.newBuffer(object.value);
                    }
                }
                if (message.type_url != null && message.hasOwnProperty("type_url"))
                    object.type_url = message.type_url;
                if (message.value != null && message.hasOwnProperty("value"))
                    object.value = options.bytes === String ? $util.base64.encode(message.value, 0, message.value.length) : options.bytes === Array ? Array.prototype.slice.call(message.value) : message.value;
                return object;
            };

            /**
             * Converts this Any to JSON.
             * @function toJSON
             * @memberof google.protobuf.Any
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            Any.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            return Any;
        })();

        return protobuf;
    })();

    return google;
})();

export const MobilePreparedFiles = $root.MobilePreparedFiles = (() => {

    /**
     * Properties of a MobilePreparedFiles.
     * @exports IMobilePreparedFiles
     * @interface IMobilePreparedFiles
     * @property {IDirectory} dir MobilePreparedFiles dir
     * @property {Object.<string,string>} pin MobilePreparedFiles pin
     */

    /**
     * Constructs a new MobilePreparedFiles.
     * @exports MobilePreparedFiles
     * @classdesc Represents a MobilePreparedFiles.
     * @implements IMobilePreparedFiles
     * @constructor
     * @param {IMobilePreparedFiles=} [properties] Properties to set
     */
    function MobilePreparedFiles(properties) {
        this.pin = {};
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * MobilePreparedFiles dir.
     * @member {IDirectory} dir
     * @memberof MobilePreparedFiles
     * @instance
     */
    MobilePreparedFiles.prototype.dir = null;

    /**
     * MobilePreparedFiles pin.
     * @member {Object.<string,string>} pin
     * @memberof MobilePreparedFiles
     * @instance
     */
    MobilePreparedFiles.prototype.pin = $util.emptyObject;

    /**
     * Creates a new MobilePreparedFiles instance using the specified properties.
     * @function create
     * @memberof MobilePreparedFiles
     * @static
     * @param {IMobilePreparedFiles=} [properties] Properties to set
     * @returns {MobilePreparedFiles} MobilePreparedFiles instance
     */
    MobilePreparedFiles.create = function create(properties) {
        return new MobilePreparedFiles(properties);
    };

    /**
     * Encodes the specified MobilePreparedFiles message. Does not implicitly {@link MobilePreparedFiles.verify|verify} messages.
     * @function encode
     * @memberof MobilePreparedFiles
     * @static
     * @param {IMobilePreparedFiles} message MobilePreparedFiles message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    MobilePreparedFiles.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.dir != null && message.hasOwnProperty("dir"))
            $root.Directory.encode(message.dir, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        if (message.pin != null && message.hasOwnProperty("pin"))
            for (let keys = Object.keys(message.pin), i = 0; i < keys.length; ++i)
                writer.uint32(/* id 2, wireType 2 =*/18).fork().uint32(/* id 1, wireType 2 =*/10).string(keys[i]).uint32(/* id 2, wireType 2 =*/18).string(message.pin[keys[i]]).ldelim();
        return writer;
    };

    /**
     * Encodes the specified MobilePreparedFiles message, length delimited. Does not implicitly {@link MobilePreparedFiles.verify|verify} messages.
     * @function encodeDelimited
     * @memberof MobilePreparedFiles
     * @static
     * @param {IMobilePreparedFiles} message MobilePreparedFiles message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    MobilePreparedFiles.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a MobilePreparedFiles message from the specified reader or buffer.
     * @function decode
     * @memberof MobilePreparedFiles
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {MobilePreparedFiles} MobilePreparedFiles
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    MobilePreparedFiles.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.MobilePreparedFiles(), key;
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.dir = $root.Directory.decode(reader, reader.uint32());
                break;
            case 2:
                reader.skip().pos++;
                if (message.pin === $util.emptyObject)
                    message.pin = {};
                key = reader.string();
                reader.pos++;
                message.pin[key] = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a MobilePreparedFiles message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof MobilePreparedFiles
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {MobilePreparedFiles} MobilePreparedFiles
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    MobilePreparedFiles.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a MobilePreparedFiles message.
     * @function verify
     * @memberof MobilePreparedFiles
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    MobilePreparedFiles.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.dir != null && message.hasOwnProperty("dir")) {
            let error = $root.Directory.verify(message.dir);
            if (error)
                return "dir." + error;
        }
        if (message.pin != null && message.hasOwnProperty("pin")) {
            if (!$util.isObject(message.pin))
                return "pin: object expected";
            let key = Object.keys(message.pin);
            for (let i = 0; i < key.length; ++i)
                if (!$util.isString(message.pin[key[i]]))
                    return "pin: string{k:string} expected";
        }
        return null;
    };

    /**
     * Creates a MobilePreparedFiles message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof MobilePreparedFiles
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {MobilePreparedFiles} MobilePreparedFiles
     */
    MobilePreparedFiles.fromObject = function fromObject(object) {
        if (object instanceof $root.MobilePreparedFiles)
            return object;
        let message = new $root.MobilePreparedFiles();
        if (object.dir != null) {
            if (typeof object.dir !== "object")
                throw TypeError(".MobilePreparedFiles.dir: object expected");
            message.dir = $root.Directory.fromObject(object.dir);
        }
        if (object.pin) {
            if (typeof object.pin !== "object")
                throw TypeError(".MobilePreparedFiles.pin: object expected");
            message.pin = {};
            for (let keys = Object.keys(object.pin), i = 0; i < keys.length; ++i)
                message.pin[keys[i]] = String(object.pin[keys[i]]);
        }
        return message;
    };

    /**
     * Creates a plain object from a MobilePreparedFiles message. Also converts values to other types if specified.
     * @function toObject
     * @memberof MobilePreparedFiles
     * @static
     * @param {MobilePreparedFiles} message MobilePreparedFiles
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    MobilePreparedFiles.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.objects || options.defaults)
            object.pin = {};
        if (options.defaults)
            object.dir = null;
        if (message.dir != null && message.hasOwnProperty("dir"))
            object.dir = $root.Directory.toObject(message.dir, options);
        let keys2;
        if (message.pin && (keys2 = Object.keys(message.pin)).length) {
            object.pin = {};
            for (let j = 0; j < keys2.length; ++j)
                object.pin[keys2[j]] = message.pin[keys2[j]];
        }
        return object;
    };

    /**
     * Converts this MobilePreparedFiles to JSON.
     * @function toJSON
     * @memberof MobilePreparedFiles
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    MobilePreparedFiles.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return MobilePreparedFiles;
})();

export const MobileFileData = $root.MobileFileData = (() => {

    /**
     * Properties of a MobileFileData.
     * @exports IMobileFileData
     * @interface IMobileFileData
     * @property {string} url MobileFileData url
     */

    /**
     * Constructs a new MobileFileData.
     * @exports MobileFileData
     * @classdesc Represents a MobileFileData.
     * @implements IMobileFileData
     * @constructor
     * @param {IMobileFileData=} [properties] Properties to set
     */
    function MobileFileData(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * MobileFileData url.
     * @member {string} url
     * @memberof MobileFileData
     * @instance
     */
    MobileFileData.prototype.url = "";

    /**
     * Creates a new MobileFileData instance using the specified properties.
     * @function create
     * @memberof MobileFileData
     * @static
     * @param {IMobileFileData=} [properties] Properties to set
     * @returns {MobileFileData} MobileFileData instance
     */
    MobileFileData.create = function create(properties) {
        return new MobileFileData(properties);
    };

    /**
     * Encodes the specified MobileFileData message. Does not implicitly {@link MobileFileData.verify|verify} messages.
     * @function encode
     * @memberof MobileFileData
     * @static
     * @param {IMobileFileData} message MobileFileData message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    MobileFileData.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.url != null && message.hasOwnProperty("url"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.url);
        return writer;
    };

    /**
     * Encodes the specified MobileFileData message, length delimited. Does not implicitly {@link MobileFileData.verify|verify} messages.
     * @function encodeDelimited
     * @memberof MobileFileData
     * @static
     * @param {IMobileFileData} message MobileFileData message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    MobileFileData.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a MobileFileData message from the specified reader or buffer.
     * @function decode
     * @memberof MobileFileData
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {MobileFileData} MobileFileData
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    MobileFileData.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.MobileFileData();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.url = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a MobileFileData message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof MobileFileData
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {MobileFileData} MobileFileData
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    MobileFileData.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a MobileFileData message.
     * @function verify
     * @memberof MobileFileData
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    MobileFileData.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.url != null && message.hasOwnProperty("url"))
            if (!$util.isString(message.url))
                return "url: string expected";
        return null;
    };

    /**
     * Creates a MobileFileData message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof MobileFileData
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {MobileFileData} MobileFileData
     */
    MobileFileData.fromObject = function fromObject(object) {
        if (object instanceof $root.MobileFileData)
            return object;
        let message = new $root.MobileFileData();
        if (object.url != null)
            message.url = String(object.url);
        return message;
    };

    /**
     * Creates a plain object from a MobileFileData message. Also converts values to other types if specified.
     * @function toObject
     * @memberof MobileFileData
     * @static
     * @param {MobileFileData} message MobileFileData
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    MobileFileData.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults)
            object.url = "";
        if (message.url != null && message.hasOwnProperty("url"))
            object.url = message.url;
        return object;
    };

    /**
     * Converts this MobileFileData to JSON.
     * @function toJSON
     * @memberof MobileFileData
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    MobileFileData.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return MobileFileData;
})();

export const AddThreadConfig = $root.AddThreadConfig = (() => {

    /**
     * Properties of an AddThreadConfig.
     * @exports IAddThreadConfig
     * @interface IAddThreadConfig
     * @property {string} key AddThreadConfig key
     * @property {string} name AddThreadConfig name
     * @property {AddThreadConfig.ISchema} schema AddThreadConfig schema
     * @property {Thread.Type} type AddThreadConfig type
     * @property {Thread.Sharing} sharing AddThreadConfig sharing
     * @property {Array.<string>} members AddThreadConfig members
     */

    /**
     * Constructs a new AddThreadConfig.
     * @exports AddThreadConfig
     * @classdesc Represents an AddThreadConfig.
     * @implements IAddThreadConfig
     * @constructor
     * @param {IAddThreadConfig=} [properties] Properties to set
     */
    function AddThreadConfig(properties) {
        this.members = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * AddThreadConfig key.
     * @member {string} key
     * @memberof AddThreadConfig
     * @instance
     */
    AddThreadConfig.prototype.key = "";

    /**
     * AddThreadConfig name.
     * @member {string} name
     * @memberof AddThreadConfig
     * @instance
     */
    AddThreadConfig.prototype.name = "";

    /**
     * AddThreadConfig schema.
     * @member {AddThreadConfig.ISchema} schema
     * @memberof AddThreadConfig
     * @instance
     */
    AddThreadConfig.prototype.schema = null;

    /**
     * AddThreadConfig type.
     * @member {Thread.Type} type
     * @memberof AddThreadConfig
     * @instance
     */
    AddThreadConfig.prototype.type = 0;

    /**
     * AddThreadConfig sharing.
     * @member {Thread.Sharing} sharing
     * @memberof AddThreadConfig
     * @instance
     */
    AddThreadConfig.prototype.sharing = 0;

    /**
     * AddThreadConfig members.
     * @member {Array.<string>} members
     * @memberof AddThreadConfig
     * @instance
     */
    AddThreadConfig.prototype.members = $util.emptyArray;

    /**
     * Creates a new AddThreadConfig instance using the specified properties.
     * @function create
     * @memberof AddThreadConfig
     * @static
     * @param {IAddThreadConfig=} [properties] Properties to set
     * @returns {AddThreadConfig} AddThreadConfig instance
     */
    AddThreadConfig.create = function create(properties) {
        return new AddThreadConfig(properties);
    };

    /**
     * Encodes the specified AddThreadConfig message. Does not implicitly {@link AddThreadConfig.verify|verify} messages.
     * @function encode
     * @memberof AddThreadConfig
     * @static
     * @param {IAddThreadConfig} message AddThreadConfig message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    AddThreadConfig.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.key != null && message.hasOwnProperty("key"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.key);
        if (message.name != null && message.hasOwnProperty("name"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.name);
        if (message.schema != null && message.hasOwnProperty("schema"))
            $root.AddThreadConfig.Schema.encode(message.schema, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        if (message.type != null && message.hasOwnProperty("type"))
            writer.uint32(/* id 4, wireType 0 =*/32).int32(message.type);
        if (message.sharing != null && message.hasOwnProperty("sharing"))
            writer.uint32(/* id 5, wireType 0 =*/40).int32(message.sharing);
        if (message.members != null && message.members.length)
            for (let i = 0; i < message.members.length; ++i)
                writer.uint32(/* id 6, wireType 2 =*/50).string(message.members[i]);
        return writer;
    };

    /**
     * Encodes the specified AddThreadConfig message, length delimited. Does not implicitly {@link AddThreadConfig.verify|verify} messages.
     * @function encodeDelimited
     * @memberof AddThreadConfig
     * @static
     * @param {IAddThreadConfig} message AddThreadConfig message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    AddThreadConfig.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes an AddThreadConfig message from the specified reader or buffer.
     * @function decode
     * @memberof AddThreadConfig
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {AddThreadConfig} AddThreadConfig
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    AddThreadConfig.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.AddThreadConfig();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.key = reader.string();
                break;
            case 2:
                message.name = reader.string();
                break;
            case 3:
                message.schema = $root.AddThreadConfig.Schema.decode(reader, reader.uint32());
                break;
            case 4:
                message.type = reader.int32();
                break;
            case 5:
                message.sharing = reader.int32();
                break;
            case 6:
                if (!(message.members && message.members.length))
                    message.members = [];
                message.members.push(reader.string());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes an AddThreadConfig message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof AddThreadConfig
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {AddThreadConfig} AddThreadConfig
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    AddThreadConfig.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies an AddThreadConfig message.
     * @function verify
     * @memberof AddThreadConfig
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    AddThreadConfig.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.key != null && message.hasOwnProperty("key"))
            if (!$util.isString(message.key))
                return "key: string expected";
        if (message.name != null && message.hasOwnProperty("name"))
            if (!$util.isString(message.name))
                return "name: string expected";
        if (message.schema != null && message.hasOwnProperty("schema")) {
            let error = $root.AddThreadConfig.Schema.verify(message.schema);
            if (error)
                return "schema." + error;
        }
        if (message.type != null && message.hasOwnProperty("type"))
            switch (message.type) {
            default:
                return "type: enum value expected";
            case 0:
            case 1:
            case 2:
            case 3:
                break;
            }
        if (message.sharing != null && message.hasOwnProperty("sharing"))
            switch (message.sharing) {
            default:
                return "sharing: enum value expected";
            case 0:
            case 1:
            case 2:
                break;
            }
        if (message.members != null && message.hasOwnProperty("members")) {
            if (!Array.isArray(message.members))
                return "members: array expected";
            for (let i = 0; i < message.members.length; ++i)
                if (!$util.isString(message.members[i]))
                    return "members: string[] expected";
        }
        return null;
    };

    /**
     * Creates an AddThreadConfig message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof AddThreadConfig
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {AddThreadConfig} AddThreadConfig
     */
    AddThreadConfig.fromObject = function fromObject(object) {
        if (object instanceof $root.AddThreadConfig)
            return object;
        let message = new $root.AddThreadConfig();
        if (object.key != null)
            message.key = String(object.key);
        if (object.name != null)
            message.name = String(object.name);
        if (object.schema != null) {
            if (typeof object.schema !== "object")
                throw TypeError(".AddThreadConfig.schema: object expected");
            message.schema = $root.AddThreadConfig.Schema.fromObject(object.schema);
        }
        switch (object.type) {
        case "Private":
        case 0:
            message.type = 0;
            break;
        case "ReadOnly":
        case 1:
            message.type = 1;
            break;
        case "Public":
        case 2:
            message.type = 2;
            break;
        case "Open":
        case 3:
            message.type = 3;
            break;
        }
        switch (object.sharing) {
        case "NotShared":
        case 0:
            message.sharing = 0;
            break;
        case "InviteOnly":
        case 1:
            message.sharing = 1;
            break;
        case "Shared":
        case 2:
            message.sharing = 2;
            break;
        }
        if (object.members) {
            if (!Array.isArray(object.members))
                throw TypeError(".AddThreadConfig.members: array expected");
            message.members = [];
            for (let i = 0; i < object.members.length; ++i)
                message.members[i] = String(object.members[i]);
        }
        return message;
    };

    /**
     * Creates a plain object from an AddThreadConfig message. Also converts values to other types if specified.
     * @function toObject
     * @memberof AddThreadConfig
     * @static
     * @param {AddThreadConfig} message AddThreadConfig
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    AddThreadConfig.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.members = [];
        if (options.defaults) {
            object.key = "";
            object.name = "";
            object.schema = null;
            object.type = options.enums === String ? "Private" : 0;
            object.sharing = options.enums === String ? "NotShared" : 0;
        }
        if (message.key != null && message.hasOwnProperty("key"))
            object.key = message.key;
        if (message.name != null && message.hasOwnProperty("name"))
            object.name = message.name;
        if (message.schema != null && message.hasOwnProperty("schema"))
            object.schema = $root.AddThreadConfig.Schema.toObject(message.schema, options);
        if (message.type != null && message.hasOwnProperty("type"))
            object.type = options.enums === String ? $root.Thread.Type[message.type] : message.type;
        if (message.sharing != null && message.hasOwnProperty("sharing"))
            object.sharing = options.enums === String ? $root.Thread.Sharing[message.sharing] : message.sharing;
        if (message.members && message.members.length) {
            object.members = [];
            for (let j = 0; j < message.members.length; ++j)
                object.members[j] = message.members[j];
        }
        return object;
    };

    /**
     * Converts this AddThreadConfig to JSON.
     * @function toJSON
     * @memberof AddThreadConfig
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    AddThreadConfig.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    AddThreadConfig.Schema = (function() {

        /**
         * Properties of a Schema.
         * @memberof AddThreadConfig
         * @interface ISchema
         * @property {string} id Schema id
         * @property {string} json Schema json
         * @property {AddThreadConfig.Schema.Preset} preset Schema preset
         */

        /**
         * Constructs a new Schema.
         * @memberof AddThreadConfig
         * @classdesc Represents a Schema.
         * @implements ISchema
         * @constructor
         * @param {AddThreadConfig.ISchema=} [properties] Properties to set
         */
        function Schema(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * Schema id.
         * @member {string} id
         * @memberof AddThreadConfig.Schema
         * @instance
         */
        Schema.prototype.id = "";

        /**
         * Schema json.
         * @member {string} json
         * @memberof AddThreadConfig.Schema
         * @instance
         */
        Schema.prototype.json = "";

        /**
         * Schema preset.
         * @member {AddThreadConfig.Schema.Preset} preset
         * @memberof AddThreadConfig.Schema
         * @instance
         */
        Schema.prototype.preset = 0;

        /**
         * Creates a new Schema instance using the specified properties.
         * @function create
         * @memberof AddThreadConfig.Schema
         * @static
         * @param {AddThreadConfig.ISchema=} [properties] Properties to set
         * @returns {AddThreadConfig.Schema} Schema instance
         */
        Schema.create = function create(properties) {
            return new Schema(properties);
        };

        /**
         * Encodes the specified Schema message. Does not implicitly {@link AddThreadConfig.Schema.verify|verify} messages.
         * @function encode
         * @memberof AddThreadConfig.Schema
         * @static
         * @param {AddThreadConfig.ISchema} message Schema message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        Schema.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.id != null && message.hasOwnProperty("id"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
            if (message.json != null && message.hasOwnProperty("json"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.json);
            if (message.preset != null && message.hasOwnProperty("preset"))
                writer.uint32(/* id 3, wireType 0 =*/24).int32(message.preset);
            return writer;
        };

        /**
         * Encodes the specified Schema message, length delimited. Does not implicitly {@link AddThreadConfig.Schema.verify|verify} messages.
         * @function encodeDelimited
         * @memberof AddThreadConfig.Schema
         * @static
         * @param {AddThreadConfig.ISchema} message Schema message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        Schema.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a Schema message from the specified reader or buffer.
         * @function decode
         * @memberof AddThreadConfig.Schema
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {AddThreadConfig.Schema} Schema
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        Schema.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.AddThreadConfig.Schema();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.id = reader.string();
                    break;
                case 2:
                    message.json = reader.string();
                    break;
                case 3:
                    message.preset = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a Schema message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof AddThreadConfig.Schema
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {AddThreadConfig.Schema} Schema
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        Schema.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a Schema message.
         * @function verify
         * @memberof AddThreadConfig.Schema
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        Schema.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.id != null && message.hasOwnProperty("id"))
                if (!$util.isString(message.id))
                    return "id: string expected";
            if (message.json != null && message.hasOwnProperty("json"))
                if (!$util.isString(message.json))
                    return "json: string expected";
            if (message.preset != null && message.hasOwnProperty("preset"))
                switch (message.preset) {
                default:
                    return "preset: enum value expected";
                case 0:
                case 1:
                case 2:
                    break;
                }
            return null;
        };

        /**
         * Creates a Schema message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof AddThreadConfig.Schema
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {AddThreadConfig.Schema} Schema
         */
        Schema.fromObject = function fromObject(object) {
            if (object instanceof $root.AddThreadConfig.Schema)
                return object;
            let message = new $root.AddThreadConfig.Schema();
            if (object.id != null)
                message.id = String(object.id);
            if (object.json != null)
                message.json = String(object.json);
            switch (object.preset) {
            case "NONE":
            case 0:
                message.preset = 0;
                break;
            case "CAMERA_ROLL":
            case 1:
                message.preset = 1;
                break;
            case "MEDIA":
            case 2:
                message.preset = 2;
                break;
            }
            return message;
        };

        /**
         * Creates a plain object from a Schema message. Also converts values to other types if specified.
         * @function toObject
         * @memberof AddThreadConfig.Schema
         * @static
         * @param {AddThreadConfig.Schema} message Schema
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        Schema.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.id = "";
                object.json = "";
                object.preset = options.enums === String ? "NONE" : 0;
            }
            if (message.id != null && message.hasOwnProperty("id"))
                object.id = message.id;
            if (message.json != null && message.hasOwnProperty("json"))
                object.json = message.json;
            if (message.preset != null && message.hasOwnProperty("preset"))
                object.preset = options.enums === String ? $root.AddThreadConfig.Schema.Preset[message.preset] : message.preset;
            return object;
        };

        /**
         * Converts this Schema to JSON.
         * @function toJSON
         * @memberof AddThreadConfig.Schema
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        Schema.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        /**
         * Preset enum.
         * @name AddThreadConfig.Schema.Preset
         * @enum {string}
         * @property {number} NONE=0 NONE value
         * @property {number} CAMERA_ROLL=1 CAMERA_ROLL value
         * @property {number} MEDIA=2 MEDIA value
         */
        Schema.Preset = (function() {
            const valuesById = {}, values = Object.create(valuesById);
            values[valuesById[0] = "NONE"] = 0;
            values[valuesById[1] = "CAMERA_ROLL"] = 1;
            values[valuesById[2] = "MEDIA"] = 2;
            return values;
        })();

        return Schema;
    })();

    return AddThreadConfig;
})();

export const Step = $root.Step = (() => {

    /**
     * Properties of a Step.
     * @exports IStep
     * @interface IStep
     * @property {string} Name Step Name
     * @property {ILink} link Step link
     */

    /**
     * Constructs a new Step.
     * @exports Step
     * @classdesc Represents a Step.
     * @implements IStep
     * @constructor
     * @param {IStep=} [properties] Properties to set
     */
    function Step(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Step Name.
     * @member {string} Name
     * @memberof Step
     * @instance
     */
    Step.prototype.Name = "";

    /**
     * Step link.
     * @member {ILink} link
     * @memberof Step
     * @instance
     */
    Step.prototype.link = null;

    /**
     * Creates a new Step instance using the specified properties.
     * @function create
     * @memberof Step
     * @static
     * @param {IStep=} [properties] Properties to set
     * @returns {Step} Step instance
     */
    Step.create = function create(properties) {
        return new Step(properties);
    };

    /**
     * Encodes the specified Step message. Does not implicitly {@link Step.verify|verify} messages.
     * @function encode
     * @memberof Step
     * @static
     * @param {IStep} message Step message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Step.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.Name != null && message.hasOwnProperty("Name"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.Name);
        if (message.link != null && message.hasOwnProperty("link"))
            $root.Link.encode(message.link, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified Step message, length delimited. Does not implicitly {@link Step.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Step
     * @static
     * @param {IStep} message Step message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Step.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Step message from the specified reader or buffer.
     * @function decode
     * @memberof Step
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Step} Step
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Step.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Step();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.Name = reader.string();
                break;
            case 2:
                message.link = $root.Link.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Step message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Step
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Step} Step
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Step.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Step message.
     * @function verify
     * @memberof Step
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Step.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.Name != null && message.hasOwnProperty("Name"))
            if (!$util.isString(message.Name))
                return "Name: string expected";
        if (message.link != null && message.hasOwnProperty("link")) {
            let error = $root.Link.verify(message.link);
            if (error)
                return "link." + error;
        }
        return null;
    };

    /**
     * Creates a Step message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Step
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Step} Step
     */
    Step.fromObject = function fromObject(object) {
        if (object instanceof $root.Step)
            return object;
        let message = new $root.Step();
        if (object.Name != null)
            message.Name = String(object.Name);
        if (object.link != null) {
            if (typeof object.link !== "object")
                throw TypeError(".Step.link: object expected");
            message.link = $root.Link.fromObject(object.link);
        }
        return message;
    };

    /**
     * Creates a plain object from a Step message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Step
     * @static
     * @param {Step} message Step
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Step.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.Name = "";
            object.link = null;
        }
        if (message.Name != null && message.hasOwnProperty("Name"))
            object.Name = message.Name;
        if (message.link != null && message.hasOwnProperty("link"))
            object.link = $root.Link.toObject(message.link, options);
        return object;
    };

    /**
     * Converts this Step to JSON.
     * @function toJSON
     * @memberof Step
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Step.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Step;
})();

export const Directory = $root.Directory = (() => {

    /**
     * Properties of a Directory.
     * @exports IDirectory
     * @interface IDirectory
     * @property {Object.<string,IFileIndex>} files Directory files
     */

    /**
     * Constructs a new Directory.
     * @exports Directory
     * @classdesc Represents a Directory.
     * @implements IDirectory
     * @constructor
     * @param {IDirectory=} [properties] Properties to set
     */
    function Directory(properties) {
        this.files = {};
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Directory files.
     * @member {Object.<string,IFileIndex>} files
     * @memberof Directory
     * @instance
     */
    Directory.prototype.files = $util.emptyObject;

    /**
     * Creates a new Directory instance using the specified properties.
     * @function create
     * @memberof Directory
     * @static
     * @param {IDirectory=} [properties] Properties to set
     * @returns {Directory} Directory instance
     */
    Directory.create = function create(properties) {
        return new Directory(properties);
    };

    /**
     * Encodes the specified Directory message. Does not implicitly {@link Directory.verify|verify} messages.
     * @function encode
     * @memberof Directory
     * @static
     * @param {IDirectory} message Directory message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Directory.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.files != null && message.hasOwnProperty("files"))
            for (let keys = Object.keys(message.files), i = 0; i < keys.length; ++i) {
                writer.uint32(/* id 1, wireType 2 =*/10).fork().uint32(/* id 1, wireType 2 =*/10).string(keys[i]);
                $root.FileIndex.encode(message.files[keys[i]], writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim().ldelim();
            }
        return writer;
    };

    /**
     * Encodes the specified Directory message, length delimited. Does not implicitly {@link Directory.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Directory
     * @static
     * @param {IDirectory} message Directory message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Directory.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Directory message from the specified reader or buffer.
     * @function decode
     * @memberof Directory
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Directory} Directory
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Directory.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Directory(), key;
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                reader.skip().pos++;
                if (message.files === $util.emptyObject)
                    message.files = {};
                key = reader.string();
                reader.pos++;
                message.files[key] = $root.FileIndex.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Directory message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Directory
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Directory} Directory
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Directory.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Directory message.
     * @function verify
     * @memberof Directory
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Directory.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.files != null && message.hasOwnProperty("files")) {
            if (!$util.isObject(message.files))
                return "files: object expected";
            let key = Object.keys(message.files);
            for (let i = 0; i < key.length; ++i) {
                let error = $root.FileIndex.verify(message.files[key[i]]);
                if (error)
                    return "files." + error;
            }
        }
        return null;
    };

    /**
     * Creates a Directory message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Directory
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Directory} Directory
     */
    Directory.fromObject = function fromObject(object) {
        if (object instanceof $root.Directory)
            return object;
        let message = new $root.Directory();
        if (object.files) {
            if (typeof object.files !== "object")
                throw TypeError(".Directory.files: object expected");
            message.files = {};
            for (let keys = Object.keys(object.files), i = 0; i < keys.length; ++i) {
                if (typeof object.files[keys[i]] !== "object")
                    throw TypeError(".Directory.files: object expected");
                message.files[keys[i]] = $root.FileIndex.fromObject(object.files[keys[i]]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from a Directory message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Directory
     * @static
     * @param {Directory} message Directory
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Directory.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.objects || options.defaults)
            object.files = {};
        let keys2;
        if (message.files && (keys2 = Object.keys(message.files)).length) {
            object.files = {};
            for (let j = 0; j < keys2.length; ++j)
                object.files[keys2[j]] = $root.FileIndex.toObject(message.files[keys2[j]], options);
        }
        return object;
    };

    /**
     * Converts this Directory to JSON.
     * @function toJSON
     * @memberof Directory
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Directory.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Directory;
})();

export const DirectoryList = $root.DirectoryList = (() => {

    /**
     * Properties of a DirectoryList.
     * @exports IDirectoryList
     * @interface IDirectoryList
     * @property {Array.<IDirectory>} items DirectoryList items
     */

    /**
     * Constructs a new DirectoryList.
     * @exports DirectoryList
     * @classdesc Represents a DirectoryList.
     * @implements IDirectoryList
     * @constructor
     * @param {IDirectoryList=} [properties] Properties to set
     */
    function DirectoryList(properties) {
        this.items = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * DirectoryList items.
     * @member {Array.<IDirectory>} items
     * @memberof DirectoryList
     * @instance
     */
    DirectoryList.prototype.items = $util.emptyArray;

    /**
     * Creates a new DirectoryList instance using the specified properties.
     * @function create
     * @memberof DirectoryList
     * @static
     * @param {IDirectoryList=} [properties] Properties to set
     * @returns {DirectoryList} DirectoryList instance
     */
    DirectoryList.create = function create(properties) {
        return new DirectoryList(properties);
    };

    /**
     * Encodes the specified DirectoryList message. Does not implicitly {@link DirectoryList.verify|verify} messages.
     * @function encode
     * @memberof DirectoryList
     * @static
     * @param {IDirectoryList} message DirectoryList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    DirectoryList.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.items != null && message.items.length)
            for (let i = 0; i < message.items.length; ++i)
                $root.Directory.encode(message.items[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified DirectoryList message, length delimited. Does not implicitly {@link DirectoryList.verify|verify} messages.
     * @function encodeDelimited
     * @memberof DirectoryList
     * @static
     * @param {IDirectoryList} message DirectoryList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    DirectoryList.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a DirectoryList message from the specified reader or buffer.
     * @function decode
     * @memberof DirectoryList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {DirectoryList} DirectoryList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    DirectoryList.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.DirectoryList();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                if (!(message.items && message.items.length))
                    message.items = [];
                message.items.push($root.Directory.decode(reader, reader.uint32()));
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a DirectoryList message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof DirectoryList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {DirectoryList} DirectoryList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    DirectoryList.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a DirectoryList message.
     * @function verify
     * @memberof DirectoryList
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    DirectoryList.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.items != null && message.hasOwnProperty("items")) {
            if (!Array.isArray(message.items))
                return "items: array expected";
            for (let i = 0; i < message.items.length; ++i) {
                let error = $root.Directory.verify(message.items[i]);
                if (error)
                    return "items." + error;
            }
        }
        return null;
    };

    /**
     * Creates a DirectoryList message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof DirectoryList
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {DirectoryList} DirectoryList
     */
    DirectoryList.fromObject = function fromObject(object) {
        if (object instanceof $root.DirectoryList)
            return object;
        let message = new $root.DirectoryList();
        if (object.items) {
            if (!Array.isArray(object.items))
                throw TypeError(".DirectoryList.items: array expected");
            message.items = [];
            for (let i = 0; i < object.items.length; ++i) {
                if (typeof object.items[i] !== "object")
                    throw TypeError(".DirectoryList.items: object expected");
                message.items[i] = $root.Directory.fromObject(object.items[i]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from a DirectoryList message. Also converts values to other types if specified.
     * @function toObject
     * @memberof DirectoryList
     * @static
     * @param {DirectoryList} message DirectoryList
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    DirectoryList.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.items = [];
        if (message.items && message.items.length) {
            object.items = [];
            for (let j = 0; j < message.items.length; ++j)
                object.items[j] = $root.Directory.toObject(message.items[j], options);
        }
        return object;
    };

    /**
     * Converts this DirectoryList to JSON.
     * @function toJSON
     * @memberof DirectoryList
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    DirectoryList.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return DirectoryList;
})();

export const Keys = $root.Keys = (() => {

    /**
     * Properties of a Keys.
     * @exports IKeys
     * @interface IKeys
     * @property {Object.<string,string>} files Keys files
     */

    /**
     * Constructs a new Keys.
     * @exports Keys
     * @classdesc Represents a Keys.
     * @implements IKeys
     * @constructor
     * @param {IKeys=} [properties] Properties to set
     */
    function Keys(properties) {
        this.files = {};
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Keys files.
     * @member {Object.<string,string>} files
     * @memberof Keys
     * @instance
     */
    Keys.prototype.files = $util.emptyObject;

    /**
     * Creates a new Keys instance using the specified properties.
     * @function create
     * @memberof Keys
     * @static
     * @param {IKeys=} [properties] Properties to set
     * @returns {Keys} Keys instance
     */
    Keys.create = function create(properties) {
        return new Keys(properties);
    };

    /**
     * Encodes the specified Keys message. Does not implicitly {@link Keys.verify|verify} messages.
     * @function encode
     * @memberof Keys
     * @static
     * @param {IKeys} message Keys message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Keys.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.files != null && message.hasOwnProperty("files"))
            for (let keys = Object.keys(message.files), i = 0; i < keys.length; ++i)
                writer.uint32(/* id 1, wireType 2 =*/10).fork().uint32(/* id 1, wireType 2 =*/10).string(keys[i]).uint32(/* id 2, wireType 2 =*/18).string(message.files[keys[i]]).ldelim();
        return writer;
    };

    /**
     * Encodes the specified Keys message, length delimited. Does not implicitly {@link Keys.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Keys
     * @static
     * @param {IKeys} message Keys message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Keys.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Keys message from the specified reader or buffer.
     * @function decode
     * @memberof Keys
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Keys} Keys
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Keys.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Keys(), key;
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                reader.skip().pos++;
                if (message.files === $util.emptyObject)
                    message.files = {};
                key = reader.string();
                reader.pos++;
                message.files[key] = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Keys message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Keys
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Keys} Keys
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Keys.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Keys message.
     * @function verify
     * @memberof Keys
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Keys.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.files != null && message.hasOwnProperty("files")) {
            if (!$util.isObject(message.files))
                return "files: object expected";
            let key = Object.keys(message.files);
            for (let i = 0; i < key.length; ++i)
                if (!$util.isString(message.files[key[i]]))
                    return "files: string{k:string} expected";
        }
        return null;
    };

    /**
     * Creates a Keys message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Keys
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Keys} Keys
     */
    Keys.fromObject = function fromObject(object) {
        if (object instanceof $root.Keys)
            return object;
        let message = new $root.Keys();
        if (object.files) {
            if (typeof object.files !== "object")
                throw TypeError(".Keys.files: object expected");
            message.files = {};
            for (let keys = Object.keys(object.files), i = 0; i < keys.length; ++i)
                message.files[keys[i]] = String(object.files[keys[i]]);
        }
        return message;
    };

    /**
     * Creates a plain object from a Keys message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Keys
     * @static
     * @param {Keys} message Keys
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Keys.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.objects || options.defaults)
            object.files = {};
        let keys2;
        if (message.files && (keys2 = Object.keys(message.files)).length) {
            object.files = {};
            for (let j = 0; j < keys2.length; ++j)
                object.files[keys2[j]] = message.files[keys2[j]];
        }
        return object;
    };

    /**
     * Converts this Keys to JSON.
     * @function toJSON
     * @memberof Keys
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Keys.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Keys;
})();

export const NewInvite = $root.NewInvite = (() => {

    /**
     * Properties of a NewInvite.
     * @exports INewInvite
     * @interface INewInvite
     * @property {string} id NewInvite id
     * @property {string} key NewInvite key
     * @property {string} inviter NewInvite inviter
     */

    /**
     * Constructs a new NewInvite.
     * @exports NewInvite
     * @classdesc Represents a NewInvite.
     * @implements INewInvite
     * @constructor
     * @param {INewInvite=} [properties] Properties to set
     */
    function NewInvite(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * NewInvite id.
     * @member {string} id
     * @memberof NewInvite
     * @instance
     */
    NewInvite.prototype.id = "";

    /**
     * NewInvite key.
     * @member {string} key
     * @memberof NewInvite
     * @instance
     */
    NewInvite.prototype.key = "";

    /**
     * NewInvite inviter.
     * @member {string} inviter
     * @memberof NewInvite
     * @instance
     */
    NewInvite.prototype.inviter = "";

    /**
     * Creates a new NewInvite instance using the specified properties.
     * @function create
     * @memberof NewInvite
     * @static
     * @param {INewInvite=} [properties] Properties to set
     * @returns {NewInvite} NewInvite instance
     */
    NewInvite.create = function create(properties) {
        return new NewInvite(properties);
    };

    /**
     * Encodes the specified NewInvite message. Does not implicitly {@link NewInvite.verify|verify} messages.
     * @function encode
     * @memberof NewInvite
     * @static
     * @param {INewInvite} message NewInvite message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    NewInvite.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.key != null && message.hasOwnProperty("key"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.key);
        if (message.inviter != null && message.hasOwnProperty("inviter"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.inviter);
        return writer;
    };

    /**
     * Encodes the specified NewInvite message, length delimited. Does not implicitly {@link NewInvite.verify|verify} messages.
     * @function encodeDelimited
     * @memberof NewInvite
     * @static
     * @param {INewInvite} message NewInvite message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    NewInvite.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a NewInvite message from the specified reader or buffer.
     * @function decode
     * @memberof NewInvite
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {NewInvite} NewInvite
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    NewInvite.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.NewInvite();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.key = reader.string();
                break;
            case 3:
                message.inviter = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a NewInvite message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof NewInvite
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {NewInvite} NewInvite
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    NewInvite.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a NewInvite message.
     * @function verify
     * @memberof NewInvite
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    NewInvite.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.key != null && message.hasOwnProperty("key"))
            if (!$util.isString(message.key))
                return "key: string expected";
        if (message.inviter != null && message.hasOwnProperty("inviter"))
            if (!$util.isString(message.inviter))
                return "inviter: string expected";
        return null;
    };

    /**
     * Creates a NewInvite message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof NewInvite
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {NewInvite} NewInvite
     */
    NewInvite.fromObject = function fromObject(object) {
        if (object instanceof $root.NewInvite)
            return object;
        let message = new $root.NewInvite();
        if (object.id != null)
            message.id = String(object.id);
        if (object.key != null)
            message.key = String(object.key);
        if (object.inviter != null)
            message.inviter = String(object.inviter);
        return message;
    };

    /**
     * Creates a plain object from a NewInvite message. Also converts values to other types if specified.
     * @function toObject
     * @memberof NewInvite
     * @static
     * @param {NewInvite} message NewInvite
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    NewInvite.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            object.key = "";
            object.inviter = "";
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.key != null && message.hasOwnProperty("key"))
            object.key = message.key;
        if (message.inviter != null && message.hasOwnProperty("inviter"))
            object.inviter = message.inviter;
        return object;
    };

    /**
     * Converts this NewInvite to JSON.
     * @function toJSON
     * @memberof NewInvite
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    NewInvite.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return NewInvite;
})();

export const InviteView = $root.InviteView = (() => {

    /**
     * Properties of an InviteView.
     * @exports IInviteView
     * @interface IInviteView
     * @property {string} id InviteView id
     * @property {string} name InviteView name
     * @property {IUser} inviter InviteView inviter
     * @property {google.protobuf.ITimestamp} date InviteView date
     */

    /**
     * Constructs a new InviteView.
     * @exports InviteView
     * @classdesc Represents an InviteView.
     * @implements IInviteView
     * @constructor
     * @param {IInviteView=} [properties] Properties to set
     */
    function InviteView(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * InviteView id.
     * @member {string} id
     * @memberof InviteView
     * @instance
     */
    InviteView.prototype.id = "";

    /**
     * InviteView name.
     * @member {string} name
     * @memberof InviteView
     * @instance
     */
    InviteView.prototype.name = "";

    /**
     * InviteView inviter.
     * @member {IUser} inviter
     * @memberof InviteView
     * @instance
     */
    InviteView.prototype.inviter = null;

    /**
     * InviteView date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof InviteView
     * @instance
     */
    InviteView.prototype.date = null;

    /**
     * Creates a new InviteView instance using the specified properties.
     * @function create
     * @memberof InviteView
     * @static
     * @param {IInviteView=} [properties] Properties to set
     * @returns {InviteView} InviteView instance
     */
    InviteView.create = function create(properties) {
        return new InviteView(properties);
    };

    /**
     * Encodes the specified InviteView message. Does not implicitly {@link InviteView.verify|verify} messages.
     * @function encode
     * @memberof InviteView
     * @static
     * @param {IInviteView} message InviteView message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    InviteView.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.name != null && message.hasOwnProperty("name"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.name);
        if (message.inviter != null && message.hasOwnProperty("inviter"))
            $root.User.encode(message.inviter, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified InviteView message, length delimited. Does not implicitly {@link InviteView.verify|verify} messages.
     * @function encodeDelimited
     * @memberof InviteView
     * @static
     * @param {IInviteView} message InviteView message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    InviteView.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes an InviteView message from the specified reader or buffer.
     * @function decode
     * @memberof InviteView
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {InviteView} InviteView
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    InviteView.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.InviteView();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.name = reader.string();
                break;
            case 3:
                message.inviter = $root.User.decode(reader, reader.uint32());
                break;
            case 4:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes an InviteView message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof InviteView
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {InviteView} InviteView
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    InviteView.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies an InviteView message.
     * @function verify
     * @memberof InviteView
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    InviteView.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.name != null && message.hasOwnProperty("name"))
            if (!$util.isString(message.name))
                return "name: string expected";
        if (message.inviter != null && message.hasOwnProperty("inviter")) {
            let error = $root.User.verify(message.inviter);
            if (error)
                return "inviter." + error;
        }
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        return null;
    };

    /**
     * Creates an InviteView message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof InviteView
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {InviteView} InviteView
     */
    InviteView.fromObject = function fromObject(object) {
        if (object instanceof $root.InviteView)
            return object;
        let message = new $root.InviteView();
        if (object.id != null)
            message.id = String(object.id);
        if (object.name != null)
            message.name = String(object.name);
        if (object.inviter != null) {
            if (typeof object.inviter !== "object")
                throw TypeError(".InviteView.inviter: object expected");
            message.inviter = $root.User.fromObject(object.inviter);
        }
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".InviteView.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        return message;
    };

    /**
     * Creates a plain object from an InviteView message. Also converts values to other types if specified.
     * @function toObject
     * @memberof InviteView
     * @static
     * @param {InviteView} message InviteView
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    InviteView.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            object.name = "";
            object.inviter = null;
            object.date = null;
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.name != null && message.hasOwnProperty("name"))
            object.name = message.name;
        if (message.inviter != null && message.hasOwnProperty("inviter"))
            object.inviter = $root.User.toObject(message.inviter, options);
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        return object;
    };

    /**
     * Converts this InviteView to JSON.
     * @function toJSON
     * @memberof InviteView
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    InviteView.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return InviteView;
})();

export const InviteViewList = $root.InviteViewList = (() => {

    /**
     * Properties of an InviteViewList.
     * @exports IInviteViewList
     * @interface IInviteViewList
     * @property {Array.<IInviteView>} items InviteViewList items
     */

    /**
     * Constructs a new InviteViewList.
     * @exports InviteViewList
     * @classdesc Represents an InviteViewList.
     * @implements IInviteViewList
     * @constructor
     * @param {IInviteViewList=} [properties] Properties to set
     */
    function InviteViewList(properties) {
        this.items = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * InviteViewList items.
     * @member {Array.<IInviteView>} items
     * @memberof InviteViewList
     * @instance
     */
    InviteViewList.prototype.items = $util.emptyArray;

    /**
     * Creates a new InviteViewList instance using the specified properties.
     * @function create
     * @memberof InviteViewList
     * @static
     * @param {IInviteViewList=} [properties] Properties to set
     * @returns {InviteViewList} InviteViewList instance
     */
    InviteViewList.create = function create(properties) {
        return new InviteViewList(properties);
    };

    /**
     * Encodes the specified InviteViewList message. Does not implicitly {@link InviteViewList.verify|verify} messages.
     * @function encode
     * @memberof InviteViewList
     * @static
     * @param {IInviteViewList} message InviteViewList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    InviteViewList.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.items != null && message.items.length)
            for (let i = 0; i < message.items.length; ++i)
                $root.InviteView.encode(message.items[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified InviteViewList message, length delimited. Does not implicitly {@link InviteViewList.verify|verify} messages.
     * @function encodeDelimited
     * @memberof InviteViewList
     * @static
     * @param {IInviteViewList} message InviteViewList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    InviteViewList.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes an InviteViewList message from the specified reader or buffer.
     * @function decode
     * @memberof InviteViewList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {InviteViewList} InviteViewList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    InviteViewList.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.InviteViewList();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                if (!(message.items && message.items.length))
                    message.items = [];
                message.items.push($root.InviteView.decode(reader, reader.uint32()));
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes an InviteViewList message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof InviteViewList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {InviteViewList} InviteViewList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    InviteViewList.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies an InviteViewList message.
     * @function verify
     * @memberof InviteViewList
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    InviteViewList.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.items != null && message.hasOwnProperty("items")) {
            if (!Array.isArray(message.items))
                return "items: array expected";
            for (let i = 0; i < message.items.length; ++i) {
                let error = $root.InviteView.verify(message.items[i]);
                if (error)
                    return "items." + error;
            }
        }
        return null;
    };

    /**
     * Creates an InviteViewList message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof InviteViewList
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {InviteViewList} InviteViewList
     */
    InviteViewList.fromObject = function fromObject(object) {
        if (object instanceof $root.InviteViewList)
            return object;
        let message = new $root.InviteViewList();
        if (object.items) {
            if (!Array.isArray(object.items))
                throw TypeError(".InviteViewList.items: array expected");
            message.items = [];
            for (let i = 0; i < object.items.length; ++i) {
                if (typeof object.items[i] !== "object")
                    throw TypeError(".InviteViewList.items: object expected");
                message.items[i] = $root.InviteView.fromObject(object.items[i]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from an InviteViewList message. Also converts values to other types if specified.
     * @function toObject
     * @memberof InviteViewList
     * @static
     * @param {InviteViewList} message InviteViewList
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    InviteViewList.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.items = [];
        if (message.items && message.items.length) {
            object.items = [];
            for (let j = 0; j < message.items.length; ++j)
                object.items[j] = $root.InviteView.toObject(message.items[j], options);
        }
        return object;
    };

    /**
     * Converts this InviteViewList to JSON.
     * @function toJSON
     * @memberof InviteViewList
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    InviteViewList.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return InviteViewList;
})();

/**
 * FeedMode enum.
 * @exports FeedMode
 * @enum {string}
 * @property {number} CHRONO=0 CHRONO value
 * @property {number} ANNOTATED=1 ANNOTATED value
 * @property {number} STACKS=2 STACKS value
 */
$root.FeedMode = (function() {
    const valuesById = {}, values = Object.create(valuesById);
    values[valuesById[0] = "CHRONO"] = 0;
    values[valuesById[1] = "ANNOTATED"] = 1;
    values[valuesById[2] = "STACKS"] = 2;
    return values;
})();

export const FeedItem = $root.FeedItem = (() => {

    /**
     * Properties of a FeedItem.
     * @exports IFeedItem
     * @interface IFeedItem
     * @property {string} block FeedItem block
     * @property {string} thread FeedItem thread
     * @property {google.protobuf.IAny} payload FeedItem payload
     */

    /**
     * Constructs a new FeedItem.
     * @exports FeedItem
     * @classdesc Represents a FeedItem.
     * @implements IFeedItem
     * @constructor
     * @param {IFeedItem=} [properties] Properties to set
     */
    function FeedItem(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * FeedItem block.
     * @member {string} block
     * @memberof FeedItem
     * @instance
     */
    FeedItem.prototype.block = "";

    /**
     * FeedItem thread.
     * @member {string} thread
     * @memberof FeedItem
     * @instance
     */
    FeedItem.prototype.thread = "";

    /**
     * FeedItem payload.
     * @member {google.protobuf.IAny} payload
     * @memberof FeedItem
     * @instance
     */
    FeedItem.prototype.payload = null;

    /**
     * Creates a new FeedItem instance using the specified properties.
     * @function create
     * @memberof FeedItem
     * @static
     * @param {IFeedItem=} [properties] Properties to set
     * @returns {FeedItem} FeedItem instance
     */
    FeedItem.create = function create(properties) {
        return new FeedItem(properties);
    };

    /**
     * Encodes the specified FeedItem message. Does not implicitly {@link FeedItem.verify|verify} messages.
     * @function encode
     * @memberof FeedItem
     * @static
     * @param {IFeedItem} message FeedItem message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    FeedItem.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.block != null && message.hasOwnProperty("block"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.block);
        if (message.thread != null && message.hasOwnProperty("thread"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.thread);
        if (message.payload != null && message.hasOwnProperty("payload"))
            $root.google.protobuf.Any.encode(message.payload, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified FeedItem message, length delimited. Does not implicitly {@link FeedItem.verify|verify} messages.
     * @function encodeDelimited
     * @memberof FeedItem
     * @static
     * @param {IFeedItem} message FeedItem message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    FeedItem.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a FeedItem message from the specified reader or buffer.
     * @function decode
     * @memberof FeedItem
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {FeedItem} FeedItem
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    FeedItem.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.FeedItem();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.block = reader.string();
                break;
            case 2:
                message.thread = reader.string();
                break;
            case 3:
                message.payload = $root.google.protobuf.Any.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a FeedItem message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof FeedItem
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {FeedItem} FeedItem
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    FeedItem.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a FeedItem message.
     * @function verify
     * @memberof FeedItem
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    FeedItem.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.block != null && message.hasOwnProperty("block"))
            if (!$util.isString(message.block))
                return "block: string expected";
        if (message.thread != null && message.hasOwnProperty("thread"))
            if (!$util.isString(message.thread))
                return "thread: string expected";
        if (message.payload != null && message.hasOwnProperty("payload")) {
            let error = $root.google.protobuf.Any.verify(message.payload);
            if (error)
                return "payload." + error;
        }
        return null;
    };

    /**
     * Creates a FeedItem message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof FeedItem
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {FeedItem} FeedItem
     */
    FeedItem.fromObject = function fromObject(object) {
        if (object instanceof $root.FeedItem)
            return object;
        let message = new $root.FeedItem();
        if (object.block != null)
            message.block = String(object.block);
        if (object.thread != null)
            message.thread = String(object.thread);
        if (object.payload != null) {
            if (typeof object.payload !== "object")
                throw TypeError(".FeedItem.payload: object expected");
            message.payload = $root.google.protobuf.Any.fromObject(object.payload);
        }
        return message;
    };

    /**
     * Creates a plain object from a FeedItem message. Also converts values to other types if specified.
     * @function toObject
     * @memberof FeedItem
     * @static
     * @param {FeedItem} message FeedItem
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    FeedItem.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.block = "";
            object.thread = "";
            object.payload = null;
        }
        if (message.block != null && message.hasOwnProperty("block"))
            object.block = message.block;
        if (message.thread != null && message.hasOwnProperty("thread"))
            object.thread = message.thread;
        if (message.payload != null && message.hasOwnProperty("payload"))
            object.payload = $root.google.protobuf.Any.toObject(message.payload, options);
        return object;
    };

    /**
     * Converts this FeedItem to JSON.
     * @function toJSON
     * @memberof FeedItem
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    FeedItem.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return FeedItem;
})();

export const FeedItemList = $root.FeedItemList = (() => {

    /**
     * Properties of a FeedItemList.
     * @exports IFeedItemList
     * @interface IFeedItemList
     * @property {Array.<IFeedItem>} items FeedItemList items
     * @property {number} count FeedItemList count
     * @property {string} next FeedItemList next
     */

    /**
     * Constructs a new FeedItemList.
     * @exports FeedItemList
     * @classdesc Represents a FeedItemList.
     * @implements IFeedItemList
     * @constructor
     * @param {IFeedItemList=} [properties] Properties to set
     */
    function FeedItemList(properties) {
        this.items = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * FeedItemList items.
     * @member {Array.<IFeedItem>} items
     * @memberof FeedItemList
     * @instance
     */
    FeedItemList.prototype.items = $util.emptyArray;

    /**
     * FeedItemList count.
     * @member {number} count
     * @memberof FeedItemList
     * @instance
     */
    FeedItemList.prototype.count = 0;

    /**
     * FeedItemList next.
     * @member {string} next
     * @memberof FeedItemList
     * @instance
     */
    FeedItemList.prototype.next = "";

    /**
     * Creates a new FeedItemList instance using the specified properties.
     * @function create
     * @memberof FeedItemList
     * @static
     * @param {IFeedItemList=} [properties] Properties to set
     * @returns {FeedItemList} FeedItemList instance
     */
    FeedItemList.create = function create(properties) {
        return new FeedItemList(properties);
    };

    /**
     * Encodes the specified FeedItemList message. Does not implicitly {@link FeedItemList.verify|verify} messages.
     * @function encode
     * @memberof FeedItemList
     * @static
     * @param {IFeedItemList} message FeedItemList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    FeedItemList.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.items != null && message.items.length)
            for (let i = 0; i < message.items.length; ++i)
                $root.FeedItem.encode(message.items[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        if (message.count != null && message.hasOwnProperty("count"))
            writer.uint32(/* id 2, wireType 0 =*/16).int32(message.count);
        if (message.next != null && message.hasOwnProperty("next"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.next);
        return writer;
    };

    /**
     * Encodes the specified FeedItemList message, length delimited. Does not implicitly {@link FeedItemList.verify|verify} messages.
     * @function encodeDelimited
     * @memberof FeedItemList
     * @static
     * @param {IFeedItemList} message FeedItemList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    FeedItemList.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a FeedItemList message from the specified reader or buffer.
     * @function decode
     * @memberof FeedItemList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {FeedItemList} FeedItemList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    FeedItemList.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.FeedItemList();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                if (!(message.items && message.items.length))
                    message.items = [];
                message.items.push($root.FeedItem.decode(reader, reader.uint32()));
                break;
            case 2:
                message.count = reader.int32();
                break;
            case 3:
                message.next = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a FeedItemList message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof FeedItemList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {FeedItemList} FeedItemList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    FeedItemList.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a FeedItemList message.
     * @function verify
     * @memberof FeedItemList
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    FeedItemList.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.items != null && message.hasOwnProperty("items")) {
            if (!Array.isArray(message.items))
                return "items: array expected";
            for (let i = 0; i < message.items.length; ++i) {
                let error = $root.FeedItem.verify(message.items[i]);
                if (error)
                    return "items." + error;
            }
        }
        if (message.count != null && message.hasOwnProperty("count"))
            if (!$util.isInteger(message.count))
                return "count: integer expected";
        if (message.next != null && message.hasOwnProperty("next"))
            if (!$util.isString(message.next))
                return "next: string expected";
        return null;
    };

    /**
     * Creates a FeedItemList message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof FeedItemList
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {FeedItemList} FeedItemList
     */
    FeedItemList.fromObject = function fromObject(object) {
        if (object instanceof $root.FeedItemList)
            return object;
        let message = new $root.FeedItemList();
        if (object.items) {
            if (!Array.isArray(object.items))
                throw TypeError(".FeedItemList.items: array expected");
            message.items = [];
            for (let i = 0; i < object.items.length; ++i) {
                if (typeof object.items[i] !== "object")
                    throw TypeError(".FeedItemList.items: object expected");
                message.items[i] = $root.FeedItem.fromObject(object.items[i]);
            }
        }
        if (object.count != null)
            message.count = object.count | 0;
        if (object.next != null)
            message.next = String(object.next);
        return message;
    };

    /**
     * Creates a plain object from a FeedItemList message. Also converts values to other types if specified.
     * @function toObject
     * @memberof FeedItemList
     * @static
     * @param {FeedItemList} message FeedItemList
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    FeedItemList.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.items = [];
        if (options.defaults) {
            object.count = 0;
            object.next = "";
        }
        if (message.items && message.items.length) {
            object.items = [];
            for (let j = 0; j < message.items.length; ++j)
                object.items[j] = $root.FeedItem.toObject(message.items[j], options);
        }
        if (message.count != null && message.hasOwnProperty("count"))
            object.count = message.count;
        if (message.next != null && message.hasOwnProperty("next"))
            object.next = message.next;
        return object;
    };

    /**
     * Converts this FeedItemList to JSON.
     * @function toJSON
     * @memberof FeedItemList
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    FeedItemList.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return FeedItemList;
})();

export const Merge = $root.Merge = (() => {

    /**
     * Properties of a Merge.
     * @exports IMerge
     * @interface IMerge
     * @property {string} block Merge block
     * @property {google.protobuf.ITimestamp} date Merge date
     * @property {Array.<IFeedItem>} targets Merge targets
     */

    /**
     * Constructs a new Merge.
     * @exports Merge
     * @classdesc Represents a Merge.
     * @implements IMerge
     * @constructor
     * @param {IMerge=} [properties] Properties to set
     */
    function Merge(properties) {
        this.targets = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Merge block.
     * @member {string} block
     * @memberof Merge
     * @instance
     */
    Merge.prototype.block = "";

    /**
     * Merge date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof Merge
     * @instance
     */
    Merge.prototype.date = null;

    /**
     * Merge targets.
     * @member {Array.<IFeedItem>} targets
     * @memberof Merge
     * @instance
     */
    Merge.prototype.targets = $util.emptyArray;

    /**
     * Creates a new Merge instance using the specified properties.
     * @function create
     * @memberof Merge
     * @static
     * @param {IMerge=} [properties] Properties to set
     * @returns {Merge} Merge instance
     */
    Merge.create = function create(properties) {
        return new Merge(properties);
    };

    /**
     * Encodes the specified Merge message. Does not implicitly {@link Merge.verify|verify} messages.
     * @function encode
     * @memberof Merge
     * @static
     * @param {IMerge} message Merge message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Merge.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.block != null && message.hasOwnProperty("block"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.block);
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        if (message.targets != null && message.targets.length)
            for (let i = 0; i < message.targets.length; ++i)
                $root.FeedItem.encode(message.targets[i], writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified Merge message, length delimited. Does not implicitly {@link Merge.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Merge
     * @static
     * @param {IMerge} message Merge message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Merge.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Merge message from the specified reader or buffer.
     * @function decode
     * @memberof Merge
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Merge} Merge
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Merge.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Merge();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.block = reader.string();
                break;
            case 2:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 3:
                if (!(message.targets && message.targets.length))
                    message.targets = [];
                message.targets.push($root.FeedItem.decode(reader, reader.uint32()));
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Merge message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Merge
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Merge} Merge
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Merge.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Merge message.
     * @function verify
     * @memberof Merge
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Merge.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.block != null && message.hasOwnProperty("block"))
            if (!$util.isString(message.block))
                return "block: string expected";
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        if (message.targets != null && message.hasOwnProperty("targets")) {
            if (!Array.isArray(message.targets))
                return "targets: array expected";
            for (let i = 0; i < message.targets.length; ++i) {
                let error = $root.FeedItem.verify(message.targets[i]);
                if (error)
                    return "targets." + error;
            }
        }
        return null;
    };

    /**
     * Creates a Merge message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Merge
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Merge} Merge
     */
    Merge.fromObject = function fromObject(object) {
        if (object instanceof $root.Merge)
            return object;
        let message = new $root.Merge();
        if (object.block != null)
            message.block = String(object.block);
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".Merge.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        if (object.targets) {
            if (!Array.isArray(object.targets))
                throw TypeError(".Merge.targets: array expected");
            message.targets = [];
            for (let i = 0; i < object.targets.length; ++i) {
                if (typeof object.targets[i] !== "object")
                    throw TypeError(".Merge.targets: object expected");
                message.targets[i] = $root.FeedItem.fromObject(object.targets[i]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from a Merge message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Merge
     * @static
     * @param {Merge} message Merge
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Merge.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.targets = [];
        if (options.defaults) {
            object.block = "";
            object.date = null;
        }
        if (message.block != null && message.hasOwnProperty("block"))
            object.block = message.block;
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        if (message.targets && message.targets.length) {
            object.targets = [];
            for (let j = 0; j < message.targets.length; ++j)
                object.targets[j] = $root.FeedItem.toObject(message.targets[j], options);
        }
        return object;
    };

    /**
     * Converts this Merge to JSON.
     * @function toJSON
     * @memberof Merge
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Merge.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Merge;
})();

export const Ignore = $root.Ignore = (() => {

    /**
     * Properties of an Ignore.
     * @exports IIgnore
     * @interface IIgnore
     * @property {string} block Ignore block
     * @property {google.protobuf.ITimestamp} date Ignore date
     * @property {IUser} user Ignore user
     * @property {IFeedItem} target Ignore target
     */

    /**
     * Constructs a new Ignore.
     * @exports Ignore
     * @classdesc Represents an Ignore.
     * @implements IIgnore
     * @constructor
     * @param {IIgnore=} [properties] Properties to set
     */
    function Ignore(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Ignore block.
     * @member {string} block
     * @memberof Ignore
     * @instance
     */
    Ignore.prototype.block = "";

    /**
     * Ignore date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof Ignore
     * @instance
     */
    Ignore.prototype.date = null;

    /**
     * Ignore user.
     * @member {IUser} user
     * @memberof Ignore
     * @instance
     */
    Ignore.prototype.user = null;

    /**
     * Ignore target.
     * @member {IFeedItem} target
     * @memberof Ignore
     * @instance
     */
    Ignore.prototype.target = null;

    /**
     * Creates a new Ignore instance using the specified properties.
     * @function create
     * @memberof Ignore
     * @static
     * @param {IIgnore=} [properties] Properties to set
     * @returns {Ignore} Ignore instance
     */
    Ignore.create = function create(properties) {
        return new Ignore(properties);
    };

    /**
     * Encodes the specified Ignore message. Does not implicitly {@link Ignore.verify|verify} messages.
     * @function encode
     * @memberof Ignore
     * @static
     * @param {IIgnore} message Ignore message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Ignore.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.block != null && message.hasOwnProperty("block"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.block);
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        if (message.user != null && message.hasOwnProperty("user"))
            $root.User.encode(message.user, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        if (message.target != null && message.hasOwnProperty("target"))
            $root.FeedItem.encode(message.target, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified Ignore message, length delimited. Does not implicitly {@link Ignore.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Ignore
     * @static
     * @param {IIgnore} message Ignore message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Ignore.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes an Ignore message from the specified reader or buffer.
     * @function decode
     * @memberof Ignore
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Ignore} Ignore
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Ignore.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Ignore();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.block = reader.string();
                break;
            case 2:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 3:
                message.user = $root.User.decode(reader, reader.uint32());
                break;
            case 4:
                message.target = $root.FeedItem.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes an Ignore message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Ignore
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Ignore} Ignore
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Ignore.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies an Ignore message.
     * @function verify
     * @memberof Ignore
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Ignore.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.block != null && message.hasOwnProperty("block"))
            if (!$util.isString(message.block))
                return "block: string expected";
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        if (message.user != null && message.hasOwnProperty("user")) {
            let error = $root.User.verify(message.user);
            if (error)
                return "user." + error;
        }
        if (message.target != null && message.hasOwnProperty("target")) {
            let error = $root.FeedItem.verify(message.target);
            if (error)
                return "target." + error;
        }
        return null;
    };

    /**
     * Creates an Ignore message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Ignore
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Ignore} Ignore
     */
    Ignore.fromObject = function fromObject(object) {
        if (object instanceof $root.Ignore)
            return object;
        let message = new $root.Ignore();
        if (object.block != null)
            message.block = String(object.block);
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".Ignore.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        if (object.user != null) {
            if (typeof object.user !== "object")
                throw TypeError(".Ignore.user: object expected");
            message.user = $root.User.fromObject(object.user);
        }
        if (object.target != null) {
            if (typeof object.target !== "object")
                throw TypeError(".Ignore.target: object expected");
            message.target = $root.FeedItem.fromObject(object.target);
        }
        return message;
    };

    /**
     * Creates a plain object from an Ignore message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Ignore
     * @static
     * @param {Ignore} message Ignore
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Ignore.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.block = "";
            object.date = null;
            object.user = null;
            object.target = null;
        }
        if (message.block != null && message.hasOwnProperty("block"))
            object.block = message.block;
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        if (message.user != null && message.hasOwnProperty("user"))
            object.user = $root.User.toObject(message.user, options);
        if (message.target != null && message.hasOwnProperty("target"))
            object.target = $root.FeedItem.toObject(message.target, options);
        return object;
    };

    /**
     * Converts this Ignore to JSON.
     * @function toJSON
     * @memberof Ignore
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Ignore.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Ignore;
})();

export const Flag = $root.Flag = (() => {

    /**
     * Properties of a Flag.
     * @exports IFlag
     * @interface IFlag
     * @property {string} block Flag block
     * @property {google.protobuf.ITimestamp} date Flag date
     * @property {IUser} user Flag user
     * @property {IFeedItem} target Flag target
     */

    /**
     * Constructs a new Flag.
     * @exports Flag
     * @classdesc Represents a Flag.
     * @implements IFlag
     * @constructor
     * @param {IFlag=} [properties] Properties to set
     */
    function Flag(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Flag block.
     * @member {string} block
     * @memberof Flag
     * @instance
     */
    Flag.prototype.block = "";

    /**
     * Flag date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof Flag
     * @instance
     */
    Flag.prototype.date = null;

    /**
     * Flag user.
     * @member {IUser} user
     * @memberof Flag
     * @instance
     */
    Flag.prototype.user = null;

    /**
     * Flag target.
     * @member {IFeedItem} target
     * @memberof Flag
     * @instance
     */
    Flag.prototype.target = null;

    /**
     * Creates a new Flag instance using the specified properties.
     * @function create
     * @memberof Flag
     * @static
     * @param {IFlag=} [properties] Properties to set
     * @returns {Flag} Flag instance
     */
    Flag.create = function create(properties) {
        return new Flag(properties);
    };

    /**
     * Encodes the specified Flag message. Does not implicitly {@link Flag.verify|verify} messages.
     * @function encode
     * @memberof Flag
     * @static
     * @param {IFlag} message Flag message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Flag.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.block != null && message.hasOwnProperty("block"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.block);
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        if (message.user != null && message.hasOwnProperty("user"))
            $root.User.encode(message.user, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        if (message.target != null && message.hasOwnProperty("target"))
            $root.FeedItem.encode(message.target, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified Flag message, length delimited. Does not implicitly {@link Flag.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Flag
     * @static
     * @param {IFlag} message Flag message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Flag.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Flag message from the specified reader or buffer.
     * @function decode
     * @memberof Flag
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Flag} Flag
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Flag.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Flag();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.block = reader.string();
                break;
            case 2:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 3:
                message.user = $root.User.decode(reader, reader.uint32());
                break;
            case 4:
                message.target = $root.FeedItem.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Flag message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Flag
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Flag} Flag
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Flag.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Flag message.
     * @function verify
     * @memberof Flag
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Flag.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.block != null && message.hasOwnProperty("block"))
            if (!$util.isString(message.block))
                return "block: string expected";
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        if (message.user != null && message.hasOwnProperty("user")) {
            let error = $root.User.verify(message.user);
            if (error)
                return "user." + error;
        }
        if (message.target != null && message.hasOwnProperty("target")) {
            let error = $root.FeedItem.verify(message.target);
            if (error)
                return "target." + error;
        }
        return null;
    };

    /**
     * Creates a Flag message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Flag
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Flag} Flag
     */
    Flag.fromObject = function fromObject(object) {
        if (object instanceof $root.Flag)
            return object;
        let message = new $root.Flag();
        if (object.block != null)
            message.block = String(object.block);
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".Flag.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        if (object.user != null) {
            if (typeof object.user !== "object")
                throw TypeError(".Flag.user: object expected");
            message.user = $root.User.fromObject(object.user);
        }
        if (object.target != null) {
            if (typeof object.target !== "object")
                throw TypeError(".Flag.target: object expected");
            message.target = $root.FeedItem.fromObject(object.target);
        }
        return message;
    };

    /**
     * Creates a plain object from a Flag message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Flag
     * @static
     * @param {Flag} message Flag
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Flag.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.block = "";
            object.date = null;
            object.user = null;
            object.target = null;
        }
        if (message.block != null && message.hasOwnProperty("block"))
            object.block = message.block;
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        if (message.user != null && message.hasOwnProperty("user"))
            object.user = $root.User.toObject(message.user, options);
        if (message.target != null && message.hasOwnProperty("target"))
            object.target = $root.FeedItem.toObject(message.target, options);
        return object;
    };

    /**
     * Converts this Flag to JSON.
     * @function toJSON
     * @memberof Flag
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Flag.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Flag;
})();

export const Announce = $root.Announce = (() => {

    /**
     * Properties of an Announce.
     * @exports IAnnounce
     * @interface IAnnounce
     * @property {string} block Announce block
     * @property {google.protobuf.ITimestamp} date Announce date
     * @property {IUser} user Announce user
     * @property {IFeedItem} target Announce target
     */

    /**
     * Constructs a new Announce.
     * @exports Announce
     * @classdesc Represents an Announce.
     * @implements IAnnounce
     * @constructor
     * @param {IAnnounce=} [properties] Properties to set
     */
    function Announce(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Announce block.
     * @member {string} block
     * @memberof Announce
     * @instance
     */
    Announce.prototype.block = "";

    /**
     * Announce date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof Announce
     * @instance
     */
    Announce.prototype.date = null;

    /**
     * Announce user.
     * @member {IUser} user
     * @memberof Announce
     * @instance
     */
    Announce.prototype.user = null;

    /**
     * Announce target.
     * @member {IFeedItem} target
     * @memberof Announce
     * @instance
     */
    Announce.prototype.target = null;

    /**
     * Creates a new Announce instance using the specified properties.
     * @function create
     * @memberof Announce
     * @static
     * @param {IAnnounce=} [properties] Properties to set
     * @returns {Announce} Announce instance
     */
    Announce.create = function create(properties) {
        return new Announce(properties);
    };

    /**
     * Encodes the specified Announce message. Does not implicitly {@link Announce.verify|verify} messages.
     * @function encode
     * @memberof Announce
     * @static
     * @param {IAnnounce} message Announce message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Announce.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.block != null && message.hasOwnProperty("block"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.block);
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        if (message.user != null && message.hasOwnProperty("user"))
            $root.User.encode(message.user, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        if (message.target != null && message.hasOwnProperty("target"))
            $root.FeedItem.encode(message.target, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified Announce message, length delimited. Does not implicitly {@link Announce.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Announce
     * @static
     * @param {IAnnounce} message Announce message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Announce.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes an Announce message from the specified reader or buffer.
     * @function decode
     * @memberof Announce
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Announce} Announce
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Announce.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Announce();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.block = reader.string();
                break;
            case 2:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 3:
                message.user = $root.User.decode(reader, reader.uint32());
                break;
            case 4:
                message.target = $root.FeedItem.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes an Announce message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Announce
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Announce} Announce
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Announce.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies an Announce message.
     * @function verify
     * @memberof Announce
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Announce.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.block != null && message.hasOwnProperty("block"))
            if (!$util.isString(message.block))
                return "block: string expected";
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        if (message.user != null && message.hasOwnProperty("user")) {
            let error = $root.User.verify(message.user);
            if (error)
                return "user." + error;
        }
        if (message.target != null && message.hasOwnProperty("target")) {
            let error = $root.FeedItem.verify(message.target);
            if (error)
                return "target." + error;
        }
        return null;
    };

    /**
     * Creates an Announce message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Announce
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Announce} Announce
     */
    Announce.fromObject = function fromObject(object) {
        if (object instanceof $root.Announce)
            return object;
        let message = new $root.Announce();
        if (object.block != null)
            message.block = String(object.block);
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".Announce.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        if (object.user != null) {
            if (typeof object.user !== "object")
                throw TypeError(".Announce.user: object expected");
            message.user = $root.User.fromObject(object.user);
        }
        if (object.target != null) {
            if (typeof object.target !== "object")
                throw TypeError(".Announce.target: object expected");
            message.target = $root.FeedItem.fromObject(object.target);
        }
        return message;
    };

    /**
     * Creates a plain object from an Announce message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Announce
     * @static
     * @param {Announce} message Announce
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Announce.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.block = "";
            object.date = null;
            object.user = null;
            object.target = null;
        }
        if (message.block != null && message.hasOwnProperty("block"))
            object.block = message.block;
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        if (message.user != null && message.hasOwnProperty("user"))
            object.user = $root.User.toObject(message.user, options);
        if (message.target != null && message.hasOwnProperty("target"))
            object.target = $root.FeedItem.toObject(message.target, options);
        return object;
    };

    /**
     * Converts this Announce to JSON.
     * @function toJSON
     * @memberof Announce
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Announce.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Announce;
})();

export const Join = $root.Join = (() => {

    /**
     * Properties of a Join.
     * @exports IJoin
     * @interface IJoin
     * @property {string} block Join block
     * @property {google.protobuf.ITimestamp} date Join date
     * @property {IUser} user Join user
     * @property {Array.<ILike>} likes Join likes
     */

    /**
     * Constructs a new Join.
     * @exports Join
     * @classdesc Represents a Join.
     * @implements IJoin
     * @constructor
     * @param {IJoin=} [properties] Properties to set
     */
    function Join(properties) {
        this.likes = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Join block.
     * @member {string} block
     * @memberof Join
     * @instance
     */
    Join.prototype.block = "";

    /**
     * Join date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof Join
     * @instance
     */
    Join.prototype.date = null;

    /**
     * Join user.
     * @member {IUser} user
     * @memberof Join
     * @instance
     */
    Join.prototype.user = null;

    /**
     * Join likes.
     * @member {Array.<ILike>} likes
     * @memberof Join
     * @instance
     */
    Join.prototype.likes = $util.emptyArray;

    /**
     * Creates a new Join instance using the specified properties.
     * @function create
     * @memberof Join
     * @static
     * @param {IJoin=} [properties] Properties to set
     * @returns {Join} Join instance
     */
    Join.create = function create(properties) {
        return new Join(properties);
    };

    /**
     * Encodes the specified Join message. Does not implicitly {@link Join.verify|verify} messages.
     * @function encode
     * @memberof Join
     * @static
     * @param {IJoin} message Join message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Join.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.block != null && message.hasOwnProperty("block"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.block);
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        if (message.user != null && message.hasOwnProperty("user"))
            $root.User.encode(message.user, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        if (message.likes != null && message.likes.length)
            for (let i = 0; i < message.likes.length; ++i)
                $root.Like.encode(message.likes[i], writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified Join message, length delimited. Does not implicitly {@link Join.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Join
     * @static
     * @param {IJoin} message Join message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Join.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Join message from the specified reader or buffer.
     * @function decode
     * @memberof Join
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Join} Join
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Join.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Join();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.block = reader.string();
                break;
            case 2:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 3:
                message.user = $root.User.decode(reader, reader.uint32());
                break;
            case 4:
                if (!(message.likes && message.likes.length))
                    message.likes = [];
                message.likes.push($root.Like.decode(reader, reader.uint32()));
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Join message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Join
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Join} Join
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Join.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Join message.
     * @function verify
     * @memberof Join
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Join.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.block != null && message.hasOwnProperty("block"))
            if (!$util.isString(message.block))
                return "block: string expected";
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        if (message.user != null && message.hasOwnProperty("user")) {
            let error = $root.User.verify(message.user);
            if (error)
                return "user." + error;
        }
        if (message.likes != null && message.hasOwnProperty("likes")) {
            if (!Array.isArray(message.likes))
                return "likes: array expected";
            for (let i = 0; i < message.likes.length; ++i) {
                let error = $root.Like.verify(message.likes[i]);
                if (error)
                    return "likes." + error;
            }
        }
        return null;
    };

    /**
     * Creates a Join message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Join
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Join} Join
     */
    Join.fromObject = function fromObject(object) {
        if (object instanceof $root.Join)
            return object;
        let message = new $root.Join();
        if (object.block != null)
            message.block = String(object.block);
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".Join.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        if (object.user != null) {
            if (typeof object.user !== "object")
                throw TypeError(".Join.user: object expected");
            message.user = $root.User.fromObject(object.user);
        }
        if (object.likes) {
            if (!Array.isArray(object.likes))
                throw TypeError(".Join.likes: array expected");
            message.likes = [];
            for (let i = 0; i < object.likes.length; ++i) {
                if (typeof object.likes[i] !== "object")
                    throw TypeError(".Join.likes: object expected");
                message.likes[i] = $root.Like.fromObject(object.likes[i]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from a Join message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Join
     * @static
     * @param {Join} message Join
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Join.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.likes = [];
        if (options.defaults) {
            object.block = "";
            object.date = null;
            object.user = null;
        }
        if (message.block != null && message.hasOwnProperty("block"))
            object.block = message.block;
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        if (message.user != null && message.hasOwnProperty("user"))
            object.user = $root.User.toObject(message.user, options);
        if (message.likes && message.likes.length) {
            object.likes = [];
            for (let j = 0; j < message.likes.length; ++j)
                object.likes[j] = $root.Like.toObject(message.likes[j], options);
        }
        return object;
    };

    /**
     * Converts this Join to JSON.
     * @function toJSON
     * @memberof Join
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Join.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Join;
})();

export const Leave = $root.Leave = (() => {

    /**
     * Properties of a Leave.
     * @exports ILeave
     * @interface ILeave
     * @property {string} block Leave block
     * @property {google.protobuf.ITimestamp} date Leave date
     * @property {IUser} user Leave user
     * @property {Array.<ILike>} likes Leave likes
     */

    /**
     * Constructs a new Leave.
     * @exports Leave
     * @classdesc Represents a Leave.
     * @implements ILeave
     * @constructor
     * @param {ILeave=} [properties] Properties to set
     */
    function Leave(properties) {
        this.likes = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Leave block.
     * @member {string} block
     * @memberof Leave
     * @instance
     */
    Leave.prototype.block = "";

    /**
     * Leave date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof Leave
     * @instance
     */
    Leave.prototype.date = null;

    /**
     * Leave user.
     * @member {IUser} user
     * @memberof Leave
     * @instance
     */
    Leave.prototype.user = null;

    /**
     * Leave likes.
     * @member {Array.<ILike>} likes
     * @memberof Leave
     * @instance
     */
    Leave.prototype.likes = $util.emptyArray;

    /**
     * Creates a new Leave instance using the specified properties.
     * @function create
     * @memberof Leave
     * @static
     * @param {ILeave=} [properties] Properties to set
     * @returns {Leave} Leave instance
     */
    Leave.create = function create(properties) {
        return new Leave(properties);
    };

    /**
     * Encodes the specified Leave message. Does not implicitly {@link Leave.verify|verify} messages.
     * @function encode
     * @memberof Leave
     * @static
     * @param {ILeave} message Leave message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Leave.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.block != null && message.hasOwnProperty("block"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.block);
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        if (message.user != null && message.hasOwnProperty("user"))
            $root.User.encode(message.user, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        if (message.likes != null && message.likes.length)
            for (let i = 0; i < message.likes.length; ++i)
                $root.Like.encode(message.likes[i], writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified Leave message, length delimited. Does not implicitly {@link Leave.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Leave
     * @static
     * @param {ILeave} message Leave message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Leave.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Leave message from the specified reader or buffer.
     * @function decode
     * @memberof Leave
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Leave} Leave
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Leave.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Leave();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.block = reader.string();
                break;
            case 2:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 3:
                message.user = $root.User.decode(reader, reader.uint32());
                break;
            case 4:
                if (!(message.likes && message.likes.length))
                    message.likes = [];
                message.likes.push($root.Like.decode(reader, reader.uint32()));
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Leave message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Leave
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Leave} Leave
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Leave.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Leave message.
     * @function verify
     * @memberof Leave
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Leave.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.block != null && message.hasOwnProperty("block"))
            if (!$util.isString(message.block))
                return "block: string expected";
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        if (message.user != null && message.hasOwnProperty("user")) {
            let error = $root.User.verify(message.user);
            if (error)
                return "user." + error;
        }
        if (message.likes != null && message.hasOwnProperty("likes")) {
            if (!Array.isArray(message.likes))
                return "likes: array expected";
            for (let i = 0; i < message.likes.length; ++i) {
                let error = $root.Like.verify(message.likes[i]);
                if (error)
                    return "likes." + error;
            }
        }
        return null;
    };

    /**
     * Creates a Leave message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Leave
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Leave} Leave
     */
    Leave.fromObject = function fromObject(object) {
        if (object instanceof $root.Leave)
            return object;
        let message = new $root.Leave();
        if (object.block != null)
            message.block = String(object.block);
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".Leave.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        if (object.user != null) {
            if (typeof object.user !== "object")
                throw TypeError(".Leave.user: object expected");
            message.user = $root.User.fromObject(object.user);
        }
        if (object.likes) {
            if (!Array.isArray(object.likes))
                throw TypeError(".Leave.likes: array expected");
            message.likes = [];
            for (let i = 0; i < object.likes.length; ++i) {
                if (typeof object.likes[i] !== "object")
                    throw TypeError(".Leave.likes: object expected");
                message.likes[i] = $root.Like.fromObject(object.likes[i]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from a Leave message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Leave
     * @static
     * @param {Leave} message Leave
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Leave.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.likes = [];
        if (options.defaults) {
            object.block = "";
            object.date = null;
            object.user = null;
        }
        if (message.block != null && message.hasOwnProperty("block"))
            object.block = message.block;
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        if (message.user != null && message.hasOwnProperty("user"))
            object.user = $root.User.toObject(message.user, options);
        if (message.likes && message.likes.length) {
            object.likes = [];
            for (let j = 0; j < message.likes.length; ++j)
                object.likes[j] = $root.Like.toObject(message.likes[j], options);
        }
        return object;
    };

    /**
     * Converts this Leave to JSON.
     * @function toJSON
     * @memberof Leave
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Leave.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Leave;
})();

export const Text = $root.Text = (() => {

    /**
     * Properties of a Text.
     * @exports IText
     * @interface IText
     * @property {string} block Text block
     * @property {google.protobuf.ITimestamp} date Text date
     * @property {IUser} user Text user
     * @property {string} body Text body
     * @property {Array.<IComment>} comments Text comments
     * @property {Array.<ILike>} likes Text likes
     */

    /**
     * Constructs a new Text.
     * @exports Text
     * @classdesc Represents a Text.
     * @implements IText
     * @constructor
     * @param {IText=} [properties] Properties to set
     */
    function Text(properties) {
        this.comments = [];
        this.likes = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Text block.
     * @member {string} block
     * @memberof Text
     * @instance
     */
    Text.prototype.block = "";

    /**
     * Text date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof Text
     * @instance
     */
    Text.prototype.date = null;

    /**
     * Text user.
     * @member {IUser} user
     * @memberof Text
     * @instance
     */
    Text.prototype.user = null;

    /**
     * Text body.
     * @member {string} body
     * @memberof Text
     * @instance
     */
    Text.prototype.body = "";

    /**
     * Text comments.
     * @member {Array.<IComment>} comments
     * @memberof Text
     * @instance
     */
    Text.prototype.comments = $util.emptyArray;

    /**
     * Text likes.
     * @member {Array.<ILike>} likes
     * @memberof Text
     * @instance
     */
    Text.prototype.likes = $util.emptyArray;

    /**
     * Creates a new Text instance using the specified properties.
     * @function create
     * @memberof Text
     * @static
     * @param {IText=} [properties] Properties to set
     * @returns {Text} Text instance
     */
    Text.create = function create(properties) {
        return new Text(properties);
    };

    /**
     * Encodes the specified Text message. Does not implicitly {@link Text.verify|verify} messages.
     * @function encode
     * @memberof Text
     * @static
     * @param {IText} message Text message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Text.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.block != null && message.hasOwnProperty("block"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.block);
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        if (message.user != null && message.hasOwnProperty("user"))
            $root.User.encode(message.user, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        if (message.body != null && message.hasOwnProperty("body"))
            writer.uint32(/* id 4, wireType 2 =*/34).string(message.body);
        if (message.comments != null && message.comments.length)
            for (let i = 0; i < message.comments.length; ++i)
                $root.Comment.encode(message.comments[i], writer.uint32(/* id 5, wireType 2 =*/42).fork()).ldelim();
        if (message.likes != null && message.likes.length)
            for (let i = 0; i < message.likes.length; ++i)
                $root.Like.encode(message.likes[i], writer.uint32(/* id 6, wireType 2 =*/50).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified Text message, length delimited. Does not implicitly {@link Text.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Text
     * @static
     * @param {IText} message Text message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Text.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Text message from the specified reader or buffer.
     * @function decode
     * @memberof Text
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Text} Text
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Text.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Text();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.block = reader.string();
                break;
            case 2:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 3:
                message.user = $root.User.decode(reader, reader.uint32());
                break;
            case 4:
                message.body = reader.string();
                break;
            case 5:
                if (!(message.comments && message.comments.length))
                    message.comments = [];
                message.comments.push($root.Comment.decode(reader, reader.uint32()));
                break;
            case 6:
                if (!(message.likes && message.likes.length))
                    message.likes = [];
                message.likes.push($root.Like.decode(reader, reader.uint32()));
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Text message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Text
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Text} Text
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Text.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Text message.
     * @function verify
     * @memberof Text
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Text.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.block != null && message.hasOwnProperty("block"))
            if (!$util.isString(message.block))
                return "block: string expected";
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        if (message.user != null && message.hasOwnProperty("user")) {
            let error = $root.User.verify(message.user);
            if (error)
                return "user." + error;
        }
        if (message.body != null && message.hasOwnProperty("body"))
            if (!$util.isString(message.body))
                return "body: string expected";
        if (message.comments != null && message.hasOwnProperty("comments")) {
            if (!Array.isArray(message.comments))
                return "comments: array expected";
            for (let i = 0; i < message.comments.length; ++i) {
                let error = $root.Comment.verify(message.comments[i]);
                if (error)
                    return "comments." + error;
            }
        }
        if (message.likes != null && message.hasOwnProperty("likes")) {
            if (!Array.isArray(message.likes))
                return "likes: array expected";
            for (let i = 0; i < message.likes.length; ++i) {
                let error = $root.Like.verify(message.likes[i]);
                if (error)
                    return "likes." + error;
            }
        }
        return null;
    };

    /**
     * Creates a Text message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Text
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Text} Text
     */
    Text.fromObject = function fromObject(object) {
        if (object instanceof $root.Text)
            return object;
        let message = new $root.Text();
        if (object.block != null)
            message.block = String(object.block);
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".Text.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        if (object.user != null) {
            if (typeof object.user !== "object")
                throw TypeError(".Text.user: object expected");
            message.user = $root.User.fromObject(object.user);
        }
        if (object.body != null)
            message.body = String(object.body);
        if (object.comments) {
            if (!Array.isArray(object.comments))
                throw TypeError(".Text.comments: array expected");
            message.comments = [];
            for (let i = 0; i < object.comments.length; ++i) {
                if (typeof object.comments[i] !== "object")
                    throw TypeError(".Text.comments: object expected");
                message.comments[i] = $root.Comment.fromObject(object.comments[i]);
            }
        }
        if (object.likes) {
            if (!Array.isArray(object.likes))
                throw TypeError(".Text.likes: array expected");
            message.likes = [];
            for (let i = 0; i < object.likes.length; ++i) {
                if (typeof object.likes[i] !== "object")
                    throw TypeError(".Text.likes: object expected");
                message.likes[i] = $root.Like.fromObject(object.likes[i]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from a Text message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Text
     * @static
     * @param {Text} message Text
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Text.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults) {
            object.comments = [];
            object.likes = [];
        }
        if (options.defaults) {
            object.block = "";
            object.date = null;
            object.user = null;
            object.body = "";
        }
        if (message.block != null && message.hasOwnProperty("block"))
            object.block = message.block;
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        if (message.user != null && message.hasOwnProperty("user"))
            object.user = $root.User.toObject(message.user, options);
        if (message.body != null && message.hasOwnProperty("body"))
            object.body = message.body;
        if (message.comments && message.comments.length) {
            object.comments = [];
            for (let j = 0; j < message.comments.length; ++j)
                object.comments[j] = $root.Comment.toObject(message.comments[j], options);
        }
        if (message.likes && message.likes.length) {
            object.likes = [];
            for (let j = 0; j < message.likes.length; ++j)
                object.likes[j] = $root.Like.toObject(message.likes[j], options);
        }
        return object;
    };

    /**
     * Converts this Text to JSON.
     * @function toJSON
     * @memberof Text
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Text.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Text;
})();

export const TextList = $root.TextList = (() => {

    /**
     * Properties of a TextList.
     * @exports ITextList
     * @interface ITextList
     * @property {Array.<IText>} items TextList items
     */

    /**
     * Constructs a new TextList.
     * @exports TextList
     * @classdesc Represents a TextList.
     * @implements ITextList
     * @constructor
     * @param {ITextList=} [properties] Properties to set
     */
    function TextList(properties) {
        this.items = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * TextList items.
     * @member {Array.<IText>} items
     * @memberof TextList
     * @instance
     */
    TextList.prototype.items = $util.emptyArray;

    /**
     * Creates a new TextList instance using the specified properties.
     * @function create
     * @memberof TextList
     * @static
     * @param {ITextList=} [properties] Properties to set
     * @returns {TextList} TextList instance
     */
    TextList.create = function create(properties) {
        return new TextList(properties);
    };

    /**
     * Encodes the specified TextList message. Does not implicitly {@link TextList.verify|verify} messages.
     * @function encode
     * @memberof TextList
     * @static
     * @param {ITextList} message TextList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    TextList.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.items != null && message.items.length)
            for (let i = 0; i < message.items.length; ++i)
                $root.Text.encode(message.items[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified TextList message, length delimited. Does not implicitly {@link TextList.verify|verify} messages.
     * @function encodeDelimited
     * @memberof TextList
     * @static
     * @param {ITextList} message TextList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    TextList.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a TextList message from the specified reader or buffer.
     * @function decode
     * @memberof TextList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {TextList} TextList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    TextList.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.TextList();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                if (!(message.items && message.items.length))
                    message.items = [];
                message.items.push($root.Text.decode(reader, reader.uint32()));
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a TextList message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof TextList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {TextList} TextList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    TextList.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a TextList message.
     * @function verify
     * @memberof TextList
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    TextList.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.items != null && message.hasOwnProperty("items")) {
            if (!Array.isArray(message.items))
                return "items: array expected";
            for (let i = 0; i < message.items.length; ++i) {
                let error = $root.Text.verify(message.items[i]);
                if (error)
                    return "items." + error;
            }
        }
        return null;
    };

    /**
     * Creates a TextList message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof TextList
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {TextList} TextList
     */
    TextList.fromObject = function fromObject(object) {
        if (object instanceof $root.TextList)
            return object;
        let message = new $root.TextList();
        if (object.items) {
            if (!Array.isArray(object.items))
                throw TypeError(".TextList.items: array expected");
            message.items = [];
            for (let i = 0; i < object.items.length; ++i) {
                if (typeof object.items[i] !== "object")
                    throw TypeError(".TextList.items: object expected");
                message.items[i] = $root.Text.fromObject(object.items[i]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from a TextList message. Also converts values to other types if specified.
     * @function toObject
     * @memberof TextList
     * @static
     * @param {TextList} message TextList
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    TextList.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.items = [];
        if (message.items && message.items.length) {
            object.items = [];
            for (let j = 0; j < message.items.length; ++j)
                object.items[j] = $root.Text.toObject(message.items[j], options);
        }
        return object;
    };

    /**
     * Converts this TextList to JSON.
     * @function toJSON
     * @memberof TextList
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    TextList.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return TextList;
})();

export const File = $root.File = (() => {

    /**
     * Properties of a File.
     * @exports IFile
     * @interface IFile
     * @property {number} index File index
     * @property {IFileIndex} file File file
     * @property {Object.<string,IFileIndex>} links File links
     */

    /**
     * Constructs a new File.
     * @exports File
     * @classdesc Represents a File.
     * @implements IFile
     * @constructor
     * @param {IFile=} [properties] Properties to set
     */
    function File(properties) {
        this.links = {};
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * File index.
     * @member {number} index
     * @memberof File
     * @instance
     */
    File.prototype.index = 0;

    /**
     * File file.
     * @member {IFileIndex} file
     * @memberof File
     * @instance
     */
    File.prototype.file = null;

    /**
     * File links.
     * @member {Object.<string,IFileIndex>} links
     * @memberof File
     * @instance
     */
    File.prototype.links = $util.emptyObject;

    /**
     * Creates a new File instance using the specified properties.
     * @function create
     * @memberof File
     * @static
     * @param {IFile=} [properties] Properties to set
     * @returns {File} File instance
     */
    File.create = function create(properties) {
        return new File(properties);
    };

    /**
     * Encodes the specified File message. Does not implicitly {@link File.verify|verify} messages.
     * @function encode
     * @memberof File
     * @static
     * @param {IFile} message File message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    File.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.index != null && message.hasOwnProperty("index"))
            writer.uint32(/* id 1, wireType 0 =*/8).int32(message.index);
        if (message.file != null && message.hasOwnProperty("file"))
            $root.FileIndex.encode(message.file, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        if (message.links != null && message.hasOwnProperty("links"))
            for (let keys = Object.keys(message.links), i = 0; i < keys.length; ++i) {
                writer.uint32(/* id 3, wireType 2 =*/26).fork().uint32(/* id 1, wireType 2 =*/10).string(keys[i]);
                $root.FileIndex.encode(message.links[keys[i]], writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim().ldelim();
            }
        return writer;
    };

    /**
     * Encodes the specified File message, length delimited. Does not implicitly {@link File.verify|verify} messages.
     * @function encodeDelimited
     * @memberof File
     * @static
     * @param {IFile} message File message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    File.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a File message from the specified reader or buffer.
     * @function decode
     * @memberof File
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {File} File
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    File.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.File(), key;
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.index = reader.int32();
                break;
            case 2:
                message.file = $root.FileIndex.decode(reader, reader.uint32());
                break;
            case 3:
                reader.skip().pos++;
                if (message.links === $util.emptyObject)
                    message.links = {};
                key = reader.string();
                reader.pos++;
                message.links[key] = $root.FileIndex.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a File message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof File
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {File} File
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    File.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a File message.
     * @function verify
     * @memberof File
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    File.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.index != null && message.hasOwnProperty("index"))
            if (!$util.isInteger(message.index))
                return "index: integer expected";
        if (message.file != null && message.hasOwnProperty("file")) {
            let error = $root.FileIndex.verify(message.file);
            if (error)
                return "file." + error;
        }
        if (message.links != null && message.hasOwnProperty("links")) {
            if (!$util.isObject(message.links))
                return "links: object expected";
            let key = Object.keys(message.links);
            for (let i = 0; i < key.length; ++i) {
                let error = $root.FileIndex.verify(message.links[key[i]]);
                if (error)
                    return "links." + error;
            }
        }
        return null;
    };

    /**
     * Creates a File message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof File
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {File} File
     */
    File.fromObject = function fromObject(object) {
        if (object instanceof $root.File)
            return object;
        let message = new $root.File();
        if (object.index != null)
            message.index = object.index | 0;
        if (object.file != null) {
            if (typeof object.file !== "object")
                throw TypeError(".File.file: object expected");
            message.file = $root.FileIndex.fromObject(object.file);
        }
        if (object.links) {
            if (typeof object.links !== "object")
                throw TypeError(".File.links: object expected");
            message.links = {};
            for (let keys = Object.keys(object.links), i = 0; i < keys.length; ++i) {
                if (typeof object.links[keys[i]] !== "object")
                    throw TypeError(".File.links: object expected");
                message.links[keys[i]] = $root.FileIndex.fromObject(object.links[keys[i]]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from a File message. Also converts values to other types if specified.
     * @function toObject
     * @memberof File
     * @static
     * @param {File} message File
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    File.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.objects || options.defaults)
            object.links = {};
        if (options.defaults) {
            object.index = 0;
            object.file = null;
        }
        if (message.index != null && message.hasOwnProperty("index"))
            object.index = message.index;
        if (message.file != null && message.hasOwnProperty("file"))
            object.file = $root.FileIndex.toObject(message.file, options);
        let keys2;
        if (message.links && (keys2 = Object.keys(message.links)).length) {
            object.links = {};
            for (let j = 0; j < keys2.length; ++j)
                object.links[keys2[j]] = $root.FileIndex.toObject(message.links[keys2[j]], options);
        }
        return object;
    };

    /**
     * Converts this File to JSON.
     * @function toJSON
     * @memberof File
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    File.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return File;
})();

export const Files = $root.Files = (() => {

    /**
     * Properties of a Files.
     * @exports IFiles
     * @interface IFiles
     * @property {string} block Files block
     * @property {string} target Files target
     * @property {google.protobuf.ITimestamp} date Files date
     * @property {IUser} user Files user
     * @property {string} caption Files caption
     * @property {Array.<IFile>} files Files files
     * @property {Array.<IComment>} comments Files comments
     * @property {Array.<ILike>} likes Files likes
     * @property {Array.<string>} threads Files threads
     */

    /**
     * Constructs a new Files.
     * @exports Files
     * @classdesc Represents a Files.
     * @implements IFiles
     * @constructor
     * @param {IFiles=} [properties] Properties to set
     */
    function Files(properties) {
        this.files = [];
        this.comments = [];
        this.likes = [];
        this.threads = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Files block.
     * @member {string} block
     * @memberof Files
     * @instance
     */
    Files.prototype.block = "";

    /**
     * Files target.
     * @member {string} target
     * @memberof Files
     * @instance
     */
    Files.prototype.target = "";

    /**
     * Files date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof Files
     * @instance
     */
    Files.prototype.date = null;

    /**
     * Files user.
     * @member {IUser} user
     * @memberof Files
     * @instance
     */
    Files.prototype.user = null;

    /**
     * Files caption.
     * @member {string} caption
     * @memberof Files
     * @instance
     */
    Files.prototype.caption = "";

    /**
     * Files files.
     * @member {Array.<IFile>} files
     * @memberof Files
     * @instance
     */
    Files.prototype.files = $util.emptyArray;

    /**
     * Files comments.
     * @member {Array.<IComment>} comments
     * @memberof Files
     * @instance
     */
    Files.prototype.comments = $util.emptyArray;

    /**
     * Files likes.
     * @member {Array.<ILike>} likes
     * @memberof Files
     * @instance
     */
    Files.prototype.likes = $util.emptyArray;

    /**
     * Files threads.
     * @member {Array.<string>} threads
     * @memberof Files
     * @instance
     */
    Files.prototype.threads = $util.emptyArray;

    /**
     * Creates a new Files instance using the specified properties.
     * @function create
     * @memberof Files
     * @static
     * @param {IFiles=} [properties] Properties to set
     * @returns {Files} Files instance
     */
    Files.create = function create(properties) {
        return new Files(properties);
    };

    /**
     * Encodes the specified Files message. Does not implicitly {@link Files.verify|verify} messages.
     * @function encode
     * @memberof Files
     * @static
     * @param {IFiles} message Files message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Files.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.block != null && message.hasOwnProperty("block"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.block);
        if (message.target != null && message.hasOwnProperty("target"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.target);
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        if (message.user != null && message.hasOwnProperty("user"))
            $root.User.encode(message.user, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
        if (message.caption != null && message.hasOwnProperty("caption"))
            writer.uint32(/* id 5, wireType 2 =*/42).string(message.caption);
        if (message.files != null && message.files.length)
            for (let i = 0; i < message.files.length; ++i)
                $root.File.encode(message.files[i], writer.uint32(/* id 6, wireType 2 =*/50).fork()).ldelim();
        if (message.comments != null && message.comments.length)
            for (let i = 0; i < message.comments.length; ++i)
                $root.Comment.encode(message.comments[i], writer.uint32(/* id 7, wireType 2 =*/58).fork()).ldelim();
        if (message.likes != null && message.likes.length)
            for (let i = 0; i < message.likes.length; ++i)
                $root.Like.encode(message.likes[i], writer.uint32(/* id 8, wireType 2 =*/66).fork()).ldelim();
        if (message.threads != null && message.threads.length)
            for (let i = 0; i < message.threads.length; ++i)
                writer.uint32(/* id 9, wireType 2 =*/74).string(message.threads[i]);
        return writer;
    };

    /**
     * Encodes the specified Files message, length delimited. Does not implicitly {@link Files.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Files
     * @static
     * @param {IFiles} message Files message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Files.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Files message from the specified reader or buffer.
     * @function decode
     * @memberof Files
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Files} Files
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Files.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Files();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.block = reader.string();
                break;
            case 2:
                message.target = reader.string();
                break;
            case 3:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 4:
                message.user = $root.User.decode(reader, reader.uint32());
                break;
            case 5:
                message.caption = reader.string();
                break;
            case 6:
                if (!(message.files && message.files.length))
                    message.files = [];
                message.files.push($root.File.decode(reader, reader.uint32()));
                break;
            case 7:
                if (!(message.comments && message.comments.length))
                    message.comments = [];
                message.comments.push($root.Comment.decode(reader, reader.uint32()));
                break;
            case 8:
                if (!(message.likes && message.likes.length))
                    message.likes = [];
                message.likes.push($root.Like.decode(reader, reader.uint32()));
                break;
            case 9:
                if (!(message.threads && message.threads.length))
                    message.threads = [];
                message.threads.push(reader.string());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Files message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Files
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Files} Files
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Files.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Files message.
     * @function verify
     * @memberof Files
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Files.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.block != null && message.hasOwnProperty("block"))
            if (!$util.isString(message.block))
                return "block: string expected";
        if (message.target != null && message.hasOwnProperty("target"))
            if (!$util.isString(message.target))
                return "target: string expected";
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        if (message.user != null && message.hasOwnProperty("user")) {
            let error = $root.User.verify(message.user);
            if (error)
                return "user." + error;
        }
        if (message.caption != null && message.hasOwnProperty("caption"))
            if (!$util.isString(message.caption))
                return "caption: string expected";
        if (message.files != null && message.hasOwnProperty("files")) {
            if (!Array.isArray(message.files))
                return "files: array expected";
            for (let i = 0; i < message.files.length; ++i) {
                let error = $root.File.verify(message.files[i]);
                if (error)
                    return "files." + error;
            }
        }
        if (message.comments != null && message.hasOwnProperty("comments")) {
            if (!Array.isArray(message.comments))
                return "comments: array expected";
            for (let i = 0; i < message.comments.length; ++i) {
                let error = $root.Comment.verify(message.comments[i]);
                if (error)
                    return "comments." + error;
            }
        }
        if (message.likes != null && message.hasOwnProperty("likes")) {
            if (!Array.isArray(message.likes))
                return "likes: array expected";
            for (let i = 0; i < message.likes.length; ++i) {
                let error = $root.Like.verify(message.likes[i]);
                if (error)
                    return "likes." + error;
            }
        }
        if (message.threads != null && message.hasOwnProperty("threads")) {
            if (!Array.isArray(message.threads))
                return "threads: array expected";
            for (let i = 0; i < message.threads.length; ++i)
                if (!$util.isString(message.threads[i]))
                    return "threads: string[] expected";
        }
        return null;
    };

    /**
     * Creates a Files message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Files
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Files} Files
     */
    Files.fromObject = function fromObject(object) {
        if (object instanceof $root.Files)
            return object;
        let message = new $root.Files();
        if (object.block != null)
            message.block = String(object.block);
        if (object.target != null)
            message.target = String(object.target);
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".Files.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        if (object.user != null) {
            if (typeof object.user !== "object")
                throw TypeError(".Files.user: object expected");
            message.user = $root.User.fromObject(object.user);
        }
        if (object.caption != null)
            message.caption = String(object.caption);
        if (object.files) {
            if (!Array.isArray(object.files))
                throw TypeError(".Files.files: array expected");
            message.files = [];
            for (let i = 0; i < object.files.length; ++i) {
                if (typeof object.files[i] !== "object")
                    throw TypeError(".Files.files: object expected");
                message.files[i] = $root.File.fromObject(object.files[i]);
            }
        }
        if (object.comments) {
            if (!Array.isArray(object.comments))
                throw TypeError(".Files.comments: array expected");
            message.comments = [];
            for (let i = 0; i < object.comments.length; ++i) {
                if (typeof object.comments[i] !== "object")
                    throw TypeError(".Files.comments: object expected");
                message.comments[i] = $root.Comment.fromObject(object.comments[i]);
            }
        }
        if (object.likes) {
            if (!Array.isArray(object.likes))
                throw TypeError(".Files.likes: array expected");
            message.likes = [];
            for (let i = 0; i < object.likes.length; ++i) {
                if (typeof object.likes[i] !== "object")
                    throw TypeError(".Files.likes: object expected");
                message.likes[i] = $root.Like.fromObject(object.likes[i]);
            }
        }
        if (object.threads) {
            if (!Array.isArray(object.threads))
                throw TypeError(".Files.threads: array expected");
            message.threads = [];
            for (let i = 0; i < object.threads.length; ++i)
                message.threads[i] = String(object.threads[i]);
        }
        return message;
    };

    /**
     * Creates a plain object from a Files message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Files
     * @static
     * @param {Files} message Files
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Files.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults) {
            object.files = [];
            object.comments = [];
            object.likes = [];
            object.threads = [];
        }
        if (options.defaults) {
            object.block = "";
            object.target = "";
            object.date = null;
            object.user = null;
            object.caption = "";
        }
        if (message.block != null && message.hasOwnProperty("block"))
            object.block = message.block;
        if (message.target != null && message.hasOwnProperty("target"))
            object.target = message.target;
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        if (message.user != null && message.hasOwnProperty("user"))
            object.user = $root.User.toObject(message.user, options);
        if (message.caption != null && message.hasOwnProperty("caption"))
            object.caption = message.caption;
        if (message.files && message.files.length) {
            object.files = [];
            for (let j = 0; j < message.files.length; ++j)
                object.files[j] = $root.File.toObject(message.files[j], options);
        }
        if (message.comments && message.comments.length) {
            object.comments = [];
            for (let j = 0; j < message.comments.length; ++j)
                object.comments[j] = $root.Comment.toObject(message.comments[j], options);
        }
        if (message.likes && message.likes.length) {
            object.likes = [];
            for (let j = 0; j < message.likes.length; ++j)
                object.likes[j] = $root.Like.toObject(message.likes[j], options);
        }
        if (message.threads && message.threads.length) {
            object.threads = [];
            for (let j = 0; j < message.threads.length; ++j)
                object.threads[j] = message.threads[j];
        }
        return object;
    };

    /**
     * Converts this Files to JSON.
     * @function toJSON
     * @memberof Files
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Files.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Files;
})();

export const FilesList = $root.FilesList = (() => {

    /**
     * Properties of a FilesList.
     * @exports IFilesList
     * @interface IFilesList
     * @property {Array.<IFiles>} items FilesList items
     */

    /**
     * Constructs a new FilesList.
     * @exports FilesList
     * @classdesc Represents a FilesList.
     * @implements IFilesList
     * @constructor
     * @param {IFilesList=} [properties] Properties to set
     */
    function FilesList(properties) {
        this.items = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * FilesList items.
     * @member {Array.<IFiles>} items
     * @memberof FilesList
     * @instance
     */
    FilesList.prototype.items = $util.emptyArray;

    /**
     * Creates a new FilesList instance using the specified properties.
     * @function create
     * @memberof FilesList
     * @static
     * @param {IFilesList=} [properties] Properties to set
     * @returns {FilesList} FilesList instance
     */
    FilesList.create = function create(properties) {
        return new FilesList(properties);
    };

    /**
     * Encodes the specified FilesList message. Does not implicitly {@link FilesList.verify|verify} messages.
     * @function encode
     * @memberof FilesList
     * @static
     * @param {IFilesList} message FilesList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    FilesList.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.items != null && message.items.length)
            for (let i = 0; i < message.items.length; ++i)
                $root.Files.encode(message.items[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified FilesList message, length delimited. Does not implicitly {@link FilesList.verify|verify} messages.
     * @function encodeDelimited
     * @memberof FilesList
     * @static
     * @param {IFilesList} message FilesList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    FilesList.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a FilesList message from the specified reader or buffer.
     * @function decode
     * @memberof FilesList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {FilesList} FilesList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    FilesList.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.FilesList();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                if (!(message.items && message.items.length))
                    message.items = [];
                message.items.push($root.Files.decode(reader, reader.uint32()));
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a FilesList message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof FilesList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {FilesList} FilesList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    FilesList.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a FilesList message.
     * @function verify
     * @memberof FilesList
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    FilesList.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.items != null && message.hasOwnProperty("items")) {
            if (!Array.isArray(message.items))
                return "items: array expected";
            for (let i = 0; i < message.items.length; ++i) {
                let error = $root.Files.verify(message.items[i]);
                if (error)
                    return "items." + error;
            }
        }
        return null;
    };

    /**
     * Creates a FilesList message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof FilesList
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {FilesList} FilesList
     */
    FilesList.fromObject = function fromObject(object) {
        if (object instanceof $root.FilesList)
            return object;
        let message = new $root.FilesList();
        if (object.items) {
            if (!Array.isArray(object.items))
                throw TypeError(".FilesList.items: array expected");
            message.items = [];
            for (let i = 0; i < object.items.length; ++i) {
                if (typeof object.items[i] !== "object")
                    throw TypeError(".FilesList.items: object expected");
                message.items[i] = $root.Files.fromObject(object.items[i]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from a FilesList message. Also converts values to other types if specified.
     * @function toObject
     * @memberof FilesList
     * @static
     * @param {FilesList} message FilesList
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    FilesList.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.items = [];
        if (message.items && message.items.length) {
            object.items = [];
            for (let j = 0; j < message.items.length; ++j)
                object.items[j] = $root.Files.toObject(message.items[j], options);
        }
        return object;
    };

    /**
     * Converts this FilesList to JSON.
     * @function toJSON
     * @memberof FilesList
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    FilesList.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return FilesList;
})();

export const Comment = $root.Comment = (() => {

    /**
     * Properties of a Comment.
     * @exports IComment
     * @interface IComment
     * @property {string} id Comment id
     * @property {google.protobuf.ITimestamp} date Comment date
     * @property {IUser} user Comment user
     * @property {string} body Comment body
     * @property {IFeedItem} target Comment target
     */

    /**
     * Constructs a new Comment.
     * @exports Comment
     * @classdesc Represents a Comment.
     * @implements IComment
     * @constructor
     * @param {IComment=} [properties] Properties to set
     */
    function Comment(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Comment id.
     * @member {string} id
     * @memberof Comment
     * @instance
     */
    Comment.prototype.id = "";

    /**
     * Comment date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof Comment
     * @instance
     */
    Comment.prototype.date = null;

    /**
     * Comment user.
     * @member {IUser} user
     * @memberof Comment
     * @instance
     */
    Comment.prototype.user = null;

    /**
     * Comment body.
     * @member {string} body
     * @memberof Comment
     * @instance
     */
    Comment.prototype.body = "";

    /**
     * Comment target.
     * @member {IFeedItem} target
     * @memberof Comment
     * @instance
     */
    Comment.prototype.target = null;

    /**
     * Creates a new Comment instance using the specified properties.
     * @function create
     * @memberof Comment
     * @static
     * @param {IComment=} [properties] Properties to set
     * @returns {Comment} Comment instance
     */
    Comment.create = function create(properties) {
        return new Comment(properties);
    };

    /**
     * Encodes the specified Comment message. Does not implicitly {@link Comment.verify|verify} messages.
     * @function encode
     * @memberof Comment
     * @static
     * @param {IComment} message Comment message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Comment.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        if (message.user != null && message.hasOwnProperty("user"))
            $root.User.encode(message.user, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        if (message.body != null && message.hasOwnProperty("body"))
            writer.uint32(/* id 4, wireType 2 =*/34).string(message.body);
        if (message.target != null && message.hasOwnProperty("target"))
            $root.FeedItem.encode(message.target, writer.uint32(/* id 5, wireType 2 =*/42).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified Comment message, length delimited. Does not implicitly {@link Comment.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Comment
     * @static
     * @param {IComment} message Comment message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Comment.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Comment message from the specified reader or buffer.
     * @function decode
     * @memberof Comment
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Comment} Comment
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Comment.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Comment();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 3:
                message.user = $root.User.decode(reader, reader.uint32());
                break;
            case 4:
                message.body = reader.string();
                break;
            case 5:
                message.target = $root.FeedItem.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Comment message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Comment
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Comment} Comment
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Comment.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Comment message.
     * @function verify
     * @memberof Comment
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Comment.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        if (message.user != null && message.hasOwnProperty("user")) {
            let error = $root.User.verify(message.user);
            if (error)
                return "user." + error;
        }
        if (message.body != null && message.hasOwnProperty("body"))
            if (!$util.isString(message.body))
                return "body: string expected";
        if (message.target != null && message.hasOwnProperty("target")) {
            let error = $root.FeedItem.verify(message.target);
            if (error)
                return "target." + error;
        }
        return null;
    };

    /**
     * Creates a Comment message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Comment
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Comment} Comment
     */
    Comment.fromObject = function fromObject(object) {
        if (object instanceof $root.Comment)
            return object;
        let message = new $root.Comment();
        if (object.id != null)
            message.id = String(object.id);
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".Comment.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        if (object.user != null) {
            if (typeof object.user !== "object")
                throw TypeError(".Comment.user: object expected");
            message.user = $root.User.fromObject(object.user);
        }
        if (object.body != null)
            message.body = String(object.body);
        if (object.target != null) {
            if (typeof object.target !== "object")
                throw TypeError(".Comment.target: object expected");
            message.target = $root.FeedItem.fromObject(object.target);
        }
        return message;
    };

    /**
     * Creates a plain object from a Comment message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Comment
     * @static
     * @param {Comment} message Comment
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Comment.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            object.date = null;
            object.user = null;
            object.body = "";
            object.target = null;
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        if (message.user != null && message.hasOwnProperty("user"))
            object.user = $root.User.toObject(message.user, options);
        if (message.body != null && message.hasOwnProperty("body"))
            object.body = message.body;
        if (message.target != null && message.hasOwnProperty("target"))
            object.target = $root.FeedItem.toObject(message.target, options);
        return object;
    };

    /**
     * Converts this Comment to JSON.
     * @function toJSON
     * @memberof Comment
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Comment.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Comment;
})();

export const CommentList = $root.CommentList = (() => {

    /**
     * Properties of a CommentList.
     * @exports ICommentList
     * @interface ICommentList
     * @property {Array.<IComment>} items CommentList items
     */

    /**
     * Constructs a new CommentList.
     * @exports CommentList
     * @classdesc Represents a CommentList.
     * @implements ICommentList
     * @constructor
     * @param {ICommentList=} [properties] Properties to set
     */
    function CommentList(properties) {
        this.items = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CommentList items.
     * @member {Array.<IComment>} items
     * @memberof CommentList
     * @instance
     */
    CommentList.prototype.items = $util.emptyArray;

    /**
     * Creates a new CommentList instance using the specified properties.
     * @function create
     * @memberof CommentList
     * @static
     * @param {ICommentList=} [properties] Properties to set
     * @returns {CommentList} CommentList instance
     */
    CommentList.create = function create(properties) {
        return new CommentList(properties);
    };

    /**
     * Encodes the specified CommentList message. Does not implicitly {@link CommentList.verify|verify} messages.
     * @function encode
     * @memberof CommentList
     * @static
     * @param {ICommentList} message CommentList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CommentList.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.items != null && message.items.length)
            for (let i = 0; i < message.items.length; ++i)
                $root.Comment.encode(message.items[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified CommentList message, length delimited. Does not implicitly {@link CommentList.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CommentList
     * @static
     * @param {ICommentList} message CommentList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CommentList.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CommentList message from the specified reader or buffer.
     * @function decode
     * @memberof CommentList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CommentList} CommentList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CommentList.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.CommentList();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                if (!(message.items && message.items.length))
                    message.items = [];
                message.items.push($root.Comment.decode(reader, reader.uint32()));
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CommentList message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CommentList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CommentList} CommentList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CommentList.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CommentList message.
     * @function verify
     * @memberof CommentList
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CommentList.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.items != null && message.hasOwnProperty("items")) {
            if (!Array.isArray(message.items))
                return "items: array expected";
            for (let i = 0; i < message.items.length; ++i) {
                let error = $root.Comment.verify(message.items[i]);
                if (error)
                    return "items." + error;
            }
        }
        return null;
    };

    /**
     * Creates a CommentList message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CommentList
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CommentList} CommentList
     */
    CommentList.fromObject = function fromObject(object) {
        if (object instanceof $root.CommentList)
            return object;
        let message = new $root.CommentList();
        if (object.items) {
            if (!Array.isArray(object.items))
                throw TypeError(".CommentList.items: array expected");
            message.items = [];
            for (let i = 0; i < object.items.length; ++i) {
                if (typeof object.items[i] !== "object")
                    throw TypeError(".CommentList.items: object expected");
                message.items[i] = $root.Comment.fromObject(object.items[i]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from a CommentList message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CommentList
     * @static
     * @param {CommentList} message CommentList
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CommentList.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.items = [];
        if (message.items && message.items.length) {
            object.items = [];
            for (let j = 0; j < message.items.length; ++j)
                object.items[j] = $root.Comment.toObject(message.items[j], options);
        }
        return object;
    };

    /**
     * Converts this CommentList to JSON.
     * @function toJSON
     * @memberof CommentList
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CommentList.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CommentList;
})();

export const Like = $root.Like = (() => {

    /**
     * Properties of a Like.
     * @exports ILike
     * @interface ILike
     * @property {string} id Like id
     * @property {google.protobuf.ITimestamp} date Like date
     * @property {IUser} user Like user
     * @property {IFeedItem} target Like target
     */

    /**
     * Constructs a new Like.
     * @exports Like
     * @classdesc Represents a Like.
     * @implements ILike
     * @constructor
     * @param {ILike=} [properties] Properties to set
     */
    function Like(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Like id.
     * @member {string} id
     * @memberof Like
     * @instance
     */
    Like.prototype.id = "";

    /**
     * Like date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof Like
     * @instance
     */
    Like.prototype.date = null;

    /**
     * Like user.
     * @member {IUser} user
     * @memberof Like
     * @instance
     */
    Like.prototype.user = null;

    /**
     * Like target.
     * @member {IFeedItem} target
     * @memberof Like
     * @instance
     */
    Like.prototype.target = null;

    /**
     * Creates a new Like instance using the specified properties.
     * @function create
     * @memberof Like
     * @static
     * @param {ILike=} [properties] Properties to set
     * @returns {Like} Like instance
     */
    Like.create = function create(properties) {
        return new Like(properties);
    };

    /**
     * Encodes the specified Like message. Does not implicitly {@link Like.verify|verify} messages.
     * @function encode
     * @memberof Like
     * @static
     * @param {ILike} message Like message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Like.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        if (message.user != null && message.hasOwnProperty("user"))
            $root.User.encode(message.user, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        if (message.target != null && message.hasOwnProperty("target"))
            $root.FeedItem.encode(message.target, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified Like message, length delimited. Does not implicitly {@link Like.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Like
     * @static
     * @param {ILike} message Like message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Like.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Like message from the specified reader or buffer.
     * @function decode
     * @memberof Like
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Like} Like
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Like.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Like();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 3:
                message.user = $root.User.decode(reader, reader.uint32());
                break;
            case 4:
                message.target = $root.FeedItem.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Like message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Like
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Like} Like
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Like.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Like message.
     * @function verify
     * @memberof Like
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Like.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        if (message.user != null && message.hasOwnProperty("user")) {
            let error = $root.User.verify(message.user);
            if (error)
                return "user." + error;
        }
        if (message.target != null && message.hasOwnProperty("target")) {
            let error = $root.FeedItem.verify(message.target);
            if (error)
                return "target." + error;
        }
        return null;
    };

    /**
     * Creates a Like message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Like
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Like} Like
     */
    Like.fromObject = function fromObject(object) {
        if (object instanceof $root.Like)
            return object;
        let message = new $root.Like();
        if (object.id != null)
            message.id = String(object.id);
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".Like.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        if (object.user != null) {
            if (typeof object.user !== "object")
                throw TypeError(".Like.user: object expected");
            message.user = $root.User.fromObject(object.user);
        }
        if (object.target != null) {
            if (typeof object.target !== "object")
                throw TypeError(".Like.target: object expected");
            message.target = $root.FeedItem.fromObject(object.target);
        }
        return message;
    };

    /**
     * Creates a plain object from a Like message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Like
     * @static
     * @param {Like} message Like
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Like.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            object.date = null;
            object.user = null;
            object.target = null;
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        if (message.user != null && message.hasOwnProperty("user"))
            object.user = $root.User.toObject(message.user, options);
        if (message.target != null && message.hasOwnProperty("target"))
            object.target = $root.FeedItem.toObject(message.target, options);
        return object;
    };

    /**
     * Converts this Like to JSON.
     * @function toJSON
     * @memberof Like
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Like.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Like;
})();

export const LikeList = $root.LikeList = (() => {

    /**
     * Properties of a LikeList.
     * @exports ILikeList
     * @interface ILikeList
     * @property {Array.<ILike>} items LikeList items
     */

    /**
     * Constructs a new LikeList.
     * @exports LikeList
     * @classdesc Represents a LikeList.
     * @implements ILikeList
     * @constructor
     * @param {ILikeList=} [properties] Properties to set
     */
    function LikeList(properties) {
        this.items = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * LikeList items.
     * @member {Array.<ILike>} items
     * @memberof LikeList
     * @instance
     */
    LikeList.prototype.items = $util.emptyArray;

    /**
     * Creates a new LikeList instance using the specified properties.
     * @function create
     * @memberof LikeList
     * @static
     * @param {ILikeList=} [properties] Properties to set
     * @returns {LikeList} LikeList instance
     */
    LikeList.create = function create(properties) {
        return new LikeList(properties);
    };

    /**
     * Encodes the specified LikeList message. Does not implicitly {@link LikeList.verify|verify} messages.
     * @function encode
     * @memberof LikeList
     * @static
     * @param {ILikeList} message LikeList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    LikeList.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.items != null && message.items.length)
            for (let i = 0; i < message.items.length; ++i)
                $root.Like.encode(message.items[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified LikeList message, length delimited. Does not implicitly {@link LikeList.verify|verify} messages.
     * @function encodeDelimited
     * @memberof LikeList
     * @static
     * @param {ILikeList} message LikeList message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    LikeList.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a LikeList message from the specified reader or buffer.
     * @function decode
     * @memberof LikeList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {LikeList} LikeList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    LikeList.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.LikeList();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                if (!(message.items && message.items.length))
                    message.items = [];
                message.items.push($root.Like.decode(reader, reader.uint32()));
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a LikeList message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof LikeList
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {LikeList} LikeList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    LikeList.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a LikeList message.
     * @function verify
     * @memberof LikeList
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    LikeList.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.items != null && message.hasOwnProperty("items")) {
            if (!Array.isArray(message.items))
                return "items: array expected";
            for (let i = 0; i < message.items.length; ++i) {
                let error = $root.Like.verify(message.items[i]);
                if (error)
                    return "items." + error;
            }
        }
        return null;
    };

    /**
     * Creates a LikeList message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof LikeList
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {LikeList} LikeList
     */
    LikeList.fromObject = function fromObject(object) {
        if (object instanceof $root.LikeList)
            return object;
        let message = new $root.LikeList();
        if (object.items) {
            if (!Array.isArray(object.items))
                throw TypeError(".LikeList.items: array expected");
            message.items = [];
            for (let i = 0; i < object.items.length; ++i) {
                if (typeof object.items[i] !== "object")
                    throw TypeError(".LikeList.items: object expected");
                message.items[i] = $root.Like.fromObject(object.items[i]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from a LikeList message. Also converts values to other types if specified.
     * @function toObject
     * @memberof LikeList
     * @static
     * @param {LikeList} message LikeList
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    LikeList.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.items = [];
        if (message.items && message.items.length) {
            object.items = [];
            for (let j = 0; j < message.items.length; ++j)
                object.items[j] = $root.Like.toObject(message.items[j], options);
        }
        return object;
    };

    /**
     * Converts this LikeList to JSON.
     * @function toJSON
     * @memberof LikeList
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    LikeList.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return LikeList;
})();

export const Summary = $root.Summary = (() => {

    /**
     * Properties of a Summary.
     * @exports ISummary
     * @interface ISummary
     * @property {number} accountPeerCount Summary accountPeerCount
     * @property {number} threadCount Summary threadCount
     * @property {number} fileCount Summary fileCount
     * @property {number} contactCount Summary contactCount
     */

    /**
     * Constructs a new Summary.
     * @exports Summary
     * @classdesc Represents a Summary.
     * @implements ISummary
     * @constructor
     * @param {ISummary=} [properties] Properties to set
     */
    function Summary(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Summary accountPeerCount.
     * @member {number} accountPeerCount
     * @memberof Summary
     * @instance
     */
    Summary.prototype.accountPeerCount = 0;

    /**
     * Summary threadCount.
     * @member {number} threadCount
     * @memberof Summary
     * @instance
     */
    Summary.prototype.threadCount = 0;

    /**
     * Summary fileCount.
     * @member {number} fileCount
     * @memberof Summary
     * @instance
     */
    Summary.prototype.fileCount = 0;

    /**
     * Summary contactCount.
     * @member {number} contactCount
     * @memberof Summary
     * @instance
     */
    Summary.prototype.contactCount = 0;

    /**
     * Creates a new Summary instance using the specified properties.
     * @function create
     * @memberof Summary
     * @static
     * @param {ISummary=} [properties] Properties to set
     * @returns {Summary} Summary instance
     */
    Summary.create = function create(properties) {
        return new Summary(properties);
    };

    /**
     * Encodes the specified Summary message. Does not implicitly {@link Summary.verify|verify} messages.
     * @function encode
     * @memberof Summary
     * @static
     * @param {ISummary} message Summary message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Summary.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.accountPeerCount != null && message.hasOwnProperty("accountPeerCount"))
            writer.uint32(/* id 1, wireType 0 =*/8).int32(message.accountPeerCount);
        if (message.threadCount != null && message.hasOwnProperty("threadCount"))
            writer.uint32(/* id 2, wireType 0 =*/16).int32(message.threadCount);
        if (message.fileCount != null && message.hasOwnProperty("fileCount"))
            writer.uint32(/* id 3, wireType 0 =*/24).int32(message.fileCount);
        if (message.contactCount != null && message.hasOwnProperty("contactCount"))
            writer.uint32(/* id 4, wireType 0 =*/32).int32(message.contactCount);
        return writer;
    };

    /**
     * Encodes the specified Summary message, length delimited. Does not implicitly {@link Summary.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Summary
     * @static
     * @param {ISummary} message Summary message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Summary.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Summary message from the specified reader or buffer.
     * @function decode
     * @memberof Summary
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Summary} Summary
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Summary.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Summary();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.accountPeerCount = reader.int32();
                break;
            case 2:
                message.threadCount = reader.int32();
                break;
            case 3:
                message.fileCount = reader.int32();
                break;
            case 4:
                message.contactCount = reader.int32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Summary message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Summary
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Summary} Summary
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Summary.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Summary message.
     * @function verify
     * @memberof Summary
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Summary.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.accountPeerCount != null && message.hasOwnProperty("accountPeerCount"))
            if (!$util.isInteger(message.accountPeerCount))
                return "accountPeerCount: integer expected";
        if (message.threadCount != null && message.hasOwnProperty("threadCount"))
            if (!$util.isInteger(message.threadCount))
                return "threadCount: integer expected";
        if (message.fileCount != null && message.hasOwnProperty("fileCount"))
            if (!$util.isInteger(message.fileCount))
                return "fileCount: integer expected";
        if (message.contactCount != null && message.hasOwnProperty("contactCount"))
            if (!$util.isInteger(message.contactCount))
                return "contactCount: integer expected";
        return null;
    };

    /**
     * Creates a Summary message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Summary
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Summary} Summary
     */
    Summary.fromObject = function fromObject(object) {
        if (object instanceof $root.Summary)
            return object;
        let message = new $root.Summary();
        if (object.accountPeerCount != null)
            message.accountPeerCount = object.accountPeerCount | 0;
        if (object.threadCount != null)
            message.threadCount = object.threadCount | 0;
        if (object.fileCount != null)
            message.fileCount = object.fileCount | 0;
        if (object.contactCount != null)
            message.contactCount = object.contactCount | 0;
        return message;
    };

    /**
     * Creates a plain object from a Summary message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Summary
     * @static
     * @param {Summary} message Summary
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Summary.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.accountPeerCount = 0;
            object.threadCount = 0;
            object.fileCount = 0;
            object.contactCount = 0;
        }
        if (message.accountPeerCount != null && message.hasOwnProperty("accountPeerCount"))
            object.accountPeerCount = message.accountPeerCount;
        if (message.threadCount != null && message.hasOwnProperty("threadCount"))
            object.threadCount = message.threadCount;
        if (message.fileCount != null && message.hasOwnProperty("fileCount"))
            object.fileCount = message.fileCount;
        if (message.contactCount != null && message.hasOwnProperty("contactCount"))
            object.contactCount = message.contactCount;
        return object;
    };

    /**
     * Converts this Summary to JSON.
     * @function toJSON
     * @memberof Summary
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Summary.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Summary;
})();

/**
 * QueryType enum.
 * @exports QueryType
 * @enum {string}
 * @property {number} THREAD_BACKUPS=0 THREAD_BACKUPS value
 * @property {number} CONTACTS=1 CONTACTS value
 */
$root.QueryType = (function() {
    const valuesById = {}, values = Object.create(valuesById);
    values[valuesById[0] = "THREAD_BACKUPS"] = 0;
    values[valuesById[1] = "CONTACTS"] = 1;
    return values;
})();

export const QueryOptions = $root.QueryOptions = (() => {

    /**
     * Properties of a QueryOptions.
     * @exports IQueryOptions
     * @interface IQueryOptions
     * @property {boolean} local QueryOptions local
     * @property {number} limit QueryOptions limit
     * @property {number} wait QueryOptions wait
     * @property {QueryOptions.FilterType} filter QueryOptions filter
     * @property {Array.<string>} exclude QueryOptions exclude
     */

    /**
     * Constructs a new QueryOptions.
     * @exports QueryOptions
     * @classdesc Represents a QueryOptions.
     * @implements IQueryOptions
     * @constructor
     * @param {IQueryOptions=} [properties] Properties to set
     */
    function QueryOptions(properties) {
        this.exclude = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * QueryOptions local.
     * @member {boolean} local
     * @memberof QueryOptions
     * @instance
     */
    QueryOptions.prototype.local = false;

    /**
     * QueryOptions limit.
     * @member {number} limit
     * @memberof QueryOptions
     * @instance
     */
    QueryOptions.prototype.limit = 0;

    /**
     * QueryOptions wait.
     * @member {number} wait
     * @memberof QueryOptions
     * @instance
     */
    QueryOptions.prototype.wait = 0;

    /**
     * QueryOptions filter.
     * @member {QueryOptions.FilterType} filter
     * @memberof QueryOptions
     * @instance
     */
    QueryOptions.prototype.filter = 0;

    /**
     * QueryOptions exclude.
     * @member {Array.<string>} exclude
     * @memberof QueryOptions
     * @instance
     */
    QueryOptions.prototype.exclude = $util.emptyArray;

    /**
     * Creates a new QueryOptions instance using the specified properties.
     * @function create
     * @memberof QueryOptions
     * @static
     * @param {IQueryOptions=} [properties] Properties to set
     * @returns {QueryOptions} QueryOptions instance
     */
    QueryOptions.create = function create(properties) {
        return new QueryOptions(properties);
    };

    /**
     * Encodes the specified QueryOptions message. Does not implicitly {@link QueryOptions.verify|verify} messages.
     * @function encode
     * @memberof QueryOptions
     * @static
     * @param {IQueryOptions} message QueryOptions message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    QueryOptions.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.local != null && message.hasOwnProperty("local"))
            writer.uint32(/* id 1, wireType 0 =*/8).bool(message.local);
        if (message.limit != null && message.hasOwnProperty("limit"))
            writer.uint32(/* id 2, wireType 0 =*/16).int32(message.limit);
        if (message.wait != null && message.hasOwnProperty("wait"))
            writer.uint32(/* id 3, wireType 0 =*/24).int32(message.wait);
        if (message.filter != null && message.hasOwnProperty("filter"))
            writer.uint32(/* id 4, wireType 0 =*/32).int32(message.filter);
        if (message.exclude != null && message.exclude.length)
            for (let i = 0; i < message.exclude.length; ++i)
                writer.uint32(/* id 5, wireType 2 =*/42).string(message.exclude[i]);
        return writer;
    };

    /**
     * Encodes the specified QueryOptions message, length delimited. Does not implicitly {@link QueryOptions.verify|verify} messages.
     * @function encodeDelimited
     * @memberof QueryOptions
     * @static
     * @param {IQueryOptions} message QueryOptions message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    QueryOptions.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a QueryOptions message from the specified reader or buffer.
     * @function decode
     * @memberof QueryOptions
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {QueryOptions} QueryOptions
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    QueryOptions.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.QueryOptions();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.local = reader.bool();
                break;
            case 2:
                message.limit = reader.int32();
                break;
            case 3:
                message.wait = reader.int32();
                break;
            case 4:
                message.filter = reader.int32();
                break;
            case 5:
                if (!(message.exclude && message.exclude.length))
                    message.exclude = [];
                message.exclude.push(reader.string());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a QueryOptions message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof QueryOptions
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {QueryOptions} QueryOptions
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    QueryOptions.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a QueryOptions message.
     * @function verify
     * @memberof QueryOptions
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    QueryOptions.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.local != null && message.hasOwnProperty("local"))
            if (typeof message.local !== "boolean")
                return "local: boolean expected";
        if (message.limit != null && message.hasOwnProperty("limit"))
            if (!$util.isInteger(message.limit))
                return "limit: integer expected";
        if (message.wait != null && message.hasOwnProperty("wait"))
            if (!$util.isInteger(message.wait))
                return "wait: integer expected";
        if (message.filter != null && message.hasOwnProperty("filter"))
            switch (message.filter) {
            default:
                return "filter: enum value expected";
            case 0:
            case 1:
                break;
            }
        if (message.exclude != null && message.hasOwnProperty("exclude")) {
            if (!Array.isArray(message.exclude))
                return "exclude: array expected";
            for (let i = 0; i < message.exclude.length; ++i)
                if (!$util.isString(message.exclude[i]))
                    return "exclude: string[] expected";
        }
        return null;
    };

    /**
     * Creates a QueryOptions message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof QueryOptions
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {QueryOptions} QueryOptions
     */
    QueryOptions.fromObject = function fromObject(object) {
        if (object instanceof $root.QueryOptions)
            return object;
        let message = new $root.QueryOptions();
        if (object.local != null)
            message.local = Boolean(object.local);
        if (object.limit != null)
            message.limit = object.limit | 0;
        if (object.wait != null)
            message.wait = object.wait | 0;
        switch (object.filter) {
        case "NO_FILTER":
        case 0:
            message.filter = 0;
            break;
        case "HIDE_OLDER":
        case 1:
            message.filter = 1;
            break;
        }
        if (object.exclude) {
            if (!Array.isArray(object.exclude))
                throw TypeError(".QueryOptions.exclude: array expected");
            message.exclude = [];
            for (let i = 0; i < object.exclude.length; ++i)
                message.exclude[i] = String(object.exclude[i]);
        }
        return message;
    };

    /**
     * Creates a plain object from a QueryOptions message. Also converts values to other types if specified.
     * @function toObject
     * @memberof QueryOptions
     * @static
     * @param {QueryOptions} message QueryOptions
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    QueryOptions.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.exclude = [];
        if (options.defaults) {
            object.local = false;
            object.limit = 0;
            object.wait = 0;
            object.filter = options.enums === String ? "NO_FILTER" : 0;
        }
        if (message.local != null && message.hasOwnProperty("local"))
            object.local = message.local;
        if (message.limit != null && message.hasOwnProperty("limit"))
            object.limit = message.limit;
        if (message.wait != null && message.hasOwnProperty("wait"))
            object.wait = message.wait;
        if (message.filter != null && message.hasOwnProperty("filter"))
            object.filter = options.enums === String ? $root.QueryOptions.FilterType[message.filter] : message.filter;
        if (message.exclude && message.exclude.length) {
            object.exclude = [];
            for (let j = 0; j < message.exclude.length; ++j)
                object.exclude[j] = message.exclude[j];
        }
        return object;
    };

    /**
     * Converts this QueryOptions to JSON.
     * @function toJSON
     * @memberof QueryOptions
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    QueryOptions.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    /**
     * FilterType enum.
     * @name QueryOptions.FilterType
     * @enum {string}
     * @property {number} NO_FILTER=0 NO_FILTER value
     * @property {number} HIDE_OLDER=1 HIDE_OLDER value
     */
    QueryOptions.FilterType = (function() {
        const valuesById = {}, values = Object.create(valuesById);
        values[valuesById[0] = "NO_FILTER"] = 0;
        values[valuesById[1] = "HIDE_OLDER"] = 1;
        return values;
    })();

    return QueryOptions;
})();

export const Query = $root.Query = (() => {

    /**
     * Properties of a Query.
     * @exports IQuery
     * @interface IQuery
     * @property {string} id Query id
     * @property {string} token Query token
     * @property {QueryType} type Query type
     * @property {IQueryOptions} options Query options
     * @property {google.protobuf.IAny} payload Query payload
     */

    /**
     * Constructs a new Query.
     * @exports Query
     * @classdesc Represents a Query.
     * @implements IQuery
     * @constructor
     * @param {IQuery=} [properties] Properties to set
     */
    function Query(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Query id.
     * @member {string} id
     * @memberof Query
     * @instance
     */
    Query.prototype.id = "";

    /**
     * Query token.
     * @member {string} token
     * @memberof Query
     * @instance
     */
    Query.prototype.token = "";

    /**
     * Query type.
     * @member {QueryType} type
     * @memberof Query
     * @instance
     */
    Query.prototype.type = 0;

    /**
     * Query options.
     * @member {IQueryOptions} options
     * @memberof Query
     * @instance
     */
    Query.prototype.options = null;

    /**
     * Query payload.
     * @member {google.protobuf.IAny} payload
     * @memberof Query
     * @instance
     */
    Query.prototype.payload = null;

    /**
     * Creates a new Query instance using the specified properties.
     * @function create
     * @memberof Query
     * @static
     * @param {IQuery=} [properties] Properties to set
     * @returns {Query} Query instance
     */
    Query.create = function create(properties) {
        return new Query(properties);
    };

    /**
     * Encodes the specified Query message. Does not implicitly {@link Query.verify|verify} messages.
     * @function encode
     * @memberof Query
     * @static
     * @param {IQuery} message Query message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Query.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.token != null && message.hasOwnProperty("token"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.token);
        if (message.type != null && message.hasOwnProperty("type"))
            writer.uint32(/* id 3, wireType 0 =*/24).int32(message.type);
        if (message.options != null && message.hasOwnProperty("options"))
            $root.QueryOptions.encode(message.options, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
        if (message.payload != null && message.hasOwnProperty("payload"))
            $root.google.protobuf.Any.encode(message.payload, writer.uint32(/* id 5, wireType 2 =*/42).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified Query message, length delimited. Does not implicitly {@link Query.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Query
     * @static
     * @param {IQuery} message Query message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Query.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Query message from the specified reader or buffer.
     * @function decode
     * @memberof Query
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Query} Query
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Query.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Query();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.token = reader.string();
                break;
            case 3:
                message.type = reader.int32();
                break;
            case 4:
                message.options = $root.QueryOptions.decode(reader, reader.uint32());
                break;
            case 5:
                message.payload = $root.google.protobuf.Any.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Query message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Query
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Query} Query
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Query.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Query message.
     * @function verify
     * @memberof Query
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Query.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.token != null && message.hasOwnProperty("token"))
            if (!$util.isString(message.token))
                return "token: string expected";
        if (message.type != null && message.hasOwnProperty("type"))
            switch (message.type) {
            default:
                return "type: enum value expected";
            case 0:
            case 1:
                break;
            }
        if (message.options != null && message.hasOwnProperty("options")) {
            let error = $root.QueryOptions.verify(message.options);
            if (error)
                return "options." + error;
        }
        if (message.payload != null && message.hasOwnProperty("payload")) {
            let error = $root.google.protobuf.Any.verify(message.payload);
            if (error)
                return "payload." + error;
        }
        return null;
    };

    /**
     * Creates a Query message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Query
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Query} Query
     */
    Query.fromObject = function fromObject(object) {
        if (object instanceof $root.Query)
            return object;
        let message = new $root.Query();
        if (object.id != null)
            message.id = String(object.id);
        if (object.token != null)
            message.token = String(object.token);
        switch (object.type) {
        case "THREAD_BACKUPS":
        case 0:
            message.type = 0;
            break;
        case "CONTACTS":
        case 1:
            message.type = 1;
            break;
        }
        if (object.options != null) {
            if (typeof object.options !== "object")
                throw TypeError(".Query.options: object expected");
            message.options = $root.QueryOptions.fromObject(object.options);
        }
        if (object.payload != null) {
            if (typeof object.payload !== "object")
                throw TypeError(".Query.payload: object expected");
            message.payload = $root.google.protobuf.Any.fromObject(object.payload);
        }
        return message;
    };

    /**
     * Creates a plain object from a Query message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Query
     * @static
     * @param {Query} message Query
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Query.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            object.token = "";
            object.type = options.enums === String ? "THREAD_BACKUPS" : 0;
            object.options = null;
            object.payload = null;
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.token != null && message.hasOwnProperty("token"))
            object.token = message.token;
        if (message.type != null && message.hasOwnProperty("type"))
            object.type = options.enums === String ? $root.QueryType[message.type] : message.type;
        if (message.options != null && message.hasOwnProperty("options"))
            object.options = $root.QueryOptions.toObject(message.options, options);
        if (message.payload != null && message.hasOwnProperty("payload"))
            object.payload = $root.google.protobuf.Any.toObject(message.payload, options);
        return object;
    };

    /**
     * Converts this Query to JSON.
     * @function toJSON
     * @memberof Query
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Query.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Query;
})();

export const PubSubQuery = $root.PubSubQuery = (() => {

    /**
     * Properties of a PubSubQuery.
     * @exports IPubSubQuery
     * @interface IPubSubQuery
     * @property {string} id PubSubQuery id
     * @property {QueryType} type PubSubQuery type
     * @property {google.protobuf.IAny} payload PubSubQuery payload
     * @property {PubSubQuery.ResponseType} responseType PubSubQuery responseType
     */

    /**
     * Constructs a new PubSubQuery.
     * @exports PubSubQuery
     * @classdesc Represents a PubSubQuery.
     * @implements IPubSubQuery
     * @constructor
     * @param {IPubSubQuery=} [properties] Properties to set
     */
    function PubSubQuery(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * PubSubQuery id.
     * @member {string} id
     * @memberof PubSubQuery
     * @instance
     */
    PubSubQuery.prototype.id = "";

    /**
     * PubSubQuery type.
     * @member {QueryType} type
     * @memberof PubSubQuery
     * @instance
     */
    PubSubQuery.prototype.type = 0;

    /**
     * PubSubQuery payload.
     * @member {google.protobuf.IAny} payload
     * @memberof PubSubQuery
     * @instance
     */
    PubSubQuery.prototype.payload = null;

    /**
     * PubSubQuery responseType.
     * @member {PubSubQuery.ResponseType} responseType
     * @memberof PubSubQuery
     * @instance
     */
    PubSubQuery.prototype.responseType = 0;

    /**
     * Creates a new PubSubQuery instance using the specified properties.
     * @function create
     * @memberof PubSubQuery
     * @static
     * @param {IPubSubQuery=} [properties] Properties to set
     * @returns {PubSubQuery} PubSubQuery instance
     */
    PubSubQuery.create = function create(properties) {
        return new PubSubQuery(properties);
    };

    /**
     * Encodes the specified PubSubQuery message. Does not implicitly {@link PubSubQuery.verify|verify} messages.
     * @function encode
     * @memberof PubSubQuery
     * @static
     * @param {IPubSubQuery} message PubSubQuery message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PubSubQuery.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.type != null && message.hasOwnProperty("type"))
            writer.uint32(/* id 2, wireType 0 =*/16).int32(message.type);
        if (message.payload != null && message.hasOwnProperty("payload"))
            $root.google.protobuf.Any.encode(message.payload, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        if (message.responseType != null && message.hasOwnProperty("responseType"))
            writer.uint32(/* id 4, wireType 0 =*/32).int32(message.responseType);
        return writer;
    };

    /**
     * Encodes the specified PubSubQuery message, length delimited. Does not implicitly {@link PubSubQuery.verify|verify} messages.
     * @function encodeDelimited
     * @memberof PubSubQuery
     * @static
     * @param {IPubSubQuery} message PubSubQuery message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PubSubQuery.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a PubSubQuery message from the specified reader or buffer.
     * @function decode
     * @memberof PubSubQuery
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {PubSubQuery} PubSubQuery
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PubSubQuery.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.PubSubQuery();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.type = reader.int32();
                break;
            case 3:
                message.payload = $root.google.protobuf.Any.decode(reader, reader.uint32());
                break;
            case 4:
                message.responseType = reader.int32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a PubSubQuery message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof PubSubQuery
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {PubSubQuery} PubSubQuery
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PubSubQuery.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a PubSubQuery message.
     * @function verify
     * @memberof PubSubQuery
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    PubSubQuery.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.type != null && message.hasOwnProperty("type"))
            switch (message.type) {
            default:
                return "type: enum value expected";
            case 0:
            case 1:
                break;
            }
        if (message.payload != null && message.hasOwnProperty("payload")) {
            let error = $root.google.protobuf.Any.verify(message.payload);
            if (error)
                return "payload." + error;
        }
        if (message.responseType != null && message.hasOwnProperty("responseType"))
            switch (message.responseType) {
            default:
                return "responseType: enum value expected";
            case 0:
            case 1:
                break;
            }
        return null;
    };

    /**
     * Creates a PubSubQuery message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof PubSubQuery
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {PubSubQuery} PubSubQuery
     */
    PubSubQuery.fromObject = function fromObject(object) {
        if (object instanceof $root.PubSubQuery)
            return object;
        let message = new $root.PubSubQuery();
        if (object.id != null)
            message.id = String(object.id);
        switch (object.type) {
        case "THREAD_BACKUPS":
        case 0:
            message.type = 0;
            break;
        case "CONTACTS":
        case 1:
            message.type = 1;
            break;
        }
        if (object.payload != null) {
            if (typeof object.payload !== "object")
                throw TypeError(".PubSubQuery.payload: object expected");
            message.payload = $root.google.protobuf.Any.fromObject(object.payload);
        }
        switch (object.responseType) {
        case "P2P":
        case 0:
            message.responseType = 0;
            break;
        case "PUBSUB":
        case 1:
            message.responseType = 1;
            break;
        }
        return message;
    };

    /**
     * Creates a plain object from a PubSubQuery message. Also converts values to other types if specified.
     * @function toObject
     * @memberof PubSubQuery
     * @static
     * @param {PubSubQuery} message PubSubQuery
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    PubSubQuery.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            object.type = options.enums === String ? "THREAD_BACKUPS" : 0;
            object.payload = null;
            object.responseType = options.enums === String ? "P2P" : 0;
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.type != null && message.hasOwnProperty("type"))
            object.type = options.enums === String ? $root.QueryType[message.type] : message.type;
        if (message.payload != null && message.hasOwnProperty("payload"))
            object.payload = $root.google.protobuf.Any.toObject(message.payload, options);
        if (message.responseType != null && message.hasOwnProperty("responseType"))
            object.responseType = options.enums === String ? $root.PubSubQuery.ResponseType[message.responseType] : message.responseType;
        return object;
    };

    /**
     * Converts this PubSubQuery to JSON.
     * @function toJSON
     * @memberof PubSubQuery
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    PubSubQuery.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    /**
     * ResponseType enum.
     * @name PubSubQuery.ResponseType
     * @enum {string}
     * @property {number} P2P=0 P2P value
     * @property {number} PUBSUB=1 PUBSUB value
     */
    PubSubQuery.ResponseType = (function() {
        const valuesById = {}, values = Object.create(valuesById);
        values[valuesById[0] = "P2P"] = 0;
        values[valuesById[1] = "PUBSUB"] = 1;
        return values;
    })();

    return PubSubQuery;
})();

export const QueryResult = $root.QueryResult = (() => {

    /**
     * Properties of a QueryResult.
     * @exports IQueryResult
     * @interface IQueryResult
     * @property {string} id QueryResult id
     * @property {google.protobuf.ITimestamp} date QueryResult date
     * @property {boolean} local QueryResult local
     * @property {google.protobuf.IAny} value QueryResult value
     */

    /**
     * Constructs a new QueryResult.
     * @exports QueryResult
     * @classdesc Represents a QueryResult.
     * @implements IQueryResult
     * @constructor
     * @param {IQueryResult=} [properties] Properties to set
     */
    function QueryResult(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * QueryResult id.
     * @member {string} id
     * @memberof QueryResult
     * @instance
     */
    QueryResult.prototype.id = "";

    /**
     * QueryResult date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof QueryResult
     * @instance
     */
    QueryResult.prototype.date = null;

    /**
     * QueryResult local.
     * @member {boolean} local
     * @memberof QueryResult
     * @instance
     */
    QueryResult.prototype.local = false;

    /**
     * QueryResult value.
     * @member {google.protobuf.IAny} value
     * @memberof QueryResult
     * @instance
     */
    QueryResult.prototype.value = null;

    /**
     * Creates a new QueryResult instance using the specified properties.
     * @function create
     * @memberof QueryResult
     * @static
     * @param {IQueryResult=} [properties] Properties to set
     * @returns {QueryResult} QueryResult instance
     */
    QueryResult.create = function create(properties) {
        return new QueryResult(properties);
    };

    /**
     * Encodes the specified QueryResult message. Does not implicitly {@link QueryResult.verify|verify} messages.
     * @function encode
     * @memberof QueryResult
     * @static
     * @param {IQueryResult} message QueryResult message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    QueryResult.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        if (message.local != null && message.hasOwnProperty("local"))
            writer.uint32(/* id 3, wireType 0 =*/24).bool(message.local);
        if (message.value != null && message.hasOwnProperty("value"))
            $root.google.protobuf.Any.encode(message.value, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified QueryResult message, length delimited. Does not implicitly {@link QueryResult.verify|verify} messages.
     * @function encodeDelimited
     * @memberof QueryResult
     * @static
     * @param {IQueryResult} message QueryResult message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    QueryResult.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a QueryResult message from the specified reader or buffer.
     * @function decode
     * @memberof QueryResult
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {QueryResult} QueryResult
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    QueryResult.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.QueryResult();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 3:
                message.local = reader.bool();
                break;
            case 4:
                message.value = $root.google.protobuf.Any.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a QueryResult message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof QueryResult
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {QueryResult} QueryResult
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    QueryResult.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a QueryResult message.
     * @function verify
     * @memberof QueryResult
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    QueryResult.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        if (message.local != null && message.hasOwnProperty("local"))
            if (typeof message.local !== "boolean")
                return "local: boolean expected";
        if (message.value != null && message.hasOwnProperty("value")) {
            let error = $root.google.protobuf.Any.verify(message.value);
            if (error)
                return "value." + error;
        }
        return null;
    };

    /**
     * Creates a QueryResult message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof QueryResult
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {QueryResult} QueryResult
     */
    QueryResult.fromObject = function fromObject(object) {
        if (object instanceof $root.QueryResult)
            return object;
        let message = new $root.QueryResult();
        if (object.id != null)
            message.id = String(object.id);
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".QueryResult.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        if (object.local != null)
            message.local = Boolean(object.local);
        if (object.value != null) {
            if (typeof object.value !== "object")
                throw TypeError(".QueryResult.value: object expected");
            message.value = $root.google.protobuf.Any.fromObject(object.value);
        }
        return message;
    };

    /**
     * Creates a plain object from a QueryResult message. Also converts values to other types if specified.
     * @function toObject
     * @memberof QueryResult
     * @static
     * @param {QueryResult} message QueryResult
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    QueryResult.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            object.date = null;
            object.local = false;
            object.value = null;
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        if (message.local != null && message.hasOwnProperty("local"))
            object.local = message.local;
        if (message.value != null && message.hasOwnProperty("value"))
            object.value = $root.google.protobuf.Any.toObject(message.value, options);
        return object;
    };

    /**
     * Converts this QueryResult to JSON.
     * @function toJSON
     * @memberof QueryResult
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    QueryResult.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return QueryResult;
})();

export const QueryResults = $root.QueryResults = (() => {

    /**
     * Properties of a QueryResults.
     * @exports IQueryResults
     * @interface IQueryResults
     * @property {QueryType} type QueryResults type
     * @property {Array.<IQueryResult>} items QueryResults items
     */

    /**
     * Constructs a new QueryResults.
     * @exports QueryResults
     * @classdesc Represents a QueryResults.
     * @implements IQueryResults
     * @constructor
     * @param {IQueryResults=} [properties] Properties to set
     */
    function QueryResults(properties) {
        this.items = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * QueryResults type.
     * @member {QueryType} type
     * @memberof QueryResults
     * @instance
     */
    QueryResults.prototype.type = 0;

    /**
     * QueryResults items.
     * @member {Array.<IQueryResult>} items
     * @memberof QueryResults
     * @instance
     */
    QueryResults.prototype.items = $util.emptyArray;

    /**
     * Creates a new QueryResults instance using the specified properties.
     * @function create
     * @memberof QueryResults
     * @static
     * @param {IQueryResults=} [properties] Properties to set
     * @returns {QueryResults} QueryResults instance
     */
    QueryResults.create = function create(properties) {
        return new QueryResults(properties);
    };

    /**
     * Encodes the specified QueryResults message. Does not implicitly {@link QueryResults.verify|verify} messages.
     * @function encode
     * @memberof QueryResults
     * @static
     * @param {IQueryResults} message QueryResults message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    QueryResults.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.type != null && message.hasOwnProperty("type"))
            writer.uint32(/* id 1, wireType 0 =*/8).int32(message.type);
        if (message.items != null && message.items.length)
            for (let i = 0; i < message.items.length; ++i)
                $root.QueryResult.encode(message.items[i], writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified QueryResults message, length delimited. Does not implicitly {@link QueryResults.verify|verify} messages.
     * @function encodeDelimited
     * @memberof QueryResults
     * @static
     * @param {IQueryResults} message QueryResults message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    QueryResults.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a QueryResults message from the specified reader or buffer.
     * @function decode
     * @memberof QueryResults
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {QueryResults} QueryResults
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    QueryResults.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.QueryResults();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.type = reader.int32();
                break;
            case 2:
                if (!(message.items && message.items.length))
                    message.items = [];
                message.items.push($root.QueryResult.decode(reader, reader.uint32()));
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a QueryResults message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof QueryResults
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {QueryResults} QueryResults
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    QueryResults.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a QueryResults message.
     * @function verify
     * @memberof QueryResults
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    QueryResults.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.type != null && message.hasOwnProperty("type"))
            switch (message.type) {
            default:
                return "type: enum value expected";
            case 0:
            case 1:
                break;
            }
        if (message.items != null && message.hasOwnProperty("items")) {
            if (!Array.isArray(message.items))
                return "items: array expected";
            for (let i = 0; i < message.items.length; ++i) {
                let error = $root.QueryResult.verify(message.items[i]);
                if (error)
                    return "items." + error;
            }
        }
        return null;
    };

    /**
     * Creates a QueryResults message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof QueryResults
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {QueryResults} QueryResults
     */
    QueryResults.fromObject = function fromObject(object) {
        if (object instanceof $root.QueryResults)
            return object;
        let message = new $root.QueryResults();
        switch (object.type) {
        case "THREAD_BACKUPS":
        case 0:
            message.type = 0;
            break;
        case "CONTACTS":
        case 1:
            message.type = 1;
            break;
        }
        if (object.items) {
            if (!Array.isArray(object.items))
                throw TypeError(".QueryResults.items: array expected");
            message.items = [];
            for (let i = 0; i < object.items.length; ++i) {
                if (typeof object.items[i] !== "object")
                    throw TypeError(".QueryResults.items: object expected");
                message.items[i] = $root.QueryResult.fromObject(object.items[i]);
            }
        }
        return message;
    };

    /**
     * Creates a plain object from a QueryResults message. Also converts values to other types if specified.
     * @function toObject
     * @memberof QueryResults
     * @static
     * @param {QueryResults} message QueryResults
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    QueryResults.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.items = [];
        if (options.defaults)
            object.type = options.enums === String ? "THREAD_BACKUPS" : 0;
        if (message.type != null && message.hasOwnProperty("type"))
            object.type = options.enums === String ? $root.QueryType[message.type] : message.type;
        if (message.items && message.items.length) {
            object.items = [];
            for (let j = 0; j < message.items.length; ++j)
                object.items[j] = $root.QueryResult.toObject(message.items[j], options);
        }
        return object;
    };

    /**
     * Converts this QueryResults to JSON.
     * @function toJSON
     * @memberof QueryResults
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    QueryResults.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return QueryResults;
})();

export const PubSubQueryResults = $root.PubSubQueryResults = (() => {

    /**
     * Properties of a PubSubQueryResults.
     * @exports IPubSubQueryResults
     * @interface IPubSubQueryResults
     * @property {string} id PubSubQueryResults id
     * @property {IQueryResults} results PubSubQueryResults results
     */

    /**
     * Constructs a new PubSubQueryResults.
     * @exports PubSubQueryResults
     * @classdesc Represents a PubSubQueryResults.
     * @implements IPubSubQueryResults
     * @constructor
     * @param {IPubSubQueryResults=} [properties] Properties to set
     */
    function PubSubQueryResults(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * PubSubQueryResults id.
     * @member {string} id
     * @memberof PubSubQueryResults
     * @instance
     */
    PubSubQueryResults.prototype.id = "";

    /**
     * PubSubQueryResults results.
     * @member {IQueryResults} results
     * @memberof PubSubQueryResults
     * @instance
     */
    PubSubQueryResults.prototype.results = null;

    /**
     * Creates a new PubSubQueryResults instance using the specified properties.
     * @function create
     * @memberof PubSubQueryResults
     * @static
     * @param {IPubSubQueryResults=} [properties] Properties to set
     * @returns {PubSubQueryResults} PubSubQueryResults instance
     */
    PubSubQueryResults.create = function create(properties) {
        return new PubSubQueryResults(properties);
    };

    /**
     * Encodes the specified PubSubQueryResults message. Does not implicitly {@link PubSubQueryResults.verify|verify} messages.
     * @function encode
     * @memberof PubSubQueryResults
     * @static
     * @param {IPubSubQueryResults} message PubSubQueryResults message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PubSubQueryResults.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.results != null && message.hasOwnProperty("results"))
            $root.QueryResults.encode(message.results, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified PubSubQueryResults message, length delimited. Does not implicitly {@link PubSubQueryResults.verify|verify} messages.
     * @function encodeDelimited
     * @memberof PubSubQueryResults
     * @static
     * @param {IPubSubQueryResults} message PubSubQueryResults message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PubSubQueryResults.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a PubSubQueryResults message from the specified reader or buffer.
     * @function decode
     * @memberof PubSubQueryResults
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {PubSubQueryResults} PubSubQueryResults
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PubSubQueryResults.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.PubSubQueryResults();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.results = $root.QueryResults.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a PubSubQueryResults message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof PubSubQueryResults
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {PubSubQueryResults} PubSubQueryResults
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PubSubQueryResults.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a PubSubQueryResults message.
     * @function verify
     * @memberof PubSubQueryResults
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    PubSubQueryResults.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.results != null && message.hasOwnProperty("results")) {
            let error = $root.QueryResults.verify(message.results);
            if (error)
                return "results." + error;
        }
        return null;
    };

    /**
     * Creates a PubSubQueryResults message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof PubSubQueryResults
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {PubSubQueryResults} PubSubQueryResults
     */
    PubSubQueryResults.fromObject = function fromObject(object) {
        if (object instanceof $root.PubSubQueryResults)
            return object;
        let message = new $root.PubSubQueryResults();
        if (object.id != null)
            message.id = String(object.id);
        if (object.results != null) {
            if (typeof object.results !== "object")
                throw TypeError(".PubSubQueryResults.results: object expected");
            message.results = $root.QueryResults.fromObject(object.results);
        }
        return message;
    };

    /**
     * Creates a plain object from a PubSubQueryResults message. Also converts values to other types if specified.
     * @function toObject
     * @memberof PubSubQueryResults
     * @static
     * @param {PubSubQueryResults} message PubSubQueryResults
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    PubSubQueryResults.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            object.results = null;
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.results != null && message.hasOwnProperty("results"))
            object.results = $root.QueryResults.toObject(message.results, options);
        return object;
    };

    /**
     * Converts this PubSubQueryResults to JSON.
     * @function toJSON
     * @memberof PubSubQueryResults
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    PubSubQueryResults.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return PubSubQueryResults;
})();

export const QueryEvent = $root.QueryEvent = (() => {

    /**
     * Properties of a QueryEvent.
     * @exports IQueryEvent
     * @interface IQueryEvent
     * @property {QueryEvent.Type} type QueryEvent type
     * @property {IQueryResult} data QueryEvent data
     */

    /**
     * Constructs a new QueryEvent.
     * @exports QueryEvent
     * @classdesc Represents a QueryEvent.
     * @implements IQueryEvent
     * @constructor
     * @param {IQueryEvent=} [properties] Properties to set
     */
    function QueryEvent(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * QueryEvent type.
     * @member {QueryEvent.Type} type
     * @memberof QueryEvent
     * @instance
     */
    QueryEvent.prototype.type = 0;

    /**
     * QueryEvent data.
     * @member {IQueryResult} data
     * @memberof QueryEvent
     * @instance
     */
    QueryEvent.prototype.data = null;

    /**
     * Creates a new QueryEvent instance using the specified properties.
     * @function create
     * @memberof QueryEvent
     * @static
     * @param {IQueryEvent=} [properties] Properties to set
     * @returns {QueryEvent} QueryEvent instance
     */
    QueryEvent.create = function create(properties) {
        return new QueryEvent(properties);
    };

    /**
     * Encodes the specified QueryEvent message. Does not implicitly {@link QueryEvent.verify|verify} messages.
     * @function encode
     * @memberof QueryEvent
     * @static
     * @param {IQueryEvent} message QueryEvent message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    QueryEvent.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.type != null && message.hasOwnProperty("type"))
            writer.uint32(/* id 1, wireType 0 =*/8).int32(message.type);
        if (message.data != null && message.hasOwnProperty("data"))
            $root.QueryResult.encode(message.data, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified QueryEvent message, length delimited. Does not implicitly {@link QueryEvent.verify|verify} messages.
     * @function encodeDelimited
     * @memberof QueryEvent
     * @static
     * @param {IQueryEvent} message QueryEvent message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    QueryEvent.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a QueryEvent message from the specified reader or buffer.
     * @function decode
     * @memberof QueryEvent
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {QueryEvent} QueryEvent
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    QueryEvent.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.QueryEvent();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.type = reader.int32();
                break;
            case 2:
                message.data = $root.QueryResult.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a QueryEvent message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof QueryEvent
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {QueryEvent} QueryEvent
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    QueryEvent.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a QueryEvent message.
     * @function verify
     * @memberof QueryEvent
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    QueryEvent.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.type != null && message.hasOwnProperty("type"))
            switch (message.type) {
            default:
                return "type: enum value expected";
            case 0:
            case 1:
                break;
            }
        if (message.data != null && message.hasOwnProperty("data")) {
            let error = $root.QueryResult.verify(message.data);
            if (error)
                return "data." + error;
        }
        return null;
    };

    /**
     * Creates a QueryEvent message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof QueryEvent
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {QueryEvent} QueryEvent
     */
    QueryEvent.fromObject = function fromObject(object) {
        if (object instanceof $root.QueryEvent)
            return object;
        let message = new $root.QueryEvent();
        switch (object.type) {
        case "DATA":
        case 0:
            message.type = 0;
            break;
        case "DONE":
        case 1:
            message.type = 1;
            break;
        }
        if (object.data != null) {
            if (typeof object.data !== "object")
                throw TypeError(".QueryEvent.data: object expected");
            message.data = $root.QueryResult.fromObject(object.data);
        }
        return message;
    };

    /**
     * Creates a plain object from a QueryEvent message. Also converts values to other types if specified.
     * @function toObject
     * @memberof QueryEvent
     * @static
     * @param {QueryEvent} message QueryEvent
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    QueryEvent.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.type = options.enums === String ? "DATA" : 0;
            object.data = null;
        }
        if (message.type != null && message.hasOwnProperty("type"))
            object.type = options.enums === String ? $root.QueryEvent.Type[message.type] : message.type;
        if (message.data != null && message.hasOwnProperty("data"))
            object.data = $root.QueryResult.toObject(message.data, options);
        return object;
    };

    /**
     * Converts this QueryEvent to JSON.
     * @function toJSON
     * @memberof QueryEvent
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    QueryEvent.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    /**
     * Type enum.
     * @name QueryEvent.Type
     * @enum {string}
     * @property {number} DATA=0 DATA value
     * @property {number} DONE=1 DONE value
     */
    QueryEvent.Type = (function() {
        const valuesById = {}, values = Object.create(valuesById);
        values[valuesById[0] = "DATA"] = 0;
        values[valuesById[1] = "DONE"] = 1;
        return values;
    })();

    return QueryEvent;
})();

export const ContactQuery = $root.ContactQuery = (() => {

    /**
     * Properties of a ContactQuery.
     * @exports IContactQuery
     * @interface IContactQuery
     * @property {string} id ContactQuery id
     * @property {string} address ContactQuery address
     * @property {string} username ContactQuery username
     */

    /**
     * Constructs a new ContactQuery.
     * @exports ContactQuery
     * @classdesc Represents a ContactQuery.
     * @implements IContactQuery
     * @constructor
     * @param {IContactQuery=} [properties] Properties to set
     */
    function ContactQuery(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ContactQuery id.
     * @member {string} id
     * @memberof ContactQuery
     * @instance
     */
    ContactQuery.prototype.id = "";

    /**
     * ContactQuery address.
     * @member {string} address
     * @memberof ContactQuery
     * @instance
     */
    ContactQuery.prototype.address = "";

    /**
     * ContactQuery username.
     * @member {string} username
     * @memberof ContactQuery
     * @instance
     */
    ContactQuery.prototype.username = "";

    /**
     * Creates a new ContactQuery instance using the specified properties.
     * @function create
     * @memberof ContactQuery
     * @static
     * @param {IContactQuery=} [properties] Properties to set
     * @returns {ContactQuery} ContactQuery instance
     */
    ContactQuery.create = function create(properties) {
        return new ContactQuery(properties);
    };

    /**
     * Encodes the specified ContactQuery message. Does not implicitly {@link ContactQuery.verify|verify} messages.
     * @function encode
     * @memberof ContactQuery
     * @static
     * @param {IContactQuery} message ContactQuery message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ContactQuery.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.id != null && message.hasOwnProperty("id"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
        if (message.address != null && message.hasOwnProperty("address"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.address);
        if (message.username != null && message.hasOwnProperty("username"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.username);
        return writer;
    };

    /**
     * Encodes the specified ContactQuery message, length delimited. Does not implicitly {@link ContactQuery.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ContactQuery
     * @static
     * @param {IContactQuery} message ContactQuery message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ContactQuery.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ContactQuery message from the specified reader or buffer.
     * @function decode
     * @memberof ContactQuery
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ContactQuery} ContactQuery
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ContactQuery.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ContactQuery();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.id = reader.string();
                break;
            case 2:
                message.address = reader.string();
                break;
            case 3:
                message.username = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ContactQuery message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ContactQuery
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ContactQuery} ContactQuery
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ContactQuery.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ContactQuery message.
     * @function verify
     * @memberof ContactQuery
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ContactQuery.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.id != null && message.hasOwnProperty("id"))
            if (!$util.isString(message.id))
                return "id: string expected";
        if (message.address != null && message.hasOwnProperty("address"))
            if (!$util.isString(message.address))
                return "address: string expected";
        if (message.username != null && message.hasOwnProperty("username"))
            if (!$util.isString(message.username))
                return "username: string expected";
        return null;
    };

    /**
     * Creates a ContactQuery message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ContactQuery
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ContactQuery} ContactQuery
     */
    ContactQuery.fromObject = function fromObject(object) {
        if (object instanceof $root.ContactQuery)
            return object;
        let message = new $root.ContactQuery();
        if (object.id != null)
            message.id = String(object.id);
        if (object.address != null)
            message.address = String(object.address);
        if (object.username != null)
            message.username = String(object.username);
        return message;
    };

    /**
     * Creates a plain object from a ContactQuery message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ContactQuery
     * @static
     * @param {ContactQuery} message ContactQuery
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ContactQuery.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.id = "";
            object.address = "";
            object.username = "";
        }
        if (message.id != null && message.hasOwnProperty("id"))
            object.id = message.id;
        if (message.address != null && message.hasOwnProperty("address"))
            object.address = message.address;
        if (message.username != null && message.hasOwnProperty("username"))
            object.username = message.username;
        return object;
    };

    /**
     * Converts this ContactQuery to JSON.
     * @function toJSON
     * @memberof ContactQuery
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ContactQuery.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ContactQuery;
})();

export const ThreadBackupQuery = $root.ThreadBackupQuery = (() => {

    /**
     * Properties of a ThreadBackupQuery.
     * @exports IThreadBackupQuery
     * @interface IThreadBackupQuery
     * @property {string} address ThreadBackupQuery address
     */

    /**
     * Constructs a new ThreadBackupQuery.
     * @exports ThreadBackupQuery
     * @classdesc Represents a ThreadBackupQuery.
     * @implements IThreadBackupQuery
     * @constructor
     * @param {IThreadBackupQuery=} [properties] Properties to set
     */
    function ThreadBackupQuery(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ThreadBackupQuery address.
     * @member {string} address
     * @memberof ThreadBackupQuery
     * @instance
     */
    ThreadBackupQuery.prototype.address = "";

    /**
     * Creates a new ThreadBackupQuery instance using the specified properties.
     * @function create
     * @memberof ThreadBackupQuery
     * @static
     * @param {IThreadBackupQuery=} [properties] Properties to set
     * @returns {ThreadBackupQuery} ThreadBackupQuery instance
     */
    ThreadBackupQuery.create = function create(properties) {
        return new ThreadBackupQuery(properties);
    };

    /**
     * Encodes the specified ThreadBackupQuery message. Does not implicitly {@link ThreadBackupQuery.verify|verify} messages.
     * @function encode
     * @memberof ThreadBackupQuery
     * @static
     * @param {IThreadBackupQuery} message ThreadBackupQuery message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadBackupQuery.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.address != null && message.hasOwnProperty("address"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.address);
        return writer;
    };

    /**
     * Encodes the specified ThreadBackupQuery message, length delimited. Does not implicitly {@link ThreadBackupQuery.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ThreadBackupQuery
     * @static
     * @param {IThreadBackupQuery} message ThreadBackupQuery message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadBackupQuery.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ThreadBackupQuery message from the specified reader or buffer.
     * @function decode
     * @memberof ThreadBackupQuery
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ThreadBackupQuery} ThreadBackupQuery
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadBackupQuery.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ThreadBackupQuery();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.address = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ThreadBackupQuery message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ThreadBackupQuery
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ThreadBackupQuery} ThreadBackupQuery
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadBackupQuery.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ThreadBackupQuery message.
     * @function verify
     * @memberof ThreadBackupQuery
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ThreadBackupQuery.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.address != null && message.hasOwnProperty("address"))
            if (!$util.isString(message.address))
                return "address: string expected";
        return null;
    };

    /**
     * Creates a ThreadBackupQuery message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ThreadBackupQuery
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ThreadBackupQuery} ThreadBackupQuery
     */
    ThreadBackupQuery.fromObject = function fromObject(object) {
        if (object instanceof $root.ThreadBackupQuery)
            return object;
        let message = new $root.ThreadBackupQuery();
        if (object.address != null)
            message.address = String(object.address);
        return message;
    };

    /**
     * Creates a plain object from a ThreadBackupQuery message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ThreadBackupQuery
     * @static
     * @param {ThreadBackupQuery} message ThreadBackupQuery
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ThreadBackupQuery.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults)
            object.address = "";
        if (message.address != null && message.hasOwnProperty("address"))
            object.address = message.address;
        return object;
    };

    /**
     * Converts this ThreadBackupQuery to JSON.
     * @function toJSON
     * @memberof ThreadBackupQuery
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ThreadBackupQuery.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ThreadBackupQuery;
})();

export const ThreadEnvelope = $root.ThreadEnvelope = (() => {

    /**
     * Properties of a ThreadEnvelope.
     * @exports IThreadEnvelope
     * @interface IThreadEnvelope
     * @property {string} thread ThreadEnvelope thread
     * @property {string} hash ThreadEnvelope hash
     * @property {Uint8Array} ciphertext ThreadEnvelope ciphertext
     */

    /**
     * Constructs a new ThreadEnvelope.
     * @exports ThreadEnvelope
     * @classdesc Represents a ThreadEnvelope.
     * @implements IThreadEnvelope
     * @constructor
     * @param {IThreadEnvelope=} [properties] Properties to set
     */
    function ThreadEnvelope(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ThreadEnvelope thread.
     * @member {string} thread
     * @memberof ThreadEnvelope
     * @instance
     */
    ThreadEnvelope.prototype.thread = "";

    /**
     * ThreadEnvelope hash.
     * @member {string} hash
     * @memberof ThreadEnvelope
     * @instance
     */
    ThreadEnvelope.prototype.hash = "";

    /**
     * ThreadEnvelope ciphertext.
     * @member {Uint8Array} ciphertext
     * @memberof ThreadEnvelope
     * @instance
     */
    ThreadEnvelope.prototype.ciphertext = $util.newBuffer([]);

    /**
     * Creates a new ThreadEnvelope instance using the specified properties.
     * @function create
     * @memberof ThreadEnvelope
     * @static
     * @param {IThreadEnvelope=} [properties] Properties to set
     * @returns {ThreadEnvelope} ThreadEnvelope instance
     */
    ThreadEnvelope.create = function create(properties) {
        return new ThreadEnvelope(properties);
    };

    /**
     * Encodes the specified ThreadEnvelope message. Does not implicitly {@link ThreadEnvelope.verify|verify} messages.
     * @function encode
     * @memberof ThreadEnvelope
     * @static
     * @param {IThreadEnvelope} message ThreadEnvelope message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadEnvelope.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.thread != null && message.hasOwnProperty("thread"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.thread);
        if (message.hash != null && message.hasOwnProperty("hash"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.hash);
        if (message.ciphertext != null && message.hasOwnProperty("ciphertext"))
            writer.uint32(/* id 3, wireType 2 =*/26).bytes(message.ciphertext);
        return writer;
    };

    /**
     * Encodes the specified ThreadEnvelope message, length delimited. Does not implicitly {@link ThreadEnvelope.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ThreadEnvelope
     * @static
     * @param {IThreadEnvelope} message ThreadEnvelope message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadEnvelope.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ThreadEnvelope message from the specified reader or buffer.
     * @function decode
     * @memberof ThreadEnvelope
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ThreadEnvelope} ThreadEnvelope
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadEnvelope.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ThreadEnvelope();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.thread = reader.string();
                break;
            case 2:
                message.hash = reader.string();
                break;
            case 3:
                message.ciphertext = reader.bytes();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ThreadEnvelope message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ThreadEnvelope
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ThreadEnvelope} ThreadEnvelope
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadEnvelope.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ThreadEnvelope message.
     * @function verify
     * @memberof ThreadEnvelope
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ThreadEnvelope.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.thread != null && message.hasOwnProperty("thread"))
            if (!$util.isString(message.thread))
                return "thread: string expected";
        if (message.hash != null && message.hasOwnProperty("hash"))
            if (!$util.isString(message.hash))
                return "hash: string expected";
        if (message.ciphertext != null && message.hasOwnProperty("ciphertext"))
            if (!(message.ciphertext && typeof message.ciphertext.length === "number" || $util.isString(message.ciphertext)))
                return "ciphertext: buffer expected";
        return null;
    };

    /**
     * Creates a ThreadEnvelope message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ThreadEnvelope
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ThreadEnvelope} ThreadEnvelope
     */
    ThreadEnvelope.fromObject = function fromObject(object) {
        if (object instanceof $root.ThreadEnvelope)
            return object;
        let message = new $root.ThreadEnvelope();
        if (object.thread != null)
            message.thread = String(object.thread);
        if (object.hash != null)
            message.hash = String(object.hash);
        if (object.ciphertext != null)
            if (typeof object.ciphertext === "string")
                $util.base64.decode(object.ciphertext, message.ciphertext = $util.newBuffer($util.base64.length(object.ciphertext)), 0);
            else if (object.ciphertext.length)
                message.ciphertext = object.ciphertext;
        return message;
    };

    /**
     * Creates a plain object from a ThreadEnvelope message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ThreadEnvelope
     * @static
     * @param {ThreadEnvelope} message ThreadEnvelope
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ThreadEnvelope.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.thread = "";
            object.hash = "";
            if (options.bytes === String)
                object.ciphertext = "";
            else {
                object.ciphertext = [];
                if (options.bytes !== Array)
                    object.ciphertext = $util.newBuffer(object.ciphertext);
            }
        }
        if (message.thread != null && message.hasOwnProperty("thread"))
            object.thread = message.thread;
        if (message.hash != null && message.hasOwnProperty("hash"))
            object.hash = message.hash;
        if (message.ciphertext != null && message.hasOwnProperty("ciphertext"))
            object.ciphertext = options.bytes === String ? $util.base64.encode(message.ciphertext, 0, message.ciphertext.length) : options.bytes === Array ? Array.prototype.slice.call(message.ciphertext) : message.ciphertext;
        return object;
    };

    /**
     * Converts this ThreadEnvelope to JSON.
     * @function toJSON
     * @memberof ThreadEnvelope
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ThreadEnvelope.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ThreadEnvelope;
})();

export const ThreadBlock = $root.ThreadBlock = (() => {

    /**
     * Properties of a ThreadBlock.
     * @exports IThreadBlock
     * @interface IThreadBlock
     * @property {IThreadBlockHeader} header ThreadBlock header
     * @property {Block.BlockType} type ThreadBlock type
     * @property {google.protobuf.IAny} payload ThreadBlock payload
     */

    /**
     * Constructs a new ThreadBlock.
     * @exports ThreadBlock
     * @classdesc Represents a ThreadBlock.
     * @implements IThreadBlock
     * @constructor
     * @param {IThreadBlock=} [properties] Properties to set
     */
    function ThreadBlock(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ThreadBlock header.
     * @member {IThreadBlockHeader} header
     * @memberof ThreadBlock
     * @instance
     */
    ThreadBlock.prototype.header = null;

    /**
     * ThreadBlock type.
     * @member {Block.BlockType} type
     * @memberof ThreadBlock
     * @instance
     */
    ThreadBlock.prototype.type = 0;

    /**
     * ThreadBlock payload.
     * @member {google.protobuf.IAny} payload
     * @memberof ThreadBlock
     * @instance
     */
    ThreadBlock.prototype.payload = null;

    /**
     * Creates a new ThreadBlock instance using the specified properties.
     * @function create
     * @memberof ThreadBlock
     * @static
     * @param {IThreadBlock=} [properties] Properties to set
     * @returns {ThreadBlock} ThreadBlock instance
     */
    ThreadBlock.create = function create(properties) {
        return new ThreadBlock(properties);
    };

    /**
     * Encodes the specified ThreadBlock message. Does not implicitly {@link ThreadBlock.verify|verify} messages.
     * @function encode
     * @memberof ThreadBlock
     * @static
     * @param {IThreadBlock} message ThreadBlock message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadBlock.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.header != null && message.hasOwnProperty("header"))
            $root.ThreadBlockHeader.encode(message.header, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        if (message.type != null && message.hasOwnProperty("type"))
            writer.uint32(/* id 2, wireType 0 =*/16).int32(message.type);
        if (message.payload != null && message.hasOwnProperty("payload"))
            $root.google.protobuf.Any.encode(message.payload, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified ThreadBlock message, length delimited. Does not implicitly {@link ThreadBlock.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ThreadBlock
     * @static
     * @param {IThreadBlock} message ThreadBlock message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadBlock.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ThreadBlock message from the specified reader or buffer.
     * @function decode
     * @memberof ThreadBlock
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ThreadBlock} ThreadBlock
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadBlock.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ThreadBlock();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.header = $root.ThreadBlockHeader.decode(reader, reader.uint32());
                break;
            case 2:
                message.type = reader.int32();
                break;
            case 3:
                message.payload = $root.google.protobuf.Any.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ThreadBlock message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ThreadBlock
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ThreadBlock} ThreadBlock
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadBlock.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ThreadBlock message.
     * @function verify
     * @memberof ThreadBlock
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ThreadBlock.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.header != null && message.hasOwnProperty("header")) {
            let error = $root.ThreadBlockHeader.verify(message.header);
            if (error)
                return "header." + error;
        }
        if (message.type != null && message.hasOwnProperty("type"))
            switch (message.type) {
            default:
                return "type: enum value expected";
            case 0:
            case 1:
            case 2:
            case 3:
            case 4:
            case 5:
            case 6:
            case 7:
            case 8:
            case 9:
            case 50:
                break;
            }
        if (message.payload != null && message.hasOwnProperty("payload")) {
            let error = $root.google.protobuf.Any.verify(message.payload);
            if (error)
                return "payload." + error;
        }
        return null;
    };

    /**
     * Creates a ThreadBlock message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ThreadBlock
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ThreadBlock} ThreadBlock
     */
    ThreadBlock.fromObject = function fromObject(object) {
        if (object instanceof $root.ThreadBlock)
            return object;
        let message = new $root.ThreadBlock();
        if (object.header != null) {
            if (typeof object.header !== "object")
                throw TypeError(".ThreadBlock.header: object expected");
            message.header = $root.ThreadBlockHeader.fromObject(object.header);
        }
        switch (object.type) {
        case "MERGE":
        case 0:
            message.type = 0;
            break;
        case "IGNORE":
        case 1:
            message.type = 1;
            break;
        case "FLAG":
        case 2:
            message.type = 2;
            break;
        case "JOIN":
        case 3:
            message.type = 3;
            break;
        case "ANNOUNCE":
        case 4:
            message.type = 4;
            break;
        case "LEAVE":
        case 5:
            message.type = 5;
            break;
        case "MESSAGE":
        case 6:
            message.type = 6;
            break;
        case "FILES":
        case 7:
            message.type = 7;
            break;
        case "COMMENT":
        case 8:
            message.type = 8;
            break;
        case "LIKE":
        case 9:
            message.type = 9;
            break;
        case "INVITE":
        case 50:
            message.type = 50;
            break;
        }
        if (object.payload != null) {
            if (typeof object.payload !== "object")
                throw TypeError(".ThreadBlock.payload: object expected");
            message.payload = $root.google.protobuf.Any.fromObject(object.payload);
        }
        return message;
    };

    /**
     * Creates a plain object from a ThreadBlock message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ThreadBlock
     * @static
     * @param {ThreadBlock} message ThreadBlock
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ThreadBlock.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.header = null;
            object.type = options.enums === String ? "MERGE" : 0;
            object.payload = null;
        }
        if (message.header != null && message.hasOwnProperty("header"))
            object.header = $root.ThreadBlockHeader.toObject(message.header, options);
        if (message.type != null && message.hasOwnProperty("type"))
            object.type = options.enums === String ? $root.Block.BlockType[message.type] : message.type;
        if (message.payload != null && message.hasOwnProperty("payload"))
            object.payload = $root.google.protobuf.Any.toObject(message.payload, options);
        return object;
    };

    /**
     * Converts this ThreadBlock to JSON.
     * @function toJSON
     * @memberof ThreadBlock
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ThreadBlock.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ThreadBlock;
})();

export const ThreadBlockHeader = $root.ThreadBlockHeader = (() => {

    /**
     * Properties of a ThreadBlockHeader.
     * @exports IThreadBlockHeader
     * @interface IThreadBlockHeader
     * @property {google.protobuf.ITimestamp} date ThreadBlockHeader date
     * @property {Array.<string>} parents ThreadBlockHeader parents
     * @property {string} author ThreadBlockHeader author
     * @property {string} address ThreadBlockHeader address
     */

    /**
     * Constructs a new ThreadBlockHeader.
     * @exports ThreadBlockHeader
     * @classdesc Represents a ThreadBlockHeader.
     * @implements IThreadBlockHeader
     * @constructor
     * @param {IThreadBlockHeader=} [properties] Properties to set
     */
    function ThreadBlockHeader(properties) {
        this.parents = [];
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ThreadBlockHeader date.
     * @member {google.protobuf.ITimestamp} date
     * @memberof ThreadBlockHeader
     * @instance
     */
    ThreadBlockHeader.prototype.date = null;

    /**
     * ThreadBlockHeader parents.
     * @member {Array.<string>} parents
     * @memberof ThreadBlockHeader
     * @instance
     */
    ThreadBlockHeader.prototype.parents = $util.emptyArray;

    /**
     * ThreadBlockHeader author.
     * @member {string} author
     * @memberof ThreadBlockHeader
     * @instance
     */
    ThreadBlockHeader.prototype.author = "";

    /**
     * ThreadBlockHeader address.
     * @member {string} address
     * @memberof ThreadBlockHeader
     * @instance
     */
    ThreadBlockHeader.prototype.address = "";

    /**
     * Creates a new ThreadBlockHeader instance using the specified properties.
     * @function create
     * @memberof ThreadBlockHeader
     * @static
     * @param {IThreadBlockHeader=} [properties] Properties to set
     * @returns {ThreadBlockHeader} ThreadBlockHeader instance
     */
    ThreadBlockHeader.create = function create(properties) {
        return new ThreadBlockHeader(properties);
    };

    /**
     * Encodes the specified ThreadBlockHeader message. Does not implicitly {@link ThreadBlockHeader.verify|verify} messages.
     * @function encode
     * @memberof ThreadBlockHeader
     * @static
     * @param {IThreadBlockHeader} message ThreadBlockHeader message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadBlockHeader.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.date != null && message.hasOwnProperty("date"))
            $root.google.protobuf.Timestamp.encode(message.date, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        if (message.parents != null && message.parents.length)
            for (let i = 0; i < message.parents.length; ++i)
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.parents[i]);
        if (message.author != null && message.hasOwnProperty("author"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.author);
        if (message.address != null && message.hasOwnProperty("address"))
            writer.uint32(/* id 4, wireType 2 =*/34).string(message.address);
        return writer;
    };

    /**
     * Encodes the specified ThreadBlockHeader message, length delimited. Does not implicitly {@link ThreadBlockHeader.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ThreadBlockHeader
     * @static
     * @param {IThreadBlockHeader} message ThreadBlockHeader message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadBlockHeader.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ThreadBlockHeader message from the specified reader or buffer.
     * @function decode
     * @memberof ThreadBlockHeader
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ThreadBlockHeader} ThreadBlockHeader
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadBlockHeader.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ThreadBlockHeader();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.date = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                break;
            case 2:
                if (!(message.parents && message.parents.length))
                    message.parents = [];
                message.parents.push(reader.string());
                break;
            case 3:
                message.author = reader.string();
                break;
            case 4:
                message.address = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ThreadBlockHeader message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ThreadBlockHeader
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ThreadBlockHeader} ThreadBlockHeader
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadBlockHeader.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ThreadBlockHeader message.
     * @function verify
     * @memberof ThreadBlockHeader
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ThreadBlockHeader.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.date != null && message.hasOwnProperty("date")) {
            let error = $root.google.protobuf.Timestamp.verify(message.date);
            if (error)
                return "date." + error;
        }
        if (message.parents != null && message.hasOwnProperty("parents")) {
            if (!Array.isArray(message.parents))
                return "parents: array expected";
            for (let i = 0; i < message.parents.length; ++i)
                if (!$util.isString(message.parents[i]))
                    return "parents: string[] expected";
        }
        if (message.author != null && message.hasOwnProperty("author"))
            if (!$util.isString(message.author))
                return "author: string expected";
        if (message.address != null && message.hasOwnProperty("address"))
            if (!$util.isString(message.address))
                return "address: string expected";
        return null;
    };

    /**
     * Creates a ThreadBlockHeader message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ThreadBlockHeader
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ThreadBlockHeader} ThreadBlockHeader
     */
    ThreadBlockHeader.fromObject = function fromObject(object) {
        if (object instanceof $root.ThreadBlockHeader)
            return object;
        let message = new $root.ThreadBlockHeader();
        if (object.date != null) {
            if (typeof object.date !== "object")
                throw TypeError(".ThreadBlockHeader.date: object expected");
            message.date = $root.google.protobuf.Timestamp.fromObject(object.date);
        }
        if (object.parents) {
            if (!Array.isArray(object.parents))
                throw TypeError(".ThreadBlockHeader.parents: array expected");
            message.parents = [];
            for (let i = 0; i < object.parents.length; ++i)
                message.parents[i] = String(object.parents[i]);
        }
        if (object.author != null)
            message.author = String(object.author);
        if (object.address != null)
            message.address = String(object.address);
        return message;
    };

    /**
     * Creates a plain object from a ThreadBlockHeader message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ThreadBlockHeader
     * @static
     * @param {ThreadBlockHeader} message ThreadBlockHeader
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ThreadBlockHeader.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.arrays || options.defaults)
            object.parents = [];
        if (options.defaults) {
            object.date = null;
            object.author = "";
            object.address = "";
        }
        if (message.date != null && message.hasOwnProperty("date"))
            object.date = $root.google.protobuf.Timestamp.toObject(message.date, options);
        if (message.parents && message.parents.length) {
            object.parents = [];
            for (let j = 0; j < message.parents.length; ++j)
                object.parents[j] = message.parents[j];
        }
        if (message.author != null && message.hasOwnProperty("author"))
            object.author = message.author;
        if (message.address != null && message.hasOwnProperty("address"))
            object.address = message.address;
        return object;
    };

    /**
     * Converts this ThreadBlockHeader to JSON.
     * @function toJSON
     * @memberof ThreadBlockHeader
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ThreadBlockHeader.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ThreadBlockHeader;
})();

export const ThreadInvite = $root.ThreadInvite = (() => {

    /**
     * Properties of a ThreadInvite.
     * @exports IThreadInvite
     * @interface IThreadInvite
     * @property {IContact} inviter ThreadInvite inviter
     * @property {IThread} thread ThreadInvite thread
     */

    /**
     * Constructs a new ThreadInvite.
     * @exports ThreadInvite
     * @classdesc Represents a ThreadInvite.
     * @implements IThreadInvite
     * @constructor
     * @param {IThreadInvite=} [properties] Properties to set
     */
    function ThreadInvite(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ThreadInvite inviter.
     * @member {IContact} inviter
     * @memberof ThreadInvite
     * @instance
     */
    ThreadInvite.prototype.inviter = null;

    /**
     * ThreadInvite thread.
     * @member {IThread} thread
     * @memberof ThreadInvite
     * @instance
     */
    ThreadInvite.prototype.thread = null;

    /**
     * Creates a new ThreadInvite instance using the specified properties.
     * @function create
     * @memberof ThreadInvite
     * @static
     * @param {IThreadInvite=} [properties] Properties to set
     * @returns {ThreadInvite} ThreadInvite instance
     */
    ThreadInvite.create = function create(properties) {
        return new ThreadInvite(properties);
    };

    /**
     * Encodes the specified ThreadInvite message. Does not implicitly {@link ThreadInvite.verify|verify} messages.
     * @function encode
     * @memberof ThreadInvite
     * @static
     * @param {IThreadInvite} message ThreadInvite message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadInvite.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.inviter != null && message.hasOwnProperty("inviter"))
            $root.Contact.encode(message.inviter, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        if (message.thread != null && message.hasOwnProperty("thread"))
            $root.Thread.encode(message.thread, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified ThreadInvite message, length delimited. Does not implicitly {@link ThreadInvite.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ThreadInvite
     * @static
     * @param {IThreadInvite} message ThreadInvite message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadInvite.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ThreadInvite message from the specified reader or buffer.
     * @function decode
     * @memberof ThreadInvite
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ThreadInvite} ThreadInvite
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadInvite.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ThreadInvite();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.inviter = $root.Contact.decode(reader, reader.uint32());
                break;
            case 2:
                message.thread = $root.Thread.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ThreadInvite message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ThreadInvite
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ThreadInvite} ThreadInvite
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadInvite.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ThreadInvite message.
     * @function verify
     * @memberof ThreadInvite
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ThreadInvite.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.inviter != null && message.hasOwnProperty("inviter")) {
            let error = $root.Contact.verify(message.inviter);
            if (error)
                return "inviter." + error;
        }
        if (message.thread != null && message.hasOwnProperty("thread")) {
            let error = $root.Thread.verify(message.thread);
            if (error)
                return "thread." + error;
        }
        return null;
    };

    /**
     * Creates a ThreadInvite message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ThreadInvite
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ThreadInvite} ThreadInvite
     */
    ThreadInvite.fromObject = function fromObject(object) {
        if (object instanceof $root.ThreadInvite)
            return object;
        let message = new $root.ThreadInvite();
        if (object.inviter != null) {
            if (typeof object.inviter !== "object")
                throw TypeError(".ThreadInvite.inviter: object expected");
            message.inviter = $root.Contact.fromObject(object.inviter);
        }
        if (object.thread != null) {
            if (typeof object.thread !== "object")
                throw TypeError(".ThreadInvite.thread: object expected");
            message.thread = $root.Thread.fromObject(object.thread);
        }
        return message;
    };

    /**
     * Creates a plain object from a ThreadInvite message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ThreadInvite
     * @static
     * @param {ThreadInvite} message ThreadInvite
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ThreadInvite.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.inviter = null;
            object.thread = null;
        }
        if (message.inviter != null && message.hasOwnProperty("inviter"))
            object.inviter = $root.Contact.toObject(message.inviter, options);
        if (message.thread != null && message.hasOwnProperty("thread"))
            object.thread = $root.Thread.toObject(message.thread, options);
        return object;
    };

    /**
     * Converts this ThreadInvite to JSON.
     * @function toJSON
     * @memberof ThreadInvite
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ThreadInvite.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ThreadInvite;
})();

export const ThreadIgnore = $root.ThreadIgnore = (() => {

    /**
     * Properties of a ThreadIgnore.
     * @exports IThreadIgnore
     * @interface IThreadIgnore
     * @property {string} target ThreadIgnore target
     */

    /**
     * Constructs a new ThreadIgnore.
     * @exports ThreadIgnore
     * @classdesc Represents a ThreadIgnore.
     * @implements IThreadIgnore
     * @constructor
     * @param {IThreadIgnore=} [properties] Properties to set
     */
    function ThreadIgnore(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ThreadIgnore target.
     * @member {string} target
     * @memberof ThreadIgnore
     * @instance
     */
    ThreadIgnore.prototype.target = "";

    /**
     * Creates a new ThreadIgnore instance using the specified properties.
     * @function create
     * @memberof ThreadIgnore
     * @static
     * @param {IThreadIgnore=} [properties] Properties to set
     * @returns {ThreadIgnore} ThreadIgnore instance
     */
    ThreadIgnore.create = function create(properties) {
        return new ThreadIgnore(properties);
    };

    /**
     * Encodes the specified ThreadIgnore message. Does not implicitly {@link ThreadIgnore.verify|verify} messages.
     * @function encode
     * @memberof ThreadIgnore
     * @static
     * @param {IThreadIgnore} message ThreadIgnore message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadIgnore.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.target != null && message.hasOwnProperty("target"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.target);
        return writer;
    };

    /**
     * Encodes the specified ThreadIgnore message, length delimited. Does not implicitly {@link ThreadIgnore.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ThreadIgnore
     * @static
     * @param {IThreadIgnore} message ThreadIgnore message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadIgnore.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ThreadIgnore message from the specified reader or buffer.
     * @function decode
     * @memberof ThreadIgnore
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ThreadIgnore} ThreadIgnore
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadIgnore.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ThreadIgnore();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.target = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ThreadIgnore message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ThreadIgnore
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ThreadIgnore} ThreadIgnore
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadIgnore.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ThreadIgnore message.
     * @function verify
     * @memberof ThreadIgnore
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ThreadIgnore.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.target != null && message.hasOwnProperty("target"))
            if (!$util.isString(message.target))
                return "target: string expected";
        return null;
    };

    /**
     * Creates a ThreadIgnore message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ThreadIgnore
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ThreadIgnore} ThreadIgnore
     */
    ThreadIgnore.fromObject = function fromObject(object) {
        if (object instanceof $root.ThreadIgnore)
            return object;
        let message = new $root.ThreadIgnore();
        if (object.target != null)
            message.target = String(object.target);
        return message;
    };

    /**
     * Creates a plain object from a ThreadIgnore message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ThreadIgnore
     * @static
     * @param {ThreadIgnore} message ThreadIgnore
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ThreadIgnore.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults)
            object.target = "";
        if (message.target != null && message.hasOwnProperty("target"))
            object.target = message.target;
        return object;
    };

    /**
     * Converts this ThreadIgnore to JSON.
     * @function toJSON
     * @memberof ThreadIgnore
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ThreadIgnore.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ThreadIgnore;
})();

export const ThreadFlag = $root.ThreadFlag = (() => {

    /**
     * Properties of a ThreadFlag.
     * @exports IThreadFlag
     * @interface IThreadFlag
     * @property {string} target ThreadFlag target
     */

    /**
     * Constructs a new ThreadFlag.
     * @exports ThreadFlag
     * @classdesc Represents a ThreadFlag.
     * @implements IThreadFlag
     * @constructor
     * @param {IThreadFlag=} [properties] Properties to set
     */
    function ThreadFlag(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ThreadFlag target.
     * @member {string} target
     * @memberof ThreadFlag
     * @instance
     */
    ThreadFlag.prototype.target = "";

    /**
     * Creates a new ThreadFlag instance using the specified properties.
     * @function create
     * @memberof ThreadFlag
     * @static
     * @param {IThreadFlag=} [properties] Properties to set
     * @returns {ThreadFlag} ThreadFlag instance
     */
    ThreadFlag.create = function create(properties) {
        return new ThreadFlag(properties);
    };

    /**
     * Encodes the specified ThreadFlag message. Does not implicitly {@link ThreadFlag.verify|verify} messages.
     * @function encode
     * @memberof ThreadFlag
     * @static
     * @param {IThreadFlag} message ThreadFlag message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadFlag.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.target != null && message.hasOwnProperty("target"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.target);
        return writer;
    };

    /**
     * Encodes the specified ThreadFlag message, length delimited. Does not implicitly {@link ThreadFlag.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ThreadFlag
     * @static
     * @param {IThreadFlag} message ThreadFlag message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadFlag.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ThreadFlag message from the specified reader or buffer.
     * @function decode
     * @memberof ThreadFlag
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ThreadFlag} ThreadFlag
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadFlag.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ThreadFlag();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.target = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ThreadFlag message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ThreadFlag
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ThreadFlag} ThreadFlag
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadFlag.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ThreadFlag message.
     * @function verify
     * @memberof ThreadFlag
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ThreadFlag.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.target != null && message.hasOwnProperty("target"))
            if (!$util.isString(message.target))
                return "target: string expected";
        return null;
    };

    /**
     * Creates a ThreadFlag message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ThreadFlag
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ThreadFlag} ThreadFlag
     */
    ThreadFlag.fromObject = function fromObject(object) {
        if (object instanceof $root.ThreadFlag)
            return object;
        let message = new $root.ThreadFlag();
        if (object.target != null)
            message.target = String(object.target);
        return message;
    };

    /**
     * Creates a plain object from a ThreadFlag message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ThreadFlag
     * @static
     * @param {ThreadFlag} message ThreadFlag
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ThreadFlag.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults)
            object.target = "";
        if (message.target != null && message.hasOwnProperty("target"))
            object.target = message.target;
        return object;
    };

    /**
     * Converts this ThreadFlag to JSON.
     * @function toJSON
     * @memberof ThreadFlag
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ThreadFlag.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ThreadFlag;
})();

export const ThreadJoin = $root.ThreadJoin = (() => {

    /**
     * Properties of a ThreadJoin.
     * @exports IThreadJoin
     * @interface IThreadJoin
     * @property {string} inviter ThreadJoin inviter
     * @property {IContact} contact ThreadJoin contact
     */

    /**
     * Constructs a new ThreadJoin.
     * @exports ThreadJoin
     * @classdesc Represents a ThreadJoin.
     * @implements IThreadJoin
     * @constructor
     * @param {IThreadJoin=} [properties] Properties to set
     */
    function ThreadJoin(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ThreadJoin inviter.
     * @member {string} inviter
     * @memberof ThreadJoin
     * @instance
     */
    ThreadJoin.prototype.inviter = "";

    /**
     * ThreadJoin contact.
     * @member {IContact} contact
     * @memberof ThreadJoin
     * @instance
     */
    ThreadJoin.prototype.contact = null;

    /**
     * Creates a new ThreadJoin instance using the specified properties.
     * @function create
     * @memberof ThreadJoin
     * @static
     * @param {IThreadJoin=} [properties] Properties to set
     * @returns {ThreadJoin} ThreadJoin instance
     */
    ThreadJoin.create = function create(properties) {
        return new ThreadJoin(properties);
    };

    /**
     * Encodes the specified ThreadJoin message. Does not implicitly {@link ThreadJoin.verify|verify} messages.
     * @function encode
     * @memberof ThreadJoin
     * @static
     * @param {IThreadJoin} message ThreadJoin message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadJoin.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.inviter != null && message.hasOwnProperty("inviter"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.inviter);
        if (message.contact != null && message.hasOwnProperty("contact"))
            $root.Contact.encode(message.contact, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified ThreadJoin message, length delimited. Does not implicitly {@link ThreadJoin.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ThreadJoin
     * @static
     * @param {IThreadJoin} message ThreadJoin message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadJoin.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ThreadJoin message from the specified reader or buffer.
     * @function decode
     * @memberof ThreadJoin
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ThreadJoin} ThreadJoin
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadJoin.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ThreadJoin();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.inviter = reader.string();
                break;
            case 2:
                message.contact = $root.Contact.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ThreadJoin message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ThreadJoin
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ThreadJoin} ThreadJoin
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadJoin.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ThreadJoin message.
     * @function verify
     * @memberof ThreadJoin
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ThreadJoin.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.inviter != null && message.hasOwnProperty("inviter"))
            if (!$util.isString(message.inviter))
                return "inviter: string expected";
        if (message.contact != null && message.hasOwnProperty("contact")) {
            let error = $root.Contact.verify(message.contact);
            if (error)
                return "contact." + error;
        }
        return null;
    };

    /**
     * Creates a ThreadJoin message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ThreadJoin
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ThreadJoin} ThreadJoin
     */
    ThreadJoin.fromObject = function fromObject(object) {
        if (object instanceof $root.ThreadJoin)
            return object;
        let message = new $root.ThreadJoin();
        if (object.inviter != null)
            message.inviter = String(object.inviter);
        if (object.contact != null) {
            if (typeof object.contact !== "object")
                throw TypeError(".ThreadJoin.contact: object expected");
            message.contact = $root.Contact.fromObject(object.contact);
        }
        return message;
    };

    /**
     * Creates a plain object from a ThreadJoin message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ThreadJoin
     * @static
     * @param {ThreadJoin} message ThreadJoin
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ThreadJoin.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.inviter = "";
            object.contact = null;
        }
        if (message.inviter != null && message.hasOwnProperty("inviter"))
            object.inviter = message.inviter;
        if (message.contact != null && message.hasOwnProperty("contact"))
            object.contact = $root.Contact.toObject(message.contact, options);
        return object;
    };

    /**
     * Converts this ThreadJoin to JSON.
     * @function toJSON
     * @memberof ThreadJoin
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ThreadJoin.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ThreadJoin;
})();

export const ThreadAnnounce = $root.ThreadAnnounce = (() => {

    /**
     * Properties of a ThreadAnnounce.
     * @exports IThreadAnnounce
     * @interface IThreadAnnounce
     * @property {IContact} contact ThreadAnnounce contact
     */

    /**
     * Constructs a new ThreadAnnounce.
     * @exports ThreadAnnounce
     * @classdesc Represents a ThreadAnnounce.
     * @implements IThreadAnnounce
     * @constructor
     * @param {IThreadAnnounce=} [properties] Properties to set
     */
    function ThreadAnnounce(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ThreadAnnounce contact.
     * @member {IContact} contact
     * @memberof ThreadAnnounce
     * @instance
     */
    ThreadAnnounce.prototype.contact = null;

    /**
     * Creates a new ThreadAnnounce instance using the specified properties.
     * @function create
     * @memberof ThreadAnnounce
     * @static
     * @param {IThreadAnnounce=} [properties] Properties to set
     * @returns {ThreadAnnounce} ThreadAnnounce instance
     */
    ThreadAnnounce.create = function create(properties) {
        return new ThreadAnnounce(properties);
    };

    /**
     * Encodes the specified ThreadAnnounce message. Does not implicitly {@link ThreadAnnounce.verify|verify} messages.
     * @function encode
     * @memberof ThreadAnnounce
     * @static
     * @param {IThreadAnnounce} message ThreadAnnounce message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadAnnounce.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.contact != null && message.hasOwnProperty("contact"))
            $root.Contact.encode(message.contact, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified ThreadAnnounce message, length delimited. Does not implicitly {@link ThreadAnnounce.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ThreadAnnounce
     * @static
     * @param {IThreadAnnounce} message ThreadAnnounce message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadAnnounce.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ThreadAnnounce message from the specified reader or buffer.
     * @function decode
     * @memberof ThreadAnnounce
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ThreadAnnounce} ThreadAnnounce
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadAnnounce.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ThreadAnnounce();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.contact = $root.Contact.decode(reader, reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ThreadAnnounce message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ThreadAnnounce
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ThreadAnnounce} ThreadAnnounce
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadAnnounce.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ThreadAnnounce message.
     * @function verify
     * @memberof ThreadAnnounce
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ThreadAnnounce.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.contact != null && message.hasOwnProperty("contact")) {
            let error = $root.Contact.verify(message.contact);
            if (error)
                return "contact." + error;
        }
        return null;
    };

    /**
     * Creates a ThreadAnnounce message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ThreadAnnounce
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ThreadAnnounce} ThreadAnnounce
     */
    ThreadAnnounce.fromObject = function fromObject(object) {
        if (object instanceof $root.ThreadAnnounce)
            return object;
        let message = new $root.ThreadAnnounce();
        if (object.contact != null) {
            if (typeof object.contact !== "object")
                throw TypeError(".ThreadAnnounce.contact: object expected");
            message.contact = $root.Contact.fromObject(object.contact);
        }
        return message;
    };

    /**
     * Creates a plain object from a ThreadAnnounce message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ThreadAnnounce
     * @static
     * @param {ThreadAnnounce} message ThreadAnnounce
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ThreadAnnounce.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults)
            object.contact = null;
        if (message.contact != null && message.hasOwnProperty("contact"))
            object.contact = $root.Contact.toObject(message.contact, options);
        return object;
    };

    /**
     * Converts this ThreadAnnounce to JSON.
     * @function toJSON
     * @memberof ThreadAnnounce
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ThreadAnnounce.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ThreadAnnounce;
})();

export const ThreadMessage = $root.ThreadMessage = (() => {

    /**
     * Properties of a ThreadMessage.
     * @exports IThreadMessage
     * @interface IThreadMessage
     * @property {string} body ThreadMessage body
     */

    /**
     * Constructs a new ThreadMessage.
     * @exports ThreadMessage
     * @classdesc Represents a ThreadMessage.
     * @implements IThreadMessage
     * @constructor
     * @param {IThreadMessage=} [properties] Properties to set
     */
    function ThreadMessage(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ThreadMessage body.
     * @member {string} body
     * @memberof ThreadMessage
     * @instance
     */
    ThreadMessage.prototype.body = "";

    /**
     * Creates a new ThreadMessage instance using the specified properties.
     * @function create
     * @memberof ThreadMessage
     * @static
     * @param {IThreadMessage=} [properties] Properties to set
     * @returns {ThreadMessage} ThreadMessage instance
     */
    ThreadMessage.create = function create(properties) {
        return new ThreadMessage(properties);
    };

    /**
     * Encodes the specified ThreadMessage message. Does not implicitly {@link ThreadMessage.verify|verify} messages.
     * @function encode
     * @memberof ThreadMessage
     * @static
     * @param {IThreadMessage} message ThreadMessage message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadMessage.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.body != null && message.hasOwnProperty("body"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.body);
        return writer;
    };

    /**
     * Encodes the specified ThreadMessage message, length delimited. Does not implicitly {@link ThreadMessage.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ThreadMessage
     * @static
     * @param {IThreadMessage} message ThreadMessage message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadMessage.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ThreadMessage message from the specified reader or buffer.
     * @function decode
     * @memberof ThreadMessage
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ThreadMessage} ThreadMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadMessage.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ThreadMessage();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.body = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ThreadMessage message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ThreadMessage
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ThreadMessage} ThreadMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadMessage.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ThreadMessage message.
     * @function verify
     * @memberof ThreadMessage
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ThreadMessage.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.body != null && message.hasOwnProperty("body"))
            if (!$util.isString(message.body))
                return "body: string expected";
        return null;
    };

    /**
     * Creates a ThreadMessage message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ThreadMessage
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ThreadMessage} ThreadMessage
     */
    ThreadMessage.fromObject = function fromObject(object) {
        if (object instanceof $root.ThreadMessage)
            return object;
        let message = new $root.ThreadMessage();
        if (object.body != null)
            message.body = String(object.body);
        return message;
    };

    /**
     * Creates a plain object from a ThreadMessage message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ThreadMessage
     * @static
     * @param {ThreadMessage} message ThreadMessage
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ThreadMessage.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults)
            object.body = "";
        if (message.body != null && message.hasOwnProperty("body"))
            object.body = message.body;
        return object;
    };

    /**
     * Converts this ThreadMessage to JSON.
     * @function toJSON
     * @memberof ThreadMessage
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ThreadMessage.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ThreadMessage;
})();

export const ThreadFiles = $root.ThreadFiles = (() => {

    /**
     * Properties of a ThreadFiles.
     * @exports IThreadFiles
     * @interface IThreadFiles
     * @property {string} target ThreadFiles target
     * @property {string} body ThreadFiles body
     * @property {Object.<string,string>} keys ThreadFiles keys
     */

    /**
     * Constructs a new ThreadFiles.
     * @exports ThreadFiles
     * @classdesc Represents a ThreadFiles.
     * @implements IThreadFiles
     * @constructor
     * @param {IThreadFiles=} [properties] Properties to set
     */
    function ThreadFiles(properties) {
        this.keys = {};
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ThreadFiles target.
     * @member {string} target
     * @memberof ThreadFiles
     * @instance
     */
    ThreadFiles.prototype.target = "";

    /**
     * ThreadFiles body.
     * @member {string} body
     * @memberof ThreadFiles
     * @instance
     */
    ThreadFiles.prototype.body = "";

    /**
     * ThreadFiles keys.
     * @member {Object.<string,string>} keys
     * @memberof ThreadFiles
     * @instance
     */
    ThreadFiles.prototype.keys = $util.emptyObject;

    /**
     * Creates a new ThreadFiles instance using the specified properties.
     * @function create
     * @memberof ThreadFiles
     * @static
     * @param {IThreadFiles=} [properties] Properties to set
     * @returns {ThreadFiles} ThreadFiles instance
     */
    ThreadFiles.create = function create(properties) {
        return new ThreadFiles(properties);
    };

    /**
     * Encodes the specified ThreadFiles message. Does not implicitly {@link ThreadFiles.verify|verify} messages.
     * @function encode
     * @memberof ThreadFiles
     * @static
     * @param {IThreadFiles} message ThreadFiles message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadFiles.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.target != null && message.hasOwnProperty("target"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.target);
        if (message.body != null && message.hasOwnProperty("body"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.body);
        if (message.keys != null && message.hasOwnProperty("keys"))
            for (let keys = Object.keys(message.keys), i = 0; i < keys.length; ++i)
                writer.uint32(/* id 3, wireType 2 =*/26).fork().uint32(/* id 1, wireType 2 =*/10).string(keys[i]).uint32(/* id 2, wireType 2 =*/18).string(message.keys[keys[i]]).ldelim();
        return writer;
    };

    /**
     * Encodes the specified ThreadFiles message, length delimited. Does not implicitly {@link ThreadFiles.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ThreadFiles
     * @static
     * @param {IThreadFiles} message ThreadFiles message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadFiles.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ThreadFiles message from the specified reader or buffer.
     * @function decode
     * @memberof ThreadFiles
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ThreadFiles} ThreadFiles
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadFiles.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ThreadFiles(), key;
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.target = reader.string();
                break;
            case 2:
                message.body = reader.string();
                break;
            case 3:
                reader.skip().pos++;
                if (message.keys === $util.emptyObject)
                    message.keys = {};
                key = reader.string();
                reader.pos++;
                message.keys[key] = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ThreadFiles message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ThreadFiles
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ThreadFiles} ThreadFiles
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadFiles.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ThreadFiles message.
     * @function verify
     * @memberof ThreadFiles
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ThreadFiles.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.target != null && message.hasOwnProperty("target"))
            if (!$util.isString(message.target))
                return "target: string expected";
        if (message.body != null && message.hasOwnProperty("body"))
            if (!$util.isString(message.body))
                return "body: string expected";
        if (message.keys != null && message.hasOwnProperty("keys")) {
            if (!$util.isObject(message.keys))
                return "keys: object expected";
            let key = Object.keys(message.keys);
            for (let i = 0; i < key.length; ++i)
                if (!$util.isString(message.keys[key[i]]))
                    return "keys: string{k:string} expected";
        }
        return null;
    };

    /**
     * Creates a ThreadFiles message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ThreadFiles
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ThreadFiles} ThreadFiles
     */
    ThreadFiles.fromObject = function fromObject(object) {
        if (object instanceof $root.ThreadFiles)
            return object;
        let message = new $root.ThreadFiles();
        if (object.target != null)
            message.target = String(object.target);
        if (object.body != null)
            message.body = String(object.body);
        if (object.keys) {
            if (typeof object.keys !== "object")
                throw TypeError(".ThreadFiles.keys: object expected");
            message.keys = {};
            for (let keys = Object.keys(object.keys), i = 0; i < keys.length; ++i)
                message.keys[keys[i]] = String(object.keys[keys[i]]);
        }
        return message;
    };

    /**
     * Creates a plain object from a ThreadFiles message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ThreadFiles
     * @static
     * @param {ThreadFiles} message ThreadFiles
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ThreadFiles.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.objects || options.defaults)
            object.keys = {};
        if (options.defaults) {
            object.target = "";
            object.body = "";
        }
        if (message.target != null && message.hasOwnProperty("target"))
            object.target = message.target;
        if (message.body != null && message.hasOwnProperty("body"))
            object.body = message.body;
        let keys2;
        if (message.keys && (keys2 = Object.keys(message.keys)).length) {
            object.keys = {};
            for (let j = 0; j < keys2.length; ++j)
                object.keys[keys2[j]] = message.keys[keys2[j]];
        }
        return object;
    };

    /**
     * Converts this ThreadFiles to JSON.
     * @function toJSON
     * @memberof ThreadFiles
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ThreadFiles.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ThreadFiles;
})();

export const ThreadComment = $root.ThreadComment = (() => {

    /**
     * Properties of a ThreadComment.
     * @exports IThreadComment
     * @interface IThreadComment
     * @property {string} target ThreadComment target
     * @property {string} body ThreadComment body
     */

    /**
     * Constructs a new ThreadComment.
     * @exports ThreadComment
     * @classdesc Represents a ThreadComment.
     * @implements IThreadComment
     * @constructor
     * @param {IThreadComment=} [properties] Properties to set
     */
    function ThreadComment(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ThreadComment target.
     * @member {string} target
     * @memberof ThreadComment
     * @instance
     */
    ThreadComment.prototype.target = "";

    /**
     * ThreadComment body.
     * @member {string} body
     * @memberof ThreadComment
     * @instance
     */
    ThreadComment.prototype.body = "";

    /**
     * Creates a new ThreadComment instance using the specified properties.
     * @function create
     * @memberof ThreadComment
     * @static
     * @param {IThreadComment=} [properties] Properties to set
     * @returns {ThreadComment} ThreadComment instance
     */
    ThreadComment.create = function create(properties) {
        return new ThreadComment(properties);
    };

    /**
     * Encodes the specified ThreadComment message. Does not implicitly {@link ThreadComment.verify|verify} messages.
     * @function encode
     * @memberof ThreadComment
     * @static
     * @param {IThreadComment} message ThreadComment message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadComment.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.target != null && message.hasOwnProperty("target"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.target);
        if (message.body != null && message.hasOwnProperty("body"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.body);
        return writer;
    };

    /**
     * Encodes the specified ThreadComment message, length delimited. Does not implicitly {@link ThreadComment.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ThreadComment
     * @static
     * @param {IThreadComment} message ThreadComment message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadComment.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ThreadComment message from the specified reader or buffer.
     * @function decode
     * @memberof ThreadComment
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ThreadComment} ThreadComment
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadComment.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ThreadComment();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.target = reader.string();
                break;
            case 2:
                message.body = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ThreadComment message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ThreadComment
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ThreadComment} ThreadComment
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadComment.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ThreadComment message.
     * @function verify
     * @memberof ThreadComment
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ThreadComment.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.target != null && message.hasOwnProperty("target"))
            if (!$util.isString(message.target))
                return "target: string expected";
        if (message.body != null && message.hasOwnProperty("body"))
            if (!$util.isString(message.body))
                return "body: string expected";
        return null;
    };

    /**
     * Creates a ThreadComment message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ThreadComment
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ThreadComment} ThreadComment
     */
    ThreadComment.fromObject = function fromObject(object) {
        if (object instanceof $root.ThreadComment)
            return object;
        let message = new $root.ThreadComment();
        if (object.target != null)
            message.target = String(object.target);
        if (object.body != null)
            message.body = String(object.body);
        return message;
    };

    /**
     * Creates a plain object from a ThreadComment message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ThreadComment
     * @static
     * @param {ThreadComment} message ThreadComment
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ThreadComment.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.target = "";
            object.body = "";
        }
        if (message.target != null && message.hasOwnProperty("target"))
            object.target = message.target;
        if (message.body != null && message.hasOwnProperty("body"))
            object.body = message.body;
        return object;
    };

    /**
     * Converts this ThreadComment to JSON.
     * @function toJSON
     * @memberof ThreadComment
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ThreadComment.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ThreadComment;
})();

export const ThreadLike = $root.ThreadLike = (() => {

    /**
     * Properties of a ThreadLike.
     * @exports IThreadLike
     * @interface IThreadLike
     * @property {string} target ThreadLike target
     */

    /**
     * Constructs a new ThreadLike.
     * @exports ThreadLike
     * @classdesc Represents a ThreadLike.
     * @implements IThreadLike
     * @constructor
     * @param {IThreadLike=} [properties] Properties to set
     */
    function ThreadLike(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ThreadLike target.
     * @member {string} target
     * @memberof ThreadLike
     * @instance
     */
    ThreadLike.prototype.target = "";

    /**
     * Creates a new ThreadLike instance using the specified properties.
     * @function create
     * @memberof ThreadLike
     * @static
     * @param {IThreadLike=} [properties] Properties to set
     * @returns {ThreadLike} ThreadLike instance
     */
    ThreadLike.create = function create(properties) {
        return new ThreadLike(properties);
    };

    /**
     * Encodes the specified ThreadLike message. Does not implicitly {@link ThreadLike.verify|verify} messages.
     * @function encode
     * @memberof ThreadLike
     * @static
     * @param {IThreadLike} message ThreadLike message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadLike.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.target != null && message.hasOwnProperty("target"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.target);
        return writer;
    };

    /**
     * Encodes the specified ThreadLike message, length delimited. Does not implicitly {@link ThreadLike.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ThreadLike
     * @static
     * @param {IThreadLike} message ThreadLike message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ThreadLike.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ThreadLike message from the specified reader or buffer.
     * @function decode
     * @memberof ThreadLike
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ThreadLike} ThreadLike
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadLike.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ThreadLike();
        while (reader.pos < end) {
            let tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.target = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ThreadLike message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ThreadLike
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ThreadLike} ThreadLike
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ThreadLike.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ThreadLike message.
     * @function verify
     * @memberof ThreadLike
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ThreadLike.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.target != null && message.hasOwnProperty("target"))
            if (!$util.isString(message.target))
                return "target: string expected";
        return null;
    };

    /**
     * Creates a ThreadLike message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ThreadLike
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ThreadLike} ThreadLike
     */
    ThreadLike.fromObject = function fromObject(object) {
        if (object instanceof $root.ThreadLike)
            return object;
        let message = new $root.ThreadLike();
        if (object.target != null)
            message.target = String(object.target);
        return message;
    };

    /**
     * Creates a plain object from a ThreadLike message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ThreadLike
     * @static
     * @param {ThreadLike} message ThreadLike
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ThreadLike.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults)
            object.target = "";
        if (message.target != null && message.hasOwnProperty("target"))
            object.target = message.target;
        return object;
    };

    /**
     * Converts this ThreadLike to JSON.
     * @function toJSON
     * @memberof ThreadLike
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ThreadLike.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ThreadLike;
})();

export { $root as default };
