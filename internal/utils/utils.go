package utils

import (
	"github.com/lucsky/cuid"
	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

func GetCUID() string {
	return cuid.New()
}

// ParserContextMapper maps ParserType to ContextType
var ParserContextMapper = map[sharedpb.ParserType]sharedpb.ContextType{
	sharedpb.ParserType_PARSER_TYPE_COLLEGES:       sharedpb.ContextType_Colleges,
	sharedpb.ParserType_PARSER_TYPE_CERTIFICATIONS: sharedpb.ContextType_Certifications,
	sharedpb.ParserType_PARSER_TYPE_ENTRANCE_EXAMS: sharedpb.ContextType_EntranceExams,
	sharedpb.ParserType_PARSER_TYPE_AP_EXAMS:       sharedpb.ContextType_APExams,
	sharedpb.ParserType_PARSER_TYPE_DOD:            sharedpb.ContextType_DoD,
}

// GetContextTypeForParser returns the corresponding ContextType for a given ParserType
func GetContextTypeForParser(parserType sharedpb.ParserType) (sharedpb.ContextType, bool) {
	contextType, exists := ParserContextMapper[parserType]
	return contextType, exists
}

// GetParserTypeForContext returns the corresponding ParserType for a given ContextType
func GetParserTypeForContext(contextType sharedpb.ContextType) (sharedpb.ParserType, bool) {
	for parserType, ctxType := range ParserContextMapper {
		if ctxType == contextType {
			return parserType, true
		}
	}
	return sharedpb.ParserType_PARSER_TYPE_UNSPECIFIED, false
}
