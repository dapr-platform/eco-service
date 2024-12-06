dp-cli gen --connstr "postgresql://things:things2024@ali4:40207/thingsdb?sslmode=disable" \
--tables=v_eco_building_info_1d,v_eco_building_info_1m,v_eco_building_info_1y --model_naming "{{ toUpperCamelCase ( replace . \"v_\" \"\") }}"  \
--file_naming "{{ toLowerCamelCase ( replace . \"v_\" \"\") }}" \
--module eco-service

