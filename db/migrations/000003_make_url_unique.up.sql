alter table URLS
add constraint UNQ_URLS_ORIGINAL_URL unique (URLS_ORIGINAL_URL);