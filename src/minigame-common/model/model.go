package model

type User struct {
	Userid             int64  `json:"userid" db:"userid" `                         // 会员ID
	Openid             string `json:"openid" db:"openid" `                         // 第三方唯一ID
	Nickname           string `json:"nickname" db:"nickname" `                     // 第三方会员昵称
	Avatar             string `json:"avatar" db:"avatar" `                         // 第三方会员昵称
	ClientIp           string `json:"client_ip" db:"client_ip" `                   // 第三方会员昵称
	Coin               int64  `json:"coin" db:"coin"`                              // 金币
	Score              int64  `json:"score" db:"score"`                            // 分数
	Level              int64  `json:"level" db:"level"`                            // 关卡
	Medal              int64  `json:"medal" db:"medal"`                            // 关卡
	Diamond            int64  `json:"diamond" db:"diamond"`                        // 钻石
	Createtime         int64  `json:"createtime" db:"createtime" `                 // 创建时间
	Updatetime         int64  `json:"updatetime" db:"updatetime" `                 // 更新时间
	Logintime          int64  `json:"logintime" db:"logintime" `                   // 登录时间
	Logouttime         int64  `json:"logouttime" db:"logouttime" `                 // 登出时间
	RoundStartTime     int64  `json:"roundstarttime" db:"roundstarttime" `         // 游戏开始时间
	RoundEndTime       int64  `json:"roundendtime" db:"roundendtime" `             // 游戏结束时间
	LastRoundStartTime int64  `json:"lastroundstarttime" db:"lastroundstarttime" ` // 游戏开始时间
	LastRoundEndTime   int64  `json:"lastroundendtime" db:"lastroundendtime" `     // 游戏结束时间
}

type GateStat struct {
	NodeName     string `json:"nodename" db:"nodename" `         // 节点名
	NodeFullName string `json:"nodefullname" db:"nodefullname" ` // 节点全名
	Host         string `json:"host" db:"host" `                 // Host
	Port         int32  `json:"port" db:"port" `                 // Port
	OnlineNum    int32  `json:"onlinenum" db:"onlinenum" `       // 在线人数
	Createtime   int64  `json:"createtime" db:"createtime" `     // 创建时间
	Updatetime   int64  `json:"updatetime" db:"updatetime" `     // 更新时间
}

type Daxigua struct {
	Userid int64 `json:"userid" db:"userid" ` // 会员ID
}

type Benpaobaxiaoqie struct {
	Userid       int64  `json:"userid" db:"userid" `            // 会员ID
	Skin         int32  `json:"skin" db:"skin"`                 // 皮肤
	SkinArr      string `json:"skinarr" db:"skinarr"`           // 皮肤列表
	LastSignTime int64  `json:"lastsigntime" db:"lastsigntime"` //签到时间
	SignTimes    int64  `json:"signtimes" db:"signtimes"`       //签到次数
}

type Caichengyu struct {
	Userid     int64  `json:"userid" db:"userid" `        // 会员ID
	RoleLevel  int32  `json:"rolelevel" db:"rolelevel"`   // 角色等级
	BuildLevel int32  `json:"buildlevel" db:"buildlevel"` // 建筑等级
	Level      int64  `json:"level" db:"level"`           //当前关卡
	LevelType  int32  `json:"leveltype" db:"leveltype"`   //关卡状态
	LevelTip   int32  `json:"leveltip" db:"leveltip"`     //关卡提示次数
	Hp         int32  `json:"hp" db:"hp"`                 //体力
	HpDate     int64  `json:"hpdate" db:"hpdate"`         //上次体力恢复时间
	GetHpCount int32  `json:"gethpcount" db:"gethpcount"` //领取体力次数
	GetHpDay   string `json:"gethpday" db:"gethpday"`     //领取体力日期
}

type Fangkuainiao struct {
	Userid       int64  `json:"userid" db:"userid" `            // 会员ID
	Level        int64  `json:"level" db:"level"`               //当前关卡
	BirdId       int32  `json:"birdid" db:"birdid"`             // 皮肤
	BirdArr      string `json:"birdarr" db:"birdarr"`           // 皮肤列表
	SignTime     int64  `json:"signtime" db:"signtime"`         //签到时间
	SignDay      int64  `json:"signday" db:"signday"`           //签到次数
	GetGoldCount int32  `json:"getgoldcount" db:"getgoldcount"` //领取金币次数
	GetGoldTime  int64  `json:"getgoldtime" db:"getgoldtime"`   //领取金币日期
}

type Gongjianchuanshu struct {
	Userid  int64  `json:"userid" db:"userid" `  // 会员ID
	Level   int64  `json:"level" db:"level"`     //当前关卡
	SkinId  int32  `json:"skinid" db:"skinid"`   // 皮肤
	SkinArr string `json:"skinarr" db:"skinarr"` // 皮肤列表
	ShopArr string `json:"shoparr" db:"shoparr"` // 商品列表
}

type Paopaolong struct {
	Userid               int64  `json:"userid" db:"userid" `                                  //会员ID
	ItemArr              string `json:"itemarr" db:"itemarr"`                                 //道具列表
	Level                int64  `json:"level" db:"level"`                                     //当前关卡
	Hp                   int32  `json:"hp" db:"hp"`                                           //体力
	FreedrawTime         int64  `json:"freedraw_time" db:"freedraw_time"`                     //免费抽奖时间
	NewPackRedeemed      int32  `json:"new_pack_redeemed" db:"new_pack_redeemed"`             //新手礼包领取状态:0=末领取
	ShopFreeDiamondTime  int64  `json:"shop_free_diamond_time" db:"shop_free_diamond_time"`   //商城免费领取金币状态
	ShopFreeDiamondTime2 int64  `json:"shop_free_diamond_time2" db:"shop_free_diamond_time2"` //商城免费领取金币状态2
	LastSignTime         int64  `json:"last_sign_time" db:"last_sign_time"`                   //上次签到时间
	SignedTime           int32  `json:"signed_time" db:"signed_time"`                         //签到天数
	FirstSignTime        int64  `json:"first_sign_time" db:"first_sign_time"`                 //签到开启时间
}

type PaopaolongLevel struct {
	Userid int64 `json:"userid" db:"userid" ` //会员ID
	Level  int64 `json:"level" db:"level"`    //关卡
	Sec    int64 `json:"sec" db:"sec"`        //时间
	Lose   int32 `json:"lose" db:"lose"`      //失败次数
	Score  int64 `json:"score" db:"score"`    //分数
	Star   int64 `json:"star" db:"star"`      //星星数
}

type Tanchishedazuozhan struct {
	Userid  int64  `json:"userid" db:"userid" `  // 会员ID
	SkinId  int32  `json:"skinid" db:"skinid"`   // 皮肤
	SkinArr string `json:"skinarr" db:"skinarr"` // 皮肤列表
}

type Tiantianpaoku struct {
	Userid  int64  `json:"userid" db:"userid" `  // 会员ID
	SkinId  int32  `json:"skinid" db:"skinid"`   // 皮肤
	SkinArr string `json:"skinarr" db:"skinarr"` // 皮肤列表
}

type Huanlemaomibei struct {
	Userid       int64  `json:"userid" db:"userid" `              // 会员ID
	FreedrawTime int64  `json:"freedraw_time" db:"freedraw_time"` //免费抽奖时间
	InkArr       string `json:"inkarr" db:"inkarr"`               //墨水皮肤列表
	CupArr       string `json:"cuparr" db:"cuparr"`               //酒杯皮肤列表
	InkId        int32  `json:"inkid" db:"inkid"`                 //墨水皮肤
	CupId        int32  `json:"cupid" db:"cupid"`                 //酒杯皮肤
	SignDay      int32  `json:"signday" db:"signday"`             //签到天数
	LastSignTime int64  `json:"lastsigntime" db:"lastsigntime"`   //上次签到时间
	SignChecked  int32  `json:"signchecked" db:"signchecked"`     //是否检查签到跨天
	FlyTime      int64  `json:"fly_time" db:"fly_time"`           //幸运蝴蝶抽奖时间
	OfflineTime  int64  `json:"offline_time" db:"offline_time"`   //离线收益时间
	Hp           int32  `json:"hp" db:"hp"`                       //体力
}

type HuanlemaomibeiLevel struct {
	Userid  int64 `json:"userid" db:"userid" `    //会员ID
	Section int32 `json:"section" db:"section"`   //章节
	Level   int32 `json:"level" db:"level"`       //关卡
	Star    int32 `json:"star" db:"star"`         //星星数
	Unlock  int32 `json:"unlocked" db:"unlocked"` //解锁
	Coin    int64 `json:"coin" db:"coin"`         //领取金币
}

type Yangzhunongchang struct {
	Userid       int64  `json:"userid" db:"userid" `                   // 会员ID
	FarmLv       int32  `json:"farm_lv" db:"farm_lv" `                 // 猪场等级
	FarmLvName   string `json:"farm_lv_name" db:"farm_lv_name" `       // 猪场等级称号
	FarmLvExp    int32  `json:"farm_lv_exp" db:"farm_lv_exp" `         // 猪场等级累计经验
	FarmLvExpCur int32  `json:"farm_lv_exp_cur" db:"farm_lv_exp_cur" ` // 猪场等级当前等级累计经验
	AwardNum     int32  `json:"award_num" db:"award_num" `             // 抽奖次数
	AwardTime    int64  `json:"award_time" db:"award_time" `           // 抽奖时间
}

type YangzhunongchangItem struct {
	Userid int64 `json:"userid" db:"userid" ` //会员ID
	Id     int32 `json:"id" db:"id"`          //道具id
	Num    int32 `json:"num" db:"num"`        //道具数量
}

type YangzhunongchangPig struct {
	Userid     int64  `json:"userid" db:"userid" `        //会员ID
	Id         string `json:"id" db:"id"`                 //猪id
	Data       string `json:"data" db:"data"`             //猪数据
	Createtime int64  `json:"createtime" db:"createtime"` //创建时间
}

type YangzhunongchangBreedPig struct {
	Userid int64  `json:"userid" db:"userid" ` //会员ID
	Id     string `json:"id" db:"id"`          //猪id
	Data   string `json:"data" db:"data"`      //猪数据
}

type YangzhunongchangUsu struct {
	Userid int64 `json:"userid" db:"userid" ` //会员ID
	Id     int32 `json:"id" db:"id"`          //图鉴id
}

type YangzhunongchangFood struct {
	Userid int64  `json:"userid" db:"userid" ` //会员ID
	Id     string `json:"id" db:"id"`          //食物id
	Data   string `json:"data" db:"data"`      //食物数据
}

type YangzhunongchangTask struct {
	Userid int64  `json:"userid" db:"userid" ` //会员ID
	Id     string `json:"id" db:"id"`          //任务id
	Index  int32  `json:"index" db:"index"`    //下标
	Count  string `json:"count" db:"count"`    //进度
}
