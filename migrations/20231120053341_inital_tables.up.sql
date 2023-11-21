BEGIN;

CREATE TABLE IF NOT EXISTS "product" (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    count INT NOT NULL,
    deleted BOOLEAN DEFAULT FALSE
);

INSERT INTO product (name, description, price, count) VALUES
        ('Смартфон Samsung Galaxy S121', 'Современный смартфон с высокой производительностью.', 799.99, 50),
        ('Ноутбук HP Envy 13', 'Легкий и мощный ноутбук для работы и развлечений.', 1099.99, 30),
        ('Фотокамера Canon EOS 90D', 'Профессиональная зеркальная фотокамера с 32 Мп сенсором.', 1299.99, 20),
        ('Телевизор LG OLED65CXPUA', '65-дюймовый OLED-телевизор с 4K разрешением и HDR.', 1799.99, 10),
        ('Наушники Sony WH-1000XM4', 'Беспроводные наушники с шумоподавлением.', 349.99, 40),
        ('Электрический чайник Bosch TWK7203', 'Мощный чайник с регулируемой температурой.', 59.99, 100),
        ('Кофеварка Philips 3200 Series', 'Автоматическая кофеварка с функцией капучино.', 349.99, 30),
        ('Пылесос Dyson V11 Absolute', 'Беспроводной пылесос с высокой мощностью всасывания.', 499.99, 15),
        ('Умные часы Apple Watch Series 6', 'Смарт-часы с измерением уровня кислорода в крови.', 399.99, 25),
        ('Игровая консоль PlayStation 5', 'Мощная игровая консоль нового поколения.', 499.99, 10),
        ('Холодильник Samsung RF28R7201SR', 'Большой холодильник с функцией автоматической очистки.', 1499.99, 5),
        ('Стиральная машина LG WM4000HWA', 'Стиральная машина с большой загрузкой и инверторным двигателем.', 799.99, 15),
        ('Микроволновая печь Panasonic NN-SU696S', 'Микроволновая печь с грилем и функцией размораживания.', 119.99, 20),
        ('Пылесос-робот iRobot Roomba i7+', 'Робот-пылесос с навигацией и самоочисткой контейнера.', 699.99, 8),
        ('Автомобиль Toyota Camry', 'Седан с бензиновым двигателем и автоматической коробкой передач.', 24999.99, 3);

COMMIT;


BEGIN;

CREATE TABLE IF NOT EXISTS "user" (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    phone VARCHAR(13) NOT NULL,
    login VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50),
    address TEXT,
    coordinate_address_x DECIMAL(10, 6),
    coordnates_address_y DECIMAL(10, 6),
    deleted BOOLEAN DEFAULT FALSE
);

INSERT INTO "user" (name, surname, phone, login, password_hash, role, address, coordinate_address_x, coordnates_address_y) VALUES
        ('user1', 'perviy', '+998805553535', 'user1', 'nothash', 'user', 'prospekt mira', 123.235, 112.236),
        ('user2', 'vtoroy', '+998805553535', 'user2', 'nothash', 'user', 'prospekt nemira', 223.11, 122.00)
;

CREATE TABLE IF NOT EXISTS "order" (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES "user" (id),
    address TEXT,
    coordinate_address_x DECIMAL(10, 6) NOT NULL,
    coordinate_address_y DECIMAL(10, 6) NOT NULL,
    coordinate_point_x DECIMAL(10, 6) NOT NULL,
    coordinate_point_y DECIMAL(10, 6) NOT NULL,
    create_at TIMESTAMP,
    start_at TIMESTAMP,
    delivery_at TIMESTAMP,
    courier_id INT,
    delivery_status VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS "order_product" (
    id SERIAL PRIMARY KEY,
    order_id INT,
    product_id INT REFERENCES "product" (id),
	count INT NOT NULL,

	FOREIGN KEY (order_id) REFERENCES "order" (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "basket" (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES "user" (id),
    product_id INT REFERENCES "product" (id),
    count INT NOT NULL
);

CREATE TABLE IF NOT EXISTS "delivery" (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES "user" (id),
    order_id INT REFERENCES "order" (id)
);

COMMIT;

-- migrate -path migrations -database "postgresql://delivery:delivery@localhost:5432/delivery?sslmode=disable" up