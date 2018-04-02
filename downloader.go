package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"sync"
)

type DownloadPayload struct {
	cookieStr string
	titlePDF  string
	dst       string
}

func DownloadPDF(dp *DownloadPayload, filepath string) string {
	// Trace.Printf(`
	// 	******** DOWNLOADPDF ***************************
	// 			dp (DownloadPayload) 		%T -> 	%#v
	// 			filepath		%T -> 	%#v
	// 		****************
	// 	`, dp, dp, filepath, filepath)
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

	_, err := io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("\n\n\n :( %s miow :(\n\n\n", filepath)
	}

	// Trace.Printf("\n :) %s salvo com tamanho %v.\n", filepath, n)
	return filepath

}

func Downloader(dp *DownloadPayload, wg *sync.WaitGroup) {
	// 	Trace.Printf(`
	// +	DOWNLOAD -------------------------------------------------------------------------
	// 		dp (DownloadPayload) 		%T -> 	%#v
	// 		wg (WaitGroup)		%T -> 	%#v
	// 	----------------------------------------------------------------------------------
	// `, dp, dp, wg, wg)
	temporario := DownloadPDF(dp, os.TempDir()+`\`+dp.titlePDF)
	err := os.Rename(temporario, dp.dst)

	if err != nil {
		fmt.Println(err)
		return
	}
	Info.Printf("\n :) %s ok!.\n", dp.dst)
	wg.Done()
}

// func main() {
// 	DownloadPDF(
// 		`CAMWebPIN=e56eeac514c7b3c76f5c01d1faae8166; cw=1.0; E%2DProcessoCache=fg%99%5B%B4%AA%B2%A7%9E%AB%B9pS%A7%5C%A0%A6%B4%9C%9D%A3%A3y%5C%9Do%B3%B8%A7%97%97%AA%B8; E%2DProcessoID=%81%87%D2%93%DE%DB%E2%D7%D7%D4%D7%AF%95%B3%8B%D3%DB%DE%D1%C6%DF%B0%82i%C6%96%E3%D7%95%AA%CA%E5%D8%AE%86%CE%8B%DC%C0%D0%C6%CE%E2%E1%A1%8F%A2l%B6%D3%DB%D6%CA%99%B4%A4%90%CE%98%D9%E5%E3%D5%C6%D7%E2%B2w%CE%9A%DF%C7%DD%CC%C9%D4%D7%A5%60%A7p%D1%DE%E2%C8%8B%B6%E2%AE%89%CE%91%E5%E4%D0%C7%D4%E5%C8%AE%8C%C9%8B%D4%D7%AC%A5%AB%D4%DF%B3%88%8Bq%D5%E4%D4%D1%C8%DC%D4%ACx%D3%93%D4%D3%D3%C8%A2%B5%C7%B2%98%CAP%C6%DB%E2%C4%D4%BA%D8%B2%88%D3%8D%D9%D3%DB%B1%C6%D6%DC%AF%91%C6%96%AD%B4%B5%C4%D1%E6%D8fw%CE%9A%DF%C7%DD%CC%C9%D4%D7%A5%60%B8P%C5%E0%D8%C7%C6%D7%D8%7Dv%8Bo%E1%E7%D8%D3%CA%B0%C6ff%C6%8D%D8%D7%B2%CF%CE%D8%E1%B4%88%C4s%DE%DB%D2%CC%C6%DF%DC%BA%84%C8%8B%DF%D1%C5%C4%D1%DC%D7%A1%87%D4g%C3%A3%95%B7%B9%BF%D2%83d%A8r%B5%D1%C3%B3%C4%C8%C1%89g%A6n%B5%AF%B8%94%95%99%B7%94%82%A6%7E%C5%B3%BB%AC%BF%B4%B6%81r%C4%7E%C0%D1%C4%B1%AE%B7%B4%84h%A2n%A0%AA%9E%93%98%A2%A5pT%9DJ%A1%A4%A9%94%9D%AD%A6qI%A6%9F%E4%E1%C3%C8%D8%E7%D8%7De%AB%8B%DC%E5%D4%89%A8%B4%C0%97%88%C7z%D9%E0%AC%B6%CA%A8%A9%A5%88%C6%8D%A5%A3%A3%C6%9C%D5%A6%A3Z%9B%90%A5%D5%9F%94%C9%A4%D9%A1%84%CAb%A1%A8%A5%89%AA%E4%E8%A9%93%CA%7D%D5%DE%D4%C6%CE%E2%E1%A1%87%C6m%D1%DB%E7%C4%B9%E5%D4%A2%84%D1%92%DF%AF%C2%96%99%EF%AB%BCT%8Bp%D9%DE%E3%D5%D4%C6%D8%AC%88%C8%93%DF%E0%D0%C7%D4%B6%D4%A9%9B%C6%7E%E2%D3%D1%C4%D1%DB%E2%7Dv%9Da%A2%A4%95%B1%DA%E0%D8%B2%92%B5%8B%D7%DB%DD%C4%A2%C6%A4ft%DA%8B%DE%E6%D8%C7%C6%D7%D8%90%95%D4%8D%D5%E5%E2%D2%D8%C3%D4%A7%8C%D3%8B%D3%D3%DE%A0%B8%A4%A3fv%D4%96%D9%D5%B9%D8%D3%E7%D4%A4%84%A2%7D%C3%98%B4%96%9B%A9%A8%84V%96b%A3%B6%B5%A8%AA%B8%A6%84X%9E%5E%B6%A2%A8%97%A7%A9%ACsX%99l%B5; verificarRequisicaoAjax=0`,
// 		`https://eprocesso.suiterfb.receita.fazenda/downloadArquivo/10010012218031860_COPIA_20180309151704842.pdf`,
// 		`C:\tmp\teste_.pdf`,
// 	)
// }

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
