package fees
import(
	"helper/dbs"
	"helper/configs"
)

type Fees struct{
	Sku string
}

type FeeRate struct{
	Code string `json:"code"`
	Rate float32 `json:"rate"`
	Remarks string `josn:"remarks"`
}

func (fee *Fees)addFree(price float32,labal,currency,remarks string){
	db:=dbs.Def()
	db.Insert("product_fees",configs.M{
		"sku":fee.Sku,
		"price":price,
		"currency":currency,
		"labal":labal,
		"remarks":remarks,
	})
}

//费用清单
func Listing(sku,platform string)(fees []*FeeOpt,rate_count float32){
	fees=[]*FeeOpt{}
	fees=append(fees,FixedFee(sku)...)
	
	if platform =="ezbuy" {
		ezfee:=FeeRate{"ezbuyServerFee",10.0,"ezbuy技术服务费"}
		rate_count=ezfee.Rate
	}
	return 
}

func NewFeeOpt(code string ,price float32)*FeeOpt{
	return &FeeOpt{fess[code],code,price}
}

var fess =map[string]string{
	"purchasePrice":"进货价格",
	"shippingCharge":"快递费",
}

func FixedFee(sku string) []*FeeOpt{
	db:=dbs.Def()
	var buying_price float32
	db.One("SELECT buying_price FROM product WHERE sku=?",sku).Scan(&buying_price)
	
	fixedfee:=[]*FeeOpt{}

	fixedfee=append(fixedfee,NewFeeOpt("purchasePrice",buying_price))

	//快递费
	fixedfee=append(fixedfee,NewFeeOpt("shippingCharge",8))
	return fixedfee
}

type FeeOpt struct{
	Name string `json:"name"`
	Labal string `json:"labal"`
	Fee float32 `json:"fee"`
}

func (fee *Fees)AddPlatformServiceFeeRate(){
	
}
