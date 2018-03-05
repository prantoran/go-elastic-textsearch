package utilities

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	es "gopkg.in/olivere/elastic.v5"

	"github.com/prantoran/go-elastic-textsearch/data"
)

// LaunchESConnectionTest encapsulates the launching of go routines
// to test whether ElasticSearch is working
func LaunchESConnectionTest() {
	errc := make(chan error)

	go func() {
		err := showNodes(data.Escon.Client)
		if err != nil {
			log.Printf("nodes info failed: %v", err)
		}

		t := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-t.C:
				err := showNodes(data.Escon.Client)
				if err != nil {
					log.Printf("nodes info failed: %v", err)
				}
			}
		}
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		log.Printf("existing with signal %v", fmt.Sprint(<-c))
		errc <- nil
	}()

	if err := <-errc; err != nil {
		log.Printf("LaunchESConnectionTest: %v\n", err)
		os.Exit(1)
	}

}

func showNodes(client *es.Client) error {
	ctx := context.Background()
	info, err := client.NodesInfo().Do(ctx)
	if err != nil {
		return err
	}
	log.Printf("Cluster %q with %d node(s)", info.ClusterName, len(info.Nodes))
	for id, node := range info.Nodes {
		log.Printf("- Node %s with IP %s", id, node.IP)
	}
	return nil
}
