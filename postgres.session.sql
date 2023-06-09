SELECT count(*) FROM userswithjob;
Select * from userswithjob where id between 1 and 2000 Order by id asc;
SELECT * FROM userswithjob WHERE id = 4 and role::varchar LIKE '%Frontend%';
Select * from userswithjob where role='Frontend' AND role::varchar LIKE '%Frontend%'
SELECT '5'::json ;
SELECT role ->> 'role' as role from userswithjob;
SELECT role->>'[0]' from userswithjob;
SELECT * from userswithjob where id between 01 and 10;
insert into userswithjob (id, first_name, last_name, role, title) values (1001, 'Arliene', 'Cotterell','["Frontend"]' , 'Nurse Practicioner');
Select * from userswithjob ;
SELECT
    *
FROM
    userswithjob
ORDER BY
    id 
Limit 10 offset 6;