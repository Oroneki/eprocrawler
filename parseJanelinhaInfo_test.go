package main

import (
	"fmt"
	"os"
	"testing"
)

func TestParseJanelinhaInfo(t *testing.T) {
	setUpLoggers(os.Stderr, os.Stdout)
	StringTest := ` Sequencial@@@@@function onclick()
{
obterIndice(false)
}@@@@@

@@@@@@@@@@
false@@@@@false@@@@@false

@@@@@@@@@@
Ficha de Identificação@@@@@function onclick()
{
exibirDocumento(41944042, 'Ficha de Identificação', 07|54|369|559, 'CONFIRMADO',1)
}@@@@@function onmouseover()
{
exibirDivInfoDocumento(1, 'Ficha de Identificação','1@##@Fl. 1 a 1', 'CONFIRMADO', '', '0,020 MB')
}

@@@@@@@@@@
Telas e Extratos@@@@@function onclick()
{
exibirDocumento(41944058, 'Telas e Extratos', 07|52|345|490, 'AUTENTICADO',2)
}@@@@@function onmouseover()
{
exibirDivInfoDocumento(1, 'Telas e Extratos','2@##@Fl. 2 a 2', 'AUTENTICADO', '', '0,018 MB')
}

@@@@@@@@@@
Despacho de Encaminhamento@@@@@function onclick()
{
exibirDocumento(41944251, 'Despacho de Encaminhamento', 07|51|339|474, 'CONFIRMADO',3)
}@@@@@function onmouseover()
{
exibirDivInfoDocumento(1, 'Despacho de Encaminhamento','3@##@Fl. 3 a 3', 'CONFIRMADO', '', '0,020 MB')
}

@@@@@@@@@@
Despacho de Encaminhamento@@@@@function onclick()
{
exibirDocumento(41977950, 'Despacho de Encaminhamento', 07|51|339|474, 'CONFIRMADO',4)
}@@@@@function onmouseover()
{
exibirDivInfoDocumento(1, 'Despacho de Encaminhamento','4@##@Fl. 4 a 4', 'CONFIRMADO', '', '0,020 MB')
}

@@@@@@@@@@
Telas e Extratos@@@@@function onclick()
{
exibirDocumento(43547530, 'Telas e Extratos', 07|52|345|490, 'AUTENTICADO',5)
}@@@@@function onmouseover()
{
exibirDivInfoDocumento(1, 'Telas e Extratos','5@##@Fl. 5 a 5', 'AUTENTICADO', '', '0,015 MB')
}

@@@@@@@@@@
Telas e Extratos@@@@@function onclick()
{
exibirDocumento(43547531, 'Telas e Extratos', 07|52|345|490, 'AUTENTICADO',6)
}@@@@@function onmouseover()
{
exibirDivInfoDocumento(1, 'Telas e Extratos','6@##@Fl. 6 a 6', 'AUTENTICADO', '', '0,015 MB')
}

@@@@@@@@@@
Telas e Extratos@@@@@function onclick()
{
exibirDocumento(43547532, 'Telas e Extratos', 07|52|345|490, 'AUTENTICADO',7)
}@@@@@function onmouseover()
{
exibirDivInfoDocumento(1, 'Telas e Extratos','7@##@Fl. 7 a 7', 'AUTENTICADO', '', '0,015 MB')
}

@@@@@@@@@@
Telas e Extratos@@@@@function onclick()
{
exibirDocumento(43547535, 'Telas e Extratos', 07|52|345|490, 'AUTENTICADO',8)
}@@@@@function onmouseover()
{
exibirDivInfoDocumento(1, 'Telas e Extratos','8@##@Fl. 8 a 8', 'AUTENTICADO', '', '0,015 MB')
}

@@@@@@@@@@
Despacho de Encaminhamento@@@@@function onclick()
{
exibirDocumento(43547676, 'Despacho de Encaminhamento', 07|51|339|474, 'CONFIRMADO',9)
}@@@@@function onmouseover()
{
exibirDivInfoDocumento(1, 'Despacho de Encaminhamento','9@##@Fl. 9 a 9', 'CONFIRMADO', '', '0,020 MB')
}

@@@@@@@@@@
Telas e Extratos@@@@@function onclick()
{
exibirDocumento(43869351, 'Telas e Extratos', 07|52|345|490, 'AUTENTICADO',10)
}@@@@@function onmouseover()
{
exibirDivInfoDocumento(1, 'Telas e Extratos','10@##@Fl. 10 a 12', 'AUTENTICADO', '', '0,083 MB')
}

@@@@@@@@@@
Telas e Extratos@@@@@function onclick()
{
exibirDocumento(43869356, 'Telas e Extratos', 07|52|345|490, 'AUTENTICADO',13)
}@@@@@function onmouseover()
{
exibirDivInfoDocumento(1, 'Telas e Extratos','11@##@Fl. 13 a 15', 'AUTENTICADO', '', '0,082 MB')
}

@@@@@@@@@@
Telas e Extratos@@@@@function onclick()
{
exibirDocumento(43869358, 'Telas e Extratos', 07|52|345|490, 'AUTENTICADO',16)
}@@@@@function onmouseover()
{
exibirDivInfoDocumento(1, 'Telas e Extratos','12@##@Fl. 16 a 18', 'AUTENTICADO', '', '0,087 MB')
}

@@@@@@@@@@
Telas e Extratos@@@@@function onclick()
{
exibirDocumento(43869360, 'Telas e Extratos', 07|52|345|490, 'AUTENTICADO',19)
}@@@@@function onmouseover()
{
exibirDivInfoDocumento(1, 'Telas e Extratos','13@##@Fl. 19 a 21', 'AUTENTICADO', '', '0,083 MB')
}

@@@@@@@@@@
Despacho de Encaminhamento@@@@@function onclick()
{
exibirDocumento(43869898, 'Despacho de Encaminhamento', 07|51|339|474, 'CONFIRMADO',22)
}@@@@@function onmouseover()
{
exibirDivInfoDocumento(1, 'Despacho de Encaminhamento','14@##@Fl. 22 a 22', 'CONFIRMADO', '', '0,020 MB')
}

@@@@@@@@@@
Despacho de Encaminhamento@@@@@function onclick()
{
exibirDocumento(45449352, 'Despacho de Encaminhamento', 07|51|339|474, 'CONFIRMADO',23)
}@@@@@function onmouseover()
{
exibirDivInfoDocumento(1, 'Despacho de Encaminhamento','15@##@Fl. 23 a 23', 'CONFIRMADO', '', '0,020 MB')
}
`
	res := parseJanelinhaInfo(StringTest)
	fmt.Println(res, len(res))

}
