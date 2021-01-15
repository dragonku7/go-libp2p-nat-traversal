package main

import (
	"context"
	"flag"
	"fmt"
	secio "github.com/libp2p/go-libp2p-secio"

	ntraversal "github.com/dragonku7/go-libp2p-nat-traversal"
	logging "github.com/ipfs/go-log"
	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	inet "github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peerstore"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	ma "github.com/multiformats/go-multiaddr"
)

var log = logging.Logger("nat-traversal")

type netNotifiee struct{}

func (nn *netNotifiee) Connected(n inet.Network, c inet.Conn) {
	h.Peerstore().AddAddr(c.RemotePeer(), c.RemoteMultiaddr(), peerstore.PermanentAddrTTL)
	fmt.Printf("Connected to: %s/p2p/%s\n", c.RemoteMultiaddr(), c.RemotePeer().Pretty())
}

func (nn *netNotifiee) Disconnected(n inet.Network, v inet.Conn)   {}
func (nn *netNotifiee) OpenedStream(n inet.Network, v inet.Stream) {}
func (nn *netNotifiee) ClosedStream(n inet.Network, v inet.Stream) {}
func (nn *netNotifiee) Listen(n inet.Network, a ma.Multiaddr)      {}
func (nn *netNotifiee) ListenClose(n inet.Network, a ma.Multiaddr) {}

var h host.Host

func main() {
	logging.SetLogLevel("nat-traversal", "DEBUG")

	port := flag.Int("p", 3000, "port number")
	flag.Parse()

	ctx := context.Background()

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	sourceMultiAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *port))

	var err error
	h, err = libp2p.New(ctx, libp2p.ListenAddrs(sourceMultiAddr), libp2p.Security(secio.ID, secio.New))
	if err != nil {
		panic(err)
	}

	no := &netNotifiee{}
	h.Network().Notify(no)

	fmt.Println("This node: ", h.ID().Pretty(), " ", h.Addrs())

	d, err := dht.New(ctx, h)
	if err != nil {
		panic(err)
	}

	ntraversal.NewNatTraversal(ctx, &h, d)

	select {}
}
