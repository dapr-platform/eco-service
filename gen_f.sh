dp-cli gen --connstr "postgresql://things:things2024@ali4:37432/thingsdb?sslmode=disable" \
--tables=f_eco_gateway_1h,f_eco_gateway_1d,f_eco_floor_1d,f_eco_building_1d,f_eco_gateway_1m,f_eco_floor_1m,f_eco_building_1m,f_eco_gateway_1y,f_eco_floor_1y,f_eco_building_1y \
--model_naming "{{ toUpperCamelCase ( replace . \"f_\" \"\") }}"  \
--file_naming "{{ toLowerCamelCase ( replace . \"f_\" \"\") }}" \
--module eco-service --api true

