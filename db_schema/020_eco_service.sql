-- +goose Up
-- +goose StatementBegin

CREATE TABLE o_eco_gateway (
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
    floor_no VARCHAR(128) NOT NULL,
    building_no VARCHAR(128) NOT NULL
);
ALTER TABLE o_eco_gateway ADD CONSTRAINT uk_eco_gateway_channel_no UNIQUE(channel_no);
COMMENT ON TABLE o_eco_gateway IS '配电网关信息表';
COMMENT ON COLUMN o_eco_gateway.id IS '主键ID';
COMMENT ON COLUMN o_eco_gateway.model_name IS '型号名称';
COMMENT ON COLUMN o_eco_gateway.dev_name IS '设备名称';
COMMENT ON COLUMN o_eco_gateway.channel_no IS '通道号';
COMMENT ON COLUMN o_eco_gateway.cm_code IS '通信码';
COMMENT ON COLUMN o_eco_gateway.location IS '组织名称';
COMMENT ON COLUMN o_eco_gateway.created_time IS '创建时间';
COMMENT ON COLUMN o_eco_gateway.updated_time IS '更新时间';
COMMENT ON COLUMN o_eco_gateway.created_by IS '创建人';
COMMENT ON COLUMN o_eco_gateway.updated_by IS '更新人';
COMMENT ON COLUMN o_eco_gateway.floor_no IS '楼层号';
COMMENT ON COLUMN o_eco_gateway.building_no IS '楼栋号';

INSERT INTO o_eco_gateway (id, created_by, updated_by, model_name, dev_name, channel_no, cm_code, location, floor_no, building_no)
VALUES
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_B-AL-04-1_98CC4D151E44', '98CC4D151E44', '20000000000416', 'B栋_四层', '四层', 'B栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_B-AL-04-2_98CC4D152976', '98CC4D152976', '20000000000445', 'B栋_四层', '四层', 'B栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_B-AL-05-2_98CC4D152988', '98CC4D152988', '20000000000458', 'B栋_五层', '五层', 'B栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_A-AL-05-1_98CC4D151D8C', '98CC4D151D8C', '20000000000246', 'A栋_五层', '五层', 'A栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_B-AL-06-2_98CC4D151F92', '98CC4D151F92', '20000000000471', 'B栋_六层', '六层', 'B栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_A-AL-01-3_98CC4D150BC8', '98CC4D150BC8', '20000000000096', 'A栋_一层', '一层', 'A栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_A-AL-01-2_98CC4D151DF8', '98CC4D151DF8', '20000000000112', 'A栋_一层', '一层', 'A栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_A-AL-02-2_98CC4D1528BA', '98CC4D1528BA', '20000000000133', 'A栋_二层', '二层', 'A栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_A-AL-03-1_98CC4D151DD0', '98CC4D151DD0', '20000000000162', 'A栋_三层', '三层', 'A栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_A-AL-05-2_98CC4D151E0C', '98CC4D151E0C', '20000000000275', 'A栋_五层', '五层', 'A栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_A-AL-03-2_98CC4D151F3E', '98CC4D151F3E', '20000000000191', 'A栋_三层', '三层', 'A栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_C-AL-01-1_98CC4D151F46', '98CC4D151F46', '20000000000484', 'C栋_一层', '一层', 'C栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_A-AL-04-1_98CC4D151D9E', '98CC4D151D9E', '20000000000217', 'A栋_四层', '四层', 'A栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_A-AL-06-2_98CC4D151D84', '98CC4D151D84', '20000000000298', 'A栋_六层', '六层', 'A栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_C-AL-02-1_98CC4D151C58', '98CC4D151C58', '20000000000513', 'C栋_二层', '二层', 'C栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_B-AL-01-1_98CC4D1528C4', '98CC4D1528C4', '20000000000321', 'B栋_一层', '一层', 'B栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_B-AL-01-2_98CC4D151D68', '98CC4D151D68', '20000000000350', 'B栋_一层', '一层', 'B栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_B-AL-02-1_98CC4D152AD8', '98CC4D152AD8', '20000000000368', 'B栋_二层', '二层', 'B栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_B-AL-02-2_98CC4D145C24', '98CC4D145C24', '20000000000397', 'B栋_二层', '二层', 'B栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_C-AL-02-2_98CC4D151E30', '98CC4D151E30', '20000000000542', 'C栋_二层', '二层', 'C栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_C-AL-03-1_98CC4D151DA0', '98CC4D151DA0', '20000000000560', 'C栋_三层', '三层', 'C栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_C-AL-03-2_98CC4D151D7A', '98CC4D151D7A', '20000000000589', 'C栋_三层', '三层', 'C栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_C-AL-04-1_98CC4D151DD8', '98CC4D151DD8', '20000000000605', 'C栋_四层', '四层', 'C栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_C-AL-04-2_98CC4D152922', '98CC4D152922', '20000000000634', 'C栋_四层', '四层', 'C栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_C-AL-06-2_98CC4D152986', '98CC4D152986', '20000000000643', 'C栋_六层', '六层', 'C栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_E-AP-02_98CC4D1528E4', '98CC4D1528E4', '20000000000662', 'E栋_二层', '二层', 'E栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_E-AL-01_98CC4D150E66', '98CC4D150E66', '20000000000661', 'E栋_一层', '一层', 'E栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_E-AL-02_98CC4D150A3C', '98CC4D150A3C', '20000000000657', 'E栋_二层', '二层', 'E栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_C-AL-01_98CC4D149A06', '98CC4D149A06', '20000000000660', 'C栋_一层', '一层', 'C栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_C-AP-02_98CC4D151E00', '98CC4D151E00', '20000000000664', 'C栋_二层', '二层', 'C栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_A-AL-03_98CC4D152928', '98CC4D152928', '20000000000663', 'A栋_三层', '三层', 'A栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_A-AP-02_98CC4D151E04', '98CC4D151E04', '20000000000659', 'A栋_二层', '二层', 'A栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_B-AL-02_98CC4D152AC8', '98CC4D152AC8', '20000000000656', 'B栋_二层', '二层', 'B栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_A-AP-01_98CC4D152AE4', '98CC4D152AE4', '20000000000654', 'A栋_一层', '一层', 'A栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_B-AP-02_98CC4D151DF2', '98CC4D151DF2', '20000000000653', 'B栋_二层', '二层', 'B栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_A-AL-02_98CC4D1528E2', '98CC4D1528E2', '20000000000665', 'A栋_二层', '二层', 'A栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_C-AL-02_98CC4D151D3A', '98CC4D151D3A', '20000000000658', 'C栋_二层', '二层', 'C栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_B-AP-01_98CC4D150E1A', '98CC4D150E1A', '20000000000652', 'B栋_一层', '一层', 'B栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_C-AP-01_98CC4D151C56', '98CC4D151C56', '20000000000666', 'C栋_一层', '一层', 'C栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_B-AL-01_98CC4D152990', '98CC4D152990', '20000000000655', 'B栋_一层', '一层', 'B栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_A-AL-01_98CC4D150A00', '98CC4D150A00', '20000000000668', 'A栋_一层', '一层', 'A栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_E-AP-01_98CC4D149A0C', '98CC4D149A0C', '20000000000667', 'E栋_一层', '一层', 'E栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_C-AL-07-1_98CC4D15209A', '98CC4D15209A', '20000000000703', 'C栋_七层', '七层', 'C栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_C-AL-07-2_98CC4D15236A', '98CC4D15236A', '20000000000704', 'C栋_七层', '七层', 'C栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_A-AL-06-1_98CC4D149A0A', '98CC4D149A0A', '20000000000705', 'A栋_六层', '六层', 'A栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_A-AL-04-2_98CC4D1524AA', '98CC4D1524AA', '20000000000706', 'A栋_四层', '四层', 'A栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_B-AL-03-1_98CC4D15298C', '98CC4D15298C', '20000000000707', 'B栋_三层', '三层', 'B栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_A-AL-02-1_98CC4D150BD4', '98CC4D150BD4', '20000000000708', 'A栋_二层', '二层', 'A栋'),
(select nanoid(), 'admin', 'admin', '配电网关', '配电网关_B-AL-05-1_98CC4D151C54', '98CC4D151C54', '20000000000709', 'B栋_五层', '五层', 'B栋');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS o_eco_gateway;

-- +goose StatementEnd
