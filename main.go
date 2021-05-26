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
	atrasoSeconds    int64
	wsChannel        chan WebSocketMessage
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
		trace.Printf("\n Proceso %s %v", j.id, processo.numStrImpuro)
		j.wsChannel <- WebSocketMessage{
			Tipo:    "JANELINHA_EVENT",
			Payload: fmt.Sprintf("%s|%s|RECEBEU|1", j.id, processo.numStrImpuro),
		}
		filepath := processoPath(dst, processo.numStrImpuro)
		trace.Printf("\n------\nJanelinha %v com processo %v (%s)", j.id, processo.numStrImpuro, filepath)
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			trace.Printf("\nBaixar: %s", filepath)
			j.wsChannel <- WebSocketMessage{
				Tipo:    "JANELINHA_EVENT",
				Payload: fmt.Sprintf("%s|%s|VAI_ABRIR|2", j.id, processo.numStrImpuro),
			}
			time.Sleep(time.Duration(j.atrasoSeconds) * time.Second)
			j.oleAPI.abreProcesso(j.id, processo)
			j.wsChannel <- WebSocketMessage{
				Tipo:    "JANELINHA_EVENT",
				Payload: fmt.Sprintf("%s|%s|ABRIU|3", j.id, processo.numStrImpuro),
			}
			time.Sleep(time.Duration(j.atrasoSeconds) * time.Second)
			time.Sleep(250 * time.Millisecond)
			trace.Printf(" %v - %v-0", j.id, processo.numStrImpuro)
			for {
				testa := j.oleAPI.paginaProcessoCarregou(j.id, processo)
				if testa {
					// fmt.Printf(" %v - %v-1-A", j.id, processo.numStrImpuro)
					break
				}
				time.Sleep(8 * time.Second)
				// TODO : melhorar o testa carregou...
			}
			j.wsChannel <- WebSocketMessage{
				Tipo:    "JANELINHA_EVENT",
				Payload: fmt.Sprintf("%s|%s|CARREGOU_PROCESSO|4", j.id, processo.numStrImpuro),
			}
			trace.Println("\n\n\n ------------------ JANELINHA_INFO -------------------")
			rawStringInfo := j.oleAPI.getRawStringInfoFromProcessoJanelinha(j.id)
			if j.atrasoSeconds > 0 {
				info.Println("rawString returned - time to pause")
				time.Sleep(time.Duration(j.atrasoSeconds) * time.Second)
			}
			janelinhasInfo := parseJanelinhaInfo(rawStringInfo)
			trace.Printf("\n JANELINHA_INFO parseJanelinhaInfo -->  \n%s\n%#v\n\n", j.id, janelinhasInfo)
			formatAndSendToWebsocket(janelinhasInfo, processo.numStrImpuro, j.wsChannel)
			time.Sleep(time.Duration(j.atrasoSeconds) * time.Second)
			trace.Printf("\n Carregou ¨¨¨¨ %s ¨¨¨¨ %v", j.id, processo.numStrImpuro)
			time.Sleep(250 * time.Millisecond)
			// time.Sleep(2 * time.Second)
			j.oleAPI.paginaProcessoPatcheVaiProDownload(j.id, processo)
			j.wsChannel <- WebSocketMessage{
				Tipo:    "JANELINHA_EVENT",
				Payload: fmt.Sprintf("%s|%s|CLICOULINK_PROCESSO|5", j.id, processo.numStrImpuro),
			}
			time.Sleep(time.Duration(j.atrasoSeconds) * time.Second)

			time.Sleep(250 * time.Millisecond)
			for {
				testa := j.oleAPI.paginaDocumentosCarregou(j.id, processo)
				// fmt.Printf(" %v - %v-2", j.id, processo.numStrImpuro)
				if testa {
					// fmt.Printf(" %v - %v-2-A", j.id, processo.numStrImpuro)
					break
				}
				time.Sleep(500 * time.Millisecond)
			}
			time.Sleep(time.Duration(j.atrasoSeconds) * time.Second)
			j.wsChannel <- WebSocketMessage{
				Tipo:    "JANELINHA_EVENT",
				Payload: fmt.Sprintf("%s|%s|CARREGOU_DOCS|6", j.id, processo.numStrImpuro),
			}
			trace.Printf("\n Carregou+++ ¨¨¨¨ %s ¨¨¨¨ %v", j.id, processo.numStrImpuro)
			time.Sleep(250 * time.Millisecond)
			time.Sleep(time.Duration(j.atrasoSeconds) * time.Second)
			j.oleAPI.clicaParaGerarPDF(j.id, processo)
			j.wsChannel <- WebSocketMessage{
				Tipo:    "JANELINHA_EVENT",
				Payload: fmt.Sprintf("%s|%s|CLICKOU_GERAR_PDF|7", j.id, processo.numStrImpuro),
			}
			time.Sleep(time.Duration(j.atrasoSeconds) * time.Second)

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
					j.wsChannel <- WebSocketMessage{
						Tipo:    "JANELINHA_EVENT",
						Payload: fmt.Sprintf("%s|%s|ENVIOU_PRA_DOWNLOAD|8", j.id, processo.numStrImpuro),
					}
					break
				}
				time.Sleep(500 * time.Millisecond)
			}
			time.Sleep(1 * time.Second)
			time.Sleep(time.Duration(j.atrasoSeconds) * time.Second)
		} else {
			trace.Printf("\nArquivo já existe: %s", filepath)
			j.wsChannel <- WebSocketMessage{
				Tipo:    "JANELINHA_EVENT",
				Payload: fmt.Sprintf("%s|%s|EXISTE|0", j.id, processo.numStrImpuro),
			}
			time.Sleep(time.Duration(j.atrasoSeconds) * time.Second)

			j.waitGroup.Done()
		}
		j.wsChannel <- WebSocketMessage{
			Tipo:    "JANELINHA_EVENT",
			Payload: fmt.Sprintf("%s|%s|FIM|9", j.id, processo.numStrImpuro),
		}
		trace.Printf("Fim do loop do processo %v (%s)", processo.numStrImpuro, filepath)

	}
	trace.Printf("%s - die!", j.id)
	j.wsChannel <- WebSocketMessage{
		Tipo:    "JANELINHA_EVENT",
		Payload: fmt.Sprintf("%s|_|END_JANELINHA|10", j.id),
	}
}

func startDownloader(id int, ch chan *downloadPayload, wg *sync.WaitGroup, cc chan bool, ci chan downloadInfo, wsWriteChannel chan WebSocketMessage) {
	trace.Printf("\nIniciando Downloader %d ", id)
	for payload := range ch {
		trace.Printf("\nDownloader %d recebeu download %s", id, payload.titlePDF)
		downloader(payload, wg, cc, ci, wsWriteChannel)
	}
}

func esperarDownloads(wg *sync.WaitGroup) {
	time.Sleep(3 * time.Second)
	wg.Wait()
	info.Println(`
===========================================================================
		Todos os processos enviados para download
===========================================================================`)
	trace.Println(` === Todos os processos enviados para download === apos o wg.Wait()`)
}

func baixarProcessosDoEprocessoPrincipal(diretorioDownload string, numJanelinhas int, numDownloaders int, api *apiConn, wg *sync.WaitGroup, wsWrite chan WebSocketMessage, atraso int64) {
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
	trace.Printf("Aguardar  segundos...")
	time.Sleep(1 * time.Second)
	trace.Printf("Aguardou ... seguir")

	chP := make(chan *Processo)
	chDownload := make(chan *downloadPayload)
	chDownloadComplete := make(chan bool)
	chDownloadInfo := make(chan downloadInfo)
	trace.Printf("-")
	// api := instantiateNewAPIConn()

	num_procs := api.sendProcessosDaJanelaToChannel(chP)
	wg.Add(num_procs)
	trace.Printf("%d processos", num_procs)
	if numDownloaders > 10 {
		numDownloaders = 10
	}
	trace.Printf("-")
	for i := 0; i < numDownloaders; i++ {
		go startDownloader(i, chDownload, wg, chDownloadComplete, chDownloadInfo, wsWrite)
	}
	trace.Printf("-")
	if numJanelinhas > 10 {
		numJanelinhas = 10
	}
	trace.Printf("-")
	for i := 0; i < numJanelinhas; i++ {
		jan := janelinha{nomesDeJanela[i], api, chP, wg, chDownload, atraso, wsWrite}
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

	defer close(chP)
	defer close(chDownload)
	defer close(chDownloadComplete)
	defer close(chDownloadInfo)

	time.Sleep(3 * time.Second)

	wsWrite <- WebSocketMessage{
		Tipo:    "ALL_DOWNLOADS_FINISHED",
		Payload: "",
	}
	wg.Wait()

	info.Printf("\nFim dos downloads :)")

}

func DownloadReporter(ch chan downloadInfo, wsWrite chan WebSocketMessage) {
	trace.Printf("\nFim dos downloads :)")

	dados := make(map[string]uint64)
	for {
		pld, more := <-ch
		if !more {
			break
		}
		dados[pld.processo] = pld.bytes
		var tot uint64
		for k := range dados {
			tot += dados[k]
		}
		if pld.bytes%10 == 0 {
			go func() {
				wsWrite <- WebSocketMessage{
					Tipo:    "D_REPORTER",
					Payload: fmt.Sprintf("%s|%d", pld.processo, pld.bytes),
				}
			}()
			if pld.bytes%50 == 0 {
				fmt.Printf("\r                                                                      ")
				fmt.Printf("\r%16d [ %-20s ]", tot, pld.processo)
			}

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
	var atrasoJanelinha int

	trace.Printf("-")

	flag.StringVar(&diretorioDownload, "pasta", defaultDownloadFolder, `Pasta de Destino dos Processos Baixados`)
	flag.StringVar(&portServer, "porta", "9090", `Porta do Servidor`)
	flag.IntVar(&num_janelinhas, "janelas", 3, `Número de janelas simultâneas que devem ser abertas`)
	flag.IntVar(&num_downloaders, "downloads", 5, `Número máximo de downloads simultâneos`)
	flag.IntVar(&atrasoJanelinha, "atraso", 0, `Atraso da janelinha`)
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
			time.Sleep(60 * time.Second)
			WSChannelWrite <- WebSocketMessage{
				Tipo:    "im_alive",
				Payload: "",
			}
		}
	}()
	if baixarProcessos {
		trace.Printf(">")
		go baixarProcessosDoEprocessoPrincipal(diretorioDownload, num_janelinhas, num_downloaders, api, wg, WSChannelWrite, int64(atrasoJanelinha))
	}

	if serveData {
		trace.Printf(">")
		go esperarDownloads(wg)
		serveHttp(api, diretorioDownload, portServer, WSChannelWrite)
	} else {
		go func() {
			for msg := range WSChannelWrite {
				trace.Printf("drain '%s'", msg.Tipo)
			}
		}()
		time.Sleep(4 * time.Second)
		esperarDownloads(wg)
	}

}
