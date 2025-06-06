-- +goose Up
-- +goose StatementBegin

CREATE TABLE o_eco_park (
    id VARCHAR(32) PRIMARY KEY,
    created_by VARCHAR(32) NOT NULL,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(32) NOT NULL,
    updated_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    park_name VARCHAR(128) NOT NULL
);
ALTER TABLE o_eco_park ADD CONSTRAINT uk_eco_park_name UNIQUE(park_name);
COMMENT ON TABLE o_eco_park IS '园区信息表';
COMMENT ON COLUMN o_eco_park.id IS '主键ID';
COMMENT ON COLUMN o_eco_park.park_name IS '园区名称';
COMMENT ON COLUMN o_eco_park.created_time IS '创建时间';
COMMENT ON COLUMN o_eco_park.updated_time IS '更新时间';
COMMENT ON COLUMN o_eco_park.created_by IS '创建人';
COMMENT ON COLUMN o_eco_park.updated_by IS '更新人';

CREATE TABLE o_eco_building (
    id VARCHAR(32) PRIMARY KEY,
    created_by VARCHAR(32) NOT NULL,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(32) NOT NULL,
    updated_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    building_name VARCHAR(128) NOT NULL,
    park_id VARCHAR(32) NOT NULL,
    index INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (park_id) REFERENCES o_eco_park(id)
);
ALTER TABLE o_eco_building ADD CONSTRAINT uk_eco_building_name_park UNIQUE(building_name, park_id);
COMMENT ON TABLE o_eco_building IS '楼栋信息表';
COMMENT ON COLUMN o_eco_building.id IS '主键ID';
COMMENT ON COLUMN o_eco_building.building_name IS '楼栋名称';
COMMENT ON COLUMN o_eco_building.park_id IS '园区ID';
COMMENT ON COLUMN o_eco_building.created_time IS '创建时间';
COMMENT ON COLUMN o_eco_building.updated_time IS '更新时间';
COMMENT ON COLUMN o_eco_building.created_by IS '创建人';
COMMENT ON COLUMN o_eco_building.updated_by IS '更新人';
COMMENT ON COLUMN o_eco_building.index IS '排序索引';

CREATE TABLE o_eco_floor (
    id VARCHAR(32) PRIMARY KEY,
    created_by VARCHAR(32) NOT NULL,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(32) NOT NULL,
    updated_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    floor_name VARCHAR(128) NOT NULL,
    building_id VARCHAR(32) NOT NULL,
    park_id VARCHAR(32) NOT NULL,
    index INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (building_id) REFERENCES o_eco_building(id),
    FOREIGN KEY (park_id) REFERENCES o_eco_park(id)
);
ALTER TABLE o_eco_floor ADD CONSTRAINT uk_eco_floor_name_building UNIQUE(floor_name, building_id);
COMMENT ON TABLE o_eco_floor IS '楼层信息表';
COMMENT ON COLUMN o_eco_floor.id IS '主键ID';
COMMENT ON COLUMN o_eco_floor.floor_name IS '楼层名称';
COMMENT ON COLUMN o_eco_floor.building_id IS '楼栋ID';
COMMENT ON COLUMN o_eco_floor.park_id IS '园区ID';
COMMENT ON COLUMN o_eco_floor.created_time IS '创建时间';
COMMENT ON COLUMN o_eco_floor.updated_time IS '更新时间';
COMMENT ON COLUMN o_eco_floor.created_by IS '创建人';
COMMENT ON COLUMN o_eco_floor.updated_by IS '更新人';
COMMENT ON COLUMN o_eco_floor.index IS '排序索引';

CREATE TABLE o_eco_gateway (
    id VARCHAR(32) PRIMARY KEY,
    created_by VARCHAR(32) NOT NULL,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(32) NOT NULL,
    updated_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    model_name VARCHAR(128) NOT NULL,
    dev_name VARCHAR(128) NOT NULL,
    mac_addr VARCHAR(64) NOT NULL,
    cm_code VARCHAR(64) NOT NULL,
    project_code VARCHAR(64) NOT NULL DEFAULT '',
    location VARCHAR(128) NOT NULL,
    floor_id VARCHAR(32) NOT NULL DEFAULT '',
    building_id VARCHAR(32) NOT NULL DEFAULT '',
    park_id VARCHAR(32) NOT NULL DEFAULT '',
    type INTEGER NOT NULL,
    level INTEGER NOT NULL DEFAULT 0,
    collect_type INTEGER NOT NULL DEFAULT 0,
    real_data_value DECIMAL(20,2) NOT NULL DEFAULT 0
);
ALTER TABLE o_eco_gateway ADD CONSTRAINT uk_eco_gateway_mac_addr UNIQUE(mac_addr);
COMMENT ON TABLE o_eco_gateway IS '配电网关信息表';
COMMENT ON COLUMN o_eco_gateway.id IS '主键ID';
COMMENT ON COLUMN o_eco_gateway.model_name IS '型号名称';
COMMENT ON COLUMN o_eco_gateway.dev_name IS '设备名称';
COMMENT ON COLUMN o_eco_gateway.mac_addr IS 'MAC地址';
COMMENT ON COLUMN o_eco_gateway.cm_code IS '通信码';
COMMENT ON COLUMN o_eco_gateway.project_code IS '项目编码';
COMMENT ON COLUMN o_eco_gateway.location IS '组织名称';
COMMENT ON COLUMN o_eco_gateway.created_time IS '创建时间';
COMMENT ON COLUMN o_eco_gateway.updated_time IS '更新时间';
COMMENT ON COLUMN o_eco_gateway.created_by IS '创建人';
COMMENT ON COLUMN o_eco_gateway.updated_by IS '更新人';
COMMENT ON COLUMN o_eco_gateway.floor_id IS '楼层ID';
COMMENT ON COLUMN o_eco_gateway.building_id IS '楼栋ID';
COMMENT ON COLUMN o_eco_gateway.park_id IS '园区ID';
COMMENT ON COLUMN o_eco_gateway.type IS '网关类型(1:AL,2:AP)';
COMMENT ON COLUMN o_eco_gateway.collect_type IS '采集类型(0:配电平台,1:IOT)';
COMMENT ON COLUMN o_eco_gateway.level IS '层级(0:园区,1:楼栋,2:楼层)';
COMMENT ON COLUMN o_eco_gateway.real_data_value IS '实时数据值';




CREATE TABLE o_eco_water_meter (
    id VARCHAR(32) PRIMARY KEY,
    created_by VARCHAR(32) NOT NULL,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(32) NOT NULL,
    updated_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    model_name VARCHAR(128) NOT NULL,
    dev_name VARCHAR(128) NOT NULL,
    channel_no VARCHAR(64) NOT NULL,
    cm_code VARCHAR(64) NOT NULL,
    location VARCHAR(128) NOT NULL,
    building_id VARCHAR(32) NOT NULL,
    park_id VARCHAR(32) NOT NULL,
    type INTEGER NOT NULL,
    total_value DECIMAL(20,2) NOT NULL DEFAULT 0,
    FOREIGN KEY (building_id) REFERENCES o_eco_building(id),
    FOREIGN KEY (park_id) REFERENCES o_eco_park(id)
);
ALTER TABLE o_eco_water_meter ADD CONSTRAINT uk_eco_water_meter_cm_code UNIQUE(cm_code);
COMMENT ON TABLE o_eco_water_meter IS '水表信息表';
COMMENT ON COLUMN o_eco_water_meter.id IS '主键ID';
COMMENT ON COLUMN o_eco_water_meter.model_name IS '型号名称';
COMMENT ON COLUMN o_eco_gateway.dev_name IS '设备名称';
COMMENT ON COLUMN o_eco_water_meter.channel_no IS '通道号';
COMMENT ON COLUMN o_eco_water_meter.cm_code IS '通信码';
COMMENT ON COLUMN o_eco_water_meter.location IS '组织名称';
COMMENT ON COLUMN o_eco_water_meter.created_time IS '创建时间';
COMMENT ON COLUMN o_eco_water_meter.updated_time IS '更新时间';
COMMENT ON COLUMN o_eco_water_meter.created_by IS '创建人';
COMMENT ON COLUMN o_eco_water_meter.updated_by IS '更新人';
COMMENT ON COLUMN o_eco_water_meter.building_id IS '楼栋ID';
COMMENT ON COLUMN o_eco_water_meter.park_id IS '园区ID';
COMMENT ON COLUMN o_eco_water_meter.type IS '水表类型(1:低区,2:高区)';
COMMENT ON COLUMN o_eco_water_meter.total_value IS '总用水量';

-- Function to generate test data
CREATE OR REPLACE FUNCTION generate_gateway_test_data(start_date DATE, end_date DATE)
RETURNS void AS $$
DECLARE
    curr_time TIMESTAMP;
    gateway RECORD;
    random_consumption DECIMAL(20,2);
BEGIN
    -- Loop through each gateway
    FOR gateway IN SELECT id, floor_id, building_id, park_id, type FROM o_eco_gateway
    LOOP
        -- Loop through each day
        curr_time := start_date;
        WHILE curr_time < end_date + INTERVAL '1 day' LOOP
            -- Insert 24 records for each hour
            FOR i IN 0..23 LOOP
                -- Generate random consumption between 10 and 100 kWh
                random_consumption := (random() * 90 + 10)::DECIMAL(20,2);
                
                INSERT INTO f_eco_gateway_1h (
                    id,
                    time,
                    gateway_id,
                    floor_id,
                    building_id,
                    park_id,
                    type,
                    power_consumption
                ) VALUES (
                    md5(random()::text)::varchar(32),
                    curr_time + (i || ' hours')::INTERVAL,
                    gateway.id,
                    gateway.floor_id,
                    gateway.building_id,
                    gateway.park_id,
                    gateway.type,
                    random_consumption
                );
            END LOOP;
            curr_time := curr_time + INTERVAL '1 day';
        END LOOP;
    END LOOP;
    
END;
$$ LANGUAGE plpgsql;

CREATE TABLE f_eco_water_meter_1h (
   id VARCHAR(32),
    time TIMESTAMP NOT NULL,
    water_meter_id VARCHAR(32) NOT NULL,
    building_id VARCHAR(32) NOT NULL,
    park_id VARCHAR(32) NOT NULL,
    type INTEGER NOT NULL,
    water_consumption DECIMAL(20,2) NOT NULL,
    FOREIGN KEY (water_meter_id) REFERENCES o_eco_water_meter(id),
    FOREIGN KEY (building_id) REFERENCES o_eco_building(id),
    FOREIGN KEY (park_id) REFERENCES o_eco_park(id),
    PRIMARY KEY (id, time)
);
SELECT create_hypertable('f_eco_water_meter_1h', 'time');
CREATE INDEX idx_water_meter_1h_water_meter_id ON f_eco_water_meter_1h(water_meter_id, time DESC);  
CREATE INDEX idx_water_meter_1h_building_id ON f_eco_water_meter_1h(building_id, time DESC);
CREATE INDEX idx_water_meter_1h_park_id ON f_eco_water_meter_1h(park_id, time DESC);    

-- Create continuous aggregates for park water hourly metrics
CREATE MATERIALIZED VIEW f_eco_park_water_1h
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 hour', time) AS time,
       park_id,
       type,
       sum(water_consumption) as water_consumption
FROM f_eco_water_meter_1h
GROUP BY time_bucket(INTERVAL '1 hour', time), park_id, type
WITH NO DATA;

-- Create continuous aggregates for park water daily metrics
CREATE MATERIALIZED VIEW f_eco_park_water_1d
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 day', time) AS time,
       park_id,
       type,
       sum(water_consumption) as water_consumption
FROM f_eco_water_meter_1h
GROUP BY time_bucket(INTERVAL '1 day', time), park_id, type
WITH NO DATA;

-- Create continuous aggregates for park water monthly metrics
CREATE MATERIALIZED VIEW f_eco_park_water_1m
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 month', time) AS time,
       park_id,
       type,
       sum(water_consumption) as water_consumption
FROM f_eco_water_meter_1h
GROUP BY time_bucket(INTERVAL '1 month', time), park_id, type
WITH NO DATA;

-- Create continuous aggregates for park water yearly metrics
CREATE MATERIALIZED VIEW f_eco_park_water_1y
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 year', time) AS time,
       park_id,
       type,
       sum(water_consumption) as water_consumption
FROM f_eco_water_meter_1h
GROUP BY time_bucket(INTERVAL '1 year', time), park_id, type
WITH NO DATA;

-- Create continuous aggregates for building water hourly metrics
CREATE MATERIALIZED VIEW f_eco_water_building_1h
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 hour', time) AS time,
       building_id,
       park_id,
       type,
       sum(water_consumption) as water_consumption
FROM f_eco_water_meter_1h
GROUP BY time_bucket(INTERVAL '1 hour', time), building_id, park_id, type
WITH NO DATA;

-- Create continuous aggregates for building water daily metrics
CREATE MATERIALIZED VIEW f_eco_water_building_1d
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 day', time) AS time,
       building_id,
       park_id,
       type,
       sum(water_consumption) as water_consumption
FROM f_eco_water_meter_1h
GROUP BY time_bucket(INTERVAL '1 day', time), building_id, park_id, type
WITH NO DATA;

-- Create continuous aggregates for building water monthly metrics
CREATE MATERIALIZED VIEW f_eco_water_building_1m
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 month', time) AS time,
       building_id,
       park_id,
       type,
       sum(water_consumption) as water_consumption
FROM f_eco_water_meter_1h
GROUP BY time_bucket(INTERVAL '1 month', time), building_id, park_id, type
WITH NO DATA;

-- Create continuous aggregates for building water yearly metrics
CREATE MATERIALIZED VIEW f_eco_water_building_1y
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 year', time) AS time,
       building_id,
       park_id,
       type,
       sum(water_consumption) as water_consumption
FROM f_eco_water_meter_1h
GROUP BY time_bucket(INTERVAL '1 year', time), building_id, park_id, type
WITH NO DATA;

-- Create hypertable for gateway hourly metrics
CREATE TABLE f_eco_gateway_1h (
    id VARCHAR(32),
    time TIMESTAMP NOT NULL,
    gateway_id VARCHAR(32) NOT NULL,
    floor_id VARCHAR(32) NOT NULL DEFAULT '',
    building_id VARCHAR(32) NOT NULL DEFAULT '',
    park_id VARCHAR(32) NOT NULL DEFAULT '',
    type INTEGER NOT NULL,
    level INTEGER NOT NULL DEFAULT 0,
    power_consumption DECIMAL(20,2) NOT NULL,
    FOREIGN KEY (gateway_id) REFERENCES o_eco_gateway(id),
    PRIMARY KEY (id, time)
);
SELECT create_hypertable('f_eco_gateway_1h', 'time');
CREATE INDEX idx_gateway_1h_gateway_id ON f_eco_gateway_1h(gateway_id, time DESC);
CREATE INDEX idx_gateway_1h_floor_id ON f_eco_gateway_1h(floor_id, time DESC);
CREATE INDEX idx_gateway_1h_building_id ON f_eco_gateway_1h(building_id, time DESC);
CREATE INDEX idx_gateway_1h_park_id ON f_eco_gateway_1h(park_id, time DESC);
COMMENT ON TABLE f_eco_gateway_1h IS '配电网关小时粒度性能数据';
COMMENT ON COLUMN f_eco_gateway_1h.time IS '时间';
COMMENT ON COLUMN f_eco_gateway_1h.gateway_id IS '网关ID';
COMMENT ON COLUMN f_eco_gateway_1h.floor_id IS '楼层ID';
COMMENT ON COLUMN f_eco_gateway_1h.building_id IS '楼栋ID';
COMMENT ON COLUMN f_eco_gateway_1h.park_id IS '园区ID';
COMMENT ON COLUMN f_eco_gateway_1h.type IS '网关类型(1:AL,2:AP)';
COMMENT ON COLUMN f_eco_gateway_1h.level IS '层级(0:园区,1:楼栋,2:楼层)';
COMMENT ON COLUMN f_eco_gateway_1h.power_consumption IS '用电量(kWh)';

CREATE INDEX idx_gateway_1h_level ON f_eco_gateway_1h(level, time DESC);

-- Create continuous aggregates for gateway daily metrics
CREATE MATERIALIZED VIEW f_eco_gateway_1d
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 day', time) AS time,
       gateway_id,
       floor_id,
       building_id,
       park_id,
       type,
       level,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h
GROUP BY time_bucket(INTERVAL '1 day', time), gateway_id, floor_id, building_id, park_id, type, level
WITH NO DATA;

-- Create continuous aggregates for gateway monthly metrics
CREATE MATERIALIZED VIEW f_eco_gateway_1m
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 month', time) AS time,
       gateway_id,
       floor_id,
       building_id,
       park_id,
       type,
       level,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h
GROUP BY time_bucket(INTERVAL '1 month', time), gateway_id, floor_id, building_id, park_id, type, level
WITH NO DATA;

-- Create continuous aggregates for gateway yearly metrics
CREATE MATERIALIZED VIEW f_eco_gateway_1y
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 year', time) AS time,
       gateway_id,
       floor_id,
       building_id,
       park_id,
       type,    
       level,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h
GROUP BY time_bucket(INTERVAL '1 year', time), gateway_id, floor_id, building_id, park_id, type, level
WITH NO DATA;

-- Create continuous aggregates for floor hourly metrics
CREATE MATERIALIZED VIEW f_eco_floor_1h
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 hour', time) AS time,
       floor_id,
       building_id,
       park_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h where level=2
GROUP BY time_bucket(INTERVAL '1 hour', time), floor_id, building_id, park_id, type, level
WITH NO DATA;

-- Create continuous aggregates for floor daily metrics
CREATE MATERIALIZED VIEW f_eco_floor_1d
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 day', time) AS time,
       floor_id,
       building_id,
       park_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h where level=2
GROUP BY time_bucket(INTERVAL '1 day', time), floor_id, building_id, park_id, type
WITH NO DATA;

-- Create continuous aggregates for floor monthly metrics
CREATE MATERIALIZED VIEW f_eco_floor_1m
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 month', time) AS time,
       floor_id,
       building_id,
       park_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h where level=2
GROUP BY time_bucket(INTERVAL '1 month', time), floor_id, building_id, park_id, type
WITH NO DATA;

-- Create continuous aggregates for floor yearly metrics
CREATE MATERIALIZED VIEW f_eco_floor_1y
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 year', time) AS time,
       floor_id,
       building_id,
       park_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h where level=2
GROUP BY time_bucket(INTERVAL '1 year', time), floor_id, building_id, park_id, type
WITH NO DATA;

-- Create continuous aggregates for building hourly metrics
CREATE MATERIALIZED VIEW f_eco_building_1h
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 hour', time) AS time,
       building_id,
       park_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h where level=1
GROUP BY time_bucket(INTERVAL '1 hour', time), building_id, park_id, type
WITH NO DATA;

-- Create continuous aggregates for building daily metrics
CREATE MATERIALIZED VIEW f_eco_building_1d
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 day', time) AS time,
       building_id,
       park_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h where level=1
GROUP BY time_bucket(INTERVAL '1 day', time), building_id, park_id, type
WITH NO DATA;

-- Create continuous aggregates for building monthly metrics
CREATE MATERIALIZED VIEW f_eco_building_1m
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 month', time) AS time,
       building_id,
       park_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h where level=1
GROUP BY time_bucket(INTERVAL '1 month', time), building_id, park_id, type
WITH NO DATA;

-- Create continuous aggregates for building yearly metrics
CREATE MATERIALIZED VIEW f_eco_building_1y
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 year', time) AS time,
       building_id,
       park_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h where level=1
GROUP BY time_bucket(INTERVAL '1 year', time), building_id, park_id, type
WITH NO DATA;

-- Create continuous aggregates for park hourly metrics
CREATE MATERIALIZED VIEW f_eco_park_1h
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 hour', time) AS time,
       park_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h where level=1
GROUP BY time_bucket(INTERVAL '1 hour', time), park_id, type
WITH NO DATA;


-- Create continuous aggregates for park daily metrics
CREATE MATERIALIZED VIEW f_eco_park_1d
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 day', time) AS time,
       park_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h where level=1
GROUP BY time_bucket(INTERVAL '1 day', time), park_id, type
WITH NO DATA;

-- Create continuous aggregates for park monthly metrics
CREATE MATERIALIZED VIEW f_eco_park_1m
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 month', time) AS time,
       park_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h where level=1
GROUP BY time_bucket(INTERVAL '1 month', time), park_id, type
WITH NO DATA;

-- Create continuous aggregates for park yearly metrics
CREATE MATERIALIZED VIEW f_eco_park_1y
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 year', time) AS time,
       park_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h where level=1
GROUP BY time_bucket(INTERVAL '1 year', time), park_id, type
WITH NO DATA;


-- Insert park
INSERT INTO o_eco_park (id, created_by, updated_by, park_name)
VALUES (md5('教科院'), 'admin', 'admin', '教科院');

-- Insert buildings
INSERT INTO o_eco_building (id, created_by, updated_by, building_name, park_id, index)
VALUES 
(md5('A栋'), 'admin', 'admin', 'A栋', md5('教科院'), 1),
(md5('B栋'), 'admin', 'admin', 'B栋', md5('教科院'), 2),
(md5('C栋'), 'admin', 'admin', 'C栋', md5('教科院'), 3),
(md5('E栋'), 'admin', 'admin', 'E/G栋', md5('教科院'), 4),
(md5('H栋'), 'admin', 'admin', 'H栋', md5('教科院'), 6),
(md5('功能厅'), 'admin', 'admin', '功能厅', md5('教科院'), 7),
(md5('充电桩'), 'admin', 'admin', '充电桩', md5('教科院'), 8);

-- Insert floors
INSERT INTO o_eco_floor (id, created_by, updated_by, floor_name, building_id, park_id, index)
VALUES
(md5('A栋_一层'), 'admin', 'admin', '一层', md5('A栋'), md5('教科院'), 1),
(md5('A栋_二层'), 'admin', 'admin', '二层', md5('A栋'), md5('教科院'), 2),
(md5('A栋_三层'), 'admin', 'admin', '三层', md5('A栋'), md5('教科院'), 3),
(md5('A栋_四层'), 'admin', 'admin', '四层', md5('A栋'), md5('教科院'), 4),
(md5('A栋_五层'), 'admin', 'admin', '五层', md5('A栋'), md5('教科院'), 5),
(md5('A栋_六层'), 'admin', 'admin', '六层', md5('A栋'), md5('教科院'), 6),
(md5('B栋_一层'), 'admin', 'admin', '一层', md5('B栋'), md5('教科院'), 1),
(md5('B栋_二层'), 'admin', 'admin', '二层', md5('B栋'), md5('教科院'), 2),
(md5('B栋_三层'), 'admin', 'admin', '三层', md5('B栋'), md5('教科院'), 3),
(md5('B栋_四层'), 'admin', 'admin', '四层', md5('B栋'), md5('教科院'), 4),
(md5('B栋_五层'), 'admin', 'admin', '五层', md5('B栋'), md5('教科院'), 5),
(md5('B栋_六层'), 'admin', 'admin', '六层', md5('B栋'), md5('教科院'), 6),
(md5('C栋_一层'), 'admin', 'admin', '一层', md5('C栋'), md5('教科院'), 1),
(md5('C栋_二层'), 'admin', 'admin', '二层', md5('C栋'), md5('教科院'), 2),
(md5('C栋_三层'), 'admin', 'admin', '三层', md5('C栋'), md5('教科院'), 3),
(md5('C栋_四层'), 'admin', 'admin', '四层', md5('C栋'), md5('教科院'), 4),
(md5('C栋_五层'), 'admin', 'admin', '五层', md5('C栋'), md5('教科院'), 5),
(md5('C栋_六层'), 'admin', 'admin', '六层', md5('C栋'), md5('教科院'), 6),
(md5('C栋_七层'), 'admin', 'admin', '七层', md5('C栋'), md5('教科院'), 7),
(md5('E栋'), 'admin', 'admin', '整栋', md5('E栋'), md5('教科院'), 1),
(md5('H栋'), 'admin', 'admin', '整栋', md5('H栋'), md5('教科院'), 1),
(md5('功能厅'), 'admin', 'admin', '整栋', md5('功能厅'), md5('教科院'), 1),
(md5('充电桩'), 'admin', 'admin', '整栋', md5('充电桩'), md5('教科院'), 1);





-- Insert gateways
INSERT INTO o_eco_gateway (id, created_by, updated_by, mac_addr, model_name, dev_name, cm_code, location, floor_id, building_id, park_id, type,level)
VALUES
('98CC4D150A00', 'admin', 'admin', '98CC4D150A00', '配电网关', '配电网关_A-AL-01_98CC4D150A00', '20000000000668', 'A栋', '', md5('A栋'), md5('教科院'), 1, 1),
('98CC4D1528E2', 'admin', 'admin', '98CC4D1528E2', '配电网关', '配电网关_A-AL-02_98CC4D1528E2', '20000000000665', 'A栋', '', md5('A栋'), md5('教科院'), 1, 1),
('98CC4D152928', 'admin', 'admin', '98CC4D152928', '配电网关', '配电网关_A-AL-03_98CC4D152928', '20000000000663', 'A栋', '', md5('A栋'), md5('教科院'), 1, 1),
('98CC4D152AE4', 'admin', 'admin', '98CC4D152AE4', '配电网关', '配电网关_A-AP-01_98CC4D152AE4', '20000000000654', 'A栋', '', md5('A栋'), md5('教科院'), 2, 1),
('98CC4D151E04', 'admin', 'admin', '98CC4D151E04', '配电网关', '配电网关_A-AP-02_98CC4D151E04', '20000000000659', 'A栋', '', md5('A栋'), md5('教科院'), 2, 1),
('98CC4D152990', 'admin', 'admin', '98CC4D152990', '配电网关', '配电网关_B-AL-01_98CC4D152990', '20000000000655', 'B栋', '', md5('B栋'), md5('教科院'), 1, 1),
('98CC4D152AC8', 'admin', 'admin', '98CC4D152AC8', '配电网关', '配电网关_B-AL-02_98CC4D152AC8', '20000000000656', 'B栋', '', md5('B栋'), md5('教科院'), 1, 1),
('98CC4D150E1A', 'admin', 'admin', '98CC4D150E1A', '配电网关', '配电网关_B-AP-01_98CC4D150E1A', '20000000000652', 'B栋', '', md5('B栋'), md5('教科院'), 2, 1),
('98CC4D151DF2', 'admin', 'admin', '98CC4D151DF2', '配电网关', '配电网关_B-AP-02_98CC4D151DF2', '20000000000653', 'B栋', '', md5('B栋'), md5('教科院'), 2, 1),
('98CC4D149A06', 'admin', 'admin', '98CC4D149A06', '配电网关', '配电网关_C-AL-01_98CC4D149A06', '20000000000660', 'C栋', '', md5('C栋'), md5('教科院'), 1, 1),
('98CC4D151D3A', 'admin', 'admin', '98CC4D151D3A', '配电网关', '配电网关_C-AL-02_98CC4D151D3A', '20000000000658', 'C栋', '', md5('C栋'), md5('教科院'), 1, 1),
('98CC4D151C56', 'admin', 'admin', '98CC4D151C56', '配电网关', '配电网关_C-AP-01_98CC4D151C56', '20000000000666', 'C栋', '', md5('C栋'), md5('教科院'), 2, 1),
('98CC4D151E00', 'admin', 'admin', '98CC4D151E00', '配电网关', '配电网关_C-AP-02_98CC4D151E00', '20000000000664', 'C栋', '', md5('C栋'), md5('教科院'), 2, 1),
('98CC4D150E66', 'admin', 'admin', '98CC4D150E66', '配电网关', '配电网关_E-AL-01_98CC4D150E66', '20000000000661', 'E栋', '', md5('E栋'), md5('教科院'), 2, 1),
('98CC4D150A3C', 'admin', 'admin', '98CC4D150A3C', '配电网关', '配电网关_E-AL-02_98CC4D150A3C', '20000000000657', 'E栋', '', md5('E栋'), md5('教科院'), 2, 1),
('98CC4D149A0C', 'admin', 'admin', '98CC4D149A0C', '配电网关', '配电网关_E-AP-01_98CC4D149A0C', '20000000000667', 'H栋', '', md5('H栋'), md5('教科院'), 2, 1),
('98CC4D1528E4', 'admin', 'admin', '98CC4D1528E4', '配电网关', '配电网关_E-AP-02_98CC4D1528E4', '20000000000662', 'H栋', '', md5('H栋'), md5('教科院'), 2, 1),
('98CC4D151D88', 'admin', 'admin', '98CC4D151D88', '配电网关', '配电网关_A-AL-01_98CC4D151D88', '20000000000036', 'A栋_一层', md5('A栋_一层'), md5('A栋'), md5('教科院'), 1,2),
('98CC4D150BD4', 'admin', 'admin', '98CC4D150BD4', '配电网关', '配电网关_A-AL-02-1_98CC4D150BD4', '20000000000708', 'A栋_二层', md5('A栋_二层'), md5('A栋'), md5('教科院'), 1, 2),
('98CC4D151DD0', 'admin', 'admin', '98CC4D151DD0', '配电网关', '配电网关_A-AL-03-1_98CC4D151DD0', '20000000000162', 'A栋_三层', md5('A栋_三层'), md5('A栋'), md5('教科院'), 1, 2),
('98CC4D151D9E', 'admin', 'admin', '98CC4D151D9E', '配电网关', '配电网关_A-AL-04-1_98CC4D151D9E', '20000000000217', 'A栋_四层', md5('A栋_四层'), md5('A栋'), md5('教科院'), 1, 2),
('98CC4D151D8C', 'admin', 'admin', '98CC4D151D8C', '配电网关', '配电网关_A-AL-05-1_98CC4D151D8C', '20000000000246', 'A栋_五层', md5('A栋_五层'), md5('A栋'), md5('教科院'), 1, 2),
('98CC4D149A0A', 'admin', 'admin', '98CC4D149A0A', '配电网关', '配电网关_A-AL-06-1_98CC4D149A0A', '20000000000705', 'A栋_六层', md5('A栋_六层'), md5('A栋'), md5('教科院'), 1, 2),
('98CC4D1528C4', 'admin', 'admin', '98CC4D1528C4', '配电网关', '配电网关_B-AL-01-1_98CC4D1528C4', '20000000000321', 'B栋_一层', md5('B栋_一层'), md5('B栋'), md5('教科院'), 1, 2),
('98CC4D152AD8', 'admin', 'admin', '98CC4D152AD8', '配电网关', '配电网关_B-AL-02-1_98CC4D152AD8', '20000000000368', 'B栋_二层', md5('B栋_二层'), md5('B栋'), md5('教科院'), 1, 2),
('98CC4D15298C', 'admin', 'admin', '98CC4D15298C', '配电网关', '配电网关_B-AL-03-1_98CC4D15298C', '20000000000707', 'B栋_三层', md5('B栋_三层'), md5('B栋'), md5('教科院'), 1, 2),
('98CC4D151E44', 'admin', 'admin', '98CC4D151E44', '配电网关', '配电网关_B-AL-04-1_98CC4D151E44', '20000000000416', 'B栋_四层', md5('B栋_四层'), md5('B栋'), md5('教科院'), 1, 2),
('98CC4D151C54', 'admin', 'admin', '98CC4D151C54', '配电网关', '配电网关_B-AL-05-1_98CC4D151C54', '20000000000709', 'B栋_五层', md5('B栋_五层'), md5('B栋'), md5('教科院'), 1, 2),
('98CC4D151F48', 'admin', 'admin', '98CC4D151F48', '配电网关', '配电网关_B-AL-06-1_98CC4D151F48', '20000000000710', 'B栋_六层', md5('B栋_六层'), md5('B栋'), md5('教科院'), 1, 2),
('98CC4D151F46', 'admin', 'admin', '98CC4D151F46', '配电网关', '配电网关_C-AL-01-1_98CC4D151F46', '20000000000484', 'C栋_一层', md5('C栋_一层'), md5('C栋'), md5('教科院'), 1, 2),
('98CC4D151C58', 'admin', 'admin', '98CC4D151C58', '配电网关', '配电网关_C-AL-02-1_98CC4D151C58', '20000000000513', 'C栋_二层', md5('C栋_二层'), md5('C栋'), md5('教科院'), 1, 2),
('98CC4D151DA0', 'admin', 'admin', '98CC4D151DA0', '配电网关', '配电网关_C-AL-03-1_98CC4D151DA0', '20000000000560', 'C栋_三层', md5('C栋_三层'), md5('C栋'), md5('教科院'), 1, 2),
('98CC4D151DD8', 'admin', 'admin', '98CC4D151DD8', '配电网关', '配电网关_C-AL-04-1_98CC4D151DD8', '20000000000605', 'C栋_四层', md5('C栋_四层'), md5('C栋'), md5('教科院'), 1, 2),
('98CC4D151E02', 'admin', 'admin', '98CC4D151E02', '配电网关', '配电网关_C-AL-05-1_98CC4D151E02', '21000000000416', 'C栋_五层', md5('C栋_五层'), md5('C栋'), md5('教科院'), 1, 2),
('98CC4D151DD6', 'admin', 'admin', '98CC4D151DD6', '配电网关', '配电网关_C-AL-06-1_98CC4D151DD6', '21000000000445', 'C栋_六层', md5('C栋_六层'), md5('C栋'), md5('教科院'), 1, 2);


INSERT INTO o_eco_gateway (id, created_by, updated_by, mac_addr, model_name, dev_name, cm_code, location, floor_id, building_id, park_id, type,level,collect_type)
VALUES
('98CC4D15209A', 'admin', 'admin', '98CC4D15209A', '配电网关', '配电网关_C-AL-07-1_98CC4D15209A', '23100000000022', 'C栋_七层',  md5('C栋_七层'), md5('C栋'), md5('教科院'), 1, 2, 0);


--充电桩
INSERT INTO o_eco_gateway (id, created_by, updated_by, mac_addr, model_name, dev_name, cm_code, location, floor_id, building_id, park_id, type,level,collect_type)
VALUES
('23000000000022', 'admin', 'admin', '23000000000022', '配电网关', '配电网关_充电桩_23000000000022', '23000000000022', '充电桩', '', md5('充电桩'), md5('教科院'), 2, 1, 1);

--多功能厅
INSERT INTO o_eco_gateway (id, created_by, updated_by, mac_addr, model_name, dev_name, cm_code, location, floor_id, building_id, park_id, type,level,collect_type)
VALUES
('98CC4D151DA4', 'admin', 'admin', '98CC4D151DA4', '配电网关', '配电网关_功能厅_98CC4D151DA4', '23100000000022', '功能厅', '', md5('功能厅'), md5('教科院'), 2, 1, 0);

-- Insert water meters
INSERT INTO o_eco_water_meter (id, created_by, updated_by, model_name, dev_name, channel_no, cm_code, location, building_id, park_id, type)
VALUES
('30002050', 'admin', 'admin', '水表', 'A座高区水表', '30002050', '24000000000001', 'A栋', md5('A栋'), md5('教科院'), 2),
('3000205C', 'admin', 'admin', '水表', 'A座低区水表', '3000205C', '24000000000002', 'A栋', md5('A栋'), md5('教科院'), 1);

-- 2025-5-23 begin
update o_eco_building set building_name = 'E栋' where id = md5('E栋');

-- Insert buildings
INSERT INTO o_eco_building (id, created_by, updated_by, building_name, park_id, index)
VALUES 
(md5('G栋'), 'admin', 'admin', 'G栋', md5('教科院'), 5);

-- Insert floors
INSERT INTO o_eco_floor (id, created_by, updated_by, floor_name, building_id, park_id, index)
VALUES
(md5('G栋'), 'admin', 'admin', '整栋', md5('G栋'), md5('教科院'), 1);

ALTER TABLE o_eco_gateway ADD COLUMN factor decimal(10,4) DEFAULT 1;


--宿舍楼
INSERT INTO o_eco_gateway (id, created_by, updated_by, mac_addr, model_name, dev_name, cm_code, location, floor_id, building_id, park_id, type,level,collect_type,factor)
VALUES
('23000000000023', 'admin', 'admin', '23000000000023', '配电网关', '配电网关_宿舍楼_23000000000023', '23000000000023', '宿舍楼', '', md5('G栋'), md5('教科院'), 1, 1, 1,0.625);

-- 2025-5-23 end
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP MATERIALIZED VIEW IF EXISTS f_eco_building_water_1h;
DROP MATERIALIZED VIEW IF EXISTS f_eco_building_water_1d;
DROP MATERIALIZED VIEW IF EXISTS f_eco_building_water_1m;
DROP MATERIALIZED VIEW IF EXISTS f_eco_building_water_1y;
DROP MATERIALIZED VIEW IF EXISTS f_eco_park_water_1m;
DROP MATERIALIZED VIEW IF EXISTS f_eco_park_water_1y;
DROP MATERIALIZED VIEW IF EXISTS f_eco_park_water_1d;
DROP MATERIALIZED VIEW IF EXISTS f_eco_park_water_1h;
DROP TABLE IF EXISTS f_eco_water_meter_1h cascade;

DROP MATERIALIZED VIEW IF EXISTS f_eco_park_1h;
DROP MATERIALIZED VIEW IF EXISTS f_eco_park_1d;
DROP MATERIALIZED VIEW IF EXISTS f_eco_park_1m;
DROP MATERIALIZED VIEW IF EXISTS f_eco_park_1y;
DROP FUNCTION IF EXISTS generate_gateway_test_data;
DROP MATERIALIZED VIEW IF EXISTS f_eco_building_1y;
DROP MATERIALIZED VIEW IF EXISTS f_eco_building_1m;
DROP MATERIALIZED VIEW IF EXISTS f_eco_building_1d;
DROP MATERIALIZED VIEW IF EXISTS f_eco_building_1h;
DROP MATERIALIZED VIEW IF EXISTS f_eco_floor_1y;
DROP MATERIALIZED VIEW IF EXISTS f_eco_floor_1m;
DROP MATERIALIZED VIEW IF EXISTS f_eco_floor_1d;
DROP MATERIALIZED VIEW IF EXISTS f_eco_floor_1h;
DROP MATERIALIZED VIEW IF EXISTS f_eco_gateway_1y;
DROP MATERIALIZED VIEW IF EXISTS f_eco_gateway_1m;
DROP MATERIALIZED VIEW IF EXISTS f_eco_gateway_1d;
DROP TABLE IF EXISTS f_eco_gateway_1h cascade;
DROP TABLE IF EXISTS o_eco_water_meter cascade;
DROP TABLE IF EXISTS o_eco_gateway cascade;
DROP TABLE IF EXISTS o_eco_floor cascade;
DROP TABLE IF EXISTS o_eco_building cascade;
DROP TABLE IF EXISTS o_eco_park cascade;
-- +goose StatementEnd
