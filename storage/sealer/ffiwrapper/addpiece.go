package ffiwrapper

import (
	"bufio"
	"github.com/ipfs/go-cid"
	"io"
	"io/ioutil"
	"os"
)

type AddPieceTemplate struct {
	Exist      bool
	ParentPath string
	DataPath   string
	CidPath    string
}

var APT = AddPieceTemplate{
	Exist: false,
}

func (a *AddPieceTemplate) Config(psp string) {
	a.Exist = false
	a.ParentPath = psp
	a.DataPath = a.ParentPath + "/sector"
	a.CidPath = a.ParentPath + "/cid"
}

func (a *AddPieceTemplate) IsExist() bool {
	if a.Exist {
		return true
	}

	_, err := os.Stat(a.DataPath)
	if err != nil {
		log.Errorf("pledge sector does not exist: %v", err)
		return false
	}
	_, err = os.Stat(a.CidPath)
	if err != nil {
		log.Errorf("pledge cid does not exist: %v", err)
		return false
	}

	a.Exist = true
	return a.Exist
}

func (a *AddPieceTemplate) GetPieceCID() (cid.Cid, error) {
	cidData, err := ioutil.ReadFile(a.CidPath)
	if err != nil {
		return cid.Undef, err
	}
	_, pieceCID, err := cid.CidFromBytes(cidData)
	if err != nil {
		return cid.Undef, err
	}
	return pieceCID, nil
}

func (a *AddPieceTemplate) SetPieceCID(pieceCID cid.Cid) error {
	return ioutil.WriteFile(a.CidPath, pieceCID.Bytes(), 0644)
}

func (a *AddPieceTemplate) GetData(unsealed string) error {
	return CopySector(a.DataPath, unsealed)
}

func (a *AddPieceTemplate) SetData(unsealed string) error {
	return CopySector(unsealed, a.DataPath)
}

func CopySector(source string, dest string) error {
	log.Errorf("copy pledge sector source: %s, dest: %s", source, dest)
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
