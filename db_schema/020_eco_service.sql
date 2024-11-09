-- +goose Up
-- +goose StatementBegin

CREATE TABLE o_eco_building (
    id VARCHAR(32) PRIMARY KEY,
    created_by VARCHAR(32) NOT NULL,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(32) NOT NULL,
    updated_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    building_name VARCHAR(128) NOT NULL
);
ALTER TABLE o_eco_building ADD CONSTRAINT uk_eco_building_name UNIQUE(building_name);
COMMENT ON TABLE o_eco_building IS '楼栋信息表';
COMMENT ON COLUMN o_eco_building.id IS '主键ID';
COMMENT ON COLUMN o_eco_building.building_name IS '楼栋名称';
COMMENT ON COLUMN o_eco_building.created_time IS '创建时间';
COMMENT ON COLUMN o_eco_building.updated_time IS '更新时间';
COMMENT ON COLUMN o_eco_building.created_by IS '创建人';
COMMENT ON COLUMN o_eco_building.updated_by IS '更新人';

CREATE TABLE o_eco_floor (
    id VARCHAR(32) PRIMARY KEY,
    created_by VARCHAR(32) NOT NULL,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(32) NOT NULL,
    updated_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    floor_name VARCHAR(128) NOT NULL,
    building_id VARCHAR(32) NOT NULL,
    FOREIGN KEY (building_id) REFERENCES o_eco_building(id)
);
ALTER TABLE o_eco_floor ADD CONSTRAINT uk_eco_floor_name_building UNIQUE(floor_name, building_id);
COMMENT ON TABLE o_eco_floor IS '楼层信息表';
COMMENT ON COLUMN o_eco_floor.id IS '主键ID';
COMMENT ON COLUMN o_eco_floor.floor_name IS '楼层名称';
COMMENT ON COLUMN o_eco_floor.building_id IS '楼栋ID';
COMMENT ON COLUMN o_eco_floor.created_time IS '创建时间';
COMMENT ON COLUMN o_eco_floor.updated_time IS '更新时间';
COMMENT ON COLUMN o_eco_floor.created_by IS '创建人';
COMMENT ON COLUMN o_eco_floor.updated_by IS '更新人';

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
    floor_id VARCHAR(32) NOT NULL,
    building_id VARCHAR(32) NOT NULL,
    type INTEGER NOT NULL,
    FOREIGN KEY (floor_id) REFERENCES o_eco_floor(id),
    FOREIGN KEY (building_id) REFERENCES o_eco_building(id)
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
COMMENT ON COLUMN o_eco_gateway.type IS '网关类型(1:AL,2:AP)';

-- Function to generate test data
CREATE OR REPLACE FUNCTION generate_gateway_test_data(start_date DATE, end_date DATE)
RETURNS void AS $$
DECLARE
    curr_time TIMESTAMP;
    gateway RECORD;
    random_consumption DECIMAL(20,2);
BEGIN
    -- Loop through each gateway
    FOR gateway IN SELECT id, floor_id, building_id, type FROM o_eco_gateway
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
                    type,
                    power_consumption
                ) VALUES (
                    md5(random()::text)::varchar(32),
                    curr_time + (i || ' hours')::INTERVAL,
                    gateway.id,
                    gateway.floor_id,
                    gateway.building_id,
                    gateway.type,
                    random_consumption
                );
            END LOOP;
            curr_time := curr_time + INTERVAL '1 day';
        END LOOP;
    END LOOP;
    
END;
$$ LANGUAGE plpgsql;

-- Create hypertable for gateway hourly metrics
CREATE TABLE f_eco_gateway_1h (
    id VARCHAR(32),
    time TIMESTAMP NOT NULL,
    gateway_id VARCHAR(32) NOT NULL,
    floor_id VARCHAR(32) NOT NULL,
    building_id VARCHAR(32) NOT NULL,
    type INTEGER NOT NULL,
    power_consumption DECIMAL(20,2) NOT NULL,
    FOREIGN KEY (gateway_id) REFERENCES o_eco_gateway(id),
    FOREIGN KEY (floor_id) REFERENCES o_eco_floor(id),
    FOREIGN KEY (building_id) REFERENCES o_eco_building(id),
    PRIMARY KEY (id, time)
);
SELECT create_hypertable('f_eco_gateway_1h', 'time');
CREATE INDEX idx_gateway_1h_gateway_id ON f_eco_gateway_1h(gateway_id, time DESC);
CREATE INDEX idx_gateway_1h_floor_id ON f_eco_gateway_1h(floor_id, time DESC);
CREATE INDEX idx_gateway_1h_building_id ON f_eco_gateway_1h(building_id, time DESC);
COMMENT ON TABLE f_eco_gateway_1h IS '配电网关小时粒度性能数据';
COMMENT ON COLUMN f_eco_gateway_1h.time IS '时间';
COMMENT ON COLUMN f_eco_gateway_1h.gateway_id IS '网关ID';
COMMENT ON COLUMN f_eco_gateway_1h.floor_id IS '楼层ID';
COMMENT ON COLUMN f_eco_gateway_1h.building_id IS '楼栋ID';
COMMENT ON COLUMN f_eco_gateway_1h.type IS '网关类型(1:AL,2:AP)';
COMMENT ON COLUMN f_eco_gateway_1h.power_consumption IS '用电量(kWh)';

-- Create continuous aggregates for gateway daily metrics
CREATE MATERIALIZED VIEW f_eco_gateway_1d
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 day', time) AS time,
       gateway_id,
       floor_id,
       building_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h
GROUP BY time_bucket(INTERVAL '1 day', time), gateway_id, floor_id, building_id, type
WITH NO DATA;

-- Create continuous aggregates for gateway monthly metrics
CREATE MATERIALIZED VIEW f_eco_gateway_1m
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 month', time) AS time,
       gateway_id,
       floor_id,
       building_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h
GROUP BY time_bucket(INTERVAL '1 month', time), gateway_id, floor_id, building_id, type
WITH NO DATA;

-- Create continuous aggregates for gateway yearly metrics
CREATE MATERIALIZED VIEW f_eco_gateway_1y
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 year', time) AS time,
       gateway_id,
       floor_id,
       building_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h
GROUP BY time_bucket(INTERVAL '1 year', time), gateway_id, floor_id, building_id, type
WITH NO DATA;

-- Create continuous aggregates for floor daily metrics
CREATE MATERIALIZED VIEW f_eco_floor_1d
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 day', time) AS time,
       floor_id,
       building_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h
GROUP BY time_bucket(INTERVAL '1 day', time), floor_id, building_id, type
WITH NO DATA;

-- Create continuous aggregates for floor monthly metrics
CREATE MATERIALIZED VIEW f_eco_floor_1m
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 month', time) AS time,
       floor_id,
       building_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h
GROUP BY time_bucket(INTERVAL '1 month', time), floor_id, building_id, type
WITH NO DATA;

-- Create continuous aggregates for floor yearly metrics
CREATE MATERIALIZED VIEW f_eco_floor_1y
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 year', time) AS time,
       floor_id,
       building_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h
GROUP BY time_bucket(INTERVAL '1 year', time), floor_id, building_id, type
WITH NO DATA;

-- Create continuous aggregates for building daily metrics
CREATE MATERIALIZED VIEW f_eco_building_1d
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 day', time) AS time,
       building_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h
GROUP BY time_bucket(INTERVAL '1 day', time), building_id, type
WITH NO DATA;

-- Create continuous aggregates for building monthly metrics
CREATE MATERIALIZED VIEW f_eco_building_1m
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 month', time) AS time,
       building_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h
GROUP BY time_bucket(INTERVAL '1 month', time), building_id, type
WITH NO DATA;

-- Create continuous aggregates for building yearly metrics
CREATE MATERIALIZED VIEW f_eco_building_1y
WITH (timescaledb.continuous) AS
SELECT time_bucket(INTERVAL '1 year', time) AS time,
       building_id,
       type,
       sum(power_consumption) as power_consumption
FROM f_eco_gateway_1h
GROUP BY time_bucket(INTERVAL '1 year', time), building_id, type
WITH NO DATA;

-- Create view for building daily metrics with floor details
CREATE OR REPLACE VIEW v_eco_building_info_1d AS
WITH floor_stats AS (
    SELECT 
        f.time,
        f.floor_id,
        f.building_id,
        fl.floor_name,
        sum(f.power_consumption) as floor_total,
        json_agg(json_build_object(
            'id', f.type,
            'total', f.power_consumption
        )) as floor_types
    FROM f_eco_floor_1d f
    LEFT JOIN o_eco_floor fl ON f.floor_id = fl.id
    GROUP BY f.time, f.floor_id, f.building_id, fl.floor_name
)
SELECT 
    b.time,
    b.building_id,
    bld.building_name,
    concat(b.building_id, '_', b.time) as id,
    sum(b.power_consumption) as total,
    json_agg(json_build_object(
        'id', b.type,
        'total', b.power_consumption
    )) as types,
    (
        SELECT json_agg(json_build_object(
            'id', concat(fs.floor_id, '_', b.time),
            'floor_id', fs.floor_id,
            'floor_name', fs.floor_name,
            'total', fs.floor_total,
            'types', fs.floor_types
        ))
        FROM floor_stats fs
        WHERE fs.time = b.time 
        AND fs.building_id = b.building_id
    ) as floors
FROM f_eco_building_1d b
LEFT JOIN o_eco_building bld ON b.building_id = bld.id
GROUP BY b.time, b.building_id, bld.building_name;

-- Create view for building monthly metrics with floor details
CREATE OR REPLACE VIEW v_eco_building_info_1m AS
WITH floor_stats AS (
    SELECT 
        f.time,
        f.floor_id,
        f.building_id,
        fl.floor_name,
        sum(f.power_consumption) as floor_total,
        json_agg(json_build_object(
            'id', f.type,
            'total', f.power_consumption
        )) as floor_types
    FROM f_eco_floor_1m f
    LEFT JOIN o_eco_floor fl ON f.floor_id = fl.id
    GROUP BY f.time, f.floor_id, f.building_id, fl.floor_name
)
SELECT 
    b.time,
    b.building_id,
    bld.building_name,
    concat(b.building_id, '_', b.time) as id,
    sum(b.power_consumption) as total,
    json_agg(json_build_object(
        'id', b.type,
        'total', b.power_consumption
    )) as types,
    (
        SELECT json_agg(json_build_object(
            'id', concat(fs.floor_id, '_', b.time),
            'floor_id', fs.floor_id,
            'floor_name', fs.floor_name,
            'total', fs.floor_total,
            'types', fs.floor_types
        ))
        FROM floor_stats fs
        WHERE fs.time = b.time 
        AND fs.building_id = b.building_id
    ) as floors
FROM f_eco_building_1m b
LEFT JOIN o_eco_building bld ON b.building_id = bld.id
GROUP BY b.time, b.building_id, bld.building_name;

-- Create view for building yearly metrics with floor details
CREATE OR REPLACE VIEW v_eco_building_info_1y AS
WITH floor_stats AS (
    SELECT 
        f.time,
        f.floor_id,
        f.building_id,
        fl.floor_name,
        sum(f.power_consumption) as floor_total,
        json_agg(json_build_object(
            'id', f.type,
            'total', f.power_consumption
        )) as floor_types
    FROM f_eco_floor_1y f
    LEFT JOIN o_eco_floor fl ON f.floor_id = fl.id
    GROUP BY f.time, f.floor_id, f.building_id, fl.floor_name
)
SELECT 
    b.time,
    b.building_id,
    bld.building_name,
    concat(b.building_id, '_', b.time) as id,
    sum(b.power_consumption) as total,
    json_agg(json_build_object(
        'id', b.type,
        'total', b.power_consumption
    )) as types,
    (
        SELECT json_agg(json_build_object(
            'id', concat(fs.floor_id, '_', b.time),
            'floor_id', fs.floor_id,
            'floor_name', fs.floor_name,
            'total', fs.floor_total,
            'types', fs.floor_types
        ))
        FROM floor_stats fs
        WHERE fs.time = b.time 
        AND fs.building_id = b.building_id
    ) as floors
FROM f_eco_building_1y b
LEFT JOIN o_eco_building bld ON b.building_id = bld.id
GROUP BY b.time, b.building_id, bld.building_name;

-- Insert buildings
INSERT INTO o_eco_building (id, created_by, updated_by, building_name)
VALUES 
(md5('A栋'), 'admin', 'admin', 'A栋'),
(md5('B栋'), 'admin', 'admin', 'B栋'),
(md5('C栋'), 'admin', 'admin', 'C栋'),
(md5('E栋'), 'admin', 'admin', 'E栋');

-- Insert floors
INSERT INTO o_eco_floor (id, created_by, updated_by, floor_name, building_id)
VALUES
(md5('A栋_一层'), 'admin', 'admin', '一层', md5('A栋')),
(md5('A栋_二层'), 'admin', 'admin', '二层', md5('A栋')),
(md5('A栋_三层'), 'admin', 'admin', '三层', md5('A栋')),
(md5('A栋_四层'), 'admin', 'admin', '四层', md5('A栋')),
(md5('A栋_五层'), 'admin', 'admin', '五层', md5('A栋')),
(md5('A栋_六层'), 'admin', 'admin', '六层', md5('A栋')),
(md5('B栋_一层'), 'admin', 'admin', '一层', md5('B栋')),
(md5('B栋_二层'), 'admin', 'admin', '二层', md5('B栋')),
(md5('B栋_三层'), 'admin', 'admin', '三层', md5('B栋')),
(md5('B栋_四层'), 'admin', 'admin', '四层', md5('B栋')),
(md5('B栋_五层'), 'admin', 'admin', '五层', md5('B栋')),
(md5('B栋_六层'), 'admin', 'admin', '六层', md5('B栋')),
(md5('C栋_一层'), 'admin', 'admin', '一层', md5('C栋')),
(md5('C栋_二层'), 'admin', 'admin', '二层', md5('C栋')),
(md5('C栋_三层'), 'admin', 'admin', '三层', md5('C栋')),
(md5('C栋_四层'), 'admin', 'admin', '四层', md5('C栋')),
(md5('C栋_六层'), 'admin', 'admin', '六层', md5('C栋')),
(md5('C栋_七层'), 'admin', 'admin', '七层', md5('C栋')),
(md5('E栋_一层'), 'admin', 'admin', '一层', md5('E栋')),
(md5('E栋_二层'), 'admin', 'admin', '二层', md5('E栋'));

-- Insert gateways with references to building and floor IDs
INSERT INTO o_eco_gateway (id, created_by, updated_by, model_name, dev_name, mac_addr, cm_code, location, floor_id, building_id, type)
VALUES
('98CC4D151E44', 'admin', 'admin', '配电网关', '配电网关_B-AL-04-1_98CC4D151E44', '98CC4D151E44', '20000000000416', 'B栋_四层', md5('B栋_四层'), md5('B栋'), 1),
('98CC4D152976', 'admin', 'admin', '配电网关', '配电网关_B-AL-04-2_98CC4D152976', '98CC4D152976', '20000000000445', 'B栋_四层', md5('B栋_四层'), md5('B栋'), 1),
('98CC4D152988', 'admin', 'admin', '配电网关', '配电网关_B-AL-05-2_98CC4D152988', '98CC4D152988', '20000000000458', 'B栋_五层', md5('B栋_五层'), md5('B栋'), 1),
('98CC4D151D8C', 'admin', 'admin', '配电网关', '配电网关_A-AL-05-1_98CC4D151D8C', '98CC4D151D8C', '20000000000246', 'A栋_五层', md5('A栋_五层'), md5('A栋'), 1),
('98CC4D151D88', 'admin', 'admin', '配电网关', '配电网关_A-AL-01_98CC4D151D88', '98CC4D151D88', '20000000000036', 'A栋_一层', md5('A栋_一层'), md5('A栋'), 1),
('98CC4D151F92', 'admin', 'admin', '配电网关', '配电网关_B-AL-06-2_98CC4D151F92', '98CC4D151F92', '20000000000471', 'B栋_六层', md5('B栋_六层'), md5('B栋'), 1),
('98CC4D150BC8', 'admin', 'admin', '配电网关', '配电网关_A-AL-01-3_98CC4D150BC8', '98CC4D150BC8', '20000000000096', 'A栋_一层', md5('A栋_一层'), md5('A栋'), 1),
('98CC4D151DF8', 'admin', 'admin', '配电网关', '配电网关_A-AL-01-2_98CC4D151DF8', '98CC4D151DF8', '20000000000112', 'A栋_一层', md5('A栋_一层'), md5('A栋'), 1),
('98CC4D1528BA', 'admin', 'admin', '配电网关', '配电网关_A-AL-02-2_98CC4D1528BA', '98CC4D1528BA', '20000000000133', 'A栋_二层', md5('A栋_二层'), md5('A栋'), 1),
('98CC4D151DD0', 'admin', 'admin', '配电网关', '配电网关_A-AL-03-1_98CC4D151DD0', '98CC4D151DD0', '20000000000162', 'A栋_三层', md5('A栋_三层'), md5('A栋'), 1),
('98CC4D151E0C', 'admin', 'admin', '配电网关', '配电网关_A-AL-05-2_98CC4D151E0C', '98CC4D151E0C', '20000000000275', 'A栋_五层', md5('A栋_五层'), md5('A栋'), 1),
('98CC4D151F3E', 'admin', 'admin', '配电网关', '配电网关_A-AL-03-2_98CC4D151F3E', '98CC4D151F3E', '20000000000191', 'A栋_三层', md5('A栋_三层'), md5('A栋'), 1),
('98CC4D151F46', 'admin', 'admin', '配电网关', '配电网关_C-AL-01-1_98CC4D151F46', '98CC4D151F46', '20000000000484', 'C栋_一层', md5('C栋_一层'), md5('C栋'), 1),
('98CC4D151D9E', 'admin', 'admin', '配电网关', '配电网关_A-AL-04-1_98CC4D151D9E', '98CC4D151D9E', '20000000000217', 'A栋_四层', md5('A栋_四层'), md5('A栋'), 1),
('98CC4D151D84', 'admin', 'admin', '配电网关', '配电网关_A-AL-06-2_98CC4D151D84', '98CC4D151D84', '20000000000298', 'A栋_六层', md5('A栋_六层'), md5('A栋'), 1),
('98CC4D151C58', 'admin', 'admin', '配电网关', '配电网关_C-AL-02-1_98CC4D151C58', '98CC4D151C58', '20000000000513', 'C栋_二层', md5('C栋_二层'), md5('C栋'), 1),
('98CC4D1528C4', 'admin', 'admin', '配电网关', '配电网关_B-AL-01-1_98CC4D1528C4', '98CC4D1528C4', '20000000000321', 'B栋_一层', md5('B栋_一层'), md5('B栋'), 1),
('98CC4D151D68', 'admin', 'admin', '配电网关', '配电网关_B-AL-01-2_98CC4D151D68', '98CC4D151D68', '20000000000350', 'B栋_一层', md5('B栋_一层'), md5('B栋'), 1),
('98CC4D152AD8', 'admin', 'admin', '配电网关', '配电网关_B-AL-02-1_98CC4D152AD8', '98CC4D152AD8', '20000000000368', 'B栋_二层', md5('B栋_二层'), md5('B栋'), 1),
('98CC4D145C24', 'admin', 'admin', '配电网关', '配电网关_B-AL-02-2_98CC4D145C24', '98CC4D145C24', '20000000000397', 'B栋_二层', md5('B栋_二层'), md5('B栋'), 1),
('98CC4D151E30', 'admin', 'admin', '配电网关', '配电网关_C-AL-02-2_98CC4D151E30', '98CC4D151E30', '20000000000542', 'C栋_二层', md5('C栋_二层'), md5('C栋'), 1),
('98CC4D151DA0', 'admin', 'admin', '配电网关', '配电网关_C-AL-03-1_98CC4D151DA0', '98CC4D151DA0', '20000000000560', 'C栋_三层', md5('C栋_三层'), md5('C栋'), 1),
('98CC4D151D7A', 'admin', 'admin', '配电网关', '配电网关_C-AL-03-2_98CC4D151D7A', '98CC4D151D7A', '20000000000589', 'C栋_三层', md5('C栋_三层'), md5('C栋'), 1),
('98CC4D151DD8', 'admin', 'admin', '配电网关', '配电网关_C-AL-04-1_98CC4D151DD8', '98CC4D151DD8', '20000000000605', 'C栋_四层', md5('C栋_四层'), md5('C栋'), 1),
('98CC4D152922', 'admin', 'admin', '配电网关', '配电网关_C-AL-04-2_98CC4D152922', '98CC4D152922', '20000000000634', 'C栋_四层', md5('C栋_四层'), md5('C栋'), 1),
('98CC4D152986', 'admin', 'admin', '配电网关', '配电网关_C-AL-06-2_98CC4D152986', '98CC4D152986', '20000000000643', 'C栋_六层', md5('C栋_六层'), md5('C栋'), 1),
('98CC4D1528E4', 'admin', 'admin', '配电网关', '配电网关_E-AP-02_98CC4D1528E4', '98CC4D1528E4', '20000000000662', 'E栋_二层', md5('E栋_二层'), md5('E栋'), 2),
('98CC4D150E66', 'admin', 'admin', '配电网关', '配电网关_E-AL-01_98CC4D150E66', '98CC4D150E66', '20000000000661', 'E栋_一层', md5('E栋_一层'), md5('E栋'), 1),
('98CC4D150A3C', 'admin', 'admin', '配电网关', '配电网关_E-AL-02_98CC4D150A3C', '98CC4D150A3C', '20000000000657', 'E栋_二层', md5('E栋_二层'), md5('E栋'), 1),
('98CC4D149A06', 'admin', 'admin', '配电网关', '配电网关_C-AL-01_98CC4D149A06', '98CC4D149A06', '20000000000660', 'C栋_一层', md5('C栋_一层'), md5('C栋'), 1),
('98CC4D151E00', 'admin', 'admin', '配电网关', '配电网关_C-AP-02_98CC4D151E00', '98CC4D151E00', '20000000000664', 'C栋_二层', md5('C栋_二层'), md5('C栋'), 2),
('98CC4D152928', 'admin', 'admin', '配电网关', '配电网关_A-AL-03_98CC4D152928', '98CC4D152928', '20000000000663', 'A栋_三层', md5('A栋_三层'), md5('A栋'), 1),
('98CC4D151E04', 'admin', 'admin', '配电网关', '配电网关_A-AP-02_98CC4D151E04', '98CC4D151E04', '20000000000659', 'A栋_二层', md5('A栋_二层'), md5('A栋'), 2),
('98CC4D152AC8', 'admin', 'admin', '配电网关', '配电网关_B-AL-02_98CC4D152AC8', '98CC4D152AC8', '20000000000656', 'B栋_二层', md5('B栋_二层'), md5('B栋'), 1),
('98CC4D152AE4', 'admin', 'admin', '配电网关', '配电网关_A-AP-01_98CC4D152AE4', '98CC4D152AE4', '20000000000654', 'A栋_一层', md5('A栋_一层'), md5('A栋'), 2),
('98CC4D151DF2', 'admin', 'admin', '配电网关', '配电网关_B-AP-02_98CC4D151DF2', '98CC4D151DF2', '20000000000653', 'B栋_二层', md5('B栋_二层'), md5('B栋'), 2),
('98CC4D1528E2', 'admin', 'admin', '配电网关', '配电网关_A-AL-02_98CC4D1528E2', '98CC4D1528E2', '20000000000665', 'A栋_二层', md5('A栋_二层'), md5('A栋'), 1),
('98CC4D151D3A', 'admin', 'admin', '配电网关', '配电网关_C-AL-02_98CC4D151D3A', '98CC4D151D3A', '20000000000658', 'C栋_二层', md5('C栋_二层'), md5('C栋'), 1),
('98CC4D150E1A', 'admin', 'admin', '配电网关', '配电网关_B-AP-01_98CC4D150E1A', '98CC4D150E1A', '20000000000652', 'B栋_一层', md5('B栋_一层'), md5('B栋'), 2),
('98CC4D151C56', 'admin', 'admin', '配电网关', '配电网关_C-AP-01_98CC4D151C56', '98CC4D151C56', '20000000000666', 'C栋_一层', md5('C栋_一层'), md5('C栋'), 2),
('98CC4D152990', 'admin', 'admin', '配电网关', '配电网关_B-AL-01_98CC4D152990', '98CC4D152990', '20000000000655', 'B栋_一层', md5('B栋_一层'), md5('B栋'), 1),
('98CC4D150A00', 'admin', 'admin', '配电网关', '配电网关_A-AL-01_98CC4D150A00', '98CC4D150A00', '20000000000668', 'A栋_一层', md5('A栋_一层'), md5('A栋'), 1),
('98CC4D149A0C', 'admin', 'admin', '配电网关', '配电网关_E-AP-01_98CC4D149A0C', '98CC4D149A0C', '20000000000667', 'E栋_一层', md5('E栋_一层'), md5('E栋'), 2),
('98CC4D15209A', 'admin', 'admin', '配电网关', '配电网关_C-AL-07-1_98CC4D15209A', '98CC4D15209A', '20000000000703', 'C栋_七层', md5('C栋_七层'), md5('C栋'), 1),
('98CC4D15236A', 'admin', 'admin', '配电网关', '配电网关_C-AL-07-2_98CC4D15236A', '98CC4D15236A', '20000000000704', 'C栋_七层', md5('C栋_七层'), md5('C栋'), 1),
('98CC4D149A0A', 'admin', 'admin', '配电网关', '配电网关_A-AL-06-1_98CC4D149A0A', '98CC4D149A0A', '20000000000705', 'A栋_六层', md5('A栋_六层'), md5('A栋'), 1),
('98CC4D1524AA', 'admin', 'admin', '配电网关', '配电网关_A-AL-04-2_98CC4D1524AA', '98CC4D1524AA', '20000000000706', 'A栋_四层', md5('A栋_四层'), md5('A栋'), 1),
('98CC4D15298C', 'admin', 'admin', '配电网关', '配电网关_B-AL-03-1_98CC4D15298C', '98CC4D15298C', '20000000000707', 'B栋_三层', md5('B栋_三层'), md5('B栋'), 1),
('98CC4D150BD4', 'admin', 'admin', '配电网关', '配电网关_A-AL-02-1_98CC4D150BD4', '98CC4D150BD4', '20000000000708', 'A栋_二层', md5('A栋_二层'), md5('A栋'), 1),
('98CC4D151C54', 'admin', 'admin', '配电网关', '配电网关_B-AL-05-1_98CC4D151C54', '98CC4D151C54', '20000000000709', 'B栋_五层', md5('B栋_五层'), md5('B栋'), 1);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS v_eco_building_info_1y;
DROP VIEW IF EXISTS v_eco_building_info_1m;
DROP VIEW IF EXISTS v_eco_building_info_1d;
DROP FUNCTION IF EXISTS generate_gateway_test_data;
DROP MATERIALIZED VIEW IF EXISTS f_eco_building_1y;
DROP MATERIALIZED VIEW IF EXISTS f_eco_building_1m;
DROP MATERIALIZED VIEW IF EXISTS f_eco_building_1d;
DROP MATERIALIZED VIEW IF EXISTS f_eco_floor_1y;
DROP MATERIALIZED VIEW IF EXISTS f_eco_floor_1m;
DROP MATERIALIZED VIEW IF EXISTS f_eco_floor_1d;
DROP MATERIALIZED VIEW IF EXISTS f_eco_gateway_1y;
DROP MATERIALIZED VIEW IF EXISTS f_eco_gateway_1m;
DROP MATERIALIZED VIEW IF EXISTS f_eco_gateway_1d;
DROP TABLE IF EXISTS f_eco_gateway_1h;
DROP TABLE IF EXISTS o_eco_gateway;
DROP TABLE IF EXISTS o_eco_floor;
DROP TABLE IF EXISTS o_eco_building;
-- +goose StatementEnd
