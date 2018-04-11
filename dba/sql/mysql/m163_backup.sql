insert into `boiler_main`.`boiler_m163` ( `Analog2_channel`, `Reserved4_filler`, `C5_calculate_parm`, `Temper6_channel`, `Boiler_term_id`, `Analog3_channel`, `C3_calculate_parm`, `Switch_in_1_16_channel`, `C12_calculate_parm`, `Temper7_channel`, `Analog4_channel`, `C1_calculate_parm`, `C8_calculate_parm`, `Temper8_channel`, `Analog10_channel`, `Temper10_channel`, `Analog5_channel`, `Boiler_data_fmt_ver`, `Temper1_channel`, `Analog11_channel`, `Boiler_sn`, `C6_calculate_parm`, `Temper11_channel`, `Temper9_channel`, `need_reload`, `Analog12_channel`, `Analog6_channel`, `Temper12_channel`, `Temper2_channel`, `C4_calculate_parm`, `Boiler_status_code`, `Analog7_channel`, `Switch_out_1_16_channel`, `Boiler_term_err_code`, `C10_calculate_parm`, `Temper3_channel`, `C2_calculate_parm`, `Term_sys_time`, `uid`, `C9_calculate_parm`, `Analog8_channel`, `Temper4_channel`, `Analog1_channel`, `Switch_in_17_32_channel`, `TS`, `C7_calculate_parm`, `Reserved1_filler`, `Analog9_channel`, `C11_calculate_parm`, `Reserved2_filler`, `Temper5_channel`, `Reserved3_filler`, `Boiler_boiler_id`) values ( '0', '0', '0.0000', '0', '030070', '0', '0.0000', '0', '0.0000', '0', '0', '0.0000', '0.0000', '0', '0', '0', '0', '10', '0', '0', '6390', '0.0000', '0', '0', '0', '0', '0', '0', '0', '0.0000', 'a3', '0', '0', '00', '0.0000', '0', '0.0000', '2017-04-13 10:44:15', '4eb8da99-a402-11e7-addd-7cd30ac4f6d2', '0.0000', '0', '0', '0', '0', '2017-04-13 10:44:31', '0.0000', '0', '0', '0.0000', '0', '0', '0', '02');

INSERT INTO `message_16bit` 
( 
	`uid`, `created_date`, `server_date`, 
	`terminal_code`, `terminal_set_id`, `serial_number`, `message_type`, `version`, 
	`error_code`, `status`, 
	`channel_temperature_1`, `channel_temperature_2`, `channel_temperature_3`, `channel_temperature_4`, `channel_temperature_5`, `channel_temperature6`, `channel_temperature7`, `channel_temperature_8`, `channel_temperature_9`, `channel_temperature_10`, `channel_temperature_11`, `channel_temperature_12`, 
	`channel_analog_1`, `channel_analog_2`, `channel_analog_3`, `channel_analog_4`, `channel_analog_5`, `channel_analog_6`, `channel_analog_7`, `channel_analog_8`, `channel_analog_9`, `channel_analog_10`, `channel_analog_11`, `channel_analog_12`, 
	`channel_calculate_1`, `channel_calculate_2`, `channel_calculate_3`, `channel_calculate_4`, `channel_calculate_5`, `channel_calculate_6`, `channel_calculate_7`, `channel_calculate_8`, `channel_calculate_9`, `channel_calculate_10`, `channel_calculate_11`, `channel_calculate_12`,
	`channel_switch_in_1_16`, `channel_switch_in_17_32`, `channel_switch_out_1_16`,
	`reserved_1`, `reserved_2`, `reserved_3`, `reserved_4`
) 
SELECT	`uid`, NOW(), `TS`,
		`Boiler_term_id`, `Boiler_boiler_id`, `Boiler_sn`, `Boiler_status_code`, `Boiler_data_fmt_ver`, 
		`Boiler_term_err_code`, 0
		`Temper1_channel`, `Temper2_channel`, `Temper3_channel`, `Temper4_channel`, `Temper5_channel`, `Temper6_channel`, `Temper7_channel`, `Temper8_channel`, `Temper9_channel`, `Temper10_channel`, `Temper11_channel`, `Temper12_channel`, 
		`Analog1_channel`, `Analog2_channel`, `Analog3_channel`, `Analog4_channel`, `Analog5_channel`, `Analog6_channel`, `Analog7_channel`, `Analog8_channel`, `Analog9_channel`, `Analog10_channel`, `Analog11_channel`, `Analog12_channel`, 
		`C1_calculate_parm`, `C2_calculate_parm`, `C3_calculate_parm`, `C4_calculate_parm`, `C5_calculate_parm`, `C6_calculate_parm`, `C7_calculate_parm`, `C8_calculate_parm`, `C9_calculate_parm`, `C10_calculate_parm`, `C11_calculate_parm`, `C12_calculate_parm`, 
		`Switch_in_1_16_channel`, `Switch_in_17_32_channel`, `Switch_out_1_16_channel`,
		`Reserved1_filler`, `Reserved2_filler`, `Reserved3_filler`, `Reserved4_filler`
FROM	`boiler_m163`
WHERE	`need_reload` = 0;