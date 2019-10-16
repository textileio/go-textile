package db

import (
	"database/sql"
	"path"
	"strings"
	"sync"

	"github.com/golang/protobuf/jsonpb"
	logging "github.com/ipfs/go-log"
	_ "github.com/mutecomm/go-sqlcipher"
	"github.com/textileio/go-textile/repo"
)

var log = logging.Logger("tex-datastore")

var pbMarshaler = jsonpb.Marshaler{
	OrigName: true,
}

var pbUnmarshaler = jsonpb.Unmarshaler{
	AllowUnknownFields: true,
}

type SQLiteDatastore struct {
	config             repo.ConfigStore
	peers              repo.PeerStore
	files              repo.FileStore
	threads            repo.ThreadStore
	threadPeers        repo.ThreadPeerStore
	blocks             repo.BlockStore
	blockMessages      repo.BlockMessageStore
	invites            repo.InviteStore
	notifications      repo.NotificationStore
	cafeSessions       repo.CafeSessionStore
	cafeRequests       repo.CafeRequestStore
	cafeMessages       repo.CafeMessageStore
	cafeClientNonces   repo.CafeClientNonceStore
	cafeClients        repo.CafeClientStore
	cafeTokens         repo.CafeTokenStore
	cafeClientThreads  repo.CafeClientThreadStore
	cafeClientMessages repo.CafeClientMessageStore
	botsStore          repo.Botstore
	db                 *sql.DB
	lock               *sync.Mutex
}

func Create(repoPath, pin string) (*SQLiteDatastore, error) {
	dbPath := path.Join(repoPath, "datastore", "mainnet.db")
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	conn.SetMaxIdleConns(2)
	conn.SetMaxOpenConns(4)
	if pin != "" {
		p := "pragma key='" + strings.Replace(pin, "'", "''", -1) + "';"
		if _, err := conn.Exec(p); err != nil {
			return nil, err
		}
	}
	lock := new(sync.Mutex)
	return &SQLiteDatastore{
		config:             NewConfigStore(conn, lock, dbPath),
		peers:              NewPeerStore(conn, lock),
		files:              NewFileStore(conn, lock),
		threads:            NewThreadStore(conn, lock),
		threadPeers:        NewThreadPeerStore(conn, lock),
		blocks:             NewBlockStore(conn, lock),
		blockMessages:      NewBlockMessageStore(conn, lock),
		invites:            NewInviteStore(conn, lock),
		notifications:      NewNotificationStore(conn, lock),
		cafeSessions:       NewCafeSessionStore(conn, lock),
		cafeRequests:       NewCafeRequestStore(conn, lock),
		cafeMessages:       NewCafeMessageStore(conn, lock),
		cafeClientNonces:   NewCafeClientNonceStore(conn, lock),
		cafeClients:        NewCafeClientStore(conn, lock),
		cafeTokens:         NewCafeTokenStore(conn, lock),
		cafeClientThreads:  NewCafeClientThreadStore(conn, lock),
		cafeClientMessages: NewCafeClientMessageStore(conn, lock),
		botsStore:          NewBotstore(conn, lock),
		db:                 conn,
		lock:               lock,
	}, nil
}

func (d *SQLiteDatastore) Ping() error {
	return d.db.Ping()
}

func (d *SQLiteDatastore) Close() {
	_ = d.db.Close()
}

func (d *SQLiteDatastore) Config() repo.ConfigStore {
	return d.config
}

func (d *SQLiteDatastore) Peers() repo.PeerStore {
	return d.peers
}

func (d *SQLiteDatastore) Files() repo.FileStore {
	return d.files
}

func (d *SQLiteDatastore) Threads() repo.ThreadStore {
	return d.threads
}

func (d *SQLiteDatastore) ThreadPeers() repo.ThreadPeerStore {
	return d.threadPeers
}

func (d *SQLiteDatastore) Blocks() repo.BlockStore {
	return d.blocks
}

func (d *SQLiteDatastore) BlockMessages() repo.BlockMessageStore {
	return d.blockMessages
}

func (d *SQLiteDatastore) Invites() repo.InviteStore {
	return d.invites
}

func (d *SQLiteDatastore) Notifications() repo.NotificationStore {
	return d.notifications
}

func (d *SQLiteDatastore) CafeSessions() repo.CafeSessionStore {
	return d.cafeSessions
}

func (d *SQLiteDatastore) CafeRequests() repo.CafeRequestStore {
	return d.cafeRequests
}

func (d *SQLiteDatastore) CafeMessages() repo.CafeMessageStore {
	return d.cafeMessages
}

func (d *SQLiteDatastore) CafeClientNonces() repo.CafeClientNonceStore {
	return d.cafeClientNonces
}

func (d *SQLiteDatastore) CafeClients() repo.CafeClientStore {
	return d.cafeClients
}

func (d *SQLiteDatastore) CafeTokens() repo.CafeTokenStore {
	return d.cafeTokens
}

func (d *SQLiteDatastore) CafeClientThreads() repo.CafeClientThreadStore {
	return d.cafeClientThreads
}

func (d *SQLiteDatastore) CafeClientMessages() repo.CafeClientMessageStore {
	return d.cafeClientMessages
}

func (d *SQLiteDatastore) Bots() repo.Botstore {
	return d.botsStore
}

func (d *SQLiteDatastore) Copy(dbPath string, pin string) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	var cp string
	stmt := "select name from sqlite_master where type='table'"
	rows, err := d.db.Query(stmt)
	if err != nil {
		log.Errorf("error in copy: %s", err)
		return err
	}
	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		tables = append(tables, name)
	}
	if pin == "" {
		cp = `attach database '` + dbPath + `' as plaintext key '';`
		for _, name := range tables {
			cp = cp + "insert into plaintext." + name + " select * from main." + name + ";"
		}
	} else {
		cp = `attach database '` + dbPath + `' as encrypted key '` + pin + `';`
		for _, name := range tables {
			cp = cp + "insert into encrypted." + name + " select * from main." + name + ";"
		}
	}
	_, err = d.db.Exec(cp)
	if err != nil {
		return err
	}
	return nil
}

func (d *SQLiteDatastore) InitTables(pin string) error {
	return initDatabaseTables(d.db, pin)
}

func ConflictError(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}

func initDatabaseTables(db *sql.DB, pin string) error {
	var sqlStmt string
	if pin != "" {
		sqlStmt = "pragma key = '" + strings.Replace(pin, "'", "''", -1) + "';"
	}
	sqlStmt += `
    create table config (key text primary key not null, value blob);

    create table peers (id text primary key not null, address text not null, username text not null, avatar text not null, inboxes blob not null, created integer not null, updated integer not null);
    create index peer_address on peers (address);
    create index peer_username on peers (username);
    create index peer_updated on peers (updated);

    create table files (mill text not null, checksum text not null, source text not null, opts text not null, hash text not null, key text not null, media text not null, name text not null, size integer not null, added integer not null, meta blob, targets text, primary key (mill, checksum));
    create index file_hash on files (hash);
    create unique index file_mill_source_opts on files (mill, source, opts);

    create table threads (id text primary key not null, key text not null, sk blob not null, name text not null, schema text not null, initiator text not null, type integer not null, state integer not null, head text not null, members text not null, sharing integer not null);
    create unique index thread_key on threads (key);

    create table thread_peers (id text not null, threadId text not null, welcomed integer not null, primary key (id, threadId));
    create index thread_peer_id on thread_peers (id);
    create index thread_peer_threadId on thread_peers (threadId);
    create index thread_peer_welcomed on thread_peers (welcomed);

    create table blocks (id text primary key not null, threadId text not null, authorId text not null, type integer not null, date integer not null, parents text not null, target text not null, body text not null, data text not null, status integer not null, attempts integer not null);
    create index block_threadId on blocks (threadId);
    create index block_type on blocks (type);
    create index block_date on blocks (date);
    create index block_target on blocks (target);
    create index block_data on blocks (data);
    create index block_status on blocks (status);

    create table block_messages (id text primary key not null, peerId text not null, envelope blob not null, date integer not null);
    create index block_message_date on block_messages (date);

    create table invites (id text primary key not null, block blob not null, name text not null, inviter blob not null, date integer not null, parents text not null);
    create index invite_date on invites (date);

    create table notifications (id text primary key not null, date integer not null, actorId text not null, subject text not null, subjectId text not null, blockId text, target text, type integer not null, body text not null, read integer not null);
    create index notification_date on notifications (date);
    create index notification_actorId on notifications (actorId);
    create index notification_subjectId on notifications (subjectId);
    create index notification_blockId on notifications (blockId);
    create index notification_read on notifications (read);

    create table cafe_sessions (cafeId text primary key not null, access text not null, refresh text not null, expiry integer not null, cafe blob not null);

    create table cafe_requests (id text primary key not null, peerId text not null, targetId text not null, cafeId text not null, cafe blob not null, groupId text not null, syncGroupId text not null, type integer not null, date integer not null, size integer not null, status integer not null, attempts integer not null, groupSize integer not null, groupTransferred integer not null);
    create index cafe_request_cafeId on cafe_requests (cafeId);
    create index cafe_request_groupId on cafe_requests (groupId);
    create index cafe_request_syncGroupId on cafe_requests (syncGroupId);
    create index cafe_request_date on cafe_requests (date);
    create index cafe_request_status on cafe_requests (status);

    create table cafe_messages (id text primary key not null, peerId text not null, date integer not null, attempts integer not null);
    create index cafe_message_date on cafe_messages (date);

    create table cafe_client_nonces (value text primary key not null, address text not null, date integer not null);

    create table cafe_clients (id text primary key not null, address text not null, created integer not null, lastSeen integer not null, tokenId text not null);
    create index cafe_client_address on cafe_clients (address);
    create index cafe_client_lastSeen on cafe_clients (lastSeen);

    create table cafe_client_threads (id text not null, clientId text not null, ciphertext blob not null, primary key (id, clientId));
    create index cafe_client_thread_clientId on cafe_client_threads (clientId);

    create table cafe_client_messages (id text not null, peerId text not null, clientId text not null, date integer not null, primary key (id, clientId));
    create index cafe_client_message_clientId on cafe_client_messages (clientId);
    create index cafe_client_message_date on cafe_client_messages (date);

		create table cafe_tokens (id text primary key not null, token text not null, date integer not null);
		
		create table bots_store (id text primary key not null, value blob, created integer not null, updated integer not null);
    `
	if _, err := db.Exec(sqlStmt); err != nil {
		return err
	}
	return nil
}
