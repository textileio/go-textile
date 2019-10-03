# Textile REST API
Textile's HTTP REST API Documentation

## Version: 0

### Terms of service
https://github.com/textileio/go-textile/blob/master/TERMS

**Contact information:**  
Textile  
https://textile.io/  
contact@textile.io  

**License:** [MIT License](https://github.com/textileio/go-textile/blob/master/LICENSE)

### Security
**BasicAuth**  

|basic|*Basic*|
|---|---|

### /account

#### GET
##### Summary:

Show account contact

##### Description:

Shows the local peer's account info as a contact

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | contact | [pb.Contact](#pb.contact) |
| 400 | Bad Request | string |

### /account/address

#### GET
##### Summary:

Show account address

##### Description:

Shows the local peer's account address

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | address | string |

### /account/seed

#### GET
##### Summary:

Show account seed

##### Description:

Shows the local peer's account seed

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | seed | string |

### /blocks

#### GET
##### Summary:

Paginates blocks in a thread

##### Description:

Paginates blocks in a thread. Blocks are the raw components in a thread.
Think of them as an append-only log of thread updates where each update is
hash-linked to its parent(s). New / recovering peers can sync history by simply
traversing the hash tree.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| X-Textile-Opts | header | thread: Thread ID, offset: Offset ID to start listing from (omit for latest), limit: List page size (default: 5) | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | blocks | [pb.BlockList](#pb.blocklist) |
| 400 | Bad Request | string |
| 404 | Not Found | string |
| 500 | Internal Server Error | string |

### /blocks/{id}

#### DELETE
##### Summary:

Remove thread block

##### Description:

Removes a thread block by ID

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | block id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | block | [pb.Block](#pb.block) |
| 400 | Bad Request | string |
| 404 | Not Found | string |
| 500 | Internal Server Error | string |

### /blocks/{id}/comment

#### GET
##### Summary:

Get thread comment

##### Description:

Gets a thread comment by block ID

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | block id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | comment | [pb.Comment](#pb.comment) |
| 400 | Bad Request | string |

### /blocks/{id}/comments

#### GET
##### Summary:

List comments

##### Description:

Lists comments on a thread block

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | block id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | comments | [pb.CommentList](#pb.commentlist) |
| 500 | Internal Server Error | string |

#### POST
##### Summary:

Add a comment

##### Description:

Adds a comment to a thread block

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | block id | Yes | string |
| X-Textile-Args | header | urlescaped comment body | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | comment | [pb.Comment](#pb.comment) |
| 400 | Bad Request | string |
| 404 | Not Found | string |
| 500 | Internal Server Error | string |

### /blocks/{id}/files

#### GET
##### Summary:

Gets the metadata for a files block

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | block id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | files | [pb.Files](#pb.files) |
| 404 | Not Found | string |

### /blocks/{id}/files/{index}/{path}/content

#### GET
##### Summary:

Gets the decrypted file content of a file within a files block

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | block id | Yes | string |
| index | path | file index | Yes | string |
| path | path | file path | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | string |
| 400 | Bad Request | string |
| 404 | Not Found | string |

### /blocks/{id}/files/{index}/{path}/meta

#### GET
##### Summary:

Gets the metadata of a file within a files block

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | block id | Yes | string |
| index | path | file index | Yes | string |
| path | path | file path | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | file | [pb.FileIndex](#pb.fileindex) |
| 400 | Bad Request | string |
| 404 | Not Found | string |

### /blocks/{id}/like

#### GET
##### Summary:

Get thread like

##### Description:

Gets a thread like by block ID

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | block id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | like | [pb.Like](#pb.like) |
| 400 | Bad Request | string |

### /blocks/{id}/likes

#### GET
##### Summary:

List likes

##### Description:

Lists likes on a thread block

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | block id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | likes | [pb.LikeList](#pb.likelist) |
| 500 | Internal Server Error | string |

#### POST
##### Summary:

Add a like

##### Description:

Adds a like to a thread block

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | block id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | like | [pb.Like](#pb.like) |
| 400 | Bad Request | string |
| 404 | Not Found | string |
| 500 | Internal Server Error | string |

### /blocks/{id}/meta

#### GET
##### Summary:

Gets the metadata for a block

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | block id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | block | [pb.Block](#pb.block) |
| 404 | Not Found | string |

### /cafes

#### GET
##### Summary:

List info about all active cafe sessions

##### Description:

List info about all active cafe sessions. Cafes are other peers on the network
who offer pinning, backup, and inbox services

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | cafe sessions | [pb.CafeSessionList](#pb.cafesessionlist) |
| 500 | Internal Server Error | string |

#### POST
##### Summary:

Register with a Cafe

##### Description:

Registers with a cafe and saves an expiring service session token. An access
token is required to register, and should be obtained separately from the target
Cafe

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| X-Textile-Args | header | cafe id | Yes | string |
| X-Textile-Opts | header | token: An access token supplied by the Cafe | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | cafe session | [pb.CafeSession](#pb.cafesession) |
| 400 | Bad Request | string |
| 500 | Internal Server Error | string |

### /cafes/{id}

#### DELETE
##### Summary:

Deregisters a cafe

##### Description:

Deregisters with a cafe (content will expire based on the cafe's service rules)

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | cafe id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | ok | string |
| 500 | Internal Server Error | string |

#### GET
##### Summary:

Gets and displays info about a cafe session

##### Description:

Gets and displays info about a cafe session. Cafes are other peers on the network
who offer pinning, backup, and inbox services

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | cafe id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | cafe session | [pb.CafeSession](#pb.cafesession) |
| 404 | Not Found | string |
| 500 | Internal Server Error | string |

### /cafes/messages

#### POST
##### Summary:

Check for messages at all cafes

##### Description:

Check for messages at all cafes. New messages are downloaded and processed
opportunistically.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | ok | string |
| 500 | Internal Server Error | string |

### /config

#### PATCH
##### Summary:

Set/update config settings

##### Description:

When patching config values, valid JSON types must be used. For example, a string
should be escaped or wrapped in single quotes (e.g., \"127.0.0.1:40600\") and
arrays and objects work fine (e.g. '{"API": "127.0.0.1:40600"}') but should be
wrapped in single quotes. Be sure to restart the daemon for changes to take effect.
See https://tools.ietf.org/html/rfc6902 for details on RFC6902 JSON patch format.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| patch | body | An RFC6902 JSON patch (array of ops) | Yes | [mill.Json](#mill.json) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | No Content | string |
| 400 | Bad Request | string |

#### PUT
##### Summary:

Replace config settings.

##### Description:

Replace entire config file contents. The config command controls configuration
variables. It works much like 'git config'. The configuration values are stored
in a config file inside the Textile repository.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| config | body | JSON document | Yes | [mill.Json](#mill.json) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | No Content | string |
| 400 | Bad Request | string |

### /config/{path}

#### GET
##### Summary:

Get active config settings

##### Description:

Report the currently active config settings, which may differ from the values
specifed when setting/patching values.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| path | path | config path (e.g., Addresses/API) | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | new config value | [mill.Json](#mill.json) |
| 400 | Bad Request | string |

### /contacts

#### GET
##### Summary:

List known contacts

##### Description:

Lists known contacts.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | contacts | [pb.ContactList](#pb.contactlist) |
| 404 | Not Found | string |
| 500 | Internal Server Error | string |

### /contacts/{address}

#### DELETE
##### Summary:

Remove a contact

##### Description:

Removes a known contact

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| address | path | address | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | ok | string |
| 404 | Not Found | string |
| 500 | Internal Server Error | string |

#### GET
##### Summary:

Get a known contact

##### Description:

Gets a known contact

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| address | path | address | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | contact | [pb.Contact](#pb.contact) |
| 404 | Not Found | string |

#### PUT
##### Summary:

Add to known contacts

##### Description:

Adds a contact by username or account address to known contacts.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| address | path | address | Yes | string |
| contact | body | contact | Yes | [pb.Contact](#pb.contact) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | ok | string |
| 400 | Bad Request | string |

### /contacts/search

#### POST
##### Summary:

Search for contacts

##### Description:

Search for contacts known locally and on the network

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| X-Textile-Opts | header | local: Whether to only search local contacts, remote: Whether to only search remote contacts, limit: Stops searching after limit results are found, wait: Stops searching after 'wait' seconds have elapsed (max 30s), username: search by username string, address: search by account address string, events: Whether to emit Server-Sent Events (SSEvent) or plain JSON | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | results stream | [pb.QueryResult](#pb.queryresult) |
| 404 | Not Found | string |
| 500 | Internal Server Error | string |

### /feed

#### GET
##### Summary:

Paginates post and annotation block types

##### Description:

Paginates post (join|leave|files|message) and annotation (comment|like) block types
The mode option dictates how the feed is displayed:
"chrono": All feed block types are shown. Annotations always nest their target post,
i.e., the post a comment is about.
"annotated": Annotations are nested under post targets, but are not shown in the
top-level feed.
"stacks": Related blocks are chronologically grouped into "stacks". A new stack is
started if an unrelated block breaks continuity. This mode is used by Textile
Photos. Stacks may include:
* The initial post with some nested annotations. Newer annotations may have already
been listed.
* One or more annotations about a post. The newest annotation assumes the "top"
position in the stack. Additional annotations are nested under the target.
Newer annotations may have already been listed in the case as well.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| X-Textile-Opts | header | thread: Thread ID (can also use 'default'), offset: Offset ID to start listing from (omit for latest), limit: List page size (default: 5), mode: Feed mode (one of 'chrono', 'annotated', or 'stacks') | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | feed | [pb.FeedItemList](#pb.feeditemlist) |
| 400 | Bad Request | string |
| 404 | Not Found | string |
| 500 | Internal Server Error | string |

### /file/{hash}/content

#### GET
##### Summary:

File content at hash

##### Description:

Returns decrypted raw content for file

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hash | path | file hash | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | string |
| 404 | Not Found | string |

### /file/{target}/meta

#### GET
##### Summary:

File metadata at hash

##### Description:

Returns the metadata for file

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hash | path | file hash | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | file | [pb.FileIndex](#pb.fileindex) |
| 404 | Not Found | string |

### /files

#### GET
##### Summary:

Paginates thread files

##### Description:

Paginates thread files. If thread id not provided, paginate all files.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| X-Textile-Opts | header | thread: Thread ID. Omit for all, offset: Offset ID to start listing from. Omit for latest, limit: List page size. (default: 5) | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | files | [pb.FilesList](#pb.fileslist) |
| 400 | Bad Request | string |
| 404 | Not Found | string |
| 500 | Internal Server Error | string |

### /invites

#### GET
##### Summary:

List invites

##### Description:

Lists all pending thread invites

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | invites | [pb.InviteViewList](#pb.inviteviewlist) |

#### POST
##### Summary:

Create an invite to a thread

##### Description:

Creates a direct account-to-account or external invite to a thread

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| X-Textile-Opts | header | thread: Thread ID (can also use 'default'), address: Account Address (omit to create an external invite) | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | invite | [pb.ExternalInvite](#pb.externalinvite) |
| 400 | Bad Request | string |
| 404 | Not Found | string |
| 500 | Internal Server Error | string |

### /invites/{id}/accept

#### POST
##### Summary:

Accept a thread invite

##### Description:

Accepts a direct peer-to-peer or external invite to a thread. Use the key option
with an external invite

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | invite id | Yes | string |
| X-Textile-Opts | header | key: key for an external invite | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | join block | [pb.Block](#pb.block) |
| 400 | Bad Request | string |
| 409 | Conflict | string |
| 500 | Internal Server Error | string |

### /invites/{id}/ignore

#### POST
##### Summary:

Ignore a thread invite

##### Description:

Ignores a direct peer-to-peer invite to a thread

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | invite id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | ok | string |
| 400 | Bad Request | string |

### /ipfs/cat/{path}

#### GET
##### Summary:

Cat IPFS data

##### Description:

Displays the data behind an IPFS CID (hash) or Path

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| path | path | ipfs/ipns cid | Yes | string |
| X-Textile-Opts | header | key: Key to decrypt data on-the-fly | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | data | [ integer ] |
| 400 | Bad Request | string |
| 401 | Unauthorized | string |
| 404 | Not Found | string |
| 500 | Internal Server Error | string |

### /ipfs/id

#### GET
##### Summary:

Get IPFS peer ID

##### Description:

Displays underlying IPFS peer ID

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | peer id | string |
| 500 | Internal Server Error | string |

### /ipfs/swarm/connect

#### POST
##### Summary:

Opens a new direct connection to a peer address

##### Description:

Opens a new direct connection to a peer using an IPFS multiaddr

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| X-Textile-Args | header | peer address | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | ok | [ string ] |
| 400 | Bad Request | string |
| 500 | Internal Server Error | string |

### /ipfs/swarm/peers

#### GET
##### Summary:

List swarm peers

##### Description:

Lists the set of peers this node is connected to

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| X-Textile-Opts | header | verbose: Display all extra information, latency: Also list information about latency to each peer, streams: Also list information about open streams for each peer, direction: Also list information about the direction of connection | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | connection | [ipfs.ConnInfos](#ipfs.conninfos) |
| 400 | Bad Request | string |
| 500 | Internal Server Error | string |

### /keys/{target}

#### GET
##### Summary:

Show file keys

##### Description:

Shows file keys under the given target from an add

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| target | path | target id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | keys | [pb.Keys](#pb.keys) |
| 400 | Bad Request | string |

### /logs/{subsystem}

#### POST
##### Summary:

Access subsystem logs

##### Description:

List or change the verbosity of one or all subsystems log output. Textile logs
piggyback on the IPFS event logs

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| subsystem | path | subsystem logging identifier (omit for all) | No | string |
| X-Textile-Opts | header | level: Log-level (one of: debug, info, warning, error, critical, or  | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | subsystems | [core.SubsystemInfo](#core.subsysteminfo) |
| 400 | Bad Request | string |
| 500 | Internal Server Error | string |

### /messages

#### GET
##### Summary:

Paginates thread messages

##### Description:

Paginates thread messages

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| X-Textile-Opts | header | thread: Thread ID (can also use 'default', omit for all), offset: Offset ID to start listing from (omit for latest), limit: List page size (default: 5) | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | messages | [pb.TextList](#pb.textlist) |
| 400 | Bad Request | string |
| 404 | Not Found | string |
| 500 | Internal Server Error | string |

### /messages/{block}

#### GET
##### Summary:

Get thread message

##### Description:

Gets a thread message by block ID

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| block | path | block id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | message | [pb.Text](#pb.text) |
| 400 | Bad Request | string |

### /mills/blob

#### POST
##### Summary:

Process raw data blobs

##### Description:

Takes a binary data blob, and optionally encrypts it, before adding to IPFS,
and returns a file object

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| file | formData | multipart/form-data file | No | file |
| X-Textile-Opts | header | plaintext: whether to leave unencrypted), use: if empty, assumes body contains multipart form file data, otherwise, will attempt to fetch given CID from IPFS | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | file | [pb.FileIndex](#pb.fileindex) |
| 400 | Bad Request | string |
| 500 | Internal Server Error | string |

### /mills/image/exif

#### POST
##### Summary:

Extract EXIF data from image

##### Description:

Takes an input image, and extracts its EXIF data (optionally encrypting output),
before adding to IPFS, and returns a file object

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| file | formData | multipart/form-data file | No | file |
| X-Textile-Opts | header | plaintext: whether to leave unencrypted, use: if empty, assumes body contains multipart form file data, otherwise, will attempt to fetch given CID from IPFS | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | file | [pb.FileIndex](#pb.fileindex) |
| 400 | Bad Request | string |
| 500 | Internal Server Error | string |

### /mills/image/resize

#### POST
##### Summary:

Resize an image

##### Description:

Takes an input image, and resizes/resamples it (optionally encrypting output),
before adding to IPFS, and returns a file object

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| file | formData | multipart/form-data file | No | file |
| X-Textile-Opts | header | plaintext: whether to leave unencrypted, use: if empty, assumes body contains multipart form file data, otherwise, will attempt to fetch given CID from IPFS, width: the requested image width (required), quality: the requested JPEG image quality | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | file | [pb.FileIndex](#pb.fileindex) |
| 400 | Bad Request | string |
| 500 | Internal Server Error | string |

### /mills/json

#### POST
##### Summary:

Process input JSON data

##### Description:

Takes an input JSON document, validates it according to its json-schema.org definition,
optionally encrypts the output before adding to IPFS, and returns a file object

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| file | formData | multipart/form-data file | No | file |
| X-Textile-Opts | header | plaintext: whether to leave unencrypted, use: if empty, assumes body contains multipart form file data, otherwise, will attempt to fetch given CID from IPFS | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | file | [pb.FileIndex](#pb.fileindex) |
| 400 | Bad Request | string |
| 500 | Internal Server Error | string |

### /mills/schema

#### POST
##### Summary:

Validate, add, and pin a new Schema

##### Description:

Takes a JSON-based Schema, validates it, adds it to IPFS, and returns a file object

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| schema | body | schema | Yes | [pb.Node](#pb.node) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | file | [pb.FileIndex](#pb.fileindex) |
| 400 | Bad Request | string |

### /notifications

#### GET
##### Summary:

List notifications

##### Description:

Lists all notifications generated by thread and account activity.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | notifications | [pb.NotificationList](#pb.notificationlist) |

### /notifications/{id}/read

#### POST
##### Summary:

Mark notifiction as read

##### Description:

Marks a notifiction as read by ID. Use 'all' to mark all as read.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | notification id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | ok | string |
| 400 | Bad Request | string |

### /ping

#### GET
##### Summary:

Ping a network peer

##### Description:

Pings another peer on the network, returning online|offline.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| X-Textile-Args | header | peerid | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | One of online|offline | string |
| 400 | Bad Request | string |
| 500 | Internal Server Error | string |

### /profile

#### GET
##### Summary:

Get public profile

##### Description:

Gets the local node's public profile

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | peer | [pb.Peer](#pb.peer) |
| 400 | Bad Request | string |

### /profile/avatar

#### POST
##### Summary:

Set avatar

##### Description:

Forces local node to update avatar image to latest image added to 'account' thread

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | ok | string |
| 400 | Bad Request | string |
| 500 | Internal Server Error | string |

### /profile/name

#### POST
##### Summary:

Set display name

##### Description:

Sets public profile display name to given string

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| X-Textile-Args | header | name | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | ok | string |
| 400 | Bad Request | string |
| 500 | Internal Server Error | string |

### /publish

#### POST
##### Summary:

Publish payload to topic

##### Description:

Publishes payload bytes to a topic on the network.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| X-Textile-Args | header | topic | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | ok | string |
| 500 | Internal Server Error | string |

### /snapshots

#### POST
##### Summary:

Create thread snapshots

##### Description:

Snapshots all threads and pushes to registered cafes

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | ok | string |
| 400 | Bad Request | string |
| 500 | Internal Server Error | string |

### /snapshots/search

#### POST
##### Summary:

Search for thread snapshots

##### Description:

Searches the network for thread snapshots

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| X-Textile-Opts | header | wait: Stops searching after 'wait' seconds have elapsed (max 30s), events: Whether to emit Server-Sent Events (SSEvent) or plain JSON | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | results stream | [pb.QueryResult](#pb.queryresult) |
| 400 | Bad Request | string |
| 500 | Internal Server Error | string |

### /subscribe/{id}

#### GET
##### Summary:

Observe to thread updates

##### Description:

Observes updates in a thread or all threads. An update is generated
when a new block is added to a thread. There are several update types:
MERGE, IGNORE, FLAG, JOIN, ANNOUNCE, LEAVE, TEXT, FILES, COMMENT, LIKE

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| thread | path | thread id, omit to stream all events | No | string |
| X-Textile-Opts | header | type: Or'd list of event types (e.g., FILES|COMMENTS|LIKES) or empty to include all types, events: Whether to emit Server-Sent Events (SSEvent) or plain JSON | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | stream of updates | [pb.FeedItem](#pb.feeditem) |
| 500 | Internal Server Error | string |

### /summary

#### GET
##### Summary:

Get a summary of node data

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | summary | [pb.Summary](#pb.summary) |

### /threads

#### GET
##### Summary:

Lists info on all threads

##### Description:

Lists all local threads, returning a ThreadList object

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | threads | [pb.ThreadList](#pb.threadlist) |
| 400 | Bad Request | string |
| 500 | Internal Server Error | string |

#### POST
##### Summary:

Adds and joins a new thread

##### Description:

Adds a new Thread with given name, type, and sharing and whitelist options, returning
a Thread object

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| X-Textile-Args | header | name | Yes | string |
| X-Textile-Opts | header | key: A locally unique key used by an app to identify this thread on recovery, schema: Existing Thread Schema IPFS CID, type: Set the thread type to one of 'private', 'read_only', 'public', or 'open', sharing: Set the thread sharing style to one of 'not_shared','invite_only', or 'shared', whitelist: An array of contact addresses. When supplied, the thread will not allow additional peers beyond those in array, useful for 1-1 chat/file sharing | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | thread | [pb.Thread](#pb.thread) |
| 400 | Bad Request | string |
| 500 | Internal Server Error | string |

### /threads/{id}

#### DELETE
##### Summary:

Abandons a thread.

##### Description:

Abandons a thread, and if no one else is participating, then the thread dissipates.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | thread id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | ok | string |
| 404 | Not Found | string |
| 500 | Internal Server Error | string |

#### GET
##### Summary:

Gets a thread

##### Description:

Gets and displays info about a thread

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | thread id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | thread | [pb.Thread](#pb.thread) |
| 400 | Bad Request | string |
| 500 | Internal Server Error | string |

#### PUT
##### Summary:

Add or update a thread directly

##### Description:

Adds or updates a thread directly, usually from a backup

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | id | Yes | string |
| thread | body | thread | Yes | [pb.Thread](#pb.thread) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | ok | string |
| 400 | Bad Request | string |

### /threads/{id}/files

#### POST
##### Summary:

Adds a file or directory of files to a thread

##### Description:

Adds a file or directory of files to a thread. Files not supported by the thread
schema are ignored. Nested directories are included. An existing file hash may
also be used as input.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| dir | body | list of milled dirs (output from mill endpoint) | Yes | [pb.DirectoryList](#pb.directorylist) |
| X-Textile-Opts | header | caption: Caption to add to file(s) | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | file | [pb.Files](#pb.files) |
| 400 | Bad Request | string |
| 404 | Not Found | string |
| 500 | Internal Server Error | string |

### /threads/{id}/messages

#### POST
##### Summary:

Add a message

##### Description:

Adds a message to a thread

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| X-Textile-Args | header | urlescaped message body | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | message | [pb.Text](#pb.text) |
| 400 | Bad Request | string |
| 404 | Not Found | string |
| 500 | Internal Server Error | string |

### /threads/{id}/name

#### PUT
##### Summary:

Rename a thread

##### Description:

Renames a thread. Only initiators can rename a thread.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | id | Yes | string |
| X-Textile-Args | header | name | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | ok | string |
| 400 | Bad Request | string |

### /threads/{id}/peers

#### GET
##### Summary:

List all thread peers

##### Description:

Lists all peers in a thread

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | thread id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | contacts | [pb.ContactList](#pb.contactlist) |
| 404 | Not Found | string |

### /tokens

#### GET
##### Summary:

List local tokens

##### Description:

List info about all stored cafe tokens

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | tokens | [ string ] |
| 500 | Internal Server Error | string |

#### POST
##### Summary:

Create an access token

##### Description:

Generates an access token (44 random bytes) and saves a bcrypt hashed version for
future lookup. The response contains a base58 encoded version of the random bytes
token. If the 'store' option is set to false, the token is generated, but not
stored in the local Cafe db. Alternatively, an existing token can be added using
by specifying the 'token' option.
Tokens allow other peers to register with a Cafe peer.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| X-Textile-Opts | header | token: Use existing token, rather than creating a new one, store: Whether to store the added/generated token to the local db | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | token | string |
| 400 | Bad Request | string |
| 500 | Internal Server Error | string |

### /tokens/{id}

#### DELETE
##### Summary:

Removes a cafe token

##### Description:

Removes an existing cafe token

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| token | path | token | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 | ok | string |
| 500 | Internal Server Error | string |

#### GET
##### Summary:

Check token validity

##### Description:

Check validity of existing cafe access token

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| token | path | invite id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | ok | string |
| 401 | Unauthorized | string |

### Models


#### core.SubsystemInfo

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| core.SubsystemInfo | object |  |  |

#### ipfs.ConnInfos

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| peers | [ [ipfs.connInfo](#ipfs.conninfo) ] |  | No |

#### ipfs.connInfo

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| addr | string |  | No |
| direction | string |  | No |
| latency | string |  | No |
| muxer | string |  | No |
| peer | string |  | No |
| streams | [ [ipfs.streamInfo](#ipfs.streaminfo) ] |  | No |

#### ipfs.streamInfo

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| protocol | string |  | No |

#### mill.Json

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| mill.Json | object |  |  |

#### pb.Block

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| attempts | integer |  | No |
| author | string |  | No |
| body | string |  | No |
| data | string |  | No |
| date | string |  | No |
| id | string |  | No |
| parents | [ string ] |  | No |
| status | integer |  | No |
| target | string |  | No |
| thread | string |  | No |
| type | integer |  | No |
| user | [pb.User](#pb.user) | view info | No |

#### pb.BlockList

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| items | [ [pb.Block](#pb.block) ] |  | No |

#### pb.Cafe

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| address | string |  | No |
| api | string |  | No |
| node | string |  | No |
| peer | string |  | No |
| protocol | string |  | No |
| url | string |  | No |

#### pb.CafeSession

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| access | string |  | No |
| cafe | [pb.Cafe](#pb.cafe) |  | No |
| exp | string |  | No |
| id | string |  | No |
| refresh | string |  | No |
| rexp | string |  | No |
| subject | string |  | No |
| type | string |  | No |

#### pb.CafeSessionList

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| items | [ [pb.CafeSession](#pb.cafesession) ] |  | No |

#### pb.Comment

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| body | string |  | No |
| date | string |  | No |
| id | string |  | No |
| target | [pb.FeedItem](#pb.feeditem) |  | No |
| user | [pb.User](#pb.user) |  | No |

#### pb.CommentList

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| items | [ [pb.Comment](#pb.comment) ] |  | No |

#### pb.Contact

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| address | string |  | No |
| avatar | string |  | No |
| name | string |  | No |
| peers | [ [pb.Peer](#pb.peer) ] |  | No |
| threads | [ string ] |  | No |

#### pb.ContactList

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| items | [ [pb.Contact](#pb.contact) ] |  | No |

#### pb.Directory

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| files | object |  | No |

#### pb.DirectoryList

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| items | [ [pb.Directory](#pb.directory) ] |  | No |

#### pb.ExternalInvite

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string |  | No |
| inviter | string |  | No |
| key | string |  | No |

#### pb.FeedItem

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| block | string |  | No |
| payload | string |  | No |
| thread | string |  | No |

#### pb.FeedItemList

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| count | integer |  | No |
| items | [ [pb.FeedItem](#pb.feeditem) ] |  | No |
| next | string |  | No |

#### pb.File

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| file | [pb.FileIndex](#pb.fileindex) |  | No |
| index | integer |  | No |
| links | object |  | No |

#### pb.FileIndex

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| added | string |  | No |
| checksum | string |  | No |
| hash | string |  | No |
| key | string |  | No |
| media | string |  | No |
| meta | string |  | No |
| mill | string |  | No |
| name | string |  | No |
| opts | string |  | No |
| size | integer |  | No |
| source | string |  | No |
| targets | [ string ] |  | No |

#### pb.Files

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| block | string |  | No |
| caption | string |  | No |
| comments | [ [pb.Comment](#pb.comment) ] |  | No |
| data | string |  | No |
| date | string |  | No |
| files | [ [pb.File](#pb.file) ] |  | No |
| likes | [ [pb.Like](#pb.like) ] |  | No |
| target | string |  | No |
| threads | [ string ] |  | No |
| user | [pb.User](#pb.user) |  | No |

#### pb.FilesList

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| items | [ [pb.Files](#pb.files) ] |  | No |

#### pb.InviteView

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| date | string |  | No |
| id | string |  | No |
| inviter | [pb.User](#pb.user) |  | No |
| name | string |  | No |

#### pb.InviteViewList

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| items | [ [pb.InviteView](#pb.inviteview) ] |  | No |

#### pb.Keys

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| files | object |  | No |

#### pb.Like

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| date | string |  | No |
| id | string |  | No |
| target | [pb.FeedItem](#pb.feeditem) |  | No |
| user | [pb.User](#pb.user) |  | No |

#### pb.LikeList

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| items | [ [pb.Like](#pb.like) ] |  | No |

#### pb.Node

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| json_schema | string |  | No |
| links | object |  | No |
| mill | string |  | No |
| name | string |  | No |
| opts | object |  | No |
| pin | boolean |  | No |
| plaintext | boolean |  | No |

#### pb.Notification

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| actor | string |  | No |
| block | string |  | No |
| body | string |  | No |
| date | string |  | No |
| id | string |  | No |
| read | boolean |  | No |
| subject | string |  | No |
| subject_desc | string |  | No |
| target | string |  | No |
| type | integer |  | No |
| user | [pb.User](#pb.user) | view info | No |

#### pb.NotificationList

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| items | [ [pb.Notification](#pb.notification) ] |  | No |

#### pb.Peer

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| address | string |  | No |
| avatar | string |  | No |
| created | string |  | No |
| id | string |  | No |
| inboxes | [ [pb.Cafe](#pb.cafe) ] |  | No |
| name | string |  | No |
| updated | string |  | No |

#### pb.QueryResult

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| date | string |  | No |
| id | string |  | No |
| local | boolean |  | No |
| value | string |  | No |

#### pb.Summary

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| account_peer_count | integer |  | No |
| address | string |  | No |
| contact_count | integer |  | No |
| files_count | integer |  | No |
| id | string |  | No |
| thread_count | integer |  | No |

#### pb.Text

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| block | string |  | No |
| body | string |  | No |
| comments | [ [pb.Comment](#pb.comment) ] |  | No |
| date | string |  | No |
| likes | [ [pb.Like](#pb.like) ] |  | No |
| user | [pb.User](#pb.user) |  | No |

#### pb.TextList

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| items | [ [pb.Text](#pb.text) ] |  | No |

#### pb.Thread

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| block_count | integer |  | No |
| head | string |  | No |
| head_blocks | [ [pb.Block](#pb.block) ] | view info | No |
| id | string |  | No |
| initiator | string |  | No |
| key | string |  | No |
| name | string |  | No |
| peer_count | integer |  | No |
| schema | string |  | No |
| schema_node | [pb.Node](#pb.node) |  | No |
| sharing | integer |  | No |
| sk | [ integer ] |  | No |
| state | integer | Deprecated: Do not use. | No |
| type | integer |  | No |
| whitelist | [ string ] |  | No |

#### pb.ThreadList

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| items | [ [pb.Thread](#pb.thread) ] |  | No |

#### pb.User

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| address | string |  | No |
| avatar | string |  | No |
| name | string |  | No |