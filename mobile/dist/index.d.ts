import * as $protobuf from "protobufjs";
/** Properties of a CafeChallenge. */
export interface ICafeChallenge {

    /** CafeChallenge address */
    address: string;
}

/** Represents a CafeChallenge. */
export class CafeChallenge implements ICafeChallenge {

    /**
     * Constructs a new CafeChallenge.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeChallenge);

    /** CafeChallenge address. */
    public address: string;

    /**
     * Creates a new CafeChallenge instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeChallenge instance
     */
    public static create(properties?: ICafeChallenge): CafeChallenge;

    /**
     * Encodes the specified CafeChallenge message. Does not implicitly {@link CafeChallenge.verify|verify} messages.
     * @param message CafeChallenge message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeChallenge, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeChallenge message, length delimited. Does not implicitly {@link CafeChallenge.verify|verify} messages.
     * @param message CafeChallenge message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeChallenge, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeChallenge message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeChallenge
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeChallenge;

    /**
     * Decodes a CafeChallenge message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeChallenge
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeChallenge;

    /**
     * Verifies a CafeChallenge message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeChallenge message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeChallenge
     */
    public static fromObject(object: { [k: string]: any }): CafeChallenge;

    /**
     * Creates a plain object from a CafeChallenge message. Also converts values to other types if specified.
     * @param message CafeChallenge
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeChallenge, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeChallenge to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeNonce. */
export interface ICafeNonce {

    /** CafeNonce value */
    value: string;
}

/** Represents a CafeNonce. */
export class CafeNonce implements ICafeNonce {

    /**
     * Constructs a new CafeNonce.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeNonce);

    /** CafeNonce value. */
    public value: string;

    /**
     * Creates a new CafeNonce instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeNonce instance
     */
    public static create(properties?: ICafeNonce): CafeNonce;

    /**
     * Encodes the specified CafeNonce message. Does not implicitly {@link CafeNonce.verify|verify} messages.
     * @param message CafeNonce message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeNonce, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeNonce message, length delimited. Does not implicitly {@link CafeNonce.verify|verify} messages.
     * @param message CafeNonce message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeNonce, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeNonce message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeNonce
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeNonce;

    /**
     * Decodes a CafeNonce message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeNonce
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeNonce;

    /**
     * Verifies a CafeNonce message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeNonce message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeNonce
     */
    public static fromObject(object: { [k: string]: any }): CafeNonce;

    /**
     * Creates a plain object from a CafeNonce message. Also converts values to other types if specified.
     * @param message CafeNonce
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeNonce, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeNonce to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeRegistration. */
export interface ICafeRegistration {

    /** CafeRegistration address */
    address: string;

    /** CafeRegistration value */
    value: string;

    /** CafeRegistration nonce */
    nonce: string;

    /** CafeRegistration sig */
    sig: Uint8Array;

    /** CafeRegistration token */
    token: string;
}

/** Represents a CafeRegistration. */
export class CafeRegistration implements ICafeRegistration {

    /**
     * Constructs a new CafeRegistration.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeRegistration);

    /** CafeRegistration address. */
    public address: string;

    /** CafeRegistration value. */
    public value: string;

    /** CafeRegistration nonce. */
    public nonce: string;

    /** CafeRegistration sig. */
    public sig: Uint8Array;

    /** CafeRegistration token. */
    public token: string;

    /**
     * Creates a new CafeRegistration instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeRegistration instance
     */
    public static create(properties?: ICafeRegistration): CafeRegistration;

    /**
     * Encodes the specified CafeRegistration message. Does not implicitly {@link CafeRegistration.verify|verify} messages.
     * @param message CafeRegistration message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeRegistration, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeRegistration message, length delimited. Does not implicitly {@link CafeRegistration.verify|verify} messages.
     * @param message CafeRegistration message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeRegistration, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeRegistration message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeRegistration
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeRegistration;

    /**
     * Decodes a CafeRegistration message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeRegistration
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeRegistration;

    /**
     * Verifies a CafeRegistration message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeRegistration message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeRegistration
     */
    public static fromObject(object: { [k: string]: any }): CafeRegistration;

    /**
     * Creates a plain object from a CafeRegistration message. Also converts values to other types if specified.
     * @param message CafeRegistration
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeRegistration, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeRegistration to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeRefreshSession. */
export interface ICafeRefreshSession {

    /** CafeRefreshSession access */
    access: string;

    /** CafeRefreshSession refresh */
    refresh: string;
}

/** Represents a CafeRefreshSession. */
export class CafeRefreshSession implements ICafeRefreshSession {

    /**
     * Constructs a new CafeRefreshSession.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeRefreshSession);

    /** CafeRefreshSession access. */
    public access: string;

    /** CafeRefreshSession refresh. */
    public refresh: string;

    /**
     * Creates a new CafeRefreshSession instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeRefreshSession instance
     */
    public static create(properties?: ICafeRefreshSession): CafeRefreshSession;

    /**
     * Encodes the specified CafeRefreshSession message. Does not implicitly {@link CafeRefreshSession.verify|verify} messages.
     * @param message CafeRefreshSession message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeRefreshSession, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeRefreshSession message, length delimited. Does not implicitly {@link CafeRefreshSession.verify|verify} messages.
     * @param message CafeRefreshSession message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeRefreshSession, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeRefreshSession message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeRefreshSession
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeRefreshSession;

    /**
     * Decodes a CafeRefreshSession message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeRefreshSession
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeRefreshSession;

    /**
     * Verifies a CafeRefreshSession message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeRefreshSession message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeRefreshSession
     */
    public static fromObject(object: { [k: string]: any }): CafeRefreshSession;

    /**
     * Creates a plain object from a CafeRefreshSession message. Also converts values to other types if specified.
     * @param message CafeRefreshSession
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeRefreshSession, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeRefreshSession to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafePublishContact. */
export interface ICafePublishContact {

    /** CafePublishContact token */
    token: string;

    /** CafePublishContact contact */
    contact: IContact;
}

/** Represents a CafePublishContact. */
export class CafePublishContact implements ICafePublishContact {

    /**
     * Constructs a new CafePublishContact.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafePublishContact);

    /** CafePublishContact token. */
    public token: string;

    /** CafePublishContact contact. */
    public contact: IContact;

    /**
     * Creates a new CafePublishContact instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafePublishContact instance
     */
    public static create(properties?: ICafePublishContact): CafePublishContact;

    /**
     * Encodes the specified CafePublishContact message. Does not implicitly {@link CafePublishContact.verify|verify} messages.
     * @param message CafePublishContact message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafePublishContact, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafePublishContact message, length delimited. Does not implicitly {@link CafePublishContact.verify|verify} messages.
     * @param message CafePublishContact message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafePublishContact, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafePublishContact message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafePublishContact
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafePublishContact;

    /**
     * Decodes a CafePublishContact message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafePublishContact
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafePublishContact;

    /**
     * Verifies a CafePublishContact message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafePublishContact message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafePublishContact
     */
    public static fromObject(object: { [k: string]: any }): CafePublishContact;

    /**
     * Creates a plain object from a CafePublishContact message. Also converts values to other types if specified.
     * @param message CafePublishContact
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafePublishContact, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafePublishContact to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafePublishContactAck. */
export interface ICafePublishContactAck {

    /** CafePublishContactAck id */
    id: string;
}

/** Represents a CafePublishContactAck. */
export class CafePublishContactAck implements ICafePublishContactAck {

    /**
     * Constructs a new CafePublishContactAck.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafePublishContactAck);

    /** CafePublishContactAck id. */
    public id: string;

    /**
     * Creates a new CafePublishContactAck instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafePublishContactAck instance
     */
    public static create(properties?: ICafePublishContactAck): CafePublishContactAck;

    /**
     * Encodes the specified CafePublishContactAck message. Does not implicitly {@link CafePublishContactAck.verify|verify} messages.
     * @param message CafePublishContactAck message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafePublishContactAck, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafePublishContactAck message, length delimited. Does not implicitly {@link CafePublishContactAck.verify|verify} messages.
     * @param message CafePublishContactAck message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafePublishContactAck, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafePublishContactAck message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafePublishContactAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafePublishContactAck;

    /**
     * Decodes a CafePublishContactAck message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafePublishContactAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafePublishContactAck;

    /**
     * Verifies a CafePublishContactAck message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafePublishContactAck message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafePublishContactAck
     */
    public static fromObject(object: { [k: string]: any }): CafePublishContactAck;

    /**
     * Creates a plain object from a CafePublishContactAck message. Also converts values to other types if specified.
     * @param message CafePublishContactAck
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafePublishContactAck, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafePublishContactAck to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeStore. */
export interface ICafeStore {

    /** CafeStore token */
    token: string;

    /** CafeStore cids */
    cids: string[];
}

/** Represents a CafeStore. */
export class CafeStore implements ICafeStore {

    /**
     * Constructs a new CafeStore.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeStore);

    /** CafeStore token. */
    public token: string;

    /** CafeStore cids. */
    public cids: string[];

    /**
     * Creates a new CafeStore instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeStore instance
     */
    public static create(properties?: ICafeStore): CafeStore;

    /**
     * Encodes the specified CafeStore message. Does not implicitly {@link CafeStore.verify|verify} messages.
     * @param message CafeStore message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeStore, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeStore message, length delimited. Does not implicitly {@link CafeStore.verify|verify} messages.
     * @param message CafeStore message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeStore, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeStore message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeStore
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeStore;

    /**
     * Decodes a CafeStore message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeStore
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeStore;

    /**
     * Verifies a CafeStore message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeStore message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeStore
     */
    public static fromObject(object: { [k: string]: any }): CafeStore;

    /**
     * Creates a plain object from a CafeStore message. Also converts values to other types if specified.
     * @param message CafeStore
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeStore, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeStore to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeObjectList. */
export interface ICafeObjectList {

    /** CafeObjectList cids */
    cids: string[];
}

/** Represents a CafeObjectList. */
export class CafeObjectList implements ICafeObjectList {

    /**
     * Constructs a new CafeObjectList.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeObjectList);

    /** CafeObjectList cids. */
    public cids: string[];

    /**
     * Creates a new CafeObjectList instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeObjectList instance
     */
    public static create(properties?: ICafeObjectList): CafeObjectList;

    /**
     * Encodes the specified CafeObjectList message. Does not implicitly {@link CafeObjectList.verify|verify} messages.
     * @param message CafeObjectList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeObjectList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeObjectList message, length delimited. Does not implicitly {@link CafeObjectList.verify|verify} messages.
     * @param message CafeObjectList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeObjectList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeObjectList message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeObjectList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeObjectList;

    /**
     * Decodes a CafeObjectList message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeObjectList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeObjectList;

    /**
     * Verifies a CafeObjectList message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeObjectList message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeObjectList
     */
    public static fromObject(object: { [k: string]: any }): CafeObjectList;

    /**
     * Creates a plain object from a CafeObjectList message. Also converts values to other types if specified.
     * @param message CafeObjectList
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeObjectList, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeObjectList to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeObject. */
export interface ICafeObject {

    /** CafeObject token */
    token: string;

    /** CafeObject cid */
    cid: string;

    /** CafeObject data */
    data: Uint8Array;

    /** CafeObject node */
    node: Uint8Array;
}

/** Represents a CafeObject. */
export class CafeObject implements ICafeObject {

    /**
     * Constructs a new CafeObject.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeObject);

    /** CafeObject token. */
    public token: string;

    /** CafeObject cid. */
    public cid: string;

    /** CafeObject data. */
    public data: Uint8Array;

    /** CafeObject node. */
    public node: Uint8Array;

    /**
     * Creates a new CafeObject instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeObject instance
     */
    public static create(properties?: ICafeObject): CafeObject;

    /**
     * Encodes the specified CafeObject message. Does not implicitly {@link CafeObject.verify|verify} messages.
     * @param message CafeObject message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeObject, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeObject message, length delimited. Does not implicitly {@link CafeObject.verify|verify} messages.
     * @param message CafeObject message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeObject, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeObject message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeObject
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeObject;

    /**
     * Decodes a CafeObject message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeObject
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeObject;

    /**
     * Verifies a CafeObject message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeObject message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeObject
     */
    public static fromObject(object: { [k: string]: any }): CafeObject;

    /**
     * Creates a plain object from a CafeObject message. Also converts values to other types if specified.
     * @param message CafeObject
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeObject, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeObject to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeStoreThread. */
export interface ICafeStoreThread {

    /** CafeStoreThread token */
    token: string;

    /** CafeStoreThread id */
    id: string;

    /** CafeStoreThread ciphertext */
    ciphertext: Uint8Array;
}

/** Represents a CafeStoreThread. */
export class CafeStoreThread implements ICafeStoreThread {

    /**
     * Constructs a new CafeStoreThread.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeStoreThread);

    /** CafeStoreThread token. */
    public token: string;

    /** CafeStoreThread id. */
    public id: string;

    /** CafeStoreThread ciphertext. */
    public ciphertext: Uint8Array;

    /**
     * Creates a new CafeStoreThread instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeStoreThread instance
     */
    public static create(properties?: ICafeStoreThread): CafeStoreThread;

    /**
     * Encodes the specified CafeStoreThread message. Does not implicitly {@link CafeStoreThread.verify|verify} messages.
     * @param message CafeStoreThread message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeStoreThread, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeStoreThread message, length delimited. Does not implicitly {@link CafeStoreThread.verify|verify} messages.
     * @param message CafeStoreThread message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeStoreThread, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeStoreThread message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeStoreThread
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeStoreThread;

    /**
     * Decodes a CafeStoreThread message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeStoreThread
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeStoreThread;

    /**
     * Verifies a CafeStoreThread message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeStoreThread message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeStoreThread
     */
    public static fromObject(object: { [k: string]: any }): CafeStoreThread;

    /**
     * Creates a plain object from a CafeStoreThread message. Also converts values to other types if specified.
     * @param message CafeStoreThread
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeStoreThread, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeStoreThread to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeStored. */
export interface ICafeStored {

    /** CafeStored id */
    id: string;
}

/** Represents a CafeStored. */
export class CafeStored implements ICafeStored {

    /**
     * Constructs a new CafeStored.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeStored);

    /** CafeStored id. */
    public id: string;

    /**
     * Creates a new CafeStored instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeStored instance
     */
    public static create(properties?: ICafeStored): CafeStored;

    /**
     * Encodes the specified CafeStored message. Does not implicitly {@link CafeStored.verify|verify} messages.
     * @param message CafeStored message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeStored, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeStored message, length delimited. Does not implicitly {@link CafeStored.verify|verify} messages.
     * @param message CafeStored message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeStored, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeStored message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeStored
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeStored;

    /**
     * Decodes a CafeStored message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeStored
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeStored;

    /**
     * Verifies a CafeStored message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeStored message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeStored
     */
    public static fromObject(object: { [k: string]: any }): CafeStored;

    /**
     * Creates a plain object from a CafeStored message. Also converts values to other types if specified.
     * @param message CafeStored
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeStored, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeStored to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeDeliverMessage. */
export interface ICafeDeliverMessage {

    /** CafeDeliverMessage id */
    id: string;

    /** CafeDeliverMessage client */
    client: string;
}

/** Represents a CafeDeliverMessage. */
export class CafeDeliverMessage implements ICafeDeliverMessage {

    /**
     * Constructs a new CafeDeliverMessage.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeDeliverMessage);

    /** CafeDeliverMessage id. */
    public id: string;

    /** CafeDeliverMessage client. */
    public client: string;

    /**
     * Creates a new CafeDeliverMessage instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeDeliverMessage instance
     */
    public static create(properties?: ICafeDeliverMessage): CafeDeliverMessage;

    /**
     * Encodes the specified CafeDeliverMessage message. Does not implicitly {@link CafeDeliverMessage.verify|verify} messages.
     * @param message CafeDeliverMessage message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeDeliverMessage, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeDeliverMessage message, length delimited. Does not implicitly {@link CafeDeliverMessage.verify|verify} messages.
     * @param message CafeDeliverMessage message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeDeliverMessage, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeDeliverMessage message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeDeliverMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeDeliverMessage;

    /**
     * Decodes a CafeDeliverMessage message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeDeliverMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeDeliverMessage;

    /**
     * Verifies a CafeDeliverMessage message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeDeliverMessage message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeDeliverMessage
     */
    public static fromObject(object: { [k: string]: any }): CafeDeliverMessage;

    /**
     * Creates a plain object from a CafeDeliverMessage message. Also converts values to other types if specified.
     * @param message CafeDeliverMessage
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeDeliverMessage, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeDeliverMessage to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeCheckMessages. */
export interface ICafeCheckMessages {

    /** CafeCheckMessages token */
    token: string;
}

/** Represents a CafeCheckMessages. */
export class CafeCheckMessages implements ICafeCheckMessages {

    /**
     * Constructs a new CafeCheckMessages.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeCheckMessages);

    /** CafeCheckMessages token. */
    public token: string;

    /**
     * Creates a new CafeCheckMessages instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeCheckMessages instance
     */
    public static create(properties?: ICafeCheckMessages): CafeCheckMessages;

    /**
     * Encodes the specified CafeCheckMessages message. Does not implicitly {@link CafeCheckMessages.verify|verify} messages.
     * @param message CafeCheckMessages message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeCheckMessages, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeCheckMessages message, length delimited. Does not implicitly {@link CafeCheckMessages.verify|verify} messages.
     * @param message CafeCheckMessages message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeCheckMessages, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeCheckMessages message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeCheckMessages
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeCheckMessages;

    /**
     * Decodes a CafeCheckMessages message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeCheckMessages
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeCheckMessages;

    /**
     * Verifies a CafeCheckMessages message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeCheckMessages message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeCheckMessages
     */
    public static fromObject(object: { [k: string]: any }): CafeCheckMessages;

    /**
     * Creates a plain object from a CafeCheckMessages message. Also converts values to other types if specified.
     * @param message CafeCheckMessages
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeCheckMessages, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeCheckMessages to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeMessages. */
export interface ICafeMessages {

    /** CafeMessages messages */
    messages: ICafeMessage[];
}

/** Represents a CafeMessages. */
export class CafeMessages implements ICafeMessages {

    /**
     * Constructs a new CafeMessages.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeMessages);

    /** CafeMessages messages. */
    public messages: ICafeMessage[];

    /**
     * Creates a new CafeMessages instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeMessages instance
     */
    public static create(properties?: ICafeMessages): CafeMessages;

    /**
     * Encodes the specified CafeMessages message. Does not implicitly {@link CafeMessages.verify|verify} messages.
     * @param message CafeMessages message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeMessages, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeMessages message, length delimited. Does not implicitly {@link CafeMessages.verify|verify} messages.
     * @param message CafeMessages message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeMessages, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeMessages message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeMessages
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeMessages;

    /**
     * Decodes a CafeMessages message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeMessages
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeMessages;

    /**
     * Verifies a CafeMessages message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeMessages message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeMessages
     */
    public static fromObject(object: { [k: string]: any }): CafeMessages;

    /**
     * Creates a plain object from a CafeMessages message. Also converts values to other types if specified.
     * @param message CafeMessages
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeMessages, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeMessages to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeDeleteMessages. */
export interface ICafeDeleteMessages {

    /** CafeDeleteMessages token */
    token: string;
}

/** Represents a CafeDeleteMessages. */
export class CafeDeleteMessages implements ICafeDeleteMessages {

    /**
     * Constructs a new CafeDeleteMessages.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeDeleteMessages);

    /** CafeDeleteMessages token. */
    public token: string;

    /**
     * Creates a new CafeDeleteMessages instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeDeleteMessages instance
     */
    public static create(properties?: ICafeDeleteMessages): CafeDeleteMessages;

    /**
     * Encodes the specified CafeDeleteMessages message. Does not implicitly {@link CafeDeleteMessages.verify|verify} messages.
     * @param message CafeDeleteMessages message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeDeleteMessages, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeDeleteMessages message, length delimited. Does not implicitly {@link CafeDeleteMessages.verify|verify} messages.
     * @param message CafeDeleteMessages message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeDeleteMessages, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeDeleteMessages message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeDeleteMessages
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeDeleteMessages;

    /**
     * Decodes a CafeDeleteMessages message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeDeleteMessages
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeDeleteMessages;

    /**
     * Verifies a CafeDeleteMessages message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeDeleteMessages message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeDeleteMessages
     */
    public static fromObject(object: { [k: string]: any }): CafeDeleteMessages;

    /**
     * Creates a plain object from a CafeDeleteMessages message. Also converts values to other types if specified.
     * @param message CafeDeleteMessages
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeDeleteMessages, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeDeleteMessages to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeDeleteMessagesAck. */
export interface ICafeDeleteMessagesAck {

    /** CafeDeleteMessagesAck more */
    more: boolean;
}

/** Represents a CafeDeleteMessagesAck. */
export class CafeDeleteMessagesAck implements ICafeDeleteMessagesAck {

    /**
     * Constructs a new CafeDeleteMessagesAck.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeDeleteMessagesAck);

    /** CafeDeleteMessagesAck more. */
    public more: boolean;

    /**
     * Creates a new CafeDeleteMessagesAck instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeDeleteMessagesAck instance
     */
    public static create(properties?: ICafeDeleteMessagesAck): CafeDeleteMessagesAck;

    /**
     * Encodes the specified CafeDeleteMessagesAck message. Does not implicitly {@link CafeDeleteMessagesAck.verify|verify} messages.
     * @param message CafeDeleteMessagesAck message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeDeleteMessagesAck, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeDeleteMessagesAck message, length delimited. Does not implicitly {@link CafeDeleteMessagesAck.verify|verify} messages.
     * @param message CafeDeleteMessagesAck message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeDeleteMessagesAck, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeDeleteMessagesAck message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeDeleteMessagesAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeDeleteMessagesAck;

    /**
     * Decodes a CafeDeleteMessagesAck message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeDeleteMessagesAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeDeleteMessagesAck;

    /**
     * Verifies a CafeDeleteMessagesAck message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeDeleteMessagesAck message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeDeleteMessagesAck
     */
    public static fromObject(object: { [k: string]: any }): CafeDeleteMessagesAck;

    /**
     * Creates a plain object from a CafeDeleteMessagesAck message. Also converts values to other types if specified.
     * @param message CafeDeleteMessagesAck
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeDeleteMessagesAck, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeDeleteMessagesAck to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a Contact. */
export interface IContact {

    /** Contact id */
    id: string;

    /** Contact address */
    address: string;

    /** Contact username */
    username: string;

    /** Contact avatar */
    avatar: string;

    /** Contact inboxes */
    inboxes: ICafe[];

    /** Contact created */
    created: google.protobuf.ITimestamp;

    /** Contact updated */
    updated: google.protobuf.ITimestamp;

    /** Contact threads */
    threads: string[];
}

/** Represents a Contact. */
export class Contact implements IContact {

    /**
     * Constructs a new Contact.
     * @param [properties] Properties to set
     */
    constructor(properties?: IContact);

    /** Contact id. */
    public id: string;

    /** Contact address. */
    public address: string;

    /** Contact username. */
    public username: string;

    /** Contact avatar. */
    public avatar: string;

    /** Contact inboxes. */
    public inboxes: ICafe[];

    /** Contact created. */
    public created: google.protobuf.ITimestamp;

    /** Contact updated. */
    public updated: google.protobuf.ITimestamp;

    /** Contact threads. */
    public threads: string[];

    /**
     * Creates a new Contact instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Contact instance
     */
    public static create(properties?: IContact): Contact;

    /**
     * Encodes the specified Contact message. Does not implicitly {@link Contact.verify|verify} messages.
     * @param message Contact message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IContact, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Contact message, length delimited. Does not implicitly {@link Contact.verify|verify} messages.
     * @param message Contact message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IContact, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Contact message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Contact
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Contact;

    /**
     * Decodes a Contact message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Contact
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Contact;

    /**
     * Verifies a Contact message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Contact message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Contact
     */
    public static fromObject(object: { [k: string]: any }): Contact;

    /**
     * Creates a plain object from a Contact message. Also converts values to other types if specified.
     * @param message Contact
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Contact, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Contact to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a ContactList. */
export interface IContactList {

    /** ContactList items */
    items: IContact[];
}

/** Represents a ContactList. */
export class ContactList implements IContactList {

    /**
     * Constructs a new ContactList.
     * @param [properties] Properties to set
     */
    constructor(properties?: IContactList);

    /** ContactList items. */
    public items: IContact[];

    /**
     * Creates a new ContactList instance using the specified properties.
     * @param [properties] Properties to set
     * @returns ContactList instance
     */
    public static create(properties?: IContactList): ContactList;

    /**
     * Encodes the specified ContactList message. Does not implicitly {@link ContactList.verify|verify} messages.
     * @param message ContactList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IContactList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified ContactList message, length delimited. Does not implicitly {@link ContactList.verify|verify} messages.
     * @param message ContactList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IContactList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a ContactList message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns ContactList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): ContactList;

    /**
     * Decodes a ContactList message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns ContactList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): ContactList;

    /**
     * Verifies a ContactList message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a ContactList message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns ContactList
     */
    public static fromObject(object: { [k: string]: any }): ContactList;

    /**
     * Creates a plain object from a ContactList message. Also converts values to other types if specified.
     * @param message ContactList
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: ContactList, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this ContactList to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a User. */
export interface IUser {

    /** User address */
    address: string;

    /** User name */
    name: string;

    /** User avatar */
    avatar: string;
}

/** Represents a User. */
export class User implements IUser {

    /**
     * Constructs a new User.
     * @param [properties] Properties to set
     */
    constructor(properties?: IUser);

    /** User address. */
    public address: string;

    /** User name. */
    public name: string;

    /** User avatar. */
    public avatar: string;

    /**
     * Creates a new User instance using the specified properties.
     * @param [properties] Properties to set
     * @returns User instance
     */
    public static create(properties?: IUser): User;

    /**
     * Encodes the specified User message. Does not implicitly {@link User.verify|verify} messages.
     * @param message User message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IUser, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified User message, length delimited. Does not implicitly {@link User.verify|verify} messages.
     * @param message User message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IUser, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a User message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns User
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): User;

    /**
     * Decodes a User message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns User
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): User;

    /**
     * Verifies a User message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a User message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns User
     */
    public static fromObject(object: { [k: string]: any }): User;

    /**
     * Creates a plain object from a User message. Also converts values to other types if specified.
     * @param message User
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: User, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this User to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a Thread. */
export interface IThread {

    /** Thread id */
    id: string;

    /** Thread key */
    key: string;

    /** Thread sk */
    sk: Uint8Array;

    /** Thread name */
    name: string;

    /** Thread schema */
    schema: string;

    /** Thread initiator */
    initiator: string;

    /** Thread type */
    type: Thread.Type;

    /** Thread sharing */
    sharing: Thread.Sharing;

    /** Thread members */
    members: string[];

    /** Thread state */
    state: Thread.State;

    /** Thread head */
    head: string;

    /** Thread headBlock */
    headBlock: IBlock;

    /** Thread schemaNode */
    schemaNode: INode;

    /** Thread blockCount */
    blockCount: number;

    /** Thread peerCount */
    peerCount: number;
}

/** Represents a Thread. */
export class Thread implements IThread {

    /**
     * Constructs a new Thread.
     * @param [properties] Properties to set
     */
    constructor(properties?: IThread);

    /** Thread id. */
    public id: string;

    /** Thread key. */
    public key: string;

    /** Thread sk. */
    public sk: Uint8Array;

    /** Thread name. */
    public name: string;

    /** Thread schema. */
    public schema: string;

    /** Thread initiator. */
    public initiator: string;

    /** Thread type. */
    public type: Thread.Type;

    /** Thread sharing. */
    public sharing: Thread.Sharing;

    /** Thread members. */
    public members: string[];

    /** Thread state. */
    public state: Thread.State;

    /** Thread head. */
    public head: string;

    /** Thread headBlock. */
    public headBlock: IBlock;

    /** Thread schemaNode. */
    public schemaNode: INode;

    /** Thread blockCount. */
    public blockCount: number;

    /** Thread peerCount. */
    public peerCount: number;

    /**
     * Creates a new Thread instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Thread instance
     */
    public static create(properties?: IThread): Thread;

    /**
     * Encodes the specified Thread message. Does not implicitly {@link Thread.verify|verify} messages.
     * @param message Thread message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IThread, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Thread message, length delimited. Does not implicitly {@link Thread.verify|verify} messages.
     * @param message Thread message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IThread, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Thread message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Thread
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Thread;

    /**
     * Decodes a Thread message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Thread
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Thread;

    /**
     * Verifies a Thread message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Thread message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Thread
     */
    public static fromObject(object: { [k: string]: any }): Thread;

    /**
     * Creates a plain object from a Thread message. Also converts values to other types if specified.
     * @param message Thread
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Thread, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Thread to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

export namespace Thread {

    /** Type enum. */
    enum Type {
        Private = 0,
        ReadOnly = 1,
        Public = 2,
        Open = 3
    }

    /** Sharing enum. */
    enum Sharing {
        NotShared = 0,
        InviteOnly = 1,
        Shared = 2
    }

    /** State enum. */
    enum State {
        LoadingBehind = 0,
        Loaded = 1,
        LoadingAhead = 2
    }
}

/** Properties of a ThreadList. */
export interface IThreadList {

    /** ThreadList items */
    items: IThread[];
}

/** Represents a ThreadList. */
export class ThreadList implements IThreadList {

    /**
     * Constructs a new ThreadList.
     * @param [properties] Properties to set
     */
    constructor(properties?: IThreadList);

    /** ThreadList items. */
    public items: IThread[];

    /**
     * Creates a new ThreadList instance using the specified properties.
     * @param [properties] Properties to set
     * @returns ThreadList instance
     */
    public static create(properties?: IThreadList): ThreadList;

    /**
     * Encodes the specified ThreadList message. Does not implicitly {@link ThreadList.verify|verify} messages.
     * @param message ThreadList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IThreadList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified ThreadList message, length delimited. Does not implicitly {@link ThreadList.verify|verify} messages.
     * @param message ThreadList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IThreadList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a ThreadList message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns ThreadList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): ThreadList;

    /**
     * Decodes a ThreadList message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns ThreadList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): ThreadList;

    /**
     * Verifies a ThreadList message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a ThreadList message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns ThreadList
     */
    public static fromObject(object: { [k: string]: any }): ThreadList;

    /**
     * Creates a plain object from a ThreadList message. Also converts values to other types if specified.
     * @param message ThreadList
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: ThreadList, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this ThreadList to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a ThreadPeer. */
export interface IThreadPeer {

    /** ThreadPeer id */
    id: string;

    /** ThreadPeer thread */
    thread: string;

    /** ThreadPeer welcomed */
    welcomed: boolean;
}

/** Represents a ThreadPeer. */
export class ThreadPeer implements IThreadPeer {

    /**
     * Constructs a new ThreadPeer.
     * @param [properties] Properties to set
     */
    constructor(properties?: IThreadPeer);

    /** ThreadPeer id. */
    public id: string;

    /** ThreadPeer thread. */
    public thread: string;

    /** ThreadPeer welcomed. */
    public welcomed: boolean;

    /**
     * Creates a new ThreadPeer instance using the specified properties.
     * @param [properties] Properties to set
     * @returns ThreadPeer instance
     */
    public static create(properties?: IThreadPeer): ThreadPeer;

    /**
     * Encodes the specified ThreadPeer message. Does not implicitly {@link ThreadPeer.verify|verify} messages.
     * @param message ThreadPeer message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IThreadPeer, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified ThreadPeer message, length delimited. Does not implicitly {@link ThreadPeer.verify|verify} messages.
     * @param message ThreadPeer message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IThreadPeer, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a ThreadPeer message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns ThreadPeer
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): ThreadPeer;

    /**
     * Decodes a ThreadPeer message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns ThreadPeer
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): ThreadPeer;

    /**
     * Verifies a ThreadPeer message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a ThreadPeer message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns ThreadPeer
     */
    public static fromObject(object: { [k: string]: any }): ThreadPeer;

    /**
     * Creates a plain object from a ThreadPeer message. Also converts values to other types if specified.
     * @param message ThreadPeer
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: ThreadPeer, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this ThreadPeer to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a Block. */
export interface IBlock {

    /** Block id */
    id: string;

    /** Block thread */
    thread: string;

    /** Block author */
    author: string;

    /** Block type */
    type: Block.BlockType;

    /** Block date */
    date: google.protobuf.ITimestamp;

    /** Block parents */
    parents: string[];

    /** Block target */
    target: string;

    /** Block body */
    body: string;

    /** Block user */
    user: IUser;
}

/** Represents a Block. */
export class Block implements IBlock {

    /**
     * Constructs a new Block.
     * @param [properties] Properties to set
     */
    constructor(properties?: IBlock);

    /** Block id. */
    public id: string;

    /** Block thread. */
    public thread: string;

    /** Block author. */
    public author: string;

    /** Block type. */
    public type: Block.BlockType;

    /** Block date. */
    public date: google.protobuf.ITimestamp;

    /** Block parents. */
    public parents: string[];

    /** Block target. */
    public target: string;

    /** Block body. */
    public body: string;

    /** Block user. */
    public user: IUser;

    /**
     * Creates a new Block instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Block instance
     */
    public static create(properties?: IBlock): Block;

    /**
     * Encodes the specified Block message. Does not implicitly {@link Block.verify|verify} messages.
     * @param message Block message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IBlock, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Block message, length delimited. Does not implicitly {@link Block.verify|verify} messages.
     * @param message Block message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IBlock, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Block message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Block
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Block;

    /**
     * Decodes a Block message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Block
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Block;

    /**
     * Verifies a Block message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Block message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Block
     */
    public static fromObject(object: { [k: string]: any }): Block;

    /**
     * Creates a plain object from a Block message. Also converts values to other types if specified.
     * @param message Block
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Block, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Block to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

export namespace Block {

    /** BlockType enum. */
    enum BlockType {
        MERGE = 0,
        IGNORE = 1,
        FLAG = 2,
        JOIN = 3,
        ANNOUNCE = 4,
        LEAVE = 5,
        MESSAGE = 6,
        FILES = 7,
        COMMENT = 8,
        LIKE = 9,
        INVITE = 50
    }
}

/** Properties of a BlockList. */
export interface IBlockList {

    /** BlockList items */
    items: IBlock[];
}

/** Represents a BlockList. */
export class BlockList implements IBlockList {

    /**
     * Constructs a new BlockList.
     * @param [properties] Properties to set
     */
    constructor(properties?: IBlockList);

    /** BlockList items. */
    public items: IBlock[];

    /**
     * Creates a new BlockList instance using the specified properties.
     * @param [properties] Properties to set
     * @returns BlockList instance
     */
    public static create(properties?: IBlockList): BlockList;

    /**
     * Encodes the specified BlockList message. Does not implicitly {@link BlockList.verify|verify} messages.
     * @param message BlockList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IBlockList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified BlockList message, length delimited. Does not implicitly {@link BlockList.verify|verify} messages.
     * @param message BlockList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IBlockList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a BlockList message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns BlockList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): BlockList;

    /**
     * Decodes a BlockList message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns BlockList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): BlockList;

    /**
     * Verifies a BlockList message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a BlockList message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns BlockList
     */
    public static fromObject(object: { [k: string]: any }): BlockList;

    /**
     * Creates a plain object from a BlockList message. Also converts values to other types if specified.
     * @param message BlockList
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: BlockList, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this BlockList to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a BlockMessage. */
export interface IBlockMessage {

    /** BlockMessage id */
    id: string;

    /** BlockMessage peer */
    peer: string;

    /** BlockMessage env */
    env: IEnvelope;

    /** BlockMessage date */
    date: google.protobuf.ITimestamp;
}

/** Represents a BlockMessage. */
export class BlockMessage implements IBlockMessage {

    /**
     * Constructs a new BlockMessage.
     * @param [properties] Properties to set
     */
    constructor(properties?: IBlockMessage);

    /** BlockMessage id. */
    public id: string;

    /** BlockMessage peer. */
    public peer: string;

    /** BlockMessage env. */
    public env: IEnvelope;

    /** BlockMessage date. */
    public date: google.protobuf.ITimestamp;

    /**
     * Creates a new BlockMessage instance using the specified properties.
     * @param [properties] Properties to set
     * @returns BlockMessage instance
     */
    public static create(properties?: IBlockMessage): BlockMessage;

    /**
     * Encodes the specified BlockMessage message. Does not implicitly {@link BlockMessage.verify|verify} messages.
     * @param message BlockMessage message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IBlockMessage, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified BlockMessage message, length delimited. Does not implicitly {@link BlockMessage.verify|verify} messages.
     * @param message BlockMessage message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IBlockMessage, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a BlockMessage message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns BlockMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): BlockMessage;

    /**
     * Decodes a BlockMessage message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns BlockMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): BlockMessage;

    /**
     * Verifies a BlockMessage message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a BlockMessage message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns BlockMessage
     */
    public static fromObject(object: { [k: string]: any }): BlockMessage;

    /**
     * Creates a plain object from a BlockMessage message. Also converts values to other types if specified.
     * @param message BlockMessage
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: BlockMessage, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this BlockMessage to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of an Invite. */
export interface IInvite {

    /** Invite id */
    id: string;

    /** Invite block */
    block: Uint8Array;

    /** Invite name */
    name: string;

    /** Invite inviter */
    inviter: IContact;

    /** Invite date */
    date: google.protobuf.ITimestamp;
}

/** Represents an Invite. */
export class Invite implements IInvite {

    /**
     * Constructs a new Invite.
     * @param [properties] Properties to set
     */
    constructor(properties?: IInvite);

    /** Invite id. */
    public id: string;

    /** Invite block. */
    public block: Uint8Array;

    /** Invite name. */
    public name: string;

    /** Invite inviter. */
    public inviter: IContact;

    /** Invite date. */
    public date: google.protobuf.ITimestamp;

    /**
     * Creates a new Invite instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Invite instance
     */
    public static create(properties?: IInvite): Invite;

    /**
     * Encodes the specified Invite message. Does not implicitly {@link Invite.verify|verify} messages.
     * @param message Invite message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IInvite, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Invite message, length delimited. Does not implicitly {@link Invite.verify|verify} messages.
     * @param message Invite message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IInvite, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes an Invite message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Invite
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Invite;

    /**
     * Decodes an Invite message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Invite
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Invite;

    /**
     * Verifies an Invite message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates an Invite message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Invite
     */
    public static fromObject(object: { [k: string]: any }): Invite;

    /**
     * Creates a plain object from an Invite message. Also converts values to other types if specified.
     * @param message Invite
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Invite, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Invite to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of an InviteList. */
export interface IInviteList {

    /** InviteList items */
    items: IInvite[];
}

/** Represents an InviteList. */
export class InviteList implements IInviteList {

    /**
     * Constructs a new InviteList.
     * @param [properties] Properties to set
     */
    constructor(properties?: IInviteList);

    /** InviteList items. */
    public items: IInvite[];

    /**
     * Creates a new InviteList instance using the specified properties.
     * @param [properties] Properties to set
     * @returns InviteList instance
     */
    public static create(properties?: IInviteList): InviteList;

    /**
     * Encodes the specified InviteList message. Does not implicitly {@link InviteList.verify|verify} messages.
     * @param message InviteList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IInviteList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified InviteList message, length delimited. Does not implicitly {@link InviteList.verify|verify} messages.
     * @param message InviteList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IInviteList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes an InviteList message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns InviteList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): InviteList;

    /**
     * Decodes an InviteList message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns InviteList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): InviteList;

    /**
     * Verifies an InviteList message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates an InviteList message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns InviteList
     */
    public static fromObject(object: { [k: string]: any }): InviteList;

    /**
     * Creates a plain object from an InviteList message. Also converts values to other types if specified.
     * @param message InviteList
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: InviteList, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this InviteList to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a FileIndex. */
export interface IFileIndex {

    /** FileIndex mill */
    mill: string;

    /** FileIndex checksum */
    checksum: string;

    /** FileIndex source */
    source: string;

    /** FileIndex opts */
    opts: string;

    /** FileIndex hash */
    hash: string;

    /** FileIndex key */
    key: string;

    /** FileIndex media */
    media: string;

    /** FileIndex name */
    name: string;

    /** FileIndex size */
    size: (number|Long);

    /** FileIndex added */
    added: google.protobuf.ITimestamp;

    /** FileIndex meta */
    meta: google.protobuf.IStruct;

    /** FileIndex targets */
    targets: string[];
}

/** Represents a FileIndex. */
export class FileIndex implements IFileIndex {

    /**
     * Constructs a new FileIndex.
     * @param [properties] Properties to set
     */
    constructor(properties?: IFileIndex);

    /** FileIndex mill. */
    public mill: string;

    /** FileIndex checksum. */
    public checksum: string;

    /** FileIndex source. */
    public source: string;

    /** FileIndex opts. */
    public opts: string;

    /** FileIndex hash. */
    public hash: string;

    /** FileIndex key. */
    public key: string;

    /** FileIndex media. */
    public media: string;

    /** FileIndex name. */
    public name: string;

    /** FileIndex size. */
    public size: (number|Long);

    /** FileIndex added. */
    public added: google.protobuf.ITimestamp;

    /** FileIndex meta. */
    public meta: google.protobuf.IStruct;

    /** FileIndex targets. */
    public targets: string[];

    /**
     * Creates a new FileIndex instance using the specified properties.
     * @param [properties] Properties to set
     * @returns FileIndex instance
     */
    public static create(properties?: IFileIndex): FileIndex;

    /**
     * Encodes the specified FileIndex message. Does not implicitly {@link FileIndex.verify|verify} messages.
     * @param message FileIndex message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IFileIndex, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified FileIndex message, length delimited. Does not implicitly {@link FileIndex.verify|verify} messages.
     * @param message FileIndex message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IFileIndex, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a FileIndex message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns FileIndex
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): FileIndex;

    /**
     * Decodes a FileIndex message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns FileIndex
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): FileIndex;

    /**
     * Verifies a FileIndex message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a FileIndex message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns FileIndex
     */
    public static fromObject(object: { [k: string]: any }): FileIndex;

    /**
     * Creates a plain object from a FileIndex message. Also converts values to other types if specified.
     * @param message FileIndex
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: FileIndex, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this FileIndex to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a Node. */
export interface INode {

    /** Node name */
    name: string;

    /** Node pin */
    pin: boolean;

    /** Node plaintext */
    plaintext: boolean;

    /** Node mill */
    mill: string;

    /** Node opts */
    opts: { [k: string]: string };

    /** Node jsonSchema */
    jsonSchema: google.protobuf.IStruct;

    /** Node links */
    links: { [k: string]: ILink };
}

/** Represents a Node. */
export class Node implements INode {

    /**
     * Constructs a new Node.
     * @param [properties] Properties to set
     */
    constructor(properties?: INode);

    /** Node name. */
    public name: string;

    /** Node pin. */
    public pin: boolean;

    /** Node plaintext. */
    public plaintext: boolean;

    /** Node mill. */
    public mill: string;

    /** Node opts. */
    public opts: { [k: string]: string };

    /** Node jsonSchema. */
    public jsonSchema: google.protobuf.IStruct;

    /** Node links. */
    public links: { [k: string]: ILink };

    /**
     * Creates a new Node instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Node instance
     */
    public static create(properties?: INode): Node;

    /**
     * Encodes the specified Node message. Does not implicitly {@link Node.verify|verify} messages.
     * @param message Node message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: INode, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Node message, length delimited. Does not implicitly {@link Node.verify|verify} messages.
     * @param message Node message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: INode, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Node message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Node
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Node;

    /**
     * Decodes a Node message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Node
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Node;

    /**
     * Verifies a Node message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Node message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Node
     */
    public static fromObject(object: { [k: string]: any }): Node;

    /**
     * Creates a plain object from a Node message. Also converts values to other types if specified.
     * @param message Node
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Node, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Node to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a Link. */
export interface ILink {

    /** Link use */
    use: string;

    /** Link pin */
    pin: boolean;

    /** Link plaintext */
    plaintext: boolean;

    /** Link mill */
    mill: string;

    /** Link opts */
    opts: { [k: string]: string };

    /** Link jsonSchema */
    jsonSchema: google.protobuf.IStruct;
}

/** Represents a Link. */
export class Link implements ILink {

    /**
     * Constructs a new Link.
     * @param [properties] Properties to set
     */
    constructor(properties?: ILink);

    /** Link use. */
    public use: string;

    /** Link pin. */
    public pin: boolean;

    /** Link plaintext. */
    public plaintext: boolean;

    /** Link mill. */
    public mill: string;

    /** Link opts. */
    public opts: { [k: string]: string };

    /** Link jsonSchema. */
    public jsonSchema: google.protobuf.IStruct;

    /**
     * Creates a new Link instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Link instance
     */
    public static create(properties?: ILink): Link;

    /**
     * Encodes the specified Link message. Does not implicitly {@link Link.verify|verify} messages.
     * @param message Link message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ILink, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Link message, length delimited. Does not implicitly {@link Link.verify|verify} messages.
     * @param message Link message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ILink, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Link message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Link
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Link;

    /**
     * Decodes a Link message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Link
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Link;

    /**
     * Verifies a Link message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Link message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Link
     */
    public static fromObject(object: { [k: string]: any }): Link;

    /**
     * Creates a plain object from a Link message. Also converts values to other types if specified.
     * @param message Link
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Link, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Link to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a Notification. */
export interface INotification {

    /** Notification id */
    id: string;

    /** Notification date */
    date: google.protobuf.ITimestamp;

    /** Notification actor */
    actor: string;

    /** Notification subject */
    subject: string;

    /** Notification subjectDesc */
    subjectDesc: string;

    /** Notification block */
    block: string;

    /** Notification target */
    target: string;

    /** Notification type */
    type: Notification.Type;

    /** Notification body */
    body: string;

    /** Notification read */
    read: boolean;

    /** Notification user */
    user: IUser;
}

/** Represents a Notification. */
export class Notification implements INotification {

    /**
     * Constructs a new Notification.
     * @param [properties] Properties to set
     */
    constructor(properties?: INotification);

    /** Notification id. */
    public id: string;

    /** Notification date. */
    public date: google.protobuf.ITimestamp;

    /** Notification actor. */
    public actor: string;

    /** Notification subject. */
    public subject: string;

    /** Notification subjectDesc. */
    public subjectDesc: string;

    /** Notification block. */
    public block: string;

    /** Notification target. */
    public target: string;

    /** Notification type. */
    public type: Notification.Type;

    /** Notification body. */
    public body: string;

    /** Notification read. */
    public read: boolean;

    /** Notification user. */
    public user: IUser;

    /**
     * Creates a new Notification instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Notification instance
     */
    public static create(properties?: INotification): Notification;

    /**
     * Encodes the specified Notification message. Does not implicitly {@link Notification.verify|verify} messages.
     * @param message Notification message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: INotification, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Notification message, length delimited. Does not implicitly {@link Notification.verify|verify} messages.
     * @param message Notification message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: INotification, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Notification message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Notification
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Notification;

    /**
     * Decodes a Notification message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Notification
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Notification;

    /**
     * Verifies a Notification message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Notification message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Notification
     */
    public static fromObject(object: { [k: string]: any }): Notification;

    /**
     * Creates a plain object from a Notification message. Also converts values to other types if specified.
     * @param message Notification
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Notification, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Notification to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

export namespace Notification {

    /** Type enum. */
    enum Type {
        INVITE_RECEIVED = 0,
        ACCOUNT_PEER_JOINED = 1,
        PEER_JOINED = 2,
        PEER_LEFT = 3,
        MESSAGE_ADDED = 4,
        FILES_ADDED = 5,
        COMMENT_ADDED = 6,
        LIKE_ADDED = 7
    }
}

/** Properties of a NotificationList. */
export interface INotificationList {

    /** NotificationList items */
    items: INotification[];
}

/** Represents a NotificationList. */
export class NotificationList implements INotificationList {

    /**
     * Constructs a new NotificationList.
     * @param [properties] Properties to set
     */
    constructor(properties?: INotificationList);

    /** NotificationList items. */
    public items: INotification[];

    /**
     * Creates a new NotificationList instance using the specified properties.
     * @param [properties] Properties to set
     * @returns NotificationList instance
     */
    public static create(properties?: INotificationList): NotificationList;

    /**
     * Encodes the specified NotificationList message. Does not implicitly {@link NotificationList.verify|verify} messages.
     * @param message NotificationList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: INotificationList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified NotificationList message, length delimited. Does not implicitly {@link NotificationList.verify|verify} messages.
     * @param message NotificationList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: INotificationList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a NotificationList message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns NotificationList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): NotificationList;

    /**
     * Decodes a NotificationList message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns NotificationList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): NotificationList;

    /**
     * Verifies a NotificationList message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a NotificationList message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns NotificationList
     */
    public static fromObject(object: { [k: string]: any }): NotificationList;

    /**
     * Creates a plain object from a NotificationList message. Also converts values to other types if specified.
     * @param message NotificationList
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: NotificationList, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this NotificationList to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a Cafe. */
export interface ICafe {

    /** Cafe peer */
    peer: string;

    /** Cafe address */
    address: string;

    /** Cafe api */
    api: string;

    /** Cafe protocol */
    protocol: string;

    /** Cafe node */
    node: string;

    /** Cafe url */
    url: string;

    /** Cafe swarm */
    swarm: string[];
}

/** Represents a Cafe. */
export class Cafe implements ICafe {

    /**
     * Constructs a new Cafe.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafe);

    /** Cafe peer. */
    public peer: string;

    /** Cafe address. */
    public address: string;

    /** Cafe api. */
    public api: string;

    /** Cafe protocol. */
    public protocol: string;

    /** Cafe node. */
    public node: string;

    /** Cafe url. */
    public url: string;

    /** Cafe swarm. */
    public swarm: string[];

    /**
     * Creates a new Cafe instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Cafe instance
     */
    public static create(properties?: ICafe): Cafe;

    /**
     * Encodes the specified Cafe message. Does not implicitly {@link Cafe.verify|verify} messages.
     * @param message Cafe message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafe, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Cafe message, length delimited. Does not implicitly {@link Cafe.verify|verify} messages.
     * @param message Cafe message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafe, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Cafe message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Cafe
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Cafe;

    /**
     * Decodes a Cafe message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Cafe
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Cafe;

    /**
     * Verifies a Cafe message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Cafe message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Cafe
     */
    public static fromObject(object: { [k: string]: any }): Cafe;

    /**
     * Creates a plain object from a Cafe message. Also converts values to other types if specified.
     * @param message Cafe
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Cafe, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Cafe to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeSession. */
export interface ICafeSession {

    /** CafeSession id */
    id: string;

    /** CafeSession access */
    access: string;

    /** CafeSession exp */
    exp: google.protobuf.ITimestamp;

    /** CafeSession refresh */
    refresh: string;

    /** CafeSession rexp */
    rexp: google.protobuf.ITimestamp;

    /** CafeSession subject */
    subject: string;

    /** CafeSession type */
    type: string;

    /** CafeSession cafe */
    cafe: ICafe;
}

/** Represents a CafeSession. */
export class CafeSession implements ICafeSession {

    /**
     * Constructs a new CafeSession.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeSession);

    /** CafeSession id. */
    public id: string;

    /** CafeSession access. */
    public access: string;

    /** CafeSession exp. */
    public exp: google.protobuf.ITimestamp;

    /** CafeSession refresh. */
    public refresh: string;

    /** CafeSession rexp. */
    public rexp: google.protobuf.ITimestamp;

    /** CafeSession subject. */
    public subject: string;

    /** CafeSession type. */
    public type: string;

    /** CafeSession cafe. */
    public cafe: ICafe;

    /**
     * Creates a new CafeSession instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeSession instance
     */
    public static create(properties?: ICafeSession): CafeSession;

    /**
     * Encodes the specified CafeSession message. Does not implicitly {@link CafeSession.verify|verify} messages.
     * @param message CafeSession message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeSession, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeSession message, length delimited. Does not implicitly {@link CafeSession.verify|verify} messages.
     * @param message CafeSession message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeSession, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeSession message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeSession
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeSession;

    /**
     * Decodes a CafeSession message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeSession
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeSession;

    /**
     * Verifies a CafeSession message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeSession message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeSession
     */
    public static fromObject(object: { [k: string]: any }): CafeSession;

    /**
     * Creates a plain object from a CafeSession message. Also converts values to other types if specified.
     * @param message CafeSession
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeSession, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeSession to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeSessionList. */
export interface ICafeSessionList {

    /** CafeSessionList items */
    items: ICafeSession[];
}

/** Represents a CafeSessionList. */
export class CafeSessionList implements ICafeSessionList {

    /**
     * Constructs a new CafeSessionList.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeSessionList);

    /** CafeSessionList items. */
    public items: ICafeSession[];

    /**
     * Creates a new CafeSessionList instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeSessionList instance
     */
    public static create(properties?: ICafeSessionList): CafeSessionList;

    /**
     * Encodes the specified CafeSessionList message. Does not implicitly {@link CafeSessionList.verify|verify} messages.
     * @param message CafeSessionList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeSessionList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeSessionList message, length delimited. Does not implicitly {@link CafeSessionList.verify|verify} messages.
     * @param message CafeSessionList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeSessionList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeSessionList message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeSessionList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeSessionList;

    /**
     * Decodes a CafeSessionList message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeSessionList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeSessionList;

    /**
     * Verifies a CafeSessionList message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeSessionList message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeSessionList
     */
    public static fromObject(object: { [k: string]: any }): CafeSessionList;

    /**
     * Creates a plain object from a CafeSessionList message. Also converts values to other types if specified.
     * @param message CafeSessionList
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeSessionList, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeSessionList to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeRequest. */
export interface ICafeRequest {

    /** CafeRequest id */
    id: string;

    /** CafeRequest peer */
    peer: string;

    /** CafeRequest target */
    target: string;

    /** CafeRequest cafe */
    cafe: ICafe;

    /** CafeRequest type */
    type: CafeRequest.Type;

    /** CafeRequest date */
    date: google.protobuf.ITimestamp;
}

/** Represents a CafeRequest. */
export class CafeRequest implements ICafeRequest {

    /**
     * Constructs a new CafeRequest.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeRequest);

    /** CafeRequest id. */
    public id: string;

    /** CafeRequest peer. */
    public peer: string;

    /** CafeRequest target. */
    public target: string;

    /** CafeRequest cafe. */
    public cafe: ICafe;

    /** CafeRequest type. */
    public type: CafeRequest.Type;

    /** CafeRequest date. */
    public date: google.protobuf.ITimestamp;

    /**
     * Creates a new CafeRequest instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeRequest instance
     */
    public static create(properties?: ICafeRequest): CafeRequest;

    /**
     * Encodes the specified CafeRequest message. Does not implicitly {@link CafeRequest.verify|verify} messages.
     * @param message CafeRequest message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeRequest, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeRequest message, length delimited. Does not implicitly {@link CafeRequest.verify|verify} messages.
     * @param message CafeRequest message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeRequest, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeRequest message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeRequest;

    /**
     * Decodes a CafeRequest message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeRequest;

    /**
     * Verifies a CafeRequest message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeRequest message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeRequest
     */
    public static fromObject(object: { [k: string]: any }): CafeRequest;

    /**
     * Creates a plain object from a CafeRequest message. Also converts values to other types if specified.
     * @param message CafeRequest
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeRequest to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

export namespace CafeRequest {

    /** Type enum. */
    enum Type {
        STORE = 0,
        STORE_THREAD = 1,
        INBOX = 2
    }
}

/** Properties of a CafeMessage. */
export interface ICafeMessage {

    /** CafeMessage id */
    id: string;

    /** CafeMessage peer */
    peer: string;

    /** CafeMessage date */
    date: google.protobuf.ITimestamp;

    /** CafeMessage attempts */
    attempts: number;
}

/** Represents a CafeMessage. */
export class CafeMessage implements ICafeMessage {

    /**
     * Constructs a new CafeMessage.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeMessage);

    /** CafeMessage id. */
    public id: string;

    /** CafeMessage peer. */
    public peer: string;

    /** CafeMessage date. */
    public date: google.protobuf.ITimestamp;

    /** CafeMessage attempts. */
    public attempts: number;

    /**
     * Creates a new CafeMessage instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeMessage instance
     */
    public static create(properties?: ICafeMessage): CafeMessage;

    /**
     * Encodes the specified CafeMessage message. Does not implicitly {@link CafeMessage.verify|verify} messages.
     * @param message CafeMessage message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeMessage, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeMessage message, length delimited. Does not implicitly {@link CafeMessage.verify|verify} messages.
     * @param message CafeMessage message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeMessage, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeMessage message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeMessage;

    /**
     * Decodes a CafeMessage message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeMessage;

    /**
     * Verifies a CafeMessage message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeMessage message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeMessage
     */
    public static fromObject(object: { [k: string]: any }): CafeMessage;

    /**
     * Creates a plain object from a CafeMessage message. Also converts values to other types if specified.
     * @param message CafeMessage
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeMessage, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeMessage to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeClientNonce. */
export interface ICafeClientNonce {

    /** CafeClientNonce value */
    value: string;

    /** CafeClientNonce address */
    address: string;

    /** CafeClientNonce date */
    date: google.protobuf.ITimestamp;
}

/** Represents a CafeClientNonce. */
export class CafeClientNonce implements ICafeClientNonce {

    /**
     * Constructs a new CafeClientNonce.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeClientNonce);

    /** CafeClientNonce value. */
    public value: string;

    /** CafeClientNonce address. */
    public address: string;

    /** CafeClientNonce date. */
    public date: google.protobuf.ITimestamp;

    /**
     * Creates a new CafeClientNonce instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeClientNonce instance
     */
    public static create(properties?: ICafeClientNonce): CafeClientNonce;

    /**
     * Encodes the specified CafeClientNonce message. Does not implicitly {@link CafeClientNonce.verify|verify} messages.
     * @param message CafeClientNonce message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeClientNonce, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeClientNonce message, length delimited. Does not implicitly {@link CafeClientNonce.verify|verify} messages.
     * @param message CafeClientNonce message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeClientNonce, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeClientNonce message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeClientNonce
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeClientNonce;

    /**
     * Decodes a CafeClientNonce message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeClientNonce
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeClientNonce;

    /**
     * Verifies a CafeClientNonce message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeClientNonce message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeClientNonce
     */
    public static fromObject(object: { [k: string]: any }): CafeClientNonce;

    /**
     * Creates a plain object from a CafeClientNonce message. Also converts values to other types if specified.
     * @param message CafeClientNonce
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeClientNonce, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeClientNonce to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeClient. */
export interface ICafeClient {

    /** CafeClient id */
    id: string;

    /** CafeClient address */
    address: string;

    /** CafeClient created */
    created: google.protobuf.ITimestamp;

    /** CafeClient seen */
    seen: google.protobuf.ITimestamp;

    /** CafeClient token */
    token: string;
}

/** Represents a CafeClient. */
export class CafeClient implements ICafeClient {

    /**
     * Constructs a new CafeClient.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeClient);

    /** CafeClient id. */
    public id: string;

    /** CafeClient address. */
    public address: string;

    /** CafeClient created. */
    public created: google.protobuf.ITimestamp;

    /** CafeClient seen. */
    public seen: google.protobuf.ITimestamp;

    /** CafeClient token. */
    public token: string;

    /**
     * Creates a new CafeClient instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeClient instance
     */
    public static create(properties?: ICafeClient): CafeClient;

    /**
     * Encodes the specified CafeClient message. Does not implicitly {@link CafeClient.verify|verify} messages.
     * @param message CafeClient message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeClient, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeClient message, length delimited. Does not implicitly {@link CafeClient.verify|verify} messages.
     * @param message CafeClient message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeClient, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeClient message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeClient
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeClient;

    /**
     * Decodes a CafeClient message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeClient
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeClient;

    /**
     * Verifies a CafeClient message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeClient message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeClient
     */
    public static fromObject(object: { [k: string]: any }): CafeClient;

    /**
     * Creates a plain object from a CafeClient message. Also converts values to other types if specified.
     * @param message CafeClient
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeClient, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeClient to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeClientList. */
export interface ICafeClientList {

    /** CafeClientList items */
    items: ICafeClient[];
}

/** Represents a CafeClientList. */
export class CafeClientList implements ICafeClientList {

    /**
     * Constructs a new CafeClientList.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeClientList);

    /** CafeClientList items. */
    public items: ICafeClient[];

    /**
     * Creates a new CafeClientList instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeClientList instance
     */
    public static create(properties?: ICafeClientList): CafeClientList;

    /**
     * Encodes the specified CafeClientList message. Does not implicitly {@link CafeClientList.verify|verify} messages.
     * @param message CafeClientList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeClientList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeClientList message, length delimited. Does not implicitly {@link CafeClientList.verify|verify} messages.
     * @param message CafeClientList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeClientList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeClientList message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeClientList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeClientList;

    /**
     * Decodes a CafeClientList message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeClientList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeClientList;

    /**
     * Verifies a CafeClientList message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeClientList message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeClientList
     */
    public static fromObject(object: { [k: string]: any }): CafeClientList;

    /**
     * Creates a plain object from a CafeClientList message. Also converts values to other types if specified.
     * @param message CafeClientList
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeClientList, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeClientList to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeToken. */
export interface ICafeToken {

    /** CafeToken id */
    id: string;

    /** CafeToken value */
    value: Uint8Array;

    /** CafeToken date */
    date: google.protobuf.ITimestamp;
}

/** Represents a CafeToken. */
export class CafeToken implements ICafeToken {

    /**
     * Constructs a new CafeToken.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeToken);

    /** CafeToken id. */
    public id: string;

    /** CafeToken value. */
    public value: Uint8Array;

    /** CafeToken date. */
    public date: google.protobuf.ITimestamp;

    /**
     * Creates a new CafeToken instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeToken instance
     */
    public static create(properties?: ICafeToken): CafeToken;

    /**
     * Encodes the specified CafeToken message. Does not implicitly {@link CafeToken.verify|verify} messages.
     * @param message CafeToken message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeToken, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeToken message, length delimited. Does not implicitly {@link CafeToken.verify|verify} messages.
     * @param message CafeToken message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeToken, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeToken message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeToken
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeToken;

    /**
     * Decodes a CafeToken message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeToken
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeToken;

    /**
     * Verifies a CafeToken message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeToken message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeToken
     */
    public static fromObject(object: { [k: string]: any }): CafeToken;

    /**
     * Creates a plain object from a CafeToken message. Also converts values to other types if specified.
     * @param message CafeToken
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeToken, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeToken to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeClientThread. */
export interface ICafeClientThread {

    /** CafeClientThread id */
    id: string;

    /** CafeClientThread client */
    client: string;

    /** CafeClientThread ciphertext */
    ciphertext: Uint8Array;
}

/** Represents a CafeClientThread. */
export class CafeClientThread implements ICafeClientThread {

    /**
     * Constructs a new CafeClientThread.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeClientThread);

    /** CafeClientThread id. */
    public id: string;

    /** CafeClientThread client. */
    public client: string;

    /** CafeClientThread ciphertext. */
    public ciphertext: Uint8Array;

    /**
     * Creates a new CafeClientThread instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeClientThread instance
     */
    public static create(properties?: ICafeClientThread): CafeClientThread;

    /**
     * Encodes the specified CafeClientThread message. Does not implicitly {@link CafeClientThread.verify|verify} messages.
     * @param message CafeClientThread message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeClientThread, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeClientThread message, length delimited. Does not implicitly {@link CafeClientThread.verify|verify} messages.
     * @param message CafeClientThread message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeClientThread, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeClientThread message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeClientThread
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeClientThread;

    /**
     * Decodes a CafeClientThread message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeClientThread
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeClientThread;

    /**
     * Verifies a CafeClientThread message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeClientThread message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeClientThread
     */
    public static fromObject(object: { [k: string]: any }): CafeClientThread;

    /**
     * Creates a plain object from a CafeClientThread message. Also converts values to other types if specified.
     * @param message CafeClientThread
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeClientThread, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeClientThread to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CafeClientMessage. */
export interface ICafeClientMessage {

    /** CafeClientMessage id */
    id: string;

    /** CafeClientMessage peer */
    peer: string;

    /** CafeClientMessage client */
    client: string;

    /** CafeClientMessage date */
    date: google.protobuf.ITimestamp;
}

/** Represents a CafeClientMessage. */
export class CafeClientMessage implements ICafeClientMessage {

    /**
     * Constructs a new CafeClientMessage.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICafeClientMessage);

    /** CafeClientMessage id. */
    public id: string;

    /** CafeClientMessage peer. */
    public peer: string;

    /** CafeClientMessage client. */
    public client: string;

    /** CafeClientMessage date. */
    public date: google.protobuf.ITimestamp;

    /**
     * Creates a new CafeClientMessage instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CafeClientMessage instance
     */
    public static create(properties?: ICafeClientMessage): CafeClientMessage;

    /**
     * Encodes the specified CafeClientMessage message. Does not implicitly {@link CafeClientMessage.verify|verify} messages.
     * @param message CafeClientMessage message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICafeClientMessage, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CafeClientMessage message, length delimited. Does not implicitly {@link CafeClientMessage.verify|verify} messages.
     * @param message CafeClientMessage message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICafeClientMessage, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CafeClientMessage message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CafeClientMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CafeClientMessage;

    /**
     * Decodes a CafeClientMessage message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CafeClientMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CafeClientMessage;

    /**
     * Verifies a CafeClientMessage message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CafeClientMessage message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CafeClientMessage
     */
    public static fromObject(object: { [k: string]: any }): CafeClientMessage;

    /**
     * Creates a plain object from a CafeClientMessage message. Also converts values to other types if specified.
     * @param message CafeClientMessage
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CafeClientMessage, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CafeClientMessage to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a Message. */
export interface IMessage {

    /** Message type */
    type: Message.Type;

    /** Message payload */
    payload: google.protobuf.IAny;

    /** Message requestId */
    requestId: number;

    /** Message isResponse */
    isResponse: boolean;
}

/** Represents a Message. */
export class Message implements IMessage {

    /**
     * Constructs a new Message.
     * @param [properties] Properties to set
     */
    constructor(properties?: IMessage);

    /** Message type. */
    public type: Message.Type;

    /** Message payload. */
    public payload: google.protobuf.IAny;

    /** Message requestId. */
    public requestId: number;

    /** Message isResponse. */
    public isResponse: boolean;

    /**
     * Creates a new Message instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Message instance
     */
    public static create(properties?: IMessage): Message;

    /**
     * Encodes the specified Message message. Does not implicitly {@link Message.verify|verify} messages.
     * @param message Message message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IMessage, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Message message, length delimited. Does not implicitly {@link Message.verify|verify} messages.
     * @param message Message message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IMessage, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Message message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Message
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Message;

    /**
     * Decodes a Message message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Message
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Message;

    /**
     * Verifies a Message message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Message message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Message
     */
    public static fromObject(object: { [k: string]: any }): Message;

    /**
     * Creates a plain object from a Message message. Also converts values to other types if specified.
     * @param message Message
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Message, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Message to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

export namespace Message {

    /** Type enum. */
    enum Type {
        PING = 0,
        PONG = 1,
        THREAD_ENVELOPE = 10,
        CAFE_CHALLENGE = 50,
        CAFE_NONCE = 51,
        CAFE_REGISTRATION = 52,
        CAFE_SESSION = 53,
        CAFE_REFRESH_SESSION = 54,
        CAFE_STORE = 55,
        CAFE_OBJECT = 56,
        CAFE_OBJECT_LIST = 57,
        CAFE_STORE_THREAD = 58,
        CAFE_STORED = 59,
        CAFE_DELIVER_MESSAGE = 60,
        CAFE_CHECK_MESSAGES = 61,
        CAFE_MESSAGES = 62,
        CAFE_DELETE_MESSAGES = 63,
        CAFE_DELETE_MESSAGES_ACK = 64,
        CAFE_YOU_HAVE_MAIL = 65,
        CAFE_PUBLISH_CONTACT = 66,
        CAFE_PUBLISH_CONTACT_ACK = 67,
        CAFE_QUERY = 70,
        CAFE_QUERY_RES = 71,
        CAFE_PUBSUB_QUERY = 102,
        CAFE_PUBSUB_QUERY_RES = 103,
        ERROR = 500,
        CAFE_CONTACT_QUERY = 68,
        CAFE_CONTACT_QUERY_RES = 69,
        CAFE_PUBSUB_CONTACT_QUERY = 100,
        CAFE_PUBSUB_CONTACT_QUERY_RES = 101
    }
}

/** Properties of an Envelope. */
export interface IEnvelope {

    /** Envelope message */
    message: IMessage;

    /** Envelope sig */
    sig: Uint8Array;
}

/** Represents an Envelope. */
export class Envelope implements IEnvelope {

    /**
     * Constructs a new Envelope.
     * @param [properties] Properties to set
     */
    constructor(properties?: IEnvelope);

    /** Envelope message. */
    public message: IMessage;

    /** Envelope sig. */
    public sig: Uint8Array;

    /**
     * Creates a new Envelope instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Envelope instance
     */
    public static create(properties?: IEnvelope): Envelope;

    /**
     * Encodes the specified Envelope message. Does not implicitly {@link Envelope.verify|verify} messages.
     * @param message Envelope message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IEnvelope, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Envelope message, length delimited. Does not implicitly {@link Envelope.verify|verify} messages.
     * @param message Envelope message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IEnvelope, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes an Envelope message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Envelope
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Envelope;

    /**
     * Decodes an Envelope message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Envelope
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Envelope;

    /**
     * Verifies an Envelope message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates an Envelope message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Envelope
     */
    public static fromObject(object: { [k: string]: any }): Envelope;

    /**
     * Creates a plain object from an Envelope message. Also converts values to other types if specified.
     * @param message Envelope
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Envelope, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Envelope to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of an Error. */
export interface IError {

    /** Error code */
    code: number;

    /** Error message */
    message: string;
}

/** Represents an Error. */
export class Error implements IError {

    /**
     * Constructs a new Error.
     * @param [properties] Properties to set
     */
    constructor(properties?: IError);

    /** Error code. */
    public code: number;

    /** Error message. */
    public message: string;

    /**
     * Creates a new Error instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Error instance
     */
    public static create(properties?: IError): Error;

    /**
     * Encodes the specified Error message. Does not implicitly {@link Error.verify|verify} messages.
     * @param message Error message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IError, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Error message, length delimited. Does not implicitly {@link Error.verify|verify} messages.
     * @param message Error message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IError, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes an Error message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Error
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Error;

    /**
     * Decodes an Error message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Error
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Error;

    /**
     * Verifies an Error message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates an Error message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Error
     */
    public static fromObject(object: { [k: string]: any }): Error;

    /**
     * Creates a plain object from an Error message. Also converts values to other types if specified.
     * @param message Error
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Error, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Error to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Namespace google. */
export namespace google {

    /** Namespace protobuf. */
    namespace protobuf {

        /** Properties of a Timestamp. */
        interface ITimestamp {

            /** Timestamp seconds */
            seconds: (number|Long);

            /** Timestamp nanos */
            nanos: number;
        }

        /** Represents a Timestamp. */
        class Timestamp implements ITimestamp {

            /**
             * Constructs a new Timestamp.
             * @param [properties] Properties to set
             */
            constructor(properties?: google.protobuf.ITimestamp);

            /** Timestamp seconds. */
            public seconds: (number|Long);

            /** Timestamp nanos. */
            public nanos: number;

            /**
             * Creates a new Timestamp instance using the specified properties.
             * @param [properties] Properties to set
             * @returns Timestamp instance
             */
            public static create(properties?: google.protobuf.ITimestamp): google.protobuf.Timestamp;

            /**
             * Encodes the specified Timestamp message. Does not implicitly {@link google.protobuf.Timestamp.verify|verify} messages.
             * @param message Timestamp message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encode(message: google.protobuf.ITimestamp, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Encodes the specified Timestamp message, length delimited. Does not implicitly {@link google.protobuf.Timestamp.verify|verify} messages.
             * @param message Timestamp message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encodeDelimited(message: google.protobuf.ITimestamp, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Decodes a Timestamp message from the specified reader or buffer.
             * @param reader Reader or buffer to decode from
             * @param [length] Message length if known beforehand
             * @returns Timestamp
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): google.protobuf.Timestamp;

            /**
             * Decodes a Timestamp message from the specified reader or buffer, length delimited.
             * @param reader Reader or buffer to decode from
             * @returns Timestamp
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): google.protobuf.Timestamp;

            /**
             * Verifies a Timestamp message.
             * @param message Plain object to verify
             * @returns `null` if valid, otherwise the reason why it is not
             */
            public static verify(message: { [k: string]: any }): (string|null);

            /**
             * Creates a Timestamp message from a plain object. Also converts values to their respective internal types.
             * @param object Plain object
             * @returns Timestamp
             */
            public static fromObject(object: { [k: string]: any }): google.protobuf.Timestamp;

            /**
             * Creates a plain object from a Timestamp message. Also converts values to other types if specified.
             * @param message Timestamp
             * @param [options] Conversion options
             * @returns Plain object
             */
            public static toObject(message: google.protobuf.Timestamp, options?: $protobuf.IConversionOptions): { [k: string]: any };

            /**
             * Converts this Timestamp to JSON.
             * @returns JSON object
             */
            public toJSON(): { [k: string]: any };
        }

        /** Properties of a Struct. */
        interface IStruct {

            /** Struct fields */
            fields: { [k: string]: google.protobuf.IValue };
        }

        /** Represents a Struct. */
        class Struct implements IStruct {

            /**
             * Constructs a new Struct.
             * @param [properties] Properties to set
             */
            constructor(properties?: google.protobuf.IStruct);

            /** Struct fields. */
            public fields: { [k: string]: google.protobuf.IValue };

            /**
             * Creates a new Struct instance using the specified properties.
             * @param [properties] Properties to set
             * @returns Struct instance
             */
            public static create(properties?: google.protobuf.IStruct): google.protobuf.Struct;

            /**
             * Encodes the specified Struct message. Does not implicitly {@link google.protobuf.Struct.verify|verify} messages.
             * @param message Struct message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encode(message: google.protobuf.IStruct, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Encodes the specified Struct message, length delimited. Does not implicitly {@link google.protobuf.Struct.verify|verify} messages.
             * @param message Struct message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encodeDelimited(message: google.protobuf.IStruct, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Decodes a Struct message from the specified reader or buffer.
             * @param reader Reader or buffer to decode from
             * @param [length] Message length if known beforehand
             * @returns Struct
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): google.protobuf.Struct;

            /**
             * Decodes a Struct message from the specified reader or buffer, length delimited.
             * @param reader Reader or buffer to decode from
             * @returns Struct
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): google.protobuf.Struct;

            /**
             * Verifies a Struct message.
             * @param message Plain object to verify
             * @returns `null` if valid, otherwise the reason why it is not
             */
            public static verify(message: { [k: string]: any }): (string|null);

            /**
             * Creates a Struct message from a plain object. Also converts values to their respective internal types.
             * @param object Plain object
             * @returns Struct
             */
            public static fromObject(object: { [k: string]: any }): google.protobuf.Struct;

            /**
             * Creates a plain object from a Struct message. Also converts values to other types if specified.
             * @param message Struct
             * @param [options] Conversion options
             * @returns Plain object
             */
            public static toObject(message: google.protobuf.Struct, options?: $protobuf.IConversionOptions): { [k: string]: any };

            /**
             * Converts this Struct to JSON.
             * @returns JSON object
             */
            public toJSON(): { [k: string]: any };
        }

        /** Properties of a Value. */
        interface IValue {

            /** Value nullValue */
            nullValue: google.protobuf.NullValue;

            /** Value numberValue */
            numberValue: number;

            /** Value stringValue */
            stringValue: string;

            /** Value boolValue */
            boolValue: boolean;

            /** Value structValue */
            structValue: google.protobuf.IStruct;

            /** Value listValue */
            listValue: google.protobuf.IListValue;
        }

        /** Represents a Value. */
        class Value implements IValue {

            /**
             * Constructs a new Value.
             * @param [properties] Properties to set
             */
            constructor(properties?: google.protobuf.IValue);

            /** Value nullValue. */
            public nullValue: google.protobuf.NullValue;

            /** Value numberValue. */
            public numberValue: number;

            /** Value stringValue. */
            public stringValue: string;

            /** Value boolValue. */
            public boolValue: boolean;

            /** Value structValue. */
            public structValue: google.protobuf.IStruct;

            /** Value listValue. */
            public listValue: google.protobuf.IListValue;

            /** Value kind. */
            public kind?: ("nullValue"|"numberValue"|"stringValue"|"boolValue"|"structValue"|"listValue");

            /**
             * Creates a new Value instance using the specified properties.
             * @param [properties] Properties to set
             * @returns Value instance
             */
            public static create(properties?: google.protobuf.IValue): google.protobuf.Value;

            /**
             * Encodes the specified Value message. Does not implicitly {@link google.protobuf.Value.verify|verify} messages.
             * @param message Value message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encode(message: google.protobuf.IValue, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Encodes the specified Value message, length delimited. Does not implicitly {@link google.protobuf.Value.verify|verify} messages.
             * @param message Value message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encodeDelimited(message: google.protobuf.IValue, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Decodes a Value message from the specified reader or buffer.
             * @param reader Reader or buffer to decode from
             * @param [length] Message length if known beforehand
             * @returns Value
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): google.protobuf.Value;

            /**
             * Decodes a Value message from the specified reader or buffer, length delimited.
             * @param reader Reader or buffer to decode from
             * @returns Value
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): google.protobuf.Value;

            /**
             * Verifies a Value message.
             * @param message Plain object to verify
             * @returns `null` if valid, otherwise the reason why it is not
             */
            public static verify(message: { [k: string]: any }): (string|null);

            /**
             * Creates a Value message from a plain object. Also converts values to their respective internal types.
             * @param object Plain object
             * @returns Value
             */
            public static fromObject(object: { [k: string]: any }): google.protobuf.Value;

            /**
             * Creates a plain object from a Value message. Also converts values to other types if specified.
             * @param message Value
             * @param [options] Conversion options
             * @returns Plain object
             */
            public static toObject(message: google.protobuf.Value, options?: $protobuf.IConversionOptions): { [k: string]: any };

            /**
             * Converts this Value to JSON.
             * @returns JSON object
             */
            public toJSON(): { [k: string]: any };
        }

        /** NullValue enum. */
        enum NullValue {
            NULL_VALUE = 0
        }

        /** Properties of a ListValue. */
        interface IListValue {

            /** ListValue values */
            values: google.protobuf.IValue[];
        }

        /** Represents a ListValue. */
        class ListValue implements IListValue {

            /**
             * Constructs a new ListValue.
             * @param [properties] Properties to set
             */
            constructor(properties?: google.protobuf.IListValue);

            /** ListValue values. */
            public values: google.protobuf.IValue[];

            /**
             * Creates a new ListValue instance using the specified properties.
             * @param [properties] Properties to set
             * @returns ListValue instance
             */
            public static create(properties?: google.protobuf.IListValue): google.protobuf.ListValue;

            /**
             * Encodes the specified ListValue message. Does not implicitly {@link google.protobuf.ListValue.verify|verify} messages.
             * @param message ListValue message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encode(message: google.protobuf.IListValue, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Encodes the specified ListValue message, length delimited. Does not implicitly {@link google.protobuf.ListValue.verify|verify} messages.
             * @param message ListValue message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encodeDelimited(message: google.protobuf.IListValue, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Decodes a ListValue message from the specified reader or buffer.
             * @param reader Reader or buffer to decode from
             * @param [length] Message length if known beforehand
             * @returns ListValue
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): google.protobuf.ListValue;

            /**
             * Decodes a ListValue message from the specified reader or buffer, length delimited.
             * @param reader Reader or buffer to decode from
             * @returns ListValue
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): google.protobuf.ListValue;

            /**
             * Verifies a ListValue message.
             * @param message Plain object to verify
             * @returns `null` if valid, otherwise the reason why it is not
             */
            public static verify(message: { [k: string]: any }): (string|null);

            /**
             * Creates a ListValue message from a plain object. Also converts values to their respective internal types.
             * @param object Plain object
             * @returns ListValue
             */
            public static fromObject(object: { [k: string]: any }): google.protobuf.ListValue;

            /**
             * Creates a plain object from a ListValue message. Also converts values to other types if specified.
             * @param message ListValue
             * @param [options] Conversion options
             * @returns Plain object
             */
            public static toObject(message: google.protobuf.ListValue, options?: $protobuf.IConversionOptions): { [k: string]: any };

            /**
             * Converts this ListValue to JSON.
             * @returns JSON object
             */
            public toJSON(): { [k: string]: any };
        }

        /** Properties of an Any. */
        interface IAny {

            /** Any type_url */
            type_url: string;

            /** Any value */
            value: Uint8Array;
        }

        /** Represents an Any. */
        class Any implements IAny {

            /**
             * Constructs a new Any.
             * @param [properties] Properties to set
             */
            constructor(properties?: google.protobuf.IAny);

            /** Any type_url. */
            public type_url: string;

            /** Any value. */
            public value: Uint8Array;

            /**
             * Creates a new Any instance using the specified properties.
             * @param [properties] Properties to set
             * @returns Any instance
             */
            public static create(properties?: google.protobuf.IAny): google.protobuf.Any;

            /**
             * Encodes the specified Any message. Does not implicitly {@link google.protobuf.Any.verify|verify} messages.
             * @param message Any message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encode(message: google.protobuf.IAny, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Encodes the specified Any message, length delimited. Does not implicitly {@link google.protobuf.Any.verify|verify} messages.
             * @param message Any message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encodeDelimited(message: google.protobuf.IAny, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Decodes an Any message from the specified reader or buffer.
             * @param reader Reader or buffer to decode from
             * @param [length] Message length if known beforehand
             * @returns Any
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): google.protobuf.Any;

            /**
             * Decodes an Any message from the specified reader or buffer, length delimited.
             * @param reader Reader or buffer to decode from
             * @returns Any
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): google.protobuf.Any;

            /**
             * Verifies an Any message.
             * @param message Plain object to verify
             * @returns `null` if valid, otherwise the reason why it is not
             */
            public static verify(message: { [k: string]: any }): (string|null);

            /**
             * Creates an Any message from a plain object. Also converts values to their respective internal types.
             * @param object Plain object
             * @returns Any
             */
            public static fromObject(object: { [k: string]: any }): google.protobuf.Any;

            /**
             * Creates a plain object from an Any message. Also converts values to other types if specified.
             * @param message Any
             * @param [options] Conversion options
             * @returns Plain object
             */
            public static toObject(message: google.protobuf.Any, options?: $protobuf.IConversionOptions): { [k: string]: any };

            /**
             * Converts this Any to JSON.
             * @returns JSON object
             */
            public toJSON(): { [k: string]: any };
        }
    }
}

/** Properties of a MobilePreparedFiles. */
export interface IMobilePreparedFiles {

    /** MobilePreparedFiles dir */
    dir: IDirectory;

    /** MobilePreparedFiles pin */
    pin: { [k: string]: string };
}

/** Represents a MobilePreparedFiles. */
export class MobilePreparedFiles implements IMobilePreparedFiles {

    /**
     * Constructs a new MobilePreparedFiles.
     * @param [properties] Properties to set
     */
    constructor(properties?: IMobilePreparedFiles);

    /** MobilePreparedFiles dir. */
    public dir: IDirectory;

    /** MobilePreparedFiles pin. */
    public pin: { [k: string]: string };

    /**
     * Creates a new MobilePreparedFiles instance using the specified properties.
     * @param [properties] Properties to set
     * @returns MobilePreparedFiles instance
     */
    public static create(properties?: IMobilePreparedFiles): MobilePreparedFiles;

    /**
     * Encodes the specified MobilePreparedFiles message. Does not implicitly {@link MobilePreparedFiles.verify|verify} messages.
     * @param message MobilePreparedFiles message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IMobilePreparedFiles, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified MobilePreparedFiles message, length delimited. Does not implicitly {@link MobilePreparedFiles.verify|verify} messages.
     * @param message MobilePreparedFiles message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IMobilePreparedFiles, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a MobilePreparedFiles message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns MobilePreparedFiles
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): MobilePreparedFiles;

    /**
     * Decodes a MobilePreparedFiles message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns MobilePreparedFiles
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): MobilePreparedFiles;

    /**
     * Verifies a MobilePreparedFiles message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a MobilePreparedFiles message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns MobilePreparedFiles
     */
    public static fromObject(object: { [k: string]: any }): MobilePreparedFiles;

    /**
     * Creates a plain object from a MobilePreparedFiles message. Also converts values to other types if specified.
     * @param message MobilePreparedFiles
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: MobilePreparedFiles, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this MobilePreparedFiles to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a MobileFileData. */
export interface IMobileFileData {

    /** MobileFileData url */
    url: string;
}

/** Represents a MobileFileData. */
export class MobileFileData implements IMobileFileData {

    /**
     * Constructs a new MobileFileData.
     * @param [properties] Properties to set
     */
    constructor(properties?: IMobileFileData);

    /** MobileFileData url. */
    public url: string;

    /**
     * Creates a new MobileFileData instance using the specified properties.
     * @param [properties] Properties to set
     * @returns MobileFileData instance
     */
    public static create(properties?: IMobileFileData): MobileFileData;

    /**
     * Encodes the specified MobileFileData message. Does not implicitly {@link MobileFileData.verify|verify} messages.
     * @param message MobileFileData message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IMobileFileData, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified MobileFileData message, length delimited. Does not implicitly {@link MobileFileData.verify|verify} messages.
     * @param message MobileFileData message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IMobileFileData, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a MobileFileData message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns MobileFileData
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): MobileFileData;

    /**
     * Decodes a MobileFileData message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns MobileFileData
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): MobileFileData;

    /**
     * Verifies a MobileFileData message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a MobileFileData message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns MobileFileData
     */
    public static fromObject(object: { [k: string]: any }): MobileFileData;

    /**
     * Creates a plain object from a MobileFileData message. Also converts values to other types if specified.
     * @param message MobileFileData
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: MobileFileData, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this MobileFileData to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of an AddThreadConfig. */
export interface IAddThreadConfig {

    /** AddThreadConfig key */
    key: string;

    /** AddThreadConfig name */
    name: string;

    /** AddThreadConfig schema */
    schema: AddThreadConfig.ISchema;

    /** AddThreadConfig type */
    type: Thread.Type;

    /** AddThreadConfig sharing */
    sharing: Thread.Sharing;

    /** AddThreadConfig members */
    members: string[];
}

/** Represents an AddThreadConfig. */
export class AddThreadConfig implements IAddThreadConfig {

    /**
     * Constructs a new AddThreadConfig.
     * @param [properties] Properties to set
     */
    constructor(properties?: IAddThreadConfig);

    /** AddThreadConfig key. */
    public key: string;

    /** AddThreadConfig name. */
    public name: string;

    /** AddThreadConfig schema. */
    public schema: AddThreadConfig.ISchema;

    /** AddThreadConfig type. */
    public type: Thread.Type;

    /** AddThreadConfig sharing. */
    public sharing: Thread.Sharing;

    /** AddThreadConfig members. */
    public members: string[];

    /**
     * Creates a new AddThreadConfig instance using the specified properties.
     * @param [properties] Properties to set
     * @returns AddThreadConfig instance
     */
    public static create(properties?: IAddThreadConfig): AddThreadConfig;

    /**
     * Encodes the specified AddThreadConfig message. Does not implicitly {@link AddThreadConfig.verify|verify} messages.
     * @param message AddThreadConfig message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IAddThreadConfig, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified AddThreadConfig message, length delimited. Does not implicitly {@link AddThreadConfig.verify|verify} messages.
     * @param message AddThreadConfig message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IAddThreadConfig, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes an AddThreadConfig message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns AddThreadConfig
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): AddThreadConfig;

    /**
     * Decodes an AddThreadConfig message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns AddThreadConfig
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): AddThreadConfig;

    /**
     * Verifies an AddThreadConfig message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates an AddThreadConfig message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns AddThreadConfig
     */
    public static fromObject(object: { [k: string]: any }): AddThreadConfig;

    /**
     * Creates a plain object from an AddThreadConfig message. Also converts values to other types if specified.
     * @param message AddThreadConfig
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: AddThreadConfig, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this AddThreadConfig to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

export namespace AddThreadConfig {

    /** Properties of a Schema. */
    interface ISchema {

        /** Schema id */
        id: string;

        /** Schema json */
        json: string;

        /** Schema preset */
        preset: AddThreadConfig.Schema.Preset;
    }

    /** Represents a Schema. */
    class Schema implements ISchema {

        /**
         * Constructs a new Schema.
         * @param [properties] Properties to set
         */
        constructor(properties?: AddThreadConfig.ISchema);

        /** Schema id. */
        public id: string;

        /** Schema json. */
        public json: string;

        /** Schema preset. */
        public preset: AddThreadConfig.Schema.Preset;

        /**
         * Creates a new Schema instance using the specified properties.
         * @param [properties] Properties to set
         * @returns Schema instance
         */
        public static create(properties?: AddThreadConfig.ISchema): AddThreadConfig.Schema;

        /**
         * Encodes the specified Schema message. Does not implicitly {@link AddThreadConfig.Schema.verify|verify} messages.
         * @param message Schema message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: AddThreadConfig.ISchema, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified Schema message, length delimited. Does not implicitly {@link AddThreadConfig.Schema.verify|verify} messages.
         * @param message Schema message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: AddThreadConfig.ISchema, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a Schema message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns Schema
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): AddThreadConfig.Schema;

        /**
         * Decodes a Schema message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns Schema
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): AddThreadConfig.Schema;

        /**
         * Verifies a Schema message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a Schema message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns Schema
         */
        public static fromObject(object: { [k: string]: any }): AddThreadConfig.Schema;

        /**
         * Creates a plain object from a Schema message. Also converts values to other types if specified.
         * @param message Schema
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: AddThreadConfig.Schema, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this Schema to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };
    }

    namespace Schema {

        /** Preset enum. */
        enum Preset {
            NONE = 0,
            CAMERA_ROLL = 1,
            MEDIA = 2
        }
    }
}

/** Properties of a Step. */
export interface IStep {

    /** Step Name */
    Name: string;

    /** Step link */
    link: ILink;
}

/** Represents a Step. */
export class Step implements IStep {

    /**
     * Constructs a new Step.
     * @param [properties] Properties to set
     */
    constructor(properties?: IStep);

    /** Step Name. */
    public Name: string;

    /** Step link. */
    public link: ILink;

    /**
     * Creates a new Step instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Step instance
     */
    public static create(properties?: IStep): Step;

    /**
     * Encodes the specified Step message. Does not implicitly {@link Step.verify|verify} messages.
     * @param message Step message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IStep, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Step message, length delimited. Does not implicitly {@link Step.verify|verify} messages.
     * @param message Step message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IStep, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Step message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Step
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Step;

    /**
     * Decodes a Step message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Step
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Step;

    /**
     * Verifies a Step message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Step message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Step
     */
    public static fromObject(object: { [k: string]: any }): Step;

    /**
     * Creates a plain object from a Step message. Also converts values to other types if specified.
     * @param message Step
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Step, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Step to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a Directory. */
export interface IDirectory {

    /** Directory files */
    files: { [k: string]: IFileIndex };
}

/** Represents a Directory. */
export class Directory implements IDirectory {

    /**
     * Constructs a new Directory.
     * @param [properties] Properties to set
     */
    constructor(properties?: IDirectory);

    /** Directory files. */
    public files: { [k: string]: IFileIndex };

    /**
     * Creates a new Directory instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Directory instance
     */
    public static create(properties?: IDirectory): Directory;

    /**
     * Encodes the specified Directory message. Does not implicitly {@link Directory.verify|verify} messages.
     * @param message Directory message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IDirectory, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Directory message, length delimited. Does not implicitly {@link Directory.verify|verify} messages.
     * @param message Directory message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IDirectory, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Directory message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Directory
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Directory;

    /**
     * Decodes a Directory message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Directory
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Directory;

    /**
     * Verifies a Directory message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Directory message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Directory
     */
    public static fromObject(object: { [k: string]: any }): Directory;

    /**
     * Creates a plain object from a Directory message. Also converts values to other types if specified.
     * @param message Directory
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Directory, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Directory to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a DirectoryList. */
export interface IDirectoryList {

    /** DirectoryList items */
    items: IDirectory[];
}

/** Represents a DirectoryList. */
export class DirectoryList implements IDirectoryList {

    /**
     * Constructs a new DirectoryList.
     * @param [properties] Properties to set
     */
    constructor(properties?: IDirectoryList);

    /** DirectoryList items. */
    public items: IDirectory[];

    /**
     * Creates a new DirectoryList instance using the specified properties.
     * @param [properties] Properties to set
     * @returns DirectoryList instance
     */
    public static create(properties?: IDirectoryList): DirectoryList;

    /**
     * Encodes the specified DirectoryList message. Does not implicitly {@link DirectoryList.verify|verify} messages.
     * @param message DirectoryList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IDirectoryList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified DirectoryList message, length delimited. Does not implicitly {@link DirectoryList.verify|verify} messages.
     * @param message DirectoryList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IDirectoryList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a DirectoryList message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns DirectoryList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): DirectoryList;

    /**
     * Decodes a DirectoryList message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns DirectoryList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): DirectoryList;

    /**
     * Verifies a DirectoryList message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a DirectoryList message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns DirectoryList
     */
    public static fromObject(object: { [k: string]: any }): DirectoryList;

    /**
     * Creates a plain object from a DirectoryList message. Also converts values to other types if specified.
     * @param message DirectoryList
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: DirectoryList, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this DirectoryList to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a Keys. */
export interface IKeys {

    /** Keys files */
    files: { [k: string]: string };
}

/** Represents a Keys. */
export class Keys implements IKeys {

    /**
     * Constructs a new Keys.
     * @param [properties] Properties to set
     */
    constructor(properties?: IKeys);

    /** Keys files. */
    public files: { [k: string]: string };

    /**
     * Creates a new Keys instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Keys instance
     */
    public static create(properties?: IKeys): Keys;

    /**
     * Encodes the specified Keys message. Does not implicitly {@link Keys.verify|verify} messages.
     * @param message Keys message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IKeys, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Keys message, length delimited. Does not implicitly {@link Keys.verify|verify} messages.
     * @param message Keys message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IKeys, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Keys message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Keys
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Keys;

    /**
     * Decodes a Keys message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Keys
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Keys;

    /**
     * Verifies a Keys message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Keys message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Keys
     */
    public static fromObject(object: { [k: string]: any }): Keys;

    /**
     * Creates a plain object from a Keys message. Also converts values to other types if specified.
     * @param message Keys
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Keys, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Keys to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a NewInvite. */
export interface INewInvite {

    /** NewInvite id */
    id: string;

    /** NewInvite key */
    key: string;

    /** NewInvite inviter */
    inviter: string;
}

/** Represents a NewInvite. */
export class NewInvite implements INewInvite {

    /**
     * Constructs a new NewInvite.
     * @param [properties] Properties to set
     */
    constructor(properties?: INewInvite);

    /** NewInvite id. */
    public id: string;

    /** NewInvite key. */
    public key: string;

    /** NewInvite inviter. */
    public inviter: string;

    /**
     * Creates a new NewInvite instance using the specified properties.
     * @param [properties] Properties to set
     * @returns NewInvite instance
     */
    public static create(properties?: INewInvite): NewInvite;

    /**
     * Encodes the specified NewInvite message. Does not implicitly {@link NewInvite.verify|verify} messages.
     * @param message NewInvite message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: INewInvite, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified NewInvite message, length delimited. Does not implicitly {@link NewInvite.verify|verify} messages.
     * @param message NewInvite message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: INewInvite, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a NewInvite message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns NewInvite
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): NewInvite;

    /**
     * Decodes a NewInvite message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns NewInvite
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): NewInvite;

    /**
     * Verifies a NewInvite message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a NewInvite message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns NewInvite
     */
    public static fromObject(object: { [k: string]: any }): NewInvite;

    /**
     * Creates a plain object from a NewInvite message. Also converts values to other types if specified.
     * @param message NewInvite
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: NewInvite, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this NewInvite to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of an InviteView. */
export interface IInviteView {

    /** InviteView id */
    id: string;

    /** InviteView name */
    name: string;

    /** InviteView inviter */
    inviter: IUser;

    /** InviteView date */
    date: google.protobuf.ITimestamp;
}

/** Represents an InviteView. */
export class InviteView implements IInviteView {

    /**
     * Constructs a new InviteView.
     * @param [properties] Properties to set
     */
    constructor(properties?: IInviteView);

    /** InviteView id. */
    public id: string;

    /** InviteView name. */
    public name: string;

    /** InviteView inviter. */
    public inviter: IUser;

    /** InviteView date. */
    public date: google.protobuf.ITimestamp;

    /**
     * Creates a new InviteView instance using the specified properties.
     * @param [properties] Properties to set
     * @returns InviteView instance
     */
    public static create(properties?: IInviteView): InviteView;

    /**
     * Encodes the specified InviteView message. Does not implicitly {@link InviteView.verify|verify} messages.
     * @param message InviteView message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IInviteView, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified InviteView message, length delimited. Does not implicitly {@link InviteView.verify|verify} messages.
     * @param message InviteView message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IInviteView, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes an InviteView message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns InviteView
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): InviteView;

    /**
     * Decodes an InviteView message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns InviteView
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): InviteView;

    /**
     * Verifies an InviteView message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates an InviteView message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns InviteView
     */
    public static fromObject(object: { [k: string]: any }): InviteView;

    /**
     * Creates a plain object from an InviteView message. Also converts values to other types if specified.
     * @param message InviteView
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: InviteView, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this InviteView to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of an InviteViewList. */
export interface IInviteViewList {

    /** InviteViewList items */
    items: IInviteView[];
}

/** Represents an InviteViewList. */
export class InviteViewList implements IInviteViewList {

    /**
     * Constructs a new InviteViewList.
     * @param [properties] Properties to set
     */
    constructor(properties?: IInviteViewList);

    /** InviteViewList items. */
    public items: IInviteView[];

    /**
     * Creates a new InviteViewList instance using the specified properties.
     * @param [properties] Properties to set
     * @returns InviteViewList instance
     */
    public static create(properties?: IInviteViewList): InviteViewList;

    /**
     * Encodes the specified InviteViewList message. Does not implicitly {@link InviteViewList.verify|verify} messages.
     * @param message InviteViewList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IInviteViewList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified InviteViewList message, length delimited. Does not implicitly {@link InviteViewList.verify|verify} messages.
     * @param message InviteViewList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IInviteViewList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes an InviteViewList message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns InviteViewList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): InviteViewList;

    /**
     * Decodes an InviteViewList message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns InviteViewList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): InviteViewList;

    /**
     * Verifies an InviteViewList message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates an InviteViewList message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns InviteViewList
     */
    public static fromObject(object: { [k: string]: any }): InviteViewList;

    /**
     * Creates a plain object from an InviteViewList message. Also converts values to other types if specified.
     * @param message InviteViewList
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: InviteViewList, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this InviteViewList to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** FeedMode enum. */
export enum FeedMode {
    CHRONO = 0,
    ANNOTATED = 1,
    STACKS = 2
}

/** Properties of a FeedItem. */
export interface IFeedItem {

    /** FeedItem block */
    block: string;

    /** FeedItem thread */
    thread: string;

    /** FeedItem payload */
    payload: google.protobuf.IAny;
}

/** Represents a FeedItem. */
export class FeedItem implements IFeedItem {

    /**
     * Constructs a new FeedItem.
     * @param [properties] Properties to set
     */
    constructor(properties?: IFeedItem);

    /** FeedItem block. */
    public block: string;

    /** FeedItem thread. */
    public thread: string;

    /** FeedItem payload. */
    public payload: google.protobuf.IAny;

    /**
     * Creates a new FeedItem instance using the specified properties.
     * @param [properties] Properties to set
     * @returns FeedItem instance
     */
    public static create(properties?: IFeedItem): FeedItem;

    /**
     * Encodes the specified FeedItem message. Does not implicitly {@link FeedItem.verify|verify} messages.
     * @param message FeedItem message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IFeedItem, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified FeedItem message, length delimited. Does not implicitly {@link FeedItem.verify|verify} messages.
     * @param message FeedItem message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IFeedItem, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a FeedItem message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns FeedItem
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): FeedItem;

    /**
     * Decodes a FeedItem message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns FeedItem
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): FeedItem;

    /**
     * Verifies a FeedItem message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a FeedItem message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns FeedItem
     */
    public static fromObject(object: { [k: string]: any }): FeedItem;

    /**
     * Creates a plain object from a FeedItem message. Also converts values to other types if specified.
     * @param message FeedItem
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: FeedItem, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this FeedItem to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a FeedItemList. */
export interface IFeedItemList {

    /** FeedItemList items */
    items: IFeedItem[];

    /** FeedItemList count */
    count: number;

    /** FeedItemList next */
    next: string;
}

/** Represents a FeedItemList. */
export class FeedItemList implements IFeedItemList {

    /**
     * Constructs a new FeedItemList.
     * @param [properties] Properties to set
     */
    constructor(properties?: IFeedItemList);

    /** FeedItemList items. */
    public items: IFeedItem[];

    /** FeedItemList count. */
    public count: number;

    /** FeedItemList next. */
    public next: string;

    /**
     * Creates a new FeedItemList instance using the specified properties.
     * @param [properties] Properties to set
     * @returns FeedItemList instance
     */
    public static create(properties?: IFeedItemList): FeedItemList;

    /**
     * Encodes the specified FeedItemList message. Does not implicitly {@link FeedItemList.verify|verify} messages.
     * @param message FeedItemList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IFeedItemList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified FeedItemList message, length delimited. Does not implicitly {@link FeedItemList.verify|verify} messages.
     * @param message FeedItemList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IFeedItemList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a FeedItemList message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns FeedItemList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): FeedItemList;

    /**
     * Decodes a FeedItemList message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns FeedItemList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): FeedItemList;

    /**
     * Verifies a FeedItemList message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a FeedItemList message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns FeedItemList
     */
    public static fromObject(object: { [k: string]: any }): FeedItemList;

    /**
     * Creates a plain object from a FeedItemList message. Also converts values to other types if specified.
     * @param message FeedItemList
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: FeedItemList, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this FeedItemList to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a Merge. */
export interface IMerge {

    /** Merge block */
    block: string;

    /** Merge date */
    date: google.protobuf.ITimestamp;

    /** Merge targets */
    targets: IFeedItem[];
}

/** Represents a Merge. */
export class Merge implements IMerge {

    /**
     * Constructs a new Merge.
     * @param [properties] Properties to set
     */
    constructor(properties?: IMerge);

    /** Merge block. */
    public block: string;

    /** Merge date. */
    public date: google.protobuf.ITimestamp;

    /** Merge targets. */
    public targets: IFeedItem[];

    /**
     * Creates a new Merge instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Merge instance
     */
    public static create(properties?: IMerge): Merge;

    /**
     * Encodes the specified Merge message. Does not implicitly {@link Merge.verify|verify} messages.
     * @param message Merge message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IMerge, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Merge message, length delimited. Does not implicitly {@link Merge.verify|verify} messages.
     * @param message Merge message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IMerge, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Merge message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Merge
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Merge;

    /**
     * Decodes a Merge message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Merge
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Merge;

    /**
     * Verifies a Merge message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Merge message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Merge
     */
    public static fromObject(object: { [k: string]: any }): Merge;

    /**
     * Creates a plain object from a Merge message. Also converts values to other types if specified.
     * @param message Merge
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Merge, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Merge to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of an Ignore. */
export interface IIgnore {

    /** Ignore block */
    block: string;

    /** Ignore date */
    date: google.protobuf.ITimestamp;

    /** Ignore user */
    user: IUser;

    /** Ignore target */
    target: IFeedItem;
}

/** Represents an Ignore. */
export class Ignore implements IIgnore {

    /**
     * Constructs a new Ignore.
     * @param [properties] Properties to set
     */
    constructor(properties?: IIgnore);

    /** Ignore block. */
    public block: string;

    /** Ignore date. */
    public date: google.protobuf.ITimestamp;

    /** Ignore user. */
    public user: IUser;

    /** Ignore target. */
    public target: IFeedItem;

    /**
     * Creates a new Ignore instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Ignore instance
     */
    public static create(properties?: IIgnore): Ignore;

    /**
     * Encodes the specified Ignore message. Does not implicitly {@link Ignore.verify|verify} messages.
     * @param message Ignore message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IIgnore, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Ignore message, length delimited. Does not implicitly {@link Ignore.verify|verify} messages.
     * @param message Ignore message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IIgnore, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes an Ignore message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Ignore
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Ignore;

    /**
     * Decodes an Ignore message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Ignore
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Ignore;

    /**
     * Verifies an Ignore message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates an Ignore message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Ignore
     */
    public static fromObject(object: { [k: string]: any }): Ignore;

    /**
     * Creates a plain object from an Ignore message. Also converts values to other types if specified.
     * @param message Ignore
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Ignore, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Ignore to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a Flag. */
export interface IFlag {

    /** Flag block */
    block: string;

    /** Flag date */
    date: google.protobuf.ITimestamp;

    /** Flag user */
    user: IUser;

    /** Flag target */
    target: IFeedItem;
}

/** Represents a Flag. */
export class Flag implements IFlag {

    /**
     * Constructs a new Flag.
     * @param [properties] Properties to set
     */
    constructor(properties?: IFlag);

    /** Flag block. */
    public block: string;

    /** Flag date. */
    public date: google.protobuf.ITimestamp;

    /** Flag user. */
    public user: IUser;

    /** Flag target. */
    public target: IFeedItem;

    /**
     * Creates a new Flag instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Flag instance
     */
    public static create(properties?: IFlag): Flag;

    /**
     * Encodes the specified Flag message. Does not implicitly {@link Flag.verify|verify} messages.
     * @param message Flag message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IFlag, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Flag message, length delimited. Does not implicitly {@link Flag.verify|verify} messages.
     * @param message Flag message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IFlag, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Flag message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Flag
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Flag;

    /**
     * Decodes a Flag message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Flag
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Flag;

    /**
     * Verifies a Flag message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Flag message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Flag
     */
    public static fromObject(object: { [k: string]: any }): Flag;

    /**
     * Creates a plain object from a Flag message. Also converts values to other types if specified.
     * @param message Flag
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Flag, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Flag to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of an Announce. */
export interface IAnnounce {

    /** Announce block */
    block: string;

    /** Announce date */
    date: google.protobuf.ITimestamp;

    /** Announce user */
    user: IUser;

    /** Announce target */
    target: IFeedItem;
}

/** Represents an Announce. */
export class Announce implements IAnnounce {

    /**
     * Constructs a new Announce.
     * @param [properties] Properties to set
     */
    constructor(properties?: IAnnounce);

    /** Announce block. */
    public block: string;

    /** Announce date. */
    public date: google.protobuf.ITimestamp;

    /** Announce user. */
    public user: IUser;

    /** Announce target. */
    public target: IFeedItem;

    /**
     * Creates a new Announce instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Announce instance
     */
    public static create(properties?: IAnnounce): Announce;

    /**
     * Encodes the specified Announce message. Does not implicitly {@link Announce.verify|verify} messages.
     * @param message Announce message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IAnnounce, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Announce message, length delimited. Does not implicitly {@link Announce.verify|verify} messages.
     * @param message Announce message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IAnnounce, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes an Announce message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Announce
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Announce;

    /**
     * Decodes an Announce message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Announce
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Announce;

    /**
     * Verifies an Announce message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates an Announce message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Announce
     */
    public static fromObject(object: { [k: string]: any }): Announce;

    /**
     * Creates a plain object from an Announce message. Also converts values to other types if specified.
     * @param message Announce
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Announce, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Announce to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a Join. */
export interface IJoin {

    /** Join block */
    block: string;

    /** Join date */
    date: google.protobuf.ITimestamp;

    /** Join user */
    user: IUser;

    /** Join likes */
    likes: ILike[];
}

/** Represents a Join. */
export class Join implements IJoin {

    /**
     * Constructs a new Join.
     * @param [properties] Properties to set
     */
    constructor(properties?: IJoin);

    /** Join block. */
    public block: string;

    /** Join date. */
    public date: google.protobuf.ITimestamp;

    /** Join user. */
    public user: IUser;

    /** Join likes. */
    public likes: ILike[];

    /**
     * Creates a new Join instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Join instance
     */
    public static create(properties?: IJoin): Join;

    /**
     * Encodes the specified Join message. Does not implicitly {@link Join.verify|verify} messages.
     * @param message Join message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IJoin, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Join message, length delimited. Does not implicitly {@link Join.verify|verify} messages.
     * @param message Join message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IJoin, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Join message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Join
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Join;

    /**
     * Decodes a Join message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Join
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Join;

    /**
     * Verifies a Join message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Join message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Join
     */
    public static fromObject(object: { [k: string]: any }): Join;

    /**
     * Creates a plain object from a Join message. Also converts values to other types if specified.
     * @param message Join
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Join, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Join to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a Leave. */
export interface ILeave {

    /** Leave block */
    block: string;

    /** Leave date */
    date: google.protobuf.ITimestamp;

    /** Leave user */
    user: IUser;

    /** Leave likes */
    likes: ILike[];
}

/** Represents a Leave. */
export class Leave implements ILeave {

    /**
     * Constructs a new Leave.
     * @param [properties] Properties to set
     */
    constructor(properties?: ILeave);

    /** Leave block. */
    public block: string;

    /** Leave date. */
    public date: google.protobuf.ITimestamp;

    /** Leave user. */
    public user: IUser;

    /** Leave likes. */
    public likes: ILike[];

    /**
     * Creates a new Leave instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Leave instance
     */
    public static create(properties?: ILeave): Leave;

    /**
     * Encodes the specified Leave message. Does not implicitly {@link Leave.verify|verify} messages.
     * @param message Leave message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ILeave, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Leave message, length delimited. Does not implicitly {@link Leave.verify|verify} messages.
     * @param message Leave message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ILeave, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Leave message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Leave
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Leave;

    /**
     * Decodes a Leave message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Leave
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Leave;

    /**
     * Verifies a Leave message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Leave message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Leave
     */
    public static fromObject(object: { [k: string]: any }): Leave;

    /**
     * Creates a plain object from a Leave message. Also converts values to other types if specified.
     * @param message Leave
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Leave, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Leave to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a Text. */
export interface IText {

    /** Text block */
    block: string;

    /** Text date */
    date: google.protobuf.ITimestamp;

    /** Text user */
    user: IUser;

    /** Text body */
    body: string;

    /** Text comments */
    comments: IComment[];

    /** Text likes */
    likes: ILike[];
}

/** Represents a Text. */
export class Text implements IText {

    /**
     * Constructs a new Text.
     * @param [properties] Properties to set
     */
    constructor(properties?: IText);

    /** Text block. */
    public block: string;

    /** Text date. */
    public date: google.protobuf.ITimestamp;

    /** Text user. */
    public user: IUser;

    /** Text body. */
    public body: string;

    /** Text comments. */
    public comments: IComment[];

    /** Text likes. */
    public likes: ILike[];

    /**
     * Creates a new Text instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Text instance
     */
    public static create(properties?: IText): Text;

    /**
     * Encodes the specified Text message. Does not implicitly {@link Text.verify|verify} messages.
     * @param message Text message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IText, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Text message, length delimited. Does not implicitly {@link Text.verify|verify} messages.
     * @param message Text message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IText, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Text message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Text
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Text;

    /**
     * Decodes a Text message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Text
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Text;

    /**
     * Verifies a Text message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Text message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Text
     */
    public static fromObject(object: { [k: string]: any }): Text;

    /**
     * Creates a plain object from a Text message. Also converts values to other types if specified.
     * @param message Text
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Text, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Text to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a TextList. */
export interface ITextList {

    /** TextList items */
    items: IText[];
}

/** Represents a TextList. */
export class TextList implements ITextList {

    /**
     * Constructs a new TextList.
     * @param [properties] Properties to set
     */
    constructor(properties?: ITextList);

    /** TextList items. */
    public items: IText[];

    /**
     * Creates a new TextList instance using the specified properties.
     * @param [properties] Properties to set
     * @returns TextList instance
     */
    public static create(properties?: ITextList): TextList;

    /**
     * Encodes the specified TextList message. Does not implicitly {@link TextList.verify|verify} messages.
     * @param message TextList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ITextList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified TextList message, length delimited. Does not implicitly {@link TextList.verify|verify} messages.
     * @param message TextList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ITextList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a TextList message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns TextList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): TextList;

    /**
     * Decodes a TextList message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns TextList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): TextList;

    /**
     * Verifies a TextList message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a TextList message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns TextList
     */
    public static fromObject(object: { [k: string]: any }): TextList;

    /**
     * Creates a plain object from a TextList message. Also converts values to other types if specified.
     * @param message TextList
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: TextList, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this TextList to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a File. */
export interface IFile {

    /** File index */
    index: number;

    /** File file */
    file: IFileIndex;

    /** File links */
    links: { [k: string]: IFileIndex };
}

/** Represents a File. */
export class File implements IFile {

    /**
     * Constructs a new File.
     * @param [properties] Properties to set
     */
    constructor(properties?: IFile);

    /** File index. */
    public index: number;

    /** File file. */
    public file: IFileIndex;

    /** File links. */
    public links: { [k: string]: IFileIndex };

    /**
     * Creates a new File instance using the specified properties.
     * @param [properties] Properties to set
     * @returns File instance
     */
    public static create(properties?: IFile): File;

    /**
     * Encodes the specified File message. Does not implicitly {@link File.verify|verify} messages.
     * @param message File message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IFile, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified File message, length delimited. Does not implicitly {@link File.verify|verify} messages.
     * @param message File message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IFile, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a File message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns File
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): File;

    /**
     * Decodes a File message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns File
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): File;

    /**
     * Verifies a File message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a File message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns File
     */
    public static fromObject(object: { [k: string]: any }): File;

    /**
     * Creates a plain object from a File message. Also converts values to other types if specified.
     * @param message File
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: File, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this File to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a Files. */
export interface IFiles {

    /** Files block */
    block: string;

    /** Files target */
    target: string;

    /** Files date */
    date: google.protobuf.ITimestamp;

    /** Files user */
    user: IUser;

    /** Files caption */
    caption: string;

    /** Files files */
    files: IFile[];

    /** Files comments */
    comments: IComment[];

    /** Files likes */
    likes: ILike[];

    /** Files threads */
    threads: string[];
}

/** Represents a Files. */
export class Files implements IFiles {

    /**
     * Constructs a new Files.
     * @param [properties] Properties to set
     */
    constructor(properties?: IFiles);

    /** Files block. */
    public block: string;

    /** Files target. */
    public target: string;

    /** Files date. */
    public date: google.protobuf.ITimestamp;

    /** Files user. */
    public user: IUser;

    /** Files caption. */
    public caption: string;

    /** Files files. */
    public files: IFile[];

    /** Files comments. */
    public comments: IComment[];

    /** Files likes. */
    public likes: ILike[];

    /** Files threads. */
    public threads: string[];

    /**
     * Creates a new Files instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Files instance
     */
    public static create(properties?: IFiles): Files;

    /**
     * Encodes the specified Files message. Does not implicitly {@link Files.verify|verify} messages.
     * @param message Files message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IFiles, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Files message, length delimited. Does not implicitly {@link Files.verify|verify} messages.
     * @param message Files message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IFiles, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Files message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Files
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Files;

    /**
     * Decodes a Files message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Files
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Files;

    /**
     * Verifies a Files message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Files message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Files
     */
    public static fromObject(object: { [k: string]: any }): Files;

    /**
     * Creates a plain object from a Files message. Also converts values to other types if specified.
     * @param message Files
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Files, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Files to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a FilesList. */
export interface IFilesList {

    /** FilesList items */
    items: IFiles[];
}

/** Represents a FilesList. */
export class FilesList implements IFilesList {

    /**
     * Constructs a new FilesList.
     * @param [properties] Properties to set
     */
    constructor(properties?: IFilesList);

    /** FilesList items. */
    public items: IFiles[];

    /**
     * Creates a new FilesList instance using the specified properties.
     * @param [properties] Properties to set
     * @returns FilesList instance
     */
    public static create(properties?: IFilesList): FilesList;

    /**
     * Encodes the specified FilesList message. Does not implicitly {@link FilesList.verify|verify} messages.
     * @param message FilesList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IFilesList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified FilesList message, length delimited. Does not implicitly {@link FilesList.verify|verify} messages.
     * @param message FilesList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IFilesList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a FilesList message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns FilesList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): FilesList;

    /**
     * Decodes a FilesList message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns FilesList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): FilesList;

    /**
     * Verifies a FilesList message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a FilesList message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns FilesList
     */
    public static fromObject(object: { [k: string]: any }): FilesList;

    /**
     * Creates a plain object from a FilesList message. Also converts values to other types if specified.
     * @param message FilesList
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: FilesList, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this FilesList to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a Comment. */
export interface IComment {

    /** Comment id */
    id: string;

    /** Comment date */
    date: google.protobuf.ITimestamp;

    /** Comment user */
    user: IUser;

    /** Comment body */
    body: string;

    /** Comment target */
    target: IFeedItem;
}

/** Represents a Comment. */
export class Comment implements IComment {

    /**
     * Constructs a new Comment.
     * @param [properties] Properties to set
     */
    constructor(properties?: IComment);

    /** Comment id. */
    public id: string;

    /** Comment date. */
    public date: google.protobuf.ITimestamp;

    /** Comment user. */
    public user: IUser;

    /** Comment body. */
    public body: string;

    /** Comment target. */
    public target: IFeedItem;

    /**
     * Creates a new Comment instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Comment instance
     */
    public static create(properties?: IComment): Comment;

    /**
     * Encodes the specified Comment message. Does not implicitly {@link Comment.verify|verify} messages.
     * @param message Comment message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IComment, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Comment message, length delimited. Does not implicitly {@link Comment.verify|verify} messages.
     * @param message Comment message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IComment, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Comment message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Comment
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Comment;

    /**
     * Decodes a Comment message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Comment
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Comment;

    /**
     * Verifies a Comment message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Comment message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Comment
     */
    public static fromObject(object: { [k: string]: any }): Comment;

    /**
     * Creates a plain object from a Comment message. Also converts values to other types if specified.
     * @param message Comment
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Comment, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Comment to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a CommentList. */
export interface ICommentList {

    /** CommentList items */
    items: IComment[];
}

/** Represents a CommentList. */
export class CommentList implements ICommentList {

    /**
     * Constructs a new CommentList.
     * @param [properties] Properties to set
     */
    constructor(properties?: ICommentList);

    /** CommentList items. */
    public items: IComment[];

    /**
     * Creates a new CommentList instance using the specified properties.
     * @param [properties] Properties to set
     * @returns CommentList instance
     */
    public static create(properties?: ICommentList): CommentList;

    /**
     * Encodes the specified CommentList message. Does not implicitly {@link CommentList.verify|verify} messages.
     * @param message CommentList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ICommentList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified CommentList message, length delimited. Does not implicitly {@link CommentList.verify|verify} messages.
     * @param message CommentList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ICommentList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a CommentList message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns CommentList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): CommentList;

    /**
     * Decodes a CommentList message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns CommentList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): CommentList;

    /**
     * Verifies a CommentList message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a CommentList message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns CommentList
     */
    public static fromObject(object: { [k: string]: any }): CommentList;

    /**
     * Creates a plain object from a CommentList message. Also converts values to other types if specified.
     * @param message CommentList
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: CommentList, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this CommentList to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a Like. */
export interface ILike {

    /** Like id */
    id: string;

    /** Like date */
    date: google.protobuf.ITimestamp;

    /** Like user */
    user: IUser;

    /** Like target */
    target: IFeedItem;
}

/** Represents a Like. */
export class Like implements ILike {

    /**
     * Constructs a new Like.
     * @param [properties] Properties to set
     */
    constructor(properties?: ILike);

    /** Like id. */
    public id: string;

    /** Like date. */
    public date: google.protobuf.ITimestamp;

    /** Like user. */
    public user: IUser;

    /** Like target. */
    public target: IFeedItem;

    /**
     * Creates a new Like instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Like instance
     */
    public static create(properties?: ILike): Like;

    /**
     * Encodes the specified Like message. Does not implicitly {@link Like.verify|verify} messages.
     * @param message Like message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ILike, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Like message, length delimited. Does not implicitly {@link Like.verify|verify} messages.
     * @param message Like message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ILike, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Like message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Like
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Like;

    /**
     * Decodes a Like message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Like
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Like;

    /**
     * Verifies a Like message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Like message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Like
     */
    public static fromObject(object: { [k: string]: any }): Like;

    /**
     * Creates a plain object from a Like message. Also converts values to other types if specified.
     * @param message Like
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Like, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Like to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a LikeList. */
export interface ILikeList {

    /** LikeList items */
    items: ILike[];
}

/** Represents a LikeList. */
export class LikeList implements ILikeList {

    /**
     * Constructs a new LikeList.
     * @param [properties] Properties to set
     */
    constructor(properties?: ILikeList);

    /** LikeList items. */
    public items: ILike[];

    /**
     * Creates a new LikeList instance using the specified properties.
     * @param [properties] Properties to set
     * @returns LikeList instance
     */
    public static create(properties?: ILikeList): LikeList;

    /**
     * Encodes the specified LikeList message. Does not implicitly {@link LikeList.verify|verify} messages.
     * @param message LikeList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ILikeList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified LikeList message, length delimited. Does not implicitly {@link LikeList.verify|verify} messages.
     * @param message LikeList message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ILikeList, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a LikeList message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns LikeList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): LikeList;

    /**
     * Decodes a LikeList message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns LikeList
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): LikeList;

    /**
     * Verifies a LikeList message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a LikeList message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns LikeList
     */
    public static fromObject(object: { [k: string]: any }): LikeList;

    /**
     * Creates a plain object from a LikeList message. Also converts values to other types if specified.
     * @param message LikeList
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: LikeList, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this LikeList to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a Summary. */
export interface ISummary {

    /** Summary accountPeerCount */
    accountPeerCount: number;

    /** Summary threadCount */
    threadCount: number;

    /** Summary fileCount */
    fileCount: number;

    /** Summary contactCount */
    contactCount: number;
}

/** Represents a Summary. */
export class Summary implements ISummary {

    /**
     * Constructs a new Summary.
     * @param [properties] Properties to set
     */
    constructor(properties?: ISummary);

    /** Summary accountPeerCount. */
    public accountPeerCount: number;

    /** Summary threadCount. */
    public threadCount: number;

    /** Summary fileCount. */
    public fileCount: number;

    /** Summary contactCount. */
    public contactCount: number;

    /**
     * Creates a new Summary instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Summary instance
     */
    public static create(properties?: ISummary): Summary;

    /**
     * Encodes the specified Summary message. Does not implicitly {@link Summary.verify|verify} messages.
     * @param message Summary message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: ISummary, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Summary message, length delimited. Does not implicitly {@link Summary.verify|verify} messages.
     * @param message Summary message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: ISummary, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Summary message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Summary
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Summary;

    /**
     * Decodes a Summary message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Summary
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Summary;

    /**
     * Verifies a Summary message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Summary message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Summary
     */
    public static fromObject(object: { [k: string]: any }): Summary;

    /**
     * Creates a plain object from a Summary message. Also converts values to other types if specified.
     * @param message Summary
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Summary, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Summary to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** QueryType enum. */
export enum QueryType {
    THREAD_BACKUPS = 0,
    CONTACTS = 1
}

/** Properties of a QueryOptions. */
export interface IQueryOptions {

    /** QueryOptions local */
    local: boolean;

    /** QueryOptions limit */
    limit: number;

    /** QueryOptions wait */
    wait: number;

    /** QueryOptions filter */
    filter: QueryOptions.FilterType;

    /** QueryOptions exclude */
    exclude: string[];
}

/** Represents a QueryOptions. */
export class QueryOptions implements IQueryOptions {

    /**
     * Constructs a new QueryOptions.
     * @param [properties] Properties to set
     */
    constructor(properties?: IQueryOptions);

    /** QueryOptions local. */
    public local: boolean;

    /** QueryOptions limit. */
    public limit: number;

    /** QueryOptions wait. */
    public wait: number;

    /** QueryOptions filter. */
    public filter: QueryOptions.FilterType;

    /** QueryOptions exclude. */
    public exclude: string[];

    /**
     * Creates a new QueryOptions instance using the specified properties.
     * @param [properties] Properties to set
     * @returns QueryOptions instance
     */
    public static create(properties?: IQueryOptions): QueryOptions;

    /**
     * Encodes the specified QueryOptions message. Does not implicitly {@link QueryOptions.verify|verify} messages.
     * @param message QueryOptions message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IQueryOptions, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified QueryOptions message, length delimited. Does not implicitly {@link QueryOptions.verify|verify} messages.
     * @param message QueryOptions message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IQueryOptions, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a QueryOptions message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns QueryOptions
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): QueryOptions;

    /**
     * Decodes a QueryOptions message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns QueryOptions
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): QueryOptions;

    /**
     * Verifies a QueryOptions message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a QueryOptions message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns QueryOptions
     */
    public static fromObject(object: { [k: string]: any }): QueryOptions;

    /**
     * Creates a plain object from a QueryOptions message. Also converts values to other types if specified.
     * @param message QueryOptions
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: QueryOptions, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this QueryOptions to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

export namespace QueryOptions {

    /** FilterType enum. */
    enum FilterType {
        NO_FILTER = 0,
        HIDE_OLDER = 1
    }
}

/** Properties of a Query. */
export interface IQuery {

    /** Query id */
    id: string;

    /** Query token */
    token: string;

    /** Query type */
    type: QueryType;

    /** Query options */
    options: IQueryOptions;

    /** Query payload */
    payload: google.protobuf.IAny;
}

/** Represents a Query. */
export class Query implements IQuery {

    /**
     * Constructs a new Query.
     * @param [properties] Properties to set
     */
    constructor(properties?: IQuery);

    /** Query id. */
    public id: string;

    /** Query token. */
    public token: string;

    /** Query type. */
    public type: QueryType;

    /** Query options. */
    public options: IQueryOptions;

    /** Query payload. */
    public payload: google.protobuf.IAny;

    /**
     * Creates a new Query instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Query instance
     */
    public static create(properties?: IQuery): Query;

    /**
     * Encodes the specified Query message. Does not implicitly {@link Query.verify|verify} messages.
     * @param message Query message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IQuery, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Query message, length delimited. Does not implicitly {@link Query.verify|verify} messages.
     * @param message Query message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IQuery, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Query message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Query
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Query;

    /**
     * Decodes a Query message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Query
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Query;

    /**
     * Verifies a Query message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Query message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Query
     */
    public static fromObject(object: { [k: string]: any }): Query;

    /**
     * Creates a plain object from a Query message. Also converts values to other types if specified.
     * @param message Query
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Query, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Query to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a PubSubQuery. */
export interface IPubSubQuery {

    /** PubSubQuery id */
    id: string;

    /** PubSubQuery type */
    type: QueryType;

    /** PubSubQuery payload */
    payload: google.protobuf.IAny;

    /** PubSubQuery responseType */
    responseType: PubSubQuery.ResponseType;
}

/** Represents a PubSubQuery. */
export class PubSubQuery implements IPubSubQuery {

    /**
     * Constructs a new PubSubQuery.
     * @param [properties] Properties to set
     */
    constructor(properties?: IPubSubQuery);

    /** PubSubQuery id. */
    public id: string;

    /** PubSubQuery type. */
    public type: QueryType;

    /** PubSubQuery payload. */
    public payload: google.protobuf.IAny;

    /** PubSubQuery responseType. */
    public responseType: PubSubQuery.ResponseType;

    /**
     * Creates a new PubSubQuery instance using the specified properties.
     * @param [properties] Properties to set
     * @returns PubSubQuery instance
     */
    public static create(properties?: IPubSubQuery): PubSubQuery;

    /**
     * Encodes the specified PubSubQuery message. Does not implicitly {@link PubSubQuery.verify|verify} messages.
     * @param message PubSubQuery message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IPubSubQuery, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified PubSubQuery message, length delimited. Does not implicitly {@link PubSubQuery.verify|verify} messages.
     * @param message PubSubQuery message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IPubSubQuery, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a PubSubQuery message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns PubSubQuery
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): PubSubQuery;

    /**
     * Decodes a PubSubQuery message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns PubSubQuery
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): PubSubQuery;

    /**
     * Verifies a PubSubQuery message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a PubSubQuery message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns PubSubQuery
     */
    public static fromObject(object: { [k: string]: any }): PubSubQuery;

    /**
     * Creates a plain object from a PubSubQuery message. Also converts values to other types if specified.
     * @param message PubSubQuery
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: PubSubQuery, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this PubSubQuery to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

export namespace PubSubQuery {

    /** ResponseType enum. */
    enum ResponseType {
        P2P = 0,
        PUBSUB = 1
    }
}

/** Properties of a QueryResult. */
export interface IQueryResult {

    /** QueryResult id */
    id: string;

    /** QueryResult date */
    date: google.protobuf.ITimestamp;

    /** QueryResult local */
    local: boolean;

    /** QueryResult value */
    value: google.protobuf.IAny;
}

/** Represents a QueryResult. */
export class QueryResult implements IQueryResult {

    /**
     * Constructs a new QueryResult.
     * @param [properties] Properties to set
     */
    constructor(properties?: IQueryResult);

    /** QueryResult id. */
    public id: string;

    /** QueryResult date. */
    public date: google.protobuf.ITimestamp;

    /** QueryResult local. */
    public local: boolean;

    /** QueryResult value. */
    public value: google.protobuf.IAny;

    /**
     * Creates a new QueryResult instance using the specified properties.
     * @param [properties] Properties to set
     * @returns QueryResult instance
     */
    public static create(properties?: IQueryResult): QueryResult;

    /**
     * Encodes the specified QueryResult message. Does not implicitly {@link QueryResult.verify|verify} messages.
     * @param message QueryResult message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IQueryResult, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified QueryResult message, length delimited. Does not implicitly {@link QueryResult.verify|verify} messages.
     * @param message QueryResult message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IQueryResult, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a QueryResult message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns QueryResult
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): QueryResult;

    /**
     * Decodes a QueryResult message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns QueryResult
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): QueryResult;

    /**
     * Verifies a QueryResult message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a QueryResult message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns QueryResult
     */
    public static fromObject(object: { [k: string]: any }): QueryResult;

    /**
     * Creates a plain object from a QueryResult message. Also converts values to other types if specified.
     * @param message QueryResult
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: QueryResult, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this QueryResult to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a QueryResults. */
export interface IQueryResults {

    /** QueryResults type */
    type: QueryType;

    /** QueryResults items */
    items: IQueryResult[];
}

/** Represents a QueryResults. */
export class QueryResults implements IQueryResults {

    /**
     * Constructs a new QueryResults.
     * @param [properties] Properties to set
     */
    constructor(properties?: IQueryResults);

    /** QueryResults type. */
    public type: QueryType;

    /** QueryResults items. */
    public items: IQueryResult[];

    /**
     * Creates a new QueryResults instance using the specified properties.
     * @param [properties] Properties to set
     * @returns QueryResults instance
     */
    public static create(properties?: IQueryResults): QueryResults;

    /**
     * Encodes the specified QueryResults message. Does not implicitly {@link QueryResults.verify|verify} messages.
     * @param message QueryResults message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IQueryResults, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified QueryResults message, length delimited. Does not implicitly {@link QueryResults.verify|verify} messages.
     * @param message QueryResults message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IQueryResults, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a QueryResults message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns QueryResults
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): QueryResults;

    /**
     * Decodes a QueryResults message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns QueryResults
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): QueryResults;

    /**
     * Verifies a QueryResults message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a QueryResults message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns QueryResults
     */
    public static fromObject(object: { [k: string]: any }): QueryResults;

    /**
     * Creates a plain object from a QueryResults message. Also converts values to other types if specified.
     * @param message QueryResults
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: QueryResults, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this QueryResults to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a PubSubQueryResults. */
export interface IPubSubQueryResults {

    /** PubSubQueryResults id */
    id: string;

    /** PubSubQueryResults results */
    results: IQueryResults;
}

/** Represents a PubSubQueryResults. */
export class PubSubQueryResults implements IPubSubQueryResults {

    /**
     * Constructs a new PubSubQueryResults.
     * @param [properties] Properties to set
     */
    constructor(properties?: IPubSubQueryResults);

    /** PubSubQueryResults id. */
    public id: string;

    /** PubSubQueryResults results. */
    public results: IQueryResults;

    /**
     * Creates a new PubSubQueryResults instance using the specified properties.
     * @param [properties] Properties to set
     * @returns PubSubQueryResults instance
     */
    public static create(properties?: IPubSubQueryResults): PubSubQueryResults;

    /**
     * Encodes the specified PubSubQueryResults message. Does not implicitly {@link PubSubQueryResults.verify|verify} messages.
     * @param message PubSubQueryResults message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IPubSubQueryResults, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified PubSubQueryResults message, length delimited. Does not implicitly {@link PubSubQueryResults.verify|verify} messages.
     * @param message PubSubQueryResults message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IPubSubQueryResults, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a PubSubQueryResults message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns PubSubQueryResults
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): PubSubQueryResults;

    /**
     * Decodes a PubSubQueryResults message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns PubSubQueryResults
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): PubSubQueryResults;

    /**
     * Verifies a PubSubQueryResults message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a PubSubQueryResults message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns PubSubQueryResults
     */
    public static fromObject(object: { [k: string]: any }): PubSubQueryResults;

    /**
     * Creates a plain object from a PubSubQueryResults message. Also converts values to other types if specified.
     * @param message PubSubQueryResults
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: PubSubQueryResults, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this PubSubQueryResults to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a QueryEvent. */
export interface IQueryEvent {

    /** QueryEvent type */
    type: QueryEvent.Type;

    /** QueryEvent data */
    data: IQueryResult;
}

/** Represents a QueryEvent. */
export class QueryEvent implements IQueryEvent {

    /**
     * Constructs a new QueryEvent.
     * @param [properties] Properties to set
     */
    constructor(properties?: IQueryEvent);

    /** QueryEvent type. */
    public type: QueryEvent.Type;

    /** QueryEvent data. */
    public data: IQueryResult;

    /**
     * Creates a new QueryEvent instance using the specified properties.
     * @param [properties] Properties to set
     * @returns QueryEvent instance
     */
    public static create(properties?: IQueryEvent): QueryEvent;

    /**
     * Encodes the specified QueryEvent message. Does not implicitly {@link QueryEvent.verify|verify} messages.
     * @param message QueryEvent message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IQueryEvent, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified QueryEvent message, length delimited. Does not implicitly {@link QueryEvent.verify|verify} messages.
     * @param message QueryEvent message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IQueryEvent, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a QueryEvent message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns QueryEvent
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): QueryEvent;

    /**
     * Decodes a QueryEvent message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns QueryEvent
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): QueryEvent;

    /**
     * Verifies a QueryEvent message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a QueryEvent message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns QueryEvent
     */
    public static fromObject(object: { [k: string]: any }): QueryEvent;

    /**
     * Creates a plain object from a QueryEvent message. Also converts values to other types if specified.
     * @param message QueryEvent
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: QueryEvent, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this QueryEvent to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

export namespace QueryEvent {

    /** Type enum. */
    enum Type {
        DATA = 0,
        DONE = 1
    }
}

/** Properties of a ContactQuery. */
export interface IContactQuery {

    /** ContactQuery id */
    id: string;

    /** ContactQuery address */
    address: string;

    /** ContactQuery username */
    username: string;
}

/** Represents a ContactQuery. */
export class ContactQuery implements IContactQuery {

    /**
     * Constructs a new ContactQuery.
     * @param [properties] Properties to set
     */
    constructor(properties?: IContactQuery);

    /** ContactQuery id. */
    public id: string;

    /** ContactQuery address. */
    public address: string;

    /** ContactQuery username. */
    public username: string;

    /**
     * Creates a new ContactQuery instance using the specified properties.
     * @param [properties] Properties to set
     * @returns ContactQuery instance
     */
    public static create(properties?: IContactQuery): ContactQuery;

    /**
     * Encodes the specified ContactQuery message. Does not implicitly {@link ContactQuery.verify|verify} messages.
     * @param message ContactQuery message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IContactQuery, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified ContactQuery message, length delimited. Does not implicitly {@link ContactQuery.verify|verify} messages.
     * @param message ContactQuery message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IContactQuery, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a ContactQuery message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns ContactQuery
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): ContactQuery;

    /**
     * Decodes a ContactQuery message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns ContactQuery
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): ContactQuery;

    /**
     * Verifies a ContactQuery message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a ContactQuery message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns ContactQuery
     */
    public static fromObject(object: { [k: string]: any }): ContactQuery;

    /**
     * Creates a plain object from a ContactQuery message. Also converts values to other types if specified.
     * @param message ContactQuery
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: ContactQuery, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this ContactQuery to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a ThreadBackupQuery. */
export interface IThreadBackupQuery {

    /** ThreadBackupQuery address */
    address: string;
}

/** Represents a ThreadBackupQuery. */
export class ThreadBackupQuery implements IThreadBackupQuery {

    /**
     * Constructs a new ThreadBackupQuery.
     * @param [properties] Properties to set
     */
    constructor(properties?: IThreadBackupQuery);

    /** ThreadBackupQuery address. */
    public address: string;

    /**
     * Creates a new ThreadBackupQuery instance using the specified properties.
     * @param [properties] Properties to set
     * @returns ThreadBackupQuery instance
     */
    public static create(properties?: IThreadBackupQuery): ThreadBackupQuery;

    /**
     * Encodes the specified ThreadBackupQuery message. Does not implicitly {@link ThreadBackupQuery.verify|verify} messages.
     * @param message ThreadBackupQuery message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IThreadBackupQuery, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified ThreadBackupQuery message, length delimited. Does not implicitly {@link ThreadBackupQuery.verify|verify} messages.
     * @param message ThreadBackupQuery message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IThreadBackupQuery, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a ThreadBackupQuery message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns ThreadBackupQuery
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): ThreadBackupQuery;

    /**
     * Decodes a ThreadBackupQuery message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns ThreadBackupQuery
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): ThreadBackupQuery;

    /**
     * Verifies a ThreadBackupQuery message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a ThreadBackupQuery message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns ThreadBackupQuery
     */
    public static fromObject(object: { [k: string]: any }): ThreadBackupQuery;

    /**
     * Creates a plain object from a ThreadBackupQuery message. Also converts values to other types if specified.
     * @param message ThreadBackupQuery
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: ThreadBackupQuery, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this ThreadBackupQuery to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a ThreadEnvelope. */
export interface IThreadEnvelope {

    /** ThreadEnvelope thread */
    thread: string;

    /** ThreadEnvelope hash */
    hash: string;

    /** ThreadEnvelope ciphertext */
    ciphertext: Uint8Array;
}

/** Represents a ThreadEnvelope. */
export class ThreadEnvelope implements IThreadEnvelope {

    /**
     * Constructs a new ThreadEnvelope.
     * @param [properties] Properties to set
     */
    constructor(properties?: IThreadEnvelope);

    /** ThreadEnvelope thread. */
    public thread: string;

    /** ThreadEnvelope hash. */
    public hash: string;

    /** ThreadEnvelope ciphertext. */
    public ciphertext: Uint8Array;

    /**
     * Creates a new ThreadEnvelope instance using the specified properties.
     * @param [properties] Properties to set
     * @returns ThreadEnvelope instance
     */
    public static create(properties?: IThreadEnvelope): ThreadEnvelope;

    /**
     * Encodes the specified ThreadEnvelope message. Does not implicitly {@link ThreadEnvelope.verify|verify} messages.
     * @param message ThreadEnvelope message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IThreadEnvelope, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified ThreadEnvelope message, length delimited. Does not implicitly {@link ThreadEnvelope.verify|verify} messages.
     * @param message ThreadEnvelope message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IThreadEnvelope, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a ThreadEnvelope message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns ThreadEnvelope
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): ThreadEnvelope;

    /**
     * Decodes a ThreadEnvelope message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns ThreadEnvelope
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): ThreadEnvelope;

    /**
     * Verifies a ThreadEnvelope message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a ThreadEnvelope message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns ThreadEnvelope
     */
    public static fromObject(object: { [k: string]: any }): ThreadEnvelope;

    /**
     * Creates a plain object from a ThreadEnvelope message. Also converts values to other types if specified.
     * @param message ThreadEnvelope
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: ThreadEnvelope, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this ThreadEnvelope to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a ThreadBlock. */
export interface IThreadBlock {

    /** ThreadBlock header */
    header: IThreadBlockHeader;

    /** ThreadBlock type */
    type: Block.BlockType;

    /** ThreadBlock payload */
    payload: google.protobuf.IAny;
}

/** Represents a ThreadBlock. */
export class ThreadBlock implements IThreadBlock {

    /**
     * Constructs a new ThreadBlock.
     * @param [properties] Properties to set
     */
    constructor(properties?: IThreadBlock);

    /** ThreadBlock header. */
    public header: IThreadBlockHeader;

    /** ThreadBlock type. */
    public type: Block.BlockType;

    /** ThreadBlock payload. */
    public payload: google.protobuf.IAny;

    /**
     * Creates a new ThreadBlock instance using the specified properties.
     * @param [properties] Properties to set
     * @returns ThreadBlock instance
     */
    public static create(properties?: IThreadBlock): ThreadBlock;

    /**
     * Encodes the specified ThreadBlock message. Does not implicitly {@link ThreadBlock.verify|verify} messages.
     * @param message ThreadBlock message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IThreadBlock, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified ThreadBlock message, length delimited. Does not implicitly {@link ThreadBlock.verify|verify} messages.
     * @param message ThreadBlock message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IThreadBlock, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a ThreadBlock message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns ThreadBlock
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): ThreadBlock;

    /**
     * Decodes a ThreadBlock message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns ThreadBlock
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): ThreadBlock;

    /**
     * Verifies a ThreadBlock message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a ThreadBlock message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns ThreadBlock
     */
    public static fromObject(object: { [k: string]: any }): ThreadBlock;

    /**
     * Creates a plain object from a ThreadBlock message. Also converts values to other types if specified.
     * @param message ThreadBlock
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: ThreadBlock, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this ThreadBlock to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a ThreadBlockHeader. */
export interface IThreadBlockHeader {

    /** ThreadBlockHeader date */
    date: google.protobuf.ITimestamp;

    /** ThreadBlockHeader parents */
    parents: string[];

    /** ThreadBlockHeader author */
    author: string;

    /** ThreadBlockHeader address */
    address: string;
}

/** Represents a ThreadBlockHeader. */
export class ThreadBlockHeader implements IThreadBlockHeader {

    /**
     * Constructs a new ThreadBlockHeader.
     * @param [properties] Properties to set
     */
    constructor(properties?: IThreadBlockHeader);

    /** ThreadBlockHeader date. */
    public date: google.protobuf.ITimestamp;

    /** ThreadBlockHeader parents. */
    public parents: string[];

    /** ThreadBlockHeader author. */
    public author: string;

    /** ThreadBlockHeader address. */
    public address: string;

    /**
     * Creates a new ThreadBlockHeader instance using the specified properties.
     * @param [properties] Properties to set
     * @returns ThreadBlockHeader instance
     */
    public static create(properties?: IThreadBlockHeader): ThreadBlockHeader;

    /**
     * Encodes the specified ThreadBlockHeader message. Does not implicitly {@link ThreadBlockHeader.verify|verify} messages.
     * @param message ThreadBlockHeader message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IThreadBlockHeader, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified ThreadBlockHeader message, length delimited. Does not implicitly {@link ThreadBlockHeader.verify|verify} messages.
     * @param message ThreadBlockHeader message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IThreadBlockHeader, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a ThreadBlockHeader message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns ThreadBlockHeader
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): ThreadBlockHeader;

    /**
     * Decodes a ThreadBlockHeader message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns ThreadBlockHeader
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): ThreadBlockHeader;

    /**
     * Verifies a ThreadBlockHeader message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a ThreadBlockHeader message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns ThreadBlockHeader
     */
    public static fromObject(object: { [k: string]: any }): ThreadBlockHeader;

    /**
     * Creates a plain object from a ThreadBlockHeader message. Also converts values to other types if specified.
     * @param message ThreadBlockHeader
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: ThreadBlockHeader, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this ThreadBlockHeader to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a ThreadInvite. */
export interface IThreadInvite {

    /** ThreadInvite inviter */
    inviter: IContact;

    /** ThreadInvite thread */
    thread: IThread;
}

/** Represents a ThreadInvite. */
export class ThreadInvite implements IThreadInvite {

    /**
     * Constructs a new ThreadInvite.
     * @param [properties] Properties to set
     */
    constructor(properties?: IThreadInvite);

    /** ThreadInvite inviter. */
    public inviter: IContact;

    /** ThreadInvite thread. */
    public thread: IThread;

    /**
     * Creates a new ThreadInvite instance using the specified properties.
     * @param [properties] Properties to set
     * @returns ThreadInvite instance
     */
    public static create(properties?: IThreadInvite): ThreadInvite;

    /**
     * Encodes the specified ThreadInvite message. Does not implicitly {@link ThreadInvite.verify|verify} messages.
     * @param message ThreadInvite message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IThreadInvite, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified ThreadInvite message, length delimited. Does not implicitly {@link ThreadInvite.verify|verify} messages.
     * @param message ThreadInvite message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IThreadInvite, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a ThreadInvite message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns ThreadInvite
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): ThreadInvite;

    /**
     * Decodes a ThreadInvite message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns ThreadInvite
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): ThreadInvite;

    /**
     * Verifies a ThreadInvite message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a ThreadInvite message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns ThreadInvite
     */
    public static fromObject(object: { [k: string]: any }): ThreadInvite;

    /**
     * Creates a plain object from a ThreadInvite message. Also converts values to other types if specified.
     * @param message ThreadInvite
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: ThreadInvite, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this ThreadInvite to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a ThreadIgnore. */
export interface IThreadIgnore {

    /** ThreadIgnore target */
    target: string;
}

/** Represents a ThreadIgnore. */
export class ThreadIgnore implements IThreadIgnore {

    /**
     * Constructs a new ThreadIgnore.
     * @param [properties] Properties to set
     */
    constructor(properties?: IThreadIgnore);

    /** ThreadIgnore target. */
    public target: string;

    /**
     * Creates a new ThreadIgnore instance using the specified properties.
     * @param [properties] Properties to set
     * @returns ThreadIgnore instance
     */
    public static create(properties?: IThreadIgnore): ThreadIgnore;

    /**
     * Encodes the specified ThreadIgnore message. Does not implicitly {@link ThreadIgnore.verify|verify} messages.
     * @param message ThreadIgnore message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IThreadIgnore, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified ThreadIgnore message, length delimited. Does not implicitly {@link ThreadIgnore.verify|verify} messages.
     * @param message ThreadIgnore message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IThreadIgnore, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a ThreadIgnore message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns ThreadIgnore
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): ThreadIgnore;

    /**
     * Decodes a ThreadIgnore message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns ThreadIgnore
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): ThreadIgnore;

    /**
     * Verifies a ThreadIgnore message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a ThreadIgnore message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns ThreadIgnore
     */
    public static fromObject(object: { [k: string]: any }): ThreadIgnore;

    /**
     * Creates a plain object from a ThreadIgnore message. Also converts values to other types if specified.
     * @param message ThreadIgnore
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: ThreadIgnore, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this ThreadIgnore to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a ThreadFlag. */
export interface IThreadFlag {

    /** ThreadFlag target */
    target: string;
}

/** Represents a ThreadFlag. */
export class ThreadFlag implements IThreadFlag {

    /**
     * Constructs a new ThreadFlag.
     * @param [properties] Properties to set
     */
    constructor(properties?: IThreadFlag);

    /** ThreadFlag target. */
    public target: string;

    /**
     * Creates a new ThreadFlag instance using the specified properties.
     * @param [properties] Properties to set
     * @returns ThreadFlag instance
     */
    public static create(properties?: IThreadFlag): ThreadFlag;

    /**
     * Encodes the specified ThreadFlag message. Does not implicitly {@link ThreadFlag.verify|verify} messages.
     * @param message ThreadFlag message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IThreadFlag, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified ThreadFlag message, length delimited. Does not implicitly {@link ThreadFlag.verify|verify} messages.
     * @param message ThreadFlag message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IThreadFlag, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a ThreadFlag message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns ThreadFlag
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): ThreadFlag;

    /**
     * Decodes a ThreadFlag message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns ThreadFlag
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): ThreadFlag;

    /**
     * Verifies a ThreadFlag message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a ThreadFlag message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns ThreadFlag
     */
    public static fromObject(object: { [k: string]: any }): ThreadFlag;

    /**
     * Creates a plain object from a ThreadFlag message. Also converts values to other types if specified.
     * @param message ThreadFlag
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: ThreadFlag, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this ThreadFlag to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a ThreadJoin. */
export interface IThreadJoin {

    /** ThreadJoin inviter */
    inviter: string;

    /** ThreadJoin contact */
    contact: IContact;
}

/** Represents a ThreadJoin. */
export class ThreadJoin implements IThreadJoin {

    /**
     * Constructs a new ThreadJoin.
     * @param [properties] Properties to set
     */
    constructor(properties?: IThreadJoin);

    /** ThreadJoin inviter. */
    public inviter: string;

    /** ThreadJoin contact. */
    public contact: IContact;

    /**
     * Creates a new ThreadJoin instance using the specified properties.
     * @param [properties] Properties to set
     * @returns ThreadJoin instance
     */
    public static create(properties?: IThreadJoin): ThreadJoin;

    /**
     * Encodes the specified ThreadJoin message. Does not implicitly {@link ThreadJoin.verify|verify} messages.
     * @param message ThreadJoin message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IThreadJoin, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified ThreadJoin message, length delimited. Does not implicitly {@link ThreadJoin.verify|verify} messages.
     * @param message ThreadJoin message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IThreadJoin, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a ThreadJoin message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns ThreadJoin
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): ThreadJoin;

    /**
     * Decodes a ThreadJoin message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns ThreadJoin
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): ThreadJoin;

    /**
     * Verifies a ThreadJoin message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a ThreadJoin message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns ThreadJoin
     */
    public static fromObject(object: { [k: string]: any }): ThreadJoin;

    /**
     * Creates a plain object from a ThreadJoin message. Also converts values to other types if specified.
     * @param message ThreadJoin
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: ThreadJoin, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this ThreadJoin to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a ThreadAnnounce. */
export interface IThreadAnnounce {

    /** ThreadAnnounce contact */
    contact: IContact;
}

/** Represents a ThreadAnnounce. */
export class ThreadAnnounce implements IThreadAnnounce {

    /**
     * Constructs a new ThreadAnnounce.
     * @param [properties] Properties to set
     */
    constructor(properties?: IThreadAnnounce);

    /** ThreadAnnounce contact. */
    public contact: IContact;

    /**
     * Creates a new ThreadAnnounce instance using the specified properties.
     * @param [properties] Properties to set
     * @returns ThreadAnnounce instance
     */
    public static create(properties?: IThreadAnnounce): ThreadAnnounce;

    /**
     * Encodes the specified ThreadAnnounce message. Does not implicitly {@link ThreadAnnounce.verify|verify} messages.
     * @param message ThreadAnnounce message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IThreadAnnounce, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified ThreadAnnounce message, length delimited. Does not implicitly {@link ThreadAnnounce.verify|verify} messages.
     * @param message ThreadAnnounce message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IThreadAnnounce, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a ThreadAnnounce message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns ThreadAnnounce
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): ThreadAnnounce;

    /**
     * Decodes a ThreadAnnounce message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns ThreadAnnounce
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): ThreadAnnounce;

    /**
     * Verifies a ThreadAnnounce message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a ThreadAnnounce message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns ThreadAnnounce
     */
    public static fromObject(object: { [k: string]: any }): ThreadAnnounce;

    /**
     * Creates a plain object from a ThreadAnnounce message. Also converts values to other types if specified.
     * @param message ThreadAnnounce
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: ThreadAnnounce, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this ThreadAnnounce to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a ThreadMessage. */
export interface IThreadMessage {

    /** ThreadMessage body */
    body: string;
}

/** Represents a ThreadMessage. */
export class ThreadMessage implements IThreadMessage {

    /**
     * Constructs a new ThreadMessage.
     * @param [properties] Properties to set
     */
    constructor(properties?: IThreadMessage);

    /** ThreadMessage body. */
    public body: string;

    /**
     * Creates a new ThreadMessage instance using the specified properties.
     * @param [properties] Properties to set
     * @returns ThreadMessage instance
     */
    public static create(properties?: IThreadMessage): ThreadMessage;

    /**
     * Encodes the specified ThreadMessage message. Does not implicitly {@link ThreadMessage.verify|verify} messages.
     * @param message ThreadMessage message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IThreadMessage, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified ThreadMessage message, length delimited. Does not implicitly {@link ThreadMessage.verify|verify} messages.
     * @param message ThreadMessage message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IThreadMessage, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a ThreadMessage message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns ThreadMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): ThreadMessage;

    /**
     * Decodes a ThreadMessage message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns ThreadMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): ThreadMessage;

    /**
     * Verifies a ThreadMessage message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a ThreadMessage message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns ThreadMessage
     */
    public static fromObject(object: { [k: string]: any }): ThreadMessage;

    /**
     * Creates a plain object from a ThreadMessage message. Also converts values to other types if specified.
     * @param message ThreadMessage
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: ThreadMessage, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this ThreadMessage to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a ThreadFiles. */
export interface IThreadFiles {

    /** ThreadFiles target */
    target: string;

    /** ThreadFiles body */
    body: string;

    /** ThreadFiles keys */
    keys: { [k: string]: string };
}

/** Represents a ThreadFiles. */
export class ThreadFiles implements IThreadFiles {

    /**
     * Constructs a new ThreadFiles.
     * @param [properties] Properties to set
     */
    constructor(properties?: IThreadFiles);

    /** ThreadFiles target. */
    public target: string;

    /** ThreadFiles body. */
    public body: string;

    /** ThreadFiles keys. */
    public keys: { [k: string]: string };

    /**
     * Creates a new ThreadFiles instance using the specified properties.
     * @param [properties] Properties to set
     * @returns ThreadFiles instance
     */
    public static create(properties?: IThreadFiles): ThreadFiles;

    /**
     * Encodes the specified ThreadFiles message. Does not implicitly {@link ThreadFiles.verify|verify} messages.
     * @param message ThreadFiles message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IThreadFiles, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified ThreadFiles message, length delimited. Does not implicitly {@link ThreadFiles.verify|verify} messages.
     * @param message ThreadFiles message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IThreadFiles, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a ThreadFiles message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns ThreadFiles
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): ThreadFiles;

    /**
     * Decodes a ThreadFiles message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns ThreadFiles
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): ThreadFiles;

    /**
     * Verifies a ThreadFiles message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a ThreadFiles message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns ThreadFiles
     */
    public static fromObject(object: { [k: string]: any }): ThreadFiles;

    /**
     * Creates a plain object from a ThreadFiles message. Also converts values to other types if specified.
     * @param message ThreadFiles
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: ThreadFiles, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this ThreadFiles to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a ThreadComment. */
export interface IThreadComment {

    /** ThreadComment target */
    target: string;

    /** ThreadComment body */
    body: string;
}

/** Represents a ThreadComment. */
export class ThreadComment implements IThreadComment {

    /**
     * Constructs a new ThreadComment.
     * @param [properties] Properties to set
     */
    constructor(properties?: IThreadComment);

    /** ThreadComment target. */
    public target: string;

    /** ThreadComment body. */
    public body: string;

    /**
     * Creates a new ThreadComment instance using the specified properties.
     * @param [properties] Properties to set
     * @returns ThreadComment instance
     */
    public static create(properties?: IThreadComment): ThreadComment;

    /**
     * Encodes the specified ThreadComment message. Does not implicitly {@link ThreadComment.verify|verify} messages.
     * @param message ThreadComment message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IThreadComment, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified ThreadComment message, length delimited. Does not implicitly {@link ThreadComment.verify|verify} messages.
     * @param message ThreadComment message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IThreadComment, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a ThreadComment message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns ThreadComment
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): ThreadComment;

    /**
     * Decodes a ThreadComment message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns ThreadComment
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): ThreadComment;

    /**
     * Verifies a ThreadComment message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a ThreadComment message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns ThreadComment
     */
    public static fromObject(object: { [k: string]: any }): ThreadComment;

    /**
     * Creates a plain object from a ThreadComment message. Also converts values to other types if specified.
     * @param message ThreadComment
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: ThreadComment, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this ThreadComment to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}

/** Properties of a ThreadLike. */
export interface IThreadLike {

    /** ThreadLike target */
    target: string;
}

/** Represents a ThreadLike. */
export class ThreadLike implements IThreadLike {

    /**
     * Constructs a new ThreadLike.
     * @param [properties] Properties to set
     */
    constructor(properties?: IThreadLike);

    /** ThreadLike target. */
    public target: string;

    /**
     * Creates a new ThreadLike instance using the specified properties.
     * @param [properties] Properties to set
     * @returns ThreadLike instance
     */
    public static create(properties?: IThreadLike): ThreadLike;

    /**
     * Encodes the specified ThreadLike message. Does not implicitly {@link ThreadLike.verify|verify} messages.
     * @param message ThreadLike message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IThreadLike, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified ThreadLike message, length delimited. Does not implicitly {@link ThreadLike.verify|verify} messages.
     * @param message ThreadLike message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IThreadLike, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a ThreadLike message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns ThreadLike
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): ThreadLike;

    /**
     * Decodes a ThreadLike message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns ThreadLike
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): ThreadLike;

    /**
     * Verifies a ThreadLike message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a ThreadLike message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns ThreadLike
     */
    public static fromObject(object: { [k: string]: any }): ThreadLike;

    /**
     * Creates a plain object from a ThreadLike message. Also converts values to other types if specified.
     * @param message ThreadLike
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: ThreadLike, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this ThreadLike to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };
}
