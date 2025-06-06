dp-cli gen --connstr "postgresql://things:things2024@ali4:40207/thingsdb?sslmode=disable" \
--tables=f_eco_gateway_1h,f_eco_gateway_1d,f_eco_floor_1d,f_eco_building_1h,f_eco_building_1d,f_eco_floor_1h,f_eco_gateway_1m,f_eco_floor_1m,f_eco_building_1m,f_eco_gateway_1y,f_eco_floor_1y,f_eco_building_1y,f_eco_park_1h,f_eco_park_1d,f_eco_park_1m,f_eco_park_1y,f_eco_water_meter_1h,f_eco_park_water_1h,f_eco_park_water_1d,f_eco_park_water_1m,f_eco_park_water_1y \
--model_naming "{{ toUpperCamelCase ( replace . \"f_\" \"\") }}"  \
--file_naming "{{ toLowerCamelCase ( replace . \"f_\" \"\") }}" \
--module eco-service

