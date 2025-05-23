-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

-- 创建园区按月用电量统计视图
CREATE VIEW v_eco_park_power_consumption_monthly AS
SELECT 
    p.id AS park_id,
    p.park_name,
    date_trunc('month', f.time) AS month,
    SUM(CASE WHEN f.type = 1 THEN f.power_consumption ELSE 0 END) AS lighting_socket_power,
    SUM(CASE WHEN f.type = 2 THEN f.power_consumption ELSE 0 END) AS air_conditioning_power,
    SUM(f.power_consumption) AS total_power_consumption
FROM 
    f_eco_park_1m f
    JOIN o_eco_park p ON f.park_id = p.id
GROUP BY 
    p.id, p.park_name, date_trunc('month', f.time)
ORDER BY 
    p.park_name, date_trunc('month', f.time);

COMMENT ON VIEW v_eco_park_power_consumption_monthly IS '园区按月用电量统计视图';
COMMENT ON COLUMN v_eco_park_power_consumption_monthly.park_id IS '园区ID';
COMMENT ON COLUMN v_eco_park_power_consumption_monthly.park_name IS '园区名称';
COMMENT ON COLUMN v_eco_park_power_consumption_monthly.month IS '统计月份';
COMMENT ON COLUMN v_eco_park_power_consumption_monthly.lighting_socket_power IS '照明插座用电量(kWh)';
COMMENT ON COLUMN v_eco_park_power_consumption_monthly.air_conditioning_power IS '空调动力用电量(kWh)';
COMMENT ON COLUMN v_eco_park_power_consumption_monthly.total_power_consumption IS '总用电量(kWh)';

-- 创建楼栋按月用电量统计视图
CREATE VIEW v_eco_building_power_consumption_monthly AS
SELECT 
    b.id AS building_id,
    b.building_name,
    p.id AS park_id,
    p.park_name,
    date_trunc('month', f.time) AS month,
    SUM(CASE WHEN f.type = 1 THEN f.power_consumption ELSE 0 END) AS lighting_socket_power,
    SUM(CASE WHEN f.type = 2 THEN f.power_consumption ELSE 0 END) AS air_conditioning_power,
    SUM(f.power_consumption) AS total_power_consumption
FROM 
    f_eco_building_1m f
    JOIN o_eco_building b ON f.building_id = b.id
    JOIN o_eco_park p ON f.park_id = p.id
GROUP BY 
    b.id, b.building_name, p.id, p.park_name, date_trunc('month', f.time)
ORDER BY 
    p.park_name, b.building_name, date_trunc('month', f.time);

COMMENT ON VIEW v_eco_building_power_consumption_monthly IS '楼栋按月用电量统计视图';
COMMENT ON COLUMN v_eco_building_power_consumption_monthly.building_id IS '楼栋ID';
COMMENT ON COLUMN v_eco_building_power_consumption_monthly.building_name IS '楼栋名称';
COMMENT ON COLUMN v_eco_building_power_consumption_monthly.park_id IS '园区ID';
COMMENT ON COLUMN v_eco_building_power_consumption_monthly.park_name IS '园区名称';
COMMENT ON COLUMN v_eco_building_power_consumption_monthly.month IS '统计月份';
COMMENT ON COLUMN v_eco_building_power_consumption_monthly.lighting_socket_power IS '照明插座用电量(kWh)';
COMMENT ON COLUMN v_eco_building_power_consumption_monthly.air_conditioning_power IS '空调动力用电量(kWh)';
COMMENT ON COLUMN v_eco_building_power_consumption_monthly.total_power_consumption IS '总用电量(kWh)';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP VIEW IF EXISTS v_eco_building_power_consumption_monthly;
DROP VIEW IF EXISTS v_eco_park_power_consumption_monthly;

-- +goose StatementEnd
