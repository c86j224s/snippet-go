package alfa

import (
	"os"
)

func FileReadAll(fileName string) ([]byte, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	fcontent := make([]byte, finfo.Size())

	_, err = f.Read(fcontent)
	if err != nil {
		return nil, err
	}

	return fcontent, nil
}
