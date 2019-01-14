package main

import (
	"time"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func sSIDAjanelaInit(api *apiConn) bool {
	if api.sidaIE != nil {
		return true
	}
	unknown, _ := oleutil.CreateObject("InternetExplorer.Application")
	ie, _ := unknown.QueryInterface(ole.IID_IDispatch)
	oleutil.CallMethod(ie, "Navigate", "http://www3.pgfn.fazenda/pgfn/milenio/aplicativos2.asp")
	oleutil.PutProperty(ie, "Visible", true)
	for {
		time.Sleep(250 * time.Millisecond)
		if oleutil.MustGetProperty(ie, "Busy").Val == 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	document := oleutil.MustGetProperty(ie, "document").ToIDispatch()
	window := oleutil.MustGetProperty(document, "parentWindow").ToIDispatch()
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
	variant := oleutil.MustCallMethod(ieWindow, "eval", `window.location.pathname === "/PGFN/Milenio/PrincipalFrames.asp";`)
	trace.Println("--------")
	return variant.Value().(bool)

}

func sSIDAVaiPraConsulta(api *apiConn, processo string) string {
	trace.Println("--------")
	oleutil.CallMethod(api.sidaIEWindow, "Navigate", "http://www3.pgfn.fazenda/PGFN/Divida/Consulta/Inscricao/Cons11.asp")
	trace.Println("--------")
	time.Sleep(50 * time.Millisecond)
	trace.Println("--------")
	WaitIEWindow(api.sidaIE)
	trace.Println("--------")
	for {
		okk_, err := oleutil.CallMethod(api.sidaIEWindow, "eval", `window.location.pathname === "/PGFN/Divida/Consulta/Inscricao/Cons11.asp";`)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			trace.Println("- :( -")
			continue
		}
		okk := okk_.Value().(bool)
		trace.Println("-")
		trace.Println("okk: ", okk)
		if okk {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	trace.Println("\nCarregou consulta pro processo: ", processo)
	trace.Println("--------")

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
	trace.Println("--------")

	keyValues := oleutil.MustCallMethod(api.sidaIEWindow, "eval", SidaKeyValuesConsulta).Value().(string)
	trace.Println("\n\nResultado:\n", keyValues)

	valDepois, err := oleutil.CallMethod(api.sidaIEWindow, "eval", `window.document.getElementsByTagName("input")[20].value;`)
	if err != nil {
		trace.Println("\n\nERRO NO 'VAL DEPOIS'... ISSO EH BOM ???")
		return keyValues
	}
	valDepois_ := valDepois.Value().(string)
	return keyValues + "_CAMPO_20_||>" + valDepois_ + "\n"
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
