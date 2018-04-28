package api

import (
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	"github.com/vechain/thor/api/accounts"
	"github.com/vechain/thor/api/blocks"
	"github.com/vechain/thor/api/doc"
	"github.com/vechain/thor/api/events"
	"github.com/vechain/thor/api/node"
	"github.com/vechain/thor/api/transactions"
	"github.com/vechain/thor/api/transfers"
	"github.com/vechain/thor/chain"
	"github.com/vechain/thor/eventdb"
	"github.com/vechain/thor/state"
	"github.com/vechain/thor/transferdb"
	"github.com/vechain/thor/txpool"
)

//New return api router
func New(chain *chain.Chain, stateCreator *state.Creator, txPool *txpool.TxPool, eventDB *eventdb.EventDB, transferDB *transferdb.TransferDB, nw node.Network) http.HandlerFunc {
	router := mux.NewRouter()

	// to serve api doc and swagger-ui
	router.PathPrefix("/doc").Handler(
		http.StripPrefix("/doc/", http.FileServer(
			&assetfs.AssetFS{
				Asset:     doc.Asset,
				AssetDir:  doc.AssetDir,
				AssetInfo: doc.AssetInfo})))

	// redirect swagger-ui
	router.Path("/").HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, "doc/swagger-ui/", http.StatusTemporaryRedirect)
		})

	accounts.New(chain, stateCreator).
		Mount(router, "/accounts")
	events.New(eventDB).
		Mount(router, "/events")
	transfers.New(transferDB).
		Mount(router, "/transfers")
	blocks.New(chain).
		Mount(router, "/blocks")
	transactions.New(chain, txPool, transferDB).
		Mount(router, "/transactions")
	node.New(nw).
		Mount(router, "/node")

	return router.ServeHTTP
}