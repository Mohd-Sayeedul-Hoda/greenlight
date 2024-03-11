ALTER TABLE movies ADD CONSTRAINT movies_runtime_check CHECK (runtime >= 0);

ALTER TABLE movies add constraint movies_year_check CHECK (year between 1888 and date_part('year', now()));

alter table movies add constraint genres_lenght_check CHECK (array_length(genres, 1) between 1 and 5)
