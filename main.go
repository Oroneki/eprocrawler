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

var nomesDeJanela = [10]string{
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

type janelinha struct {
	id               string
	oleAPI           *apiConn
	entradaProcessos <-chan *Processo
	waitGroup        *sync.WaitGroup
	downloadChannel  chan *downloadPayload
}

type downloadInfo struct {
	processo string
	bytes    uint64
	// total
}

func (j *janelinha) init(dst string) {
	// runtime.LockOSThread()
	trace.Printf("\n Iniciando Janelinha %v", j.id)
	for processo := range j.entradaProcessos {
		filepath := processoPath(dst, processo.numStrImpuro)
		trace.Printf("\n------\nJanelinha %v com processo %v (%s)", j.id, processo.numStrImpuro, filepath)
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			trace.Printf("\nBaixar: %s", filepath)
			j.oleAPI.abreProcesso(j.id, processo)
			time.Sleep(900 * time.Millisecond)
			trace.Printf(" %v - %v-0", j.id, processo.numStrImpuro)
			for {
				testa := j.oleAPI.paginaProcessoCarregou(j.id, processo)
				if testa {
					// fmt.Printf(" %v - %v-1-A", j.id, processo.numStrImpuro)
					break
				}
				time.Sleep(500 * time.Millisecond)
				// fmt.Printf(" %v - %v-1", j.id, processo.numStrImpuro)
			}
			trace.Printf("\n Carregou ¨¨¨¨ %s ¨¨¨¨ %v", j.id, processo.numStrImpuro)
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
			trace.Printf("\n Carregou+++ ¨¨¨¨ %s ¨¨¨¨ %v", j.id, processo.numStrImpuro)
			time.Sleep(500 * time.Millisecond)
			j.oleAPI.clicaParaGerarPDF(j.id, processo)

			for {
				href := j.oleAPI.getHrefStringOrNot(j.id, processo)
				if len(href) > 10 {
					// Trace.Println("\n\n\nHREF encontrado")
					trace.Println(href)
					// MANDA GOROTINA DO DOWNLOAD
					cookies := j.oleAPI.getCookies()
					// Trace.Println("\n%s", cookies)
					j.downloadChannel <- &downloadPayload{
						cookieStr: cookies,
						titlePDF:  href,
						dst:       filepath,
					}
					trace.Printf("\n%s enviado ao canal", href)
					break
				}
				time.Sleep(500 * time.Millisecond)
			}
			time.Sleep(1 * time.Second)
		} else {
			trace.Printf("\nArquivo já existe: %s", filepath)
			j.waitGroup.Done()
		}
		trace.Printf("Fim do loop do processo %v (%s)", processo.numStrImpuro, filepath)

	}
	trace.Printf("%s - die!", j.id)
}

func startDownloader(id int, ch chan *downloadPayload, wg *sync.WaitGroup, cc chan bool, ci chan downloadInfo) {
	trace.Printf("\nIniciando Downloader %d ", id)
	for payload := range ch {
		trace.Printf("\nDownloader %d recebeu download %s", id, payload.titlePDF)
		downloader(payload, wg, cc, ci)
	}
}

func esperarDownloads(wg *sync.WaitGroup) {
	wg.Wait()
	info.Println(`
===========================================================================
		Todos os processos enviados para download
===========================================================================`)
	trace.Println(` === Todos os processos enviados para download === apos o wg.Wait()`)
}

func baixarProcessosDoEprocessoPrincipal(diretorioDownload string, numJanelinhas int, numDownloaders int, api *apiConn, wg *sync.WaitGroup, wsWrite chan WebSocketMessage) {
	// runtime.LockOSThread()

	trace.Printf("-")

	_, err := os.Stat(diretorioDownload)
	if err != nil {
		trace.Panicf(`Diretório %s não válido.`, diretorioDownload)
	}
	trace.Printf("-")

	trace.Printf(`
-------------------------------------------------------------------------------
	Diretório de Download: %s
	Janelas Simultâneas: %d		Downloads Simultâneos: %d
-------------------------------------------------------------------------------
`, diretorioDownload, numJanelinhas, numDownloaders)
	trace.Printf("-")

	chP := make(chan *Processo)
	chDownload := make(chan *downloadPayload)
	chDownloadComplete := make(chan bool)
	chDownloadInfo := make(chan downloadInfo)
	trace.Printf("-")
	// api := instantiateNewAPIConn()

	num_procs := api.sendProcessosDaJanelaToChannel(chP)
	wg.Add(num_procs)
	trace.Printf("-")
	if numDownloaders > 10 {
		numDownloaders = 10
	}
	trace.Printf("-")
	for i := 0; i < numDownloaders; i++ {
		go startDownloader(i, chDownload, wg, chDownloadComplete, chDownloadInfo)
	}
	trace.Printf("-")
	if numJanelinhas > 10 {
		numJanelinhas = 10
	}
	trace.Printf("-")
	for i := 0; i < numJanelinhas; i++ {
		jan := janelinha{nomesDeJanela[i], api, chP, wg, chDownload}
		go jan.init(diretorioDownload)
	}
	trace.Printf("-")
	info.Printf(" * %d processos encontrados na página * ", num_procs)
	trace.Printf("%d processos encontrados na página", num_procs)
	trace.Printf("-")

	go DownloadReporter(chDownloadInfo, wsWrite)

	for index := 0; index < num_procs; index++ {
		<-chDownloadComplete
		info.Printf("\n%d download(s) completo(s) de %d", index+1, num_procs)
		trace.Printf("\n%d download(s) completo(s) de %d", index+1, num_procs)
	}

	close(chP)
	close(chDownload)
	close(chDownloadComplete)
	close(chDownloadInfo)

	info.Printf("\nFim dos downloads :)")

}

func DownloadReporter(ch chan downloadInfo, wsWrite chan WebSocketMessage) {
	dados := make(map[string]uint64)
	for {
		pld, more := <-ch
		if more == false {
			break
		}
		dados[pld.processo] = pld.bytes
		var tot uint64
		for k := range dados {
			tot += dados[k]
		}
		if pld.bytes%15 == 0 {
			fmt.Printf("\r                                                                      ")
			fmt.Printf("\r%16d [ %-20s ]", tot, pld.processo)
			go func() {
				wsWrite <- WebSocketMessage{
					Tipo:    "D_REPORTER",
					Payload: "_",
				}
			}()

		}

	}
	trace.Println("Encerrando DownloadReporter")
}

func main() {

	setUpLoggers(os.Stderr, os.Stdout)
	defaultDownloadFolder := getUserHomeDir() + `\Downloads\`

	var diretorioDownload string
	var portServer string
	var num_janelinhas int
	var num_downloaders int
	var baixarProcessos bool
	var serveData bool
	var sidaDesajuiza bool
	var injectCode bool

	trace.Printf("-")

	flag.StringVar(&diretorioDownload, "pasta", defaultDownloadFolder, `Pasta de Destino dos Processos Baixados`)
	flag.StringVar(&portServer, "porta", "9090", `Porta do Servidor`)
	flag.IntVar(&num_janelinhas, "janelas", 3, `Número de janelas simultâneas que devem ser abertas`)
	flag.IntVar(&num_downloaders, "downloads", 5, `Número máximo de downloads simultâneos`)
	flag.BoolVar(&baixarProcessos, "baixar", true, `Baixar Processos na Pasta Indicada em -pasta`)
	flag.BoolVar(&serveData, "servir", false, `Servir dados`)
	flag.BoolVar(&sidaDesajuiza, "sida_desajuiza", false, `Iniciar desajuizamento de varios`)
	flag.BoolVar(&injectCode, "inject_code", false, `Injetar código no I.E.`)

	trace.Printf("-")

	if !strings.HasSuffix(diretorioDownload, `\`) {
		diretorioDownload = diretorioDownload + `\`
	}

	flag.Parse()

	trace.Printf(`
diretorioDownload -->  %v
portServer -->  %v
num_janelinhas -->  %v
num_downloaders -->  %v
baixarProcessos -->  %v
serveData -->  %v
sidaDesajuiza -->  %v
injectCode -->  %v
	`, diretorioDownload,
		portServer,
		num_janelinhas,
		num_downloaders,
		baixarProcessos,
		serveData,
		sidaDesajuiza,
		injectCode)

	api := instantiateNewAPIConn()

	if sidaDesajuiza {
		trace.Printf("- sida")
		api.grabSidaWindow()
	} else {
		api.janelaEprocesso()
		api.patchWinPrincipal()
	}

	time.Sleep(100 * time.Millisecond)
	wg := &sync.WaitGroup{}

	WSChannelWrite := make(chan WebSocketMessage)
	go func() {
		for {
			time.Sleep(5 * time.Second)
			WSChannelWrite <- WebSocketMessage{
				Tipo:    "im_alive",
				Payload: "",
			}
		}
	}()
	if baixarProcessos {
		trace.Printf(">")
		go baixarProcessosDoEprocessoPrincipal(diretorioDownload, num_janelinhas, num_downloaders, api, wg, WSChannelWrite)
	}

	if serveData {
		trace.Printf(">")
		go esperarDownloads(wg)
		serveHttp(api, diretorioDownload, portServer, WSChannelWrite)
	} else {
		time.Sleep(4 * time.Second)
		esperarDownloads(wg)
	}

}
