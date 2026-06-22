-- Seed data for initial development and demo environments.
-- The seed tool replaces __ADMIN_PASSWORD_HASH__ with a bcrypt hash before executing.

INSERT INTO users (id, username, password, name, role, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'admin',
    '__ADMIN_PASSWORD_HASH__',
    'Administrator',
    'admin',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT DO NOTHING;

INSERT INTO categories (id, name, sort_order, created_at, updated_at)
VALUES
    ('10000000-0000-0000-0000-000000000001', '咖啡', 10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('10000000-0000-0000-0000-000000000002', '茶饮', 20, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('10000000-0000-0000-0000-000000000003', '甜点', 30, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

INSERT INTO products (id, category_id, name, description, price, image, status, sort_order, created_at, updated_at)
VALUES
    (
        '20000000-0000-0000-0000-000000000001',
        '10000000-0000-0000-0000-000000000001',
        '拿铁咖啡',
        '浓缩咖啡搭配丝滑蒸奶，口感醇厚。',
        28.00,
        '/images/products/latte.jpg',
        'active',
        10,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '20000000-0000-0000-0000-000000000002',
        '10000000-0000-0000-0000-000000000001',
        '美式咖啡',
        '经典黑咖啡，清爽顺口，突出咖啡原香。',
        22.00,
        '/images/products/americano.jpg',
        'active',
        20,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '20000000-0000-0000-0000-000000000003',
        '10000000-0000-0000-0000-000000000002',
        '茉莉绿茶',
        '清新茉莉花香与绿茶茶汤融合，回甘自然。',
        18.00,
        '/images/products/jasmine-green-tea.jpg',
        'active',
        10,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '20000000-0000-0000-0000-000000000004',
        '10000000-0000-0000-0000-000000000002',
        '蜜桃乌龙茶',
        '乌龙茶底搭配蜜桃果香，清甜不腻。',
        24.00,
        '/images/products/peach-oolong.jpg',
        'active',
        20,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '20000000-0000-0000-0000-000000000005',
        '10000000-0000-0000-0000-000000000003',
        '提拉米苏',
        '咖啡酒香与马斯卡彭奶酪交织的经典甜点。',
        32.00,
        '/images/products/tiramisu.jpg',
        'active',
        10,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '20000000-0000-0000-0000-000000000006',
        '10000000-0000-0000-0000-000000000003',
        '芝士蛋糕',
        '绵密芝士口感，奶香浓郁，甜度适中。',
        30.00,
        '/images/products/cheesecake.jpg',
        'active',
        20,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    )
ON CONFLICT DO NOTHING;
