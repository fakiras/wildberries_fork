-- +goose Up
-- +goose StatementBegin
create schema if not exists app;

create table app.promotion (
    id serial primary key,
    date_from timestamptz not null,
    date_to timestamptz not null,
    created_at timestamptz not null,
    updated_at timestamptz not null
);

create table app.zodiac (
    id serial primary key,
    name text not null
);

create table app.horoscope (
    id serial primary key,
    promotion_id int references app.promotion(id) not null,
    zodiac_id int references app.zodiac(id) not null,
    created_at timestamptz not null,
    text text not null,
    category_id int not null
);

create table app.product (
    id serial primary key,
    nm_id bigint not null,
    category_id int not null,
    category_name text not null,
    name text not null,
    image text,
    price bigint not null constraint price_check check ( price > 0 ),
    created_at timestamptz not null,
    updated_at timestamptz not null,
    deleted_at timestamptz
);

create table app.horoscope_product(
    horoscope_id int references app.horoscope(id) not null ,
    product_id int references app.product(id) not null
);

create table app.auction (
    id serial primary key,
    date_from timestamptz not null,
    date_to timestamptz not null,
    created_at timestamptz not null,
    updated_at timestamptz not null,
    deleted_at timestamptz,
    horoscope_id int references app.horoscope(id) not null,
    position int not null
);

create table app.bet (
    id serial primary key,
    auction_id int references app.auction(id) not null,
    product_id int references app.product(id) not null,
    bet bigint not null constraint bet_check check ( bet > 0 ),
    created_at timestamptz not null,
    deleted_at timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
