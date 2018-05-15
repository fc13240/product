package ebay

import (
	"fmt"
)

type Account struct {
	DevId  string
	AppId  string
	CertId string
	Token  string
	AppUrl string
}

var Mod string = "dev"

func ProductionMod() {
	Mod = "production"
}

func (account *Account) Init() *Account {
	if Mod == "dev" {
		fmt.Println(Mod)
		account.d()
	} else if Mod == "production" {
		fmt.Println(Mod)
		account.p()
	}

	return account
}

func (account *Account) d() {
	account.DevId = "43f3e57c-fe08-412d-8aa5-162ef3de2008"
	account.AppId = "meilan407-fe83-4065-8661-e3a1cad21c0"
	account.CertId = "2ce0d386-d592-41cc-bd03-d3a4fd240117"
	account.Token = "AgAAAA**AQAAAA**aAAAAA**LjKaVQ**nY+sHZ2PrBmdj6wVnY+sEZ2PrA2dj6wFk4GhDZeEoAWdj6x9nY+seQ**JS8DAA**AAMAAA**ybVVv+MUBVzoX91G0G/6tF8RJE9D+nYoWYC8NNiu2SHP0/crqTbMYwGXN1WnzPZQq+XOcjK+yyycUDpoKdkxC4nNrvhI7W5nydNsY5mVj+LUGjYQNvu8iogeAXlNElMRr9SljLudeEsw+qUOV1V884V2gNpOD4j/dtNfv16HtgJfuurEMms4BsE5derYPNcuKNM4EKAMrtUOF+7pbfGG+zQp8GYbHEbADo4yfZoqt+KP6d7tV3qeVO8qoW/KippIYZmnoOhYRhd7tiMq+d5Cgvn4XjdYnz4LCUXYnCjl1LXG+3+a1IJ4w+gSV6HQGuwLcPkhoc1OeMUWNam880RjzNAUW35fzwuBWtRp4kPJ5yPpx3ZM5ijHWq0tmpbb0iwyh/7jNEt7l8g9Nw9D9MppEAV73jHpvrCdbu9pKMU0pzuKeI8xynwIqiZfDB/scQhMTJXADN9ZFTXNZRMavJ2//+uhbOEsX4HY4SWDou1ulPFpVBclP0OC/UkliSkrWaaOJ/MDSWdHHP+1+CYICE38/N3BoXpWTq3FVBDoaXpnVLHCoUwwfpXRYstTtwuq4UoqKhQ+dm/OKqdLwtwXBvY6bRAEwy5qrNGAA3akRDlRLbtk0ZSD4j+MnMU9F5V9+HgoaIBIGeGD5O8N4cnVlnvAJ0PapmPd5lIxUjhMphfy5HS8qqlqNlPm4b+DRvsvzrLrZmt8yQblso/FazPPH5KD0sRP0nD2qZX9fB2wFSeYaCFzYBptV4h3ghMAfImQlqMf"
	account.AppUrl = "https://api.sandbox.ebay.com/ws/api.dll"
}

func (account *Account) p() {
	account.DevId = "43f3e57c-fe08-412d-8aa5-162ef3de2008"
	account.AppId = "meilan2a2-2ac1-4b11-8641-7475ef9ba62"
	account.CertId = "87917947-7777-4b69-9e16-743006731162"
	account.Token = "AgAAAA**AQAAAA**aAAAAA**YIs+WA**nY+sHZ2PrBmdj6wVnY+sEZ2PrA2dj6AGloCjCJaCqAWdj6x9nY+seQ**TaMCAA**AAMAAA**0a0+jT2PjxHLiCJS4OhHBX5VdvWOlE20EhJwnYwxARThZnCEivZzgXk/xvBFmIl/JdtDtzQuM8yC7Y6mTI0+rnIT9L/rsc9Lrwhz0ubbl9KErCpSv7bHVcThf0p1HxgPE0LUWxxAqbm5l93zfoLlH1VDUaDiG7WN+G2yJv6iAtBkov6Iwa1KZjT+8MKhiTJyMvKxr6l/DYTQw1bFT70IDFN9dMpTsP6yitg1cMcxixOJ7qelnlKPSc1naE5IzqYW1NhH/h7ounCtaDJoYeELIKi0n9YoWVv9gSWQc0+4fYlx1btSdS6FZBjCcXKnEyJqArwVQ6lybObriXypBW8Dn8beMawA5et2yvTLiU7V5SiKL6eqfhw4sDO4mfCyt4ndjJ+ZPbpCK8gV2MpRjC4DMpO0hkb5sEW7aIqBlWI4ca0goOsek6arsfbsCEflQFEfZ94SmqSEHXWGBIDMVAD6k/c+Z6sAEabJpRDyvbcamm53ZnnQI62NpChUc6eo8hDNwFv6hk7WO0D4TpBYcRjLigeV4cisqE0ES9l7abSQwCgzdQD1CSEYyK+OFBbPuHayFeJxSBSrWL6sRIlsxN2n+mSbaHaAqC16l+EoALPLbSwi/jw/KZ06qY0c/25O4q6+ku1WNW/YYZFV/FkgXu6O09OfMM075rUnLXhM6bQGqdfOM7FzKh4vk11sKvNw+tO+y/t2uu8DGzZlfgFhhuvaobXq1rX9JXHViI+qCNl7YTBHfaXuotqefnEFxonrzOu5"
	account.AppUrl = "https://api.ebay.com/ws/api.dll"
}
