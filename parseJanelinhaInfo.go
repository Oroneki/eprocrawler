package main

import "strings"
import "regexp"
import "strconv"
import "fmt"

type JanelinhaInfo struct {
	NomeDoc    string `json:"nome_doc"`
	Ordem      int64  `json:"ordem"`
	PagInicial int64  `json:"pag_inicio"`
	PagFinal   int64  `json:"pag_fim"`
	Situacao   string `json:"situacao"`
	Tamanho    int64  `json:"tamanho"`
}

func formatAndSendToWebsocket(infos []JanelinhaInfo, processoImpuro string, wsChan chan WebSocketMessage) {
	go func(infos []JanelinhaInfo, processoImpuro string, wsChan chan WebSocketMessage) {
		var payload string = fmt.Sprintf("%s|@@|", processoImpuro)
		for _, info := range infos {
			infoFmt := fmt.Sprintf(
				"%s|%d|%d|%d|%s|%d|##|",
				info.NomeDoc,
				info.Ordem,
				info.PagInicial,
				info.PagFinal,
				info.Situacao,
				info.Tamanho,
			)
			trace.Printf("JANELINHA_INFO_PROCESSO info: %v --> %s", info, infoFmt)
			payload = fmt.Sprintf("%s%s",payload,infoFmt)
		}
		trace.Printf("JANELINHA_INFO_PROCESSO\nPayload: %s", payload)
		wsChan <- WebSocketMessage{
			Tipo:    "JANELINHA_INFO_PROCESSO",
			Payload: payload,
		}
	}(infos, processoImpuro, wsChan)
}

var rejanelinhainfo = regexp.MustCompile(`exibirDivInfoDocumento\(\d+\,\s\'(?P<nome>.*?)\'\,\'(?P<ordem>\d+)\@\#\#\@Fl\.\s(?P<i>\d+)\sa\s(?P<f>\d+)\W+\'(?P<sit>.*?)\'[^\d]*\'(?P<tam>[\d\,\.]*\s[A-Z]{2})\'`)

func splitAndFilterJanelinhaInfo(pld string) []string {
	infos := strings.Split(pld, "@@@@@@@@@@")
	filteredInfos := []string{}
	for _, info := range infos {
		if !strings.Contains(info, "exibirDivInfo") {
			continue
		}
		filteredInfos = append(filteredInfos, info)

	}
	return filteredInfos
}

func apllyRegexJanelinhaInfo(info string) JanelinhaInfo {

	all := rejanelinhainfo.FindStringSubmatch(info)

	ord, er1 := strconv.ParseInt(all[2], 10, 64)
	if er1 != nil {
		ord = -1
	}

	ini, er2 := strconv.ParseInt(all[3], 10, 64)
	if er2 != nil {
		ini = -1
	}

	fin, er3 := strconv.ParseInt(all[4], 10, 64)
	if er3 != nil {
		fin = -1
	}

	jan := JanelinhaInfo{
		NomeDoc:    all[1],
		Ordem:      ord,
		PagInicial: ini,
		PagFinal:   fin,
		Situacao:   all[5],
		Tamanho:    -1,
	}

	trace.Printf("apllyRegexJanelinhaInfo --> %v", jan)
	
	return jan
	
}

func parseJanelinhaInfo(pld string) []JanelinhaInfo {
	trace.Printf("parseJanelinhaInfo --> %v", pld)
	infos := splitAndFilterJanelinhaInfo(pld)
	filteredInfos := []JanelinhaInfo{}
	for _, info := range infos {
		reg := apllyRegexJanelinhaInfo(info)
		filteredInfos = append(filteredInfos, reg)
	}
	trace.Printf("parseJanelinhaInfo --> %v", filteredInfos)
	return filteredInfos
}
