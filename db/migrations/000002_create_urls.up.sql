create table URLS(
   URLS_ID serial primary key,
   URLS_ORIGINAL_URL text not null,
   USERS_ID bigint not null references USERS(USERS_ID)
);