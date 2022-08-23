package ffiwrapper

import (
	"bufio"
	"github.com/ipfs/go-cid"
	"io"
	"io/ioutil"
	"os"
)

type addPieceTemplate struct {
	exist      bool
	parentPath string
	dataPath   string
	cidPath    string
}

var adt *addPieceTemplate

func GetADT() *addPieceTemplate {
	if adt == nil {
		return &addPieceTemplate{}
	}
	return adt
}

func ConfigADT(pledgeSectorPath string) {
	adt := GetADT()
	adt.parentPath = pledgeSectorPath
	adt.exist = false
	adt.dataPath = adt.parentPath + "/sector"
	adt.cidPath = adt.parentPath + "/cid"
}

func IsADTExist() bool {
	adt := GetADT()
	if adt.exist {
		return true
	}

	_, err := os.Stat(adt.dataPath)
	if err != nil {
		log.Errorf("pledge sector does not exist: %v", err)
		return false
	}
	_, err = os.Stat(adt.cidPath)
	if err != nil {
		log.Errorf("pledge cid does not exist: %v", err)
		return false
	}

	adt.exist = true
	return adt.exist
}

func GetADTPieceCID() (cid.Cid, error) {
	cidData, err := ioutil.ReadFile(GetADT().cidPath)
	if err != nil {
		return cid.Undef, err
	}
	_, pieceCID, err := cid.CidFromBytes(cidData)
	if err != nil {
		return cid.Undef, err
	}
	return pieceCID, nil
}

func SetADTPieceCID(pieceCID cid.Cid) error {
	return ioutil.WriteFile(GetADT().cidPath, pieceCID.Bytes(), 0644)
}

func GetADTData(unsealed string) error {
	return CopyPledgeSector(GetADT().dataPath, unsealed)
}

func SetADTData(unsealed string) error {
	return CopyPledgeSector(unsealed, GetADT().dataPath)
}

func CopyPledgeSector(source string, dest string) error {
	log.Debugf("copy pledge sector source: %s, dest: %s", source, dest)
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
