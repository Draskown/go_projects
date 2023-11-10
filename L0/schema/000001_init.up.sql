CREATE TABLE deliveries (
    id serial not null unique,
    name text not null,
    phone text not null,
    zip text not null,
    city text not null,
    address text not null,
    region text not null,
    email text not null
);

CREATE TABLE payments (
    id serial not null unique,
    transaction text not null,
    request_id text,
    currency text not null,
    provider text not null,
    amount int not null,
    payment_dt int not null,
    bank text not null,
    delivery_cost int not null,
    goods_total int not null,
    custom_fee int not null
);

CREATE TABLE items (
    id serial not null unique,
    track_number text not null,
    price int not null,
    rid text not null,
    name text not null,
    sale int not null,
    size text not null,
    total_price int not null,
    nm_id int not null,
    brand text not null,
    status int not null,
    chrt_id int not null
);

CREATE TABLE orders (
    id serial not null unique,
    order_uid text not null,
    track_number text not null,
    entry text not null,
    delivery int references deliveries (id) on delete cascade not null,
    payment int references payments (id) on delete cascade not null,
    items int [] not null,
    locale text not null,
    internal_signature text,
    customer_id text not null,
    delivery_service text not null,
    shardkey text not null,
    sm_id int not null,
    oof_shard text not null,
    date_created text not null
);

CREATE TABLE tests (
    id serial not null unique,
    value int,
    text text,
    arr int[]
);