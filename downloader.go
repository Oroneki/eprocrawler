package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
)

type writeCounter struct {
	Processo             string
	chDownloadInfoReport chan downloadInfo
	Bytes                uint64
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Bytes += uint64(n)
	wc.Report()
	return n, nil
}

func (wc writeCounter) Report() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	wc.chDownloadInfoReport <- downloadInfo{wc.Processo, wc.Bytes}
	// fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	// fmt.Printf("\rDownloading... %v complete", wc.Bytes)
}

type downloadPayload struct {
	cookieStr string
	titlePDF  string
	dst       string
}

func downloadPDF(dp *downloadPayload, filepath string, ci chan downloadInfo) string {
	trace.Printf("\nInicio do download: %s para %s", dp.titlePDF, dp.dst)
	info.Printf("\nInicio do download: %s para %s", dp.titlePDF, dp.dst)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	// Declare HTTP Method and Url
	PROCESSO := strings.Split(dp.dst, `\`)
	DOWNLOAD_URL := `https://eprocesso.suiterfb.receita.fazenda/ControleDownloadArquivoZip.asp?psAcao=download&psNumeroProcesso=` + strings.Replace(PROCESSO[len(PROCESSO)-1], ".pdf", "", 1) + `&psNomeArquivoZip=` + dp.titlePDF
	trace.Printf("\n\n\n\nDownload URL: '%s'\n\n\n\n", DOWNLOAD_URL)
	req, _ := http.NewRequest("GET", DOWNLOAD_URL, nil)
	// Set cookie
	req.Header.Set("Cookie", dp.cookieStr)
	// Read response

	out, _ := os.Create(filepath)
	defer out.Close()

	resp, e := client.Do(req)
	if e != nil {
		panic(e)
	}
	trace.Printf("\nResposta: %s para %s\n   RESP %s %d [%d] \n  HEADER: %s", dp.titlePDF, dp.dst, resp.Status, resp.StatusCode, resp.ContentLength, resp.Header)

	defer resp.Body.Close()

	counter := &writeCounter{dp.titlePDF, ci, 0}
	n, err := io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		info.Printf("\n\n\n :( %s miow :(\n\n\n", filepath)
		trace.Printf("\n\n\n :( %s miow :(\n\n\n", filepath)
	}

	trace.Printf("\n :) %s salvo com tamanho %v.\n", filepath, n)

	return filepath

}

func downloader(dp *downloadPayload, wg *sync.WaitGroup, cc chan bool, ci chan downloadInfo, wsWriteChannel chan WebSocketMessage) {
	trace.Printf(`
	+	DOWNLOAD -------------------------------------------------------------------------
			dp (DownloadPayload) 		%T -> 	%#v
		----------------------------------------------------------------------------------
	`, dp, dp)
	temporario := downloadPDF(dp, os.TempDir()+`\`+dp.titlePDF, ci)
	err := os.Rename(temporario, dp.dst)

	if err != nil {
		info.Println(err)
		trace.Println(err)
		return
	}
	info.Printf("\n\n:) %s ok!.\n", dp.dst)
	trace.Printf("\n\n:) %s ok!.\n", dp.dst)
	cc <- true
	go func(t string, d string) {
		wsWriteChannel <- WebSocketMessage{
			Tipo:    "DOWNLOAD_FINISHED",
			Payload: fmt.Sprintf("%s|%s", t, d),
		}
	}(dp.titlePDF, dp.dst)
	wg.Done()
}

func getUserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
