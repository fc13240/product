package ebay

import (
	"encoding/xml"
)

type ReviseItem struct {
	XMLName xml.Name `xml:"ReviseItem"`
	Ebay
	ItemID int `xml:"Item>ItemID"`
}

type ReviseItemDescription struct {
	ReviseItem
	Description string `xml:"Item>Description"`
}

type ReviseItemSKU struct {
	ReviseItem
	SKU string `xml:"Item>SKU"`
}

func UpdateSKU(itemId int, SKU string) *Read {
	var sku ReviseItemSKU
	sku.ItemID = itemId
	sku.SKU = SKU
	return Up(&sku, "ReviseItem")
}

func UpdateDescript(itemId int, descript string) *Read {
	var desc ReviseItemDescription
	desc.ItemID = itemId
	desc.Description = descript
	return Up(&desc, "ReviseItem")
}

type ReviseItemQuantity struct {
	ReviseItem
	Quantity int `xml:"Item>Quantity"`
}

func UpdateQuantity(itemId, quantity int) *Read {
	var ebay ReviseItemQuantity
	ebay.ItemID = itemId
	ebay.Quantity = quantity
	return Up(&ebay, "ReviseItem")
}
