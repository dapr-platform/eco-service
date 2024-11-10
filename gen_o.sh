dp-cli gen --connstr "postgresql://things:things2024@ali4:37432/thingsdb?sslmode=disable" \
--tables=o_eco_park,o_eco_gateway,o_eco_building,o_eco_floor --model_naming "{{ toUpperCamelCase ( replace . \"o_\" \"\") }}"  \
--file_naming "{{ toLowerCamelCase ( replace . \"o_\" \"\") }}" \
--module eco-service --api true

