create table if not exists wallet
(
    id    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    amount integer
);