package houses

// House representation
type House struct {
	ID       int    `json:"id"`
	Address  string `json:"address"`
	Owner    string `json:"homeowner"`
	PhotoURL string `json:"photoURL"`
}
