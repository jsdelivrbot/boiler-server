SELECT 	COUNT(*) AS `count`,
        IF(`fuel_type_id` = 1 OR `fuel_type_id` = 4, 'coal', 'gas') AS `fuel_cate`,
        CASE 	WHEN `evaporte` BETWEEN 0 AND 1 THEN 'D≤1'
        WHEN `evaporte` BETWEEN 1 AND 2 THEN '1＜D≤2'
        WHEN `evaporte` BETWEEN 2 AND 8 THEN '2＜D≤8'
        WHEN `evaporte` BETWEEN 8 AND 20 THEN '8＜D≤20'
        ELSE 'D＞20'
        END AS 	`evaporate_cate`,
        CASE 	WHEN `evaporte` BETWEEN 0 AND 1 THEN 'c0'
        WHEN `evaporte` BETWEEN 1 AND 2 THEN 'c1'
        WHEN `evaporte` BETWEEN 2 AND 8 THEN 'c2'
        WHEN `evaporte` BETWEEN 8 AND 20 THEN 'c3'
        ELSE 'c4'
        END AS 	`evaporate_id`,
        IF(`heats` >= (SELECT CASE
                              WHEN `evaporte` BETWEEN 0 AND 1 THEN 61
                              WHEN `evaporte` BETWEEN 1 AND 2 THEN 69
                              WHEN `evaporte` BETWEEN 2 AND 8 THEN 71
                              WHEN `evaporte` BETWEEN 8 AND 20 THEN 72
                              ELSE 72 END),
           'success', 'failed') AS `rank`
FROM
  (SELECT `boiler`.`uid`, `boiler`.`name`, `boiler`.`evaporating_capacity` AS `evaporte`, `fuel`.`type_id` AS `fuel_type_id`, AVG(`heat`.`value`) AS `heats`, YEARWEEK(`heat`.`created_date`) AS `week`
   FROM	`boiler`, `fuel`, `boiler_runtime_cache_heat` AS `heat`
   WHERE	`boiler`.`fuel_id` = `fuel`.`uid` AND `boiler`.`uid` = `heat`.`boiler_id` AND YEARWEEK(`heat`.`created_date`) = YEARWEEK(CURDATE())
   GROUP BY `heat`.`boiler_id`) AS `bFuelEva`
WHERE	`fuel_type_id` = 1 OR `fuel_type_id` = 4
GROUP BY `fuel_cate`, `evaporate_cate`, `rank`
ORDER BY `fuel_cate`, `evaporate_cate`, `rank`;