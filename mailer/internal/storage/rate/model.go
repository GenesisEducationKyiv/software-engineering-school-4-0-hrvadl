package rate

type Exchange struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Rate float32 `json:"rate"`
}
