package services

import (
	sharedv1 "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"github.com/studyguides-com/study-guides-api/internal/types"
)


func ToProtoContextType(internal types.ContextType) sharedv1.ContextType	 {
	switch internal {
	case types.ContextTypeColleges:
		return sharedv1.ContextType_CONTEXT_TYPE_COLLEGES
	case types.ContextTypeCertifications:
		return sharedv1.ContextType_CONTEXT_TYPE_CERTIFICATIONS
	case types.ContextTypeEntranceExams:
		return sharedv1.ContextType_CONTEXT_TYPE_ENTRANCE_EXAMS
	case types.ContextTypeAPExams:
		return sharedv1.ContextType_CONTEXT_TYPE_AP_EXAMS
	case types.ContextTypeDoD:
		return sharedv1.ContextType_CONTEXT_TYPE_DOD
	case types.ContextTypeUserGeneratedContent:
		return sharedv1.ContextType_CONTEXT_TYPE_USER_GENERATED_CONTENT
	case types.ContextTypeAll:
		return sharedv1.ContextType_CONTEXT_TYPE_ALL
	default:
		return sharedv1.ContextType_CONTEXT_TYPE_ALL
	}
}

func FromProtoContextType(proto sharedv1.ContextType) types.ContextType {
	switch proto {
	case sharedv1.ContextType_CONTEXT_TYPE_COLLEGES:
		return types.ContextTypeColleges
	case sharedv1.ContextType_CONTEXT_TYPE_CERTIFICATIONS:
		return types.ContextTypeCertifications
	case sharedv1.ContextType_CONTEXT_TYPE_ENTRANCE_EXAMS:
		return types.ContextTypeEntranceExams
	case sharedv1.ContextType_CONTEXT_TYPE_AP_EXAMS:
		return types.ContextTypeAPExams
	case sharedv1.ContextType_CONTEXT_TYPE_DOD:
		return types.ContextTypeDoD
	case sharedv1.ContextType_CONTEXT_TYPE_USER_GENERATED_CONTENT:
		return types.ContextTypeUserGeneratedContent
	case sharedv1.ContextType_CONTEXT_TYPE_ALL:
		return types.ContextTypeAll
	default:
		return types.ContextTypeAll
	}
}

