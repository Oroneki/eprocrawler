function retornoProcessarDocumentos(retorno, tipoExibicao, idParte, lsTamanhoParte) {
		
    //Se ocorreram erros na geração/gravação do arquivo
    if (retorno.indexOf("<erro>") >= 0){
        falhaProcessarDocumentos(retorno, idParte, lsTamanhoParte);
    }else{

        var lsParteLink = "";
    //remover mensagem "Processando..." da parte atual
    $("#cabecalhoParte"+idParte).empty();
    
    //atualiza accordion com link: retorno
    $("#cabecalhoParte"+idParte).off("click"); //limpa a pilha de eventos onclick
        lsParteLink = ($("h3[name^=parte]").length > 1 ) ? "da " + idParte + "a. parte" : "";
    if (tipoExibicao == "Zip") {
            $("#cabecalhoParte"+idParte).append("&nbsp;&nbsp;<a href='#' title='"+retorno.split('=')[1]+"' id='linkDownloadParte"+idParte+"'>Clique aqui para obter o arquivo ZIP gerado com os documentos " + lsParteLink + " - " + lsTamanhoParte + " MB</a>");
        $("#cabecalhoParte"+idParte).click(function(event){
                trocarCorLink('linkDownloadParte'+idParte);
                openDialog(retorno, 500, 180, false, ", scrollbars=yes, resizable=yes", 0, 0 ,'Salvando Arquivos',true); 
                //window.location.href = retorno; 
                ativabotaoTornarFisico(idParte)
            return false;
        });
    } else {  //pdf
            $("#cabecalhoParte"+idParte).append("&nbsp;&nbsp;<a href='#' title='"+retorno.split('=')[1]+"' id='linkDownloadParte"+idParte+"'>Clique aqui para obter o arquivo PDF gerado " + lsParteLink + " - " + lsTamanhoParte + " MB</a>");
        $("#cabecalhoParte"+idParte).click(function(event){
                trocarCorLink('linkDownloadParte'+idParte);
            window._SECRET_URL_ = retorno;
            console.log(retorno);
            window.open(retorno, this.name, "toolbar=no, resizable=yes");
                ativabotaoTornarFisico(idParte)
            return false;
        });
    }
        lnTemLink++
    }
    
    var lnProxParte = idParte+1;
    var lsPararProcessamento = $("#hidPararProcessamento").val();			
    
    //Verifica se é para parar e se tem a próxima parte
    if (lsPararProcessamento=="False" && ($("h3[id="+lnProxParte+"]")).length ){
        processarDocumentos(tipoExibicao, lnProxParte);
    } else {
        habilitarBotoesPlay();
        $("#hidPararProcessamento").val('False');
    }		
}

jQuery.extend({
    stringify  : function stringify(obj) {         
        if ("JSON" in window) {
            return JSON.stringify(obj);
        }

        var t = typeof (obj);
        if (t != "object" || obj === null) {
            // simple data type
            if (t == "string") obj = '"' + obj + '"';

            return String(obj);
        } else {
            // recurse array or object
            var n, v, json = [], arr = (obj && obj.constructor == Array);

            for (n in obj) {
                v = obj[n];
                t = typeof(v);
                if (obj.hasOwnProperty(n)) {
                    if (t == "string") {
                        v = '"' + v + '"';
                    } else if (t == "object" && v !== null){
                        v = jQuery.stringify(v);
                    }

                    json.push((arr ? "" : '"' + n + '":') + String(v));
                }
            }

            return (arr ? "[" : "{") + String(json) + (arr ? "]" : "}");
        }
    }
});

window.getJsonData = function() {
    var td_regex = /ddrivetip\(\'(.*?)\'\,/m;
    var tableCol = document.getElementsByTagName("table");
    var num_campos;
    var headers;
    var map_headers = {};
    var map_final = {};
    console.log(tableCol.length);
    for (i = 0; i < tableCol.length; i++) {
        if (tableCol[i].id !== "tblProcessos") {
            continue;
        }
        headers = tableCol[i].getElementsByTagName("th");
        num_campos = headers.length;
        console.log(num_campos, " campos");
        for (n = 0; n < num_campos; n++) {
            //console.log(n, " -> ", headers[n].innerText);
            map_headers[n] = headers[n].innerText.replace(/[^a-zA-Z ]/g, "_").toLowerCase().replace(/[^a-zA-Z0-9\_]+$/gm, "");
            //console.log(n, " -> ", map_headers[n]);
        }
        var trs = tableCol[i].getElementsByTagName("tr");
        var trlen = trs.length;
        //console.log();
        //console.log("trlen", trlen);
        //console.log(trs);
        for (t = 0; t < trlen; t++) {
            var tds = trs[t].getElementsByTagName("td");
            var mapinha = {};
            if (tds.length !== num_campos) { continue };
            //console.log(t, tds);
            for (k = 0; k < tds.length; k++) {
                //console.log("---------------------------------------");
                var value
                if (!tds[k].onmouseover) {
                    value = tds[k].innerText
                } else {
                    var match = td_regex.exec(tds[k].outerHTML);
                    value = match[1];
                }
                //console.log(k, " - ", map_headers[k], " --> ", value);
                value = value;
                if (value === "" || value === "-") {
                    console.log("+Descartando: ", value);
                    continue
                }
                mapinha[map_headers[k]] = value;
            }
            console.log(Object.keys(mapinha));
            key_p = mapinha["n_mero processo"].replace(/\(\d+\)/g, "");
            key_p = key_p.replace(/\D/g, "");
            //console.log("key_p ", key_p);
            map_final[key_p] = mapinha;
        }

    }
    return jQuery.stringify(map_final);
};

trans = {
    "todos": "Todos",
    "informa__es": "Informações",
    "indicadores": "Indicadores",
    "n_mero processo": "Número Processo",
    "ni contribuinte": "NI Contribuinte",
    "data entrada atividade": "Data Entrada Atividade",
    "nome respons_vel": "Nome Responsável",
    "nome equipe _ltima": "Nome Equipe Última",
    "assunto comprot": "Assunto Comprot",
    "cpf respons_vel _ltimo": "CPF Responsável Último",
    "nome atividade _ltima": "Nome Atividade Última",
    "cpf respons_vel atual": "CPF Responsável Atual",
    "indicador dossi_": "Indicador Dossiê",
    "valor atualizado da inscri__o": "Valor Atualizado da Inscrição",
    "assuntos_objetos": "Assuntos Objetos",
    "n_mero do requerimento _sicar_pgfn_": "Número Requerimento SICAR",
    "nome _ltimo documento confirmado": "Último Documento",
    "nome unidade _ltima": "Última Atividade",
    "indicador grande devedor": "GD",
    "nome contribuinte": "Contribuinte",
    "n_mero de inscri__o": "Número Inscrição",
    "situa__o da inscri__o": "Situação Inscrição",
}