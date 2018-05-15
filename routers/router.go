package routers

import (
	"helper/configs"
	"helper/net/igin"

	//_ "product/shopee/router"
	"routers/member"
	"routers/factory"
	"routers/manage"
	doc "doc/action"

	"routers/passport"
	"github.com/gin-gonic/contrib/static"

	product "product/action"
	productCategory "product/category/action"
	productFees "product/fee/action"
	productCountry "product/country/action"
	productGallery "product/gallery/action"
	productTitle "product/title/action"
	productAttr "product/attr/action"

	lazada  "product/lazada/action"
	shopee "product/shopee/action"
	ezbuy "product/ezbuy/action"
	alibaba "product/alibaba/action"

)
var (
	web_dir = configs.Get("web_dir")
	download_dir=configs.Get("download_dir")
)

func init() {
	r:=igin.R
	start()
	r.Use(static.Serve("/download",static.LocalFile(download_dir,true)))
	
	//passport
	r.GET("/login/callback/qq",passport.QQcallBack)
	api:=igin.R.Group("/api",passport.Auth())

	//product 
	api.POST("/items", product.Search)
	
	api.GET("/item/:sku/skuusestatus",product.GetSkuUseStatus) //sku use status
	api.GET("/item/:sku/categoryinfo",product.GetItemCategoryInfo) //sku use status
	
	api.POST("/item/:sku", product.Save)
	api.PUT("/item/:sku/eidtfield",product.SaveField)
	api.POST("/item/edit", product.EditBase)
	api.POST("/item/del",product.DelItem)
	api.GET( "/item/:sku", product.Detail) 

	//product gallery 
	api.POST("/item/images", productGallery.Images)  //图片列表
	api.POST("/item/image",productGallery.AddImage)
	api.POST("/item/image/del",productGallery.DelImage)
	api.POST("/item/image/upload", productGallery.UploadImage)
	api.POST("/item/image/flag",productGallery.SetImageFlag) //设置标签

	//product labal
	api.GET("/item/:sku/labels",product.Labels)
	api.POST("/item/label/logs",product.LabelLogs)
	api.POST("/item/label/log",product.AddLabelLog)
	
	//shop channels
	api.GET("/channels", product.GetChannelList)
	
	//platform 
	api.GET("/platform/:platform/category/:cid",product.GetPlatformCategorys)
	api.GET("/platform/:platform/category/:cid/selected",product.GetPlatformCategorySelected)
	

	//product tag
	api.POST("/item/tag", product.AddToTag)
	api.DELETE("/item/tag", product.RemToTag)        
	api.POST("/items/template/:tmpid/margecontent", product.MargeContent)

	//address
	api.GET("/address/country", product.Country)

	api.POST("/item/addSkuUploadLog",product.AddSkuUploadLog)          
	//api.GET("/product.resize", ProductApi.ResizeAllByProduct)
	
	//product category
	api.GET("/items/categorys",productCategory.Listing)
	api.POST("/items/categorys",productCategory.Add)
	api.GET("/items/categorys/childs/:pid",productCategory.Categorys)
	api.PUT("/items/categorys/attrs/:platform/:cid",productCategory.SetAttr)
	api.GET("/items/categorys/attrs/:platform/:cid",productCategory.GetAttrIds)
	api.POST("/items/categorys/attrs",product.GetCategoryAttrs)

	//product attr
	api.GET("/items/attrs",productAttr.GetAttrs)
	api.POST("/items/attrs",productAttr.AddAttr)
	api.GET("/items/attrs/packagefill",productAttr.PackageFill)

	api.POST("/item/attr/options",productAttr.AddAttrOpt)
	api.GET("/item/attr/:attid/options",productAttr.GetAttOptions)

	//product attr value
	api.POST("/item/attr/values",productAttr.SaveAttrVal)
	api.GET("/item/:sku/attr/:attrid/values",productAttr.GetAttrVal)
	api.GET("/item/:sku/attr/:attrid/values/selected",productAttr.GetOneAttrSelectedOption)
	api.GET("/item/:sku/attr/values",productAttr.GetSkuSelectedOptions)


	//product size
	api.POST("/item/size/template",product.SaveSizeTemplate)
	api.GET("/item/size/template/:item_id",product.GetSizeTemplate)

	//product title 
	api.POST("/items/titles",productTitle.TitleKeywordSearch)
	api.POST("/items/title/label",productTitle.AddTitleLabel)
	api.POST("/items/title/labelcate",productTitle.AddTitleLabelCate)	
	api.GET( "/items/title/labels/:cid",productTitle.GetTitleLabels)
	api.POST("/items/title/labels/:cid/import",productTitle.TitleKeywordImport)
	api.POST("/items/title/labelcates",productTitle.GetTitleLabelCateListing)
	//api.GET("/downtowarehouse", DownToWarehouse)
	//api.GET("/warehouse.productinfo", GetWarehouseProductInfo)
	//alibaiba

	//product fee
	api.GET("/item/:sku/fees",productFees.Listing)
	api.PUT("/item/fees",productFees.Edit)
	api.POST("/item/fees",productFees.Add)
	
	//country
	api.GET("/currencys",productCountry.GetCurrencys)

	//alibaba api
	api.POST("/alibaba/down",alibaba.Down)
	api.POST("/alibaba/downall",alibaba.DownAll)
	api.POST("/alibaba/tiankeng",alibaba.Tiankeng)
	api.POST("/alibaba/save",alibaba.Save)
	api.GET("/alibaba/get",alibaba.Get)
	api.GET("/alibaba/getsource",alibaba.BySource)
	api.POST("/alibaba/sources",alibaba.Sources)
	api.DELETE("/alibaba/delsource",alibaba.DelSource)
	api.POST("/alibaba/downsource",alibaba.DownSource)
	api.POST("/alibaba/set",alibaba.Set)
	api.POST("/alibaba/sellers",alibaba.Sellers)
	api.POST("/alibaba/add",alibaba.Add)
	api.GET("/alibaba/getattrs/:item_id",alibaba.GetAttrs)
	api.POST("/alibaba/search",alibaba.Search)
	api.DELETE("/alibaba/:id",alibaba.Delete)
	api.POST("/alibaba/checkUrlExist",alibaba.CheckUrlExist)
	api.GET("/alibaba/colorandsize/:sku",alibaba.GetColorsAndSizes)
	
	//ezbuy api
	api.POST("/ezbuy/orders",ezbuy.Orders)
	api.POST("/ezbuy/items",ezbuy.Listing)
	api.POST("/ezbuy/savesetting",ezbuy.SaveSetting)
	api.GET("/ezbuy/getsetting",ezbuy.GetStore)
	api.GET("/ezbuy/store/:storeid",ezbuy.GetStore)
	api.POST("/ezbuy/saveitems",ezbuy.SaveItems)
	api.GET("/ezbuy/orderdetal",ezbuy.OrderDetal)
	api.POST("/ezbuy/setitem",ezbuy.SetItemField)
	api.POST("/ezbuy/saveorders",ezbuy.SaveOrders)
	api.POST("/ezbuy/checkneworders",ezbuy.CheckNewOrders)
	api.GET("/ezbuy/refreshorder/:ordernum",ezbuy.RefreshOrder)
	api.GET("/ezbuy/orderlabels",ezbuy.OrderLabels)
	api.POST("/ezbuy/addLabellog",ezbuy.AddLabelLog)
	api.GET("/ezbuy/cleanitems",ezbuy.CleanItems)
	api.GET("/ezbuy/get/:sku",ezbuy.Get)
	api.GET("/ezbuy/paste/:sku",ezbuy.Paste)
	api.GET("/ezbuy/colors",ezbuy.Colors)
	api.GET("/ezbuy/sizes",ezbuy.Sizes)
	api.POST("/ezbuy/export",ezbuy.Export)
	api.POST("/ezbuy/save",ezbuy.Save)
	api.GET("/ezbuy/categorys/:parent_id",ezbuy.Categorys)
	api.GET("/ezbuy/mycategorys",ezbuy.MyCategorys)
	api.POST("/ezbuy/addmycategory",ezbuy.AddMyCategory)
	api.POST("/ezbuy/saveUserProductsFromSource",ezbuy.SaveUserProductsFromSource)
	api.POST("/ezbuy/userProductsFromSource",ezbuy.UserProductsFromSource)
	api.GET("/ezbuy/getAttrs/:cid",ezbuy.GetAttrs)
	api.GET("/ezbuy/mystores",ezbuy.MyStores)
	api.POST("/ezbuy/fissiont/:sku/bycolor",ezbuy.FissionItemByColor)
	api.POST("/ezbuy/fissiont/:sku",ezbuy.NewFission)
	api.GET("/ezbuy/fissiont/:sku",ezbuy.GetItemChilds)
	api.DELETE("/ezbuy/fissiont/:sku",ezbuy.DelItemChild)
	api.PUT("/ezbuy/setinfor",ezbuy.SetInfor)

	
	//lazada api
	api.POST("/lazada/images/upload",lazada.SaveUploadImageInfo)
	api.GET("/lazada/images/uploadall/:sku",lazada.UploadImages)
	api.POST("/lazada/item",lazada.SaveProduct)
	api.GET("/lazada/item/:sku",lazada.GetProduct)
	api.GET("/lazada/create/:sku",lazada.Create)

	//shopee api
	api.POST("/shopee/items",shopee.SaveProduct)
	api.GET("/shopee/items/:sku",Get)
	api.DELETE("/shopee/remote/item/:sku",shopee.DelRemoteItem)
	api.PUT("/shopee/remote/item/:sku",shopee.AddRemoteItem)

	//platform  ebay
	api.POST("/ebay/bind", product.BindEbay)
	api.POST("/ebay/upload",product.SaveToEbay)

	//member
	api.POST("/account/create", member.CreateAccount)
	api.POST("/account/login",member.Login)
	api.POST("/account/repassword", member.RePassword)
	api.POST("/account/sendregistercode", member.SendRegisterCodeMail)
	api.GET("/account/logout",member.Logout)
	api.GET("/account/info", member.GetAccountInfo)

	api.POST("/account/saveinfo", member.SaveAccountInfo)
	api.POST("/account/upheadimg", member.SaveAccountHeadimg)
	api.GET("/account/messages",member.MessageList)
	//api.GET("/account.addaddr", SaveShippingAddress)
	//api.GET("/account.addrs", GetShippingAddress)
	//api.GET("/account.addrdel", DelShippingAddress)
	//api.GET("/account.active/{code}/{uid}", ActiveAccount)
	//api.GET("/account.active/{code}/{uid}", ActiveAccount)
	//api.GET("/message.send", SendMessage)
	api.GET("/message/del", member.DelMessage)
	api.GET("/myserves",member.MyServeList)

	api.POST("/label/addfollow",member.AddFollowLabel)

	api.POST("/follow/labelcancel",member.FollowCancel)
	api.GET("/follow/labels",member.FollowLabels)

	//manage
	manageGroup:=api.Group("/manage")
	manageGroup.GET("/console",manage.Console)
	manageGroup.POST("/members",manage.Members)
	manageGroup.POST("/servers",manage.Servers)
	
	//doc
	docGroup:=api.Group("/doc")
	
	docGroup.POST("/labels",doc.SearchLabel)
	docGroup.POST("/label/add",doc.AddLabel)
	docGroup.POST("/create",doc.Create)
	docGroup.POST("/listing",doc.Listing)
	docGroup.GET("/get/:id",doc.Get)
	docGroup.POST("/mylist",doc.MyList)

	//factory
	factoryGroup:=api.Group("/factory")
	factoryGroup.POST("/app/save",factory.CreateApp)
	factoryGroup.POST("/apps",factory.Apps)
	factoryGroup.GET("/app/:id",factory.Get)
	factoryGroup.POST("/app/req/:id",factory.Req)
}