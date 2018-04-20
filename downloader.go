package main

import (
	"io"
	"net/http"
	"os"
	"runtime"
	"sync"
)

type WriteCounter struct {
	Processo             string
	chDownloadInfoReport chan DownloadInfo
	Bytes                uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Bytes += uint64(n)
	wc.Report()
	return n, nil
}

func (wc WriteCounter) Report() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	wc.chDownloadInfoReport <- DownloadInfo{wc.Processo, wc.Bytes}
	// fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	// fmt.Printf("\rDownloading... %v complete", wc.Bytes)
}

type DownloadPayload struct {
	cookieStr string
	titlePDF  string
	dst       string
}

func DownloadPDF(dp *DownloadPayload, filepath string, ci chan DownloadInfo) string {
	Trace.Printf("\nInicio do download: %s para %s", dp.titlePDF, dp.dst)
	Info.Printf("\nInicio do download: %s para %s", dp.titlePDF, dp.dst)
	client := &http.Client{}
	// Declare HTTP Method and Url
	req, _ := http.NewRequest("GET", `https://eprocesso.suiterfb.receita.fazenda/downloadArquivo/`+dp.titlePDF, nil)
	// Set cookie
	req.Header.Set("Cookie", dp.cookieStr)
	// Read response

	out, _ := os.Create(filepath)
	defer out.Close()

	resp, _ := client.Do(req)

	defer resp.Body.Close()

	counter := &WriteCounter{dp.titlePDF, ci, 0}
	n, err := io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		Info.Printf("\n\n\n :( %s miow :(\n\n\n", filepath)
		Trace.Printf("\n\n\n :( %s miow :(\n\n\n", filepath)
	}

	Trace.Printf("\n :) %s salvo com tamanho %v.\n", filepath, n)
	return filepath

}

func Downloader(dp *DownloadPayload, wg *sync.WaitGroup, cc chan bool, ci chan DownloadInfo) {
	// 	Trace.Printf(`
	// +	DOWNLOAD -------------------------------------------------------------------------
	// 		dp (DownloadPayload) 		%T -> 	%#v
	// 		wg (WaitGroup)		%T -> 	%#v
	// 	----------------------------------------------------------------------------------
	// `, dp, dp, wg, wg)
	temporario := DownloadPDF(dp, os.TempDir()+`\`+dp.titlePDF, ci)
	err := os.Rename(temporario, dp.dst)

	if err != nil {
		Info.Println(err)
		Trace.Println(err)
		return
	}
	Info.Printf("\n :) %s ok!.\n", dp.dst)
	Trace.Printf("\n :) %s ok!.\n", dp.dst)
	cc <- true
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
