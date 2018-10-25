package db

import (
	"database/sql"
	_ "github.com/mutecomm/go-sqlcipher"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/repo"
	"path"
	"sync"
)

var log = logging.MustGetLogger("db")

type SQLiteDatastore struct {
	config             repo.ConfigStore
	profile            repo.ProfileStore
	contacts           repo.ContactStore
	threads            repo.ThreadStore
	threadPeers        repo.ThreadPeerStore
	threadMessages     repo.ThreadMessageStore
	blocks             repo.BlockStore
	notifications      repo.NotificationStore
	cafeSessions       repo.CafeSessionStore
	cafeRequests       repo.CafeRequestStore
	cafeMessages       repo.CafeMessageStore
	cafeClientNonces   repo.CafeClientNonceStore
	cafeClients        repo.CafeClientStore
	cafeClientThreads  repo.CafeClientThreadStore
	cafeClientMessages repo.CafeClientMessageStore
	db                 *sql.DB
	lock               *sync.Mutex
}

func Create(repoPath, pin string) (*SQLiteDatastore, error) {
	var dbPath string
	dbPath = path.Join(repoPath, "datastore", "mainnet.db")
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if pin != "" {
		p := "pragma key='" + pin + "';"
		conn.Exec(p)
	}
	mux := new(sync.Mutex)
	sqliteDB := &SQLiteDatastore{
		config:             NewConfigStore(conn, mux, dbPath),
		profile:            NewProfileStore(conn, mux),
		contacts:           NewContactStore(conn, mux),
		threads:            NewThreadStore(conn, mux),
		threadPeers:        NewThreadPeerStore(conn, mux),
		threadMessages:     NewThreadMessageStore(conn, mux),
		blocks:             NewBlockStore(conn, mux),
		notifications:      NewNotificationStore(conn, mux),
		cafeSessions:       NewCafeSessionStore(conn, mux),
		cafeRequests:       NewCafeRequestStore(conn, mux),
		cafeMessages:       NewCafeMessageStore(conn, mux),
		cafeClientNonces:   NewCafeClientNonceStore(conn, mux),
		cafeClients:        NewCafeClientStore(conn, mux),
		cafeClientThreads:  NewCafeClientThreadStore(conn, mux),
		cafeClientMessages: NewCafeClientMessageStore(conn, mux),
		db:                 conn,
		lock:               mux,
	}

	return sqliteDB, nil
}

func (d *SQLiteDatastore) Ping() error {
	return d.db.Ping()
}

func (d *SQLiteDatastore) Close() {
	d.db.Close()
}

func (d *SQLiteDatastore) Config() repo.ConfigStore {
	return d.config
}

func (d *SQLiteDatastore) Profile() repo.ProfileStore {
	return d.profile
}

func (d *SQLiteDatastore) Contacts() repo.ContactStore {
	return d.contacts
}

func (d *SQLiteDatastore) Threads() repo.ThreadStore {
	return d.threads
}

func (d *SQLiteDatastore) ThreadPeers() repo.ThreadPeerStore {
	return d.threadPeers
}

func (d *SQLiteDatastore) ThreadMessages() repo.ThreadMessageStore {
	return d.threadMessages
}

func (d *SQLiteDatastore) Blocks() repo.BlockStore {
	return d.blocks
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

func (d *SQLiteDatastore) CafeClientThreads() repo.CafeClientThreadStore {
	return d.cafeClientThreads
}

func (d *SQLiteDatastore) CafeClientMessages() repo.CafeClientMessageStore {
	return d.cafeClientMessages
}

func (d *SQLiteDatastore) Copy(dbPath string, password string) error {
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
	if password == "" {
		cp = `attach database '` + dbPath + `' as plaintext key '';`
		for _, name := range tables {
			cp = cp + "insert into plaintext." + name + " select * from main." + name + ";"
		}
	} else {
		cp = `attach database '` + dbPath + `' as encrypted key '` + password + `';`
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

func (d *SQLiteDatastore) InitTables(password string) error {
	return initDatabaseTables(d.db, password)
}

func initDatabaseTables(db *sql.DB, pin string) error {
	var sqlStmt string
	if pin != "" {
		sqlStmt = "PRAGMA key = '" + pin + "';"
	}
	sqlStmt += `
	create table config (key text primary key not null, value blob);

    create table profile (key text primary key not null, value blob);

	create table contacts (id text primary key not null, username text not null, inboxes text not null, added integer not null);
    create index contact_username on contacts (username);
	create index contact_added on contacts (added);

    create table threads (id text primary key not null, name text not null, sk blob not null, head text not null);

    create table thread_peers (id text not null, threadId text not null, welcomed integer not null, primary key (id, threadId));
    create index thread_peer_id on thread_peers (id);
    create index thread_peer_threadId on thread_peers (threadId);
	create index thread_peer_welcomed on thread_peers (welcomed);

    create table blocks (id text primary key not null, date integer not null, parents text not null, threadId text not null, authorId text not null, type integer not null, dataId text, dataKey blob, dataCaption text, dataMetadata blob);
    create index block_dataId on blocks (dataId);
    create index block_threadId_type_date on blocks (threadId, type, date);

	create table thread_messages (id text primary key not null, peerId text not null, envelope blob not null, date integer not null);
	create index thread_message_date on thread_messages (date);

    create table notifications (id text primary key not null, date integer not null, actorId text not null, actorUsername text not null, subject text not null, subjectId text not null, blockId text, dataId text, type integer not null, body text not null, read integer not null);
    create index notification_date on notifications (date);
	create index notification_actorId on notifications (actorId);
    create index notification_subjectId on notifications (subjectId);
    create index notification_blockId on notifications (blockId);
    create index notification_read on notifications (read);

    create table cafe_sessions (cafeId text primary key not null, access text not null, refresh text not null, expiry integer not null);

    create table cafe_requests (id text primary key not null, peerId text not null, targetId text not null, cafeId text not null, type integer not null, date integer not null);
    create index cafe_request_cafeId on cafe_requests (cafeId);
	create index cafe_request_date on cafe_requests (date);

	create table cafe_messages (id text primary key not null, peerId text not null, date integer not null);
	create index cafe_message_date on cafe_messages (date);

	create table cafe_client_nonces (value text primary key not null, address text not null, date integer not null);

    create table cafe_clients (id text primary key not null, address text not null, created integer not null, lastSeen integer not null);
    create index cafe_client_address on cafe_clients (address);
    create index cafe_client_lastSeen on cafe_clients (lastSeen);

    create table cafe_client_threads (id text not null, clientId text not null, skCipher blob not null, headCipher blob not null, nameCipher blob not null, primary key (id, clientId));
    create index cafe_client_thread_clientId on cafe_client_threads (clientId);

	create table cafe_client_messages (id text not null, peerId text not null, clientId text not null, date integer not null, primary key (id, clientId));
    create index cafe_client_message_clientId on cafe_client_messages (clientId);
	create index cafe_client_message_date on cafe_client_messages (date);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}
