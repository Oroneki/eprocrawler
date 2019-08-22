package main

import (
	"strings"
	"time"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func sSIDAjanelaInit(api *apiConn) bool {
	trace.Println("bora abrir janelinha do SIDA...")
	if api.sidaIE != nil {
		trace.Println("já tem uma...")
		return true
	}
	unknown, _ := oleutil.CreateObject("InternetExplorer.Application")
	ie, _ := unknown.QueryInterface(ole.IID_IDispatch)
	oleutil.CallMethod(ie, "Navigate", "http://www3.pgfn.fazenda/pgfn/milenio/aplicativos2.asp")
	oleutil.PutProperty(ie, "Visible", true)
	trace.Println("abriu...")
	for {
		time.Sleep(250 * time.Millisecond)
		if oleutil.MustGetProperty(ie, "Busy").Val == 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	trace.Println("pegar info e metadata")
	document := oleutil.MustGetProperty(ie, "document").ToIDispatch()
	window := oleutil.MustGetProperty(document, "parentWindow").ToIDispatch()
	trace.Printf("\n document: %#v\n window:  %#v", document, window)
	trace.Println("pegou info e metadata... dormir um pokim e ver se o homi loga...")
	time.Sleep(3000 * time.Millisecond)
	for {
		logado := sSIDAcheckLogado(window)
		trace.Println("logado: ", logado)
		if logado {
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}

	api.sidaIE = ie
	api.sidaIEWindow = window
	return true
}

func sSIDAcheckLogado(ieWindow *ole.IDispatch) bool {

	trace.Println("--------")
	variant, err := oleutil.CallMethod(ieWindow, "eval", `window.location.pathname === "/PGFN/Milenio/PrincipalFrames.asp";`)
	if err != nil {
		trace.Println("miow sSIDAcheckLogado")
		trace.Printf("recebeu: %#v", ieWindow)
		return false
	}
	trace.Println("--------")
	return variant.Value().(bool)

}

func sidaVaiPraPaginaResultadoConsulta(api *apiConn, processo string) bool {
	trace.Println("--------")
	oleutil.CallMethod(api.sidaIEWindow, "Navigate", "http://www3.pgfn.fazenda/PGFN/Divida/Consulta/Inscricao/Cons11.asp")
	trace.Println("--------")
	time.Sleep(50 * time.Millisecond)
	trace.Println("--------")
	WaitIEWindow(api.sidaIE)
	trace.Println("--------")
	// for {

	// 	okkk, err := api.sidaIEWindow.CallMethod("eval", `window.location.pathname === "/PGFN/Divida/Consulta/Inscricao/Cons11.asp";`)
	// 	if err != nil {
	// 		time.Sleep(900 * time.Millisecond)
	// 		trace.Printf("\n :( - %v %#v\n", okkk, err)
	// 		continue
	// 	}
	// 	okk := okkk.Value().(bool)
	// 	trace.Println("-")
	// 	trace.Println("okk: ", okk)
	// 	if okk {
	// 		break
	// 	}
	// 	time.Sleep(100 * time.Millisecond)
	// }
	time.Sleep(5 * time.Second)
	trace.Println("\nCarregou consulta pro processo: ", processo)
	trace.Println("--------")

	time.Sleep(10 * time.Millisecond)
	oleutil.MustCallMethod(api.sidaIEWindow, "eval", `window.document.getElementById("op_numProcAntigo").click();`)
	trace.Println("--------")
	time.Sleep(100 * time.Millisecond)
	trace.Println("--------")
	oleutil.MustCallMethod(api.sidaIEWindow, "eval", `window.document.getElementsByTagName("input")[12].value = "`+processo+`";`)
	trace.Println("--------")

	valAntes := oleutil.MustCallMethod(api.sidaIEWindow, "eval", `window.document.getElementsByTagName("input")[20].value;`).Value().(string)
	trace.Println("--------")

	trace.Println("\nVal antes: ", valAntes)
	trace.Println("--------")

	oleutil.MustCallMethod(api.sidaIEWindow, "eval", `window.document.getElementsByName("ok")[0].click();`)
	trace.Println("--------")
	time.Sleep(50 * time.Millisecond)
	trace.Println("--------")
	WaitIEWindow(api.sidaIE)
	trace.Println("  injectar codigos...")

	evalOnIEWindowNoLock(api, jsPolyfills)
	trace.Println("  jsPolyfills...")
	evalOnIEWindowNoLock(api, jsPolyfillShimInjectScript)
	trace.Println("  jsPolyfillShimInjectScript...")
	WaitIEWindow(api.sidaIE)
	trace.Println("  evaluou busy...")
	waitForConditionOnIEWindow(api.sidaIEWindow, `(function () {
		try {
		  var el1 = Array.from(document.querySelectorAll('tr').filter(function (i) {
			return i.innerText && i.innerText.length > 1;
		  })).map(function (a) {
			return a.innerText;
		  });
		  var el2 = JSON.stringify(el1);
		} catch (e) {
		  return false;
		}
	  
		return true;
	  })();`)

	trace.Println("  RETORNOU DA FUNCAO sidaVaiPraPaginaResultadoConsulta :)")
	return true

}

func evalOnIEWindowNoLock(api *apiConn, pld string) {
	_, err := oleutil.CallMethod(
		api.sidaIEWindow,
		"eval",
		pld,
	)
	if err != nil {
		trace.Println("sida opa... erro no eval.")

	}

}

func grabConsultaInfo(api *apiConn, processo string) *grabConsultaProcessoSidaResult {
	keyValues := oleutil.MustCallMethod(api.sidaIEWindow, "eval", SidaKeyValuesConsulta).Value().(string)
	trace.Println("   grabConsultaInfo:\n", keyValues)

	valDepois, err := oleutil.CallMethod(api.sidaIEWindow, "eval", `window.document.getElementsByTagName("input")[20].value;`)
	if err != nil {
		trace.Println("     não achou o input... ta em outra pagina. (1 insc)")
		return &grabConsultaProcessoSidaResult{
			Json:                    keyValues,
			quantidadeIdentificador: UMA_INSCRICAO,
		}
	}
	valDepoisSS := valDepois.Value().(string)
	trace.Printf("    valDEpoisSS : %s", valDepoisSS)
	if strings.Contains(valDepoisSS, "FORAM LOCALIZADAS") {
		return &grabConsultaProcessoSidaResult{
			Json:                    "",
			quantidadeIdentificador: VARIAS_INSCRICOES,
		}
	} else {
		return &grabConsultaProcessoSidaResult{Json: "",
			quantidadeIdentificador: NENHUMA_INSCRICAO}
	}
}

func WaitIEWindow(ie *ole.IDispatch) {
	for {
		time.Sleep(250 * time.Millisecond)
		trace.Println("-")
		try, err := oleutil.GetProperty(ie, "Busy")
		if err != nil {
			trace.Println("ERRO! waitWindow - não pegou a propriedade Busy")
			trace.Printf("erro: %s", err)
		}
		if try.Val == 0 {
			trace.Println("Busy é 0 :)")
			break
		}
	}
	time.Sleep(100 * time.Millisecond)
}
