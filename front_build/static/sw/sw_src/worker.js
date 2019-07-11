let Database;
const JanelinhaProcessoInfoEventInfo = {
    tipo: "JANELINHA_INFO_PROCESSO",
    parseFunction: (pld) => {
        const striped_prc_events = pld.split('|@@|');
        const processo = striped_prc_events[0];
        const stripped_events = striped_prc_events[1].split('|##|');
        const events = [];
        for (const se of stripped_events) {
            if (se.length < 5) {
                return;
            }
            const splitted = se.split('|');
            const event = {
                nome_doc: splitted[0],
                ordem: parseInt(splitted[1], 10),
                pag_inicio: parseInt(splitted[2], 10),
                pag_fim: parseInt(splitted[3], 10),
                situacao: splitted[4],
                tamanho: parseInt(splitted[5], 10),
            };
            events.push(event);
        }
        return {
            processoImpuro: processo,
            infos: events
        };
    },
    callback: (input) => {
        if (Database === undefined) {
            console.error('Database is undefined');
            return;
        }
        const trans = Database.transaction(['documentos'], 'readwrite');
        const processo = input.processoImpuro
            .replace(/\s+\(\d+\)/g, '')
            .replace(/\D/, '');
        const obj = Object.assign({}, input, { processo });
        console.log(obj.processoImpuro, ' --> ', obj.processo);
        const put = trans.objectStore('documentos').put(obj);
        put.onsuccess = function () {
            console.log('salvo -> ', obj);
        };
    }
};
const DownloadConcludedEventInfo = {
    tipo: "DOWNLOAD_FINISHED",
    parseFunction: (pld) => {
        const striped = pld.split('|');
        return {
            processo_filename: striped[0],
            final_filepath: striped[1],
        };
    },
    callback: null
};
const JanelinhaEventInfo = {
    tipo: "JANELINHA_EVENT",
    parseFunction: (pld) => {
        const striped = pld.split('|');
        return {
            janId: striped[0],
            processoImpuro: striped[1],
            descricao: striped[2],
            fase: parseInt(striped[3], 10),
        };
    },
    callback: null
};
const DownloadBytesEventInfo = {
    tipo: "D_REPORTER",
    parseFunction: (pld) => {
        const striped = pld.split('|');
        return {
            processo_filename: striped[0],
            bytes: parseInt(striped[1], 10)
        };
    },
    callback: null
};
const IMALIVEInfo = {
    tipo: "im_alive",
    parseFunction: (pld) => "parsed",
    callback: null
};
const ALL_DOWNLOADS_FINISHED_EVINFO = {
    tipo: "ALL_DOWNLOADS_FINISHED",
    parseFunction: (pld) => null,
    callback: null
};
const WSEvents = {
    D_REPORTER: DownloadBytesEventInfo,
    im_alive: IMALIVEInfo,
    DOWNLOAD_FINISHED: DownloadConcludedEventInfo,
    ALL_DOWNLOADS_FINISHED: ALL_DOWNLOADS_FINISHED_EVINFO,
    JANELINHA_EVENT: JanelinhaEventInfo,
    JANELINHA_INFO_PROCESSO: JanelinhaProcessoInfoEventInfo,
};
const handleWebsocketPortHandler = (payload) => {
    console.log('WORKER: WS PORT RECEIVED :', payload);
    const ws = new WebSocket(payload);
    ws.onopen = function () {
        console.log('WS opened!');
    };
    ws.onmessage = function (e) {
        const ServerData = JSON.parse(e.data);
        //@ts-ignore
        const tratador = WSEvents[ServerData.tipo];
        const payloadParsed = tratador.parseFunction(ServerData.payload);
        //@ts-ignore
        postMessage({ tipo: ServerData.tipo, payload: payloadParsed });
        //
        if (tratador.callback !== null) {
            tratador.callback(payloadParsed);
        }
    };
};
const handleDataBaseConnect = (payload) => {
    const dbreq = indexedDB.open(payload.name, payload.version);
    dbreq.onupgradeneeded = function (e) {
        const database = dbreq.result;
        if (!database.objectStoreNames.contains('documentos')) {
            database.createObjectStore('documentos', { keyPath: 'numero' });
            console.log('updated db');
        }
        else {
        }
    };
    dbreq.onsuccess = function (e) {
        Database = dbreq.result;
        console.log('WORKER CONNECTED TO DB', payload.name, payload.version, 'Database:', Database);
    };
};
self.onmessage = e => {
    if (e.data && e.data.tipo) {
        switch (e.data.tipo) {
            // handle cases
            case "WEBSOCKET_PORT":
                handleWebsocketPortHandler(e.data.payload);
                break;
            case "DATABASE_INFO_START":
                handleDataBaseConnect(e.data.payload);
                break;
            default:
                console.error('ERR:', e.data.tipo, ' no handler on worker');
                break;
        }
    }
};
// @ts-ignore
setInterval(function () {
    // @ts-ignore
    postMessage({ tipo: 'teste', payload: { massa: 'eh bosta' } });
}, 9000);
