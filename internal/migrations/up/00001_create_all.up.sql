create table if not exists "review"
(
    id        serial primary key,
    user_id   bigint       not null,
    review_id bigint       not null,
    feedback  text         not null,
    period    varchar(255) not null
);

create table if not exists "feedback"
(
    id      serial primary key,
    user_id bigint       not null,
    score   varchar(255) not null,
    result  varchar(255) not null,
    resume  text         not null
);