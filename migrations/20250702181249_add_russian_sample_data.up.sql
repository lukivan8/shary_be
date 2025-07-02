-- Add Russian sample data with tenge currency

-- Insert sample users with Kazakhstan IIN
INSERT INTO users (first_name, last_name, identity, phone, avatar_url, verified) VALUES
('Александр', 'Иванов', '123456789012', '+7 701 123 4567', 'https://example.com/avatars/alex.jpg', true),
('Мария', 'Петрова', '234567890123', '+7 702 234 5678', 'https://example.com/avatars/maria.jpg', true),
('Дмитрий', 'Сидоров', '345678901234', '+7 703 345 6789', 'https://example.com/avatars/dmitry.jpg', false),
('Анна', 'Козлова', '456789012345', '+7 704 456 7890', 'https://example.com/avatars/anna.jpg', true),
('Сергей', 'Новиков', '567890123456', '+7 705 567 8901', 'https://example.com/avatars/sergey.jpg', false);

-- Insert sample categories in Russian
INSERT INTO categories (name) VALUES
('Электроника'),
('Одежда и обувь'),
('Спорт и отдых'),
('Книги и образование'),
('Дом и сад'),
('Транспорт'),
('Красота и здоровье'),
('Игрушки и хобби'),
('Другое');

-- Insert sample statuses (internal names in English)
INSERT INTO statuses (name) VALUES
('requested'),
('approved'),
('rejected'),
('recieved'),
('finished');

-- Insert sample items with Russian titles and tenge prices
INSERT INTO items (title, description, image_url, price, author_id, location, category_id, tags, has_photos) VALUES
('iPhone 14 Pro', 'Современный смартфон Apple в отличном состоянии. Идеально подходит для работы и развлечений.', 'https://example.com/items/iphone14.jpg', 50000, 1, 'Алматы, ул. Достык 123', 1, ARRAY['смартфон', 'apple', 'техника'], true),
('Велосипед горный', 'Качественный горный велосипед для активного отдыха. Подходит для катания по городу и за городом.', 'https://example.com/items/bike.jpg', 15000, 2, 'Астана, пр. Республики 45', 3, ARRAY['велосипед', 'спорт', 'отдых'], true),
('Книга "Война и мир"', 'Классическое произведение Льва Толстого в твердом переплете. Отличное состояние.', 'https://example.com/items/war_peace.jpg', 2000, 3, 'Алматы, ул. Абая 67', 4, ARRAY['книга', 'классика', 'литература'], false),
('Платье вечернее', 'Элегантное вечернее платье черного цвета. Размер M, идеально для торжественных мероприятий.', 'https://example.com/items/dress.jpg', 8000, 4, 'Астана, ул. Бейбитшилик 89', 2, ARRAY['платье', 'вечернее', 'элегантное'], true),
('Ноутбук Dell', 'Мощный ноутбук для работы и учебы. Intel i7, 16GB RAM, SSD 512GB.', 'https://example.com/items/laptop.jpg', 120000, 1, 'Алматы, ул. Толе би 234', 1, ARRAY['ноутбук', 'dell', 'компьютер'], true),
('Футбольный мяч', 'Профессиональный футбольный мяч для тренировок и игр. Размер 5.', 'https://example.com/items/football.jpg', 3000, 5, 'Астана, ул. Кенесары 12', 3, ARRAY['футбол', 'мяч', 'спорт'], false),
('Кофемашина', 'Автоматическая кофемашина для приготовления эспрессо и капучино. В отличном состоянии.', 'https://example.com/items/coffee.jpg', 25000, 2, 'Алматы, ул. Фурманова 56', 5, ARRAY['кофе', 'машина', 'кухня'], true),
('Гитара акустическая', 'Красивая акустическая гитара для начинающих музыкантов. Включает чехол и медиатор.', 'https://example.com/items/guitar.jpg', 18000, 3, 'Астана, ул. Сарыарка 78', 8, ARRAY['гитара', 'музыка', 'инструмент'], true);

-- Insert sample item photos
INSERT INTO item_photos (item_id, url) VALUES
(1, 'https://example.com/photos/iphone1.jpg'),
(1, 'https://example.com/photos/iphone2.jpg'),
(2, 'https://example.com/photos/bike1.jpg'),
(2, 'https://example.com/photos/bike2.jpg'),
(4, 'https://example.com/photos/dress1.jpg'),
(5, 'https://example.com/photos/laptop1.jpg'),
(5, 'https://example.com/photos/laptop2.jpg'),
(7, 'https://example.com/photos/coffee1.jpg'),
(8, 'https://example.com/photos/guitar1.jpg'),
(8, 'https://example.com/photos/guitar2.jpg');

-- Insert sample rents
INSERT INTO rents (item_id, renter_id, date_start, date_end, price, status_id) VALUES
(1, 2, '2024-07-01', '2024-07-03', 15000, 2),
(2, 1, '2024-07-05', '2024-07-07', 4500, 1),
(5, 4, '2024-07-10', '2024-07-15', 36000, 1),
(7, 3, '2024-07-12', '2024-07-14', 7500, 4);

-- Insert sample favorite items
INSERT INTO favorite_items (item_id, user_id) VALUES
(1, 2),
(1, 4),
(2, 1),
(2, 5),
(5, 3),
(7, 1),
(8, 2); 