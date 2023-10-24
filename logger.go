package main

import "log"

type CompilerStage int

const (
	STAGE_LEXER CompilerStage = iota
	STAGE_PARSER
)

func Fatal(stage CompilerStage, message string) {
	var stageStr string
	switch stage {
	case STAGE_LEXER:
		stageStr = "lexer"
	case STAGE_PARSER:
		stageStr = "parser"
	}

	log.Fatalf("[%s] FATAL: %s", stageStr, message)
}
