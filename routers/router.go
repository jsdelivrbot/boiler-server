package routers

import (
	"github.com/AzureRelease/boiler-server/controllers"
	"github.com/AzureTech/goazure"
)

func init() {



	goazure.Router("/", &controllers.MainController{})
	goazure.Router("/login/", &controllers.MainController{})
	goazure.Router("/admin", &controllers.AdminController{})

	goazure.Router("/user", &controllers.UserController{})
	goazure.Router("/user_login/", &controllers.UserController{}, "post:UserLogin")
	goazure.Router("/user_register/", &controllers.UserController{}, "post:UserRegister")
	goazure.Router("/user_logout/", &controllers.UserController{}, "post:UserLogout")

	goazure.Router("/user_list/", &controllers.UserController{}, "get:UserList")
	goazure.Router("/user_roles/", &controllers.UserController{}, "get:UserRoleList")

	goazure.Router("/data_update", &controllers.DataUpdateController{})

	goazure.Router("/user_update/", &controllers.UserController{}, "post:UserUpdate")
	goazure.Router("/user_update_admin/", &controllers.UserController{}, "post:UserUpdateAdmin")
	goazure.Router("/user_active/", &controllers.UserController{}, "post:UserActive")
	goazure.Router("/user_delete/", &controllers.UserController{}, "post:UserDelete")

	goazure.Router("/user_profile_update/", &controllers.UserController{}, "post:UserProfileUpdate")
	goazure.Router("/user_image_upload/", &controllers.UserController{}, "post:UserImageUpload")
	goazure.Router("/user_password_update/", &controllers.UserController{}, "post:UserPasswordUpdate")

	goazure.Router("/runtime_parameter_list/", &controllers.ParameterController{}, "get:RuntimeParameterList")
	goazure.Router("/runtime_parameter_update/", &controllers.ParameterController{}, "post:RuntimeParameterUpdate")
	goazure.Router("/runtime_parameter_delete/", &controllers.ParameterController{}, "post:RuntimeParameterDelete")

	goazure.Router("/channel_config_list/", &controllers.ParameterController{}, "post:ChannelConfigList")
	goazure.Router("/channel_config_matrix/", &controllers.ParameterController{}, "post:ChannelConfigMatrix")
	//修改
	goazure.Router("/channel_config_update/", &controllers.ParameterController{}, "post:ChannelIssuedUpdate")

	goazure.Router("/organization_list/", &controllers.OrganizationController{}, "get:OrganizationList")
	goazure.Router("/organization_type_list/", &controllers.OrganizationController{}, "get:OrganizationTypeList")
	goazure.Router("/organization_update/", &controllers.OrganizationController{}, "post:OrganizationUpdate")
	goazure.Router("/organization_delete/", &controllers.OrganizationController{}, "post:OrganizationDelete")

	goazure.Router("/boiler_count/", &controllers.BoilerController{}, "get:BoilerCount")
	goazure.Router("/boiler_list/", &controllers.BoilerController{}, "get:BoilerList")
	goazure.Router("/boiler_update/", &controllers.BoilerController{}, "post:BoilerUpdate")
	goazure.Router("/boiler_delete/", &controllers.BoilerController{}, "post:BoilerDelete")
	goazure.Router("/boiler_bind/", &controllers.BoilerController{}, "post:BoilerBind")
	goazure.Router("/boiler_unbind/", &controllers.BoilerController{}, "post:BoilerUnbind")

	goazure.Router("/boiler_medium_list/", &controllers.BoilerController{}, "get:BoilerMediumList")
	goazure.Router("/boiler_form_list/", &controllers.BoilerController{}, "get:BoilerFormList")
	goazure.Router("/boiler_fuel_list", &controllers.FuelController{}, "get:FuelList")
	goazure.Router("/boiler_fuel_type_list", &controllers.FuelController{}, "get:FuelTypeList")

	goazure.Router("/boiler_config/", &controllers.BoilerController{}, "post:GetBoilerConfig")
	goazure.Router("/boiler_config_set/", &controllers.BoilerController{}, "post:SetBoilerConfig")
	goazure.Router("/boiler/state/is_burning/", &controllers.BoilerController{}, "get:BoilerIsBurning")
	goazure.Router("/boiler/state/is_Online",&controllers.BoilerController{},"get:BoilerIsOnline")
	goazure.Router("/boiler/state/has_subscribed/", &controllers.BoilerController{}, "get:BoilerHasSubscribed")
	goazure.Router("/boiler/state/set_subscribe/", &controllers.BoilerController{}, "post:BoilerSetSubscribe")
	goazure.Router("/boiler/state/has_channel_custom/", &controllers.ParameterController{}, "get:BoilerHasChannelCustom")

	goazure.Router("/boiler_message_send/", &controllers.BoilerController{}, "post:BoilerMessageSend")

	//goazure.Router("/boiler_runtime/", &controllers.BoilerController{}, "post:GetBoilerRuntime")
	goazure.Router("/boiler_runtime_count/", &controllers.RuntimeController{}, "get:BoilerRuntimeCount")
	goazure.Router("/boiler_runtime_list/", &controllers.RuntimeController{}, "post:BoilerRuntimeList")
	goazure.Router("/boiler_runtime_history/", &controllers.RuntimeController{}, "post:BoilerRuntimeHistory")
	goazure.Router("/boiler_runtime_instants/", &controllers.RuntimeController{}, "post:BoilerRuntimeInstants")
	goazure.Router("/boiler_runtime_daily/", &controllers.RuntimeController{}, "post:BoilerRuntimeDaily")
	goazure.Router("/boiler_runtime_daily_total/", &controllers.RuntimeController{}, "*:BoilerRuntimeDailyTotal")

	goazure.Router("/boiler_status_running/", &controllers.RuntimeController{}, "get:BoilerStatusRunningDuration")

	goazure.Router("/boiler_evaporate_rank/", &controllers.RuntimeController{}, "get:BoilerHeatRank")


	goazure.Router("/boiler_calculate_parameter/", &controllers.BoilerController{}, "get:BoilerCalculateParameter")

	goazure.Router("/boiler_calculate/", &controllers.CalculateController{}, "post:BoilerCalculate")
	goazure.Router("/boiler_calculate_parameter_update/", &controllers.CalculateController{}, "post:BoilerCalculateParameterUpdate")

	goazure.Router("/boiler_maintain_list/", &controllers.BoilerController{}, "get:BoilerMaintainList")
	goazure.Router("/boiler_maintain_update/", &controllers.BoilerController{}, "post:BoilerMaintainUpdate")
	goazure.Router("/boiler_maintain_delete/", &controllers.BoilerController{}, "post:BoilerMaintainDelete")

	goazure.Router("/terminal_list/", &controllers.TerminalController{}, "get:TerminalList")
	goazure.Router("/terminal_reset/", &controllers.TerminalController{}, "post:TerminalReset")
	goazure.Router("/terminal_update/", &controllers.TerminalController{}, "post:TerminalUpdate")
	goazure.Router("/terminal_delete/", &controllers.TerminalController{}, "post:TerminalDelete")

	goazure.Router("/terminal_config/", &controllers.TerminalController{}, "post:TerminalConfig")

	goazure.Router("/alarm_rule_list/", &controllers.AlarmController{}, "get:AlarmRuleList")
	goazure.Router("/alarm_rule_update/", &controllers.AlarmController{}, "post:AlarmRuleUpdate")
	goazure.Router("/alarm_rule_delete/", &controllers.AlarmController{}, "post:AlarmRuleDelete")

	goazure.Router("/boiler_alarm_count/", &controllers.AlarmController{}, "get:BoilerAlarmCount")
	goazure.Router("/boiler_alarm_list/", &controllers.AlarmController{}, "get:BoilerAlarmList")
	goazure.Router("/boiler_alarm_history_list/", &controllers.AlarmController{}, "get:BoilerAlarmHistoryList")
	goazure.Router("/boiler_alarm_update/", &controllers.AlarmController{}, "post:BoilerAlarmUpdate")
	goazure.Router("/boiler_alarm_detail/", &controllers.AlarmController{}, "post:BoilerAlarmDetail")
	goazure.Router("/boiler_alarm_feedback_list/", &controllers.AlarmController{}, "post:BoilerAlarmFeedbackList")

	goazure.Router("/dialogue_list/", &controllers.DialogueController{}, "get:DialogueList")
	goazure.Router("/dialogue_comment_update/", &controllers.DialogueController{}, "post:CommentUpdate")
	goazure.Router("/dialogue_delete/", &controllers.DialogueController{}, "post:DialogueDelete")

	goazure.Router("/location_list/", &controllers.LocationController{}, "get:LocationList")

	goazure.Router("/terminal_origin_message_list/", &controllers.DeveloperController{}, "get:TerminalOriginMessageList")
	//下发
	goazure.Router("/term_function_code_list",&controllers.IssuedController{},"get:FunctionCodeList")
	goazure.Router("/term_byte_list",&controllers.IssuedController{},"get:ByteCodeList")
	goazure.Router("/baud_rate_list",&controllers.IssuedController{},"get:BaudRateList")
	goazure.Router("/correspond_type_list",&controllers.IssuedController{},"get:CorrespondTypeList")
	goazure.Router("/date_bit_list",&controllers.IssuedController{},"get:DataBitList")
	goazure.Router("/heartbeat_packet_list",&controllers.IssuedController{},"get:HeartbeatPacketList")
	goazure.Router("/parity_bit",&controllers.IssuedController{},"get:ParityBitList")
	goazure.Router("/slave_address_list",&controllers.IssuedController{},"get:SlaveAddressList")
	goazure.Router("/stop_bit_list",&controllers.IssuedController{},"get:StopBitList")

	//重启
	goazure.Router("/terminal_restart",&controllers.IssuedController{},"post:TerminalRestart")

	//bin文件上传
	goazure.Router("/bin_upload",&controllers.IssuedController{},"post:BinUpload")
	//获取bin文件路径
	goazure.Router("/bin_list",&controllers.IssuedController{},"get:BinFileList")
	//bin升级配置
	goazure.Router("/upgrade_configuration",&controllers.IssuedController{},"post:UpgradeConfiguration")

	//下方配置报文
	goazure.Router("/issued_config",&controllers.IssuedController{},"post:IssuedConfig")

	goazure.Router("/terminal_issued_list",&controllers.TerminalController{},"get:TerminalIssuedList")

	//获取通信参数
	goazure.Router("/issued_communication",&controllers.ParameterController{},"post:IssuedCommunication")

	//模板另存为
	goazure.Router("/issued_template",&controllers.TemplateController{},"post:IssuedTemplate")

	//锅炉重启
	goazure.Router("/issued_boiler",&controllers.IssuedController{},"post:IssuedBoiler")

	//下发测试按钮


	goazure.SetStaticPath("/assets", "static/assets/")
	goazure.SetStaticPath("/js", "static/js/")
	goazure.SetStaticPath("/css", "static/css/")
	goazure.SetStaticPath("/img", "static/img/")
	goazure.SetStaticPath("/images", "static/images/")
	goazure.SetStaticPath("/bower_components", "static/bower_components/")
	goazure.SetStaticPath("/res", "static/res/")
	goazure.SetStaticPath("/node_modules", "static/node_modules/")
	goazure.SetStaticPath("/views", "views/views/")
	goazure.SetStaticPath("/tpl", "views/tpl/")
	goazure.SetStaticPath("/directives", "views/directives/")

	goazure.SetStaticPath("/upload", "static/images/upload/")

	initWeixinRoutes()
}

func initWeixinRoutes() {
	goazure.Router("/user_login_weixin/", &controllers.UserThirdController{}, "get:UserLoginWeixinWeb")
	goazure.Router("/user_login_weixin/callback/?:code:state", &controllers.UserThirdController{}, "get:UserLoginWeixinWebCallback")
	goazure.Router("/user_login_bind_third/", &controllers.UserThirdController{}, "post:UserLoginBindThird")
	goazure.Router("/user_register_bind_third/", &controllers.UserThirdController{}, "post:UserRegisterBindThird")
	goazure.Router("/user_bind_weixin/", &controllers.UserThirdController{}, "get:UserBindWeixin")
	goazure.Router("/user_bind_weixin/callback/?:code:state", &controllers.UserThirdController{}, "get:UserBindWeixinCallback")
	goazure.Router("/user_unbind_weixin/", &controllers.UserThirdController{}, "post:UserUnbindWeixin")

	goazure.Router("/user_login_weixin_mini/", &controllers.UserThirdController{}, "post:UserLoginWeixinMini")
	goazure.Router("/user_bind_weixin_mini/", &controllers.UserThirdController{}, "post:UserBindWeixinMini")
	goazure.Router("/user_unbind_weixin_mini/", &controllers.UserThirdController{}, "get:UserUnbindWeixin")

	goazure.Router("/boiler_list_weixin/", &controllers.BoilerController{}, "get:BoilerListWeixin")
	goazure.Router("/boiler_runtime_daily_weixin/", &controllers.RuntimeController{}, "post:BoilerRuntimeDaily")

	goazure.Router("/location_list_weixin/", &controllers.LocationController{}, "get:LocationListWeixin")
	goazure.Router("/fuel_list_weixin/", &controllers.FuelController{}, "get:FuelListWeixin")

	goazure.Router("/fuel_record_list/", &controllers.FuelController{}, "get:FuelRecordList")
	goazure.Router("/fuel_record_update/", &controllers.FuelController{}, "post:FuelRecordUpdate")
	goazure.Router("/fuel_record_delete/", &controllers.FuelController{}, "post:FuelRecordDelete")


	//goazure.Router("/wechat-server", &controllers.WechatController{})
	goazure.Router("/wechat-server", &controllers.WechatController{}, "*:WXCallbackHandler")
	goazure.Router("/wechat-subscribe", &controllers.WechatController{})
	goazure.Router("/wechat-mini-server", &controllers.WechatController{})
}
