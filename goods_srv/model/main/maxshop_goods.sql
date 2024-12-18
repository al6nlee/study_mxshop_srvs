#先关闭外键检查
SET FOREIGN_KEY_CHECKS = 0;

#drop
TABLE banner;
#drop
table goods;
#drop
TABLE goodscategorybrand;
#drop
TABLE category;
#drop
TABLE brands;

insert into brands (`id`, `name`, `is_deleted`)
values (1, "顶端", 0),
       (2, "nba", 0),
       (3, "文果", 0),
       (4, "寻真", 0),
       (5, "轻恋", 0),
       (6, "木石", 0),
       (7, "马小二", 0),
       (8, "金山湾", 0),
       (9, "日日顺", 0),
       (10, "华为", 0),
       (11, "小米", 0),
       (12, "Apple", 0);

INSERT INTO goodscategorybrand (`category_id`, `brands_id`)
VALUES (2001, 1),
       (2002, 2),
       (2003, 3),
       (2004, 4),
       (2005, 5),
       (2006, 6),
       (2007, 7),
       (2008, 8),
       (2009, 9);

INSERT INTO category (`id`, `name`, `is_deleted`, `parent_category_id`, `LEVEL`)
VALUES (1001, "新鲜水果", 0, 0, 1),
       (1002, "新鲜蔬菜", 0, 0, 1),
       (1003, "化妆品", 0, 0, 1),

       (2001, "苹果", 0, 1001, 2),
       (2002, "香蕉", 0, 1001, 2),
       (2003, "凤梨", 0, 1001, 2),
       (2004, "空心菜", 0, 1002, 2),
       (2005, "黄花菜", 0, 1002, 2),
       (2006, "白菜", 0, 1002, 2),
       (2007, "口红", 0, 1003, 2),
       (2008, "香水", 0, 1003, 2),
       (2009, "护肤品", 0, 1003, 2),

       (3001, "红富士", 0, 2001, 3),
       (3002, "华圣果业", 0, 2001, 3),
       (3003, "潘苹果", 0, 2001, 3),
       (3004, "佳农", 0, 2002, 3),
       (3005, "都乐", 0, 2002, 3),
       (3006, "果迎鲜", 0, 2002, 3),
       (3007, "甘福园", 0, 2003, 3),
       (3008, "百果园", 0, 2003, 3),
       (3009, "鲜蜂堆", 0, 2003, 3),


       (3011, "沿海", 0, 2004, 3),
       (3012, "南方", 0, 2004, 3),
       (3013, "北方", 0, 2004, 3),
       (3014, "绿色", 0, 2005, 3),
       (3015, "蓝色", 0, 2005, 3),
       (3016, "黄色", 0, 2005, 3),
       (3017, "冬季", 0, 2006, 3),
       (3018, "夏季", 0, 2006, 3),
       (3019, "春季", 0, 2006, 3),

       (3021, "Dior", 0, 2007, 3),
       (3022, "香奈儿", 0, 2007, 3),
       (3023, "YSL", 0, 2007, 3),
       (3024, "让巴杜", 0, 2008, 3),
       (3025, "恩加罗", 0, 2008, 3),
       (3026, "龙芳", 0, 2008, 3),
       (3027, "相宜本草", 0, 2009, 3),
       (3028, "百雀羚", 0, 2009, 3),
       (3029, "妮维雅", 0, 2009, 3);

INSERT INTO goods (`category_id`, `brands_id`, `is_new`, `is_hot`, `name`, `shop_price`, `market_price`, `goods_sn`,
                   `goods_brief`, `images`, `desc_images`, `goods_front_image`)
VALUES (1001, 1, 0, 0, "芒果", 12.3, 20.1, 1, '简介1',
        '["https://example.com/mango1.jpg", "https://example.com/mango2.jpg"]',
        '["https://example.com/mango_detail1.jpg", "https://example.com/mango_detail2.jpg"]',
        "https://example.com/mango_detail2.jpg"),
       (1002, 5, 0, 0, "大白菜", 6.1, 11.2, 2, '简介2',
        '["https://example.com/cabbage1.jpg", "https://example.com/cabbage2.jpg"]',
        '["https://example.com/cabbage_detail1.jpg", "https://example.com/cabbage_detail2.jpg"]',
        "https://example.com/mango_detail2.jpg"),
       (1003, 3, 0, 0, "卸妆水", 58.8, 98.9, 3, '简介3',
        '["https://example.com/makeupremover1.jpg", "https://example.com/makeupremover2.jpg"]',
        '["https://example.com/makeupremover_detail1.jpg", "https://example.com/makeupremover_detail2.jpg"]',
        "https://example.com/mango_detail2.jpg"),
       (2001, 9, 0, 0, "香山苹果", 21.2, 33.7, 4, '简介4',
        '["https://example.com/apple1.jpg", "https://example.com/apple2.jpg"]',
        '["https://example.com/apple_detail1.jpg", "https://example.com/apple_detail2.jpg"]',
        "https://example.com/mango_detail2.jpg"),
       (2002, 8, 0, 0, "巴西香蕉", 6.7, 10.3, 5, '简介5',
        '["https://example.com/banana1.jpg", "https://example.com/banana2.jpg"]',
        '["https://example.com/banana_detail1.jpg", "https://example.com/banana_detail2.jpg"]',
        "https://example.com/mango_detail2.jpg"),
       (2003, 6, 0, 0, "台湾凤梨", 28.7, 39.2, 6, '简介6',
        '["https://example.com/pineapple1.jpg", "https://example.com/pineapple2.jpg"]',
        '["https://example.com/pineapple_detail1.jpg", "https://example.com/pineapple_detail2.jpg"]',
        "https://example.com/mango_detail2.jpg"),
       (2009, 3, 0, 0, "粉底液", 158.8, 198.9, 7, '简介7',
        '["https://example.com/foundation1.jpg", "https://example.com/foundation2.jpg"]',
        '["https://example.com/foundation_detail1.jpg", "https://example.com/foundation_detail2.jpg"]',
        "https://example.com/mango_detail2.jpg"),
       (2007, 3, 0, 0, "迪奥口红", 699.9, 1168.8, 8, '简介8',
        '["https://example.com/lipstick1.jpg", "https://example.com/lipstick2.jpg"]',
        '["https://example.com/lipstick_detail1.jpg", "https://example.com/lipstick_detail2.jpg"]',
        "https://example.com/mango_detail2.jpg");

 #插入数据后开启外键检查
SET FOREIGN_KEY_CHECKS = 1;