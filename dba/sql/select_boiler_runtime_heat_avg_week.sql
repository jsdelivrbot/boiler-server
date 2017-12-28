SELECT 	`boiler`.`uid` AS `uid`, AVG(`heat`.`value`) AS `heat`, YEARWEEK(`heat`.`created_date`) AS `week`
FROM	`boiler`, `fuel`, `boiler_runtime_cache_heat` AS `heat`
WHERE	`boiler`.`fuel_id` = `fuel`.`uid` AND `boiler`.`uid` = `heat`.`boiler_id` AND YEARWEEK(`heat`.`created_date`) = YEARWEEK(CURDATE())
GROUP BY `heat`.`boiler_id`;