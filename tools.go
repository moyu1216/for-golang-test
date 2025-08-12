package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// WeatherTool 天气查询工具
type WeatherTool struct{}

func (w *WeatherTool) Name() string {
	return "weather_query"
}

func (w *WeatherTool) Description() string {
	return "查询指定城市的天气信息，包括当前天气状况、温度、湿度等"
}

func (w *WeatherTool) Execute(args map[string]interface{}) (string, error) {
	query, ok := args["query"].(string)
	if !ok {
		return "", fmt.Errorf("缺少查询参数")
	}

	// 模拟天气数据（实际应用中应该调用真实的天气API）
	cities := map[string]map[string]interface{}{
		"北京": {"temp": "15", "condition": "晴天", "humidity": "45%", "wind": "微风"},
		"上海": {"temp": "18", "condition": "多云", "humidity": "60%", "wind": "东南风3级"},
		"广州": {"temp": "25", "condition": "小雨", "humidity": "75%", "wind": "南风2级"},
		"深圳": {"temp": "26", "condition": "晴天", "humidity": "65%", "wind": "海风"},
		"成都": {"temp": "20", "condition": "阴天", "humidity": "70%", "wind": "无风"},
		"杭州": {"temp": "22", "condition": "晴天", "humidity": "55%", "wind": "东风2级"},
	}

	// 查找城市
	for city, weather := range cities {
		if strings.Contains(query, city) {
			return fmt.Sprintf(`%s天气信息：
🌡️ 温度：%s°C
☁️ 天气：%s
💧 湿度：%s
🌪️ 风况：%s
📅 更新时间：%s`,
				city,
				weather["temp"],
				weather["condition"],
				weather["humidity"],
				weather["wind"],
				time.Now().Format("2006-01-02 15:04")), nil
		}
	}

	// 如果没有找到具体城市，返回通用建议
	return fmt.Sprintf("抱歉，暂时无法获取 '%s' 的准确天气信息。建议您：\n1. 出行前查看当地天气预报\n2. 准备雨具以防天气变化\n3. 根据季节调整着装", query), nil
}

// AttractionTool 景点推荐工具
type AttractionTool struct{}

func (a *AttractionTool) Name() string {
	return "attraction_recommend"
}

func (a *AttractionTool) Description() string {
	return "根据城市或地区推荐热门景点、特色景观和旅游路线"
}

func (a *AttractionTool) Execute(args map[string]interface{}) (string, error) {
	query, ok := args["query"].(string)
	if !ok {
		return "", fmt.Errorf("缺少查询参数")
	}

	attractions := map[string][]map[string]string{
		"北京": {
			{"name": "故宫博物院", "type": "历史文化", "rating": "⭐⭐⭐⭐⭐", "tip": "建议预约门票，游览时间3-4小时"},
			{"name": "天安门广场", "type": "地标建筑", "rating": "⭐⭐⭐⭐⭐", "tip": "早晨观看升旗仪式很震撼"},
			{"name": "长城(八达岭)", "type": "世界遗产", "rating": "⭐⭐⭐⭐⭐", "tip": "体力消耗大，建议穿舒适鞋子"},
			{"name": "颐和园", "type": "皇家园林", "rating": "⭐⭐⭐⭐", "tip": "春秋季节风景最佳"},
		},
		"上海": {
			{"name": "外滩", "type": "城市地标", "rating": "⭐⭐⭐⭐⭐", "tip": "夜景最美，建议傍晚前往"},
			{"name": "东方明珠", "type": "现代建筑", "rating": "⭐⭐⭐⭐", "tip": "登塔观景需提前购票"},
			{"name": "豫园", "type": "古典园林", "rating": "⭐⭐⭐⭐", "tip": "体验传统江南园林文化"},
			{"name": "南京路步行街", "type": "商业街区", "rating": "⭐⭐⭐", "tip": "购物和品尝小吃的好地方"},
		},
		"广州": {
			{"name": "广州塔", "type": "现代地标", "rating": "⭐⭐⭐⭐⭐", "tip": "夜晚灯光秀非常壮观"},
			{"name": "陈家祠", "type": "传统建筑", "rating": "⭐⭐⭐⭐", "tip": "岭南建筑艺术的代表"},
			{"name": "沙面岛", "type": "历史文化", "rating": "⭐⭐⭐⭐", "tip": "欧式建筑群，适合拍照"},
			{"name": "白云山", "type": "自然风光", "rating": "⭐⭐⭐", "tip": "爬山健身，俯瞰城市"},
		},
	}

	for city, places := range attractions {
		if strings.Contains(query, city) {
			result := fmt.Sprintf("🏛️ %s热门景点推荐：\n\n", city)
			for i, place := range places {
				result += fmt.Sprintf("%d. 📍 %s\n   类型：%s | 评价：%s\n   💡 %s\n\n",
					i+1, place["name"], place["type"], place["rating"], place["tip"])
			}
			return result, nil
		}
	}

	return fmt.Sprintf("暂时没有 '%s' 的具体景点信息，但我可以为您提供以下通用旅游建议：\n1. 提前查询当地著名景点\n2. 关注景点开放时间和门票信息\n3. 安排合理的游览路线\n4. 准备相机记录美好时光", query), nil
}

// HotelTool 酒店搜索工具
type HotelTool struct{}

func (h *HotelTool) Name() string {
	return "hotel_search"
}

func (h *HotelTool) Description() string {
	return "搜索推荐酒店住宿，包括不同价格区间和类型的选择"
}

func (h *HotelTool) Execute(args map[string]interface{}) (string, error) {
	query, ok := args["query"].(string)
	if !ok {
		return "", fmt.Errorf("缺少查询参数")
	}

	hotels := map[string][]map[string]string{
		"北京": {
			{"name": "北京饭店", "type": "豪华酒店", "price": "¥800-1200/晚", "location": "王府井", "rating": "⭐⭐⭐⭐⭐"},
			{"name": "如家酒店", "type": "经济型", "price": "¥200-350/晚", "location": "各区域", "rating": "⭐⭐⭐"},
			{"name": "全季酒店", "type": "中档酒店", "price": "¥400-600/晚", "location": "朝阳区", "rating": "⭐⭐⭐⭐"},
			{"name": "青年旅社", "type": "背包客", "price": "¥80-150/晚", "location": "东城区", "rating": "⭐⭐⭐"},
		},
		"上海": {
			{"name": "和平饭店", "type": "历史酒店", "price": "¥1000-1500/晚", "location": "外滩", "rating": "⭐⭐⭐⭐⭐"},
			{"name": "锦江之星", "type": "经济型", "price": "¥250-400/晚", "location": "各区域", "rating": "⭐⭐⭐"},
			{"name": "亚朵酒店", "type": "精品酒店", "price": "¥500-800/晚", "location": "静安区", "rating": "⭐⭐⭐⭐"},
			{"name": "上海青旅", "type": "青年旅社", "price": "¥100-200/晚", "location": "徐汇区", "rating": "⭐⭐⭐"},
		},
	}

	for city, hotelList := range hotels {
		if strings.Contains(query, city) {
			result := fmt.Sprintf("🏨 %s酒店推荐：\n\n", city)
			for i, hotel := range hotelList {
				result += fmt.Sprintf("%d. 🏩 %s (%s)\n   💰 价格：%s\n   📍 位置：%s | 评价：%s\n\n",
					i+1, hotel["name"], hotel["type"], hotel["price"], hotel["location"], hotel["rating"])
			}
			result += "💡 预订建议：\n• 提前预订可享受优惠\n• 查看酒店评价和设施\n• 考虑交通便利性\n• 注意退改政策"
			return result, nil
		}
	}

	return fmt.Sprintf("正在为您搜索 '%s' 的酒店信息...\n\n🏨 通用住宿建议：\n1. 根据预算选择合适类型\n2. 考虑地理位置和交通\n3. 查看用户评价和设施\n4. 比较多个平台价格", query), nil
}

// RouteTool 路线规划工具
type RouteTool struct{}

func (r *RouteTool) Name() string {
	return "route_planning"
}

func (r *RouteTool) Description() string {
	return "规划旅游路线，包括交通方式、时间安排和行程建议"
}

func (r *RouteTool) Execute(args map[string]interface{}) (string, error) {
	query, ok := args["query"].(string)
	if !ok {
		return "", fmt.Errorf("缺少查询参数")
	}

	// 简单的路线规划逻辑
	routes := map[string]map[string]string{
		"北京一日游": {
			"morning":   "09:00 天安门广场 → 10:30 故宫博物院",
			"afternoon": "14:00 北海公园 → 16:00 南锣鼓巷",
			"evening":   "18:00 王府井大街 (晚餐购物)",
			"transport": "地铁+步行，建议购买一日交通卡",
		},
		"上海两日游": {
			"day1": "外滩 → 南京路 → 豫园 → 城隍庙",
			"day2": "东方明珠 → 陆家嘴 → 田子坊 → 新天地",
			"transport": "地铁为主，部分景点间可步行",
		},
	}

	// 生成通用路线建议
	tips := []string{
		"🚇 优先使用公共交通，环保又经济",
		"⏰ 合理安排时间，避免行程过于紧张",
		"📍 相近景点可安排在同一天游览",
		"🍽️ 预留用餐和休息时间",
		"📱 下载导航和交通APP",
		"🎫 提前购买景点门票",
	}

	result := fmt.Sprintf("🗺️ 为您规划 '%s' 的旅游路线：\n\n", query)
	
	// 检查是否有预设路线
	found := false
	for routeName, details := range routes {
		if strings.Contains(strings.ToLower(query), strings.ToLower(routeName[:2])) {
			found = true
			result += fmt.Sprintf("📋 推荐行程：%s\n\n", routeName)
			for key, value := range details {
				if key == "transport" {
					result += fmt.Sprintf("🚌 交通建议：%s\n\n", value)
				} else {
					result += fmt.Sprintf("• %s：%s\n", key, value)
				}
			}
			break
		}
	}

	if !found {
		result += "📝 通用路线规划建议：\n\n"
	}

	// 添加通用建议
	result += "💡 出行小贴士：\n"
	for _, tip := range tips {
		result += fmt.Sprintf("   %s\n", tip)
	}

	return result, nil
}

// FoodTool 美食推荐工具
type FoodTool struct{}

func (f *FoodTool) Name() string {
	return "food_recommend"
}

func (f *FoodTool) Description() string {
	return "推荐当地特色美食、餐厅和小吃街"
}

func (f *FoodTool) Execute(args map[string]interface{}) (string, error) {
	query, ok := args["query"].(string)
	if !ok {
		return "", fmt.Errorf("缺少查询参数")
	}

	foods := map[string][]map[string]string{
		"北京": {
			{"name": "北京烤鸭", "place": "全聚德、便宜坊", "price": "¥150-300", "desc": "外酥内嫩，配薄饼和甜面酱"},
			{"name": "炸酱面", "place": "老北京面馆", "price": "¥25-40", "desc": "地道老北京味道"},
			{"name": "豆汁", "place": "护国寺小吃", "price": "¥8-15", "desc": "传统北京小吃，口味独特"},
			{"name": "糖葫芦", "place": "王府井小吃街", "price": "¥10-20", "desc": "酸甜开胃，冬季必尝"},
		},
		"上海": {
			{"name": "小笼包", "place": "南翔馒头店", "price": "¥30-50", "desc": "皮薄汁多，上海经典"},
			{"name": "生煎包", "place": "大壶春", "price": "¥20-35", "desc": "底部焦脆，馅料鲜美"},
			{"name": "白切鸡", "place": "沪上人家", "price": "¥40-80", "desc": "嫩滑鲜美，配姜蓉蘸料"},
			{"name": "糖醋排骨", "place": "本帮菜馆", "price": "¥35-60", "desc": "酸甜开胃，色泽红亮"},
		},
		"广州": {
			{"name": "白切鸡", "place": "陶陶居", "price": "¥45-80", "desc": "粤菜经典，嫩滑鲜甜"},
			{"name": "虾饺", "place": "点都德", "price": "¥20-35", "desc": "晶莹剔透，鲜虾饱满"},
			{"name": "干炒牛河", "place": "广州酒家", "price": "¥25-40", "desc": "河粉爽滑，牛肉嫩香"},
			{"name": "双皮奶", "place": "南信甜品", "price": "¥15-25", "desc": "奶香浓郁，口感顺滑"},
		},
	}

	for city, foodList := range foods {
		if strings.Contains(query, city) {
			result := fmt.Sprintf("🍜 %s特色美食推荐：\n\n", city)
			for i, food := range foodList {
				result += fmt.Sprintf("%d. 🥢 %s\n   📍 推荐：%s\n   💰 价格：%s\n   📝 %s\n\n",
					i+1, food["name"], food["place"], food["price"], food["desc"])
			}
			result += "🍽️ 用餐建议：\n• 选择人气较高的老字号\n• 注意用餐高峰时间\n• 尝试当地特色小吃街\n• 关注食物卫生安全"
			return result, nil
		}
	}

	// 随机推荐一些通用美食建议
	generalTips := []string{
		"🏪 寻找当地老字号餐厅",
		"🍽️ 品尝特色街边小吃",
		"👥 选择人多的店铺，通常味道不错",
		"📱 可以使用美食APP查看评价",
		"🥢 注意饮食习惯和忌口",
		"💧 选择信誉好的餐厅保证卫生",
	}

	rand.Seed(time.Now().UnixNano())
	selectedTips := make([]string, 3)
	for i := 0; i < 3; i++ {
		selectedTips[i] = generalTips[rand.Intn(len(generalTips))]
	}

	result := fmt.Sprintf("🍴 '%s' 美食探索建议：\n\n", query)
	result += "💡 通用美食攻略：\n"
	for _, tip := range selectedTips {
		result += fmt.Sprintf("   %s\n", tip)
	}

	return result, nil
}