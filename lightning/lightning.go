package lightning

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// Lightning RPC Struct
type LightningServer struct {
	// RPC filename
	Rpc_filename string

	// Unix Domain Socket
	Unix_socket net.Conn

	// Timeout Read
	Readtimeout time.Duration
}

// Connects to the RPC Server
func LightningRpc(rpc_filename string) *LightningServer {
	ln := LightningServer{
		Rpc_filename: rpc_filename,
		Unix_socket:  nil,
	}

	conn, err := net.Dial("unix", ln.Rpc_filename)
	if err != nil {
		fmt.Println("Unable to connect:", err)
		os.Exit(1)
	}

	ln.Unix_socket = conn

	return &ln
}

// Closes the connection with Unix Domain Socket
func (l *LightningServer) Destroy() {
	l.Unix_socket.Close()
}

// Writes 1024 bytes into the Unix Domain Socket
func (l *LightningServer) write(jsonb []byte) int {
	bytes_written, err := l.Unix_socket.Write(jsonb)
	if err != nil {
		fmt.Println("Unable to write:", err)
	}
	return bytes_written
}

// Read 1024 bytes to buff
func (l *LightningServer) read() []byte {
	buff := []byte{}
	tmp := make([]byte, 1024)

	// Avoid be blocked in case no buffer is available to read
	if l.Readtimeout == 0 {
		l.Readtimeout = 200 // Milliseconds
	}
	timeoutDuration := l.Readtimeout * time.Millisecond
	l.Unix_socket.SetReadDeadline(time.Now().Add(timeoutDuration))

	for {
		n, err := l.Unix_socket.Read(tmp[:])
		if err != nil {
			break
		}
		buff = append(buff, tmp[:n]...)
	}
	return buff
}

// Wrapper to write and read
func (l *LightningServer) call(method_name string, params ...string) string {
	const id int = 0

	// In case params is not needed, use {}, otherwise convert to str
	payload := `{}`
	if len(params) > 0 {
		payload = strings.Join(params, "")
	}

	// Formating json str, for later convert to []byte
	json_str := fmt.Sprintf(
		`{"method": "%s", "params": %s, "id": %d }`,
		method_name, payload, id,
	)

	jsonb := []byte(json_str)

	l.write(jsonb)

	return string(l.read())
}

// Show current block height
func (l *LightningServer) Dev_blockheight() string {
	return l.call("dev-blockheight")
}

// Struct to hold Listfunds
type JsonListFunds struct {
	Jsonrpc string
	Result  struct {
		Outputs []struct {
			Txid   string
			Output int64
			Value  int64
			Status string
		}
		Channels []struct {
			Peer_id           string
			Short_channel_id  string
			Channel_sat       int64
			Channel_total_sat int64
			Funding_txid      string
		}
	}
	Id int64 `json:"id""`
}

// Show funds available for opening channels
func (l *LightningServer) Listfunds() string {
	return l.call("listfunds")
}

// Synchronize the state of our funds with bitcoind
func (l *LightningServer) Dev_rescan_outputs() string {
	return l.call("dev-rescan-outputs")
}

// Crash lightningd by calling fatal()
func (l *LightningServer) Dev_crash() string {
	return l.call("dev-crash")
}

// Structure to hold getinfo json data
type JsonGetInfo struct {
	Jsonrpc string
	Result  struct {
		Id   string
		Port float64

		Address []struct {
			Type    string
			Address string
			Port    int64
		}

		Version     string
		Blockheight int64
		Network     string
	}
	Id int64
}

// Show information about this node
func (l *LightningServer) Getinfo() string {
	return l.call("getinfo")
}

// Show memory objects currently in use
func (l *LightningServer) Dev_memdump() string {
	return l.call("dev-memdump")
}

// Show unreferenced memory objects
func (l *LightningServer) Dev_memleak() string {
	return l.call("dev-memleak")
}

// Show available commands
func (l *LightningServer) Help() string {
	return l.call("help")
}

// Shut down the lightningd process
func (l *LightningServer) Stop() string {
	return l.call("stop")
}

// Show logs, with optional log {level} (info|unusual|debug|io)
func (l LightningServer) Getlog(level ...string) string {
	payload := ""
	if len(level) == 0 {
		payload = `{"level": "info"}`
	} else {
		payload = fmt.Sprintf(
			`{"level": "%s"}`, strings.Join(level, ""),
		)
	}
	return l.call("getlog", payload)
}

// Close the channel with peer {id}
func (l *LightningServer) Close(peerid string) string {
	payload := fmt.Sprintf(
		`{"id": "%s"}`, peerid,
	)
	return l.call("close", payload)
}

// Structure to hold ListNodes json data
type JsonListNodes struct {
	Jsonrpc string
	Result  struct {
		Nodes []struct {
			Nodeid         string
			Alias          string
			Color          string
			Last_Timestamp int64

			Addresses []struct {
				Type    string
				Address string
				Port    int64
			}
		}
	}
}

// Show all nodes in our local network view, filter on node {id}
// if provided
func (l *LightningServer) Listnodes(node_id ...string) string {
	payload := `{}`
	if len(node_id) > 0 {
		payload = fmt.Sprintf(
			`{"id": "%s"}`, node_id[0],
		)
	}
	return l.call("listnodes", payload)
}

// Connect to {peer_id} at {host} and {port}
func (l *LightningServer) Connect(peerid string, host_port ...string) string {
	payload := fmt.Sprintf(
		`{"id": "%s"}`, peerid,
	)

	switch len(host_port) {
	case 1:
		payload = fmt.Sprintf(
			`{"id": "%s", "host": "%s"}`,
			peerid, host_port[0],
		)
	case 2:
		payload = fmt.Sprintf(
			`{"id": "%s", "host": "%s", "port": "%s"}`,
			peerid, host_port[0], host_port[1],
		)
	}
	return l.call("connect", payload)
}

// Structure to hold ListPeers json data
type JsonListPeers struct {
	Jsonrpc string
	Result  struct {
		Peers []struct {
			State     string
			Id        string
			Connected bool
			Netaddr   []string
			Alias     string
			Color     string
			Owner     string
		}
	}
	Id int64 `json:"id""`
}

// NOTE: getpeers got deprecated, renamed to listpeers
func (l *LightningServer) Getpeers(peerid_level ...string) string {
	return l.Listpeers(peerid_level...)
}

// Show current peers, if {level} is set, include {log}s"
func (l *LightningServer) Listpeers(peerid_level ...string) string {
	payload := "{}"
	if len(peerid_level) == 1 {
		payload = fmt.Sprintf(
			`{"id": "%s"}`, peerid_level[0],
		)

	} else if len(peerid_level) == 2 {
		payload = fmt.Sprintf(
			`{"id":, "%s", "level": "%s"}`,
			peerid_level[0], peerid_level[1],
		)
	}
	return l.call("listpeers", payload)
}

// Send {peer_id} a ping of length {len} asking for {pongbytes}"
func (l *LightningServer) Dev_ping(
	peerid string,
	length string,
	pongbytes string) string {

	payload := fmt.Sprintf(
		`{"id": "%s", "len": "%s", "pongbytes": "%s"}`,
		peerid, length, pongbytes,
	)
	return l.call("dev-ping", payload)
}

// Sign and show the last commitment transaction with peer {id}
func (l *LightningServer) Dev_fail(peerid string) string {
	payload := fmt.Sprintf(
		`{"id": "%s"}`, peerid,
	)
	return l.call("dev-fail", payload)
}

// Sign and show the last commitment transaction with peer {id}
func (l *LightningServer) Dev_sign_last_tx(peerid string) string {
	payload := fmt.Sprintf(
		`{"id": "%s"}`, peerid,
	)
	return l.call("dev-sign-last-tx", payload)
}

// Show SHA256 of {secret}
func (l *LightningServer) Dev_rhash(secret string) string {
	payload := fmt.Sprintf(
		`{"secret": "%s"}`, secret,
	)
	return l.call("dev-rhash", payload)
}

// Send along {route} in return for preimage of {rhash}
func (l *LightningServer) Sendpay(route string, rhash string) string {
	payload := fmt.Sprintf(
		`{"route": "%s", "rhash": "%s"}`,
		route, rhash,
	)
	return l.call("sendpay", payload)
}

// Re-enable the commit timer on peer {id}
func (l *LightningServer) Dev_reenable_commit(peerid string) string {
	payload := fmt.Sprintf(
		`{"id": "%s"}`, peerid,
	)
	return l.call("dev-reenable-commit", payload)
}

// Struct to hold Listchannels
type JsonChannels struct {
	Jsonrpc string
	Result  struct {
		Channels []struct {
			Source                string
			Destination           string
			Short_channel_id      string
			Flags                 int64
			Active                bool
			Public                bool
			Last_update           uint64
			Base_fee_millisatoshi uint64
			Fee_per_millionth     uint64
			Delay                 uint64
		}
	}
}

// Show all known channels, accept optional {short_channel_id}
func (l *LightningServer) Listchannels(short_channel_id ...string) string {
	payload := "{}"
	if len(short_channel_id) > 0 {
		payload = fmt.Sprintf(
			`{"short_channel_id": "%s"}`, short_channel_id[0],
		)
	}
	return l.call("listchannels", payload)
}

// Send to {destination} address {satoshi} (or "all")
// amount via Bitcoin transaction
func (l *LightningServer) Withdraw(destination string, satoshi uint64) string {
	payload := fmt.Sprintf(
		`{"destination": "%s", "satoshi": "%d"}`,
		destination, satoshi,
	)
	return l.call("withdraw", payload)
}

// Delete unpaid invoice {label} with {status}
func (l *LightningServer) Delinvoice(label string, status string) string {
	payload := fmt.Sprintf(
		`{"label": "%s", "status": "%s"}`,
		label, status,
	)
	return l.call("delinvoice", payload)
}

// Structure to hold getinfo json data
type JsonInvoice struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Label        string
		Payment_hash string
		Msatoshi     uint64
		Status       string
		Expiry_time  uint64
		Expires_at   uint64
		Bolt11       string
	}
	Id int64 `json:"id""`
}

// Create an invoice for {msatoshi} with {label} and {description} with
// optional {expiry} seconds (default 1 hour)
func (l *LightningServer) Invoice(msatoshi uint64,
	label string,
	description string,
	opt ...string) string {

	payload := fmt.Sprintf(
		`{"msatoshi": "%d", "label": "%s", "description": "%s"`,
		msatoshi, label, description,
	)

	for i := 0; i < len(opt); i++ {
		switch i {
		case 0:
			payload += fmt.Sprintf(`, "expiry": "%s"`, opt[i])
		case 1:
			payload += fmt.Sprintf(`, "fallback": "%s"`, opt[i])
		}
	}
	payload += `}`
	return l.call("invoice", payload)
}

// Fund channel with {id} using {satoshi} satoshis"
func (l *LightningServer) Fundchannel(channel_id string, satoshi string) string {
	payload := fmt.Sprintf(
		`{"id": "%s", "satoshi": "%s"}`,
		channel_id, satoshi,
	)
	return l.call("fundchannel", payload)
}

// Wait for an incoming payment matching the invoice with {label}
func (l *LightningServer) Waitinvoice(label string) string {
	payload := fmt.Sprintf(
		`{"label": "%s"}`, label,
	)
	return l.call("waitinvoice", payload)
}

// Wait for the next invoice to be paid, after {lastpay_index}
// (if supplied)
func (l *LightningServer) Waitanyinvoice(lastpay_index ...string) string {
	payload := "{}"

	if len(lastpay_index) > 0 {
		payload = fmt.Sprintf(
			`{"lastpay_index": "%s"}`, lastpay_index[0],
		)
	}
	return l.call("waitanyinvoice", payload)
}

// Structure to hold ListInvoices json data
type JsonListInvoices struct {
	Jsonrpc string
	Result  struct {
		Invoices []struct {
			Label        string `json:"label"`
			Payment_hash string
			Msatoshi     uint64
			Status       string
			Expiry_time  uint64
			Expires_at   uint64
		}
	}
}

// Show invoice {label} (or all, if no {label))
func (l *LightningServer) Listinvoices(label ...string) string {
	payload := "{}"
	if len(label) > 0 {
		payload = fmt.Sprintf(
			`{"label": "%s"}`, label[0],
		)
	}
	return l.call("listinvoices", payload)
}

// Structure to hold newaddr json data
type JsonNewaddr struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Address string `json:"address"`
	}
	Id int64 `json:"id""`
}

// Get a new address of type {addresstype} of the internal wallet.
func (l *LightningServer) Newaddr(addresstype ...string) string {
	payload := "{}"

	if len(addresstype) > 0 {
		payload = fmt.Sprintf(
			`{"addresstype": "%s"}`, addresstype[0],
		)
	}
	return l.call("newaddr", payload)
}

// Forget the channel with id=peerid
func (l *LightningServer) Dev_forget_channel(peerid string, force ...string) string {
	force_opt := "False"

	if len(force) > 0 {
		force_opt = force[0]
	}

	payload := fmt.Sprintf(
		`{"id": "%s", "force": "%s"}`,
		peerid, force_opt,
	)
	return l.call("dev-forget-channel", payload)
}

type JsonDecodePay struct {
	Jsonrpc string
	Result  struct {
		Currency              string
		Timestamp             int64
		Created_at            int64
		Expiry                int64
		Payee                 string
		Msatoshi              int64
		Description           string
		Min_final_cltv_expiry int64
		Payment_hash          string
		Signature             string
	}
}

// Decode {bolt11}, using {description} if necessary
func (l *LightningServer) Decodepay(bolt11 string, description ...string) string {
	payload := fmt.Sprintf(
		`{"bolt11": "%s"}`, bolt11,
	)

	if len(description) > 0 {
		payload = fmt.Sprintf(
			`{"bolt11": "%s", "description": "%s"}`,
			bolt11, description[0],
		)
	}
	return l.call("decodepay", payload)
}

// Show outgoing payments, regarding {bolt11} or {payment_hash} if set
// Can only specify one of {bolt11} or {payment_hash}
func (l *LightningServer) Listpayments(bolt11_payment_hash ...string) string {
	payload := "{}"
	if len(bolt11_payment_hash) == 2 &&
		len(bolt11_payment_hash[0]) > 0 &&
		len(bolt11_payment_hash[1]) > 0 {
		l.Destroy()
		fmt.Println("Please only specify bolt11 OR payment_hash, not both")
		os.Exit(1)
	}

	if len(bolt11_payment_hash) > 0 {
		if len(bolt11_payment_hash[0]) > 0 {
			payload = fmt.Sprintf(
				`{"bolt11": "%s"}`, bolt11_payment_hash[0],
			)
		} else if len(bolt11_payment_hash[1]) > 0 {
			payload = fmt.Sprintf(
				`{"payment_hash": "%s"}`, bolt11_payment_hash[1],
			)
		}
	}
	return l.call("listpayments", payload)
}

// Show route to {id} for {msatoshi}, using {riskfactor} and optional
// {cltv} (default 9)
func (l *LightningServer) Getroute(
	peerid string,
	msatoshi uint64,
	riskfactor string,
	cltv ...string) string {

	cltv_opt := ""
	if len(cltv) == 0 {
		cltv_opt = `, "cltv": "9"`
	} else {
		cltv_opt = fmt.Sprintf(`, "cltv": "%s"`, cltv[0])
	}
	payload := fmt.Sprintf(
		`{"id": "%s", "msatoshi": "%d", "riskfactor": "%s" %s}`,
		peerid, msatoshi, riskfactor, cltv_opt,
	)
	return l.call("getroute", payload)
}

// Send payment specified by {bolt11} with optional {msatoshi}
// (if and only if {bolt11} does not have amount),
//
// {description} (required if {bolt11} uses description hash)
// and {riskfactor} (default 1.0)
func (l *LightningServer) Pay(bolt11 string, opt ...string) string {

	payload := fmt.Sprintf(`{"bolt11": "%s"`, bolt11)

	for i := 0; i < len(opt); i++ {
		switch i {
		case 0:
			payload += fmt.Sprintf(`, "msatoshi": "%d"`, opt[i])
		case 1:
			payload += fmt.Sprintf(`, "description": "%s"`, opt[i])
		case 2:
			payload += fmt.Sprintf(`, "riskfactor": "%s"`, opt[i])
		}
	}
	payload += `}`

	return l.call("pay", payload)
}

// Set feerate in satoshi-per-kw for {immediate}, {normal} and {slow}
// (each is optional, when set, separate by spaces) and show the value
// of those three feerates
func (l *LightningServer) Dev_setfees(opt ...string) string {
	payload := `{`
	for i := 0; i < len(opt); i++ {
		switch i {
		case 0:
			payload += fmt.Sprintf(`"immediate": "%s"`, opt[i])
		case 1:
			payload += fmt.Sprintf(`, "normal": "%s"`, opt[i])
		case 2:
			payload += fmt.Sprintf(`, "slow": "%s"`, opt[i])
		}
	}
	payload += `}`

	return l.call("dev-setfees", payload)
}
