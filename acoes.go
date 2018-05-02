package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"strings"
)

func apagarArquivo(path string) bool {
	err := os.Remove(path)
	ret := true
	if err != nil {
		Trace.Println("err apagarArquivo", path)
		Trace.Println(err)
		ret = false
		fmt.Println("Erro ao apagar ", path)
	}
	return ret
}

func apagarArquivos(pasta string, listacomma string) ([]string, []string) {
	arrayProcs := strings.Split(listacomma, ",")
	var arrayOk = []string{}
	var arrayErrado = []string{}
	for _, proc := range arrayProcs {
		fullpath := pasta + proc + ".pdf"
		ok := apagarArquivo(fullpath)
		if ok != true {
			arrayErrado = append(arrayErrado, proc)
			continue
		}
		arrayOk = append(arrayOk, proc)
	}
	return arrayOk, arrayErrado
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
