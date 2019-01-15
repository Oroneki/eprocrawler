package main

import "fmt"

type consultaIdentificadorDeTipo int

const (
	NENHUMA_INSCRICAO consultaIdentificadorDeTipo = 0
	UMA_INSCRICAO     consultaIdentificadorDeTipo = 1
	VARIAS_INSCRICOES consultaIdentificadorDeTipo = 2
)

type grabConsultaProcessoSidaResult struct {
	Json                    string
	quantidadeIdentificador consultaIdentificadorDeTipo
}

func (g *grabConsultaProcessoSidaResult) String() string {
	return fmt.Sprintf(`
	grabConsultaProcessoSidaResult{
		Json: %s
		QUANTIDADE ID: %v
	}
	`, g.Json, g.quantidadeIdentificador)
}
