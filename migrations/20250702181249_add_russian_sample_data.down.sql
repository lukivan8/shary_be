-- Remove Russian sample data

-- Remove sample favorite items
DELETE FROM favorite_items WHERE item_id IN (1, 2, 5, 7, 8) AND user_id IN (1, 2, 3, 4, 5);

-- Remove sample rents
DELETE FROM rents WHERE item_id IN (1, 2, 5, 7) AND renter_id IN (1, 2, 3, 4);

-- Remove sample item photos
DELETE FROM item_photos WHERE item_id IN (1, 2, 4, 5, 7, 8);

-- Remove sample items
DELETE FROM items WHERE id IN (1, 2, 3, 4, 5, 6, 7, 8);

-- Remove sample statuses
DELETE FROM statuses WHERE name IN ('requested', 'approved', 'rejected', 'recieved', 'finished');

-- Remove sample categories
DELETE FROM categories WHERE name IN ('Электроника', 'Одежда и обувь', 'Спорт и отдых', 'Книги и образование', 'Дом и сад', 'Транспорт', 'Красота и здоровье', 'Игрушки и хобби', 'Другое');

-- Remove sample users
DELETE FROM users WHERE identity IN ('123456789012', '234567890123', '345678901234', '456789012345', '567890123456'); 