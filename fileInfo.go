package main

import "os"

func getFolderInfo(folderPath string) (names []string, err error) {
	f, err := os.Open(folderPath)
	defer f.Close()
	if err != nil {
		info.Printf("Erro ao abrir o arquivo %s", folderPath)
		return nil, err
	}
	names, errr := f.Readdirnames(0)
	if errr != nil {
		info.Printf("Não foi possivel ler arquivos de %s", folderPath)
		return nil, errr
	}
	return names, nil
}
