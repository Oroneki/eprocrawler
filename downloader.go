package main

import (
	"io"
	"net/http"
	"os"
	"runtime"
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
	client := &http.Client{}
	// Declare HTTP Method and Url
	req, _ := http.NewRequest("GET", `https://eprocesso.suiterfb.receita.fazenda/downloadArquivo/`+dp.titlePDF, nil)
	// Set cookie
	req.Header.Set("Cookie", dp.cookieStr)
	// Read response

	out, _ := os.Create(filepath)
	defer out.Close()

	resp, e := client.Do(req)
	if e != nil {
		panic(e)
	}

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

func downloader(dp *downloadPayload, wg *sync.WaitGroup, cc chan bool, ci chan downloadInfo) {
	// 	Trace.Printf(`
	// +	DOWNLOAD -------------------------------------------------------------------------
	// 		dp (DownloadPayload) 		%T -> 	%#v
	// 		wg (WaitGroup)		%T -> 	%#v
	// 	----------------------------------------------------------------------------------
	// `, dp, dp, wg, wg)
	temporario := downloadPDF(dp, os.TempDir()+`\`+dp.titlePDF, ci)
	err := os.Rename(temporario, dp.dst)

	if err != nil {
		info.Println(err)
		trace.Println(err)
		return
	}
	info.Printf("\n :) %s ok!.\n", dp.dst)
	trace.Printf("\n :) %s ok!.\n", dp.dst)
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
