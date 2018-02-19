package alfa

import (
	"os"
)

// FileReadAll 은 파일이름을  받아, 해당 파일을 읽어들여 바이트열로 리턴합니다.
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
