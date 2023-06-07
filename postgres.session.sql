SELECT * FROM userswithjob;
SELECT * FROM userswithjob WHERE id = 4 and role::varchar like '%Frontend%';
Select * from userswithjob where role='Frontend' AND role::varchar like '%Frontend%'
SELECT '5'::json ;
SELECT role ->> 'role' as role from userswithjob;
SELECT role->>'[0]' from userswithjob;
SELECT * from userswithjob WHERE role;
insert into userswithjob (id, first_name, last_name, role, title) values (1001, 'Arliene', 'Cotterell','["Frontend"]' , 'Nurse Practicioner');