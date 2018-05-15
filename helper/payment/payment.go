package payment

import (
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/paypalsdk"
	"helper/configs"
	"helper/dbs"
	"helper/redisCli"
	utiltime "helper/time"
	"log"
	"net/url"
	"time"
	"order"
)

var(
	paypal_option = configs.GetSection("paypal")
	clickid = paypal_option["clickid"]
	secret = paypal_option["secret"]
	redirectURI = paypal_option["redirectURI"]
	cancelURI = paypal_option["cancelURI"]
)

type OrderPayment struct {
	Id            int
	Orderid       int
	PaymentStatus int
	TokentCode    string
	PaymentCode   string
	PayerCode     string
	Addtime       *time.Time
	Paymenttime   *time.Time
}

func (payment *OrderPayment) IsPaymentWait() bool {
	if payment.PaymentStatus == 0 {
		return true
	}
	return false
}

func (payment *OrderPayment) IsPaymentComplete() bool {
	if payment.PaymentStatus == 2 {
		return true
	}
	return false
}

func parsePayUrl(href string) (url.Values, error) {
	if route, err := url.Parse(href); err == nil {
		if m, err := url.ParseQuery(route.RawQuery); err == nil {
			return m, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func ApprovedPayment(payment_code, payer_code, token_code string) error {
	orderPayment := &OrderPayment{}

	dbs.One("SELECT id,order_id,payment_status FROM order_pay WHERE tokent_code=?", token_code).
		Scan(&orderPayment.Id, &orderPayment.Orderid, &orderPayment.PaymentStatus)

	if false == (orderPayment.Id > 0) {
		log.Println("not payment token:", token_code)
		return errors.New("No Payment Record")
	}

	if orderPayment.IsPaymentComplete() {
		return nil
	}

	if false == orderPayment.IsPaymentWait() {
		return errors.New("No Wait Payment Orders")
	}

	//更新支付表
	dbs.Update("order_pay", configs.M{
		"payer_code":   payer_code,
		"payment_code": payment_code,
		"uptime":       utiltime.String(),
	}, "id=?", orderPayment.Id)

	if c, err := paypalsdk.NewClient(clickid, secret, paypalsdk.APIBaseSandBox); err == nil {
		c.GetAccessToken()
		if res, e := c.ExecuteApprovedPayment(payment_code, payer_code); e == nil {
			//complete
			if res.State != "approved" {
				return errors.New("Payment Approved Failing")
			}

			dbs.Update("order",
				configs.M{
					"payment_status": 2,
					"order_status":   2,
				},
				"order_id=?", orderPayment.Orderid)

			dbs.Update("order_pay",
				configs.M{
					"payment_time":   utiltime.String(),
					"payment_status": 2,
				},
				"id=?", orderPayment.Id)
			return nil
		} else {
			return e
		}
	} else {
		return err
	}
}

type Item struct{
	Quantity int
	Name	string
	Price   string
	Currency string
	SKU string
	Description string
}

type Amount struct{
	Currency string
	Total string
}

type Pay struct{
	OrderId int
	Items []Item
	Total  string
	Currency string
	Remark string
}

func (p *Pay)CreatePaymentUrl() (string,error){
	var items []paypalsdk.Item
	for _, item := range p.Items {
		item := paypalsdk.Item{
			Quantity:    item.Quantity,
			Name:        item.Name,
			Price:       item.Price,
			Currency:    p.Currency,
			SKU:         item.SKU,
			Description: item.Description,
		}
		items = append(items, item)
	}

	transaction := paypalsdk.Transaction{
		Amount:      &paypalsdk.Amount{Currency: "USD", Total: p.Total},
		ItemList:    &paypalsdk.ItemList{Items: items},
		Description: p.Remark,
	}

	payment := paypalsdk.Payment{
		Intent:       "sale",
		Payer:        &paypalsdk.Payer{PaymentMethod: "paypal"},
		Transactions: []paypalsdk.Transaction{transaction},
		RedirectURLs: &paypalsdk.RedirectURLs{ReturnURL: redirectURI, CancelURL: cancelURI},
	}

	if c, err := paypalsdk.NewClient(clickid, secret, paypalsdk.APIBaseSandBox); err == nil {
		setAccessToken(c)
		if payment_result, err := c.CreatePayment(payment); err == nil {
			href := payment_result.Links[1].Href
			if param, err := parsePayUrl(href); err == nil {
				dbs.Insert("order_pay", configs.M{"order_id": p.OrderId, "tokent_code": param.Get("token"), "addtime": utiltime.String()})
				log.Println(href)
				return href, nil
			} else {
				return "", err
			}
		} else {
			data, _ := json.MarshalIndent(payment, "", "    ")
			log.Println(string(data))
			log.Println(payment)
			return "", err
		}
	} else {
		return "", err
	}
}


func New(items []Item,amount Amount) {
	return Pay{Items:items,Currency:currency}
}

func setAccessToken(c *paypalsdk.Client) string {
	r := redisCli.Conn()

	token_catch, _ := redis.String(r.Do("get", "paypal_account_token:"))

	if token_catch == "" {

		if token, err := c.GetAccessToken(); err == nil {
			r.Do("SET", "paypal_account_token:", token.Token)
			r.Do("EXPIRE", "paypal_account_token:", token.ExpiresIn)

			return token.Token
		} else {
			log.Println("paypal token failing")
			return ""
		}
	} else {
		log.Println("read is redis cache")
		token := token_catch
		c.Token = &paypalsdk.TokenResponse{Token: token, ExpiresIn: 3600}
		return token
	}
}

/*
	shipping := order.GetShippingAddress()
	shippingAddress := &paypalsdk.ShippingAddress{
		RecipientName: shipping.FirstName + shipping.LastName,
		Line1:         shipping.Address,
		Line2:         shipping.OtherAddress,
		City:          shipping.City,
		CountryCode:   shipping.CountryCode,
		PostalCode:    shipping.Zip,
		State:         shipping.State,
		Phone:         shipping.Phone,
	}
	ShippingAddress: shippingAddress
*/
