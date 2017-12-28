BEGIN
	IF NEW.`parameter_id` = 1001 THEN
		INSERT INTO `boiler_runtime_cache_steam_temperature`
		SET `runtime_id` = NEW.`id`,
				`boiler_id` = NEW.`boiler_id`,
				`parameter_id` = NEW.`parameter_id`,
				`alarm_id` = NEW.`alarm_id`,

				`created_date` = NEW.`created_date`,
				`created_by_id` = NEW.`created_by_id`,
				`updated_date` = NEW.`updated_date`,
				`updated_by_id` = NEW.`updated_by_id`,
				`is_deleted` = NEW.`is_deleted`,
				`is_demo` = NEW.`is_demo`,

				`name_en` = NEW.`name_en`,
				`remark` = NEW.`remark`,
				`name` = (SELECT `boiler`.`name`
									FROM `boiler`
									WHERE `boiler`.`uid` = NEW.`boiler_id`),

				`value` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
									FROM `runtime_parameter` AS `param`
									WHERE `param`.id = NEW.parameter_id),
				`parameter_name` = (SELECT `param`.`name`
									FROM `runtime_parameter` AS `param`
									WHERE `param`.id = NEW.parameter_id),
				`unit` = (SELECT `param`.`unit`
									FROM `runtime_parameter` AS `param`
									WHERE `param`.id = NEW.parameter_id),

				`alarm_level` = (SELECT `alarm`.`alarm_level`
									FROM `boiler_alarm` AS `alarm`
									WHERE `alarm`.uid = NEW.alarm_id),
				`alarm_description` = (SELECT `alarm`.`description`
									FROM `boiler_alarm` AS `alarm`
									WHERE `alarm`.uid = NEW.alarm_id)
				;

		UPDATE `boiler_runtime_cache_steam_temperature`
		SET `alarm_level` = 0
		WHERE `alarm_level` IS NULL;

		UPDATE `boiler_runtime_cache_steam_temperature`
		SET `alarm_description` = ''
		WHERE `alarm_description` IS NULL;

		UPDATE `boiler_runtime_cache_history` AS `history`
		SET `p1001` = (SELECT `cache`.`value`
				FROM `boiler_runtime_cache_steam_temperature` AS `cache`
				WHERE `cache`.`runtime_id` = NEW.`id`),
			`a1001` = (SELECT `cache`.`alarm_level`
				FROM `boiler_runtime_cache_steam_temperature` AS `cache`
				WHERE `cache`.`runtime_id` = NEW.`id`),
		WHERE `history`.`boiler_id` = NEW.`boiler_id`
			AND `history`.`created_date` > (NOW() - INTERVAL 5 MINUTE)
			AND `history`.`created_date` <= NOW();
	END IF;

	IF NEW.`parameter_id` = 1002 THEN
		INSERT INTO `boiler_runtime_cache_steam_pressure`
		SET `runtime_id` = NEW.id,
				`boiler_id` = NEW.boiler_id,
				`parameter_id` = NEW.parameter_id,
				`alarm_id` = NEW.alarm_id,

			created_date = NEW.created_date,
			created_by_id = NEW.created_by_id,
			updated_date = NEW.updated_date,
			updated_by_id = NEW.updated_by_id,
			is_deleted = NEW.is_deleted,
			is_demo = NEW.is_demo,

			`name_en` = NEW.`name_en`,
			`remark` = NEW.`remark`,
			`name` = (SELECT `boiler`.`name`
				FROM `boiler`
				WHERE `boiler`.uid = NEW.boiler_id),

			`value` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`parameter_name` = (SELECT `param`.`name`
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`unit` = (SELECT `param`.`unit`
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),

			`alarm_level` = (SELECT `alarm`.`alarm_level`
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id),
			`alarm_description` = (SELECT `alarm`.`description`
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id)
			;

		UPDATE boiler_runtime_cache_steam_pressure
		SET `alarm_level` = 0
		WHERE `alarm_level` IS NULL;

		UPDATE boiler_runtime_cache_steam_pressure
		SET `alarm_description` = ''
		WHERE `alarm_description` IS NULL;

		UPDATE `boiler_runtime_cache_history` AS `history`
		SET `p1002` = (SELECT `cache`.`value`
									FROM `boiler_runtime_cache_steam_pressure` AS `cache`
									WHERE `cache`.`runtime_id` = NEW.`id`),
				`a1002` = (SELECT `cache`.`alarm_level`
									FROM `boiler_runtime_cache_steam_pressure` AS `cache`
									WHERE `cache`.`runtime_id` = NEW.`id`),
		WHERE `history`.`boiler_id` = NEW.`boiler_id`
			AND `history`.`created_date` > (NOW() - INTERVAL 5 MINUTE)
			AND `history`.`created_date` <= NOW();
	END IF;

	IF NEW.parameter_id = 1003 THEN
		INSERT INTO boiler_runtime_cache_flow
		SET runtime_id = NEW.id,
			boiler_id = NEW.boiler_id,
			parameter_id = NEW.parameter_id,
			alarm_id = NEW.alarm_id,

			created_date = NEW.created_date,
			created_by_id = NEW.created_by_id,
			updated_date = NEW.updated_date,
			updated_by_id = NEW.updated_by_id,
			is_deleted = NEW.is_deleted,
			is_demo = NEW.is_demo,

			`name_en` = NEW.`name_en`,
			`remark` = NEW.`remark`,
			`name` = (SELECT `boiler`.`name`
				FROM `boiler`
				WHERE `boiler`.uid = NEW.boiler_id),

			`value` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`parameter_name` = (SELECT `param`.`name`
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`unit` = (SELECT `param`.`unit`
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),

			`alarm_level` = (SELECT `alarm`.`alarm_level`
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id),
			`alarm_description` = (SELECT `alarm`.`description`
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id)
			;

		UPDATE `boiler_runtime_cache_flow`
		SET `alarm_level` = 0
		WHERE `alarm_level` IS NULL;

		UPDATE `boiler_runtime_cache_flow`
		SET `alarm_description` = ''
		WHERE `alarm_description` IS NULL;

		UPDATE `boiler_runtime_cache_history` AS `history`
		SET `p1003` = (SELECT `cache`.`value`
									FROM `boiler_runtime_cache_flow` AS `cache`
									WHERE `cache`.`runtime_id` = NEW.`id`),
				`a1003` = (SELECT `cache`.`alarm_level`
									FROM `boiler_runtime_cache_flow` AS `cache`
									WHERE `cache`.`runtime_id` = NEW.`id`),
		WHERE `history`.`boiler_id` = NEW.`boiler_id`
			AND `history`.`created_date` > (NOW() - INTERVAL 5 MINUTE)
			AND `history`.`created_date` <= NOW();
	END IF;


	IF NEW.parameter_id = 1005 OR NEW.parameter_id = 1006 THEN
		INSERT INTO boiler_runtime_cache_water_temperature
		SET runtime_id = NEW.id,
			boiler_id = NEW.boiler_id,
			parameter_id = NEW.parameter_id,
			alarm_id = NEW.alarm_id,

			created_date = NEW.created_date,
			created_by_id = NEW.created_by_id,
			updated_date = NEW.updated_date,
			updated_by_id = NEW.updated_by_id,
			is_deleted = NEW.is_deleted,
			is_demo = NEW.is_demo,

			`name_en` = NEW.`name_en`,
			`remark` = NEW.`remark`,
			`name` = (SELECT `boiler`.`name`
				FROM `boiler`
				WHERE `boiler`.uid = NEW.boiler_id),

			`value` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`parameter_name` = (SELECT `param`.`name`
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`unit` = (SELECT `param`.`unit`
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),

			`alarm_level` = (SELECT `alarm`.`alarm_level`
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id),
			`alarm_description` = (SELECT `alarm`.`description`
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id)
			;

		UPDATE boiler_runtime_cache_water_temperature
		SET `alarm_level` = 0
		WHERE `alarm_level` IS NULL;

		UPDATE boiler_runtime_cache_water_temperature
		SET `alarm_description` = ''
		WHERE `alarm_description` IS NULL;

		IF NEW.parameter_id = 1005 THEN
			UPDATE `boiler_runtime_cache_history` AS `history`
			SET `p1005` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
					FROM `runtime_parameter` AS `param`
					WHERE `param`.id = NEW.parameter_id)
			WHERE `history`.`boiler_id` = NEW.`boiler_id`
				AND `history`.`created_date` > (NOW() - INTERVAL 5 MINUTE)
				AND `history`.`created_date` <= NOW();
		ELSEIF NEW.parameter_id = 1006 THEN
			UPDATE `boiler_runtime_cache_history` AS `history`
			SET `p1006` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
					FROM `runtime_parameter` AS `param`
					WHERE `param`.id = NEW.parameter_id)
			WHERE `history`.`boiler_id` = NEW.`boiler_id`
				AND `history`.`created_date` > (NOW() - INTERVAL 5 MINUTE)
				AND `history`.`created_date` <= NOW();
		END IF;
	END IF;

	IF NEW.parameter_id = 1014 OR NEW.parameter_id = 1015 THEN
		INSERT INTO boiler_runtime_cache_smoke_temperature
		SET runtime_id = NEW.id,
			boiler_id = NEW.boiler_id,
			parameter_id = NEW.parameter_id,
			alarm_id = NEW.alarm_id,

			created_date = NEW.created_date,
			created_by_id = NEW.created_by_id,
			updated_date = NEW.updated_date,
			updated_by_id = NEW.updated_by_id,
			is_deleted = NEW.is_deleted,
			is_demo = NEW.is_demo,

			`name_en` = NEW.`name_en`,
			`remark` = NEW.`remark`,
			`name` = (SELECT `boiler`.`name`
				FROM `boiler`
				WHERE `boiler`.uid = NEW.boiler_id),

			`value` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`parameter_name` = (SELECT `param`.`name`
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`unit` = (SELECT `param`.`unit`
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),

			`alarm_level` = (SELECT `alarm`.`alarm_level`
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id),
			`alarm_description` = (SELECT `alarm`.`description`
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id)
			;

		UPDATE boiler_runtime_cache_smoke_temperature
		SET `alarm_level` = 0
		WHERE `alarm_level` IS NULL;

		UPDATE boiler_runtime_cache_smoke_temperature
		SET `alarm_description` = ''
		WHERE `alarm_description` IS NULL;

		IF NEW.parameter_id = 1014 THEN
			UPDATE `boiler_runtime_cache_history` AS `history`
			SET `p1014` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
					FROM `runtime_parameter` AS `param`
					WHERE `param`.id = NEW.parameter_id)
			WHERE `history`.`boiler_id` = NEW.`boiler_id`
				AND `history`.`created_date` > (NOW() - INTERVAL 5 MINUTE)
				AND `history`.`created_date` <= NOW();
		ELSEIF NEW.parameter_id = 1015 THEN
			UPDATE `boiler_runtime_cache_history` AS `history`
			SET `p1015` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
					FROM `runtime_parameter` AS `param`
					WHERE `param`.id = NEW.parameter_id)
			WHERE `history`.`boiler_id` = NEW.`boiler_id`
				AND `history`.`created_date` > (NOW() - INTERVAL 5 MINUTE)
				AND `history`.`created_date` <= NOW();
		END IF;
	END IF;

	IF NEW.parameter_id = 1016 OR NEW.parameter_id = 1017 OR NEW.parameter_id = 1018 OR NEW.parameter_id = 1019 THEN
		INSERT INTO boiler_runtime_cache_smoke_component
		SET runtime_id = NEW.id,
			boiler_id = NEW.boiler_id,
			parameter_id = NEW.parameter_id,
			alarm_id = NEW.alarm_id,

			created_date = NEW.created_date,
			created_by_id = NEW.created_by_id,
			updated_date = NEW.updated_date,
			updated_by_id = NEW.updated_by_id,
			is_deleted = NEW.is_deleted,
			is_demo = NEW.is_demo,

			`name_en` = NEW.`name_en`,
			`remark` = NEW.`remark`,
			`name` = (SELECT `boiler`.`name`
				FROM `boiler`
				WHERE `boiler`.uid = NEW.boiler_id),

			`value` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`parameter_name` = (SELECT `param`.`name`
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`unit` = (SELECT `param`.`unit`
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),

			`alarm_level` = (SELECT `alarm`.`alarm_level`
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id),
			`alarm_description` = (SELECT `alarm`.`description`
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id)
			;

		UPDATE boiler_runtime_cache_smoke_component
		SET `alarm_level` = 0
		WHERE `alarm_level` IS NULL;

		UPDATE boiler_runtime_cache_smoke_component
		SET `alarm_description` = ''
		WHERE `alarm_description` IS NULL;

		IF NEW.parameter_id = 1016 THEN
			UPDATE `boiler_runtime_cache_history` AS `history`
			SET `p1016` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
					FROM `runtime_parameter` AS `param`
					WHERE `param`.id = NEW.parameter_id)
			WHERE `history`.`boiler_id` = NEW.`boiler_id`
				AND `history`.`created_date` > (NOW() - INTERVAL 5 MINUTE)
				AND `history`.`created_date` <= NOW();
		ELSEIF NEW.parameter_id = 1017 THEN
			UPDATE `boiler_runtime_cache_history` AS `history`
			SET `p1017` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
					FROM `runtime_parameter` AS `param`
					WHERE `param`.id = NEW.parameter_id)
			WHERE `history`.`boiler_id` = NEW.`boiler_id`
				AND `history`.`created_date` > (NOW() - INTERVAL 5 MINUTE)
				AND `history`.`created_date` <= NOW();
		ELSEIF NEW.parameter_id = 1018 THEN
			UPDATE `boiler_runtime_cache_history` AS `history`
			SET `p1018` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
					FROM `runtime_parameter` AS `param`
					WHERE `param`.id = NEW.parameter_id)
			WHERE `history`.`boiler_id` = NEW.`boiler_id`
				AND `history`.`created_date` > (NOW() - INTERVAL 5 MINUTE)
				AND `history`.`created_date` <= NOW();
		ELSEIF NEW.parameter_id = 1019 THEN
			UPDATE `boiler_runtime_cache_history` AS `history`
			SET `p1019` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
					FROM `runtime_parameter` AS `param`
					WHERE `param`.id = NEW.parameter_id)
			WHERE `history`.`boiler_id` = NEW.`boiler_id`
				AND `history`.`created_date` > (NOW() - INTERVAL 5 MINUTE)
				AND `history`.`created_date` <= NOW();
		END IF;
	END IF;


	IF NEW.parameter_id = 1021 OR NEW.parameter_id = 1022 THEN
		INSERT INTO boiler_runtime_cache_environment_temperature
		SET runtime_id = NEW.id,
			boiler_id = NEW.boiler_id,
			parameter_id = NEW.parameter_id,
			alarm_id = NEW.alarm_id,

			created_date = NEW.created_date,
			created_by_id = NEW.created_by_id,
			updated_date = NEW.updated_date,
			updated_by_id = NEW.updated_by_id,
			is_deleted = NEW.is_deleted,
			is_demo = NEW.is_demo,

			`name_en` = NEW.`name_en`,
			`remark` = NEW.`remark`,
			`name` = (SELECT `boiler`.`name`
				FROM `boiler`
				WHERE `boiler`.uid = NEW.boiler_id),

			`value` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`parameter_name` = (SELECT `param`.`name`
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`unit` = (SELECT `param`.`unit`
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),

			`alarm_level` = (SELECT `alarm`.`alarm_level`
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id),
			`alarm_description` = (SELECT `alarm`.`description`
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id)
			;

		UPDATE `boiler_runtime_cache_environment_temperature`
		SET `alarm_level` = 0
		WHERE `alarm_level` IS NULL;

		UPDATE `boiler_runtime_cache_environment_temperature`
		SET `alarm_description` = ''
		WHERE `alarm_description` IS NULL;

		IF NEW.parameter_id = 1021 THEN
			UPDATE `boiler_runtime_cache_history` AS `history`
			SET `p1021` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
					FROM `runtime_parameter` AS `param`
					WHERE `param`.id = NEW.parameter_id)
			WHERE `history`.`boiler_id` = NEW.`boiler_id`
				AND `history`.`created_date` > (NOW() - INTERVAL 5 MINUTE)
				AND `history`.`created_date` <= NOW();
		ELSEIF NEW.parameter_id = 1022 THEN
			UPDATE `boiler_runtime_cache_history` AS `history`
			SET `p1022` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
					FROM `runtime_parameter` AS `param`
					WHERE `param`.id = NEW.parameter_id)
			WHERE `history`.`boiler_id` = NEW.`boiler_id`
				AND `history`.`created_date` > (NOW() - INTERVAL 5 MINUTE)
				AND `history`.`created_date` <= NOW();
		END IF;
	END IF;

	IF NEW.parameter_id = 1201 THEN
		INSERT INTO boiler_runtime_cache_heat
		SET runtime_id = NEW.id,
			boiler_id = NEW.boiler_id,
			parameter_id = NEW.parameter_id,
			alarm_id = NEW.alarm_id,

			created_date = NEW.created_date,
			created_by_id = NEW.created_by_id,
			updated_date = NEW.updated_date,
			updated_by_id = NEW.updated_by_id,
			is_deleted = NEW.is_deleted,
			is_demo = NEW.is_demo,

			`name_en` = NEW.`name_en`,
			`remark` = NEW.`remark`,
			`name` = (SELECT `boiler`.`name`
				FROM `boiler`
				WHERE `boiler`.uid = NEW.boiler_id),

			`value` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`parameter_name` = (SELECT `param`.`name`
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`unit` = (SELECT `param`.`unit`
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),

			`alarm_level` = (SELECT `alarm`.`alarm_level`
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id),
			`alarm_description` = (SELECT `alarm`.`description`
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id)
			;

		UPDATE boiler_runtime_cache_heat
		SET `alarm_level` = 0
		WHERE `alarm_level` IS NULL;

		UPDATE boiler_runtime_cache_heat
		SET `alarm_description` = ''
		WHERE `alarm_description` IS NULL;

		UPDATE `boiler_runtime_cache_history` AS `history`
		SET `p1201` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id)
		WHERE `history`.`boiler_id` = NEW.`boiler_id`
			AND `history`.`created_date` > (NOW() - INTERVAL 5 MINUTE)
			AND `history`.`created_date` <= NOW();
	END IF;

	IF NEW.parameter_id = 1202 THEN
		INSERT INTO boiler_runtime_cache_excess_air
		SET runtime_id = NEW.id,
			boiler_id = NEW.boiler_id,
			parameter_id = NEW.parameter_id,
			alarm_id = NEW.alarm_id,

			created_date = NEW.created_date,
			created_by_id = NEW.created_by_id,
			updated_date = NEW.updated_date,
			updated_by_id = NEW.updated_by_id,
			is_deleted = NEW.is_deleted,
			is_demo = NEW.is_demo,

			`name_en` = NEW.`name_en`,
			`remark` = NEW.`remark`,
			`name` = (SELECT `boiler`.`name`
				FROM `boiler`
				WHERE `boiler`.uid = NEW.boiler_id),

			`value` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`parameter_name` = (SELECT `param`.`name`
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`unit` = (SELECT `param`.`unit`
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),

			`alarm_level` = (SELECT `alarm`.`alarm_level`
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id),
			`alarm_description` = (SELECT `alarm`.`description`
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id)
			;

		UPDATE boiler_runtime_cache_excess_air
		SET `alarm_level` = 0
		WHERE `alarm_level` IS NULL;

		UPDATE boiler_runtime_cache_excess_air
		SET `alarm_description` = ''
		WHERE `alarm_description` IS NULL;

		UPDATE `boiler_runtime_cache_history` AS `history`
		SET `p1202` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id)
		WHERE `history`.`boiler_id` = NEW.`boiler_id`
			AND `history`.`created_date` > (NOW() - INTERVAL 5 MINUTE)
			AND `history`.`created_date` <= NOW();
	END IF;
END