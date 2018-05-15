package main 
import(
	_ "ainit"
	"product"
)
func main(){
	AddMaterial()
}

func add(){
	product.NewAtt("品牌(Brand)","brand",product.INPUT_OPTIONS)
	product.NewAtt("颜色(Colors)","colors",product.INPUT_CHECKBOX)
	product.NewAtt("材料(Material)","material",product.INPUT_OPTIONS)
	product.NewAtt("裙子风格(Skirt Style)","skirt_style",product.INPUT_OPTIONS)
	product.NewAtt("裙长(Skirt Length)","skirt_length",product.INPUT_OPTIONS)
	product.NewAtt("场合(Occasion)","occasion",product.INPUT_OPTIONS)
}

func AddColors(){
	type Val struct{
		Name string
		CnName string
		}

	colors:=[]Val{{
		Name: "Black",
		CnName:"黑色",
	},
	{
		Name: "Beige",
		CnName: "米色",
	},
	{
		Name: "Blue",
		CnName: "蓝色",
	},
	{
		Name: "Brown",
		CnName: "棕色",
	},
	{
		Name: "Gold",
		CnName: "金色",
	},
	{
		Name: "Green",
		CnName: "绿色",
	},
	{
		Name: "Grey",
		CnName: "灰色",
	},
	{
		Name: "Multicolor",
		CnName: "彩色",
	},
	{
		Name: "Olive",
		CnName: "橄榄色",
	},
	{
		Name: "Orange",
		CnName: "橙色",
	},
	{
		Name: "Pink",
		CnName: "粉色",
	},
	{
		Name: "Purple",
		CnName: "紫色",
	},
	{
		Name: "Red",
		CnName: "红色",
	},
	{
		Name: "Silver",
		CnName: "银色",
	},
	{
		Name: "Turquoise",
		CnName: "绿松石",
	},
	{
		Name: "Violet",
		CnName: "紫色",
	},
	{
		Name: "White",
		CnName: "白色",
	},
	{
		Name: "Yellow",
		CnName: "黄色",
	},
	{
		Name: "Clear",
		CnName: "浅色的",
	},
	{
		Name: "Apricot",
		CnName: "杏",
	},
	{
		Name: "Aqua",
		CnName: "水族",
	},
	{
		Name: "Avocado",
		CnName: "鳄梨",
	},
	{
		Name: "Blueberry",
		CnName: "蓝莓",
	},
	{
		Name: "Blush Pink",
		CnName: "脸红粉红色",
	},
	{
		Name: "Bronze",
		CnName: "青铜",
	},
	{
		Name: "Charcoal",
		CnName: "木炭",
	},
	{
		Name: "Cherry",
		CnName: "樱桃",
	},
	{
		Name: "Chestnut",
		CnName: "板栗",
	},
	{
		Name: "Chili Red",
		CnName: "辣椒红",
	},
	{
		Name: "Chocolate",
		CnName: "巧克力",
	},
	{
		Name: "Cinnamon",
		CnName: "肉桂",
	},
	{
		Name: "Coffee",
		CnName: "咖啡",
	},
	{
		Name: "Cream",
		CnName: "奶油",
	},
	{
		Name: "Floral",
		CnName: "花的",
	},
	{
		Name: "Galaxy",
		CnName: "星系",
	},
	{
		Name: "Hotpink",
		CnName: "亮粉色",
	},
	{
		Name: "Ivory",
		CnName: "象牙",
	},
	{
		Name: "Jade",
		CnName: "玉",
	},
	{
		Name: "Khaki",
		CnName: "黄褐色",
	},
	{
		Name: "Lavender",
		CnName: "薰衣草",
	},
	{
		Name: "Magenta",
		CnName: "品红",
	},
	{
		Name: "Mahogany",
		CnName: "桃花心木",
	},
	{
		Name: "Mango",
		CnName: "芒果",
	},
	{
		Name: "Maroon",
		CnName: "栗色",
	},
	{
		Name: "Neon",
		CnName: "氖",
	},
	{
		Name: "Tan",
		CnName: "黄褐色",
	},
	{
		Name: "Watermelon red",
		CnName: "西瓜红",
	},
	{
		Name: "Lake Blue",
		CnName: "湖蓝色",
	},
	{
		Name: "Lemon Yellow",
		CnName: "柠檬黄",
	},
	{
		Name: "Army Green",
		CnName: "军绿色",
	},
	{
		Name: "Rose",
		CnName: "玫瑰",
	},
	{
		Name: "Dark blue",
		CnName: "深蓝",
	},
	{
		Name: "Camel",
		CnName: "骆驼",
	},
	{
		Name: "Burgundy",
		CnName: "勃艮第",
	},
	{
		Name: "Light blue",
		CnName: "浅蓝色",
	},
	{
		Name: "Champagne",
		CnName: "香槟酒",
	},
	{
		Name: "Light green",
		CnName: "浅绿色",
	},
	{
		Name: "Dark Brown",
		CnName: "深棕色",
	},
	{
		Name: "Navy Blue",
		CnName: "海军蓝",
	},
	{
		Name: "Light Grey",
		CnName: "浅灰色",
	},
	{
		Name: "Off White",
		CnName: "米白色",
	},
	{
		Name: "Light yellow",
		CnName: "淡黄色",
	},
	{
		Name: "Emerald Green",
		CnName: "翡翠绿",
	},
	{
		Name: "Fluorescent Green",
		CnName: "荧光绿色",
	},
	{
		Name: "Fluorescent Yellow",
		CnName: "荧光黄",
	},
	{
		Name: "Deep green",
		CnName: "深绿色",
	},
	{
		Name: "Rose Gold",
		CnName: "玫瑰金",
	},
	{
		Name: "Neutral",
		CnName: "中性",
	},
	{
		Name: "…",
		CnName: "中性",
	},
	{
		Name: "Peach",
		CnName: "桃子",
	},
	{
		Name: "Fuchsia",
		CnName: "紫红色",
	},
	{
		Name: "Blue Gray",
		CnName: "蓝灰色",
	},
	{
		Name: "Not Specified",
		CnName: "未标明",
	},
	{
		Name: "Orchid Grey",
		CnName: "兰花灰色",
	}}
	
	for _,color:=range colors{
		product.NewAttOption(3,color.Name,color.CnName)
	}
}

func AddSize(){
	type Size struct{ 
		Id int
		 Name string 
		 Blo bool
		}
	sizes:=[]Size{
		{1, "S",false},
		{2, "M",false},
		{3, "L",false},
		{4, "XL",false},
		{5, "XXL",false},
		{6, "3XL",false},
		{7, "4XL",false},
		{8, "5XL",false},
		{9, "6XL",false},
		{10, "7XL",false},
		{11, "8XL",false},
		{12, "9XL",false},
		{13, "10XL",false},
		{14, "One",false},
	}
	for _,size:=range sizes{
		product.NewAttOption(1,size.Name,"")
	}
}

func AddDressShape(){
	type Val struct{
		Name string
		CnName string
		}
	dress_shape:=[]Val{{
		Name: "Babydoll/Smock Dresses",
		CnName:"娃娃装/罩衫连衣裙",
	},
	{
		Name: "Cami/Slip Dresses",
		CnName:"背心/吊带式衬裙",
	},
	{
		Name: "Pencil Dresses",
		CnName: "铅笔裙",
	},
	{
		Name: "Shift Dresses",
		CnName: "宽松直筒连衣裙",
	},
	{
		Name: "Swing Dresses",
		CnName: "摆裙",
	},
	{
		Name: "T-Shirt Dresses",
		CnName:"T恤裙",
	},
	{
		Name: "Tulip Dresses",
		CnName: "郁金香连衣裙",
	},
	{
		Name: "Wrap Dresses",
		CnName: "包装礼服",
	},
	{
		Name: "Asymmetric Dresses",
		CnName: "不对称的裙子",
	},
	{
		Name: "A-Line Dresses",
		CnName: "A线连衣裙",
	},
	{
		Name: "Peplum Dresses",
		CnName: "腰部周围的装饰短裙",
	},
	{
		Name: "Formal Dresses",
		CnName: "正式的礼服",
	},
	{
		Name: "Shirt Dresses",
		CnName: "衬衫裙",
	},
	{
		Name: "Sweater Dresses",
		CnName: "毛衣连衣裙",
	},
	{
		Name: "Tunic Dresses",
		CnName: "束腰连衣裙",
	},
	{
		Name: "Bodycon Dresses",
		CnName: "合身的衣服",
	},
	{
		Name: "Wedding Dresses",
		CnName: "婚纱礼服",
	},
	{
		Name: "Not Specified",
		CnName: "未指定",
	},
}
for _,color:=range dress_shape{
	product.NewAttOption(7,color.Name,color.CnName)
}
}

func AddDressLength(){
	//连衣裙长度
	type Val struct{
		Name string
		CnName string
		}

 dress_lengtg:= []Val{
    {
        Name: "Short",
        CnName:"短裙",
    },
    {
        Name: "Long",
        CnName:"长裙",
    },
    {
        Name: "Knee Length",
        CnName:"及膝长度",

    },
    {
        Name: "Above Knee",
        CnName:"膝盖以上",
    },
    {
        Name: "Mini",
        CnName:"迷你",
    },
    {
        Name: "Below Knee",
        CnName:"膝盖以下",
    },
    {
        Name: "Maxi",
        CnName:"很长",
    },
    {
        Name: "Ankle Length",
        CnName:"脚踝的长度",
    },
    {
        Name: "Floor Length",
        CnName:"拖地长",
    },
    {
        Name: "Not Specified",
        CnName:"未规定的",
    },
}
for _,color:=range dress_lengtg{
	product.NewAttOption(5,color.Name,color.CnName)
}




}


func AddMaterial(){
	type Val struct{
		Name string
		CnName string
		}

		materials:=[]Val{
			{"Cotton","棉"},
			{"Fleece","羊毛"},
			{"Jersey","针织"},
			{"Acrylic Wool","腈纶织物"},
			{"Cashmere","羊绒"},
			{"chiffon","雪纺"},
			{"Crochet","钩针"},
			{"Denim","牛仔布"},
			{"Lace","花边"},
			{"Leather/Suede","皮衣/山羊皮"},
			{"Linen","亚麻布"},
			{"Mesh","网格"},
			{"Mohair","马海毛，安哥拉山羊毛，马海毛织物"},
			{"Other Material","其他材料"},
			{"PU Leather","仿皮"},
			{"Polyester","聚酯"},
			{"Satin","缎子绸缎做的;光滑的;似缎的"},
			{"Silk","丝绸"},
			{"Velvet","天鹅绒"},
			{"Viscose","纤维胶"},
			{"Wool","毛织品，羊毛织物，毛料衣服"},
			
		}
		for _,color:=range materials{
			product.NewAttOption(4,color.Name,color.CnName)
		}
}