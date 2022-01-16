package file

import "os"

func Write(name string, data []byte) error {
	err := os.WriteFile(name, data, 0755)
	if err != nil {
		return err
	}

	return nil
}
