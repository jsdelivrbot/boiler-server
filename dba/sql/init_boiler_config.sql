INSERT INTO `boiler_config`
(`uid`, `name`, `created_date`, `updated_date`, `boiler_id`, `is_generate_data`)
SELECT UUID(), `boiler`.`name`, NOW(), NOW(), `boiler`.`uid`, 1
FROM `boiler`
WHERE `boiler`.`uid` NOT IN (SELECT `boiler_id` FROM `boiler_config`)