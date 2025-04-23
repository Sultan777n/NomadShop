DO $$
DECLARE
u RECORD;
BEGIN
SELECT id, username, email INTO u
FROM users
WHERE email = 'asqar@mail.com'
    LIMIT 1;

RAISE NOTICE 'Қолданушы өшірілді: id=%, username=%, email=%', u.id, u.username, u.email;
END $$;

DELETE FROM users
WHERE email = 'asqar@mail.com';

