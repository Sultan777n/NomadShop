DO $$
    DECLARE
        u RECORD;
    BEGIN

        INSERT INTO users (username, email, password)
        VALUES ('asqar', 'asqar@mail.com', 'asqar12345')
        ON CONFLICT (username) DO UPDATE
            SET username = 'Asqar'
        RETURNING id, username, email INTO u;

        RAISE NOTICE 'Қолданушы қосылды/жаңартылды: id=%, username=%, email=%', u.id, u.username, u.email;
    END $$;
