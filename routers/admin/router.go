package admin

import (
	"helper/configs"
	"helper/net/igin"
	ProductApi "product/api"
	alibaba "product/alibaba/api"
	lazada "product/lazada/api"
)

var (
	web_dir = configs.Get("web_dir")
)

func init() {

	apiRoute:=igin.R.Group("/api",Auth())

	itemGroup:=apiRoute.Group("/item")

	igin.R.GET("/login/callback/qq",QQcallBack)

	itemGroup.GET("packagefillattrs",ProductApi.PackageAttrsFill)

	apiRoute.POST("/items", ProductApi.Search)

	apiRoute.POST("/item.save", ProductApi.Save)
	apiRoute.POST("/item.editbase", ProductApi.EditBase)
	apiRoute.POST("/item.upimg", ProductApi.UpImg)
	apiRoute.GET("/item.ress", ProductApi.Res)              //资源列表
	apiRoute.POST("/item.setresflag", ProductApi.SetResFlag) //设置标签
	apiRoute.POST("/item.del",ProductApi.DelItem)
	apiRoute.GET("/item.labels",ProductApi.Labels)
	apiRoute.GET("/item.label.logs",ProductApi.LabelLogs)
	apiRoute.POST("/item.addlabellog",ProductApi.AddLabelLog)
	apiRoute.GET("/item/info/:sku", ProductApi.Detail)
	apiRoute.GET("/item.channels", ProductApi.GetChannelList)
	apiRoute.GET("/item.attrs", ProductApi.GetAttrs)
	apiRoute.POST("/item.addtotag", ProductApi.AddToTag)
	apiRoute.POST("/item.remtotag", ProductApi.RemToTag)
	apiRoute.POST("/item.setunder", ProductApi.SetUnder)
	apiRoute.POST("/item.setupper", ProductApi.SetUpper)
	apiRoute.POST("/item.settemp", ProductApi.SetTemp)
	apiRoute.POST("/item.bindebay", ProductApi.BindEbay)
	apiRoute.GET("/address.country", ProductApi.Country)
	apiRoute.GET("/item.template.tags", ProductApi.TemplateOpts)
	//apiRoute.GET("/item.resize", ProductApi.ResizeAllByProduct)
	//apiRoute.GET("/item.image.del", ProductApi.RemoveImage)
	apiRoute.GET("/item.margecontent", ProductApi.MargeContent)
	apiRoute.POST("/item.savetoebay",ProductApi.SaveToEbay)

	apiRoute.GET("/item.title.keyword.search",ProductApi.TitleKeywordSearch)
	apiRoute.POST("/item.title.keyword.import",ProductApi.TitleKeywordImport)
	//apiRoute.GET("/downtowarehouse", DownToWarehouse)
	//apiRoute.GET("/warehouse.productinfo", GetWarehouseProductInfo)
	//alibaiba

	apiRoute.POST("/alibaba.down",alibaba.Down)
	apiRoute.POST("/alibaba.downall",alibaba.DownAll)
	apiRoute.POST("/alibaba.tiankeng",alibaba.Tiankeng)

	apiRoute.POST("/alibaba.save",alibaba.Save)
	apiRoute.GET("/alibaba.get",alibaba.Get)
	apiRoute.GET("/alibaba.getsource",alibaba.BySource)
	apiRoute.POST("/alibaba.sources",alibaba.Sources)
	apiRoute.DELETE("/alibaba.delsource",alibaba.DelSource)
	apiRoute.POST("/alibaba.downsource",alibaba.DownSource)
	apiRoute.POST("/alibaba.set",alibaba.Set)
	apiRoute.POST("/alibaba.sellers",alibaba.Sellers)
	apiRoute.POST("/alibaba.add",alibaba.Add)
	apiRoute.GET("/alibaba.getattrs/:item_id",alibaba.GetAttrs)


	var ezbuy ProductApi.EzBuy
	apiRoute.POST("/ezbuy.orders",ezbuy.Orders)
	apiRoute.POST("/ezbuy.items",ezbuy.Listing)
	apiRoute.POST("/ezbuy.savesetting",ezbuy.SaveSetting)
	apiRoute.GET("/ezbuy.getsetting",ezbuy.GetSetting)
	apiRoute.POST("/ezbuy.saveitems",ezbuy.SaveItems)
	apiRoute.GET("/ezbuy.orderdetal",ezbuy.OrderDetal)
	apiRoute.POST("/ezbuy.setitem",ezbuy.SetItemField)
	apiRoute.POST("/ezbuy.saveorders",ezbuy.SaveOrders)
	apiRoute.POST("/ezbuy.checkneworders",ezbuy.CheckNewOrders)
	apiRoute.GET("/ezbuy.refreshorder/:ordernum",ezbuy.RefreshOrder)
	apiRoute.GET("/ezbuy.orderlabels",ezbuy.OrderLabels)
	apiRoute.POST("/ezbuy.addLabellog",ezbuy.AddLabelLog)
	apiRoute.GET("/ezbuy.cleanitems",ezbuy.CleanItems)
	apiRoute.GET("/ezbuy.get/:sku",ezbuy.Get)
	apiRoute.GET("/ezbuy.paste/:sku",ezbuy.Paste)
	apiRoute.GET("/ezbuy.colors",ezbuy.Colors)
	apiRoute.GET("/ezbuy.sizes",ezbuy.Sizes)
	apiRoute.POST("/ezbuy.export",ezbuy.Export)
	apiRoute.POST("/ezbuy.save",ezbuy.Save)
	apiRoute.GET("/ezbuy.categorys",ezbuy.Categorys)
	apiRoute.GET("/ezbuy.mycategorys",ezbuy.MyCategorys)
	apiRoute.POST("/ezbuy.addmycategory",ezbuy.AddMyCategory)
	apiRoute.POST("/ezbuy.saveUserProductsFromSource",ezbuy.SaveUserProductsFromSource)
	apiRoute.POST("/ezbuy.userProductsFromSource",ezbuy.UserProductsFromSource)
	apiRoute.GET("/ezbuy.getAttrs/:cid",ezbuy.GetAttrs)

	//lazada
	_lazada:=apiRoute.Group("/lazada")

	_lazada.POST("/SaveUploadImageInfo",lazada.SaveUploadImageInfo)
	_lazada.GET("/UploadImages/:item_id",lazada.UploadImages)


	apiRoute.POST("/account.create", CreateAccount)
	apiRoute.POST("/account.login", Login)

	apiRoute.POST("/account.sendregistercode", SendRegisterCodeMail)
	apiRoute.GET("/account.logout", Logout)
	apiRoute.GET("/account.info", GetAccountInfo)

	apiRoute.POST("/account.saveinfo", SaveAccountInfo)
	apiRoute.POST("/account.upheadimg", SaveAccountHeadimg)
	apiRoute.GET("/account.messagelist", MessageList)
	apiRoute.POST("/account.repass", RePassword)
	//apiRoute.GET("/account.addaddr", SaveShippingAddress)
	//apiRoute.GET("/account.addrs", GetShippingAddress)
	//apiRoute.GET("/account.addrdel", DelShippingAddress)
	//apiRoute.GET("/account.active/{code}/{uid}", ActiveAccount)
	//apiRoute.GET("/account.active/{code}/{uid}", ActiveAccount)
	//apiRoute.GET("/message.send", SendMessage)
	apiRoute.GET("/message.del", DelMessage)
	apiRoute.GET("/myserves",MyServeList)

	apiRoute.POST("/label.addfollow",AddFollowLabel)

	apiRoute.POST("/follow.labelcancel",FollowCancel)
	apiRoute.GET("/follow.labels",FollowLabels)

	//manage
	manageGroup:=apiRoute.Group("/manage")
	manageGroup.GET("/console",Console)
	manageGroup.POST("/members",Members)
	manageGroup.POST("/servers",Servers)
	//doc
	docGroup:=apiRoute.Group("/doc")
	docGroup.POST("/label.add", addLabel)
	docGroup.POST("/label.search",SearchLabel)
	docGroup.POST("/create", Create)
	docGroup.POST("/listing", Listing)
	docGroup.GET("/get/:id", Get)
	docGroup.POST("/mylist",MyList)

	//factory
	factoryGroup:=apiRoute.Group("/factory")
	factoryGroup.POST("/app.save",factoryAction.CreateApp)
	factoryGroup.POST("/apps",factoryAction.Apps)
	factoryGroup.GET("/app.get/:id",factoryAction.Get)
	factoryGroup.POST("/app.req/:id",factoryAction.Req)
}