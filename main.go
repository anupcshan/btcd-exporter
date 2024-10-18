package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/btcsuite/btcd/rpcclient"
)

func main() {
	user := flag.String("user", "", "BTCD RPC User")
	password := flag.String("password", "", "BTCD RPC Password")
	listen := flag.String("listen", ":6061", "Listen address")
	flag.Parse()

	connCfg := &rpcclient.ConnConfig{
		Host:       "127.0.0.1:8334",
		Endpoint:   "ws",
		User:       *user,
		Pass:       *password,
		DisableTLS: true,
	}

	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		miningInfo, err := client.GetMiningInfo()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to read mining info: %v", err)
			return
		}
		log.Printf("Mining info: %+v", miningInfo)

		buf := new(bytes.Buffer)
		fmt.Fprintf(buf, "btcd_blockcount %d\n", miningInfo.Blocks)
		fmt.Fprintf(buf, "btcd_difficulty %f\n", miningInfo.Difficulty)
		fmt.Fprintf(buf, "btcd_networkhashrate %f\n", miningInfo.NetworkHashPS)
		fmt.Fprintf(buf, "btcd_polledtx %d\n", miningInfo.PooledTx)

		_, _ = buf.WriteTo(w)
	})

	log.Fatal(http.ListenAndServe(*listen, nil))
}
