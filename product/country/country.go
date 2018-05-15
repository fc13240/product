package country

//货币
const(
	CNY="CNY"
	RM ="RM"
	SGD="SGD"
)
type Currency struct{
	Country string `json:"country"`
	Code string `json:"code"`
	ExchangeRate float32 `json:"exchange_rate"`
}

func GetCurrencys()[]Currency{
	return []Currency{
		{"人民币",CNY,1},
		{"马来西亚",RM,1.6129},
		{"新加坡",SGD,4.7693},
		
	}
}