package main

const (
	JSconsole = `window.__TYEMPKHGFJHGFJGFJGFJHGFURYDUYBN__ = function(){
		if (window.console) {
		  console.log('ja tem console');
		  return;		
	   };
	   window.console = {
		   log: function(a){return a;}	   
			   };
	   }();`

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
		var posJan = Object.keys(window.oro_obj).length * 35 + 5;
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

	JSUnicodeHandle = `window.mapaCharNumberToUnicode = {
		33: "\u0021",
		34: "\u0022",
		35: "\u0023",
		36: "\u0024",
		37: "\u0025",
		38: "\u0026",
		39: "\u0027",
		40: "\u0028",
		41: "\u0029",
		42: "\u002A",
		43: "\u002B",
		44: "\u002C",
		45: "\u002D",
		46: "\u002E",
		47: "\u002F",
		48: "\u0030",
		49: "\u0031",
		50: "\u0032",
		51: "\u0033",
		52: "\u0034",
		53: "\u0035",
		54: "\u0036",
		55: "\u0037",
		56: "\u0038",
		57: "\u0039",
		58: "\u003A",
		59: "\u003B",
		60: "\u003C",
		61: "\u003D",
		62: "\u003E",
		63: "\u003F",
		64: "\u0040",
		65: "\u0041",
		66: "\u0042",
		67: "\u0043",
		68: "\u0044",
		69: "\u0045",
		70: "\u0046",
		71: "\u0047",
		72: "\u0048",
		73: "\u0049",
		74: "\u004A",
		75: "\u004B",
		76: "\u004C",
		77: "\u004D",
		78: "\u004E",
		79: "\u004F",
		80: "\u0050",
		81: "\u0051",
		82: "\u0052",
		83: "\u0053",
		84: "\u0054",
		85: "\u0055",
		86: "\u0056",
		87: "\u0057",
		88: "\u0058",
		89: "\u0059",
		90: "\u005A",
		91: "\u005B",
		92: "\u005C",
		93: "\u005D",
		94: "\u005E",
		95: "\u005F",
		96: "\u0060",
		97: "\u0061",
		98: "\u0062",
		99: "\u0063",
		100: "\u0064",
		101: "\u0065",
		102: "\u0066",
		103: "\u0067",
		104: "\u0068",
		105: "\u0069",
		106: "\u006A",
		107: "\u006B",
		108: "\u006C",
		109: "\u006D",
		110: "\u006E",
		111: "\u006F",
		112: "\u0070",
		113: "\u0071",
		114: "\u0072",
		115: "\u0073",
		116: "\u0074",
		117: "\u0075",
		118: "\u0076",
		119: "\u0077",
		120: "\u0078",
		121: "\u0079",
		122: "\u007A",
		123: "\u007B",
		124: "\u007C",
		125: "\u007D",
		126: "\u007E",
		161: "\u00A1",
		162: "\u00A2",
		163: "\u00A3",
		164: "\u00A4",
		165: "\u00A5",
		166: "\u00A6",
		167: "\u00A7",
		168: "\u00A8",
		169: "\u00A9",
		170: "\u00AA",
		171: "\u00AB",
		172: "\u00AC",
		174: "\u00AE",
		175: "\u00AF",
		176: "\u00B0",
		177: "\u00B1",
		178: "\u00B2",
		179: "\u00B3",
		180: "\u00B4",
		181: "\u00B5",
		182: "\u00B6",
		183: "\u00B7",
		184: "\u00B8",
		185: "\u00B9",
		186: "\u00BA",
		187: "\u00BB",
		188: "\u00BC",
		189: "\u00BD",
		190: "\u00BE",
		191: "\u00BF",
		192: "\u00C0",
		193: "\u00C1",
		194: "\u00C2",
		195: "\u00C3",
		196: "\u00C4",
		197: "\u00C5",
		198: "\u00C6",
		199: "\u00C7",
		200: "\u00C8",
		201: "\u00C9",
		202: "\u00CA",
		203: "\u00CB",
		204: "\u00CC",
		205: "\u00CD",
		206: "\u00CE",
		207: "\u00CF",
		208: "\u00D0",
		209: "\u00D1",
		210: "\u00D2",
		211: "\u00D3",
		212: "\u00D4",
		213: "\u00D5",
		214: "\u00D6",
		215: "\u00D7",
		216: "\u00D8",
		217: "\u00D9",
		218: "\u00DA",
		219: "\u00DB",
		220: "\u00DC",
		221: "\u00DD",
		222: "\u00DE",
		223: "\u00DF",
		224: "\u00E0",
		225: "\u00E1",
		226: "\u00E2",
		227: "\u00E3",
		228: "\u00E4",
		229: "\u00E5",
		230: "\u00E6",
		231: "\u00E7",
		232: "\u00E8",
		233: "\u00E9",
		234: "\u00EA",
		235: "\u00EB",
		236: "\u00EC",
		237: "\u00ED",
		238: "\u00EE",
		239: "\u00EF",
		240: "\u00F0",
		241: "\u00F1",
		242: "\u00F2",
		243: "\u00F3",
		244: "\u00F4",
		245: "\u00F5",
		246: "\u00F6",
		247: "\u00F7",
		248: "\u00F8",
		249: "\u00F9",
		250: "\u00FA",
		251: "\u00FB",
		252: "\u00FC",
		253: "\u00FD",
		254: "\u00FE",
		255: "\u00FF"
	};
	
	window.AnsiToUnicode = function (word) {
		var chars = Array.from(word);
		var newC = chars.map(function (c) {
			if (mapaCharNumberToUnicode[c.charCodeAt(0)]) {
				// console.log(c,'   ', c.charCodeAt(0), '    ', mapaCharNumberToUnicode[c.charCodeAt(0)]);
				return mapaCharNumberToUnicode[c.charCodeAt(0)];
			}
			console.log(c, ' <> ', c.charCodeAt(0), ' < ', word);
			return " ";
		});
		// console.log(newC);
		return newC.join('');
	
	};`

	JSgetJsonData = `window.getJsonData = function () {
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
				//map_headers[n] = headers[n].innerText.replace(/[^a-zA-Z ]/g, "_").toLowerCase().replace(/[^a-zA-Z0-9\_]+$/gm, "");
				map_headers[n] = window.AnsiToUnicode(headers[n].innerText.strip());
				// console.log(n, " -> ", map_headers[n]);
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
					//value = value.replace(/[^\w\(\)]/g, " ");
					value = value.replace("\\'", "'");
					value = window.AnsiToUnicode(value);
	
					if (value === "" || value === "-" || value === "\u002D" || value === "    ") {
						continue
					}
	
					mapinha[map_headers[k]] = value;
				}
				// console.log(Object.keys(mapinha));
				key_p = mapinha["Número Processo"].replace(/\(\d+\)/g, "");
				// console.log('mapinha["Número Processo"]', mapinha["Número Processo"]);
				key_p = key_p.replace(/\D/g, "");
				//console.log(' ->', key_p);
				//console.log("key_p ", key_p);
				map_final[key_p] = mapinha;
			}
	
		}
		map_final["__META__"] = { codEquipe: document.getElementById("hidEquipeSelecionadaCaixaTrabalho").value };
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
