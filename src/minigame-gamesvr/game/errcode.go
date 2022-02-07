package game

const (
	errCodeOk         = 0
	errCodeNotLogin   = 1
	errCodeLogin      = 2
	errCodeLoginOther = 3
	errCodeDbErr      = 4
	errGoldNotEnough  = 4
	errCodeHacker     = 5

	errCodeDaxiguaStart = 10001

	errCodeBpbxqDaySign = 20002
	errCodeBpbxqBuySkin = 20003
	errCodeBpbxqUseSkin = 20004
	errCodeBpbxqRevive  = 20005
	errCodeBpbxqDiefly  = 20006
	errCodeBpbxqReward  = 20007

	errCodeCaichengyuStart        = iota + 1
	errCodeCaichengyuUpgradeRole  = iota + 1
	errCodeCaichengyuUpgradeBuild = iota + 1
	errCodeCaichengyuRoundResult  = iota + 1
	errCodeCaichengyuRoundStart   = iota + 1
	errCodeCaichengyuTip          = iota + 1

	errCodeFangkuainiaoBuySkin     = iota + 1
	errCodeFangkuainiaoUseSkin     = iota + 1
	errCodeFangkuainiaoDaySign     = iota + 1
	errCodeFangkuainiaoGetCoin     = iota + 1
	errCodeFangkuainiaoRoundResult = iota + 1
	errCodeFangkuainiaoRoundStart  = iota + 1

	errCodeGongjianchuanshuUnlockSkin = 50001
	errCodeGongjianchuanshuUseSkin    = 50002

	errCodePaopaolongCostDraw     = 60001
	errCodePaopaolongFreeDraw     = 60002
	errCodePaopaolongSign         = 60003
	errCodePaopaolongShopCoinFree = 60004
	errCodePaopaolongShopBuy      = 60005
	errCodePaopaolongStart        = 60006
	errCodePaopaolongNewPack      = 60007
	errCodePaopaolongUseItem      = 60008

	errCodeTanchishedazuozhanBuySkin = 70001
	errCodeTanchishedazuozhanUseSkin = 70002

	errCodeTiantianpaokuBuySkin = 80001
	errCodeTiantianpaokuUseSkin = 80002

	errCodeHuanlemaomibeiFreeDraw    = 90001
	errCodeHuanlemaomibeiSign        = 90002
	errCodeHuanlemaomibeiFly         = 90003
	errCodeHuanlemaomibeiUseSkin     = 90004
	errCodeHuanlemaomibeiBuySkin     = 90005
	errCodeHuanlemaomibeiLevelUnlock = 90006
	errCodeHuanlemaomibeiLevelCoin   = 90007
)
