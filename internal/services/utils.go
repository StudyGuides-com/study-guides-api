package services

import (
	sharedv1 "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"github.com/studyguides-com/study-guides-api/internal/types"
)

func ToProtoContextType(internal types.ContextType) sharedv1.ContextType {
	switch internal {
	case types.ContextTypeColleges:
		return sharedv1.ContextType_Colleges
	case types.ContextTypeCertifications:
		return sharedv1.ContextType_Certifications
	case types.ContextTypeEntranceExams:
		return sharedv1.ContextType_EntranceExams
	case types.ContextTypeAPExams:
		return sharedv1.ContextType_APExams
	case types.ContextTypeDoD:
		return sharedv1.ContextType_DoD
	case types.ContextTypeUserGeneratedContent:
		return sharedv1.ContextType_UserGeneratedContent
	case types.ContextTypeAll:
		return sharedv1.ContextType_All
	default:
		return sharedv1.ContextType_All
	}
}

func FromProtoContextType(proto sharedv1.ContextType) types.ContextType {
	switch proto {
	case sharedv1.ContextType_Colleges:
		return types.ContextTypeColleges
	case sharedv1.ContextType_Certifications:
		return types.ContextTypeCertifications
	case sharedv1.ContextType_EntranceExams:
		return types.ContextTypeEntranceExams
	case sharedv1.ContextType_APExams:
		return types.ContextTypeAPExams
	case sharedv1.ContextType_DoD:
		return types.ContextTypeDoD
	case sharedv1.ContextType_UserGeneratedContent:
		return types.ContextTypeUserGeneratedContent
	case sharedv1.ContextType_All:
		return types.ContextTypeAll
	default:
		return types.ContextTypeAll
	}
}
