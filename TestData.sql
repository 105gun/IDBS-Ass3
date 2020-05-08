INSERT INTO user
VALUES	(1, 1),
	(2, 1),
	(3, 0);

INSERT INTO booktype
VALUES	('T1', 'Nintendo Power', 'Nintendo'),
	('T2', 'Ninja Slayer', 'Ninja'),
	('T3', 'You Know WHO', 'WHO');

INSERT INTO book
VALUES	(1, 'T1', 1, 0, ''),
	(2, 'T1', 0, 1, 'MISSED'),
	(3, 'T2', 1, 0, ''),
	(4, 'T3', 1, 0, ''),
	(5, 'T3', 1, 0, ''),
	(6, 'T3', 1, 0, ''),
	(7, 'T3', 1, 0, '');

INSERT INTO borrow
VALUES	(1, 1, 1, 0, 1, 0),
	(2, 3, 1, 3, 0, 0),
	(3, 4, 1, 0, 3, 0),
	(4, 5, 2, 0, 2, 0),
	(5, 6, 2, 0, 1, 0);
