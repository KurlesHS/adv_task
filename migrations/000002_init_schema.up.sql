create table adverts
(
    id               bigserial
        constraint adverts_pk
            primary key,
    title            varchar,
    description      varchar,
    price            double precision,
    create_timestamp timestamp with time zone
);

create unique index adverts_id_uindex
    on adverts (id);

create index adverts_create_timestamp_index
    on adverts (create_timestamp);

create index "adverts_price  _index"
    on adverts (price);

create table photos
(
    adv_id      integer
        constraint photos_adverts_id_fk
            references adverts
            on update cascade on delete cascade,
    photo_order serial,
    link        varchar
);