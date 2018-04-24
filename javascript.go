package main

const (
	JShackedobj = `window.hacked_visualizarProcesso = function(TARGET_JANELA, psNumeroProcesso, psNumeroEquipeAtividade, psNomeEquipeAtual, psNomeAtividadeAtual) {
		abrirPopupProcesso(TARGET_JANELA, { 'psNumeroProcesso': psNumeroProcesso, 'psNumeroEquipeAtividade': psNumeroEquipeAtividade, 'psNomeEquipeAtual': psNomeEquipeAtual, 'psNomeAtividadeAtual': psNomeAtividadeAtual });
	};`

	JSabrirJanela = `window.abrirPopupProcesso = function (TARGET, params) {
		var parametros = "psAcao=exibir&" ;
		var lnHeight = screen.height;
		var lnWidth = screen.width;
		parametros += "psNumeroProcesso=" + params["psNumeroProcesso"] ;
		parametros += "&psNumeroEquipeAtividade="+params["psNumeroEquipeAtividade"] ;
		parametros += "&psNomeEquipeAtual=" + params["psNomeEquipeAtual"] ;
		parametros += "&psNomeAtividadeAtual=" + params["psNomeAtividadeAtual"] ;
		var posJan = Object.keys(window.oro_obj).length * 25 + 5;
		if (window.oro_obj[TARGET]) {
			window.oro_obj[TARGET] = window.open("about:blank", TARGET, "width="+(lnWidth-300)+",height="+(lnHeight-270)+",scrollbars=no,resizable=yes,left=100,top=100");
		}
		window.oro_obj[TARGET] = window.open("/ControleVisualizacaoProcesso.asp?" + parametros, TARGET, "width="+(lnWidth-360)+",height="+(lnHeight-298)+",scrollbars=no,resizable=yes,left="+(posJan)+",top="+(posJan));
		console.log('Abrindo processo ', params["psNumeroProcesso"], ' na janela ', TARGET);
	};`

	JSpatchInicial = `window.test_load_processo = function(TARGET, processo_str) {
		var t_1_b_ =  window.oro_obj[TARGET] && window.oro_obj[TARGET].document && window.oro_obj[TARGET].document.readyState === "complete";
		console.log(t_1_b_, "  readystate")
		if (t_1_b_) {
			console.log('testa_readystate');
			return window.oro_obj[TARGET].document.title.replace(/\D/g, "").indexOf(processo_str.replace(/\D/g, "").slice(0, 10)) > -1 &&
			window.oro_obj[TARGET].document.getElementsByTagName("area") && 
			window.oro_obj[TARGET].document.getElementsByTagName("area").length > 1;
		} else {
			return false;
		};    
	};

	window.clica_pra_gerar_pdf = function(TARGET) {
		window.oro_obj[TARGET].document.getElementById("chkMetaDados").checked = false;
		window.oro_obj[TARGET].document.getElementById("chkMetaDados").checked = false;
		var naoPag = window.oro_obj[TARGET].document.getElementById("chkNaoPaginavel");
		if (naoPag) {
			naoPag.checked = false;
			naoPag.checked = false;
		};				
		window.oro_obj[TARGET].document.getElementById("imgPdf").click();
		window.oro_obj[TARGET].document.getElementById("imgPdf").click();
	};

	window.get_download_href_or_false = function(TARGET) {
		var string_get_donload_retiurn = window.oro_obj[TARGET].document.getElementById("linkDownloadParte1") || {title: ""};
		if (string_get_donload_retiurn.title === "") {
			if (!window.oro_obj[TARGET].document.getElementById("imgPdf").disabled) {
				window.oro_obj[TARGET].document.getElementById("imgPdf").click()
			};
		};
		return string_get_donload_retiurn.title;
		};		
		
		`

	JStestLoadPagina = `window.test_load_pagina_download = function(TARGET) {
		var t_2_a_ = window.oro_obj[TARGET].document && window.oro_obj[TARGET].document.readyState === "complete";
		if (t_2_a_) {
			return window.oro_obj[TARGET].document.getElementsByTagName("img") && window.oro_obj[TARGET].document.getElementsByTagName("img").length > 1;
		} else {
			return false
		}
	};`

	JSgetJsonData = `
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
						value = value.replace(/[^\w]/g, " ");

						//if (value === "" || value === "-") {
						//	console.log("+Descartando: ", value);
						//	continue
						//}
						
						mapinha[map_headers[k]] = value;
					}
					console.log(Object.keys(mapinha));
					key_p = mapinha["n_mero processo"].replace(/\(\d+\)/g, "");
					key_p = key_p.replace(/\D/g, "");
					//console.log("key_p ", key_p);
					map_final[key_p] = mapinha;
				}
				
			}
			map_final["__META__"] = {codEquipe: document.getElementById("hidEquipeSelecionadaCaixaTrabalho").value};
			return jQuery.stringify(map_final);
			};
	`

	JSjqueryStringify = `jQuery.extend({
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
						`
	JSpaginaProcessoPatch = `					
	window.oro_obj.%s.obterMultiplosDocumentos = function () {
	var psOperacao = "C";
	var psEscopo = "I";
	var winbind = window.oro_obj.%s;
	var gsNumeroProcesso = winbind.document.getElementById('hidNumeroProcesso').value;
	var gsNumeroEquipeAtividade = winbind.document.getElementById('hidNumeroEquipeAtividade').value;
	var lsSituacaoProcesso = winbind.document.getElementById('hidSituacaoProcesso').value;					
	var gsResponsavelProcesso = winbind.document.getElementById('hidResponsavelProcesso').value;
	var gsNomeEquipe = winbind.document.getElementById('hidNomeEquipeAtual').value;
	var gsNomeAtividade = winbind.document.getElementById('hidNomeAtividadeAtual').value;
	var gsNumeroProcessoFormatado = winbind.document.getElementById('hidNumeroProcessoFormatado').value;					
	//Pega só as chaves dos documentos
	var gaDocumentosSelecionados = winbind.gsDocSelecionados.split("@");
	//valida se selecionou no máximo 1000 documentos
	if (gaDocumentosSelecionados.length>1000){
		alert("Selecione no máximo 1000 documentos.");
		return;
	}			

	var laDocumento = new Array();
	var laNumeroDocumentos = new Array();
	for (var lnIndice = 0; lnIndice < gaDocumentosSelecionados.length; lnIndice++){

		laDocumento = gaDocumentosSelecionados[lnIndice].split("|");
		laNumeroDocumentos[lnIndice] = laDocumento[0];
	}
	var lsURL = "ControleMultiplosDocumentos.asp?psAcao=apresentarPagina&psNumeroProcesso=" + gsNumeroProcesso + "&paDocSelecionados=" +
			laNumeroDocumentos + "&psNumeroEquipeAtividade=" + gsNumeroEquipeAtividade + "&psOperacao=" + psOperacao + "&psEscopo=" + psEscopo +
			"&psSituacaoProcesso=" + lsSituacaoProcesso + "&psResponsavelProcesso=" + gsResponsavelProcesso + "&psNomeEquipe=" + gsNomeEquipe + 
			"&psNomeAtividade=" + gsNomeAtividade + "&psNumeroProcessoFormatado=" + gsNumeroProcessoFormatado;

	console.log(winbind.name)
	window.open(lsURL, winbind.name);
}();				
				
`
)
