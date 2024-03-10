create table if not exists images (
    id serial primary key,
    user_id uuid references auth.users,
    status int not null default 1,
    prompt text not null,
    image_location text,
    deleted boolean not null default 'false',
    created_at timestamp not null default now(),
    deleted_at timestamp
)