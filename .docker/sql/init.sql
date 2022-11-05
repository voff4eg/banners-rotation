start transaction;

create table public.slots
(
    id          serial             not null
        constraint "PK_slot_id" primary key,
    description varchar(50)        not null,
    status      smallint default 0 not null
);

comment on table public.slots is 'Слоты';
comment on column public.slots.id is 'ID записи';
comment on column public.slots.description is 'Описание слота';
comment on column public.slots.status is 'Статус слота';

create table public.banners
(
    id          serial             not null
        constraint "PK_banner_id" primary key,
    description varchar(50)        not null,
    status      smallint default 0 not null
);

comment on table public.banners is 'Баннеры';
comment on column public.banners.id is 'ID записи';
comment on column public.banners.description is 'Описание баннера';
comment on column public.banners.status is 'Статус баннера';

create table public.slots_banner
(
    id        serial not null
        constraint "PK_slots_banner_id" primary key,
    slot_id   int    not null
        constraint "FK_slots_banner_slot_id" references slots (id),
    banner_id int    not null
        constraint "FK_slots_banner_banner_id" references banners (id)
);
create unique index "UK_slots_banner"
    on public.slots_banner (slot_id, banner_id);

comment on table public.slots_banner is 'Связь между слотами и баннерами';
comment on column public.slots_banner.id is 'ID записи';
comment on column public.slots_banner.slot_id is 'ID слота';
comment on column public.slots_banner.banner_id is 'ID баннера';

create table public.groups
(
    id     serial             not null
        constraint "PK_group_id" primary key,
        description varchar(50) not null,
    status smallint default 0 not null
);

comment on table public.groups is 'Группы';
comment on column public.groups.id is 'ID записи';
comment on column public.groups.description is 'Описание группы';
comment on column public.groups.status is 'Статус группы';

create table public.stats
(
    id        serial not null
        constraint "PK_stat_id" primary key,
    slot_id   int    not null
        constraint "FK_stats_slot_id" references slots (id),
    banner_id int    not null
        constraint "FK_stats_banner_id" references banners (id),
    group_id  int    not null
        constraint "FK_stats_group_id" references groups (id),
    shows     int    not null default 0,
    hits      int    not null default 0
);

comment on table public.stats is 'Статистика';
comment on column public.stats.id is 'ID записи';
comment on column public.stats.slot_id is 'ID слота';
comment on column public.stats.banner_id is 'ID баннера';
comment on column public.stats.group_id is 'ID группы';
comment on column public.stats.shows is 'Кол-во показов';
comment on column public.stats.hits is 'Кол-во просмотров';

commit;
