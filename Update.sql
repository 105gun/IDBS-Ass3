UPDATE user
SET authority=0
WHERE id IN (SELECT uid
FROM book, booktype, borrow
WHERE bid=book.id AND book.ISBN=booktype.ISBN AND removed=0 AND existed=1 AND is_returned=0 AND time+(1 + extend_status)*7<17)
