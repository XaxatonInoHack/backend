create table if not exists "review"
(
    id        serial primary key,
    user_id   bigint       not null,
    review_id bigint       not null,
    feedback  text         not null,
    period    varchar(255) not null
    UNIQUE (user_id, review_id, period)
);

create table if not exists "feedback"
(
    id      serial primary key,
    user_id bigint not null,
    score   varchar(255),
    result  varchar(255),
    resume  text
);

create table if not exists "self_review"
(
    id      serial primary key,
    user_id bigint not null,
    score   varchar(255),
    result  varchar(255),
    resume  text
);