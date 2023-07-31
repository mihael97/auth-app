INSERT INTO
    auth.CUSTOMER_ROLES(ROLE_NAME, USER_ID)
select 'USER',u.id from auth.users u WHERE u.id NOT IN (SELECT x.USER_ID FROM auth.CUSTOMER_ROLES x WHERE x.ROLE_NAME = 'USER');