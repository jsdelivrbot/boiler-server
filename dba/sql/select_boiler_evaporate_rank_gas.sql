SELECT 	COUNT(*) AS `count`,
        IF(`fuel_type_id` = 1 OR `fuel_type_id` = 4, 'coal', 'gas') AS `fuel_cate`,
        IF(`evaporte` <= 2, 'D≤2', 'D＞2') AS `evaporate_cate`,
        IF(`evaporte` <= 2, 'g0', 'g1') AS `evaporate_id`,
        IF(`heats` >= IF(`evaporte` <= 2, 79, 81),
           'success', 'failed') AS `rank`
FROM
  (SELECT `boiler`.`uid`, `boiler`.`name`, `boiler`.`evaporating_capacity` AS `evaporte`, `fuel`.`type_id` AS `fuel_type_id`, AVG(`heat`.`value`) AS `heats`, YEARWEEK(`heat`.`created_date`) AS `week`
   FROM	`boiler`, `fuel`, `boiler_runtime_cache_heat` AS `heat`
   WHERE	`boiler`.`fuel_id` = `fuel`.`uid` AND `boiler`.`uid` = `heat`.`boiler_id` AND YEARWEEK(`heat`.`created_date`) = YEARWEEK(CURDATE())
   GROUP BY `heat`.`boiler_id`) AS `bFuelEva`
WHERE	`fuel_type_id` <> 1 AND `fuel_type_id` <> 4
GROUP BY `fuel_cate`, `evaporate_cate`, `rank`
ORDER BY `fuel_cate`, `evaporate_cate`, `rank`;