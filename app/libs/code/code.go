package code

import (
	"Hello/app/libs/utils"
	"github.com/wenlng/go-captcha/captcha" // https://github.com/wenlng/go-captcha/blob/master/README_zh.md
	"strings"
)

// 验证码类
/**
 * @Example:
	c := code.Code{}
	s := c.CreateStrVerifyCode()
	fmt.Println(s.Thumb)
*/
type Code struct {
	dots  map[int]captcha.CharDot
	err   error
	Id    string // 唯一Id
	Str   string // 验证码
	Thumb string // 缩略图
}

var charActer = "的是我不人在他有这个上们来到时大地为子中你说生国年着就那和要她出也得里后自以会家可下而过天去能对小多然于心学么之都好看起发当没成只如事把还用第样道想作种开美总从无情己面最女但现前些所同日手又行意动方期它头经长儿回位分爱老因很给名法间斯知世什两次使身者被高已亲其进此话常与活正感见明问力理尔点文几定本公特做外孩相西果走将月十实向声车全信重三机工物气每并别真打太新比才便夫再书部水像眼等体却加电主界门利海受听表德少克代员许稜先口由死安写性马光白或住难望教命花结乐色更拉东神记处让母父应直字场平报友关放至张认接告入笑内英军候民岁往何度山觉路带万男边风解叫任金快原吃妈变通师立象数四失满战远格士音轻目条呢病始达深完今提求清王化空业思切怎非找片罗钱紶吗语元喜曾离飞科言干流欢约各即指合反题必该论交终林请医晚制球决窢传画保读运及则房早院量苦火布品近坐产答星精视五连司巴奇管类未朋且婚台夜青北队久乎越观落尽形影红爸百令周吧识步希亚术留市半热送兴造谈容极随演收首根讲整式取照办强石古华諣拿计您装似足双妻尼转诉米称丽客南领节衣站黑刻统断福城故历惊脸选包紧争另建维绝树系伤示愿持千史谁准联妇纪基买志静阿诗独复痛消社算义竟确酒需单治卡幸兰念举仅钟怕共毛句息功官待究跟穿室易游程号居考突皮哪费倒价图具刚脑永歌响商礼细专黄块脚味灵改据般破引食仍存众注笔甚某沉血备习校默务土微娘须试怀料调广蜖苏显赛查密议底列富梦错座参八除跑亮假印设线温虽掉京初养香停际致阳纸李纳验助激够严证帝饭忘趣支春集丈木研班普导顿睡展跳获艺六波察群皇段急庭创区奥器谢弟店否害草排背止组州朝封睛板角况曲馆育忙质河续哥呼若推境遇雨标姐充围案伦护冷警贝著雪索剧啊船险烟依斗值帮汉慢佛肯闻唱沙局伯族低玩资屋击速顾泪洲团圣旁堂兵七露园牛哭旅街劳型烈姑陈莫鱼异抱宝权鲁简态级票怪寻杀律胜份汽右洋范床舞秘午登楼贵吸责例追较职属渐左录丝牙党继托赶章智冲叶胡吉卖坚喝肉遗救修松临藏担戏善卫药悲敢靠伊村戴词森耳差短祖云规窗散迷油旧适乡架恩投弹铁博雷府压超负勒杂醒洗采毫嘴毕九冰既状乱景席珍童顶派素脱农疑练野按犯拍征坏骨余承置臓彩灯巨琴免环姆暗换技翻束增忍餐洛塞缺忆判欧层付阵玛批岛项狗休懂武革良恶恋委拥娜妙探呀营退摇弄桌熟诺宣银势奖宫忽套康供优课鸟喊降夏困刘罪亡鞋健模败伴守挥鲜财孤枪禁恐伙杰迹妹藸遍盖副坦牌江顺秋萨菜划授归浪凡预奶雄升碃编典袋莱含盛济蒙棋端腿招释介烧误"
var chars []string

func init() {
	for _, v := range charActer {
		chars = append(chars, string(v))
	}
}

// 创建字符验证码
func (this *Code) CreateStrVerifyCode() *Code {
	capt := captcha.GetCaptcha()
	chars := strings.Split("abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
	capt.SetRangCheckTextLen(captcha.RangeVal{4, 4})
	_ = capt.SetRangChars(chars) // 设置字符
	// 生成验证码
	this.dots, _, this.Thumb, this.Id, this.err = capt.Generate()
	if this.err != nil {
		utils.ExitError("生成验证码失败", -1)
	}
	this.Str += this.dots[0].Text
	for i := 1; i < len(this.dots); i++ {
		this.Str += this.dots[i].Text
	}
	return this
}

// 创建中文验证码
func (this *Code) CreateCharVerifyCode() *Code {
	capt := captcha.GetCaptcha()
	capt.SetRangCheckTextLen(captcha.RangeVal{3, 3})
	_ = capt.SetRangChars(chars)
	// 生成验证码
	this.dots, _, this.Thumb, this.Id, this.err = capt.Generate()
	if this.err != nil {
		utils.ExitError("生成验证码失败", -1)
	}
	for _, d := range this.dots {
		this.Str += d.Text
	}
	return this
}
