package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	// "strings"
	"time"
	// "github.com/atotto/clipboard"
	// "github.com/scjalliance/comshim"
)

var nomesDeJanela [10]string = [10]string{
	"_____JANELINHA_____ZERO______",
	"_____JANELINHA______UM_______",
	"_____JANELINHA_____DOIS______",
	"_____JANELINHA_____TRES______",
	"_____JANELINHA____QUATRO_____",
	"_____JANELINHA_____CINCO_____",
	"_____JANELINHA_____SEIS______",
	"_____JANELINHA_____SETE______",
	"_____JANELINHA_____OITO______",
	"_____JANELINHA_____NOVE______",
}

type Janelinha struct {
	id               string
	oleAPI           *apiConn
	entradaProcessos <-chan *Processo
	waitGroup        *sync.WaitGroup
	downloadChannel  chan *DownloadPayload
}

type DownloadInfo struct {
	processo string
	bytes    uint64
	// total
}

func (j *Janelinha) init(dst string) {
	// runtime.LockOSThread()
	Trace.Printf("\n Iniciando Janelinha %v", j.id)
	for processo := range j.entradaProcessos {
		filepath := processoPath(dst, processo.numStrImpuro)
		Trace.Printf("\n------\nJanelinha %v com processo %v (%s)", j.id, processo.numStrImpuro, filepath)
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			Trace.Printf("\nBaixar: %s", filepath)
			j.oleAPI.abreProcesso(j.id, processo)
			time.Sleep(900 * time.Millisecond)
			Trace.Printf(" %v - %v-0", j.id, processo.numStrImpuro)
			for {
				testa := j.oleAPI.paginaProcessoCarregou(j.id, processo)
				if testa {
					// fmt.Printf(" %v - %v-1-A", j.id, processo.numStrImpuro)
					break
				}
				time.Sleep(500 * time.Millisecond)
				// fmt.Printf(" %v - %v-1", j.id, processo.numStrImpuro)
			}
			Trace.Printf("\n Carregou ¨¨¨¨ %s ¨¨¨¨ %v", j.id, processo.numStrImpuro)
			time.Sleep(500 * time.Millisecond)
			// time.Sleep(2 * time.Second)
			j.oleAPI.paginaProcessoPatcheVaiProDownload(j.id, processo)
			time.Sleep(500 * time.Millisecond)
			for {
				testa := j.oleAPI.paginaDocumentosCarregou(j.id, processo)
				// fmt.Printf(" %v - %v-2", j.id, processo.numStrImpuro)
				if testa {
					// fmt.Printf(" %v - %v-2-A", j.id, processo.numStrImpuro)
					break
				}
				time.Sleep(500 * time.Millisecond)
			}
			Trace.Printf("\n Carregou+++ ¨¨¨¨ %s ¨¨¨¨ %v", j.id, processo.numStrImpuro)
			time.Sleep(500 * time.Millisecond)
			j.oleAPI.clicaParaGerarPDF(j.id, processo)

			for {
				href := j.oleAPI.getHrefStringOrNot(j.id, processo)
				if len(href) > 10 {
					// Trace.Println("\n\n\nHREF encontrado")
					Trace.Println(href)
					// MANDA GOROTINA DO DOWNLOAD
					cookies := j.oleAPI.getCookies()
					// Trace.Println("\n%s", cookies)
					j.downloadChannel <- &DownloadPayload{
						cookieStr: cookies,
						titlePDF:  href,
						dst:       filepath,
					}
					Trace.Printf("\n%s enviado ao canal", href)
					break
				}
				time.Sleep(500 * time.Millisecond)
			}
			time.Sleep(1 * time.Second)
		} else {
			Trace.Printf("\nArquivo já existe: %s", filepath)
			j.waitGroup.Done()
		}
		Trace.Printf("Fim do loop do processo %v (%s)", processo.numStrImpuro, filepath)

	}
	Trace.Printf("%s - die!", j.id)
}

func startDownloader(id int, ch chan *DownloadPayload, wg *sync.WaitGroup, cc chan bool, ci chan DownloadInfo) {
	Trace.Printf("\nIniciando Downloader %d ", id)
	for payload := range ch {
		Trace.Printf("\nDownloader %d recebeu download %s", id, payload.titlePDF)
		Downloader(payload, wg, cc, ci)
	}
}

func esperarDownloads(wg *sync.WaitGroup) {
	wg.Wait()
	Info.Println(`
===========================================================================
		Todos os processos enviados para download
===========================================================================`)
	Trace.Println(` === Todos os processos enviados para download === apos o wg.Wait()`)
}

func baixarProcessosDoEprocessoPrincipal(diretorioDownload string, num_janelinhas int, num_downloaders int, api *apiConn, wg *sync.WaitGroup) {
	// runtime.LockOSThread()
	Trace.Printf("-")
	if !strings.HasSuffix(diretorioDownload, `\`) {
		diretorioDownload = diretorioDownload + `\`
	}
	Trace.Printf("-")

	_, err := os.Stat(diretorioDownload)
	if err != nil {
		Trace.Panicf(`Diretório %s não válido.`, diretorioDownload)
	}
	Trace.Printf("-")

	Trace.Printf(`
-------------------------------------------------------------------------------
	Diretório de Download: %s
	Janelas Simultâneas: %d		Downloads Simultâneos: %d
-------------------------------------------------------------------------------
`, diretorioDownload, num_janelinhas, num_downloaders)
	Trace.Printf("-")

	chP := make(chan *Processo)
	chDownload := make(chan *DownloadPayload)
	chDownloadComplete := make(chan bool)
	chDownloadInfo := make(chan DownloadInfo)
	Trace.Printf("-")
	// api := instantiateNewAPIConn()

	num_procs := api.sendProcessosDaJanelaToChannel(chP)
	wg.Add(num_procs)
	Trace.Printf("-")
	if num_downloaders > 10 {
		num_downloaders = 10
	}
	Trace.Printf("-")
	for i := 0; i < num_downloaders; i++ {
		go startDownloader(i, chDownload, wg, chDownloadComplete, chDownloadInfo)
	}
	Trace.Printf("-")
	if num_janelinhas > 10 {
		num_janelinhas = 10
	}
	Trace.Printf("-")
	for i := 0; i < num_janelinhas; i++ {
		jan := Janelinha{nomesDeJanela[i], api, chP, wg, chDownload}
		go jan.init(diretorioDownload)
	}
	Trace.Printf("-")
	Info.Printf(" * %d processos encontrados na página * ", num_procs)
	Trace.Printf("%d processos encontrados na página", num_procs)
	Trace.Printf("-")

	go DownloadReporter(chDownloadInfo)

	for index := 0; index < num_procs; index++ {
		<-chDownloadComplete
		Info.Printf("%d download(s) completo(s) de %d", index+1, num_procs)
		Trace.Printf("%d download(s) completo(s) de %d", index+1, num_procs)
	}

	close(chP)
	close(chDownload)
	close(chDownloadComplete)
	close(chDownloadInfo)

	Info.Printf("\nFim dos downloads :)")

}

func DownloadReporter(ch chan DownloadInfo) {
	dados := make(map[string]uint64)
	for {
		pld := <-ch
		dados[pld.processo] = pld.bytes
		var tot uint64
		for k := range dados {
			tot += dados[k]
		}
		fmt.Printf("\r                                                                      ")
		fmt.Printf("\r%s [ %d ]", pld.processo, tot)
		time.Sleep(200 * time.Millisecond)
	}
}

func main() {

	SetUpLoggers(os.Stderr, os.Stdout)
	defaultDownloadFolder := getUserHomeDir() + `\Downloads\`

	var diretorioDownload string
	var num_janelinhas int
	var num_downloaders int
	var baixarProcessos bool
	var serveData bool

	Trace.Printf("-")

	flag.StringVar(&diretorioDownload, "pasta", defaultDownloadFolder, `Pasta de Destino dos Processos Baixados`)
	flag.IntVar(&num_janelinhas, "janelas", 3, `Número de janelas simultâneas que devem ser abertas`)
	flag.IntVar(&num_downloaders, "downloads", 5, `Número máximo de downloads simultâneos`)
	flag.BoolVar(&baixarProcessos, "baixar", true, `Baixar Processos na Pasta Indicada em -pasta`)
	flag.BoolVar(&serveData, "servir", false, `Servir dados`)

	Trace.Printf("-")

	flag.Parse()

	api := instantiateNewAPIConn()
	api.janelaEprocesso()
	api.patchWinPrincipal()

	time.Sleep(100 * time.Millisecond)
	wg := &sync.WaitGroup{}

	if baixarProcessos {
		Trace.Printf(">")
		go baixarProcessosDoEprocessoPrincipal(diretorioDownload, num_janelinhas, num_downloaders, api, wg)
	}

	if serveData {
		Trace.Printf(">")
		go esperarDownloads(wg)
		serveHttp(api, diretorioDownload)
	} else {
		time.Sleep(4 * time.Second)
		esperarDownloads(wg)
	}

}
