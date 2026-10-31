package sdk

import "github.com/icosmos-space/llmatrix/engine/model"

// SupportsFunctionCalling mirrors litellm.supports_function_calling.
func SupportsFunctionCalling(modelName string, providerHint ...string) bool {
	return model.SupportsFunctionCalling(modelName, providerHint...)
}

// SupportsParallelFunctionCalling mirrors litellm.supports_parallel_function_calling.
func SupportsParallelFunctionCalling(modelName string, providerHint ...string) bool {
	return model.SupportsParallelFunctionCalling(modelName, providerHint...)
}

// SupportsVision mirrors litellm.supports_vision.
func SupportsVision(modelName string, providerHint ...string) bool {
	return model.SupportsVision(modelName, providerHint...)
}

// GetSupportedOpenAIParams mirrors litellm.get_supported_openai_params.
func GetSupportedOpenAIParams(modelName string, providerHint ...string) []string {
	return model.GetSupportedOpenAIParams(modelName, providerHint...)
}

// GetModelInfo mirrors litellm.get_model_info / model_info lookups.
func GetModelInfo(modelName string, providerHint ...string) model.Info {
	hint := ""
	if len(providerHint) > 0 {
		hint = providerHint[0]
	}
	return model.GetModelInfo(modelName, hint)
}
