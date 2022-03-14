package response

type ShortenUrl struct {
	Code             string  `json:"code"`
	Url              string  `json:"url"`
	RedirectionCount float64 `json:"redirectionCount"`
}

func NewShortenUrl(code, url string, count float64) ShortenUrl {
	return ShortenUrl{
		Code:             code,
		Url:              url,
		RedirectionCount: count,
	}
}
