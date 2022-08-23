package ffiwrapper

import (
	"bufio"
	"github.com/filecoin-project/go-state-types/abi"
	"io"
	"os"
)

func (sb *Sealer) isPledgeRequest(entries int, ps abi.UnpaddedPieceSize) bool {
	return entries == 0 && ps == 34091302912
}

func (sb *Sealer) copyPledgeSector(source string, dest string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	reader := bufio.NewReader(sourceFile)
	destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer destFile.Close()
	writer := bufio.NewWriter(destFile)
	_, err = io.Copy(writer, reader)
	return err
}
