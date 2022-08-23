package ffiwrapper

import (
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("ffiwrapper")

type Sealer struct {
	sectors SectorProvider

	pledgeSectorExist bool
	pledgeSectorPath  string
	pledgeSectorCid   cid.Cid

	stopping chan struct{}
}

func (sb *Sealer) Stop() {
	close(sb.stopping)
}
