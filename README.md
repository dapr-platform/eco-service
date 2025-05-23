查询楼栋和园区用电量月统计sql
园区
select park_name ,month,lighting_socket_power as 照明开关,air_conditioning_power as 空调动力, total_power_consumption as 总用电量 
from v_eco_park_power_consumption_monthly
where month>'2024-11-01'
order by month,park_name 

楼栋
select building_name ,month,lighting_socket_power as 照明开关,air_conditioning_power as 空调动力, total_power_consumption as 总用电量 
from v_eco_building_power_consumption_monthly
where month>'2024-11-01'
order by month,building_name 