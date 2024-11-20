
-- 刷新楼层用电量统计视图
CALL refresh_continuous_aggregate('f_eco_floor_1h', NULL, NULL);
CALL refresh_continuous_aggregate('f_eco_floor_1d', NULL, NULL);
CALL refresh_continuous_aggregate('f_eco_floor_1m', NULL, NULL);
CALL refresh_continuous_aggregate('f_eco_floor_1y', NULL, NULL);

-- 刷新楼栋用电量统计视图
CALL refresh_continuous_aggregate('f_eco_building_1h', NULL, NULL);
CALL refresh_continuous_aggregate('f_eco_building_1d', NULL, NULL);
CALL refresh_continuous_aggregate('f_eco_building_1m', NULL, NULL);
CALL refresh_continuous_aggregate('f_eco_building_1y', NULL, NULL);

-- 刷新园区用电量统计视图
CALL refresh_continuous_aggregate('f_eco_park_1h', NULL, NULL);
CALL refresh_continuous_aggregate('f_eco_park_1d', NULL, NULL);
CALL refresh_continuous_aggregate('f_eco_park_1m', NULL, NULL);
CALL refresh_continuous_aggregate('f_eco_park_1y', NULL, NULL);

-- 刷新园区水表用电量统计视图
CALL refresh_continuous_aggregate('f_eco_park_water_1h', NULL, NULL);
CALL refresh_continuous_aggregate('f_eco_park_water_1d', NULL, NULL);
CALL refresh_continuous_aggregate('f_eco_park_water_1m', NULL, NULL);
CALL refresh_continuous_aggregate('f_eco_park_water_1y', NULL, NULL);
