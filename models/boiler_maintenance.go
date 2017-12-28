package models

import "time"

type BoilerMaintenance struct {
	MyUidObject

	Boiler				*Boiler			`orm:"rel(fk);null;index"`

	InspectDate			time.Time		`orm:"type(datetime);index"`
	Inspector			string
	Content				string
	Attachment			string
	IsDone				bool

	/*
	BurnerCoolant			int			// 燃烧器冷却水温度是否过高
	BurnerFireBrick			int			// 燃烧器炉膛耐火砖是否有脱落
	BurnerFireHold			int			// 燃烧器火口是否有变形
	BurnerImportProcess		int			// 进料蛟龙运转是否有异响
	BurnerBlower			int			// 鼓风机转动是否有异响
	BurnerWindChain			int			// 风料连锁控制情况是否正常
	BurnerPressureDiffTransducer	int			// 燃烧器差压传感器是否稳定

	ImportGrateBucket		int			// 斗提机运转是否正常
	ImportGrateStarFeeder		int			// 星型给料器给料情况是否匹配
	ImportGrateRevolutions		int			// 炉排运转数度调节是否与负荷匹配
	ImportGrateLeak			int			// 炉排漏料是否过多
	ImportGrateBurnOut		int			// 燃料燃尽与炉排速度匹配度
	ImportGrateAirDoor		int			// 风门是否可开合
	ImportGrateAirDoorOptimize	int			// 风门开合与燃烧是否最优

	WaterSoftenerBlowdown		int			// 对软化水器进行定期的排污
	*/

	Burner			string			// 燃烧器
	/**
	 * 燃烧器冷却水温度是否过高
	 * 燃烧器炉膛耐火砖是否有脱落
	 * 燃烧器火口是否有变形
	 * 进料蛟龙运转是否有异响
	 * 鼓风机转动是否有异响
	 * 风料连锁控制情况是否正常
	 * 燃烧器差压传感器是否稳定
	 */
	ImportGrate		string			// 进料及炉排
	/**
	 * 斗提机运转是否正常
	 * 星型给料器给料情况是否匹配
	 * 炉排运转数度调节是否与负荷匹配
	 * 炉排漏料是否过多
	 * 燃料燃尽与炉排速度匹配度
	 * 风门是否可开合
	 * 风门开合与燃烧是否最优
	 */
	WaterSoftener		string			// 软水器
	/**
	 * 对软化水器进行定期的排污
	 * 对软化水器灵敏度的检查
	 * 软水箱中的水硬度超标
	 */
	WaterPump		string			// 水泵
	/**
	 * 检查泵体是否完好
	 * 检查水泵水流方向指示是否明确清晰
	 * 检查水泵有无渗漏情况
	 * 对水泵补充润滑油，若油变色、变质、有杂质应及时更换。
	 */
	BoilerBody		string			// 锅炉本体
	/**
	 * 水位计1是否可清楚辨认
	 * 水位计2是否稳定显示
	 * 压力开关是否稳定
	 * 排污阀是否锈蚀
	 * 安全阀是否有漏气现象
	 * 锅炉压力传感器是否稳定
	 */
	EnergySaver		string			// 节能器
	/**
	 * 节能器前后烟温表是否正常
	 * 节能器前后水温表是否正常
	 * 节能器是否密封严密
	 */
	AirPreHeater		string			// 空预器
	/**
	 * 前后烟温表是否正常
	 * 前后空气温是否正常
	 * 是否有密封严密
	 */
	DustCatcher		string			// 除尘器
	/**
	 * 密封性是否良好
	 * 设备前后压差值是否过高
	 * 排气电池阀是否完好
	 */
	DraughtFan		string			// 引风机
	/**
	 * 引风机异响情况
	 * 引风机出力较前是否下降
	 * 引风机冷却水是否合格
	 */
}
