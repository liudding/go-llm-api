package baidu

const (
	BaiduModelERNIEBot4                 = "completions_pro"
	BaiduModelERNIEBot                  = "completions"
	BaiduModelERNIEBotTurbo             = "eb-instant"
	BaiduModelBLOOMZ7B                  = "bloomz_7b1"
	BaiduModelQianfanBLOOMZ7BCompressed = "qianfan_bloomz_7b_compressed"
)

// Usage Represents the total token usage per request to OpenAI.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
